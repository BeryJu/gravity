package dns

import (
	"context"
	"encoding/json"
	"net"
	"strings"

	"beryju.io/ddet/pkg/roles"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

type Zone struct {
	Name string `json:"-"`

	Authoritative  bool                `json:"authoritative"`
	HandlerConfigs []map[string]string `json:"handlerConfigs"`

	h []Handler

	etcdKey string
	inst    roles.Instance
	log     *log.Entry
}

type Record struct {
	Data string `json:"data"`
	TTL  uint32 `json:"ttl"`
}

func (z *Zone) resolve(w dns.ResponseWriter, r *dns.Msg) {
	for _, handler := range z.h {
		handlerReply := handler.Handle(NweFakeDNSWriter(w), r)
		if handlerReply != nil {
			handlerReply.SetReply(r)
			w.WriteMsg(handlerReply)
			return
		}
	}
	z.log.Debug("No handler has a reply, fallback back to NX")
	fallback := new(dns.Msg)
	fallback.SetReply(r)
	fallback.SetRcode(r, dns.RcodeNameError)
	w.WriteMsg(fallback)
}

func (z *Zone) kvToDNS(qname string, kv *mvccpb.KeyValue, t uint16) dns.RR {
	rec := Record{}
	err := json.Unmarshal(kv.Value, &rec)
	if err != nil {
		z.log.WithError(err).Warning("failed to parse record")
		return nil
	}
	hdr := dns.RR_Header{
		Name:   qname,
		Rrtype: t,
		Class:  dns.ClassINET,
		Ttl:    rec.TTL,
	}
	return &dns.A{
		Hdr: hdr,
		A:   net.ParseIP(rec.Data),
	}
}

func (r *DNSRole) zoneFromKV(raw *mvccpb.KeyValue) (*Zone, error) {
	z := Zone{
		inst: r.i,
		h:    make([]Handler, 0),
	}
	err := json.Unmarshal(raw.Value, &z)
	if err != nil {
		return nil, err
	}
	prefix := r.i.GetKV().Key(KeyRole, KeyZones, "")
	z.Name = strings.TrimPrefix(string(raw.Key), prefix)
	// Get full etcd key without leading slash since this usually gets passed to Instance Key()
	z.etcdKey = string(raw.Key)[1:]
	z.log = log.WithField("zone", z.Name)

	if len(z.HandlerConfigs) < 1 {
		r.log.Trace("No handler defined, defaulting to etcd")
		z.HandlerConfigs = append(z.HandlerConfigs, map[string]string{
			"type": "etcd",
		})
	}

	for _, handlerCfg := range z.HandlerConfigs {
		t := handlerCfg["type"]
		var handler Handler
		var err error
		switch t {
		case "forward_blocky":
			handler, err = NewBlockyForwarder(z, handlerCfg)
		case "forward_ip":
			handler = NewIPForwarderHandler(z, handlerCfg)
		case "etcd":
			handler = NewEtcdHandler(z, handlerCfg)
		default:
			r.log.WithField("type", t).Warning("invalid forwarder type")
		}
		if err != nil {
			z.log.WithError(err).Warning("failed to setup handler")
			continue
		}
		z.h = append(z.h, handler)
	}
	return &z, nil
}

func (z *Zone) put() error {
	raw, err := json.Marshal(&z)
	if err != nil {
		return err
	}

	_, err = z.inst.GetKV().Put(
		context.TODO(),
		z.inst.GetKV().Key(
			KeyRole,
			KeyZones,
			z.Name,
		),
		string(raw),
	)
	if err != nil {
		return err
	}
	return nil
}
