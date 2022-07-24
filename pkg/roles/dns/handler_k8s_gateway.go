package dns

import (
	"context"

	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
	k8s_gateway "github.com/ori-edge/k8s_gateway"
	log "github.com/sirupsen/logrus"
)

type K8sGateway struct {
	c   map[string]string
	log *log.Entry
	gw  k8s_gateway.Gateway
}

func NewK8sGateway(z *Zone, rawConfig map[string]string) *K8sGateway {
	k8gw := &K8sGateway{
		c: rawConfig,
		gw: k8s_gateway.Gateway{
			Zones: []string{z.Name},
		},
	}
	k8gw.log = z.log.WithField("handler", k8gw.Identifier())
	k8gw.gw.RunKubeController(context.Background())
	k8gw.gw.ExternalAddrFunc = k8gw.gw.SelfAddress
	return k8gw
}

func (k8gw *K8sGateway) Identifier() string {
	return "k8s_gateway"
}

func (k8gw *K8sGateway) Handle(w *utils.FakeDNSWriter, r *dns.Msg) *dns.Msg {
	if !k8gw.gw.Controller.HasSynced() {
		k8gw.log.Info("K8s Gateway not synced yet, falling through to next handler")
		return nil
	}
	k8gw.gw.ServeDNS(context.TODO(), w, r)
	// fall to next handler when no record is found
	if w.Msg().Rcode == dns.RcodeNameError {
		return nil
	}
	return w.Msg()
}
