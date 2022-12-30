package dhcp

import (
	"bytes"

	"beryju.io/gravity/internal/macoui"
	"beryju.io/gravity/pkg/roles/dhcp/oui"
	"go.uber.org/zap"
)

func (r *Role) initOUI() {
	db := oui.New(bytes.NewBuffer([]byte(macoui.DB)))
	r.oui = db
	r.log.Info("loaded OUI database", zap.Int("size", len(macoui.DB)))
}
