package dns

import (
	"fmt"
	"net/netip"
	"strings"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/creasty/defaults"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"

	"github.com/0xERR0R/blocky/config"
	blockylog "github.com/0xERR0R/blocky/log"
	"github.com/0xERR0R/blocky/server"
	"github.com/miekg/dns"
)

const BlockyForwarderType = "forward_blocky"

type BlockyForwarder struct {
	*IPForwarderHandler
	c   map[string]string
	b   *server.Server
	log *zap.Logger
	st  time.Time
	cfg *config.Config
}

func NewBlockyForwarder(z *Zone, rawConfig map[string]string) *BlockyForwarder {
	bfwd := &BlockyForwarder{
		IPForwarderHandler: NewIPForwarderHandler(z, rawConfig),
		c:                  rawConfig,
		st:                 time.Now(),
	}
	bfwd.log = z.log.With(zap.String("handler", bfwd.Identifier()))
	bfwd.startBlockyListServer()
	waitForStart := func() {
		err := bfwd.setup()
		if err != nil {
			bfwd.log.Warn("failed to setup blocky, queries will fallthrough", zap.Error(err))
		}
	}
	cfg, err := bfwd.getConfig()
	if err != nil {
		bfwd.log.Warn("Failed to build blocky config", zap.Error(err))
	}
	bfwd.cfg = cfg

	if extconfig.Get().Debug {
		bfwd.log.Info("starting blocky sync")
		waitForStart()
	} else {
		bfwd.log.Info("starting blocky async")
		go waitForStart()
	}
	return bfwd
}

func (bfwd *BlockyForwarder) Identifier() string {
	return BlockyForwarderType
}

func (bfwd *BlockyForwarder) getConfig() (*config.Config, error) {
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

	blockLists := []config.BytesSource{
		config.TextBytesSource(blockyListBase + "AdGuardSDNSFilter.txt"),
		config.TextBytesSource(blockyListBase + "AdguardDNS.txt"),
		config.TextBytesSource(blockyListBase + "Easylist.txt"),
		config.TextBytesSource(blockyListBase + "StevenBlack.hosts.txt"),
		config.TextBytesSource(blockyListBase + "adaway.org.txt"),
		config.TextBytesSource(blockyListBase + "dbl.oisd.nl.txt"),
	}
	if bll, ok := bfwd.c["blocklists"]; ok {
		lists := strings.Split(bll, ";")
		blockLists = []config.BytesSource{}
		for _, list := range lists {
			blockLists = append(blockLists, config.TextBytesSource(list))
		}
	}

	cfg := config.Config{}
	err := defaults.Set(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to set config defaults: %w", err)
	}
	// Blocky uses a custom registry, so this doesn't work as expected
	// cfg.Prometheus.Enable = true
	cfg.Log.Level = blockylog.LevelDebug
	cfg.QueryLog.Type = config.QueryLogTypeNone
	if !extconfig.Get().Debug {
		cfg.Log.Format = blockylog.FormatTypeJson
		// Only log errors from blocky to prevent double-logging all queries
		cfg.Log.Level = blockylog.LevelFatal
	}
	bootstrap, err := netip.ParseAddrPort(extconfig.Get().FallbackDNS)
	if err != nil {
		return nil, err
	}
	cfg.BootstrapDNS = []config.BootstrappedUpstreamConfig{
		{
			Upstream: config.Upstream{
				Net:  config.NetProtocolTcpUdp,
				Host: bootstrap.Addr().String(),
				Port: bootstrap.Port(),
			},
		},
	}
	cfg.Upstreams = config.Upstreams{
		Groups: map[string][]config.Upstream{
			"default": upstreams,
		},
		Timeout: config.Duration(types.DefaultUpstreamTimeout),
	}
	// TODO: Blocky config
	cfg.Blocking = config.Blocking{
		BlockType: "zeroIP",
		BlackLists: map[string][]config.BytesSource{
			"block": {
				config.TextBytesSource(),
			},
		},
		Loading: config.SourceLoadingConfig{
			Downloads: config.DownloaderConfig{
				Attempts: 3,
			},
		},
		ClientGroupsBlock: map[string][]string{
			"default": {"block"},
		},
	}
	return &cfg, nil
}

func (bfwd *BlockyForwarder) setup() error {
	bfwd.log.Debug("blocky config", zap.Any("config", bfwd.cfg))

	srv, err := server.NewServer(bfwd.z.inst.Context(), bfwd.cfg)
	if err != nil {
		return fmt.Errorf("can't start server: %w", err)
	}
	bfwd.log.Info("finished blocky setup", zap.Duration("took", time.Since(bfwd.st)))
	bfwd.b = srv
	return nil
}

func (bfwd *BlockyForwarder) Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg {
	if bfwd.b == nil {
		bfwd.log.Debug("Blocky not started/setup yet, falling through to next handler")
		return nil
	}
	bs := sentry.TransactionFromContext(r.Context()).StartChild("gravity.dns.handler.forward_blocky.handle")
	bfwd.b.OnRequest(r.Context(), w, r.Msg)
	bs.Finish()
	// fall to next handler when no record is found
	if w.Msg().Rcode == dns.RcodeNameError {
		return nil
	}
	for _, query := range r.Question {
		for idx, ans := range w.Msg().Answer {
			go bfwd.cacheToEtcd(r, query, ans, idx)
		}
	}
	return w.Msg()
}
