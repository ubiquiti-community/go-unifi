# Unifi Schema Fetcher: ELF/OCI Support, Sensitive Metadata, and Full Automation

## Problem

The go-unifi generator (`cmd/fields/`) is stuck on Unifi Network 9.x. It downloads a
`.deb` from `dl.ui.com`, extracts `ace.jar`, and reads field validators. Ubiquiti stopped
shipping the `.deb` for Unifi Network 10.x; the only source is the Unifi OS installer
(ELF), which bundles an OCI image tar containing `ace.jar`. The generator cannot reach
10.x.

A proof-of-concept Python script in
`terraform-provider-unifi/emdash/extract-find-api-json-definitions-gyr8e/scripts/extract_unifi_api_defs.py`
demonstrates the ELF extraction path. It also surfaces metadata files
(`sensitive_metadata.json`, `radio_specification.json`, etc.) that the current generator
ignores.

The GitHub Actions `generate.yaml` workflow creates a daily PR but never auto-merges or
tags a release, so updates require manual intervention.

## Goals

1. Port the ELF extraction path to Go so the generator can fetch Unifi Network 10.x
   schemas directly from `fw-download.ubnt.com`.
2. Preserve the `.deb` path for 9.x legacy regeneration.
3. Extract the extra metadata files (`sensitive_metadata.json` and seven siblings) and
   wire `sensitive_metadata.json` into the spec generator so matching fields gain
   `Sensitive: true`.
4. Fully automate the release pipeline: daily cron produces a PR, auto-merges on green
   CI, tags a new version, and triggers GoReleaser.

## Non-Goals

- Wiring `radio_specification.json`, `country_codes_list.json`, `timezones.json`, or
  other metadata files into the spec generator. Only `sensitive_metadata.json`.
- Refactoring `cmd/fields/` around a `Fetcher` interface (YAGNI for two sources).
- Shipping the Python PoC; the port is pure Go.
- Branch protection configuration via workflow; that is a manual one-time repo setting.

## Architecture

```
go generate ./...
  └─ cmd/fields -latest (or -version X.Y.Z)
       │
       ├─ version >= 10.x ─→ ELF path (NEW)
       │     1. query fw-update API for unifi-os-server/linux-x64
       │     2. download ELF (self-extracting: zip appended to stub)
       │     3. read image.tar from appended zip
       │     4. parse OCI index.json → manifest → layers
       │     5. for each layer: scan for usr/lib/unifi/lib/ace.jar
       │     6. read BOOT-INF/classes/product.properties → Network version
       │     7. read BOOT-INF/lib/internal-dependencies.jar
       │     8. extract api/fields/*.json + 8 metadata files → cmd/fields/v<net-ver>/
       │
       └─ version < 10.x  ─→ .deb path (EXISTING, legacy)
             1. query fw-update API for unifi-controller/debian (maxVersion < 10.0.0)
             2. download .deb, ar→xz→tar, find ace.jar
             3. extract api/fields/*.json from ace.jar directly → cmd/fields/v<ver>/
             (no metadata files; those only ship in internal-dependencies.jar)
       │
       ├─ copyCustom() (existing, unchanged)
       │
       ├─ for each *.json: processJSON → ResourceInfo → generateCode → *.generated.go
       │
       └─ schema.go: read sensitive_metadata.json (if present)
            → MarkSensitiveFields walks each resource's type tree
            → fieldToResourceAttribute sets Sensitive:true on matching leaf attrs
```

### Version Handling

- `-latest` queries `unifi-os-server`/`linux-x64` (always latest OS; currently v5.1.21 →
  Network 10.4.57). Uses the ELF path.
- `-version X.Y.Z` where `X < 10` uses the `.deb` path (unchanged).
- `-version X.Y.Z` where `X >= 10` is rejected with a helpful error: use `-latest` for
  10.x+. The firmware API is the only source for the ELF URL, and the user cannot
  predict the OS version that corresponds to a given Network version.
