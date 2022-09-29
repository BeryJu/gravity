package backup

import (
	"context"

	"github.com/swaggest/usecase"
)

func (r *Role) apiHandlerBackupStart() usecase.Interactor {
	type backupStartInput struct {
		Wait bool `query:"wait" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input backupStartInput, output *BackupStatus) error {
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
