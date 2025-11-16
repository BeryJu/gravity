package instance

import (
	"context"
	"sync"

	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"github.com/getsentry/sentry-go"
)

func (i *Instance) bootstrap(ctx context.Context) {
	i.log.Debug("bootstrapping instance")
	i.keepAliveInstanceInfo(ctx)
	i.setupInstanceAPI()
	rootInstance := i.ForRole("root", ctx)
	for _, roleId := range i.getRoles(ctx) {
		instanceRoles.WithLabelValues(roleId).Add(1)
		rctx, cancel := context.WithCancelCause(i.rootContext)
		rc := RoleContext{
			RoleInstance:      i.ForRole(roleId, rctx),
			ContextCancelFunc: cancel,
		}
		switch roleId {
		case "etcd":
			// Special handling
			continue
		default:
			span := sentry.StartSpan(ctx, "gravity.instance.bootstrap.role")
			span.SetTag("gravity.role", roleId)
			rc.Role = roles.GetRole(roleId)(rc.RoleInstance)
			span.Finish()
		}
		i.rolesM.Lock()
		i.roles[roleId] = rc
		i.rolesM.Unlock()
	}
	rootInstance.AddEventListener(types.EventTopicRoleRestart, i.eventRoleRestart)
	rootInstance.DispatchEvent(
		types.EventTopicInstanceBootstrapped,
		roles.NewEvent(i.rootContext, map[string]interface{}{}),
	)
	i.checkFirstStart(ctx)
	wg := sync.WaitGroup{}
	for roleId := range i.roles {
		wg.Add(1)
		go i.startWatchRole(ctx, roleId, func() {
			wg.Done()
		})
	}
	go func() {
		wg.Wait()
		i.DispatchEvent(types.EventTopicRolesStarted, roles.NewEvent(ctx, map[string]interface{}{}))
		sentry.TransactionFromContext(ctx).Finish()
	}()
}
