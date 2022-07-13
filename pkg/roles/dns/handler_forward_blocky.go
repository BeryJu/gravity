package dns

import (
	"fmt"
	"net"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/creasty/defaults"
	log "github.com/sirupsen/logrus"

	"github.com/0xERR0R/blocky/config"
	blockylog "github.com/0xERR0R/blocky/log"
	"github.com/0xERR0R/blocky/server"
	"github.com/miekg/dns"
)

type BlockyForwarder struct {
	*IPForwarderHandler
	b   *server.Server
	log *log.Entry
}

func NewBlockyForwarder(z *Zone, rawConfig map[string]string) (*BlockyForwarder, error) {
	bfwd := &BlockyForwarder{
		IPForwarderHandler: NewIPForwarderHandler(z, rawConfig),
	}
	bfwd.log = z.log.WithField("handler", bfwd.Identifier())
	forwarders := strings.Split(rawConfig["to"], ";")
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
		return nil, fmt.Errorf("failed to set config defaults: %w", err)
	}
	blockylog.Silence()
	cfg.BootstrapDNS = config.BootstrapConfig{
		IPs: []net.IP{
			net.ParseIP(extconfig.Get().FallbackDNS),
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
			},
		},
		ClientGroupsBlock: map[string][]string{
			"default": {"block"},
		},
	}

	srv, err := server.NewServer(&cfg)
	if err != nil {
		return nil, fmt.Errorf("can't start server: %w", err)
	}
	bfwd.b = srv
	return bfwd, err
}

func (bfwd *BlockyForwarder) Identifier() string {
	return "forward_blocky"
}

func (bfwd *BlockyForwarder) Handle(w *fakeDNSWriter, r *dns.Msg) *dns.Msg {
	bfwd.b.OnRequest(w, r)
	// fall to next handler when no record is found
	if w.msg.Rcode == dns.RcodeNameError {
		return nil
	}
	for _, query := range r.Question {
		for idx, ans := range w.msg.Answer {
			go bfwd.cacheToEtcd(query, ans, idx)
		}
	}
	return w.msg
}
