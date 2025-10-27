package instance

import (
	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	"github.com/swaggest/rest/web"
)

func (i *Instance) setupInstanceAPI() {
	i.ForRole("instance", i.rootContext).AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/cluster", i.APIClusterInfo())
		svc.Get("/api/v1/cluster/instance", i.APIInstanceGet())
		svc.Put("/api/v1/cluster/instance", i.APIInstancePut())
		svc.Post("/api/v1/cluster/roles/restart", i.APIClusterRoleRestart())
	})
}
