package dhcp

import (
	"fmt"

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
	fmt.Println(scope.log)
	lease := &Lease{
		Identifier: ident,

		Hostname: hostname,
		Address:  address,
		ScopeKey: scope.Name,

		inst:  scope.inst,
		log:   scope.log.With(zap.String("lease", ident)),
		scope: scope,
	}
	lease.Put(ev.Context, -1)
}
