package discovery

import (
	"context"
	"encoding/json"
	"errors"

	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/discovery/types"
	clientv3 "go.etcd.io/etcd/client/v3"
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

func (d *Device) put(expiry int64, opts ...clientv3.OpOption) error {
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

	if expiry > 0 {
		exp, err := d.inst.GetKV().Lease.Grant(context.TODO(), expiry)
		if err != nil {
			return err
		}
		opts = append(opts, clientv3.WithLease(exp.ID))
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
		opts...,
	)
	if err != nil {
		return err
	}

	ev := roles.NewEvent(
		map[string]interface{}{
			"device": d,
		},
	)
	ev.Payload.RelatedObjectKey = key
	ev.Payload.RelatedObjectOptions = opts
	d.inst.DispatchEvent(types.EventTopicDiscoveryDeviceFound, ev)

	return nil
}
