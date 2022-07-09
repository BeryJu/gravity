package discovery

import "net/http"

func (ro *DiscoveryRole) apiHandlerApply(w http.ResponseWriter, r *http.Request) {
	ro.log.Debug("test")

}
