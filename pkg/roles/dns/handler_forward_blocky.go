package dns

import (
	"fmt"
	"net/netip"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/creasty/defaults"
	log "github.com/sirupsen/logrus"

	"github.com/0xERR0R/blocky/config"
	blockylog "github.com/0xERR0R/blocky/log"
	"github.com/0xERR0R/blocky/server"
	"github.com/miekg/dns"
)

type BlockyForwarder struct {
	*IPForwarderHandler
	c   map[string]string
	b   *server.Server
	log *log.Entry
}

func NewBlockyForwarder(z *Zone, rawConfig map[string]string) *BlockyForwarder {
	bfwd := &BlockyForwarder{
		IPForwarderHandler: NewIPForwarderHandler(z, rawConfig),
		c:                  rawConfig,
	}
	bfwd.log = z.log.WithField("handler", bfwd.Identifier())
	go func() {
		err := bfwd.setup()
		if err != nil {
			bfwd.log.WithError(err).Warning("failed to setup blocky, queries will fallthrough")
		}
	}()
	return bfwd
}

func (bfwd *BlockyForwarder) Identifier() string {
	return "forward_blocky"
}

func (bfwd *BlockyForwarder) setup() error {
	forwarders := strings.Split(bfwd.c["to"], ";")
	upstreams := make([]config.Upstream, len(forwarders))
	for idx, fwd := range forwarders {
		us, err := config.ParseUpstream(fwd)
		if err != nil {
			log.WithError(err).Warning("failed to parse upstream")
			continue
		}
		upstreams[idx] = us
	}
	cfg := config.Config{}
	err := defaults.Set(&cfg)
	if err != nil {
		return fmt.Errorf("failed to set config defaults: %w", err)
	}
	// Blocky uses a custom registry, so this doesn't work as expected
	// cfg.Prometheus.Enable = true
	if !extconfig.Get().Debug {
		cfg.LogFormat = blockylog.FormatTypeJson
	}
	bootstrap, err := netip.ParseAddrPort(extconfig.Get().FallbackDNS)
	if err != nil {
		return err
	}
	cfg.BootstrapDNS = config.BootstrapConfig{
		Upstream: config.Upstream{
			Net:  config.NetProtocolTcpUdp,
			Host: bootstrap.Addr().String(),
			Port: bootstrap.Port(),
		},
	}
	cfg.Upstream = config.UpstreamConfig{
		ExternalResolvers: map[string][]config.Upstream{
			"default": upstreams,
		},
	}
	// TODO: Blocky config
	cfg.Blocking = config.BlockingConfig{
		BlockType: "zeroIP",
		BlackLists: map[string][]string{
			"block": {
				"https://adaway.org/hosts.txt",
				"https://dbl.oisd.nl/",
				"https://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&showintro=0&mimetype=plaintext",
				"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
				"https://v.firebog.net/hosts/AdguardDNS.txt",
				"https://v.firebog.net/hosts/Easylist.txt",
				"https://www.github.developerdan.com/hosts/lists/ads-and-tracking-extended.txt",
				"https://adguardteam.github.io/AdGuardSDNSFilter/Filters/filter.txt",
			},
		},
		ClientGroupsBlock: map[string][]string{
			"default": {"block"},
		},
	}

	srv, err := server.NewServer(&cfg)
	if err != nil {
		return fmt.Errorf("can't start server: %w", err)
	}
	bfwd.b = srv
	return nil
}

func (bfwd *BlockyForwarder) Handle(w *utils.FakeDNSWriter, r *dns.Msg) *dns.Msg {
	if bfwd.b == nil {
		bfwd.log.Info("Blocky not started/setup yet, falling through to next handler")
		return nil
	}
	bfwd.b.OnRequest(w, r)
	// fall to next handler when no record is found
	if w.Msg().Rcode == dns.RcodeNameError {
		return nil
	}
	for _, query := range r.Question {
		for idx, ans := range w.Msg().Answer {
			go bfwd.cacheToEtcd(query, ans, idx)
		}
	}
	return w.Msg()
}
