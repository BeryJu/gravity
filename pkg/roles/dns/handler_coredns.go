package dns

import (
	"reflect"
	"unsafe"

	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"
	"github.com/miekg/dns"
	_ "github.com/ori-edge/k8s_gateway"
	"go.uber.org/zap"
)

const CoreDNSType = "coredns"

type CoreDNS struct {
	c        map[string]interface{}
	log      *zap.Logger
	instance *caddy.Instance
	srv      *dnsserver.Server
}

func init() {
	HandlerRegistry.Add(CoreDNSType, func(z *Zone, rawConfig map[string]interface{}) Handler {
		return NewCoreDNS(z, rawConfig)
	})
}

func NewCoreDNS(z *Zone, rawConfig map[string]interface{}) *CoreDNS {
	core := &CoreDNS{
		c: rawConfig,
	}
	core.log = z.log.With(zap.String("handler", core.Identifier()))
	dnsserver.Quiet = true
	corefile := caddy.CaddyfileInput{
		Contents:       []byte(core.c["config"].(string)),
		Filepath:       "in-memory",
		ServerTypeName: "dns",
	}
	core.log.Info("starting coredns", zap.String("version", coremain.CoreVersion))
	instance, err := caddy.Start(corefile)
	if err != nil {
		core.log.Warn("failed to start codedns", zap.Error(err))
		return core
	}
	core.instance = instance
	rawSrv := core.instance.Servers()
	if len(rawSrv) < 1 {
		core.log.Warn("No server configured")
		return core
	}
	for _, srv := range rawSrv {
		rs := reflect.ValueOf(&srv).Elem()
		rf := rs.FieldByName("server")
		rfs := reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem().Interface().(caddy.Server)
		if dss, ok := rfs.(*dnsserver.Server); ok {
			core.srv = dss
			return core
		}
	}
	core.log.Warn("no compatible server found")
	return core
}

func (core *CoreDNS) Identifier() string {
	return CoreDNSType
}

func (core *CoreDNS) Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg {
	if core.instance == nil {
		return nil
	}
	core.srv.ServeDNS(r.Context(), w, r.Msg)
	// fall to next handler when no record is found
	if w.Msg().Rcode == dns.RcodeNameError {
		return nil
	}
	return w.Msg()
}
