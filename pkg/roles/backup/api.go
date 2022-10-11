package backup

import (
	"context"
	"encoding/json"
	"strings"

	"beryju.io/gravity/pkg/roles/backup/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type APIBackupStartInput struct {
	Wait bool `query:"wait" required:"true"`
}

func (r *Role) APIBackupStart() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIBackupStartInput, output *BackupStatus) error {
		if input.Wait {
			o := r.SaveSnapshot()
			output.Duration = o.Duration
			output.Error = o.Error
			output.Filename = o.Filename
			output.Size = o.Size
			output.Status = o.Status
		} else {
			go r.SaveSnapshot()
			output.Status = BackupStatusStarted
		}
		return nil
	})
	u.SetName("backup.start")
	u.SetTitle("Backup start")
	u.SetTags("roles/backup")
	return u
}

type APIBackupStatusOutput struct {
	Status map[string]BackupStatus `json:"status" required:"true"`
}

func (r *Role) APIBackupStatus() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIBackupStatusOutput) error {
		prefix := r.i.KV().Key(
			types.KeyRole,
			types.KeyStatus,
		).Prefix(true).String()
		rawStatus, err := r.i.KV().Get(
			r.ctx,
			prefix,
			clientv3.WithPrefix(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		output.Status = make(map[string]BackupStatus)
		for _, rs := range rawStatus.Kvs {
			var s BackupStatus
			err := json.Unmarshal(rs.Value, &s)
			if err != nil {
				r.log.WithError(err).Warning("failed to unmarshal status")
				continue
			}
			keyParts := strings.Split(strings.TrimPrefix(string(rs.Key), prefix), "/")
			output.Status[keyParts[0]] = s
		}
		return nil
	})
	u.SetName("backup.status")
	u.SetTitle("Backup status")
	u.SetTags("roles/backup")
	u.SetExpectedErrors(status.Internal)
	return u
}
