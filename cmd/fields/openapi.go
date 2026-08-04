package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/ubiquiti-community/go-unifi/internal/fields"
)

//go:embed openapi.tmpl
var openapiTemplate string

// OpenAPIVersion is the OpenAPI document version emitted. The extended Tag
// Object fields (summary/parent/kind) were standardized in OpenAPI 3.2.0, so we
// target that release; tooling that only understands 3.1 still reads the
// name/description/externalDocs subset.
const OpenAPIVersion = "3.2.0"

// OpenAPIData is the root object handed to openapi.tmpl.
type OpenAPIData struct {
	Version   string
	OpenAPI   string
	Resources []*ResourceInfo
	Tags      []Tag
}

// Tag is an OpenAPI 3.2 Tag Object, including the hierarchical/extended fields.
type Tag struct {
	Name        string // matches the tag referenced by each operation
	Summary     string // short, human-readable label
	Description string // longer prose
	Parent      string // name of the parent tag ("" for top-level)
	Kind        string // machine-readable grouping (firewall, clients, ...)
}

// acronyms keeps these tokens upper-cased when humanizing struct names.
var acronyms = map[string]bool{
	"WLAN": true, "AP": true, "DNS": true, "DHCP": true, "VPN": true, "NAT": true,
	"RADIUS": true, "IP": true, "QOS": true, "BGP": true, "OSPF": true, "ID": true,
	"MAC": true, "UTC": true, "SSID": true, "DPI": true, "IGMP": true, "ICMP": true,
	"WAN": true, "LAN": true, "WEP": true, "WPA": true, "PMF": true, "TFTP": true,
}

var (
	reAcronymBoundary = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	reWordBoundary    = regexp.MustCompile(`([a-z0-9])([A-Z])`)
)

// humanize turns a CamelCase struct name into spaced words while keeping known
// acronyms intact: "WLANGroup" -> "WLAN Group", "ClientGroup" -> "Client Group".
func humanize(name string) string {
	s := reAcronymBoundary.ReplaceAllString(name, "$1 $2")
	s = reWordBoundary.ReplaceAllString(s, "$1 $2")
	parts := strings.Fields(s)
	for i, p := range parts {
		if acronyms[strings.ToUpper(p)] {
			parts[i] = strings.ToUpper(p)
		}
	}
	return strings.Join(parts, " ")
}

// explicitParents covers hierarchies the suffix heuristic cannot infer.
var explicitParents = map[string]string{
	"ClientGroup": "Client",
	"WLANGroup":   "WLAN",
	"APGroup":     "Device",
}

// childSuffixes are stripped to discover a parent ("FooProfile" -> "Foo").
var childSuffixes = []string{"Group", "Profile", "Override", "Schedule", "Policy"}

// tagParent resolves the parent tag for a resource. It only returns a parent
// that actually exists in the tag set (names) so references never dangle.
func tagParent(r *ResourceInfo, names map[string]bool) string {
	if p, ok := explicitParents[r.StructName]; ok && names[p] {
		return p
	}
	for _, suf := range childSuffixes {
		if base := strings.TrimSuffix(r.StructName, suf); base != r.StructName && names[base] {
			return base
		}
	}
	return ""
}

