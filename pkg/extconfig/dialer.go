package extconfig

import (
	"context"
	"fmt"
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

func Transport() http.RoundTripper {
	dialer := &net.Dialer{
		Resolver: Resolver(),
	}
	return NewUserAgentTransport(fmt.Sprintf("gravity/%s", FullVersion()), &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial:  dialer.Dial,
	})
}

type userAgentTransport struct {
	inner http.RoundTripper
	ua    string
}

func NewUserAgentTransport(ua string, inner http.RoundTripper) *userAgentTransport {
	return &userAgentTransport{inner, ua}
}

func (uat *userAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", uat.ua)
	return uat.inner.RoundTrip(r)
}
