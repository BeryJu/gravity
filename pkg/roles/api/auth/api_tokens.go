package auth

import (
	"context"
	"encoding/base64"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/securecookie"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (ap *AuthProvider) apiHandlerAuthTokenGet() usecase.Interactor {
	type token struct {
		Key      string `json:"key" required:"true"`
		Username string `json:"username" required:"true"`
	}
	type authTokensOutput struct {
		Tokens []token `json:"tokens" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *authTokensOutput) error {
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
				ap.log.WithError(err).Warning("failed to parse api key")
				continue
			}
			output.Tokens = append(output.Tokens, token{
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

func (ap *AuthProvider) apiHandlerAuthTokenPut() usecase.Interactor {
	type authTokensPutInput struct {
		Username string `query:"username" required:"true"`
	}
	type authTokensPutOutput struct {
		Key string `json:"key" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input authTokensPutInput, output *authTokensPutOutput) error {
		token := &Token{
			Key:      base64.RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)),
			Username: input.Username,
			ap:       ap,
		}
		err := token.put(ctx)
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

func (ap *AuthProvider) apiHandlerAuthTokenDelete() usecase.Interactor {
	type authTokensDeleteInput struct {
		Key string `query:"key" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input authTokensDeleteInput, output *struct{}) error {
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
