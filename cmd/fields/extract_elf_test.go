package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadELFInstaller(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	body := []byte("fake-elf-content-for-testing")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(body)
		assert.NoError(err)
	}))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	require.NoError(err)

	tmpDir := t.TempDir()
	elfPath, err := downloadELFInstaller(u, tmpDir)
	require.NoError(err)

	assert.Equal(filepath.Join(tmpDir, "unifi-os-installer"), elfPath)

	got, err := os.ReadFile(elfPath)
	require.NoError(err)
	assert.Equal(body, got)
}

func TestDownloadELFInstaller_HTTPError(t *testing.T) {
	require := require.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "boom")
	}))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	require.NoError(err)

	tmpDir := t.TempDir()
	_, err = downloadELFInstaller(u, tmpDir)
	require.Error(err)
	assert.Contains(t, err.Error(), "500")
}

// buildZip creates an in-memory zip with the given entries and returns its bytes.
func buildZip(t *testing.T, entries map[string][]byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for name, content := range entries {
		f, err := w.Create(name)
		require.NoError(t, err)
		_, err = f.Write(content)
		require.NoError(t, err)
	}
	require.NoError(t, w.Close())
	return buf.Bytes()
}

func TestExtractInternalDepsJSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Build a synthetic internal-dependencies.jar in memory.
	internalJar := buildZip(t, map[string][]byte{
		"api/fields/Account.json":            []byte(`{"name":"^[a-zA-Z]+$"}`),
		"api/fields/Device.json":             []byte(`{"name":"^[a-zA-Z]+$"}`),
		"sensitive_metadata.json":            []byte(`{"sensitive_db_fields_by_collection":{"account":["name"]}}`),
		"radio_specification.json":           []byte(`{}`),
		"country_codes_list.json":            []byte(`{}`),
		"geo_ip_country_codes_list.json":     []byte(`{}`),
		"timezones.json":                     []byte(`{}`),
		"legacy_endpoint_segments.json":      []byte(`{}`),
		"event_defs.json":                    []byte(`{}`),
		"ssl-inspection-file-extension.json": []byte(`{}`),
	})

	// Build a synthetic ace.jar containing the internal jar.
	aceJar := buildZip(t, map[string][]byte{
		"BOOT-INF/lib/internal-dependencies.jar": internalJar,
		"BOOT-INF/classes/product.properties":    []byte("version=10.4.57\n"),
	})

	// Write ace.jar to a temp file.
	acePath := filepath.Join(t.TempDir(), "ace.jar")
	require.NoError(os.WriteFile(acePath, aceJar, 0o644))

	outdir := t.TempDir()
	written, err := extractInternalDepsJSON(acePath, outdir)
	require.NoError(err)

	// Should have written 8 top-level files + 2 api/fields files = 10.
	assert.Len(written, 10)

	// Verify a few specific files landed.
	assert.FileExists(filepath.Join(outdir, "sensitive_metadata.json"))
	assert.FileExists(filepath.Join(outdir, "Account.json"))
	assert.FileExists(filepath.Join(outdir, "Device.json"))

	// Verify content.
	got, err := os.ReadFile(filepath.Join(outdir, "sensitive_metadata.json"))
	require.NoError(err)
	assert.JSONEq(`{"sensitive_db_fields_by_collection":{"account":["name"]}}`, string(got))
}

func TestExtractInternalDepsJSON_MissingInternalDeps(t *testing.T) {
	require := require.New(t)
	aceJar := buildZip(t, map[string][]byte{
		"BOOT-INF/classes/product.properties": []byte("version=10.4.57\n"),
	})
	acePath := filepath.Join(t.TempDir(), "ace.jar")
	require.NoError(os.WriteFile(acePath, aceJar, 0o644))

	_, err := extractInternalDepsJSON(acePath, t.TempDir())
	require.Error(err)
	assert.Contains(t, err.Error(), "internal-dependencies")
}

// buildTar builds an in-memory uncompressed tar with the given entries
// (path -> content). Returns the tar bytes.
func buildTar(t *testing.T, entries map[string][]byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	// Sort names for stable output.
	names := make([]string, 0, len(entries))
	for n := range entries {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, name := range names {
		content := entries[name]
		hdr := &tar.Header{
			Name: name,
			Mode: 0o644,
			Size: int64(len(content)),
		}
		require.NoError(t, tw.WriteHeader(hdr))
		_, err := tw.Write(content)
		require.NoError(t, err)
	}
	require.NoError(t, tw.Close())
	return buf.Bytes()
}

// buildGzipTar builds an in-memory gzip-compressed tar with the given entries.
func buildGzipTar(t *testing.T, entries map[string][]byte) []byte {
	t.Helper()
	tarBytes := buildTar(t, entries)
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write(tarBytes)
	require.NoError(t, err)
	require.NoError(t, gw.Close())
	return buf.Bytes()
}

// sha256Hex returns the hex-encoded sha256 of data, prefixed with "sha256:".
func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(h[:])
}

