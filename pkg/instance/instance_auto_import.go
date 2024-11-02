package instance

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"os"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles/api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (i *Instance) autoImportConfig() {
	for _, p := range extconfig.Get().ImportConfigs {
		var c []byte
		id := p
		if strings.HasPrefix(p, "file://") {
			u, err := url.Parse(p)
			if err != nil {
				i.log.Warn("failed to import config", zap.String("path", id), zap.Error(err))
				continue
			}
			_c, err := os.ReadFile(u.Path)
			if err != nil {
				i.log.Warn("failed to import config", zap.String("path", id), zap.Error(err))
				continue
			}
			c = _c
		} else {
			id := p[20:]
			dec, err := base64.StdEncoding.DecodeString(p)
			if err != nil {
				i.log.Warn("failed to import config", zap.String("path", id), zap.Error(err))
				continue
			}
			c = dec
		}
		err := i.importSingleConfig(c)
		if err != nil {
			i.log.Warn("failed to import config", zap.String("path", id), zap.Error(err))
		} else {
			i.log.Info("Successfully imported config", zap.String("path", id))
		}
	}
}

func (i *Instance) importSingleConfig(c []byte) error {
	var input api.APIImportInput
	err := json.Unmarshal(c, &input)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}
	inst := i.ForRole("root", context.Background())
	for _, entry := range input.Entries {
		val, err := base64.StdEncoding.DecodeString(entry.Value)
		if err != nil {
			i.log.Warn("failed to decode value", zap.Error(err), zap.String("key", entry.Key))
			continue
		}
		_, err = inst.KV().Put(context.Background(), entry.Key, string(val))
		if err != nil {
			i.log.Warn("failed to put value", zap.Error(err), zap.String("key", entry.Key))
			continue
		}
	}
	return nil
}
