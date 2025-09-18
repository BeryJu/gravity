package auth

import (
	"context"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type APIUsersGetInput struct {
	Username string `query:"username" description:"Optional username of a user to get"`
}
type APIUser struct {
	Username    string              `json:"username" required:"true"`
	Permissions []*types.Permission `json:"permissions" required:"true"`
}
type APIUsersGetOutput struct {
	Users []APIUser `json:"users" required:"true"`
}

func (ap *AuthProvider) APIUsersGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIUsersGetInput, output *APIUsersGetOutput) error {
		key := ap.inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
		)
		if input.Username == "" {
			key = key.Prefix(true)
		} else {
			key = key.Add(input.Username)
		}
		rawUsers, err := ap.inst.KV().Get(
			ctx,
			key.String(),
			clientv3.WithPrefix(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, ruser := range rawUsers.Kvs {
			u, err := ap.userFromKV(ruser)
			if err != nil {
				ap.log.Warn("failed to parse user", zap.Error(err))
				continue
			}
			output.Users = append(output.Users, APIUser{
				Username:    u.Username,
				Permissions: u.Permissions,
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

	Password    string              `json:"password" required:"true"`
	Permissions []*types.Permission `json:"permissions" required:"true"`
}

func (ap *AuthProvider) APIUsersPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIUsersPutInput, output *struct{}) error {
		rawUsers, err := ap.inst.KV().Get(
			ctx,
			ap.inst.KV().Key(
				types.KeyRole,
				types.KeyUsers,
				input.Username,
			).String(),
		)
		var oldUser *types.User
		if err == nil && len(rawUsers.Kvs) > 0 {
			user, err := ap.userFromKV(rawUsers.Kvs[0])
			if err != nil {
				_ = bcrypt.CompareHashAndPassword([]byte{}, []byte(input.Password))
				ap.log.Warn("failed to parse user", zap.Error(err), zap.String("user", input.Username))
				return status.Wrap(err, status.Internal)
			}
			oldUser = user
		}

		user := &types.User{
			Username:    input.Username,
			Permissions: input.Permissions,
		}
		if input.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
			if err != nil {
				return status.Wrap(err, status.Internal)
			}
			user.Password = string(hash)
		} else if oldUser != nil {
			user.Password = oldUser.Password
		}
		err = ap.putUser(user, ctx)
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