- The output directory and `UnifiVersion` const use the **Network** version (e.g.
  `10.4.57`), not the OS version (`5.1.21`). For 10.x, this comes from
  `BOOT-INF/classes/product.properties` inside `ace.jar`. For 9.x, it comes from the
  firmware API response.

### Storage

Extracted JSON (field validators + metadata) remains gitignored under
`cmd/fields/v<version>/`. The metadata files are transient build-time inputs; the PR
diff shows generated Go and `specification.json` changes only.

## Components

### 1. ELF Extractor (`cmd/fields/extract_elf.go`, NEW)

Pure Go, ~250 lines. Uses `github.com/opencontainers/image-spec` (manifest/index
types), `github.com/opencontainers/go-digest` (digest path construction), and stdlib
(`archive/tar`, `archive/zip`, `compress/gzip`, `bufio`, `io`, `os`).

Three functions called from `main.go`:

#### `downloadELFInstaller(url *url.URL, outdir string) (elfPath string, err error)`

Stream HTTP GET to a temp file. The ELF is ~880 MB; never hold it in memory.

#### `extractACFromELF(elfPath, outdir string) (aceJarPath, networkVersion string, err error)`

1. `zip.OpenReader(elfPath)` opens the self-extracting ELF directly. Go's zip reader
   handles the prepended stub via the End-of-Central-Directory record; no offset
   computation needed.
2. Read the `image.tar` zip entry to a temp file.
3. Open `image.tar` as a tar archive. Find and read `index.json`. Parse as
   `v1.Index` from `image-spec`.
4. Resolve the manifest digest from the index. Read
   `blobs/sha256/<manifest-digest>` and parse as `v1.Manifest`.
5. For each `manifest.Layers[i]`: open `blobs/sha256/<layer.digest>`. If `mediaType`
   indicates gzip, wrap in a `gzip.Reader`. Scan the layer tar for
   `usr/lib/unifi/lib/ace.jar`. Write the match to `outdir/ace.jar`.
6. Open `ace.jar` as a zip. Read `BOOT-INF/classes/product.properties`. Parse the
   `version=` line. Return the ace.jar path and the Network version.

If `index.json` is absent, the image tar is not OCI layout; return a clear error. No
silent fallback to brute-force scanning.

#### `extractInternalDepsJSON(aceJarPath, outdir string) (written []string, err error)`

1. Open `ace.jar` as a zip.
2. Read `BOOT-INF/lib/internal-dependencies.jar` (fall back to searching for any entry
   whose name contains `internal-dependencies`, matching the PoC).
3. Open the inner jar as a zip. For each name in `KEEP_TOPLEVEL` (eight files:
   `legacy_endpoint_segments.json`, `event_defs.json`, `sensitive_metadata.json`,
   `radio_specification.json`, `country_codes_list.json`,
   `geo_ip_country_codes_list.json`, `timezones.json`,
   `ssl-inspection-file-extension.json`), write to `outdir`.
4. Recurse `api/fields/`: write every `.json` entry to `outdir/api/fields/`.
5. Return the list of relative paths written.

### 2. Firmware API (`cmd/fields/fwupdate.go` + `version.go`, EDIT)

#### `fwupdate.go`

Add `unifiOSServerProduct = "unifi-os-server"` and `linuxX64Platform = "linux-x64"`
constants. Existing `firmwareUpdateApiFilter`, `firmwareUpdateApi`, and response types
are reused unchanged.

#### `version.go`

Add `latestUnifiOSVersion() (*version.Version, *url.URL, error)` as a sibling to
`latestUnifiVersion()`. It queries `firmwareUpdateApi` with filters:
`eq~~channel~~release`, `eq~~product~~unifi-os-server`, `eq~~platform~~linux-x64`.
No `maxVersion` cap. Returns `(osVersion, downloadURL)`.

`latestUnifiVersion()` stays unchanged for 9.x legacy support.

