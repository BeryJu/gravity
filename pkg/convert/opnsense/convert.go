package opnsense

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/netip"
	"strings"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"go.uber.org/zap"
)

type Converter struct {
	a      *api.APIClient
	l      *zap.Logger
	parsed Opnsense
	zone   string

	createdZones map[string]struct{}
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

func New(api *api.APIClient, raw []byte, options ...ConverterOption) (*Converter, error) {
	conv := &Converter{
		a:            api,
		l:            extconfig.Get().Logger().Named("convert.bind"),
		createdZones: make(map[string]struct{}),
	}
	var ops Opnsense
	err := xml.Unmarshal(raw, &ops)
	if err != nil {
		return nil, err
	}
	conv.parsed = ops
	for _, opt := range options {
		opt.apply(conv)
	}
	return conv, nil
}

func (c *Converter) Run(ctx context.Context) error {
	for _, record := range c.parsed.Dnsmasq.Hosts {
		err := c.checkDomainExists(ctx, record.Domain)
		if err != nil {
			return err
		}
		err = c.convertHost(ctx, record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Converter) checkDomainExists(ctx context.Context, domain string) error {
	if c.zone != "" && !strings.EqualFold(domain, c.zone) {
		return nil
	}
	if _, ok := c.createdZones[domain]; ok {
		return nil
	}
	_, err := c.a.RolesDnsApi.DnsPutZones(ctx).DnsAPIZonesPutInput(api.DnsAPIZonesPutInput{
		Authoritative: true,
		DefaultTTL:    86400,
		HandlerConfigs: []map[string]interface{}{
			{
				"type": "etcd",
			},
		},
	}).Zone(domain).Execute()
	if err != nil {
		return err
	}
	c.createdZones[domain] = struct{}{}
	c.l.Info("converted zone", zap.String("name", domain))
	return nil
}

func (c *Converter) convertHost(ctx context.Context, r OpnsenseHost) error {
	req := api.DnsAPIRecordsPutInput{
		Type: "A",
	}

	i, err := netip.ParseAddr(r.Ip)
	if err != nil {
		return err
	}
	if i.Is6() {
		req.Type = "AAAA"
	}
	req.Data = r.Ip

	_, err = c.a.RolesDnsApi.DnsPutRecords(ctx).
		DnsAPIRecordsPutInput(req).
		Zone(utils.EnsureTrailingPeriod(r.Domain)).
		Hostname(r.Host).
		Uid("opnsense-import").
		Execute()
	if err != nil {
		return err
	}
	c.l.Info("converted record", zap.String("name", r.Host))
	for _, alias := range r.Aliases.Item {
		err := c.checkDomainExists(ctx, alias.Domain)
		if err != nil {
			return err
		}
		err = c.convertAlias(ctx, r, alias)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Converter) convertAlias(ctx context.Context, h OpnsenseHost, r OpnsenseAlias) error {
	req := api.DnsAPIRecordsPutInput{
		Type: "CNAME",
	}
	req.Data = fmt.Sprintf("%s.%s.", h.Host, h.Domain)

	_, err := c.a.RolesDnsApi.DnsPutRecords(ctx).
		DnsAPIRecordsPutInput(req).
		Zone(utils.EnsureTrailingPeriod(r.Domain)).
		Hostname(r.Host).
		Uid("opnsense-import").
		Execute()
	if err != nil {
		return err
	}
	c.l.Info("converted alias", zap.String("name", r.Host))
	return nil
}
