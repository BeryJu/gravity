package dns

import (
	"context"
	"encoding/json"
	"net"
	"strings"

	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles"
	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles/dns/forward"
	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles/dns/forward/ip"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

type Zone struct {
	Name string

	Authoritative bool `json:"authoritative"`

	Forwarder *map[string]string `json:"forwarder"`

	fwd forward.Forwarder

	etcdKey string
	inst    roles.Instance
	log     *log.Entry
}

func (z *Zone) resolve(w dns.ResponseWriter, r *dns.Msg) {
	if z.fwd != nil {
		z.fwd.Handle(w, r)
		return
	}

	m := new(dns.Msg)
	m.SetReply(r)

	m.Authoritative = z.Authoritative

	ctx := context.Background()
	for _, question := range r.Question {
		relRecordName := strings.TrimSuffix(question.Name, "."+z.Name)
		fullRecordKey := z.inst.GetKV().Key(z.etcdKey, relRecordName, dns.Type(question.Qtype).String())
		// TODO: Optimise this
		res, err := z.inst.GetKV().Get(ctx, fullRecordKey)
		if err != nil || len(res.Kvs) < 1 {
			continue
		}
		for _, key := range res.Kvs {
			ans := z.kvToDNS(question.Name, key, question.Qtype)
			if ans != nil {
				m.Answer = append(m.Answer, ans)
			}
		}
	}
	if len(m.Answer) < 1 {
		m.SetRcode(r, dns.RcodeNameError)
	}
	w.WriteMsg(m)
}

type Record struct {
	Data string `json:"data"`
	TTL  uint32 `json:"ttl"`
}

func (z *Zone) kvToDNS(qname string, kv *mvccpb.KeyValue, t uint16) dns.RR {
	rec := Record{}
	err := json.Unmarshal(kv.Value, &rec)
	if err != nil {
		z.log.WithError(err).Warning("failed to parse record")
		return nil
	}
	return &dns.A{
		Hdr: dns.RR_Header{
			Name:   qname,
			Rrtype: t,
			Class:  dns.ClassINET,
			Ttl:    rec.TTL,
		},
		A: net.ParseIP(rec.Data),
	}
}

func (r *DNSServerRole) zoneFromKV(inst roles.Instance, raw *mvccpb.KeyValue) (*Zone, error) {
	z := Zone{
		inst: inst,
	}
	err := json.Unmarshal(raw.Value, &z)
	if err != nil {
		return nil, err
	}
	prefix := inst.GetKV().Key(KeyRole, KeyZones, "")
	z.Name = strings.TrimPrefix(string(raw.Key), prefix)
	// Get full etcd key without leading slash since this usually gets passed to Instance Key()
	z.etcdKey = string(raw.Key)[1:]

	if z.Forwarder != nil {
		f := *z.Forwarder
		t := f["type"]
		switch t {
		case "ip":
			z.fwd = ip.New(f)
		default:
			r.log.WithField("type", t).Warning("invalid forwarder type")
		}
	}

	z.log = log.WithField("zone", z.Name)
	return &z, nil
}
