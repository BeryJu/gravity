package forward

import "github.com/miekg/dns"

type Forwarder interface {
	Handle(w dns.ResponseWriter, r *dns.Msg)
}
