package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ulikunitz/xz"
	"github.com/xor-gate/ar"
)

const (
	testUserJSON    = `{"mac": "^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$", "name": ""}`
	testSettingJSON = `{"mgmt": {"x_ssh_enabled": "true|false"}}`

	// Tar entry paths as they appear in the controller .deb.
	debAceJarPath          = "./usr/lib/unifi/lib/ace.jar"
	debInternalDepsJarPath = "./usr/lib/unifi/lib/internal/internal-dependencies.jar"
)

// debEntry is one file inside the deb's data.tar.xz; order is preserved.
type debEntry struct {
	name string
	data []byte
}

// buildJar returns an in-memory jar (zip) with the given entries.
func buildJar(t *testing.T, entries map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range entries {
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = w.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

// stubLauncherJar mimics the 10.x ace.jar, which is only a launcher stub with
// no api/fields JSON in it.
func stubLauncherJar(t *testing.T) []byte {
	t.Helper()
	return buildJar(t, map[string]string{
		"META-INF/MANIFEST.MF":        "Main-Class: com.ubnt.ace.Launcher\n",
		"com/ubnt/ace/Launcher.class": "stub",
	})
}

// fieldsJar mimics a jar carrying the api/fields definitions (the 9.x fat
// ace.jar or the 10.x internal-dependencies.jar).
func fieldsJar(t *testing.T) []byte {
	t.Helper()
	return buildJar(t, map[string]string{
		"api/fields/User.json":    testUserJSON,
		"api/fields/Setting.json": testSettingJSON,
		"com/ubnt/some/App.class": "app",
	})
}

// buildDeb wraps the given tar entries into a data.tar.xz inside an ar
// archive, mimicking a UniFi controller .deb.
func buildDeb(t *testing.T, files []debEntry) []byte {
	t.Helper()

	var dataBuf bytes.Buffer
	xzw, err := xz.NewWriter(&dataBuf)
	require.NoError(t, err)
	tw := tar.NewWriter(xzw)
	for _, f := range files {
		require.NoError(t, tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeReg,
			Name:     f.name,
			Mode:     0o644,
			Size:     int64(len(f.data)),
		}))
		_, err = tw.Write(f.data)
		require.NoError(t, err)
	}
	require.NoError(t, tw.Close())
	require.NoError(t, xzw.Close())

	var deb bytes.Buffer
	arw := ar.NewWriter(&deb)
	require.NoError(t, arw.WriteGlobalHeader())
	for _, member := range []debEntry{
		{"debian-binary", []byte("2.0\n")},
		{"data.tar.xz", dataBuf.Bytes()},
	} {
		require.NoError(t, arw.WriteHeader(&ar.Header{
			Name:    member.name,
			ModTime: time.Unix(0, 0),
			Mode:    0o644,
			Size:    int64(len(member.data)),
		}))
		_, err = arw.Write(member.data)
		require.NoError(t, err)
	}
	return deb.Bytes()
}

// serveDeb serves the deb over a local test server and returns its URL.
func serveDeb(t *testing.T, deb []byte, requests *int) *url.URL {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requests != nil {
			*requests++
		}
		_, _ = w.Write(deb)
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)
	return u
}

func assertFieldsExtracted(t *testing.T, dir string) {
	t.Helper()

	userJSON, err := os.ReadFile(filepath.Join(dir, "User.json"))
	require.NoError(t, err, "expected User.json to be extracted")
	assert.JSONEq(t, testUserJSON, string(userJSON))

	_, err = os.Stat(filepath.Join(dir, "SettingMgmt.json"))
	require.NoError(t, err, "expected Setting.json to be split into SettingMgmt.json")

	_, err = os.Stat(filepath.Join(dir, "FirewallPolicy.json"))
	require.NoError(t, err, "expected custom overrides to be copied in")

	jars, err := filepath.Glob(filepath.Join(dir, "*.jar"))
	require.NoError(t, err)
	assert.Empty(t, jars, "extracted jars should be removed from the fields dir")
}

func TestDownloadAndExtractTenXLayout(t *testing.T) {
	t.Parallel()

	// 10.x: ace.jar is a launcher stub, fields live in internal-dependencies.jar.
	deb := buildDeb(t, []debEntry{
		{debAceJarPath, stubLauncherJar(t)},
		{debInternalDepsJarPath, fieldsJar(t)},
	})
	fieldsDir := filepath.Join(t.TempDir(), "v10.4.57")

	require.NoError(t, DownloadAndExtract(serveDeb(t, deb, nil), fieldsDir))
	assertFieldsExtracted(t, fieldsDir)
}

func TestDownloadAndExtractLegacyLayout(t *testing.T) {
	t.Parallel()

	// 9.x and earlier: ace.jar is the fat application jar carrying the fields.
	deb := buildDeb(t, []debEntry{
		{debAceJarPath, fieldsJar(t)},
	})
	fieldsDir := filepath.Join(t.TempDir(), "v9.5.21")

	require.NoError(t, DownloadAndExtract(serveDeb(t, deb, nil), fieldsDir))
	assertFieldsExtracted(t, fieldsDir)
}

