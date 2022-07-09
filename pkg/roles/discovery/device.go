package discovery

import (
	"context"
	"encoding/json"
	"errors"

	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/discovery/types"
)

type Device struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	MAC      string `json:"mac"`

	inst roles.Instance
}

func NewDevice(inst roles.Instance) *Device {
	return &Device{
		inst: inst,
	}
}

func (d *Device) save() error {
	by := ""
	identifier := ""
	if d.IP != "" {
		by = types.KeyDevicesByIP
		identifier = d.IP
	}
	if d.MAC != "" {
		by = types.KeyDevicesByMAC
		identifier = d.MAC
	}
	if by == "" {
		return errors.New("device without IP and MAC")
	}
	key := d.inst.GetKV().Key(types.KeyRole, types.KeyDevices, by, identifier)
	raw, err := json.Marshal(&d)
	if err != nil {
		return err
	}
	_, err = d.inst.GetKV().Put(
		context.Background(),
		key,
		string(raw),
	)
	if err != nil {
		return err
	}
	d.inst.DispatchEvent(types.EventTopicDiscoveryDeviceFound, roles.NewEvent(map[string]interface{}{
		"device": d,
	}))
	return nil
}
