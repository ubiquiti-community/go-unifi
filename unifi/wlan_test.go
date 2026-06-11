package unifi_test

import (
	"encoding/json"
	"testing"

	"github.com/ubiquiti-community/go-unifi/unifi"
)

func TestWLANMarshalJSONOmitEmptyWLANGroupID(t *testing.T) {
	actual, err := json.Marshal(&unifi.WLAN{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}

	if string(actual) == "{}" {
		t.Fatal("expected non-empty JSON payload")
	}

	var payload map[string]any
	if err := json.Unmarshal(actual, &payload); err != nil {
		t.Fatal(err)
	}

	if _, ok := payload["wlangroup_id"]; ok {
		t.Fatalf("wlangroup_id was serialized for empty value: %s", actual)
	}
}

func TestWLANMarshalJSONIncludesWLANGroupID(t *testing.T) {
	actual, err := json.Marshal(&unifi.WLAN{Name: "test", WLANGroupID: "default"})
	if err != nil {
		t.Fatal(err)
	}

	var payload map[string]any
	if err := json.Unmarshal(actual, &payload); err != nil {
		t.Fatal(err)
	}

	if payload["wlangroup_id"] != "default" {
		t.Fatalf("wlangroup_id = %#v, want %q; payload: %s", payload["wlangroup_id"], "default", actual)
	}
}
