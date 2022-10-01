package tests

import (
	"context"
	"net"
	"time"
)

func DNSLookup(query string, to string) []string {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * 5,
			}
			return d.DialContext(ctx, network, to)
		},
	}
	addrs, err := r.LookupHost(context.Background(), query)
	if err != nil {
		panic(err)
	}
	return addrs
}
