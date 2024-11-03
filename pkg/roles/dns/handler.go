package dns

import (
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
)

type Handler interface {
	Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg
	Identifier() string
}

type HandlerConstructor func(z *Zone, rawConfig map[string]interface{}) Handler

var HandlerRegistry = newRegistry()

type handlerRegistry struct {
	handlers map[string]HandlerConstructor
}

func newRegistry() handlerRegistry {
	return handlerRegistry{
		handlers: make(map[string]HandlerConstructor),
	}
}

func (hn handlerRegistry) Add(identifier string, h HandlerConstructor) {
	hn.handlers[identifier] = h
}

func (hn handlerRegistry) Find(identifier string) (HandlerConstructor, bool) {
	c, ok := hn.handlers[identifier]
	return c, ok
}
