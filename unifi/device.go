package unifi

import (
	"context"
	"encoding/json"
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
	PortIdx             int                  `json:"port_idx,omitempty"`
	Media               string               `json:"media,omitempty"`
	PortPoe             bool                 `json:"port_poe,omitempty"`
	PoeCaps             int                  `json:"poe_caps,omitempty"`
	SpeedCaps           int                  `json:"speed_caps,omitempty"`
	LastConnection      DeviceLastConnection `json:"last_connection,omitempty"`
	OpMode              string               `json:"op_mode,omitempty"`
	Forward             string               `json:"forward,omitempty"`
	PoeMode             string               `json:"poe_mode,omitempty"`
	Anomalies           int                  `json:"anomalies,omitempty"`
	Autoneg             bool                 `json:"autoneg,omitempty"`
	Dot1XMode           string               `json:"dot1x_mode,omitempty"`
	Dot1XStatus         string               `json:"dot1x_status,omitempty"`
	Enable              bool                 `json:"enable,omitempty"`
	FlowctrlRx          bool                 `json:"flowctrl_rx,omitempty"`
	FlowctrlTx          bool                 `json:"flowctrl_tx,omitempty"`
	FullDuplex          bool                 `json:"full_duplex,omitempty"`
	IsUplink            bool                 `json:"is_uplink,omitempty"`
	Jumbo               bool                 `json:"jumbo,omitempty"`
	MacTableCount       int                  `json:"mac_table_count,omitempty"`
	PoeClass            string               `json:"poe_class,omitempty"`
	PoeCurrent          string               `json:"poe_current,omitempty"`
	PoeEnable           bool                 `json:"poe_enable,omitempty"`
	PoeGood             bool                 `json:"poe_good,omitempty"`
	PoePower            string               `json:"poe_power,omitempty"`
	PoeVoltage          string               `json:"poe_voltage,omitempty"`
	RxBroadcast         int                  `json:"rx_broadcast,omitempty"`
	RxBytes             int                  `json:"rx_bytes,omitempty"`
	RxDropped           int                  `json:"rx_dropped,omitempty"`
	RxErrors            int                  `json:"rx_errors,omitempty"`
	RxMulticast         int                  `json:"rx_multicast,omitempty"`
	RxPackets           int                  `json:"rx_packets,omitempty"`
	Satisfaction        int                  `json:"satisfaction,omitempty"`
	SatisfactionReason  int                  `json:"satisfaction_reason,omitempty"`
	Speed               int                  `json:"speed,omitempty"`
	StpPathcost         int                  `json:"stp_pathcost,omitempty"`
	StpState            string               `json:"stp_state,omitempty"`
	TxBroadcast         int                  `json:"tx_broadcast,omitempty"`
	TxBytes             int64                `json:"tx_bytes,omitempty"`
	TxDropped           int                  `json:"tx_dropped,omitempty"`
	TxErrors            int                  `json:"tx_errors,omitempty"`
	TxMulticast         int                  `json:"tx_multicast,omitempty"`
	TxPackets           int                  `json:"tx_packets,omitempty"`
	Up                  bool                 `json:"up,omitempty"`
	TxBytesR            float64              `json:"tx_bytes-r,omitempty"`
	RxBytesR            float64              `json:"rx_bytes-r,omitempty"`
	BytesR              float64              `json:"bytes-r,omitempty"`
	FlowControlEnabled  bool                 `json:"flow_control_enabled,omitempty"`
	NativeNetworkconfID string               `json:"native_networkconf_id,omitempty"`
	Name                string               `json:"name,omitempty"`
	SettingPreference   string               `json:"setting_preference,omitempty"`
	StormctrlBcastRate  int                  `json:"stormctrl_bcast_rate,omitempty"`
	StormctrlMcastRate  int                  `json:"stormctrl_mcast_rate,omitempty"`
	StormctrlUcastRate  int                  `json:"stormctrl_ucast_rate,omitempty"`
	TaggedVlanMgmt      string               `json:"tagged_vlan_mgmt,omitempty"`
	Masked              bool                 `json:"masked,omitempty"`
	AggregatedBy        bool                 `json:"aggregated_by,omitempty"`
}

