package api

import (
	"context"
	"encoding/base64"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type APIExportOutput struct {
	Entries []APITransportEntry `json:"entries"`
}
type APITransportEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (r *Role) APIClusterExport() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIExportOutput) error {
		exps, err := r.i.KV().Get(ctx, "/", clientv3.WithPrefix())
		if err != nil {
			return err
		}
		output.Entries = make([]APITransportEntry, len(exps.Kvs))
		for idx, exp := range exps.Kvs {
			output.Entries[idx] = APITransportEntry{
				Key:   strings.TrimPrefix(string(exp.Key), extconfig.Get().Etcd.Prefix),
				Value: base64.StdEncoding.EncodeToString(exp.Value),
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
