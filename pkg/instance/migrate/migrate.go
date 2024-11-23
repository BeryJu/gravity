package migrate

import (
	"context"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
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

func (mi *Migrator) Run(ctx context.Context) error {
	cv := semver.MustParse(extconfig.FullVersion())
	for _, m := range mi.migrations {
		enabled, err := m.Check(cv, ctx)
		if err != nil {
			mi.log.Warn("failed to check if migration should be enabled", zap.String("migration", m.Name()), zap.Error(err))
			continue
		}
		if !enabled {
			continue
		}
		err = m.Hook(ctx)
		if err != nil {
			mi.log.Warn("failed to hook for migration", zap.String("migration", m.Name()), zap.Error(err))
			continue
		}
	}
	return nil
}

func (mi *Migrator) AddMigration(migration roles.Migration) {
	mi.migrations = append(mi.migrations, migration)
}
