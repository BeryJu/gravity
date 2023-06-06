package dns

import (
	"context"
	"net"
	"strings"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/storage"
	"github.com/miekg/dns"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

const TXTSeparator = "\n"

type RecordContext struct {
	*types.Record

	inst roles.Instance
	zone *ZoneContext

	uid       string
	recordKey string
}

func (z *ZoneContext) recordFromKV(kv *mvccpb.KeyValue) (*RecordContext, error) {
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

	_, err := storage.Parse(kv.Value, rec.Record)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (z *ZoneContext) newRecord(name string, t string) *RecordContext {
	return &RecordContext{
		Record: &types.Record{
			Name: strings.ToLower(name),
			Type: t,
		},
		inst: z.inst,
		zone: z,
	}
}

func (r *RecordContext) ToDNS(qname string, t uint16) dns.RR {
	hdr := dns.RR_Header{
		Name:   qname,
		Rrtype: t,
		Class:  dns.ClassINET,
		Ttl:    r.Ttl,
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
		rr.(*dns.SRV).Port = uint16(r.SrvPort)
		rr.(*dns.SRV).Priority = uint16(r.SrvPriority)
		rr.(*dns.SRV).Weight = uint16(r.SrvWeight)
	case dns.TypeMX:
		rr = &dns.MX{}
		rr.(*dns.MX).Hdr = hdr
		rr.(*dns.MX).Mx = r.Data
		rr.(*dns.MX).Preference = uint16(r.MxPreference)
	case dns.TypeCNAME:
		rr = &dns.CNAME{}
		rr.(*dns.CNAME).Hdr = hdr
		rr.(*dns.CNAME).Target = r.Data
	case dns.TypeTXT:
		rr = &dns.TXT{}
		rr.(*dns.TXT).Txt = strings.Split(r.Data, TXTSeparator)
	}
	return rr
}

func (r *RecordContext) put(ctx context.Context, expiry int64, opts ...clientv3.OpOption) error {
	raw, err := proto.Marshal(r.Record)
	if err != nil {
		return err
	}

	if expiry > 0 {
		exp, err := r.inst.KV().Lease.Grant(ctx, expiry)
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
