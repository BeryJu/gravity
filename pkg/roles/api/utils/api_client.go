package utils

import (
	"context"

	"beryju.io/gravity/api"
)

type contextKey string

const ContextKeyAPIConfig contextKey = "api-config"

func APIClientFromRequest(ctx context.Context) *api.APIClient {
	rcfg := ctx.Value(ContextKeyAPIConfig)
	if rcfg == nil {
		return nil
	}
	cfg, ok := rcfg.(*api.Configuration)
	if !ok {
		return nil
	}
	return api.NewAPIClient(cfg)
}
