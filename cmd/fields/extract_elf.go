package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
				if err := writeZipEntry(f, outdir, name); err != nil {
					return nil, err
				}
				written = append(written, name)
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
