package dns

import (
	"fmt"
	"net/netip"
	"strings"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/creasty/defaults"
	"go.uber.org/zap"

	"github.com/0xERR0R/blocky/config"
	"github.com/0xERR0R/blocky/log"
	blockylog "github.com/0xERR0R/blocky/log"
	"github.com/0xERR0R/blocky/server"
	"github.com/miekg/dns"
)

type BlockyForwarder struct {
	*IPForwarderHandler
	c   map[string]string
	b   *server.Server
	log *zap.Logger
	st  time.Time
}

func NewBlockyForwarder(z *Zone, rawConfig map[string]string) *BlockyForwarder {
	bfwd := &BlockyForwarder{
		IPForwarderHandler: NewIPForwarderHandler(z, rawConfig),
		c:                  rawConfig,
		st:                 time.Now(),
	}
	bfwd.log = z.log.With(zap.String("handler", bfwd.Identifier()))
	bfwd.log.Debug("starting blocky setup")
	waitForStart := func() {
		err := bfwd.setup()
		if err != nil {
			bfwd.log.Warn("failed to setup blocky, queries will fallthrough", zap.Error(err))
		}
	}
	if extconfig.Get().Debug {
		waitForStart()
	} else {
		go waitForStart()
	}
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
			bfwd.log.Warn("failed to parse upstream", zap.Error(err))
			continue
		}
		upstreams[idx] = us
	}

	blockLists := []string{
		"https://adaway.org/hosts.txt",
		"https://dbl.oisd.nl/",
		"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
		"https://v.firebog.net/hosts/AdguardDNS.txt",
		"https://v.firebog.net/hosts/Easylist.txt",
		"https://adguardteam.github.io/AdGuardSDNSFilter/Filters/filter.txt",
	}
	if bll, ok := bfwd.c["blocklists"]; ok {
		lists := strings.Split(bll, ";")
		blockLists = lists
	}

	cfg := config.Config{}
	err := defaults.Set(&cfg)
	if err != nil {
		return fmt.Errorf("failed to set config defaults: %w", err)
	}
	// Blocky uses a custom registry, so this doesn't work as expected
	// cfg.Prometheus.Enable = true
	cfg.LogLevel = log.LevelDebug
	if !extconfig.Get().Debug {
		cfg.LogFormat = blockylog.FormatTypeJson
		// Only log errors from blocky to prevent double-logging all queries
		cfg.LogLevel = log.LevelFatal
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
			"block": blockLists,
		},
		ClientGroupsBlock: map[string][]string{
			"default": {"block"},
		},
	}
	bfwd.log.Debug("blocky config", zap.Any("config", cfg))

	srv, err := server.NewServer(&cfg)
	if err != nil {
		return fmt.Errorf("can't start server: %w", err)
	}
	bfwd.log.Debug("finished blocky setup", zap.Duration("took", time.Since(bfwd.st)))
	bfwd.b = srv
	return nil
}

func (bfwd *BlockyForwarder) Handle(w *utils.FakeDNSWriter, r *dns.Msg) *dns.Msg {
	if bfwd.b == nil {
		bfwd.log.Debug("Blocky not started/setup yet, falling through to next handler")
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
