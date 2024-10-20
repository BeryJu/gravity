package tftp

import (
	"context"
	"strings"

	"beryju.io/gravity/pkg/roles/tftp/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type APIFile struct {
	Name string `json:"name" required:"true"`
	Host string `json:"host" required:"true"`
}
type APIFilesOutput struct {
	Files []APIFile `json:"files" required:"true"`
}

func (r *Role) APIFilesGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIFilesOutput) error {
		prefix := r.i.KV().Key(
			types.KeyRole,
			types.KeyFiles,
		).Prefix(true).String()
		rawFiles, err := r.i.KV().Get(ctx, prefix, clientv3.WithKeysOnly())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, rf := range rawFiles.Kvs {
			parts := strings.SplitN(strings.TrimSuffix(string(rf.Key), prefix), "/", 1)
			output.Files = append(output.Files, APIFile{
				Host: parts[0],
				Name: parts[1],
			})
		}
		return nil
	})
	u.SetName("tftp.get_files")
	u.SetTitle("TFTP Files")
	u.SetTags("roles/tftp")
	u.SetExpectedErrors(status.Internal)
	return u
}