func (dst *DevicePortTable) UnmarshalJSON(b []byte) error {
	type Alias struct {
		PortIdx             emptyStringInt       `json:"port_idx,omitempty"`
		Media               string               `json:"media,omitempty"`
		PortPoe             bool                 `json:"port_poe,omitempty"`
		PoeCaps             emptyStringInt       `json:"poe_caps,omitempty"`
		SpeedCaps           emptyStringInt       `json:"speed_caps,omitempty"`
		LastConnection      DeviceLastConnection `json:"last_connection,omitempty"`
		OpMode              string               `json:"op_mode,omitempty"`
		Forward             string               `json:"forward,omitempty"`
		PoeMode             string               `json:"poe_mode,omitempty"`
		Anomalies           emptyStringInt       `json:"anomalies,omitempty"`
		Autoneg             bool                 `json:"autoneg,omitempty"`
		Dot1XMode           string               `json:"dot1x_mode,omitempty"`
		Dot1XStatus         string               `json:"dot1x_status,omitempty"`
		Enable              bool                 `json:"enable,omitempty"`
		FlowctrlRx          bool                 `json:"flowctrl_rx,omitempty"`
		FlowctrlTx          bool                 `json:"flowctrl_tx,omitempty"`
		FullDuplex          bool                 `json:"full_duplex,omitempty"`
		IsUplink            bool                 `json:"is_uplink,omitempty"`
		Jumbo               bool                 `json:"jumbo,omitempty"`
		MacTableCount       emptyStringInt       `json:"mac_table_count,omitempty"`
		PoeClass            string               `json:"poe_class,omitempty"`
		PoeCurrent          string               `json:"poe_current,omitempty"`
		PoeEnable           bool                 `json:"poe_enable,omitempty"`
		PoeGood             bool                 `json:"poe_good,omitempty"`
		PoePower            string               `json:"poe_power,omitempty"`
		PoeVoltage          string               `json:"poe_voltage,omitempty"`
		RxBroadcast         emptyStringInt       `json:"rx_broadcast,omitempty"`
		RxBytes             emptyStringInt       `json:"rx_bytes,omitempty"`
		RxDropped           emptyStringInt       `json:"rx_dropped,omitempty"`
		RxErrors            emptyStringInt       `json:"rx_errors,omitempty"`
		RxMulticast         emptyStringInt       `json:"rx_multicast,omitempty"`
		RxPackets           emptyStringInt       `json:"rx_packets,omitempty"`
		Satisfaction        emptyStringInt       `json:"satisfaction,omitempty"`
		SatisfactionReason  emptyStringInt       `json:"satisfaction_reason,omitempty"`
		Speed               emptyStringInt       `json:"speed,omitempty"`
		StpPathcost         emptyStringInt       `json:"stp_pathcost,omitempty"`
		StpState            string               `json:"stp_state,omitempty"`
		TxBroadcast         emptyStringInt       `json:"tx_broadcast,omitempty"`
		TxBytes             emptyStringInt       `json:"tx_bytes,omitempty"`
		TxDropped           emptyStringInt       `json:"tx_dropped,omitempty"`
		TxErrors            emptyStringInt       `json:"tx_errors,omitempty"`
		TxMulticast         emptyStringInt       `json:"tx_multicast,omitempty"`
		TxPackets           emptyStringInt       `json:"tx_packets,omitempty"`
		Up                  bool                 `json:"up,omitempty"`
		TxBytesR            float64              `json:"tx_bytes-r,omitempty"`
		RxBytesR            float64              `json:"rx_bytes-r,omitempty"`
		BytesR              float64              `json:"bytes-r,omitempty"`
		FlowControlEnabled  bool                 `json:"flow_control_enabled,omitempty"`
		NativeNetworkconfID string               `json:"native_networkconf_id,omitempty"`
		Name                string               `json:"name,omitempty"`
		SettingPreference   string               `json:"setting_preference,omitempty"`
		StormctrlBcastRate  emptyStringInt       `json:"stormctrl_bcast_rate,omitempty"`
		StormctrlMcastRate  emptyStringInt       `json:"stormctrl_mcast_rate,omitempty"`
		StormctrlUcastRate  emptyStringInt       `json:"stormctrl_ucast_rate,omitempty"`
		TaggedVlanMgmt      string               `json:"tagged_vlan_mgmt,omitempty"`
		Masked              bool                 `json:"masked,omitempty"`
		AggregatedBy        bool                 `json:"aggregated_by,omitempty"`
	}

	var alias Alias
	if err := json.Unmarshal(b, &alias); err != nil {
		return err
	}

	dst.PortIdx = int(alias.PortIdx)
	dst.Media = alias.Media
	dst.PortPoe = alias.PortPoe
	dst.PoeCaps = int(alias.PoeCaps)
	dst.SpeedCaps = int(alias.SpeedCaps)
	dst.LastConnection = alias.LastConnection
	dst.OpMode = alias.OpMode
	dst.Forward = alias.Forward
	dst.PoeMode = alias.PoeMode
	dst.Anomalies = int(alias.Anomalies)
	dst.Autoneg = alias.Autoneg
	dst.Dot1XMode = alias.Dot1XMode
	dst.Dot1XStatus = alias.Dot1XStatus
	dst.Enable = alias.Enable
	dst.FlowctrlRx = alias.FlowctrlRx
	dst.FlowctrlTx = alias.FlowctrlTx
	dst.FullDuplex = alias.FullDuplex
	dst.IsUplink = alias.IsUplink
	dst.Jumbo = alias.Jumbo
	dst.MacTableCount = int(alias.MacTableCount)
	dst.PoeClass = alias.PoeClass
	dst.PoeCurrent = alias.PoeCurrent
	dst.PoeEnable = alias.PoeEnable
	dst.PoeGood = alias.PoeGood
	dst.PoePower = alias.PoePower
	dst.PoeVoltage = alias.PoeVoltage
	dst.RxBroadcast = int(alias.RxBroadcast)
	dst.RxBytes = int(alias.RxBytes)
	dst.RxDropped = int(alias.RxDropped)
	dst.RxErrors = int(alias.RxErrors)
	dst.RxMulticast = int(alias.RxMulticast)
	dst.RxPackets = int(alias.RxPackets)
	dst.Satisfaction = int(alias.Satisfaction)
	dst.SatisfactionReason = int(alias.SatisfactionReason)
	dst.Speed = int(alias.Speed)
	dst.StpPathcost = int(alias.StpPathcost)
	dst.StpState = alias.StpState
	dst.TxBroadcast = int(alias.TxBroadcast)
	dst.TxBytes = int64(alias.TxBytes)
	dst.TxDropped = int(alias.TxDropped)
	dst.TxErrors = int(alias.TxErrors)
	dst.TxMulticast = int(alias.TxMulticast)
	dst.TxPackets = int(alias.TxPackets)
	dst.Up = alias.Up
	dst.TxBytesR = alias.TxBytesR
	dst.RxBytesR = alias.RxBytesR
	dst.BytesR = alias.BytesR
	dst.FlowControlEnabled = alias.FlowControlEnabled
	dst.NativeNetworkconfID = alias.NativeNetworkconfID
	dst.Name = alias.Name
	dst.SettingPreference = alias.SettingPreference
	dst.StormctrlBcastRate = int(alias.StormctrlBcastRate)
	dst.StormctrlMcastRate = int(alias.StormctrlMcastRate)
	dst.StormctrlUcastRate = int(alias.StormctrlUcastRate)
	dst.TaggedVlanMgmt = alias.TaggedVlanMgmt
	dst.Masked = alias.Masked
	dst.AggregatedBy = alias.AggregatedBy

	return nil
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
