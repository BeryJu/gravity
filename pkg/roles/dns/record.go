package dns

import (
	"context"
	"encoding/json"
	"net"
	"strings"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Record struct {
	inst roles.Instance
	zone *Zone
	Name string `json:"-"`
	Type string `json:"-"`

	Data      string `json:"data"`
	uid       string
	recordKey string
	TTL       uint32 `json:"ttl,omitempty"`

	MXPreference uint16 `json:"mxPreference,omitempty"`
	SRVPort      uint16 `json:"srvPort,omitempty"`
	SRVPriority  uint16 `json:"srvPriority,omitempty"`
	SRVWeight    uint16 `json:"srvWeight,omitempty"`
}

func (z *Zone) recordFromKV(kv *mvccpb.KeyValue) (*Record, error) {
	fullRecordKey := string(kv.Key)
	// Relative key compared to zone, format of
	// host/A[/...]
	relKey := strings.TrimPrefix(fullRecordKey, z.inst.KV().Key(z.etcdKey).Prefix(true).String())
	// parts[0] is the hostname, parts[1] is the type
	// parts[2] is the remaindar
	parts := strings.SplitN(relKey, "/", 3)
	if len(parts) < 2 {
		parts = []string{"", ""}
	}
	rec := z.newRecord(parts[0], parts[1])
	if len(parts) > 2 {
		rec.uid = parts[2]
	}
	rec.recordKey = strings.TrimSuffix(fullRecordKey, "/"+rec.uid)
	err := json.Unmarshal(kv.Value, &rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (z *Zone) newRecord(name string, t string) *Record {
	return &Record{
		Name: strings.ToLower(name),
		Type: t,
		inst: z.inst,
		zone: z,
	}
}

func (r *Record) RRType() uint16 {
	switch strings.ToLower(r.Type) {
	case "a":
		return dns.TypeA
	case "aaaa":
		return dns.TypeAAAA
	case "ptr":
		return dns.TypePTR
	case "srv":
		return dns.TypeSRV
	case "mx":
		return dns.TypeMX
	case "cname":
		return dns.TypeCNAME
	case "txt":
		return dns.TypeTXT
	}
	return dns.TypeNone
}

func (r *Record) FQDN(qname string) string {
	parts := []string{}
	qnameParts := strings.Split(qname, types.DNSSep)
	for idx, lp := range strings.Split(r.Name, types.DNSSep) {
		if lp == types.DNSWildcard {
			parts = append(parts, qnameParts[idx])
		} else if lp != types.DNSRootRecord {
			parts = append(parts, lp)
		}
	}
	if r.zone.Name != types.DNSRootZone {
		parts = append(parts, r.zone.Name)
	}
	return utils.EnsureTrailingPeriod(strings.Join(parts, types.DNSSep))
}

func (r *Record) ToDNS(qname string) dns.RR {
	hdr := dns.RR_Header{
		Name:   r.FQDN(qname),
		Rrtype: r.RRType(),
		Class:  dns.ClassINET,
		Ttl:    r.TTL,
	}
	var rr dns.RR
	switch r.RRType() {
	case dns.TypeA:
		rr = &dns.A{
			Hdr: hdr,
			A:   net.ParseIP(r.Data),
		}
	case dns.TypeAAAA:
		rr = &dns.AAAA{
			Hdr:  hdr,
			AAAA: net.ParseIP(r.Data),
		}
	case dns.TypePTR:
		rr = &dns.PTR{
			Hdr: hdr,
			Ptr: r.Data,
		}
	case dns.TypeSRV:
		rr = &dns.SRV{
			Hdr:      hdr,
			Target:   r.Data,
			Port:     r.SRVPort,
			Priority: r.SRVPriority,
			Weight:   r.SRVWeight,
		}
	case dns.TypeMX:
		rr = &dns.MX{
			Hdr:        hdr,
			Mx:         r.Data,
			Preference: r.MXPreference,
		}
	case dns.TypeCNAME:
		rr = &dns.CNAME{
			Hdr:    hdr,
			Target: utils.EnsureTrailingPeriod(r.Data),
		}
	case dns.TypeTXT:
		rr = &dns.TXT{
			Hdr: hdr,
			Txt: strings.Split(r.Data, types.TXTSeparator),
		}
	}
	return rr
}

func (r *Record) put(ctx context.Context, expiry int64, opts ...clientv3.OpOption) error {
	raw, err := json.Marshal(&r)
	if err != nil {
		return err
	}

	if expiry > 0 {
		exp, err := r.inst.KV().Grant(ctx, expiry)
		if err != nil {
			return err
		}
		opts = append(opts, clientv3.WithLease(exp.ID))
	}

	recordKey := r.inst.KV().Key(
		types.KeyRole,
		types.KeyZones,
		strings.ToLower(r.zone.Name),
		strings.ToLower(r.Name),
		r.Type,
	)
	if r.uid != "" {
		recordKey.Add(r.uid)
	}
	_, err = r.inst.KV().Put(
		ctx,
		recordKey.String(),
		string(raw),
		opts...,
	)
	return err
}
