package migrate

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/storage"
	"github.com/Masterminds/semver/v3"
	clientv3 "go.etcd.io/etcd/client/v3"
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

func (mi *Migrator) GetClusterVersion(ctx context.Context) (*semver.Version, error) {
	type partialInstanceInfo struct {
		Version string `json:"version" required:"true"`
	}
	pfx := mi.ri.KV().Key(
		types.KeyInstance,
	).Prefix(true).String()
	instances, err := mi.ri.KV().Get(
		ctx,
		pfx,
		clientv3.WithPrefix(),
	)
	if err != nil {
		return nil, err
	}
	// Gather all instances in the cluster and parse their versions
	version := []*semver.Version{}
	for _, inst := range instances.Kvs {
		if strings.Count(strings.TrimPrefix(string(inst.Key), pfx), "/") > 0 {
			continue
		}
		pi := partialInstanceInfo{}
		err = json.Unmarshal(inst.Value, &pi)
		if err != nil {
			mi.log.Warn("failed to parse instance info", zap.Error(err))
			continue
		}
		v, err := semver.NewVersion(pi.Version)
		if err != nil {
			mi.log.Warn("failed to parse instance version", zap.Error(err))
			continue
		}
		version = append(version, v)
	}
	sort.Sort(semver.Collection(version))
	if len(version) < 1 {
		return semver.MustParse(extconfig.FullVersion()), nil
	}
	return version[0], nil
}

func (mi *Migrator) Run(ctx context.Context) (*storage.Client, error) {
	cv, err := mi.GetClusterVersion(ctx)
	if err != nil {
		return nil, err
	}
	mi.log.Debug("Checking migrations to activate for cluster version", zap.String("clusterVersion", cv.String()))
	cli := mi.ri.KV()
	for _, m := range mi.migrations {
		mi.log.Debug("Checking if migration needs to be run", zap.String("migration", m.Name()))
		enabled, err := m.Check(cv, ctx)
		if err != nil {
			mi.log.Warn("failed to check if migration should be enabled", zap.String("migration", m.Name()), zap.Error(err))
			return nil, err
		}
		if enabled {
			_cli, err := m.Hook(ctx)
			if err != nil {
				mi.log.Warn("failed to hook for migration", zap.String("migration", m.Name()), zap.Error(err))
				return nil, err
			}
			mi.log.Info("Enabling migration", zap.String("migration", m.Name()))
			cli = _cli
		} else {
			mi.log.Info("Running cleanup for migration", zap.String("migration", m.Name()))
			err := m.Cleanup(ctx)
			if err != nil {
				mi.log.Warn("failed to cleanup migration", zap.String("migration", m.Name()), zap.Error(err))
				continue
			}
		}
	}
	return cli, nil
}

func (mi *Migrator) AddMigration(migration roles.Migration) {
	mi.migrations = append(mi.migrations, migration)
}
