package unifi

import (
	"context"
	"fmt"
)

//go:generate go tool golang.org/x/tools/cmd/stringer -trimprefix DeviceState -type DeviceState
type DeviceState int

const (
	DeviceStateUnknown          DeviceState = 0
	DeviceStateConnected        DeviceState = 1
	DeviceStatePending          DeviceState = 2
	DeviceStateFirmwareMismatch DeviceState = 3
	DeviceStateUpgrading        DeviceState = 4
	DeviceStateProvisioning     DeviceState = 5
	DeviceStateHeartbeatMissed  DeviceState = 6
	DeviceStateAdopting         DeviceState = 7
	DeviceStateDeleting         DeviceState = 8
	DeviceStateInformError      DeviceState = 9
	DeviceStateAdoptFailed      DeviceState = 10
	DeviceStateIsolated         DeviceState = 11
)

type DeviceLastConnection struct {
	MAC      string `json:"mac,omitempty"`
	LastSeen int    `json:"last_seen,omitempty"`
}
type DevicePortTable struct {
	PortIdx             int                  `json:"port_idx"`
	Media               string               `json:"media"`
	PortPoe             bool                 `json:"port_poe"`
	PoeCaps             int                  `json:"poe_caps"`
	SpeedCaps           int                  `json:"speed_caps"`
	LastConnection      DeviceLastConnection `json:"last_connection"`
	OpMode              string               `json:"op_mode"`
	Forward             string               `json:"forward"`
	PoeMode             string               `json:"poe_mode"`
	Anomalies           int                  `json:"anomalies"`
	Autoneg             bool                 `json:"autoneg"`
	Dot1XMode           string               `json:"dot1x_mode"`
	Dot1XStatus         string               `json:"dot1x_status"`
	Enable              bool                 `json:"enable"`
	FlowctrlRx          bool                 `json:"flowctrl_rx"`
	FlowctrlTx          bool                 `json:"flowctrl_tx"`
	FullDuplex          bool                 `json:"full_duplex"`
	IsUplink            bool                 `json:"is_uplink"`
	Jumbo               bool                 `json:"jumbo"`
	MacTableCount       int                  `json:"mac_table_count"`
	PoeClass            string               `json:"poe_class"`
	PoeCurrent          string               `json:"poe_current"`
	PoeEnable           bool                 `json:"poe_enable"`
	PoeGood             bool                 `json:"poe_good"`
	PoePower            string               `json:"poe_power,omitempty"`
	PoeVoltage          string               `json:"poe_voltage"`
	RxBroadcast         int                  `json:"rx_broadcast"`
	RxBytes             int                  `json:"rx_bytes"`
	RxDropped           int                  `json:"rx_dropped"`
	RxErrors            int                  `json:"rx_errors"`
	RxMulticast         int                  `json:"rx_multicast"`
	RxPackets           int                  `json:"rx_packets"`
	Satisfaction        int                  `json:"satisfaction"`
	SatisfactionReason  int                  `json:"satisfaction_reason"`
	Speed               int                  `json:"speed"`
	StpPathcost         int                  `json:"stp_pathcost"`
	StpState            string               `json:"stp_state"`
	TxBroadcast         int                  `json:"tx_broadcast"`
	TxBytes             int64                `json:"tx_bytes"`
	TxDropped           int                  `json:"tx_dropped"`
	TxErrors            int                  `json:"tx_errors"`
	TxMulticast         int                  `json:"tx_multicast"`
	TxPackets           int                  `json:"tx_packets"`
	Up                  bool                 `json:"up"`
	TxBytesR            float64              `json:"tx_bytes-r"`
	RxBytesR            float64              `json:"rx_bytes-r"`
	BytesR              float64              `json:"bytes-r"`
	FlowControlEnabled  bool                 `json:"flow_control_enabled"`
	NativeNetworkconfID string               `json:"native_networkconf_id"`
	Name                string               `json:"name"`
	SettingPreference   string               `json:"setting_preference"`
	StormctrlBcastRate  int                  `json:"stormctrl_bcast_rate"`
	StormctrlMcastRate  int                  `json:"stormctrl_mcast_rate"`
	StormctrlUcastRate  int                  `json:"stormctrl_ucast_rate"`
	TaggedVlanMgmt      string               `json:"tagged_vlan_mgmt"`
	Masked              bool                 `json:"masked"`
	AggregatedBy        bool                 `json:"aggregated_by"`
}

func (c *Client) ListDevice(ctx context.Context, site string) ([]Device, error) {
	return c.listDevice(ctx, site)
}

func (c *Client) GetDeviceByMAC(ctx context.Context, site, mac string) (*Device, error) {
	return c.getDevice(ctx, site, mac)
}

func (c *Client) DeleteDevice(ctx context.Context, site, id string) error {
	return c.deleteDevice(ctx, site, id)
}

func (c *Client) CreateDevice(ctx context.Context, site string, d *Device) (*Device, error) {
	return c.createDevice(ctx, site, d)
}

func (c *Client) UpdateDevice(ctx context.Context, site string, d *Device) (*Device, error) {
	return c.updateDevice(ctx, site, d)
}

func (c *Client) GetDevice(ctx context.Context, site, id string) (*Device, error) {
	devices, err := c.ListDevice(ctx, site)
	if err != nil {
		return nil, err
	}

	for _, d := range devices {
		if d.ID == id {
			return &d, nil
		}
	}

	return nil, &NotFoundError{}
}

func (c *Client) AdoptDevice(ctx context.Context, site, mac string) error {
	reqBody := struct {
		Cmd string `json:"cmd"`
		MAC string `json:"mac"`
	}{
		Cmd: "adopt",
		MAC: mac,
	}

	var respBody struct {
		Meta meta `json:"meta"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("api/s/%s/cmd/devmgr", site), reqBody, &respBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ForgetDevice(ctx context.Context, site, mac string) error {
	reqBody := struct {
		Cmd  string   `json:"cmd"`
		MACs []string `json:"macs"`
	}{
		Cmd:  "delete-device",
		MACs: []string{mac},
	}

	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Device `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("api/s/%s/cmd/sitemgr", site), reqBody, &respBody)
	if err != nil {
		return err
	}

	return nil
}
