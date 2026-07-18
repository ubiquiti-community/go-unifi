package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// keepToplevel is the list of top-level JSON files inside
// internal-dependencies.jar that we extract alongside api/fields.
var keepToplevel = []string{
	"legacy_endpoint_segments.json",
	"event_defs.json",
	"sensitive_metadata.json",
	"radio_specification.json",
	"country_codes_list.json",
	"geo_ip_country_codes_list.json",
	"timezones.json",
	"ssl-inspection-file-extension.json",
}

// keepDirs is the list of directories inside internal-dependencies.jar
// whose .json files we extract recursively.
var keepDirs = []string{"api/fields"}

// aceJarPath inside the OCI image layer tar.
const aceJarLayerPath = "usr/lib/unifi/lib/ace.jar"

// downloadELFInstaller streams the Unifi OS installer (a self-extracting
// ELF with an appended zip) to a temp file and returns its path. The
// installer is ~880 MB; we never hold it in memory.
func downloadELFInstaller(u *url.URL, outdir string) (string, error) {
	dst := filepath.Join(outdir, "unifi-os-installer")
	f, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("create installer temp file: %w", err)
	}
	defer f.Close()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download installer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download installer: HTTP %d", resp.StatusCode)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", fmt.Errorf("write installer: %w", err)
	}

	return dst, nil
}

// Compile-time assertions that the OCI types we depend on exist.
var (
	_ v1.Index
	_ v1.Manifest
	_ digest.Digest
)

// extractInternalDepsJSON opens ace.jar, reads
// BOOT-INF/lib/internal-dependencies.jar, and writes the JSON definition
// files we care about (keepToplevel + keepDirs recursively) to outdir.
// Returns the list of relative paths written.
func extractInternalDepsJSON(aceJarPath, outdir string) ([]string, error) {
	aceZip, err := zip.OpenReader(aceJarPath)
	if err != nil {
		return nil, fmt.Errorf("open ace.jar: %w", err)
	}
	defer aceZip.Close()

	const internalDepsName = "BOOT-INF/lib/internal-dependencies.jar"
	var internalEntry io.ReadCloser
	internalEntry, err = aceZip.Open(internalDepsName)
	if err != nil {
		// Fall back: any entry whose name contains "internal-dependencies".
		var found *zip.File
		for _, f := range aceZip.File {
			if strings.Contains(f.Name, "internal-dependencies") {
				found = f
				break
			}
		}
		if found == nil {
			return nil, fmt.Errorf("internal-dependencies.jar not found in ace.jar")
		}
		internalEntry, err = found.Open()
		if err != nil {
			return nil, fmt.Errorf("open internal-dependencies.jar: %w", err)
		}
	}
	defer internalEntry.Close()

	internalBytes, err := io.ReadAll(internalEntry)
	if err != nil {
		return nil, fmt.Errorf("read internal-dependencies.jar: %w", err)
	}

	internalReader, err := zip.NewReader(bytes.NewReader(internalBytes), int64(len(internalBytes)))
	if err != nil {
		return nil, fmt.Errorf("open internal-dependencies.jar as zip: %w", err)
	}

	var written []string
	keepSet := make(map[string]struct{}, len(keepToplevel))
	for _, k := range keepToplevel {
		keepSet[k] = struct{}{}
	}

	for _, f := range internalReader.File {
		name := f.Name
		if _, ok := keepSet[name]; ok {
			if err := writeZipEntry(f, outdir, name); err != nil {
				return nil, err
			}
			written = append(written, name)
			continue
		}
		for _, dir := range keepDirs {
			if strings.HasPrefix(name, dir+"/") && strings.HasSuffix(name, ".json") {
				// Flatten to top level, matching the .deb path behavior
				// (extract.go uses filepath.Base for the same purpose).
				flatName := filepath.Base(name)
				if err := writeZipEntry(f, outdir, flatName); err != nil {
					return nil, err
				}
				written = append(written, flatName)
				break
			}
		}
	}

	return written, nil
}

// writeZipEntry writes a zip entry to outdir, preserving its relative path.
func writeZipEntry(f *zip.File, outdir, name string) error {
	src, err := f.Open()
	if err != nil {
		return fmt.Errorf("open zip entry %q: %w", name, err)
	}
	defer src.Close()

	dst := filepath.Join(outdir, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir for %q: %w", name, err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create %q: %w", name, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return fmt.Errorf("write %q: %w", name, err)
	}
	return nil
}

// extractACFromELF opens the self-extracting installer, reads image.tar,
// parses the OCI index + manifest, scans layers for ace.jar, writes it to
// outdir/ace.jar, and reads the Network version from product.properties.
// Returns (aceJarPath, networkVersion, error).
func extractACFromELF(installerPath, outdir string) (string, string, error) {
	zipReader, err := zip.OpenReader(installerPath)
	if err != nil {
		return "", "", fmt.Errorf("open installer as zip: %w", err)
	}
	defer zipReader.Close()

	// Read image.tar to a temp file (it can be ~840 MB).
	var imageTarEntry *zip.File
	for _, f := range zipReader.File {
		if f.Name == "image.tar" {
			imageTarEntry = f
			break
		}
	}
	if imageTarEntry == nil {
		return "", "", fmt.Errorf("image.tar not found in installer zip")
	}

	imageTarPath := filepath.Join(outdir, "image.tar")
	if err := copyZipEntryToFile(imageTarEntry, imageTarPath); err != nil {
		return "", "", fmt.Errorf("extract image.tar: %w", err)
	}
	defer os.Remove(imageTarPath)

	aceJarPath := filepath.Join(outdir, "ace.jar")
	networkVersion, err := extractAceFromImageTar(imageTarPath, aceJarPath)
	if err != nil {
		return "", "", err
	}

	return aceJarPath, networkVersion, nil
}

