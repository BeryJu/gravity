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
	"go.uber.org/zap"
)

type Device struct {
	inst       roles.Instance
	Identifier string `json:"-"`

	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
}

func (r *Role) newDevice() *Device {
	return &Device{
		inst: r.i,
	}
}

func (r *Role) deviceFromKV(kv *mvccpb.KeyValue) *Device {
	prefix := r.i.KV().Key(types.KeyRole, types.KeyDevices).Prefix(true).String()
	identifier := strings.TrimPrefix(string(kv.Key), prefix)
	dev := r.newDevice()
	dev.Identifier = identifier

	err := json.Unmarshal(kv.Value, &dev)
	if err != nil {
		r.log.Warn("failed to parse device", zap.Error(err))
		return nil
	}
	return dev
}

func (d *Device) put(ctx context.Context, expiry int64, opts ...clientv3.OpOption) error {
	if d.IP == "" && d.MAC == "" {
		return errors.New("device without IP and MAC")
	}

	if expiry > 0 {
		exp, err := d.inst.KV().Grant(ctx, expiry)
		if err != nil {
			return err
		}
		opts = append(opts, clientv3.WithLease(exp.ID))
	}

	key := d.inst.KV().Key(types.KeyRole, types.KeyDevices, d.Identifier)
	raw, err := json.Marshal(&d)
	if err != nil {
		return err
	}
	_, err = d.inst.KV().Put(
		ctx,
		key.String(),
		string(raw),
		opts...,
	)
	if err != nil {
		return err
	}

	ev := roles.NewEvent(
		ctx,
		map[string]interface{}{
			"device": d,
		},
	)
	ev.Payload.RelatedObjectKey = key
	ev.Payload.RelatedObjectOptions = opts
	d.inst.DispatchEvent(types.EventTopicDiscoveryDeviceFound, ev)

	return nil
}

func (d *Device) toDNS(ctx context.Context, zone string) error {
	hostname := strings.Split(d.Hostname, ".")
	fqdn := d.Hostname
	if zone != "" {
		fqdn = hostname[0] + "." + zone
	}
	if zone == "" && len(hostname) == 1 {
		return errors.New("device hostname has no domain and no zone given")
	}
	d.inst.DispatchEvent(dnstypes.EventTopicDNSRecordCreateForward, roles.NewEvent(ctx, map[string]interface{}{
		"fqdn":       fqdn,
		"hostname":   hostname[0],
		"address":    d.IP,
		"identifier": d.Identifier,
		"expiry":     int64(-1),
	}))
	d.inst.DispatchEvent(dnstypes.EventTopicDNSRecordCreateReverse, roles.NewEvent(ctx, map[string]interface{}{
		"fqdn":       fqdn,
		"address":    d.IP,
		"identifier": d.Identifier,
		"expiry":     int64(-1),
	}))
	// Maybe delete device? Mark as applied?
	return nil
}

func (d *Device) toDHCP(ctx context.Context, scope string) error {
	if scope == "" {
		return errors.New("blank scope")
	}
	if d.MAC == "" {
		return errors.New("mac address blank")
	}
	hostname := strings.Split(d.Hostname, ".")[0]
	d.inst.DispatchEvent(dhcptypes.EventTopicDHCPCreateLease, roles.NewEvent(ctx, map[string]interface{}{
		"mac":      d.MAC,
		"hostname": hostname,
		"address":  d.IP,
		"scope":    scope,
	}))
	return nil
}
