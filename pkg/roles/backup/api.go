package backup

import (
	"context"

	"github.com/swaggest/usecase"
)

func (r *BackupRole) apiHandlerBackupStart() usecase.Interactor {
	type backupStartInput struct {
		Wait bool `query:"wait"`
	}
	u := usecase.NewIOI(new(backupStartInput), new(BackupStatus), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*backupStartInput)
			out = output.(*BackupStatus)
		)
		if in.Wait {
			out = r.saveSnapshot()
		} else {
			go r.saveSnapshot()
			out.Status = BackupStatusStarted
		}
		return nil
	})
	u.SetName("backup.start")
	u.SetTitle("Backup start")
	u.SetTags("roles/backup")
	u.SetDescription("Start a backup.")
	return u
}
