package dns

import (
	"context"
	"encoding/json"
	"net"
	"strings"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/miekg/dns"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Record struct {
	Name string `json:"-"`
	Type string `json:"-"`

	Data string `json:"data"`
	TTL  uint32 `json:"ttl,omitempty"`

	MXPreference uint16 `json:"mxPreference,omitempty"`
	SRVPort      uint16 `json:"srvPort,omitempty"`
	SRVPriority  uint16 `json:"srvPriority,omitempty"`
	SRVWeight    uint16 `json:"srvWeight,omitempty"`

	inst      roles.Instance
	zone      *Zone
	uid       string
	recordKey string
}

func (z *Zone) recordFromKV(kv *mvccpb.KeyValue) (*Record, error) {
	rec := Record{}
	fullRecordKey := string(kv.Key)
	// Relative key compared to zone, format of
	// host/A[/...]
	relKey := strings.TrimPrefix(fullRecordKey, z.inst.KV().Key(z.etcdKey).Prefix(true).String())
	// parts[0] is the hostname, parts[1] is the type
	// parts[2] is the remaindar
	parts := strings.SplitN(relKey, "/", 3)
	if len(parts) > 2 {
		rec.uid = parts[2]
	}
	rec.recordKey = strings.TrimSuffix(fullRecordKey, "/"+rec.uid)
	err := json.Unmarshal(kv.Value, &rec)
	if err != nil {
		return &rec, err
	}
	return &rec, nil
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
		rr.(*dns.SRV).Port = r.SRVPort
		rr.(*dns.SRV).Priority = r.SRVPriority
		rr.(*dns.SRV).Weight = r.SRVWeight
	case dns.TypeMX:
		rr = &dns.MX{}
		rr.(*dns.MX).Hdr = hdr
		rr.(*dns.MX).Mx = r.Data
		rr.(*dns.MX).Preference = r.MXPreference
	case dns.TypeCNAME:
		rr = &dns.CNAME{}
		rr.(*dns.CNAME).Hdr = hdr
		rr.(*dns.CNAME).Target = r.Data
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
	if r.uid != "" {
		leaseKey.Add(r.uid)
	}
	_, err = r.inst.KV().Put(
		context.TODO(),
		leaseKey.String(),
		string(raw),
		opts...,
	)
	if err != nil {
		return err
	}
	return nil
}
