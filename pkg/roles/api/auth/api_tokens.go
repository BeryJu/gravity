package auth

import (
	"context"
	"encoding/base64"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/securecookie"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type APIToken struct {
	Key      string `json:"key" required:"true"`
	Username string `json:"username" required:"true"`
}
type APITokensGetOutput struct {
	Tokens []APIToken `json:"tokens" required:"true"`
}

func (ap *AuthProvider) APITokensGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APITokensGetOutput) error {
		rawTokens, err := ap.inst.KV().Get(
			ctx,
			ap.inst.KV().Key(
				types.KeyRole,
				types.KeyTokens,
			).Prefix(true).String(),
			clientv3.WithPrefix(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, rkey := range rawTokens.Kvs {
			u, err := ap.tokenFromKV(rkey)
			if err != nil {
				ap.log.Warn("failed to parse api key", zap.Error(err))
				continue
			}
			output.Tokens = append(output.Tokens, APIToken{
				Key:      u.Key,
				Username: u.Username,
			})
		}
		return nil
	})
	u.SetName("api.get_tokens")
	u.SetTitle("Tokens")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APITokensPutInput struct {
	Username string `query:"username" required:"true"`
}
type APITokensPutOutput struct {
	Key string `json:"key" required:"true"`
}

func (ap *AuthProvider) APITokensPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APITokensPutInput, output *APITokensPutOutput) error {
		token := &types.Token{
			Key:      base64.RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)),
			Username: input.Username,
		}
		err := ap.putToken(token, ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		output.Key = token.Key
		return nil
	})
	u.SetName("api.put_tokens")
	u.SetTitle("Tokens")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APITokensDeleteInput struct {
	Key string `query:"key" required:"true"`
}

func (ap *AuthProvider) APITokensDelete() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APITokensDeleteInput, output *struct{}) error {
		_, err := ap.inst.KV().Delete(ctx, ap.inst.KV().Key(
			types.KeyRole,
			types.KeyTokens,
			input.Key,
		).String())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("api.delete_tokens")
	u.SetTitle("Tokens")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}
