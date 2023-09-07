package dhcp

import (
	"beryju.io/gravity/pkg/roles"
	"go.uber.org/zap"
)

func (r *Role) eventCreateLease(ev *roles.Event) {
	ident := ev.Payload.Data["mac"].(string)
	hostname := ev.Payload.Data["hostname"].(string)
	address := ev.Payload.Data["address"].(string)
	scopeName := ev.Payload.Data["scope"].(string)

	r.scopesM.RLock()
	scope := r.scopes[scopeName]
	r.scopesM.RUnlock()
	if scope == nil {
		r.log.Warn("event to create lease with missing scope", zap.String("scopeName", scopeName))
		return
	}
	lease := &Lease{
		Identifier: ident,

		Hostname: hostname,
		Address:  address,
		ScopeKey: scope.Name,

		inst:  scope.inst,
		log:   scope.log.With(zap.String("lease", ident)),
		scope: scope,
	}
	err := lease.Put(ev.Context, -1)
	if err != nil {
		r.log.Warn("failed to put lease in event handler", zap.Error(err))
	}
}
