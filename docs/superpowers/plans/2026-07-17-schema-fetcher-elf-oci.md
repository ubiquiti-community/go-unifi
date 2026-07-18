# Unifi Schema Fetcher: ELF/OCI Support, Sensitive Metadata, and Full Automation — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Port the Python PoC's ELF extraction path to Go so the generator can fetch Unifi Network 10.x schemas, wire `sensitive_metadata.json` into the spec generator, and fully automate the release pipeline via GitHub Actions.

**Architecture:** A new `extract_elf.go` in `cmd/fields/` downloads the Unifi OS installer (ELF), parses the OCI image tar inside it, locates `ace.jar`, and extracts field validators + metadata files from `internal-dependencies.jar`. `main.go` routes to this path via `-latest`; the `.deb` path stays for `-version X.Y.Z` where `X < 10`. `schema.go` gains a pre-pass that marks `FieldInfo.Sensitive=true` on fields listed in `sensitive_metadata.json`, and the attribute generators propagate that to `specification.json`. Three GitHub Actions workflows handle the daily-cron → PR → auto-merge → tag → release pipeline.

**Tech Stack:** Go 1.25, `github.com/opencontainers/image-spec` + `github.com/opencontainers/go-digest` (OCI parsing), stdlib `archive/tar`/`archive/zip`/`compress/gzip`, `github.com/hashicorp/terraform-plugin-codegen-spec`, GitHub Actions (`peter-evans/create-pull-request`, `gh` CLI).

**Spec:** `docs/superpowers/specs/2026-07-17-schema-fetcher-elf-oci-design.md`

---

## File Structure

| File | Status | Responsibility |
|---|---|---|
| `cmd/fields/extract_elf.go` | NEW | `downloadELFInstaller`, `extractACFromELF`, `extractInternalDepsJSON` — the ELF/OCI pipeline. |
| `cmd/fields/extract_elf_test.go` | NEW | Tests for the three functions + programmatic fixture builders (synthetic zips/tars). |
| `cmd/fields/fwupdate.go` | EDIT | Add `unifiOSServerProduct`, `linuxX64Platform` constants. |
| `cmd/fields/version.go` | EDIT | Add `latestUnifiOSVersion()`. |
| `cmd/fields/version_test.go` | EDIT | Add `TestLatestUnifiOSVersion`. |
| `cmd/fields/main.go` | EDIT | Route by flag (not version number); call ELF path for `-latest`; load sensitive metadata; call `MarkSensitiveFields`. |
| `cmd/fields/schema.go` | EDIT | `SensitiveMetadata` struct + loader; `CollectionName()` method; `MarkSensitiveFields`/`walkAndMark`; `FieldInfo.Sensitive`; set `Sensitive: ptr(true)` on leaf attributes. |
| `cmd/fields/schema_test.go` | EDIT | `TestMarkSensitiveFields`, `TestFieldToResourceAttribute_Sensitive`. |
| `go.mod`, `go.sum` | EDIT | Add `image-spec` + `go-digest`. |
| `.github/workflows/generate.yaml` | EDIT | Add test gate, auto-merge, permissions. |
| `.github/workflows/tag.yaml` | NEW | Tag `v<unifi-version>` on push to `main` if absent. |

---

### Task 1: Verify the Unifi image.tar is OCI layout (format spike)

**Why:** The spec committed to manifest-driven OCI parsing with a hard error if `index.json` is absent. The PoC author never parsed the manifest — they brute-force scanned. Before writing 250 lines of parser, confirm the format. If `index.json` is absent, STOP and revisit the spec with the user.

**Files:**
- No file changes — this is a read-only investigation.

- [ ] **Step 1: Download the ELF installer to a temp location**

The file is ~880 MB. Download to `/tmp` (or the OS temp dir). Use the URL from the spec:

```bash
curl -L -o /tmp/unifi-os-installer \
  "https://fw-download.ubnt.com/data/unifi-os-server/f5e2-linux-x64-5.1.21-a400c9c6-8328-4634-b223-ebfcf742720a.21-x64"
```

Expected: a ~880 MB file at `/tmp/unifi-os-installer`.

- [ ] **Step 2: Open the ELF as a zip and list top-level entries**

```bash
python3 -c "
import zipfile
zf = zipfile.ZipFile('/tmp/unifi-os-installer')
for n in zf.namelist():
    print(n)
"
```

Expected: `image.tar` is among the entries. Record the full namelist.

- [ ] **Step 3: Extract `image.tar` and list its top-level tar members**

```bash
python3 -c "
import zipfile, tarfile, io
zf = zipfile.ZipFile('/tmp/unifi-os-installer')
img = zf.read('image.tar')
tf = tarfile.open(fileobj=io.BytesIO(img))
for m in tf.getmembers():
    print(m.name, m.size, 'dir' if m.isdir() else 'file')
" 2>&1 | head -40
```

Expected: if OCI layout, you'll see `index.json`, `oci-layout`, and `blobs/sha256/` directory entries. If Docker-save format, you'll see `manifest.json` and `repositories` (no `index.json`). Record what you find.

- [ ] **Step 4: If OCI layout, read `index.json` and print the manifest digests**

```bash
python3 -c "
import zipfile, tarfile, io, json
zf = zipfile.ZipFile('/tmp/unifi-os-installer')
img = zf.read('image.tar')
tf = tarfile.open(fileobj=io.BytesIO(img))
idx = tf.extractfile('index.json').read()
print(json.dumps(json.loads(idx), indent=2)[:2000])
"
```

Expected: JSON with `schemaVersion` and `manifests` array, each entry having a `digest` field like `sha256:abc123...`. Record one manifest digest.

- [ ] **Step 5: Read that manifest and print its layer digests + mediaTypes**

```bash
python3 -c "
import zipfile, tarfile, io, json
zf = zipfile.ZipFile('/tmp/unifi-os-installer')
img = zf.read('image.tar')
tf = tarfile.open(fileobj=io.BytesIO(img))
idx = json.loads(tf.extractfile('index.json').read())
digest = idx['manifests'][0]['digest'].split(':',1)[1]
man = json.loads(tf.extractfile(f'blobs/sha256/{digest}').read())
print(json.dumps(man, indent=2)[:3000])
"
```

Expected: JSON with `layers` array, each having `mediaType` (e.g. `application/vnd.oci.image.layer.v1.tar+gzip`) and `digest`. Record whether layers are gzip-compressed (mediaType ends in `+gzip`) or not.

- [ ] **Step 6: Decision point**

- If `index.json` is present and manifest parsing worked → **proceed to Task 2**. The spec is correct.
- If `index.json` is absent (Docker-save format or other) → **STOP**. Do not write the parser. Report back to the user that the image.tar is not OCI layout, and the spec needs revision (likely back to brute-force scan, which the user originally rejected). Wait for guidance.

