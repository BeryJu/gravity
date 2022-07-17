package dhcp

import (
	"beryju.io/gravity/pkg/roles"
)

func (r *Role) eventCreateLease(ev *roles.Event) {
	ident := ev.Payload.Data["mac"].(string)
	hostname := ev.Payload.Data["hostname"].(string)
	address := ev.Payload.Data["address"].(string)
	scopeName := ev.Payload.Data["scope"].(string)

	scope := r.scopes[scopeName]

	lease := &Lease{
		Identifier: ident,

		Hostname: hostname,
		Address:  address,
		ScopeKey: scope.Name,

		inst:  scope.inst,
		log:   scope.log.WithField("lease", ident),
		scope: scope,
	}
	lease.put(ev.Context, -1)
}
