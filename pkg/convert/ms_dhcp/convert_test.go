package ms_dhcp_test

import (
	"os"
	"testing"

	"beryju.io/gravity/pkg/convert/ms_dhcp"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestDHCPImport(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	// Create DHCP role to register API routes
	dhcp.New(rootInst.ForRole("dhcp", ctx))

	files := []string{
		"./fixtures/test_a.xml",
		"./fixtures/test_b.xml",
		"./fixtures/test_c.xml",
	}

	api, stop := tests.APIClient(rootInst)
	defer stop()

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			x, err := os.ReadFile(file)
			assert.NoError(t, err)
			c, err := ms_dhcp.New(api, string(x))
			assert.NoError(t, err)
			assert.NoError(t, c.Run(ctx))
		})
	}
}
