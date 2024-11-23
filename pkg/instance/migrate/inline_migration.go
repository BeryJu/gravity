package migrate

import (
	"context"

	"beryju.io/gravity/pkg/storage"
	"github.com/Masterminds/semver/v3"
)

func MustParseConstraint(input string) *semver.Constraints {
	c, err := semver.NewConstraint(input)
	if err != nil {
		panic(err)
	}
	return c
}

type InlineMigration struct {
	MigrationName     string
	ActivateOnVersion *semver.Constraints
	HookFunc          func(context.Context) (*storage.Client, error)
	CleanupFunc       func(context.Context) error
}

func (im *InlineMigration) Name() string {
	return im.MigrationName
}

func (im *InlineMigration) Check(clusterVersion *semver.Version, ctx context.Context) (bool, error) {
	if im.ActivateOnVersion.Check(clusterVersion) {
		return true, nil
	}
	return false, nil
}

func (im *InlineMigration) Hook(ctx context.Context) (*storage.Client, error) {
	return im.HookFunc(ctx)
}

func (im *InlineMigration) Cleanup(ctx context.Context) error {
	if im.CleanupFunc != nil {
		return im.CleanupFunc(ctx)
	}
	return nil
}
