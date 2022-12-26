package ms_dhcp_test

import (
	"testing"

	"beryju.io/gravity/pkg/convert/ms_dhcp"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/tests"
	"beryju.io/gravity/pkg/tests/api"
	"github.com/stretchr/testify/assert"
)

func TestDHCPImport(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	// Create DHCP role to register API routes
	dhcp.New(rootInst.ForRole("dhcp"))

	files := []string{
		"./test_a.xml",
		"./test_b.xml",
		"./test_c.xml",
	}

	api, stop := api.APIClient(rootInst)
	defer stop()

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			c, err := ms_dhcp.New(api, file)
			assert.NoError(t, err)
			errors := c.Run(ctx)
			assert.Equal(t, []error{}, errors)
		})
	}
}
