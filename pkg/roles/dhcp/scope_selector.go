package dhcp

import (
	"net"

	"go.uber.org/zap"
)

type scopeSelector func(scope *Scope) int

func (r *Role) findScopeForRequest(req *Request4, additionalSelectors ...scopeSelector) *Scope {
	var match *Scope
	longestBits := 0
	r.scopesM.RLock()
	defer r.scopesM.RUnlock()
	// To prioritise requests from a DHCP relay being matched correctly, give their subnet
	// match a 1 bit more priority
	const dhcpRelayBias = 1
	for _, scope := range r.scopes {
		// Check additional selectors (highest priority)
		for _, sel := range additionalSelectors {
			m := sel(scope)
			if m > -1 && m > longestBits {
				match = scope
				longestBits = m
			}
		}
		// Check based on gateway IP (next highest priority)
		gatewayMatchBits := scope.match(req.GatewayIPAddr)
		if gatewayMatchBits > -1 && gatewayMatchBits+dhcpRelayBias > longestBits {
			req.log.Debug("selected scope based on cidr match (gateway IP)", zap.String("scope", scope.Name))
			match = scope
			longestBits = gatewayMatchBits + dhcpRelayBias
		}
		// Handle local broadcast, check with the instance's listening IP
		// Only consider local scopes if we don't have a match already
		localMatchBits := scope.match(net.ParseIP(req.LocalIP()))
		if localMatchBits > -1 && localMatchBits > longestBits {
			req.log.Debug("selected scope based on cidr match (instance/interface IP)", zap.String("scope", scope.Name))
			match = scope
			longestBits = localMatchBits
		}
		// Fallback to default scope if we don't already have a match
		if match == nil && scope.Default {
			req.log.Debug("selected scope based on default flag", zap.String("scope", scope.Name))
			match = scope
		}
	}
	if match != nil {
		req.log.Debug("final scope selection", zap.String("scope", match.Name))
	}
	return match
}
