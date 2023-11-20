package instance

import (
	"context"

	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type APIRoleRestartInput struct {
	ID string `json:"roleId"`
}
type APIRoleRestartOutput struct{}

func (i *Instance) APIClusterRoleRestart() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIRoleRestartInput, output *APIRoleRestartOutput) error {
		i.rolesM.Lock()
		_, ok := i.roles[input.ID]
		i.rolesM.Unlock()
		if !ok {
			return status.NotFound
		}
		config, err := i.kv.KV.Get(
			ctx,
			i.kv.Key(
				types.KeyInstance,
				types.KeyRole,
				input.ID,
			).String(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		c := []byte{}
		if len(config.Kvs) > 0 {
			c = config.Kvs[0].Value
		}
		i.DispatchEvent(types.EventTopicRoleRestart, roles.NewEvent(ctx, map[string]interface{}{
			"id":     input.ID,
			"config": c,
		}))
		return nil
	})
	u.SetName("cluster.instance_role_restart")
	u.SetTitle("Instance roles")
	u.SetTags("cluster/instances")
	u.SetExpectedErrors(status.NotFound, status.Internal)
	return u
}
