# Unifi Go SDK [![GoDoc](https://godoc.org/github.com/ubiquiti-community/go-unifi?status.svg)](https://godoc.org/github.com/ubiquiti-community/go-unifi)

This was written primarily for use in my [Terraform provider for Unifi](https://github.com/ubiquiti-community/terraform-provider-unifi).

## Versioning

Many of the naming adjustments are breaking changes, but to simplify things, treating naming errors as minor changes for the 1.0.0 version (probably should have just started at 0.1.0).

## Note on Code Generation

The data models and basic REST methods are "generated" from OpenAPI Specification JSON files.

To regenerate the code, run `go generate` inside the `unifi` directory.