### 3. Main Generator (`cmd/fields/main.go`, EDIT)

Routing is based on which flag was used, not the version number (the OS version
5.1.21 is < 10, so a version-number check would misroute `-latest` to `.deb`):

- `-latest` flag: call `latestUnifiOSVersion()` → ELF path. The OS version (5.1.21) is
  only used to construct the initial `*version.Version` for the firmware API response;
  the Network version (10.4.57) is discovered during extraction from
  `product.properties` and used for `UnifiVersion` and the output directory name.
- `-version X.Y.Z` where `X < 10`: construct the `.deb` URL
  (`https://dl.ui.com/unifi/%s/unifi_sysvinit_all.deb`), use the existing `.deb` path
  (`downloadJar` + `extractJSON`). Unchanged.
- `-version X.Y.Z` where `X >= 10`: reject with a helpful error. The firmware API is
  the only source for the ELF URL; the user cannot predict the OS version that
  corresponds to a given Network version. Use `-latest` for 10.x+.

After extraction and before code generation:

- Attempt to load `sensitive_metadata.json` from the version directory via
  `loadSensitiveMetadata()`. If present, call `resource.MarkSensitiveFields(meta)`
  for each resource after `processJSON`. The metadata is consumed only during this
  pre-pass; the spec generator reads `field.Sensitive` during attribute generation and
  does not hold a reference to the metadata itself.
- If absent (9.x path), skip silently.

### 4. Sensitive Metadata Wiring (`cmd/fields/schema.go`, EDIT)

#### Data types

```go
type SensitiveMetadata struct {
    SensitiveDBFieldsByCollection          map[string][]string `json:"sensitive_db_fields_by_collection"`
    SensitiveDistinctDBFieldsByCollection map[string]string   `json:"sensitive_distinct_db_fields_by_collection"`
}
```

#### Loader

`loadSensitiveMetadata(path string) (*SensitiveMetadata, error)` returns `(nil, nil)`
if the file is absent. The 9.x path produces no metadata file; the loader treats this
as a non-error.

#### Collection-to-struct mapping

A `CollectionName()` method on `ResourceInfo` returns the Mongo collection name for a
resource by reverse-applying `fileReps` (struct name → original file name → lowercased).
Examples:

- `networkconf` → `NetworkConf` → struct `Network`
- `wlanconf` → `WlanConf` → struct `WLAN`
- `setting` → split into per-setting resources (`SettingMgmt`, `SettingUsg`, etc.)
- `device` → `Device`
- `account` → `Account`
- `radiusprofile` → `RadiusProfile` → struct `RADIUSProfile`
- `user` → `User` → struct `Client`
- `usergroup` → `UserGroup` → struct `ClientGroup`
- `dpigroup` → `DpiGroup`
- `dynamicdns` → `DynamicDNS`

Special case: all `Setting*` resources map to the `setting` collection (the metadata
file lists sensitive fields for `setting` as a single collection, but the generator
splits settings into per-resource structs).

Unmapped collections (`teleport_client`, `teleport_token`,
`ssl_inspectioncertificate`, `site`, `admin`) have no corresponding resource; log a
warning and skip.

#### Pre-pass marking

Add `Sensitive bool` to `FieldInfo`. After `processJSON` for each resource, call:

```go
func (r *ResourceInfo) MarkSensitiveFields(meta *SensitiveMetadata) {
    collection := r.CollectionName()
    paths := meta.SensitiveDBFieldsByCollection[collection]
    if len(paths) == 0 { return }
    set := make(map[string]struct{}, len(paths))
    for _, p := range paths { set[p] = struct{}{} }
    // Also fold in distinct fields (single-string entries)
    if distinct, ok := meta.SensitiveDistinctDBFieldsByCollection[collection]; ok {
        set[distinct] = struct{}{}
    }
    r.walkAndMark(r.Types[r.StructName], "", set)
}

func (r *ResourceInfo) walkAndMark(field *FieldInfo, parent string, set map[string]struct{}) {
    if field == nil { return }
    path := field.JSONName
    if parent != "" { path = parent + "." + path }
    if _, hit := set[path]; hit { field.Sensitive = true }
    for _, child := range field.Fields {
        r.walkAndMark(child, path, set)
    }
}
```

