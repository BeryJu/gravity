package dhcp

import (
	"bytes"

	"beryju.io/gravity/internal/resources"
	"beryju.io/gravity/pkg/roles/dhcp/oui"
	"go.uber.org/zap"
)

func (r *Role) initOUI() {
	db := oui.New(bytes.NewBuffer([]byte(resources.MacOUIDB)))
	r.oui = db
	r.log.Info("loaded OUI database", zap.Int("size", len(resources.MacOUIDB)))
}