// extractAceFromImageTar parses an OCI image tar, locates ace.jar in one
// of the layers, writes it to aceJarPath, and reads the Network version
// from product.properties inside ace.jar.
func extractAceFromImageTar(imageTarPath, aceJarPath string) (string, error) {
	imageFile, err := os.Open(imageTarPath)
	if err != nil {
		return "", fmt.Errorf("open image.tar: %w", err)
	}
	defer imageFile.Close()

	imageTar := tar.NewReader(imageFile)

	// First pass: find and read index.json. tar is a streaming format,
	// so we break as soon as we find it.
	var indexJSON []byte

	for {
		hdr, err := imageTar.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("read image.tar: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		if hdr.Name == "index.json" {
			indexJSON, err = io.ReadAll(imageTar)
			if err != nil {
				return "", fmt.Errorf("read index.json: %w", err)
			}
			break
		}
	}

	if indexJSON == nil {
		return "", fmt.Errorf("image.tar is not an OCI image layout (missing index.json)")
	}

	var idx v1.Index
	if err := json.Unmarshal(indexJSON, &idx); err != nil {
		return "", fmt.Errorf("parse index.json: %w", err)
	}
	if len(idx.Manifests) == 0 {
		return "", fmt.Errorf("no manifests in index.json")
	}

	// Re-open the image tar for a second pass to read blobs by name.
	// (tar is a streaming format; we can't seek within a single reader.)
	manifestDigest := idx.Manifests[0].Digest
	manifestBlobName := "blobs/sha256/" + strings.TrimPrefix(string(manifestDigest), "sha256:")

	manifestBytes, err := readTarEntry(imageTarPath, manifestBlobName)
	if err != nil {
		return "", fmt.Errorf("read manifest blob: %w", err)
	}

	var manifest v1.Manifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return "", fmt.Errorf("parse manifest: %w", err)
	}

	// Scan each layer for ace.jar.
	for _, layer := range manifest.Layers {
		layerBlobName := "blobs/sha256/" + strings.TrimPrefix(string(layer.Digest), "sha256:")
		found, err := extractAceFromLayer(imageTarPath, layerBlobName, layer.MediaType, aceJarPath)
		if err != nil {
			return "", fmt.Errorf("scan layer %s: %w", layerBlobName[:min(len(layerBlobName), 30)], err)
		}
		if found {
			version, err := readNetworkVersion(aceJarPath)
			if err != nil {
				return "", err
			}
			return version, nil
		}
	}

	return "", fmt.Errorf("ace.jar not found in any image layer")
}

// readTarEntry opens imageTarPath as a tar and returns the bytes of the
// first regular-file entry matching name.
func readTarEntry(imageTarPath, name string) ([]byte, error) {
	f, err := os.Open(imageTarPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("entry %q not found", name)
		}
		if err != nil {
			return nil, err
		}
		if hdr.Typeflag == tar.TypeReg && hdr.Name == name {
			return io.ReadAll(tr)
		}
	}
}

// extractAceFromLayer reads a layer blob from imageTarPath, decompresses
// it if mediaType indicates gzip, scans for aceJarLayerPath, and writes
// the match to aceJarPath. Returns (true, nil) if found.
func extractAceFromLayer(imageTarPath, blobName, mediaType string, aceJarPath string) (bool, error) {
	layerBytes, err := readTarEntry(imageTarPath, blobName)
	if err != nil {
		return false, err
	}

	var layerReader io.Reader = bytes.NewReader(layerBytes)
	if strings.HasSuffix(mediaType, "+gzip") || strings.HasSuffix(mediaType, ".gzip") {
		gr, err := gzip.NewReader(layerReader)
		if err != nil {
			return false, fmt.Errorf("decompress layer: %w", err)
		}
		defer gr.Close()
		layerReader = gr
	}

	layerTar := tar.NewReader(layerReader)
	for {
		hdr, err := layerTar.Next()
		if err == io.EOF {
			return false, nil
		}
		if err != nil {
			return false, fmt.Errorf("read layer tar: %w", err)
		}
		if hdr.Typeflag == tar.TypeReg && hdr.Name == aceJarLayerPath {
			out, err := os.Create(aceJarPath)
			if err != nil {
				return false, fmt.Errorf("create ace.jar: %w", err)
			}
			defer out.Close()
			if _, err := io.Copy(out, layerTar); err != nil {
				return false, fmt.Errorf("write ace.jar: %w", err)
			}
			return true, nil
		}
	}
}

// readNetworkVersion opens ace.jar and reads the version= line from
// BOOT-INF/classes/product.properties. Returns "unknown" if not found.
func readNetworkVersion(aceJarPath string) (string, error) {
	zipReader, err := zip.OpenReader(aceJarPath)
	if err != nil {
		return "", fmt.Errorf("open ace.jar for version: %w", err)
	}
	defer zipReader.Close()

	propsEntry, err := zipReader.Open("BOOT-INF/classes/product.properties")
	if err != nil {
		return "unknown", nil
	}
	defer propsEntry.Close()

	propsBytes, err := io.ReadAll(propsEntry)
	if err != nil {
		return "unknown", nil
	}

	re := regexp.MustCompile(`(?m)^version=(.+)$`)
	m := re.FindSubmatch(propsBytes)
	if m == nil {
		return "unknown", nil
	}
	return strings.TrimSpace(string(m[1])), nil
}

// copyZipEntryToFile writes a zip.File to a file on disk.
func copyZipEntryToFile(f *zip.File, dst string) error {
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