Path matching uses `JSONName` (not the Go field name). Dotted paths like
`auth_servers.x_secret` work naturally because the walk recurses through
`field.Fields`, building the path as it goes.

#### Attribute generation

In `fieldToResourceAttribute` and `fieldToDataSourceAttribute`, after constructing the
attribute, if `field.Sensitive`:

```go
if field.Sensitive {
    if attr.String != nil  { attr.String.Sensitive = ptr(true) }
    if attr.Bool != nil    { attr.Bool.Sensitive = ptr(true) }
    if attr.Int64 != nil   { attr.Int64.Sensitive = ptr(true) }
    if attr.Float64 != nil { attr.Float64.Sensitive = ptr(true) }
}
```

Container attributes (`ListNested`, `SingleNested`, `List`) have no `Sensitive` field
at the container level in codegen-spec. The pre-pass already marked their leaf
children, so recursive attribute generation picks those up automatically.

### 5. GitHub Actions Automation

#### `generate.yaml` (rewritten)

Daily cron at `0 0 * * *` plus `workflow_dispatch`. Steps:

1. Checkout.
2. Setup Go.
3. `go generate ./...` — runs `cmd/fields -latest`, regenerates `unifi/*.generated.go`,
   `specification.json`, `version.generated.go`.
4. `go test ./...` — gates auto-merge.
5. If `git diff --quiet` (no changes), exit successfully.
6. Create PR via `peter-evans/create-pull-request` with title
   `Update to latest controller version (Network X.Y.Z)`.
7. Enable auto-merge (squash, delete branch) via
   `gh pr merge --auto --squash --delete-branch`.

Required permissions: `contents: write`, `pull-requests: write`.

Branch protection on `main` (manual one-time setup): require CI status checks (lint,
test, generate-idempotent) before merge; allow auto-merge.

#### `tag.yaml` (NEW)

On push to `main`. Steps:

1. Checkout with `fetch-depth: 0`.
2. Read `UnifiVersion` from `unifi/version.generated.go`.
3. Check whether tag `v<unifi-version>` already exists:
   `git rev-parse "refs/tags/v<ver>"`. If yes, exit (idempotent).
4. Create and push the tag: `git tag v<ver> && git push origin v<ver>`.

The tag push triggers the existing `release.yaml`, which runs GoReleaser and publishes
the GitHub release.

#### `release.yaml` (unchanged)

Triggers on `v*` tag push. Runs GoReleaser. Already works.

#### Idempotency

If the firmware API has not changed since the last run, `go generate` produces no diff,
no PR is created, no tag is pushed, no release is published. The pipeline is a no-op.
This is the correct behavior for daily polling.

## Testing

### `cmd/fields/extract_elf_test.go` (NEW)

| Test | Covers |
|---|---|
| `TestDownloadELFInstaller` | HTTP test server serving a fake ELF; verifies temp file creation and content. |
| `TestExtractACFromELF` | Synthetic `image.tar` (one OCI index + one small layer containing a fake `ace.jar` with `product.properties`); wrapped in a zip as `testdata/installer`. Verifies Network version extraction and ace.jar saved. |
| `TestExtractInternalDepsJSON` | Synthetic `internal-dependencies.jar` containing one `api/fields/Foo.json` + `sensitive_metadata.json`; verifies both written. |

### `cmd/fields/version_test.go` (EDIT)

Add `TestLatestUnifiOSVersion`, mirroring the existing `TestLatestUnifiVersion` but
for the `unifi-os-server`/`linux-x64` query path.

### `cmd/fields/schema_test.go` (EDIT)

