package discovery

import (
	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	apitypes "beryju.io/ddet/pkg/roles/api/types"
	"beryju.io/ddet/pkg/roles/discovery/types"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type DiscoveryRole struct {
	log *log.Entry
	i   roles.Instance
	cfg *DiscoveryRoleConfig
}

func New(instance roles.Instance) *DiscoveryRole {
	r := &DiscoveryRole{
		log: instance.GetLogger().WithField("role", types.KeyRole),
		i:   instance,
	}
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		mux := ev.Payload.Data["mux"].(*mux.Router).Name("roles.discovery").Subrouter()
		mux.Name("v0.applyDevice").Path("/api/v0/discovery/apply").Methods("POST").HandlerFunc(r.apiHandlerApply)
	})
	return r
}

func (r *DiscoveryRole) Start(config []byte) error {
	r.cfg = r.decodeDiscoveryRoleConfig(config)
	if !r.cfg.Enabled || extconfig.Get().ListenOnlyMode {
		return nil
	}
	r.startWatchSubnets()
	return nil
}

func (r *DiscoveryRole) Stop() {
}
