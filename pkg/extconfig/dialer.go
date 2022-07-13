package extconfig

import (
	"context"
	"net"
	"net/http"
)

func Resolver() *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, network, Get().FallbackDNS)
		},
	}
}

func Transport() *http.Transport {
	dialer := &net.Dialer{
		Resolver: Resolver(),
	}
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial:  dialer.Dial,
	}
}
