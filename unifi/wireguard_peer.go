package unifi

import (
	"context"
	"fmt"
	"net/http"
)

// This is a v2 API object, manually coded.
//
// WireGuard peers ("clients" in the UI) of a WireGuard server network
// (vpn_type=wireguard-server). The controller only exposes batch
// create/update/delete endpoints, so single-peer CRUD wraps those.

type WireGuardPeer struct {
	ID        string `json:"_id,omitempty"`
	NetworkID string `json:"network_id,omitempty"`

	Name        string   `json:"name"`
	InterfaceIP string   `json:"interface_ip"`
	PublicKey   string   `json:"public_key"`
	AllowedIPs  []string `json:"allowed_ips"`
}

func (c *ApiClient) wireGuardPeersPath(site, networkID string) string {
	return fmt.Sprintf("v2/api/site/%s/wireguard/%s/users", site, networkID)
}

func (c *ApiClient) ListWireGuardPeers(ctx context.Context, site, networkID string) ([]WireGuardPeer, error) {
	var respBody []WireGuardPeer

	err := c.do(ctx, http.MethodGet, c.wireGuardPeersPath(site, networkID), nil, &respBody,
		map[string]string{"networkId": networkID})
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// GetWireGuardPeer fetches a single peer by ID. The controller has no
// single-peer endpoint, so this lists and filters.
func (c *ApiClient) GetWireGuardPeer(ctx context.Context, site, networkID, id string) (*WireGuardPeer, error) {
	peers, err := c.ListWireGuardPeers(ctx, site, networkID)
	if err != nil {
		return nil, err
	}

	for i := range peers {
		if peers[i].ID == id {
			return &peers[i], nil
		}
	}

	return nil, &NotFoundError{Type: "wireguard_peer", Attr: "_id", Value: id}
}

func (c *ApiClient) CreateWireGuardPeer(ctx context.Context, site, networkID string, d *WireGuardPeer) (*WireGuardPeer, error) {
	var respBody []WireGuardPeer
	d.ID = ""
	if d.AllowedIPs == nil {
		d.AllowedIPs = []string{}
	}

	err := c.do(ctx, http.MethodPost, c.wireGuardPeersPath(site, networkID)+"/batch", []*WireGuardPeer{d}, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody) != 1 {
		return nil, &NotFoundError{Type: "wireguard_peer"}
	}

	return &respBody[0], nil
}

func (c *ApiClient) UpdateWireGuardPeer(ctx context.Context, site, networkID string, d *WireGuardPeer) (*WireGuardPeer, error) {
	var respBody []WireGuardPeer
	if d.AllowedIPs == nil {
		d.AllowedIPs = []string{}
	}

	err := c.do(ctx, http.MethodPut, c.wireGuardPeersPath(site, networkID)+"/batch", []*WireGuardPeer{d}, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody) != 1 {
		return nil, &NotFoundError{Type: "wireguard_peer", Attr: "_id", Value: d.ID}
	}

	return &respBody[0], nil
}

func (c *ApiClient) DeleteWireGuardPeer(ctx context.Context, site, networkID, id string) error {
	return c.do(ctx, http.MethodPost, c.wireGuardPeersPath(site, networkID)+"/batch_delete", []string{id}, nil)
}