- [ ] **Step 7: Clean up the downloaded installer**

```bash
rm -f /tmp/unifi-os-installer
```

---

### Task 2: Add OCI dependencies to go.mod

**Files:**
- Modify: `go.mod`, `go.sum`

- [ ] **Step 1: Add the two dependencies**

Run from the repo root:

```bash
go get github.com/opencontainers/image-spec@latest
go get github.com/opencontainers/go-digest@latest
```

- [ ] **Step 2: Verify they resolve**

```bash
go build ./...
```

Expected: builds with no errors (nothing imports them yet, but the deps are now in go.mod).

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "deps: add opencontainers image-spec and go-digest for OCI parsing"
```

---

### Task 3: Implement `downloadELFInstaller` with TDD

**Files:**
- Create: `cmd/fields/extract_elf.go`
- Create: `cmd/fields/extract_elf_test.go`

This task creates the new file with the first function plus the `KEEP_TOPLEVEL` and `KEEP_DIRS` constants. Subsequent tasks add the other two functions.

- [ ] **Step 1: Write the failing test**

Create `cmd/fields/extract_elf_test.go`:

```go
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
```

- [ ] **Step 2: Run the test to verify it fails**

```bash
go test ./cmd/fields/ -run TestDownloadELFInstaller -v
```

Expected: FAIL with `downloadELFInstaller undefined`.

- [ ] **Step 3: Write the minimal implementation**

Create `cmd/fields/extract_elf.go`:

```go
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
```

- [ ] **Step 4: Run the test to verify it passes**

```bash
go test ./cmd/fields/ -run TestDownloadELFInstaller -v
```

Expected: PASS (both tests).

- [ ] **Step 5: Commit**

```bash
git add cmd/fields/extract_elf.go cmd/fields/extract_elf_test.go
git commit -m "fields: add downloadELFInstaller for Unifi OS ELF download"
```

---

### Task 4: Implement `extractInternalDepsJSON` with TDD

This is simpler than `extractACFromELF` (no OCI parsing), so we do it first. It also lets us build the test fixture helpers we'll reuse.

**Files:**
- Modify: `cmd/fields/extract_elf.go`
- Modify: `cmd/fields/extract_elf_test.go`

- [ ] **Step 1: Write the failing test**

Append to `cmd/fields/extract_elf_test.go`:

```go
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
```

Also add the `buildZip` helper (used by this and future tests) to the test file:

```go
import (
	"archive/zip"
	"bytes"
)

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
```

- [ ] **Step 2: Run the test to verify it fails**

```bash
go test ./cmd/fields/ -run TestExtractInternalDepsJSON -v
```

Expected: FAIL with `extractInternalDepsJSON undefined`.

- [ ] **Step 3: Write the implementation**

Add to `cmd/fields/extract_elf.go`:

```go
import (
	"archive/zip"
	"path"
	"path/filepath"
	"strings"
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
	internalEntry, err := aceZip.Open(internalDepsName)
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
```

Note: the imports block now also needs `"archive/zip"`, `"bytes"`, `"path"`, `"path/filepath"`, `"strings"`. Merge these with the existing imports. Remove the unused `path` import if `goimports` complains — `path` isn't used here, only `path/filepath`. Run `goimports -w` to fix.

- [ ] **Step 4: Run the test to verify it passes**

```bash
go test ./cmd/fields/ -run TestExtractInternalDepsJSON -v
```

Expected: PASS (both tests).

- [ ] **Step 5: Commit**

```bash
git add cmd/fields/extract_elf.go cmd/fields/extract_elf_test.go
git commit -m "fields: add extractInternalDepsJSON for metadata + api/fields extraction"
```

---

### Task 5: Implement `extractACFromELF` with TDD

This is the trickiest function — OCI manifest parsing + layer scanning + ace.jar extraction + version reading. Build it in stages with synthetic fixtures.

**Files:**
- Modify: `cmd/fields/extract_elf.go`
- Modify: `cmd/fields/extract_elf_test.go`

- [ ] **Step 1: Add a tar-builder helper to the test file**

Append to `cmd/fields/extract_elf_test.go`:

```go
import (
	"archive/tar"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// buildTar builds an in-memory uncompressed tar with the given entries
// (path → content). Returns the tar bytes.
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
// is a map of tar entries (path → content); it will be gzip-compressed
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
			Digest:    "sha256:" + strings.Repeat("0", 64),
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
				Digest:   digest.Digest(manifestDigest),
				Size:     int64(len(manifestBytes)),
			},
		},
	}
	idx.SchemaVersion = 2
	idxBytes, err := json.Marshal(idx)
	require.NoError(t, err)

	// Build the outer image.tar.
	entries := map[string][]byte{
		"index.json":  idxBytes,
		"oci-layout":   []byte(`{"imageLayoutVersion":"1.0.0"}`),
	}
	for _, b := range blobs {
		entries[b.name] = b.data
	}
	return buildTar(t, entries)
}
```

Note: add `"compress/gzip"`, `"crypto/sha256"`, `"encoding/hex"`, `"encoding/json"`, `"sort"`, `"strings"` to the test file imports. Also add `v1 "github.com/opencontainers/image-spec/specs-go/v1"` and `"github.com/opencontainers/go-digest"`.

- [ ] **Step 2: Write the failing test for ace.jar extraction**

Append to `cmd/fields/extract_elf_test.go`:

```go
func TestExtractACFromELF(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Build a synthetic ace.jar with product.properties.
	aceJar := buildZip(t, map[string][]byte{
		"BOOT-INF/classes/product.properties": []byte("version=10.4.57\nbuild=abc\n"),
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
```

- [ ] **Step 3: Run the test to verify it fails**

```bash
go test ./cmd/fields/ -run TestExtractACFromELF -v
```

Expected: FAIL with `extractACFromELF undefined`.

- [ ] **Step 4: Write the implementation**

Add to `cmd/fields/extract_elf.go`:

```go
import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"regexp"
)

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

Notes on imports: the file now needs `"archive/tar"`, `"archive/zip"`, `"bytes"`, `"compress/gzip"`, `"encoding/json"`, `"regexp"`, `"strings"` in addition to the existing imports. Run `goimports -w cmd/fields/extract_elf.go` to fix imports.

The two-pass approach (first pass finds `index.json`, second pass via `readTarEntry` reads blobs by name) is necessary because tar is a streaming format — we can't seek within a single reader. `readTarEntry` re-opens the file for each blob read. This is acceptable because there are only a handful of entries (index.json + manifest + a few layer blobs).

- [ ] **Step 5: Run the test to verify it passes**

```bash
go test ./cmd/fields/ -run TestExtractACFromELF -v
```

Expected: PASS (all three tests).

If the `NoIndexJSON` test fails because `zip.OpenReader` succeeds on a zip that has no `index.json`... that's fine, the error should come from `extractAceFromImageTar` returning the "missing index.json" error. Verify the error contains "index.json".

- [ ] **Step 6: Commit**

```bash
git add cmd/fields/extract_elf.go cmd/fields/extract_elf_test.go
git commit -m "fields: add extractACFromELF with OCI manifest parsing"
```

---

### Task 6: Add `latestUnifiOSVersion` to firmware API

**Files:**
- Modify: `cmd/fields/fwupdate.go`
- Modify: `cmd/fields/version.go`
- Modify: `cmd/fields/version_test.go`

- [ ] **Step 1: Write the failing test**

Append to `cmd/fields/version_test.go`:

```go
func TestLatestUnifiOSVersion(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	osVersion, err := version.NewVersion("v5.1.21")
	require.NoError(err)

	osDownload, err := url.Parse(
		"https://fw-download.ubnt.com/data/unifi-os-server/f5e2-linux-x64-5.1.21-a400c9c6-8328-4634-b223-ebfcf742720a.21-x64",
	)
	require.NoError(err)

	respData := firmwareUpdateApiResponse{
		Embedded: firmwareUpdateApiResponseEmbedded{
			Firmware: []firmwareUpdateApiResponseEmbeddedFirmware{
				{
					Channel:  releaseChannel,
					Platform: linuxX64Platform,
					Product:  unifiOSServerProduct,
					Version:  osVersion,
					Links: firmwareUpdateApiResponseEmbeddedFirmwareLinks{
						Data: firmwareUpdateApiResponseEmbeddedFirmwareDataLink{
							Href: osDownload,
						},
					},
				},
				{
					Channel:  releaseChannel,
					Platform: "macOS-dmg-amd64",
					Product:  unifiOSServerProduct,
					Version:  osVersion,
					Links: firmwareUpdateApiResponseEmbeddedFirmwareLinks{
						Data: firmwareUpdateApiResponseEmbeddedFirmwareDataLink{Href: nil},
					},
				},
				{
					Channel:  releaseChannel,
					Platform: linuxX64Platform,
					Product:  "unifi-controller", // wrong product
					Version:  osVersion,
					Links: firmwareUpdateApiResponseEmbeddedFirmwareLinks{
						Data: firmwareUpdateApiResponseEmbeddedFirmwareDataLink{Href: nil},
					},
				},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		assert.Contains(query["filter"], firmwareUpdateApiFilter("eq", "channel", releaseChannel))
		assert.Contains(query["filter"], firmwareUpdateApiFilter("eq", "product", unifiOSServerProduct))
		assert.Contains(query["filter"], firmwareUpdateApiFilter("eq", "platform", linuxX64Platform))

		resp, err := json.Marshal(respData)
		assert.NoError(err)
		_, err = rw.Write(resp)
		assert.NoError(err)
	}))
	defer srv.Close()

	firmwareUpdateApi = srv.URL
	gotVersion, gotDownload, err := latestUnifiOSVersion()
	require.NoError(err)

	assert.Equal(osVersion.Core(), gotVersion)
	assert.Equal(osDownload, gotDownload)
}
```

- [ ] **Step 2: Run the test to verify it fails**

```bash
go test ./cmd/fields/ -run TestLatestUnifiOSVersion -v
```

Expected: FAIL with `latestUnifiOSVersion undefined` and/or `linuxX64Platform undefined` and/or `unifiOSServerProduct undefined`.

- [ ] **Step 3: Add the constants to fwupdate.go**

In `cmd/fields/fwupdate.go`, add to the existing `const` block:

```go
const (
	debianPlatform         = "debian"
	releaseChannel         = "release"
	unifiControllerProduct = "unifi-controller"
	unifiOSServerProduct   = "unifi-os-server"
	linuxX64Platform       = "linux-x64"
	maxVersion             = "10.0.0"
)
```

(Replace the existing const block with this expanded one.)

- [ ] **Step 4: Add `latestUnifiOSVersion` to version.go**

Append to `cmd/fields/version.go`:

```go
// latestUnifiOSVersion queries the firmware update API for the latest
// unifi-os-server release on linux-x64. Unlike latestUnifiVersion, it
// has no maxVersion cap (10.x is always reached via this path). Returns
// the OS version (e.g. v5.1.21) and the download URL.
func latestUnifiOSVersion() (*version.Version, *url.URL, error) {
	u, err := url.Parse(firmwareUpdateApi)
	if err != nil {
		return nil, nil, err
	}

	q := u.Query()
	q.Add("filter", firmwareUpdateApiFilter("eq", "channel", releaseChannel))
	q.Add("filter", firmwareUpdateApiFilter("eq", "product", unifiOSServerProduct))
	q.Add("filter", firmwareUpdateApiFilter("eq", "platform", linuxX64Platform))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var respData firmwareUpdateApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, nil, err
	}

	for _, fw := range respData.Embedded.Firmware {
		if fw.Platform != linuxX64Platform {
			continue
		}
		return fw.Version.Core(), fw.Links.Data.Href, nil
	}

	return nil, nil, fmt.Errorf("no unifi-os-server firmware found for platform %q", linuxX64Platform)
}
```

Add `"fmt"` to the imports in version.go if not already present.

- [ ] **Step 5: Run the test to verify it passes**

```bash
go test ./cmd/fields/ -run TestLatestUnifiOSVersion -v
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add cmd/fields/fwupdate.go cmd/fields/version.go cmd/fields/version_test.go
git commit -m "fields: add latestUnifiOSVersion for unifi-os-server/linux-x64 firmware API"
```

---

### Task 7: Route `main.go` to the ELF path for `-latest`

**Files:**
- Modify: `cmd/fields/main.go`

This is the integration point. We change the version/URL resolution in `main()` to branch by flag, and wire the ELF extraction functions into the existing download→extract flow.

- [ ] **Step 1: Read the current main() resolution block**

Read `cmd/fields/main.go` lines 304-420 (the flag parsing + version resolution + download/extraction block). Understand the current flow:

1. Parse flags (`-output-dir`, `-download-only`, `-latest`, `-generate-spec`, `-spec-output`).
2. Resolve version: `-latest` → `latestUnifiVersion()`; else parse `-version X.Y.Z` + construct `.deb` URL.
3. Compute `fieldsDir = cmd/fields/v<version>/`.
4. If `fieldsDir` doesn't exist: `downloadJar` → `extractJSON` → `copyCustom`.
5. If `-download-only`: exit.
6. Read `fieldsDir`, iterate `*.json`, generate code.

We need to change step 2-4 to branch by flag.

- [ ] **Step 2: Edit the version resolution block**

In `cmd/fields/main.go`, find the block starting after `flag.Parse()` (around line 329) through the end of the `if fieldsInfo, err := os.Stat(fieldsDir)` block (around line 417). Replace it.

Old code (lines ~330-417, approximately):

```go
	specifiedVersion := flag.Arg(0)
	if specifiedVersion != "" && *useLatestVersion {
		fmt.Print("error: cannot specify version with latest\n\n")
		usage()
		os.Exit(1)
	} else if specifiedVersion == "" && !*useLatestVersion {
		fmt.Print("error: must specify version or latest\n\n")
		usage()
		os.Exit(1)
	}

	var unifiVersion *version.Version
	var unifiDownloadUrl *url.URL
	var err error

	if *useLatestVersion {
		unifiVersion, unifiDownloadUrl, err = latestUnifiVersion()
		if err != nil {
			panic(err)
		}
	} else {
		unifiVersion, err = version.NewVersion(specifiedVersion)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		unifiDownloadUrl, err = url.Parse(fmt.Sprintf("https://dl.ui.com/unifi/%s/unifi_sysvinit_all.deb", unifiVersion))
		if err != nil {
			panic(err)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get the current filename")
	}

	versionBaseDir := filepath.Dir(filename)

	fieldsDir := filepath.Join(versionBaseDir, fmt.Sprintf("v%s", unifiVersion))

	outDir := filepath.Join(wd, *outputDirFlag)

	fieldsInfo, err := os.Stat(fieldsDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}

		err = os.MkdirAll(fieldsDir, 0o755)
		if err != nil {
			panic(err)
		}

		// download fields, create
		jarFile, err := downloadJar(unifiDownloadUrl, fieldsDir)
		if err != nil {
			panic(err)
		}

		err = extractJSON(jarFile, fieldsDir)
		if err != nil {
			panic(err)
		}

		// defer func() {
		// 	err = os.RemoveAll(fieldsDir)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }()

		err = copyCustom(fieldsDir)
		if err != nil {
			panic(err)
		}

		fieldsInfo, err = os.Stat(fieldsDir)
		if err != nil {
			panic(err)
		}
	}
```

New code:

```go
	specifiedVersion := flag.Arg(0)
	if specifiedVersion != "" && *useLatestVersion {
		fmt.Print("error: cannot specify version with latest\n\n")
		usage()
		os.Exit(1)
	} else if specifiedVersion == "" && !*useLatestVersion {
		fmt.Print("error: must specify version or latest\n\n")
		usage()
		os.Exit(1)
	}

	// Determine which fetcher path to use.
	//
	// -latest  → ELF path (queries unifi-os-server/linux-x64). The OS
	//             version (e.g. v5.1.21) is only used to satisfy the
	//             firmware API; the Network version (e.g. 10.4.57) is
	//             discovered during extraction from product.properties.
	// -version X.Y.Z where X < 10 → .deb path (legacy, unchanged).
	// -version X.Y.Z where X >= 10 → rejected; use -latest for 10.x+.
	var useELFPath bool
	var unifiVersion *version.Version
	var unifiDownloadUrl *url.URL
	var err error

	if *useLatestVersion {
		useELFPath = true
		unifiVersion, unifiDownloadUrl, err = latestUnifiOSVersion()
		if err != nil {
			panic(err)
		}
	} else {
		unifiVersion, err = version.NewVersion(specifiedVersion)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		segments := unifiVersion.Segments()
		if len(segments) > 0 && segments[0] >= 10 {
			fmt.Print("error: explicit 10.x+ versions are not supported; use -latest\n\n")
			usage()
			os.Exit(1)
		}

		useELFPath = false
		unifiDownloadUrl, err = url.Parse(fmt.Sprintf("https://dl.ui.com/unifi/%s/unifi_sysvinit_all.deb", unifiVersion))
		if err != nil {
			panic(err)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get the current filename")
	}

	versionBaseDir := filepath.Dir(filename)

	// For the ELF path, fieldsDir uses the Network version (discovered
	// during extraction), not the OS version. We compute it after
	// extraction below. For the .deb path, use the firmware API version.
	var fieldsDir string
	if !useELFPath {
		fieldsDir = filepath.Join(versionBaseDir, fmt.Sprintf("v%s", unifiVersion))
	}

	outDir := filepath.Join(wd, *outputDirFlag)

	if useELFPath {
		// The ELF path always re-extracts (the OS version changes
		// whenever Ubiquiti ships a new release). Use the OS version
		// for the temp dir, then rename to the Network version.
		osVersionDir := filepath.Join(versionBaseDir, fmt.Sprintf("v%s", unifiVersion))
		err = os.MkdirAll(osVersionDir, 0o755)
		if err != nil {
			panic(err)
		}

		elfPath, err := downloadELFInstaller(unifiDownloadUrl, osVersionDir)
		if err != nil {
			panic(err)
		}

		aceJarPath, networkVersion, err := extractACFromELF(elfPath, osVersionDir)
		if err != nil {
			panic(err)
		}
		os.Remove(elfPath) // clean up the 880 MB installer

		_, err = extractInternalDepsJSON(aceJarPath, osVersionDir)
		if err != nil {
			panic(err)
		}
		os.Remove(aceJarPath) // clean up the 116 MB jar

		// Rename the temp dir to the Network version.
		networkVersionDir := filepath.Join(versionBaseDir, fmt.Sprintf("v%s", networkVersion))
		if osVersionDir != networkVersionDir {
			os.RemoveAll(networkVersionDir)
			if err := os.Rename(osVersionDir, networkVersionDir); err != nil {
				panic(err)
			}
		}
		fieldsDir = networkVersionDir

		err = copyCustom(fieldsDir)
		if err != nil {
			panic(err)
		}

		// Override unifiVersion with the Network version so the
		// version.generated.go const is correct.
		unifiVersion, err = version.NewVersion(networkVersion)
		if err != nil {
			panic(err)
		}
	} else {
		fieldsInfo, err := os.Stat(fieldsDir)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				panic(err)
			}

			err = os.MkdirAll(fieldsDir, 0o755)
			if err != nil {
				panic(err)
			}

			jarFile, err := downloadJar(unifiDownloadUrl, fieldsDir)
			if err != nil {
				panic(err)
			}

			err = extractJSON(jarFile, fieldsDir)
			if err != nil {
				panic(err)
			}

			err = copyCustom(fieldsDir)
			if err != nil {
				panic(err)
			}

			fieldsInfo, err = os.Stat(fieldsDir)
			if err != nil {
				panic(err)
			}
		}
		_ = fieldsInfo
	}
```

- [ ] **Step 3: Remove the now-redundant `fieldsInfo` check after the if/else**

After the if/else block above, the original code had `if !fieldsInfo.IsDir()`. For the ELF path there's no `fieldsInfo`. Replace the original `if !fieldsInfo.IsDir()` block with:

```go
	fieldsInfo, err := os.Stat(fieldsDir)
	if err != nil {
		panic(err)
	}
	if !fieldsInfo.IsDir() {
		panic("version info isn't a directory")
	}
```

This runs for both paths after `fieldsDir` is set.

- [ ] **Step 4: Build to verify it compiles**

```bash
go build ./cmd/fields/
```

Expected: builds with no errors.

- [ ] **Step 5: Run existing tests**

```bash
go test ./cmd/fields/ -v
```

Expected: all existing tests still pass (none should be affected — they test `latestUnifiVersion`, not `main()`).

- [ ] **Step 6: Commit**

```bash
git add cmd/fields/main.go
git commit -m "fields: route -latest to ELF path, keep .deb for -version <10"
```

---

### Task 8: Add `SensitiveMetadata` struct and loader

**Files:**
- Modify: `cmd/fields/schema.go`
- Modify: `cmd/fields/schema_test.go`

- [ ] **Step 1: Write the failing test**

Append to `cmd/fields/schema_test.go`:

```go
func TestLoadSensitiveMetadata(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	content := `{
		"sensitive_db_fields_by_collection": {
			"networkconf": ["name", "domain_name"],
			"wlanconf": ["x_passphrase"]
		},
		"sensitive_distinct_db_fields_by_collection": {
			"rogue": "essid"
		}
	}`

	path := filepath.Join(t.TempDir(), "sensitive_metadata.json")
	require.NoError(os.WriteFile(path, []byte(content), 0o644))

	meta, err := loadSensitiveMetadata(path)
	require.NoError(err)
	require.NotNil(meta)

	assert.Equal([]string{"name", "domain_name"}, meta.SensitiveDBFieldsByCollection["networkconf"])
	assert.Equal([]string{"x_passphrase"}, meta.SensitiveDBFieldsByCollection["wlanconf"])
	assert.Equal("essid", meta.SensitiveDistinctDBFieldsByCollection["rogue"])
}

func TestLoadSensitiveMetadata_Absent(t *testing.T) {
	require := require.New(t)
	meta, err := loadSensitiveMetadata(filepath.Join(t.TempDir(), "nonexistent.json"))
	require.NoError(err)
	require.Nil(meta)
}
```

Add `"os"`, `"path/filepath"` to the test file imports if not present.

- [ ] **Step 2: Run the test to verify it fails**

```bash
go test ./cmd/fields/ -run TestLoadSensitiveMetadata -v
```

Expected: FAIL with `loadSensitiveMetadata undefined`.

- [ ] **Step 3: Write the implementation**

Add to `cmd/fields/schema.go` (near the top, after the imports and consts):

```go
// SensitiveMetadata mirrors sensitive_metadata.json from the Unifi
// internal-dependencies.jar. It maps Mongo collection names to lists of
// sensitive field paths (dotted for nested fields).
type SensitiveMetadata struct {
	SensitiveDBFieldsByCollection          map[string][]string `json:"sensitive_db_fields_by_collection"`
	SensitiveDistinctDBFieldsByCollection map[string]string   `json:"sensitive_distinct_db_fields_by_collection"`
}

// loadSensitiveMetadata reads and parses sensitive_metadata.json. Returns
// (nil, nil) if the file does not exist (the 9.x .deb path produces no
// metadata file).
func loadSensitiveMetadata(path string) (*SensitiveMetadata, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read sensitive_metadata.json: %w", err)
	}

	var meta SensitiveMetadata
	if err := json.Unmarshal(b, &meta); err != nil {
		return nil, fmt.Errorf("parse sensitive_metadata.json: %w", err)
	}
	return &meta, nil
}
```

Add `"errors"` to the imports in schema.go if not already present (it likely is).

- [ ] **Step 4: Run the test to verify it passes**

```bash
go test ./cmd/fields/ -run TestLoadSensitiveMetadata -v
```

Expected: PASS (both tests).

- [ ] **Step 5: Commit**

```bash
git add cmd/fields/schema.go cmd/fields/schema_test.go
git commit -m "fields: add SensitiveMetadata struct and loader"
```

---

### Task 9: Add `CollectionName()` and `MarkSensitiveFields`

**Files:**
- Modify: `cmd/fields/main.go` (add `Sensitive` to `FieldInfo`)
- Modify: `cmd/fields/schema.go` (add `CollectionName`, `MarkSensitiveFields`, `walkAndMark`)
- Modify: `cmd/fields/schema_test.go`

- [ ] **Step 1: Add `Sensitive` to `FieldInfo`**

In `cmd/fields/main.go`, find the `FieldInfo` struct definition (around line 124):

```go
type FieldInfo struct {
	FieldName           string
	JSONName            string
	FieldType           string
	IsPointer           bool
	FieldValidation     string
	OmitEmpty           bool
	IsArray             bool
	Fields              map[string]*FieldInfo
	CustomUnmarshalType string
	CustomUnmarshalFunc string
}
```

Add `Sensitive bool` as the last field:

```go
type FieldInfo struct {
	FieldName           string
	JSONName            string
	FieldType           string
	IsPointer           bool
	FieldValidation     string
	OmitEmpty           bool
	IsArray             bool
	Fields              map[string]*FieldInfo
	CustomUnmarshalType string
	CustomUnmarshalFunc string
	Sensitive           bool
}
```

- [ ] **Step 2: Write the failing tests for `CollectionName` and `MarkSensitiveFields`**

Append to `cmd/fields/schema_test.go`:

```go
func TestCollectionName(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		structName string
		want       string
	}{
		{"Network", "networkconf"},
		{"WLAN", "wlanconf"},
		{"Device", "device"},
		{"Account", "account"},
		{"RADIUSProfile", "radiusprofile"},
		{"Client", "user"},
		{"ClientGroup", "usergroup"},
		{"DPIGroup", "dpigroup"},
		{"DynamicDNS", "dynamicdns"},
		{"SettingMgmt", "setting"},
		{"SettingUsg", "setting"},
		{"FirewallRule", "firewallrule"},
		{"FirewallGroup", "firewallgroup"},
	}

	for _, tt := range tests {
		t.Run(tt.structName, func(t *testing.T) {
			r := &ResourceInfo{StructName: tt.structName}
			assert.Equal(tt.want, r.CollectionName(), "struct %q", tt.structName)
		})
	}
}

func TestMarkSensitiveFields(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	meta := &SensitiveMetadata{
		SensitiveDBFieldsByCollection: map[string][]string{
			"networkconf": {"name", "domain_name", "x_wireguard_private_key"},
			"radiusprofile": {"name", "auth_servers.x_secret"},
		},
	}

	r := NewResource("Network", "network")
	// Add some fields to test.
	r.Types["Network"].Fields["Name"] = NewFieldInfo("Name", "name", fields.String, "", true, false, false, "")
	r.Types["Network"].Fields["DomainName"] = NewFieldInfo("DomainName", "domain_name", fields.String, "", true, false, false, "")
	r.Types["Network"].Fields["Purpose"] = NewFieldInfo("Purpose", "purpose", fields.String, "", false, false, false, "")

	r.MarkSensitiveFields(meta)

	assert.True(r.Types["Network"].Fields["Name"].Sensitive, "name should be sensitive")
	assert.True(r.Types["Network"].Fields["DomainName"].Sensitive, "domain_name should be sensitive")
	assert.False(r.Types["Network"].Fields["Purpose"].Sensitive, "purpose should NOT be sensitive")
}

func TestMarkSensitiveFields_DottedPath(t *testing.T) {
	assert := assert.New(t)

	meta := &SensitiveMetadata{
		SensitiveDBFieldsByCollection: map[string][]string{
			"radiusprofile": {"name", "auth_servers.x_secret"},
		},
	}

	r := NewResource("RADIUSProfile", "radiusprofile")
	// Add a nested field: auth_servers → x_secret.
	authServers := NewFieldInfo("AuthServers", "auth_servers", "AuthServers", "", true, false, true, "")
	authServers.Fields = map[string]*FieldInfo{
		"XSecret": NewFieldInfo("XSecret", "x_secret", fields.String, "", true, false, false, ""),
		"Host":    NewFieldInfo("Host", "host", fields.String, "", true, false, false, ""),
	}
	r.Types["RADIUSProfile"].Fields["AuthServers"] = authServers
	r.Types["AuthServers"] = authServers

	r.MarkSensitiveFields(meta)

	assert.True(authServers.Fields["XSecret"].Sensitive, "auth_servers.x_secret should be sensitive")
	assert.False(authServers.Fields["Host"].Sensitive, "auth_servers.host should NOT be sensitive")
}

func TestMarkSensitiveFields_DistinctCollection(t *testing.T) {
	assert := assert.New(t)

	meta := &SensitiveMetadata{
		SensitiveDBFieldsByCollection: map[string][]string{},
		SensitiveDistinctDBFieldsByCollection: map[string]string{
			"rogue": "essid",
		},
	}

	r := NewResource("Rogue", "rogue")
	r.Types["Rogue"].Fields["Essid"] = NewFieldInfo("Essid", "essid", fields.String, "", true, false, false, "")
	r.Types["Rogue"].Fields["Name"] = NewFieldInfo("Name", "name", fields.String, "", true, false, false, "")

	r.MarkSensitiveFields(meta)

	assert.True(r.Types["Rogue"].Fields["Essid"].Sensitive, "essid should be sensitive (distinct)")
	assert.False(r.Types["Rogue"].Fields["Name"].Sensitive, "name should NOT be sensitive")
}

func TestMarkSensitiveFields_NoCollection(t *testing.T) {
	assert := assert.New(t)

	meta := &SensitiveMetadata{
		SensitiveDBFieldsByCollection: map[string][]string{
			"networkconf": {"name"},
		},
	}

	r := NewResource("TeleportClient", "teleport_client")
	r.Types["TeleportClient"].Fields["Name"] = NewFieldInfo("Name", "name", fields.String, "", true, false, false, "")

	// Should not panic and should not mark anything (no mapping for teleport_client).
	r.MarkSensitiveFields(meta)
	assert.False(r.Types["TeleportClient"].Fields["Name"].Sensitive)
}
```

- [ ] **Step 3: Run the test to verify it fails**

```bash
go test ./cmd/fields/ -run "TestCollectionName|TestMarkSensitiveFields" -v
```

Expected: FAIL with `CollectionName undefined` and/or `MarkSensitiveFields undefined`.

- [ ] **Step 4: Write the implementation**

Add to `cmd/fields/schema.go` (or `main.go` — since `ResourceInfo` is defined in `main.go`, these methods can go in either file. Put them in `schema.go` to keep schema logic together, but they're methods on `ResourceInfo` which is in `main.go` — same package, so it works):

```go
// collectionToStruct maps Mongo collection names (from sensitive_metadata.json)
// to the resource struct names used by the generator. Built as a reverse of
// fileReps with special handling for settings (all Setting* map to "setting").
var collectionToStruct = map[string]string{
	"networkconf":   "Network",
	"wlanconf":      "WLAN",
	"device":        "Device",
	"account":       "Account",
	"radiusprofile":  "RADIUSProfile",
	"user":          "Client",
	"usergroup":     "ClientGroup",
	"dpigroup":      "DPIGroup",
	"dynamicdns":    "DynamicDNS",
	"firewallrule":  "FirewallRule",
	"firewallgroup": "FirewallGroup",
}

// structToCollection is the inverse, built at init time.
var structToCollection map[string]string

func init() {
	structToCollection = make(map[string]string, len(collectionToStruct))
	for coll, structName := range collectionToStruct {
		structToCollection[structName] = coll
	}
}

// CollectionName returns the Mongo collection name for this resource,
// or "" if no mapping exists.
func (r *ResourceInfo) CollectionName() string {
	// All Setting* resources map to the "setting" collection.
	if strings.HasPrefix(r.StructName, "Setting") {
		return "setting"
	}
	return structToCollection[r.StructName]
}

// MarkSensitiveFields walks the resource's type tree and sets Sensitive=true
// on any field whose JSONName path appears in the sensitive metadata for this
// resource's collection. Dotted paths (e.g. "auth_servers.x_secret") are
// matched by recursing through nested Fields.
func (r *ResourceInfo) MarkSensitiveFields(meta *SensitiveMetadata) {
	if meta == nil {
		return
	}
	collection := r.CollectionName()
	if collection == "" {
		return
	}

	paths := meta.SensitiveDBFieldsByCollection[collection]
	if len(paths) == 0 && meta.SensitiveDistinctDBFieldsByCollection == nil {
		return
	}

	set := make(map[string]struct{}, len(paths)+1)
	for _, p := range paths {
		set[p] = struct{}{}
	}
	if distinct, ok := meta.SensitiveDistinctDBFieldsByCollection[collection]; ok {
		set[distinct] = struct{}{}
	}

	base := r.Types[r.StructName]
	if base == nil {
		return
	}
	r.walkAndMark(base, "", set)
}

// walkAndMark recurses through field.Fields, building a dotted path from
// JSONName values, and sets Sensitive=true on any field whose path is in
// the set.
func (r *ResourceInfo) walkAndMark(field *FieldInfo, parent string, set map[string]struct{}) {
	if field == nil {
		return
	}
	path := field.JSONName
	if parent != "" {
		path = parent + "." + path
	}
	if _, hit := set[path]; hit {
		field.Sensitive = true
	}
	for _, child := range field.Fields {
		r.walkAndMark(child, path, set)
	}
}
```

- [ ] **Step 5: Run the test to verify it passes**

```bash
go test ./cmd/fields/ -run "TestCollectionName|TestMarkSensitiveFields" -v
```

Expected: PASS (all tests).

If `TestCollectionName` fails for a struct not in the map (e.g. `FirewallRule`), add it to `collectionToStruct`. The map above includes `firewallrule` and `firewallgroup` — verify the test expects those.

- [ ] **Step 6: Commit**

```bash
git add cmd/fields/main.go cmd/fields/schema.go cmd/fields/schema_test.go
git commit -m "fields: add CollectionName, MarkSensitiveFields, and FieldInfo.Sensitive"
```

---

### Task 10: Wire `Sensitive` into attribute generation

**Files:**
- Modify: `cmd/fields/schema.go`
- Modify: `cmd/fields/schema_test.go`

- [ ] **Step 1: Write the failing test**

Append to `cmd/fields/schema_test.go`:

```go
func TestFieldToResourceAttribute_Sensitive(t *testing.T) {
	assert := assert.New(t)

	r := NewResource("Account", "account")
	r.Types["Account"].Fields["Name"] = NewFieldInfo("Name", "name", fields.String, "", true, false, false, "")
	r.Types["Account"].Fields["XPassword"] = NewFieldInfo("XPassword", "x_password", fields.String, "", true, false, false, "")
	r.Types["Account"].Fields["XPassword"].Sensitive = true

	gen := NewSpecificationGenerator("unifi")
	attr := gen.fieldToResourceAttribute(r, r.Types["Account"].Fields["XPassword"])
	require.NotNil(t, attr)
	require.NotNil(t, attr.String)
	assert.NotNil(t, attr.String.Sensitive, "x_password should have Sensitive set")
	assert.True(*attr.String.Sensitive)

	// Non-sensitive field should not have Sensitive set.
	attr2 := gen.fieldToResourceAttribute(r, r.Types["Account"].Fields["Name"])
	require.NotNil(t, attr2)
	require.NotNil(t, attr2.String)
	assert.Nil(attr2.String.Sensitive, "name should NOT have Sensitive set")
}
```

- [ ] **Step 2: Run the test to verify it fails**

```bash
go test ./cmd/fields/ -run TestFieldToResourceAttribute_Sensitive -v
```

Expected: FAIL — `attr.String.Sensitive` is nil because the generator never sets it.

- [ ] **Step 3: Add the sensitive-setting logic to `fieldToResourceAttribute`**

In `cmd/fields/schema.go`, find `fieldToResourceAttribute`. At the end of the function, just before `return attr`, add:

```go
	if field.Sensitive {
		if attr.String != nil {
			attr.String.Sensitive = ptr(true)
		}
		if attr.Bool != nil {
			attr.Bool.Sensitive = ptr(true)
		}
		if attr.Int64 != nil {
			attr.Int64.Sensitive = ptr(true)
		}
		if attr.Float64 != nil {
			attr.Float64.Sensitive = ptr(true)
		}
	}
```

Do the same for `fieldToDataSourceAttribute` — add the same block before `return attr`.

Note: for `ListNested`/`SingleNested`/`List` attributes, we don't set Sensitive at the container level (codegen-spec doesn't have that field on containers). The pre-pass already marked the leaf children, and the recursive calls to `generateNestedResourceAttributes`/`generateNestedDataSourceAttributes` will pick those up because they call `fieldToResourceAttribute`/`fieldToDataSourceAttribute` on each child.

- [ ] **Step 4: Run the test to verify it passes**

```bash
go test ./cmd/fields/ -run TestFieldToResourceAttribute_Sensitive -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add cmd/fields/schema.go cmd/fields/schema_test.go
git commit -m "fields: propagate FieldInfo.Sensitive to spec attributes"
```

---

### Task 11: Wire sensitive metadata loading into `main.go`

**Files:**
- Modify: `cmd/fields/main.go`

- [ ] **Step 1: Add the metadata loading + marking calls**

In `cmd/fields/main.go`, find the code generation loop. It currently starts around:

```go
	fieldsFiles, err := os.ReadDir(fieldsDir)
	...
	// Initialize specification generator
	specGen := NewSpecificationGenerator("unifi")

	for _, fieldsFile := range fieldsFiles {
		...
		err = resource.processJSON(b)
		...
		// Add resource to specification generator
		specGen.AddResource(resource)
		...
	}
```

After `specGen := NewSpecificationGenerator("unifi")` and before the loop, add the metadata loading:

```go
	// Load sensitive metadata if present (10.x ELF path; absent for 9.x .deb).
	sensitiveMeta, err := loadSensitiveMetadata(filepath.Join(fieldsDir, "sensitive_metadata.json"))
	if err != nil {
		panic(err)
	}
```

Then inside the loop, after `err = resource.processJSON(b)` and before `specGen.AddResource(resource)`, add:

```go
		// Mark sensitive fields from metadata (no-op if meta is nil).
		resource.MarkSensitiveFields(sensitiveMeta)
```

- [ ] **Step 2: Build to verify it compiles**

```bash
go build ./cmd/fields/
```

Expected: builds with no errors.

- [ ] **Step 3: Run all tests**

```bash
go test ./cmd/fields/ -v
```

Expected: all tests pass.

- [ ] **Step 4: Commit**

```bash
git add cmd/fields/main.go
git commit -m "fields: load sensitive_metadata.json and mark fields before codegen"
```

---

### Task 12: Rewrite `generate.yaml` for auto-merge

**Files:**
- Modify: `.github/workflows/generate.yaml`

- [ ] **Step 1: Replace the workflow file**

Overwrite `.github/workflows/generate.yaml` with:

```yaml
---
name: Schema Generation

on:
  schedule:
    - cron: 0 0 * * *
  workflow_dispatch: {}

permissions:
  contents: write
  pull-requests: write

jobs:
  fields:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v7.0.0

      - name: Setup Go
        uses: actions/setup-go@924ae3a1cded613372ab5595356fb5720e22ba16 # v6.5.0
        with:
          go-version-file: 'go.mod'

      - name: Generate
        run: go generate ./...

      - name: Test
        run: go test ./...

      - name: Check for changes
        id: changes
        run: |
          if git diff --quiet; then
            echo "changed=false" >> $GITHUB_OUTPUT
          else
            echo "changed=true" >> $GITHUB_OUTPUT
          fi

      - name: Create PR
        if: steps.changes.outputs.changed == 'true'
        id: create-pr
        uses: peter-evans/create-pull-request@5f6978faf089d4d20b00c7766989d076bb2fc7f1 # v8.1.1
        with:
          delete-branch: true
          title: 'Update to latest controller version'
          labels: 'auto-generated'

      - name: Enable auto-merge
        if: steps.changes.outputs.changed == 'true' && steps.create-pr.outputs.pull-request-number != ''
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          gh pr merge "${{ steps.create-pr.outputs.pull-request-number }}" --auto --squash --delete-branch
```

Note: the `peter-evans/create-pull-request` action exposes `pull-request-number` as a step output, which we use directly for the `gh pr merge` call. This is more robust than looking up the PR by branch name.

- [ ] **Step 2: Lint the YAML**

```bash
yamllint .github/workflows/generate.yaml
```

Expected: no errors (or only warnings). Fix any real issues.

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/generate.yaml
git commit -m "ci: auto-merge generated schema PRs on green CI"
```

---

### Task 13: Add `tag.yaml` for automatic version tagging

**Files:**
- Create: `.github/workflows/tag.yaml`

- [ ] **Step 1: Create the workflow file**

Create `.github/workflows/tag.yaml`:

```yaml
---
name: Tag Release

on:
  push:
    branches:
      - main

permissions:
  contents: write

jobs:
  tag:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v7.0.0
        with:
          fetch-depth: 0

      - name: Read UnifiVersion
        id: version
        run: |
          VERSION=$(grep -oP '(?<=UnifiVersion = ")[^"]*' unifi/version.generated.go)
          echo "version=$VERSION" >> $GITHUB_OUTPUT

      - name: Check if tag exists
        id: exists
        run: |
          if git rev-parse "refs/tags/v${{ steps.version.outputs.version }}" >/dev/null 2>&1; then
            echo "exists=true" >> $GITHUB_OUTPUT
          else
            echo "exists=false" >> $GITHUB_OUTPUT
          fi

      - name: Create and push tag
        if: steps.exists.outputs.exists == 'false'
        run: |
          git tag "v${{ steps.version.outputs.version }}"
          git push origin "v${{ steps.version.outputs.version }}"
```

Note: `grep -oP` uses Perl-compatible regex (the `-P` flag). On macOS this isn't available by default, but GitHub Actions runners use Ubuntu, where `grep -P` is supported. If you want to test locally on macOS, use `grep -oE` with a different pattern, but keep `-P` for the workflow.

- [ ] **Step 2: Lint the YAML**

```bash
yamllint .github/workflows/tag.yaml
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/tag.yaml
git commit -m "ci: auto-tag v<unifi-version> on merge to main"
```

---

### Task 14: Full regeneration and verification

This task runs the full pipeline end-to-end to verify everything works together. It downloads the real ~880 MB installer, so it takes a few minutes. Run it on a machine with enough disk space (~2 GB free) and bandwidth.

**Files:**
- No file changes expected beyond regenerated `unifi/*.generated.go`, `unifi/version.generated.go`, `specification.json`.

- [ ] **Step 1: Run the full generator**

```bash
go generate ./...
```

Expected: downloads the ELF (~880 MB), extracts image.tar, parses OCI manifest, finds ace.jar, extracts internal-dependencies.jar, writes `cmd/fields/v10.4.57/` (or whatever the current Network version is), generates `unifi/*.generated.go`, `unifi/version.generated.go`, `specification.json`.

If this fails with "image.tar is not an OCI image layout (missing index.json)", the format spike (Task 1) was wrong or the format changed. Stop and revisit the spec.

- [ ] **Step 2: Verify the version file**

```bash
cat unifi/version.generated.go
```

Expected: `const UnifiVersion = "10.4.57"` (or the current Network version — NOT the OS version `5.1.21`).

- [ ] **Step 3: Verify specification.json has sensitive fields**

```bash
python3 -c "
import json
with open('specification.json') as f:
    spec = json.load(f)
sensitive_count = 0
for res in spec.get('resources', []):
    for attr in res.get('schema', {}).get('attributes', []):
        def check(obj):
            global sensitive_count
            if isinstance(obj, dict):
                for k, v in obj.items():
                    if k == 'sensitive' and v is True:
                        sensitive_count += 1
                    check(v)
            elif isinstance(obj, list):
                for i in obj:
                    check(i)
        check(attr)
print(f'Found {sensitive_count} sensitive attributes')
"
```

Expected: a non-zero count (e.g. ~50+). If zero, the sensitive metadata wiring isn't working — check that `sensitive_metadata.json` was extracted to `cmd/fields/v10.4.57/` and that `MarkSensitiveFields` is being called.

- [ ] **Step 4: Run the full test suite**

```bash
go test ./...
```

Expected: all tests pass. Some existing tests may fail if the generated code changed field types or names — investigate each failure. If a generated field changed type (e.g. a string became an int), the test expectations in `unifi/*_test.go` may need updating. This is expected when bumping controller versions.

- [ ] **Step 5: Run linters**

```bash
go vet ./...
golangci-lint run --timeout=10m
```

Expected: no new issues introduced by the generator code. Generated code in `unifi/*.generated.go` may have pre-existing lint issues — don't fix those here.

- [ ] **Step 6: Verify git diff is clean after re-running generation**

```bash
go generate ./...
git diff --compact-summary --exit-code
```

Expected: exit code 0 (no diff). If there's a diff, the generator is non-idempotent — investigate.

- [ ] **Step 7: Commit the regenerated files**

```bash
git add unifi/ specification.json
git commit -m "fields: regenerate for Unifi Network 10.4.57 with sensitive metadata"
```

(Replace `10.4.57` with the actual version if different.)

---

### Task 15: Update README and documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update the "Note on Code Generation" section**

In `README.md`, find the existing "Note on Code Generation" section (lines 9-17). Replace it with:

```markdown
## Note on Code Generation

The data models and basic REST methods are generated from JSON field-definition
files shipped inside the Unifi Network application's `internal-dependencies.jar`
(bundled inside `ace.jar` in the OCI image shipped by the Unifi OS installer).

To regenerate the code, run `go generate ./...` inside the repo root. This
downloads the latest Unifi OS installer, extracts the field definitions, and
regenerates `unifi/*.generated.go` and `specification.json`.

For older (pre-10.x) controller versions that shipped as `.deb` packages, use
`go run ./cmd/fields/ -version 9.5.21` instead.

The `specification.json` file includes `sensitive: true` flags on fields
identified as sensitive by `sensitive_metadata.json` from the Unifi jar.
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: update README for ELF extraction + sensitive metadata"
```

---

## Self-Review Checklist

After implementing all tasks, verify:

1. **Spec coverage:** Every section of the spec has a corresponding task.
   - ELF extractor (3 functions) → Tasks 3, 4, 5
   - Firmware API (`latestUnifiOSVersion`) → Task 6
   - `main.go` routing → Task 7
   - Sensitive metadata (struct, loader, collection mapping, marking, attribute generation) → Tasks 8, 9, 10, 11
   - GHA automation (`generate.yaml` rewrite, `tag.yaml` new) → Tasks 12, 13
   - Full verification → Task 14
   - Documentation → Task 15

2. **Type consistency:** `FieldInfo.Sensitive` (added in Task 9) is read in Task 10's attribute generation. `SensitiveMetadata` (defined in Task 8) is loaded in Task 11. `CollectionName()` and `MarkSensitiveFields` (defined in Task 9) are called in Task 11. All signatures match.

3. **No placeholders:** Every code block is complete. Every test has actual assertions. Every commit message is specified.
