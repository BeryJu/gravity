package backup

import (
	"context"
	"encoding/json"
	"strings"

	"beryju.io/gravity/pkg/roles/backup/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type APIBackupStartInput struct {
	Wait bool `query:"wait" required:"true"`
}

func (r *Role) APIBackupStart() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIBackupStartInput, output *BackupStatus) error {
		if input.Wait {
			o := r.SaveSnapshot(ctx)
			output.Duration = o.Duration
			output.Error = o.Error
			output.Filename = o.Filename
			output.Size = o.Size
			output.Status = o.Status
		} else {
			go r.SaveSnapshot(ctx)
			output.Status = BackupStatusStarted
		}
		return nil
	})
	u.SetName("backup.start")
	u.SetTitle("Backup start")
	u.SetTags("roles/backup")
	return u
}

type APIBackupStatus struct {
	BackupStatus
	Node string `json:"node"`
}

type APIBackupStatusOutput struct {
	Status []APIBackupStatus `json:"status" required:"true"`
}

func (r *Role) APIBackupStatus() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIBackupStatusOutput) error {
		prefix := r.i.KV().Key(
			types.KeyRole,
			types.KeyStatus,
		).Prefix(true).String()
		rawStatus, err := r.i.KV().Get(
			ctx,
			prefix,
			clientv3.WithPrefix(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, rs := range rawStatus.Kvs {
			var s BackupStatus
			err := json.Unmarshal(rs.Value, &s)
			if err != nil {
				r.log.Warn("failed to unmarshal status", zap.Error(err))
				continue
			}
			keyParts := strings.Split(strings.TrimPrefix(string(rs.Key), prefix), "/")
			output.Status = append(output.Status, APIBackupStatus{
				BackupStatus: s,
				Node:         keyParts[0],
			})
		}
		return nil
	})
	u.SetName("backup.status")
	u.SetTitle("Backup status")
	u.SetTags("roles/backup")
	u.SetExpectedErrors(status.Internal)
	return u
}
