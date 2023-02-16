package utils

import (
	"context"

	"github.com/miekg/dns"
)

type DNSRequest struct {
	*dns.Msg
	context context.Context
}

func NewRequest(msg *dns.Msg, ctx context.Context) *DNSRequest {
	return &DNSRequest{
		Msg:     msg,
		context: ctx,
	}
}

func (r *DNSRequest) Context() context.Context {
	return r.context
}