// tagKind assigns a coarse, machine-readable domain grouping to a tag.
func tagKind(r *ResourceInfo) string {
	if r.IsSetting() {
		return "settings"
	}
	switch r.StructName {
	case "PortForward":
		return "network"
	case "RADIUSProfile":
		return "security"
	}
	n := r.StructName
	switch {
	case containsAny(n, "Firewall", "Nat", "Acl"):
		return "firewall"
	case containsAny(n, "Client", "User", "Account"):
		return "clients"
	case containsAny(n, "WLAN", "Hotspot", "Guest", "SSID", "Ssid", "APGroup", "ApGroup", "ChannelPlan", "Radio"):
		return "wireless"
	case containsAny(n, "Device", "Port"):
		return "devices"
	case containsAny(n, "Route", "BGP", "OSPF", "VPN", "Vpn", "DynamicDNS", "Dynamicdns"):
		return "routing"
	case containsAny(n, "Network", "DHCP", "DNS"):
		return "network"
	default:
		return "resource"
	}
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// buildTags produces a sorted, hierarchical tag list from the resources.
func buildTags(resources []*ResourceInfo) []Tag {
	names := make(map[string]bool, len(resources))
	for _, r := range resources {
		names[r.StructName] = true
	}

	tags := make([]Tag, 0, len(resources))
	for _, r := range resources {
		summary := humanize(r.CleanStructName())
		kind := tagKind(r)
		var desc string
		if r.IsSetting() {
			desc = fmt.Sprintf("Configure the UniFi %s settings.", summary)
		} else {
			desc = fmt.Sprintf("Manage UniFi %s objects (%s) via the controller API.", summary, kind)
		}
		tags = append(tags, Tag{
			Name:        r.StructName,
			Summary:     summary,
			Description: desc,
			Parent:      tagParent(r, names),
			Kind:        kind,
		})
	}
	sort.Slice(tags, func(i, j int) bool { return tags[i].Name < tags[j].Name })
	return tags
}

// primitiveOpenAPIType maps a Go field type (as stored on FieldInfo.FieldType)
// to an OpenAPI primitive type. Anything not listed here is treated as a
// reference to another component schema (see isRefType).
var primitiveOpenAPIType = map[string]string{
	fields.String: "string",
	fields.Int:    "integer", // "int64"
	fields.Bool:   "boolean",
	fields.Number: "number",  // "types.Number"
	"float64":     "number",
	"int":         "integer",
}

// baseType strips a leading "[]" slice marker (some synthetic fields, e.g.
// Device.PortTable, encode the array directly in the type string).
func baseType(fieldType string) string {
	return strings.TrimPrefix(fieldType, "[]")
}

// isArrayField reports whether a field renders as an OpenAPI array, covering
// both the FieldInfo.IsArray flag and the "[]Type" inline-slice convention.
func isArrayField(f *FieldInfo) bool {
	return f.IsArray || strings.HasPrefix(f.FieldType, "[]")
}

// isRefType reports whether a (possibly slice-stripped) field type refers to
// another component schema rather than an OpenAPI primitive.
func isRefType(fieldType string) bool {
	_, primitive := primitiveOpenAPIType[baseType(fieldType)]
	return !primitive
}

// openapiType returns the OpenAPI primitive type for a field type.
func openapiType(fieldType string) string {
	if t, ok := primitiveOpenAPIType[baseType(fieldType)]; ok {
		return t
	}
	return "object"
}

// openapiFormat returns the OpenAPI "format" hint for a primitive, if any.
func openapiFormat(fieldType string) string {
	switch baseType(fieldType) {
	case fields.Int, "int":
		return "int64"
	case "float64":
		return "double"
	default:
		return ""
	}
}

// requiredFields returns the JSON names of a type's required properties, sorted
// for deterministic output. A field is considered required when it is always
// serialized (i.e. not omitempty); server-assigned/optional fields are omitted.
func requiredFields(t *FieldInfo) []string {
	var req []string
	for _, f := range t.Fields {
		if f == nil || f.OmitEmpty {
			continue
		}
		req = append(req, f.JSONName)
	}
	sort.Strings(req)
	return req
}

// collectionPath returns the REST collection path (with a {site} placeholder)
// for a resource, matching the URLs the generated Go client uses.
func collectionPath(r *ResourceInfo) string {
	switch {
	case r.IsV2():
		return "/v2/api/site/{site}/" + r.ResourcePath
	case r.IsDevice():
		return "/api/s/{site}/stat/" + r.ResourcePath
	default:
		return "/api/s/{site}/rest/" + r.ResourcePath
	}
}

// itemPath returns the single-object path for a resource.
func itemPath(r *ResourceInfo) string {
	return collectionPath(r) + "/{id}"
}

func openapiFuncMap() template.FuncMap {
	return template.FuncMap{
		"trimPrefix":     strings.TrimPrefix,
		"hasPrefix":      strings.HasPrefix,
		"lower":          strings.ToLower,
		"baseType":       baseType,
		"isArrayField":   isArrayField,
		"isRef":          isRefType,
		"openapiType":    openapiType,
		"openapiFormat":  openapiFormat,
		"requiredFields": requiredFields,
		"collectionPath": collectionPath,
		"itemPath":       itemPath,
	}
}

// WriteOpenAPI renders an OpenAPI 3.0.3 document for every resource directly
// from the field-validation template variables and writes it to outPath.
func WriteOpenAPI(resources []*ResourceInfo, version, outPath string) error {
	// Deterministic ordering by struct name.
	sorted := make([]*ResourceInfo, len(resources))
	copy(sorted, resources)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].StructName < sorted[j].StructName
	})

	tpl, err := template.New("openapi.tmpl").Funcs(openapiFuncMap()).Parse(openapiTemplate)
	if err != nil {
		return fmt.Errorf("parsing openapi template: %w", err)
	}

	var buf bytes.Buffer
	data := OpenAPIData{
		Version:   version,
		OpenAPI:   OpenAPIVersion,
		Resources: sorted,
		Tags:      buildTags(sorted),
	}
	if err := tpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("rendering openapi template: %w", err)
	}

	if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", outPath, err)
	}
	return nil
}
