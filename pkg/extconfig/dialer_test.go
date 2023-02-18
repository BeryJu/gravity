package extconfig

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolver(t *testing.T) {
	globalExtConfig.FallbackDNS = "127.0.0.1:1053"
	r := Resolver()
	addr, err := r.LookupHost(context.Background(), "gravity.beryju.io.")
	assert.NoError(t, err)
	assert.Equal(t, []string{"10.0.0.1"}, addr)
}
