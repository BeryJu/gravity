package migrate

import (
	"context"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/storage"
	"github.com/Masterminds/semver/v3"
	"go.uber.org/zap"
)

type Migrator struct {
	ri         roles.Instance
	log        *zap.Logger
	migrations []roles.Migration
}

func New(ri roles.Instance) *Migrator {
	return &Migrator{
		ri:         ri,
		log:        ri.Log().Named("migrator"),
		migrations: make([]roles.Migration, 0),
	}
}

func (mi *Migrator) Run(ctx context.Context) (*storage.Client, error) {
	cv := semver.MustParse(extconfig.FullVersion())
	cli := mi.ri.KV()
	for _, m := range mi.migrations {
		enabled, err := m.Check(cv, ctx)
		if err != nil {
			mi.log.Warn("failed to check if migration should be enabled", zap.String("migration", m.Name()), zap.Error(err))
			return nil, err
		}
		if !enabled {
			continue
		}
		_cli, err := m.Hook(ctx)
		if err != nil {
			mi.log.Warn("failed to hook for migration", zap.String("migration", m.Name()), zap.Error(err))
			return nil, err
		}
		cli = _cli
	}
	return cli, nil
}

func (mi *Migrator) AddMigration(migration roles.Migration) {
	mi.migrations = append(mi.migrations, migration)
}
