package main

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

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
	assert.FileExists(filepath.Join(outdir, "api/fields/Account.json"))
	assert.FileExists(filepath.Join(outdir, "api/fields/Device.json"))

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