func TestDownloadAndExtractPrefersInternalDepsJar(t *testing.T) {
	t.Parallel()

	// When both jars carry field definitions, internal-dependencies.jar (the
	// 10.x application jar) must win deterministically rather than whichever
	// jar happens to come later in the tar stream.
	staleUserJSON := `{"name": "", "stale": "true|false"}`
	aceWithStaleFields := buildJar(t, map[string]string{
		"api/fields/User.json":    staleUserJSON,
		"api/fields/Setting.json": testSettingJSON,
	})
	deb := buildDeb(t, []debEntry{
		{debInternalDepsJarPath, fieldsJar(t)},
		// Last in the tar: naive last-write-wins would pick the stale copy.
		{debAceJarPath, aceWithStaleFields},
	})
	fieldsDir := filepath.Join(t.TempDir(), "v10.4.57")

	require.NoError(t, DownloadAndExtract(serveDeb(t, deb, nil), fieldsDir))

	userJSON, err := os.ReadFile(filepath.Join(fieldsDir, "User.json"))
	require.NoError(t, err)
	assert.JSONEq(t, testUserJSON, string(userJSON))
}

func TestDownloadAndExtractNoFieldDefinitionsErrors(t *testing.T) {
	t.Parallel()

	// A package whose jars carry no field definitions must fail loudly, not
	// succeed with an empty fields directory (the silent no-op that froze
	// generation at 9.5.21).
	deb := buildDeb(t, []debEntry{
		{debAceJarPath, stubLauncherJar(t)},
	})
	fieldsDir := filepath.Join(t.TempDir(), "v10.4.57")

	err := DownloadAndExtract(serveDeb(t, deb, nil), fieldsDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "api/fields")
}

func TestDownloadAndExtractNoJarsErrors(t *testing.T) {
	t.Parallel()

	deb := buildDeb(t, []debEntry{
		{"./usr/lib/unifi/lib/unrelated.txt", []byte("nope")},
	})
	fieldsDir := filepath.Join(t.TempDir(), "v10.4.57")

	err := DownloadAndExtract(serveDeb(t, deb, nil), fieldsDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "jar")
}

func TestDownloadAndExtractFailedRunLeavesNoPartialCache(t *testing.T) {
	t.Parallel()

	// A fields jar without Setting.json fails extraction...
	incompleteJar := buildJar(t, map[string]string{
		"api/fields/User.json": testUserJSON,
	})
	badDeb := buildDeb(t, []debEntry{
		{debAceJarPath, incompleteJar},
	})
	fieldsDir := filepath.Join(t.TempDir(), "v10.4.57")

	err := DownloadAndExtract(serveDeb(t, badDeb, nil), fieldsDir)
	require.Error(t, err)

	// ...and must not leave a partial fields dir that poisons the next run.
	_, statErr := os.Stat(fieldsDir)
	assert.True(t, os.IsNotExist(statErr), "failed extraction should not leave a partial fields dir")

	goodDeb := buildDeb(t, []debEntry{
		{debAceJarPath, fieldsJar(t)},
	})
	require.NoError(t, DownloadAndExtract(serveDeb(t, goodDeb, nil), fieldsDir))
	assertFieldsExtracted(t, fieldsDir)
}

func TestDownloadAndExtractRefreshesStaleCacheDir(t *testing.T) {
	t.Parallel()

	// The pre-fix generator, run against a 10.x package, left behind exactly
	// this shape: a stub ace.jar (no fields extracted) plus the custom
	// overrides copied in by copyCustom. Such a dir must not be treated as a
	// valid fields cache just because it contains .json files.
	fieldsDir := filepath.Join(t.TempDir(), "v10.4.57")
	require.NoError(t, os.MkdirAll(fieldsDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fieldsDir, "ace.jar"), stubLauncherJar(t), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(fieldsDir, "FirewallPolicy.json"), []byte(`{}`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(fieldsDir, "Nat.json"), []byte(`{}`), 0o644))

	deb := buildDeb(t, []debEntry{
		{debAceJarPath, stubLauncherJar(t)},
		{debInternalDepsJarPath, fieldsJar(t)},
	})

	require.NoError(t, DownloadAndExtract(serveDeb(t, deb, nil), fieldsDir))
	assertFieldsExtracted(t, fieldsDir)
}

func TestDownloadAndExtractSkipsDownloadWhenCached(t *testing.T) {
	t.Parallel()

	// Setting.json is the marker of a completed extraction: it is required by
	// extractJSON, and unlike the custom overrides it can only come from a
	// controller package.
	fieldsDir := filepath.Join(t.TempDir(), "v10.4.57")
	require.NoError(t, os.MkdirAll(fieldsDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fieldsDir, "Setting.json"), []byte(testSettingJSON), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(fieldsDir, "User.json"), []byte(testUserJSON), 0o644))

	requests := 0
	deb := buildDeb(t, []debEntry{
		{debAceJarPath, fieldsJar(t)},
	})

	require.NoError(t, DownloadAndExtract(serveDeb(t, deb, &requests), fieldsDir))
	assert.Zero(t, requests, "cached fields dir should not trigger a download")
}

func TestDownloadAndExtractSurfacesHTTPErrors(t *testing.T) {
	t.Parallel()

	// A missing package (e.g. a mistyped version) must surface the HTTP
	// failure, not a misleading archive-parsing error.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "<html>not found</html>", http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)
	u, err := url.Parse(srv.URL)
	require.NoError(t, err)

	fieldsDir := filepath.Join(t.TempDir(), "v99.99.99")
	err = DownloadAndExtract(u, fieldsDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}
