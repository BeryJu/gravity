package dns

import (
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
)

type Handler interface {
	Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg
	Identifier() string
}
