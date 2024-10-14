package utils

import (
	"context"

	"github.com/miekg/dns"
)

type DNSRoutingMeta struct {
	HandlerIdx      int
	HasMoreHandlers bool
}

type DNSRequest struct {
	*dns.Msg
	context context.Context
	meta    DNSRoutingMeta
}

func NewRequest(msg *dns.Msg, ctx context.Context, meta DNSRoutingMeta) *DNSRequest {
	return &DNSRequest{
		Msg:     msg,
		context: ctx,
		meta:    meta,
	}
}

func (r *DNSRequest) Context() context.Context {
	return r.context
}

func (r *DNSRequest) Meta() DNSRoutingMeta {
	return r.meta
}
