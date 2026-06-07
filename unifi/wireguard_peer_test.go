package unifi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubiquiti-community/go-unifi/unifi"
)

func TestWireGuardPeerMarshalJSON(t *testing.T) {
	peer := unifi.WireGuardPeer{
		Name:        "test-peer",
		InterfaceIP: "192.0.2.10",
		PublicKey:   "ZmFrZS10ZXN0LXdpcmVndWFyZC1wdWJrZXkAAAAAAAA=",
		AllowedIPs:  []string{},
	}

	actual, err := json.Marshal(&peer)
	if err != nil {
		t.Fatal(err)
	}
	assert.JSONEq(t,
		`{"name":"test-peer","interface_ip":"192.0.2.10","public_key":"ZmFrZS10ZXN0LXdpcmVndWFyZC1wdWJrZXkAAAAAAAA=","allowed_ips":[]}`,
		string(actual))
}

func TestWireGuardPeerUnmarshalJSON(t *testing.T) {
	// shape returned by the controller (v2 wireguard users API)
	raw := `{
		"_id": "000000000000000000000001",
		"network_id": "000000000000000000000002",
		"name": "test-peer",
		"interface_ip": "192.0.2.10",
		"public_key": "ZmFrZS10ZXN0LXdpcmVndWFyZC1wdWJrZXkAAAAAAAA=",
		"allowed_ips": ["198.51.100.0/24"]
	}`

	var peer unifi.WireGuardPeer
	if err := json.Unmarshal([]byte(raw), &peer); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "000000000000000000000001", peer.ID)
	assert.Equal(t, "000000000000000000000002", peer.NetworkID)
	assert.Equal(t, "test-peer", peer.Name)
	assert.Equal(t, "192.0.2.10", peer.InterfaceIP)
	assert.Equal(t, "ZmFrZS10ZXN0LXdpcmVndWFyZC1wdWJrZXkAAAAAAAA=", peer.PublicKey)
	assert.Equal(t, []string{"198.51.100.0/24"}, peer.AllowedIPs)
}
