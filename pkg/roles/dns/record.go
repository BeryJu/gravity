package dns

import (
	"context"
	"encoding/json"
	"net"

	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/dns/types"
	"github.com/miekg/dns"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Record struct {
	Name string `json:"-"`
	Type string `json:"-"`

	Data string `json:"data"`
	TTL  uint32 `json:"ttl"`

	inst roles.Instance
	zone *Zone
}

func (z *Zone) recordFromKV(kv *mvccpb.KeyValue) *Record {
	rec := Record{}
	err := json.Unmarshal(kv.Value, &rec)
	if err != nil {
		z.log.WithError(err).Warning("failed to parse record")
		return nil
	}
	return &rec
}

func (z *Zone) newRecord(name string, t string) *Record {
	return &Record{
		Name: name,
		Type: t,
		inst: z.inst,
		zone: z,
	}
}

func (r *Record) ToDNS(qname string, t uint16) dns.RR {
	hdr := dns.RR_Header{
		Name:   qname,
		Rrtype: t,
		Class:  dns.ClassINET,
		Ttl:    r.TTL,
	}
	var rr dns.RR
	switch t {
	case dns.TypeA:
		rr = &dns.A{}
		rr.(*dns.A).Hdr = hdr
		rr.(*dns.A).A = net.ParseIP(r.Data)
	case dns.TypeAAAA:
		rr = &dns.AAAA{}
		rr.(*dns.AAAA).Hdr = hdr
		rr.(*dns.AAAA).AAAA = net.ParseIP(r.Data)
	case dns.TypePTR:
		rr = &dns.PTR{}
		rr.(*dns.PTR).Hdr = hdr
		rr.(*dns.PTR).Ptr = r.Data
	case dns.TypeSRV:
		rr = &dns.SRV{}
		rr.(*dns.SRV).Hdr = hdr
		rr.(*dns.SRV).Target = r.Data
	}
	return rr
}

func (r *Record) put(expiry int64, opts ...clientv3.OpOption) error {
	raw, err := json.Marshal(&r)
	if err != nil {
		return err
	}

	if expiry > 0 {
		exp, err := r.inst.KV().Lease.Grant(context.TODO(), expiry)
		if err != nil {
			return err
		}
		opts = append(opts, clientv3.WithLease(exp.ID))
	}

	leaseKey := r.inst.KV().Key(
		types.KeyRole,
		types.KeyZones,
		r.zone.Name,
		r.Name,
		r.Type,
	)
	_, err = r.inst.KV().Put(
		context.TODO(),
		leaseKey,
		string(raw),
		opts...,
	)
	if err != nil {
		return err
	}
	return nil
}
