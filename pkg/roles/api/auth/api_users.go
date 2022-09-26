package auth

import (
	"context"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/crypto/bcrypt"
)

func (ap *AuthProvider) apiHandlerAuthUserGet() usecase.Interactor {
	type user struct {
		Username string `json:"username" required:"true"`
	}
	type authUsersOutput struct {
		Users []user `json:"users" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *authUsersOutput) error {
		rawUsers, err := ap.inst.KV().Get(
			ctx,
			ap.inst.KV().Key(
				types.KeyRole,
				types.KeyUsers,
			).Prefix(true).String(),
			clientv3.WithPrefix(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, ruser := range rawUsers.Kvs {
			u, err := ap.userFromKV(ruser)
			if err != nil {
				ap.log.WithError(err).Warning("failed to parse user")
				continue
			}
			output.Users = append(output.Users, user{
				Username: u.Username,
			})
		}
		return nil
	})
	u.SetName("api.get_users")
	u.SetTitle("API Users")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}

func (ap *AuthProvider) apiHandlerAuthUserPut() usecase.Interactor {
	type authUsersPut struct {
		Username string `query:"username" required:"true"`

		Password string `json:"password" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input authUsersPut, output *struct{}) error {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		user := &User{
			Username: input.Username,
			Password: string(hash),
			ap:       ap,
		}
		err = user.put(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("api.put_users")
	u.SetTitle("API Users")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}

func (ap *AuthProvider) apiHandlerAuthUserDelete() usecase.Interactor {
	type authUserDeleteInput struct {
		Username string `query:"username" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input authUserDeleteInput, output *struct{}) error {
		_, err := ap.inst.KV().Delete(ctx, ap.inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
			input.Username,
		).String())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("api.delete_users")
	u.SetTitle("API Users")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}
