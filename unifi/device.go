package unifi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/unifi/types"
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
	type Alias DevicePortTable
	aux := &struct {
		PortIdx            types.Number `json:"port_idx,omitempty"`
		PoeCaps            types.Number `json:"poe_caps,omitempty"`
		SpeedCaps          types.Number `json:"speed_caps,omitempty"`
		Anomalies          types.Number `json:"anomalies,omitempty"`
		MacTableCount      types.Number `json:"mac_table_count,omitempty"`
		RxBroadcast        types.Number `json:"rx_broadcast,omitempty"`
		RxBytes            types.Number `json:"rx_bytes,omitempty"`
		RxDropped          types.Number `json:"rx_dropped,omitempty"`
		RxErrors           types.Number `json:"rx_errors,omitempty"`
		RxMulticast        types.Number `json:"rx_multicast,omitempty"`
		RxPackets          types.Number `json:"rx_packets,omitempty"`
		Satisfaction       types.Number `json:"satisfaction,omitempty"`
		SatisfactionReason types.Number `json:"satisfaction_reason,omitempty"`
		Speed              types.Number `json:"speed,omitempty"`
		StpPathcost        types.Number `json:"stp_pathcost,omitempty"`
		TxBroadcast        types.Number `json:"tx_broadcast,omitempty"`
		TxBytes            types.Number `json:"tx_bytes,omitempty"`
		TxDropped          types.Number `json:"tx_dropped,omitempty"`
		TxErrors           types.Number `json:"tx_errors,omitempty"`
		TxMulticast        types.Number `json:"tx_multicast,omitempty"`
		TxPackets          types.Number `json:"tx_packets,omitempty"`
		StormctrlBcastRate types.Number `json:"stormctrl_bcast_rate,omitempty"`
		StormctrlMcastRate types.Number `json:"stormctrl_mcast_rate,omitempty"`
		StormctrlUcastRate types.Number `json:"stormctrl_ucast_rate,omitempty"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	if portIdx, err := aux.PortIdx.Int64(); err != nil {
		dst.PortIdx = int(portIdx)
	}
	if poeCaps, err := aux.PoeCaps.Int64(); err != nil {
		dst.PoeCaps = int(poeCaps)
	}
	if speedCaps, err := aux.SpeedCaps.Int64(); err != nil {
		dst.SpeedCaps = int(speedCaps)
	}
	if anomalies, err := aux.Anomalies.Int64(); err != nil {
		dst.Anomalies = int(anomalies)
	}
	if macTableCount, err := aux.MacTableCount.Int64(); err != nil {
		dst.MacTableCount = int(macTableCount)
	}
	if rxBroadcast, err := aux.RxBroadcast.Int64(); err != nil {
		dst.RxBroadcast = int(rxBroadcast)
	}
	if rxBytes, err := aux.RxBytes.Int64(); err != nil {
		dst.RxBytes = int(rxBytes)
	}
	if rxDropped, err := aux.RxDropped.Int64(); err != nil {
		dst.RxDropped = int(rxDropped)
	}
	if rxErrors, err := aux.RxErrors.Int64(); err != nil {
		dst.RxErrors = int(rxErrors)
	}
	if rxMulticast, err := aux.RxMulticast.Int64(); err != nil {
		dst.RxMulticast = int(rxMulticast)
	}
	if rxPackets, err := aux.RxPackets.Int64(); err != nil {
		dst.RxPackets = int(rxPackets)
	}
	if satisfaction, err := aux.Satisfaction.Int64(); err != nil {
		dst.Satisfaction = int(satisfaction)
	}
	if satisfactionReason, err := aux.SatisfactionReason.Int64(); err != nil {
		dst.SatisfactionReason = int(satisfactionReason)
	}
	if speed, err := aux.Speed.Int64(); err != nil {
		dst.Speed = int(speed)
	}
	if stpPathcost, err := aux.StpPathcost.Int64(); err != nil {
		dst.StpPathcost = int(stpPathcost)
	}
	if txBroadcast, err := aux.TxBroadcast.Int64(); err != nil {
		dst.TxBroadcast = int(txBroadcast)
	}
	if txBytes, err := aux.TxBytes.Int64(); err != nil {
		dst.TxBytes = txBytes
	}
	if txDropped, err := aux.TxDropped.Int64(); err != nil {
		dst.TxDropped = int(txDropped)
	}
	if txErrors, err := aux.TxErrors.Int64(); err != nil {
		dst.TxErrors = int(txErrors)
	}
	if txMulticast, err := aux.TxMulticast.Int64(); err != nil {
		dst.TxMulticast = int(txMulticast)
	}
	if txPackets, err := aux.TxPackets.Int64(); err != nil {
		dst.TxPackets = int(txPackets)
	}
	if stormctrlBcastRate, err := aux.StormctrlBcastRate.Int64(); err != nil {
		dst.StormctrlBcastRate = int(stormctrlBcastRate)
	}
	if stormctrlMcastRate, err := aux.StormctrlMcastRate.Int64(); err != nil {
		dst.StormctrlMcastRate = int(stormctrlMcastRate)
	}
	if stormctrlUcastRate, err := aux.StormctrlUcastRate.Int64(); err != nil {
		dst.StormctrlUcastRate = int(stormctrlUcastRate)
	}

	return nil
}

func (c *ApiClient) ListDevice(ctx context.Context, site string) ([]Device, error) {
	return c.listDevice(ctx, site)
}

func (c *ApiClient) GetDeviceByMAC(ctx context.Context, site, mac string) (*Device, error) {
	return c.getDevice(ctx, site, mac)
}

func (c *ApiClient) DeleteDevice(ctx context.Context, site, id string) error {
	return c.deleteDevice(ctx, site, id)
}

func (c *ApiClient) CreateDevice(ctx context.Context, site string, d *Device) (*Device, error) {
	return c.createDevice(ctx, site, d)
}

func (c *ApiClient) UpdateDevice(ctx context.Context, site string, d *Device) (*Device, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Device `json:"data"`
	}

	// Get the existing device to compare
	existing, err := c.GetDevice(ctx, site, d.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing device: %w", err)
	}

	// Create a patch with only changed fields
	patch, err := getDeviceDiff(existing, d)
	if err != nil {
		return nil, fmt.Errorf("failed to create device diff: %w", err)
	}

	err = c.do(
		ctx,
		"PUT",
		fmt.Sprintf("api/s/%s/rest/device/%s", site, d.ID),
		patch,
		&respBody,
	)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}

// getDeviceDiff compares two Device objects and returns a map containing only changed fields.
func getDeviceDiff(original, target *Device) (map[string]any, error) {
	// Marshal both to JSON then unmarshal to maps for comparison
	origJSON, err := json.Marshal(original)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal original device: %w", err)
	}

	targetJSON, err := json.Marshal(target)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal target device: %w", err)
	}

	var origMap map[string]any
	var targetMap map[string]any

	if err := json.Unmarshal(origJSON, &origMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal original device: %w", err)
	}

	if err := json.Unmarshal(targetJSON, &targetMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal target device: %w", err)
	}

	// Build patch with only changed fields
	patch := make(map[string]any)

	for key, targetValue := range targetMap {
		// Skip read-only fields
		if key == "_id" || key == "site_id" {
			continue
		}

		origValue, exists := origMap[key]

		// Include if field doesn't exist in original or value changed
		if !exists || !deepEqualJSON(origValue, targetValue) {
			patch[key] = targetValue
		}
	}

	return patch, nil
}

// deepEqualJSON compares two values for deep equality by comparing their JSON representations.
func deepEqualJSON(a, b any) bool {
	aJSON, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bJSON, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(aJSON) == string(bJSON)
}

func (c *ApiClient) GetDevice(ctx context.Context, site, id string) (*Device, error) {
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

func (c *ApiClient) AdoptDevice(ctx context.Context, site, mac string) error {
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

func (c *ApiClient) ForgetDevice(ctx context.Context, site, mac string) error {
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
