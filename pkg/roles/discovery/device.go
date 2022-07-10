package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"beryju.io/gravity/pkg/roles"
	dhcptypes "beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/roles/discovery/types"
	dnstypes "beryju.io/gravity/pkg/roles/dns/types"
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

func (d *Device) toDNS(zone string) error {
	hostname := strings.Split(d.Hostname, ".")
	fqdn := d.Hostname
	if zone != "" {
		fqdn = hostname[0] + "." + zone
	}
	if zone == "" && len(hostname) == 1 {
		return errors.New("device hostname has no domain and no zone given")
	}
	d.inst.DispatchEvent(dnstypes.EventTopicDNSRecordCreateForward, roles.NewEvent(map[string]interface{}{
		"fqdn":     fqdn,
		"hostname": hostname[0],
		"address":  d.IP,
	}))
	d.inst.DispatchEvent(dnstypes.EventTopicDNSRecordCreateReverse, roles.NewEvent(map[string]interface{}{
		"fqdn":    fqdn,
		"address": d.IP,
	}))
	// Maybe delete device? Mark as applied?
	return nil
}

func (d *Device) toDHCP(scope string) error {
	if scope == "" {
		return errors.New("blank scope")
	}
	if d.MAC == "" {
		return errors.New("mac address blank")
	}
	hostname := strings.Split(d.Hostname, ".")[0]
	d.inst.DispatchEvent(dhcptypes.EventTopicDHCPCreateLease, roles.NewEvent(map[string]interface{}{
		"mac":      d.MAC,
		"hostname": hostname,
		"address":  d.IP,
		"scope":    scope,
	}))
	return nil
}
