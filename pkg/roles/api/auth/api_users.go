package auth

import (
	"context"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/crypto/bcrypt"
)

type APIUser struct {
	Username string `json:"username" required:"true"`
}
type APIUsersGetOutput struct {
	Users []APIUser `json:"users" required:"true"`
}

func (ap *AuthProvider) APIUsersGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIUsersGetOutput) error {
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
			output.Users = append(output.Users, APIUser{
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

type APIUsersPutInput struct {
	Username string `query:"username" required:"true"`

	Password string `json:"password" required:"true"`
}

func (ap *AuthProvider) APIUsersPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIUsersPutInput, output *struct{}) error {
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

type APIUsersDeleteInput struct {
	Username string `query:"username" required:"true"`
}

func (ap *AuthProvider) APIUsersDelete() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIUsersDeleteInput, output *struct{}) error {
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
