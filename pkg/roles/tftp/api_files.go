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
	Name      string `json:"name" required:"true"`
	Host      string `json:"host" required:"true"`
	SizeBytes int    `json:"sizeBytes" required:"true"`
}
type APIFilesGetOutput struct {
	Files []APIFile `json:"files" required:"true"`
}

func (r *Role) APIFilesGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIFilesGetOutput) error {
		prefix := r.i.KV().Key(
			types.KeyRole,
			types.KeyFiles,
		).Prefix(true).String()
		rawFiles, err := r.i.KV().Get(ctx, prefix, clientv3.WithPrefix())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, rf := range rawFiles.Kvs {
			parts := strings.SplitN(strings.TrimPrefix(string(rf.Key), prefix), "/", 2)
			output.Files = append(output.Files, APIFile{
				Host:      parts[0],
				Name:      parts[1],
				SizeBytes: len(rf.Value),
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

type APIFilesPutInput struct {
	APIFile
	Data string `json:"data" required:"true"`
}

func (r *Role) APIFilesPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIFilesPutInput, output *struct{}) error {
		_, err := r.i.KV().Put(
			ctx,
			r.i.KV().Key(
				types.KeyRole,
				types.KeyFiles,
				input.Host,
				input.Name,
			).String(),
			input.Data,
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("tftp.put_files")
	u.SetTitle("TFTP Files")
	u.SetTags("roles/tftp")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APIFilesDeleteInput struct {
	Host string `query:"host"`
	Name string `query:"name"`
}

func (r *Role) APIFilesDelete() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIFilesDeleteInput, output *struct{}) error {
		_, err := r.i.KV().Delete(
			ctx,
			r.i.KV().Key(
				types.KeyRole,
				types.KeyFiles,
				input.Host,
				input.Name,
			).String(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("tftp.delete_files")
	u.SetTitle("TFTP Files")
	u.SetTags("roles/tftp")
	u.SetExpectedErrors(status.Internal)
	return u
}
