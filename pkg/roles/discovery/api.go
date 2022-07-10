package discovery

import (
	"net/http"

	"beryju.io/ddet/pkg/roles"
	"github.com/gorilla/mux"
)

func (ro *DiscoveryRole) eventHandlerAPIMux(ev *roles.Event) {
	mux := ev.Payload.Data["mux"].(*mux.Router).Name("roles.discovery").Subrouter()
	mux.Name("v0.applyDevice").Path("/v0/discovery/apply").Methods(http.MethodPost).HandlerFunc(ro.apiHandlerDeviceApply)
	mux.Name("v0.listDevices").Path("/v0/discovery/list").Methods(http.MethodGet).HandlerFunc(ro.apiHandlerDeviceList)
	mux.Name("v0.listSubnets").Path("/v0/discovery/subnets").Methods(http.MethodGet).HandlerFunc(ro.apiHandlerSubnetList)
	mux.Name("v0.listDevices").Path("/v0/discovery/subnets").Methods(http.MethodPost).HandlerFunc(ro.apiHandlerDeviceList)
}
