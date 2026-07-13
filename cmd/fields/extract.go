package main

import (
	"archive/tar"
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/ulikunitz/xz"
	"github.com/xor-gate/ar"
)

// fieldsJarPaths are the jars inside the controller .deb that may carry the
// api/fields JSON definitions, in precedence order: 10.x ships them in
// internal-dependencies.jar (ace.jar became a launcher stub), 9.x and earlier
// in the fat ace.jar. The first jar that yields field JSON wins; definitions
// are never merged across jars.
var fieldsJarPaths = []string{
	"./usr/lib/unifi/lib/internal/internal-dependencies.jar",
	"./usr/lib/unifi/lib/ace.jar",
}

// DownloadAndExtract downloads the controller package and populates fieldsDir
// with the api/fields JSON definitions plus the custom overrides. Everything
// is staged in a temp dir and only moved into place once complete, so a
// failed run can never leave a partial fields dir that a later run mistakes
// for a valid cache.
func DownloadAndExtract(url *url.URL, fieldsDir string) error {
	fieldsInfo, err := os.Stat(fieldsDir)
	switch {
	case err == nil:
		if !fieldsInfo.IsDir() {
			return errors.New("version info isn't a directory")
		}

		cached, err := isValidFieldsCache(fieldsDir)
		if err != nil {
			return err
		}
		if cached {
			// Already downloaded and extracted.
			return nil
		}
		// Leftovers from a broken or pre-fix run (a stub jar, possibly with
		// the custom overrides already copied in): re-download.
	case errors.Is(err, os.ErrNotExist):
	default:
		return err
	}

	tmpDir := fieldsDir + ".tmp"
	if err := os.RemoveAll(tmpDir); err != nil {
		return err
	}
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	jarFiles, err := downloadJars(url, tmpDir)
	if err != nil {
		return err
	}

	if err := extractJSON(jarFiles, tmpDir); err != nil {
		return err
	}

	if err := copyCustom(tmpDir); err != nil {
		return err
	}

	// The jars have served their purpose; cache only the JSON.
	for _, jarFile := range jarFiles {
		if err := os.Remove(jarFile); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(fieldsDir); err != nil {
		return err
	}
	return os.Rename(tmpDir, fieldsDir)
}

// isValidFieldsCache reports whether dir holds a completed extraction.
// Setting.json is the marker: extractJSON fails hard without it, and unlike
// the custom overrides (which pre-fix runs copied into otherwise-empty dirs)
// it can only come from a controller package.
func isValidFieldsCache(dir string) (bool, error) {
	_, err := os.Stat(filepath.Join(dir, "Setting.json"))
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	default:
		return false, err
	}
}

func downloadJars(url *url.URL, outputDir string) ([]string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to download deb: %w", err)
	}

	debResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to download deb: %w", err)
	}
	defer debResp.Body.Close()

	if debResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to download deb: unexpected status %q from %s", debResp.Status, url)
	}

	var uncompressedReader io.Reader

	arReader := ar.NewReader(debResp.Body)
	for {
		header, err := arReader.Next()
		if errors.Is(err, io.EOF) || header == nil {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("in ar next: %w", err)
		}

		// Read the data file.
		if header.Name == "data.tar.xz" {
			uncompressedReader, err = xz.NewReader(arReader)
			if err != nil {
				return nil, fmt.Errorf("in xz reader: %w", err)
			}
			break
		}
	}
	if uncompressedReader == nil {
		return nil, errors.New("unable to find .deb data file")
	}

	tarReader := tar.NewReader(uncompressedReader)

	found := make([]string, len(fieldsJarPaths))
	foundCount := 0

	for foundCount < len(fieldsJarPaths) {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("in next: %w", err)
		}

		idx := slices.Index(fieldsJarPaths, header.Name)
		if header.Typeflag != tar.TypeReg || idx < 0 || found[idx] != "" {
			// Skipping.
			continue
		}

		dstPath := filepath.Join(outputDir, filepath.Base(header.Name))
		if err := extractTarFile(tarReader, dstPath); err != nil {
			return nil, err
		}
		found[idx] = dstPath
		foundCount++
	}

	// Keep precedence order regardless of where the jars sat in the tar.
	jarFiles := make([]string, 0, foundCount)
	for _, jarFile := range found {
		if jarFile != "" {
			jarFiles = append(jarFiles, jarFile)
		}
	}

	if len(jarFiles) == 0 {
		return nil, fmt.Errorf("unable to find any field definition jar (%s) in controller package", strings.Join(fieldsJarPaths, ", "))
	}

	return jarFiles, nil
}

