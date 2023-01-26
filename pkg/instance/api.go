package instance

import (
	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	"github.com/swaggest/rest/web"
)

func (i *Instance) setupInstanceAPI() {
	i.ForRole("instance").AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/cluster/instances", i.APIInstances())
		svc.Get("/api/v1/cluster/info", i.APIInstanceInfo())
		svc.Post("/api/v1/cluster/roles/restart", i.APIClusterRoleRestart())
	})
}
