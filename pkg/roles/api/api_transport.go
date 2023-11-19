package api

import (
	"context"
	"encoding/base64"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	apiTypes "beryju.io/gravity/pkg/roles/api/types"
	tsdbTypes "beryju.io/gravity/pkg/roles/tsdb/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type APIExportInput struct {
	Safe bool `json:"safe"`
}
type APIExportOutput struct {
	Entries []APITransportEntry `json:"entries"`
}
type APITransportEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (r *Role) ignoredPrefixes() []string {
	return []string{
		// Sensitive data (tokens, hashed passwords, sessions)
		r.i.KV().Key(apiTypes.KeyRole, apiTypes.KeySessions).String(),
		r.i.KV().Key(apiTypes.KeyRole, apiTypes.KeyTokens).String(),
		r.i.KV().Key(apiTypes.KeyRole, apiTypes.KeyUsers).String(),
		// Noisy data we don't need
		r.i.KV().Key(tsdbTypes.KeyRole, tsdbTypes.KeySystem).String(),
	}
}

func (r *Role) APIClusterExport() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIExportInput, output *APIExportOutput) error {
		exps, err := r.i.KV().Get(ctx, "/", clientv3.WithPrefix())
		if err != nil {
			return err
		}
		output.Entries = make([]APITransportEntry, 0)
		for _, exp := range exps.Kvs {
			relKey := strings.TrimPrefix(string(exp.Key), extconfig.Get().Etcd.Prefix)
			shouldExport := true
			if input.Safe {
				for _, k := range r.ignoredPrefixes() {
					if strings.HasPrefix(strings.ToLower(relKey), k) {
						shouldExport = false
					}
				}
			}
			if shouldExport {
				output.Entries = append(output.Entries, APITransportEntry{
					Key:   relKey,
					Value: base64.StdEncoding.EncodeToString(exp.Value),
				})
			}
		}
		return nil
	})
	u.SetName("api.export")
	u.SetTitle("Export Cluster")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APIImportInput struct {
	Entries []APITransportEntry `json:"entries"`
}

func (r *Role) APIClusterImport() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIImportInput, output *struct{}) error {
		for _, entry := range input.Entries {
			val, err := base64.StdEncoding.DecodeString(entry.Value)
			if err != nil {
				r.log.Warn("failed to decode value", zap.Error(err), zap.String("key", entry.Key))
				continue
			}
			_, err = r.i.KV().Put(ctx, entry.Key, string(val))
			if err != nil {
				r.log.Warn("failed to put value", zap.Error(err), zap.String("key", entry.Key))
				continue
			}
		}
		return nil
	})
	u.SetName("api.import")
	u.SetTitle("Import Cluster")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}