func extractTarFile(tarReader *tar.Reader, dstPath string) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("unable to create temp file: %w", err)
	}

	if _, err := io.Copy(dst, tarReader); err != nil {
		_ = dst.Close()
		return fmt.Errorf("unable to write %s temp file: %w", filepath.Base(dstPath), err)
	}

	if err := dst.Close(); err != nil {
		return fmt.Errorf("unable to write %s temp file: %w", filepath.Base(dstPath), err)
	}

	return nil
}

// extractJSONFromJar copies the api/fields/*.json definitions out of one jar
// and returns how many it found.
func extractJSONFromJar(jarFile, fieldsDir string) (int, error) {
	jarZip, err := zip.OpenReader(jarFile)
	if err != nil {
		return 0, fmt.Errorf("unable to open jar: %w", err)
	}
	defer jarZip.Close()

	extracted := 0
	for _, f := range jarZip.File {
		if !strings.HasPrefix(f.Name, "api/fields/") || path.Ext(f.Name) != ".json" {
			// Skip file.
			continue
		}

		err = func() error {
			src, err := f.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			dst, err := os.Create(filepath.Join(fieldsDir, filepath.Base(f.Name)))
			if err != nil {
				return err
			}
			defer dst.Close()

			_, err = io.Copy(dst, src)
			if err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return 0, fmt.Errorf("unable to write JSON file: %w", err)
		}
		extracted++
	}

	return extracted, nil
}

func extractJSON(jarFiles []string, fieldsDir string) error {
	// jarFiles is in precedence order; the first jar carrying field
	// definitions wins so two jars can never mix definitions.
	extracted := 0
	for _, jarFile := range jarFiles {
		n, err := extractJSONFromJar(jarFile, fieldsDir)
		if err != nil {
			return err
		}
		if n > 0 {
			extracted = n
			break
		}
	}

	if extracted == 0 {
		return fmt.Errorf("no api/fields/*.json field definitions found in %s; the controller package layout may have changed", strings.Join(jarFiles, ", "))
	}

	// Setting.json is required: most of the generated Setting resources come
	// from its per-section split, so a package without it means the layout
	// changed and generation must fail loudly rather than silently emit an
	// SDK missing every Setting resource.
	settingsData, err := os.ReadFile(filepath.Join(fieldsDir, "Setting.json"))
	if err != nil {
		return fmt.Errorf("unable to open settings file: %w", err)
	}

	var settings map[string]any
	err = json.Unmarshal(settingsData, &settings)
	if err != nil {
		return fmt.Errorf("unable to unmarshal settings: %w", err)
	}

	for k, v := range settings {
		fileName := fmt.Sprintf("Setting%s.json", strcase.ToCamel(k))

		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to marshal setting %q: %w", k, err)
		}

		err = os.WriteFile(filepath.Join(fieldsDir, fileName), data, 0o755)
		if err != nil {
			return fmt.Errorf("unable to write new settings file: %w", err)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootDir := findModuleRoot(wd)
	srcDir := path.Join(rootDir, "cmd", "fields")

	files, err := os.ReadDir(path.Join(srcDir, "custom"))
	if err != nil {
		return fmt.Errorf("unable to read custom directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			fs, err := os.Open(path.Join(srcDir, "custom", file.Name()))
			if err != nil {
				return fmt.Errorf("unable to open file: %w", err)
			}
			defer fs.Close()
			rf, err := os.Create(filepath.Join(fieldsDir, file.Name()))
			if err != nil {
				return fmt.Errorf("unable to create file: %w", err)
			}
			defer rf.Close()
			_, err = io.Copy(rf, fs)
			if err != nil {
				return fmt.Errorf("unable to copy file: %w", err)
			}
			_, err = io.ReadAll(fs)
			if err != nil {
				return fmt.Errorf("unable to read file: %w", err)
			}
		}
		fmt.Println(file.Name(), file.IsDir())
	}

	// TODO: cleanup JSON
	return nil
}

func copyCustom(fieldsDir string) error {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootDir := findModuleRoot(wd)
	srcDir := path.Join(rootDir, "cmd", "fields")

	files, err := os.ReadDir(path.Join(srcDir, "custom"))
	if err != nil {
		return fmt.Errorf("unable to read custom directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			fs, err := os.Open(path.Join(srcDir, "custom", file.Name()))
			if err != nil {
				return fmt.Errorf("unable to open file: %w", err)
			}
			defer fs.Close()
			rf, err := os.Create(filepath.Join(fieldsDir, file.Name()))
			if err != nil {
				return fmt.Errorf("unable to create file: %w", err)
			}
			defer rf.Close()
			_, err = io.Copy(rf, fs)
			if err != nil {
				return fmt.Errorf("unable to copy file: %w", err)
			}
		}
	}

	return nil
}

func findModuleRoot(dir string) (roots string) {
	if dir == "" {
		panic("dir not set")
	}
	dir = filepath.Clean(dir)
	// Look for enclosing go.mod.
	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		dir = d
	}
	return ""
}
