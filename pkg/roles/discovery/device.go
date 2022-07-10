package discovery

import (
	"context"
	"encoding/json"
	"errors"

	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/discovery/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
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

func (r *DiscoveryRole) deviceFromKV(kv *mvccpb.KeyValue) *Device {
	rec := Device{
		inst: r.i,
	}
	err := json.Unmarshal(kv.Value, &rec)
	if err != nil {
		r.log.WithError(err).Warning("failed to parse device")
		return nil
	}
	return &rec
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
		exp, err := d.inst.KV().Lease.Grant(context.TODO(), expiry)
		if err != nil {
			return err
		}
		opts = append(opts, clientv3.WithLease(exp.ID))
	}

	key := d.inst.KV().Key(types.KeyRole, types.KeyDevices, by, identifier)
	raw, err := json.Marshal(&d)
	if err != nil {
		return err
	}
	_, err = d.inst.KV().Put(
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

func (d *Device) toDNS(zone string) {
	// TODO: Stub method
}

func (d *Device) toDHCP() {
	// TODO: Stub method
}
