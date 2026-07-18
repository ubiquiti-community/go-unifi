package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

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
