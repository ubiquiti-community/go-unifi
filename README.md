# Unifi Go SDK [![GoDoc](https://godoc.org/github.com/ubiquiti-community/go-unifi?status.svg)](https://godoc.org/github.com/ubiquiti-community/go-unifi)

This was written primarily for use in my [Terraform provider for Unifi](https://github.com/ubiquiti-community/terraform-provider-unifi).

## Versioning

Many of the naming adjustments are breaking changes, but to simplify things, treating naming errors as minor changes for the 1.0.0 version (probably should have just started at 0.1.0).

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