// buildOCITar builds a synthetic OCI image tar (uncompressed outer tar)
// with an index.json, a manifest blob, and the given layers. Each layer
// is a map of tar entries (path -> content); it will be gzip-compressed
// and stored under blobs/sha256/<digest>. Returns the image.tar bytes.
func buildOCITar(t *testing.T, layers []map[string][]byte) []byte {
	t.Helper()

	type blobEntry struct {
		name string
		data []byte
	}
	var blobs []blobEntry

	// Build layer blobs.
	var manifestLayers []v1.Descriptor
	for _, layerEntries := range layers {
		layerTar := buildGzipTar(t, layerEntries)
		digestStr := sha256Hex(layerTar)
		blobName := "blobs/sha256/" + strings.TrimPrefix(digestStr, "sha256:")
		blobs = append(blobs, blobEntry{blobName, layerTar})
		manifestLayers = append(manifestLayers, v1.Descriptor{
			MediaType: "application/vnd.oci.image.layer.v1.tar+gzip",
			Digest:    digest.Digest(digestStr),
			Size:      int64(len(layerTar)),
		})
	}

	// Build manifest blob.
	manifest := v1.Manifest{
		Config: v1.Descriptor{
			MediaType: "application/vnd.oci.image.config.v1+json",
			Digest:    digest.Digest("sha256:" + strings.Repeat("0", 64)),
			Size:      2,
		},
		Layers: manifestLayers,
	}
	manifest.SchemaVersion = 2
	manifestBytes, err := json.Marshal(manifest)
	require.NoError(t, err)
	manifestDigest := sha256Hex(manifestBytes)
	blobs = append(blobs, blobEntry{
		"blobs/sha256/" + strings.TrimPrefix(manifestDigest, "sha256:"),
		manifestBytes,
	})

	// Build index.
	idx := v1.Index{
		Manifests: []v1.Descriptor{
			{
				MediaType: "application/vnd.oci.image.manifest.v1+json",
				Digest:    digest.Digest(manifestDigest),
				Size:      int64(len(manifestBytes)),
			},
		},
	}
	idx.SchemaVersion = 2
	idxBytes, err := json.Marshal(idx)
	require.NoError(t, err)

	// Build the outer image.tar.
	entries := map[string][]byte{
		"index.json": idxBytes,
		"oci-layout": []byte(`{"imageLayoutVersion":"1.0.0"}`),
	}
	for _, b := range blobs {
		entries[b.name] = b.data
	}
	return buildTar(t, entries)
}

func TestExtractACFromELF(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Build a synthetic ace.jar with product.properties.
	aceJar := buildZip(t, map[string][]byte{
		"BOOT-INF/classes/product.properties":    []byte("version=10.4.57\nbuild=abc\n"),
		"BOOT-INF/lib/internal-dependencies.jar": []byte("placeholder"),
	})

	// Build an OCI image tar with one layer containing ace.jar at the
	// expected path.
	imageTar := buildOCITar(t, []map[string][]byte{
		{aceJarLayerPath: aceJar},
	})

	// Wrap image.tar in a zip (this is what the installer appends).
	installer := buildZip(t, map[string][]byte{
		"image.tar": imageTar,
	})

	installerPath := filepath.Join(t.TempDir(), "installer")
	require.NoError(os.WriteFile(installerPath, installer, 0o644))

	outdir := t.TempDir()
	acePath, version, err := extractACFromELF(installerPath, outdir)
	require.NoError(err)

	assert.Equal(filepath.Join(outdir, "ace.jar"), acePath)
	assert.Equal("10.4.57", version)

	// Verify ace.jar was actually written.
	got, err := os.ReadFile(acePath)
	require.NoError(err)
	assert.Equal(aceJar, got)
}

func TestExtractACFromELF_NoIndexJSON(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	// Build a non-OCI image tar (Docker-save style: manifest.json, no index.json).
	imageTar := buildTar(t, map[string][]byte{
		"manifest.json": []byte(`[]`),
	})
	installer := buildZip(t, map[string][]byte{
		"image.tar": imageTar,
	})
	installerPath := filepath.Join(t.TempDir(), "installer")
	require.NoError(os.WriteFile(installerPath, installer, 0o644))

	_, _, err := extractACFromELF(installerPath, t.TempDir())
	require.Error(err)
	assert.Contains(err.Error(), "index.json")
}

func TestExtractACFromELF_AceJarNotFound(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	// Layer without ace.jar.
	imageTar := buildOCITar(t, []map[string][]byte{
		{"some/other/file.txt": []byte("no ace.jar here")},
	})
	installer := buildZip(t, map[string][]byte{
		"image.tar": imageTar,
	})
	installerPath := filepath.Join(t.TempDir(), "installer")
	require.NoError(os.WriteFile(installerPath, installer, 0o644))

	_, _, err := extractACFromELF(installerPath, t.TempDir())
	require.Error(err)
	assert.Contains(err.Error(), "ace.jar not found")
}
