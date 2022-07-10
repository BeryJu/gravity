package discovery

import (
	"context"
	"net/http"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/roles/discovery/types"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type DiscoveryRole struct {
	log *log.Entry
	i   roles.Instance
	cfg *DiscoveryRoleConfig
	ctx context.Context
}

func New(instance roles.Instance) *DiscoveryRole {
	r := &DiscoveryRole{
		log: instance.GetLogger().WithField("role", types.KeyRole),
		i:   instance,
	}
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		mux := ev.Payload.Data["mux"].(*mux.Router).Name("roles.discovery").Subrouter()
		mux.Name("v0.applyDevice").Path("/v0/discovery/apply").Methods(http.MethodPost).HandlerFunc(r.apiHandlerDeviceApply)
		mux.Name("v0.listDevices").Path("/v0/discovery/list").Methods(http.MethodGet).HandlerFunc(r.apiHandlerDeviceList)
		mux.Name("v0.listSubnets").Path("/v0/discovery/subnets").Methods(http.MethodGet).HandlerFunc(r.apiHandlerSubnetList)
		mux.Name("v0.listDevices").Path("/v0/discovery/subnets").Methods(http.MethodPost).HandlerFunc(r.apiHandlerDeviceList)
	})
	return r
}

func (r *DiscoveryRole) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeDiscoveryRoleConfig(config)
	if !r.cfg.Enabled || extconfig.Get().ListenOnlyMode {
		r.log.Info("Not enabling discovery")
		return nil
	}
	r.startWatchSubnets()
	return nil
}

func (r *DiscoveryRole) Stop() {
}