| Test | Covers |
|---|---|
| `TestMarkSensitiveFields` | Synthetic `ResourceInfo` with nested fields; verifies simple-name matches and dotted-path matches (`auth_servers.x_secret`) both get `Sensitive=true`; non-listed fields stay `false`. |
| `TestFieldToResourceAttribute_Sensitive` | After marking, verifies `Sensitive: ptr(true)` set on the resource attribute. |

### Fixtures

Built programmatically in `_test.go` (assemble zips and tars in Go), not committed as
binaries. A small `testdata/` directory holds human-readable expectations. The
`testdata/installer` is a bare zip containing `image.tar` — no ELF stub needed because
`zip.OpenReader` handles prepended bytes via the EOCD record.

### Error Handling

- Network and download failures: surface with `fmt.Errorf("...: %w", err)`. No
  swallowing.
- Missing `index.json`: clear error
  `"image.tar is not an OCI image layout (missing index.json)"`.
- `ace.jar` not found in any layer: clear error
  `"ace.jar not found in any image layer"`.
- `internal-dependencies.jar` missing from `ace.jar`: clear error.
- `sensitive_metadata.json` absent (9.x path): silent. Loader returns `(nil, nil)`.
- Unmapped collection in sensitive metadata: `fmt.Printf` warning, skip.
- Network version unparseable from `product.properties`: fall back to `"unknown"`
  (matching the PoC behavior). Should not happen in practice.

## Dependencies

New Go module dependencies:

- `github.com/opencontainers/image-spec` — OCI image manifest and index types.
- `github.com/opencontainers/go-digest` — digest string construction for
  `blobs/sha256/<digest>` path lookup.

Both are lightweight OCI spec libraries with minimal transitive dependencies. Add via
`go get` and update `go.mod`/`go.sum`.

## File Deliverables

| File | Status | Purpose |
|---|---|---|
| `cmd/fields/extract_elf.go` | NEW | ELF download, OCI parsing, ace.jar extraction, internal-deps JSON extraction. |
| `cmd/fields/extract_elf_test.go` | NEW | Tests and programmatic fixture builders. |
| `cmd/fields/fwupdate.go` | EDIT | Add `unifi-os-server` and `linux-x64` constants. |
| `cmd/fields/version.go` | EDIT | Add `latestUnifiOSVersion()`. |
| `cmd/fields/version_test.go` | EDIT | Add `TestLatestUnifiOSVersion`. |
| `cmd/fields/main.go` | EDIT | Version-based fetcher selection; call `MarkSensitiveFields`; set `UnifiVersion` from Network version for 10.x; reject explicit 10.x version. |
| `cmd/fields/schema.go` | EDIT | `SensitiveMetadata` struct + loader; `MarkSensitiveFields`; `FieldInfo.Sensitive`; `CollectionName()`; set `Sensitive: true` on leaf attributes. |
| `cmd/fields/schema_test.go` | EDIT | `TestMarkSensitiveFields`, `TestFieldToResourceAttribute_Sensitive`. |
| `go.mod`, `go.sum` | EDIT | Add `image-spec` + `go-digest`. |
| `.github/workflows/generate.yaml` | EDIT | Add test gate, auto-merge, permissions. |
| `.github/workflows/tag.yaml` | NEW | Tag `v<unifi-version>` on push to `main` if absent. |
| `unifi/version.generated.go` | REGENERATED | Reads `UnifiVersion = "10.4.57"` after first run. |
| `specification.json` | REGENERATED | Gains `sensitive: true` on matching fields. |

## Out of Scope

- Wiring `radio_specification.json`, `country_codes_list.json`, `timezones.json`, or
  other metadata files into the spec generator.
- Modifying the `.deb` extractor internals beyond the version-based branch.
- Branch protection and repo settings via workflow (manual one-time setup).
- Changing `cmd/fields/custom/*.json` v2 overrides (BgpConfig, FirewallPolicy, etc.).
