package bind

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

type Converter struct {
	a    *api.APIClient
	l    *zap.Logger
	i    io.Reader
	zone string
}

type ConverterOption struct {
	apply func(*Converter)
}

func WithExistingZone(name string) ConverterOption {
	return ConverterOption{
		apply: func(c *Converter) {
			c.zone = name
		},
	}
}

func New(api *api.APIClient, input io.Reader, options ...ConverterOption) (*Converter, error) {
	conv := &Converter{
		a: api,
		i: input,
		l: extconfig.Get().Logger().Named("convert.bind"),
	}
	for _, opt := range options {
		opt.apply(conv)
	}
	return conv, nil
}

func (c *Converter) Run(ctx context.Context) error {
	p := dns.NewZoneParser(c.i, "", "")

	records := []dns.RR{}
	for rr, ok := p.Next(); ok; rr, ok = p.Next() {
		records = append(records, rr)
	}

	// First create zone
	if c.zone == "" {
		for _, rr := range records {
			if rr.Header().Rrtype == dns.TypeSOA {
				err := c.createZone(rr, ctx)
				if err != nil {
					return err
				}
			}
		}
	}
	// Then create all the records
	for _, rr := range records {
		if rr.Header().Rrtype == dns.TypeSOA {
			continue
		}
		err := c.convertRecord(rr, ctx)
		if err != nil {
			return err
		}
	}
	return p.Err()
}

func (c *Converter) createZone(rr dns.RR, ctx context.Context) error {
	_r := rr.(*dns.SOA)
	_, err := c.a.RolesDnsApi.DnsPutZones(ctx).DnsAPIZonesPutInput(api.DnsAPIZonesPutInput{
		Authoritative: true,
		DefaultTTL:    int32(_r.Expire),
		HandlerConfigs: []map[string]interface{}{
			{
				"type": "etcd",
			},
		},
	}).Zone(_r.Hdr.Name).Execute()
	if err != nil {
		return err
	}
	c.zone = _r.Hdr.Name
	c.l.Info("converted zone", zap.String("name", c.zone))
	return nil
}

func (c *Converter) convertRecord(rr dns.RR, ctx context.Context) error {
	req := api.DnsAPIRecordsPutInput{
		Type: dns.TypeToString[rr.Header().Rrtype],
	}

	switch v := rr.(type) {
	case *dns.A:
		req.Data = v.A.String()
	case *dns.AAAA:
		req.Data = v.AAAA.String()
	case *dns.TXT:
		req.Data = strings.Join(v.Txt, types.TXTSeparator)
	case *dns.PTR:
		req.Data = v.Ptr
	case *dns.CNAME:
		req.Data = v.Target
	case *dns.MX:
		req.Data = v.Mx
		req.MxPreference = api.PtrInt32(int32(v.Preference))
	case *dns.SRV:
		req.Data = v.Target
		req.SrvPort = api.PtrInt32(int32(v.Port))
		req.SrvPriority = api.PtrInt32(int32(v.Priority))
		req.SrvWeight = api.PtrInt32(int32(v.Weight))
	default:
		c.l.Info("unsupported record type", zap.String("name", rr.Header().Name), zap.String("type", dns.TypeToString[rr.Header().Rrtype]))
		return nil
	}

	relName := strings.TrimSuffix(rr.Header().Name, utils.EnsureTrailingPeriod(c.zone))
	if rr.Header().Name == c.zone {
		relName = types.DNSRootRecord
	}

	hasher := sha512.New()
	hasher.Write([]byte(rr.String()))

	_, err := c.a.RolesDnsApi.DnsPutRecords(ctx).
		DnsAPIRecordsPutInput(req).
		Zone(c.zone).
		Hostname(relName).
		Uid(fmt.Sprintf("bind-import-%s", hex.EncodeToString(hasher.Sum(nil))[:8])).
		Execute()
	if err != nil {
		return err
	}
	c.l.Info("converted record", zap.String("name", relName))
	return nil
}
