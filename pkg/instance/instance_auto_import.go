package instance

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles/api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (i *Instance) autoImportConfig(ctx context.Context) error {
	for _, p := range extconfig.Get().ImportConfigs {
		var c []byte
		id := p
		if strings.HasPrefix(p, "file://") {
			u, err := url.Parse(p)
			if err != nil {
				return errors.Wrapf(err, "failed to read config %s", id)
			}
			_c, err := os.ReadFile(u.Path)
			if err != nil {
				return errors.Wrapf(err, "failed to read config %s", id)
			}
			c = _c
		} else {
			id = "inline"
			dec, err := base64.StdEncoding.DecodeString(p)
			if err != nil {
				return errors.Wrapf(err, "failed to read config %s", id)
			}
			c = dec
		}
		err := i.importSingleConfig(ctx, c)
		if err != nil {
			return errors.Wrapf(err, "failed to import config %s", id)
		}
		i.log.Info("Successfully imported config", zap.String("path", id))
	}
	return nil
}

func (i *Instance) importSingleConfig(ctx context.Context, c []byte) error {
	var input api.APIImportInput
	err := json.Unmarshal(c, &input)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}
	values := make([]string, len(input.Entries))
	for index, entry := range input.Entries {
		val, err := base64.StdEncoding.DecodeString(entry.Value)
		if err != nil {
			return fmt.Errorf("invalid value for config key %s: %w", entry.Key, err)
		}
		values[index] = string(val)
	}
	for index, entry := range input.Entries {
		// Startup owns the setup marker and its temporary lock keys.
		if key := i.kv.Key(types.KeyRole, types.KeyCluster).String(); entry.Key == key || strings.HasPrefix(entry.Key, key+"/") {
			continue
		}
		_, err := i.kv.Put(ctx, entry.Key, values[index])
		if err != nil {
			return errors.Wrapf(err, "failed to put config key %s", entry.Key)
		}
	}
	return nil
}
