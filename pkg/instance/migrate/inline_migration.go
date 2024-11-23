package migrate

import (
	"context"

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
	HookFunc          func(context.Context)
	CleanupFunc       func(context.Context)
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

func (im *InlineMigration) Hook(context.Context) error {
	return nil
}

func (im *InlineMigration) Cleanup(context.Context) error {
	return nil
}
