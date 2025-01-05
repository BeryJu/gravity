package utils

import (
	"context"

	"github.com/miekg/dns"
)

type DNSRoutingMeta struct {
	HandlerIdx      int
	HasMoreHandlers bool
	ResolveRequest  func(w dns.ResponseWriter, r *DNSRequest)
}

type DNSRequest struct {
	*dns.Msg
	context context.Context
	meta    DNSRoutingMeta
	iter    int
}

func NewRequest(msg *dns.Msg, ctx context.Context, meta DNSRoutingMeta) *DNSRequest {
	return &DNSRequest{
		Msg:     msg,
		context: ctx,
		meta:    meta,
		iter:    0,
	}
}

func (r *DNSRequest) Context() context.Context {
	return r.context
}

func (r *DNSRequest) Meta() DNSRoutingMeta {
	return r.meta
}

func (r *DNSRequest) Iteration() int {
	return r.iter
}

func (r *DNSRequest) Chain(msg *dns.Msg, ctx context.Context, meta DNSRoutingMeta) *DNSRequest {
	return &DNSRequest{
		Msg:     msg,
		context: ctx,
		meta:    meta,
		iter:    r.iter + 1,
	}
}
