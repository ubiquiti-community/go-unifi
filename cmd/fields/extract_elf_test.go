package main

import (
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
