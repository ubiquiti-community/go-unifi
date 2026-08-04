package main

import (
	"os"
	"strings"
	"testing"

	"github.com/ubiquiti-community/go-unifi/internal/fields"
	"gopkg.in/yaml.v3"
)

// sampleResources builds a couple of resources by hand (covering primitives,
// arrays, pointers, nested ref types, validation, v1 and v2) so the template can
// be exercised without downloading the controller jar.
func sampleResources() []*ResourceInfo {
	// A v1 REST resource with a nested type.
	net := NewResource("Network", "networkconf")
	base := net.Types["Network"]
	base.Fields["Name"] = NewFieldInfo("Name", "name", fields.String, "", false, false, false, "")
	base.Fields["VLAN"] = NewFieldInfo("VLAN", "vlan", fields.Int, "[0-9]+", true, false, true, fields.Number)
	base.Fields["Enabled"] = NewFieldInfo("Enabled", "enabled", fields.Bool, "", true, false, true, "")
	base.Fields["DNS"] = NewFieldInfo("DNS", "dns", fields.String, "", true, true, false, "")
	base.Fields["Gateway"] = NewFieldInfo("Gateway", "gateway", "NetworkGateway", "", true, false, false, "")

	gw := NewFieldInfo("NetworkGateway", "gateway", "struct", "", false, false, false, "")
	gw.Fields = map[string]*FieldInfo{
		"IP": NewFieldInfo("IP", "ip", fields.String, "", false, false, false, ""),
	}
	net.Types["NetworkGateway"] = gw

	// A v2 resource (FirewallPolicy is in the IsV2 allow-list).
	fp := NewResource("FirewallPolicy", "firewall-policies")
	fp.Types["FirewallPolicy"].Fields["Action"] = NewFieldInfo("Action", "action", fields.String, "allow|deny", false, false, false, "")

	return []*ResourceInfo{net, fp}
}

func TestOpenAPITemplateRendersValidYAML(t *testing.T) {
	dir := t.TempDir()
	out := dir + "/openapi.yaml"
	if err := WriteOpenAPI(sampleResources(), "9.9.9", out); err != nil {
		t.Fatalf("WriteOpenAPI: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	// 1. The document must be syntactically valid YAML.
	var doc map[string]any
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("generated OpenAPI is not valid YAML: %v\n---\n%s", err, data)
	}

	// 2. Top-level OpenAPI structure.
	if doc["openapi"] != OpenAPIVersion {
		t.Errorf("openapi version = %v, want %s", doc["openapi"], OpenAPIVersion)
	}

	// 2b. Hierarchical tags: ClientGroup -> parent Client, with kind/summary.
	clientGroup := NewResource("ClientGroup", "usergroup")
	tags := buildTags([]*ResourceInfo{
		NewResource("Client", "user"), clientGroup,
		NewResource("Network", "networkconf"),
	})
	byName := map[string]Tag{}
	for _, tg := range tags {
		byName[tg.Name] = tg
	}
	if got := byName["ClientGroup"].Parent; got != "Client" {
		t.Errorf("ClientGroup parent = %q, want Client", got)
	}
	if byName["Client"].Parent != "" {
		t.Errorf("Client should be top-level, got parent %q", byName["Client"].Parent)
	}
	if byName["ClientGroup"].Kind != "clients" {
		t.Errorf("ClientGroup kind = %q, want clients", byName["ClientGroup"].Kind)
	}
	if byName["ClientGroup"].Summary != "Client Group" {
		t.Errorf("ClientGroup summary = %q, want \"Client Group\"", byName["ClientGroup"].Summary)
	}
	comps, _ := doc["components"].(map[string]any)
	schemas, _ := comps["schemas"].(map[string]any)
	if schemas == nil {
		t.Fatalf("components.schemas missing")
	}
	for _, want := range []string{"Network", "NetworkGateway", "FirewallPolicy", "Meta"} {
		if _, ok := schemas[want]; !ok {
			t.Errorf("missing component schema %q", want)
		}
	}

	// 3. Network schema specifics: primitive type, format, ref, array, required.
	netSchema, _ := schemas["Network"].(map[string]any)
	props, _ := netSchema["properties"].(map[string]any)
	if props["name"].(map[string]any)["type"] != "string" {
		t.Errorf("Network.name should be string, got %v", props["name"])
	}
	if props["vlan"].(map[string]any)["format"] != "int64" {
		t.Errorf("Network.vlan should have format int64, got %v", props["vlan"])
	}
	if props["dns"].(map[string]any)["type"] != "array" {
		t.Errorf("Network.dns should be array, got %v", props["dns"])
	}
	gwRef, _ := props["gateway"].(map[string]any)
	if gwRef["$ref"] != "#/components/schemas/NetworkGateway" {
		t.Errorf("Network.gateway should $ref NetworkGateway, got %v", gwRef)
	}
	// "name" is not omitempty -> required; "enabled" is omitempty -> not.
	req := toStringSlice(netSchema["required"])
	if !contains(req, "name") {
		t.Errorf("Network.required should include name, got %v", req)
	}
	if contains(req, "enabled") {
		t.Errorf("Network.required should not include enabled, got %v", req)
	}

	// 4. Paths: v1 uses /rest/, v2 uses /v2/api/site/. Both expose CRUD verbs.
	paths, _ := doc["paths"].(map[string]any)
	if _, ok := paths["/api/s/{site}/rest/networkconf"]; !ok {
		t.Errorf("missing v1 collection path; have %v", keys(paths))
	}
	if _, ok := paths["/v2/api/site/{site}/firewall-policies"]; !ok {
		t.Errorf("missing v2 collection path; have %v", keys(paths))
	}
	item, _ := paths["/api/s/{site}/rest/networkconf/{id}"].(map[string]any)
	for _, verb := range []string{"get", "put", "delete"} {
		if _, ok := item[verb]; !ok {
			t.Errorf("network item path missing verb %q", verb)
		}
	}
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

func toStringSlice(v any) []string {
	arr, _ := v.([]any)
	out := make([]string, 0, len(arr))
	for _, x := range arr {
		out = append(out, strings.TrimSpace(x.(string)))
	}
	return out
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
