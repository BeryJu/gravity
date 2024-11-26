package oui_test

import (
	"bytes"
	"testing"

	"beryju.io/gravity/internal/resources"
	"beryju.io/gravity/pkg/roles/dhcp/oui"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	db := oui.New(bytes.NewBuffer([]byte(resources.MacOUIDB)))
	assert.NotNil(t, db)
	v, err := db.LookupString("00:50:56:ab:ec:bd")
	assert.NoError(t, err)
	assert.Equal(t, "VMware, Inc.", v.Organization)
}
