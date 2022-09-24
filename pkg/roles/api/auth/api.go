package auth

import (
	"context"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/sessions"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"golang.org/x/crypto/bcrypt"
)

func (ap *AuthProvider) apiHandlerAuthUserMe() usecase.Interactor {
	type userMeOutput struct {
		Authenticated bool   `json:"authenticated"`
		Username      string `json:"username"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *userMeOutput) error {
		session := ctx.Value(types.RequestSession).(*sessions.Session)
		u, ok := session.Values[types.RequestKeyUser]
		if u == nil || !ok {
			output.Authenticated = false
			return nil
		}
		user := u.(User)
		output.Authenticated = true
		output.Username = user.Username
		return nil
	})
	u.SetName("api.users_me")
	u.SetTitle("API Users")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}

func (ap *AuthProvider) apiHandlerAuthUserLogin() usecase.Interactor {
	type userLoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type usersLoginOutput struct {
		Successful bool `json:"successful"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input *userLoginInput, output *usersLoginOutput) error {
		rawUsers, err := ap.inst.KV().Get(
			ctx,
			ap.inst.KV().Key(
				types.KeyRole,
				types.KeyUsers,
				input.Username,
			).String(),
		)
		if err != nil {
			bcrypt.CompareHashAndPassword([]byte{}, []byte(input.Password))
			ap.log.WithError(err).Warning("failed to get users")
			return status.Wrap(err, status.Internal)
		}
		if len(rawUsers.Kvs) < 1 {
			bcrypt.CompareHashAndPassword([]byte{}, []byte(input.Password))
			return status.Unauthenticated
		}
		user, err := ap.userFromKV(rawUsers.Kvs[0])
		if err != nil {
			bcrypt.CompareHashAndPassword([]byte{}, []byte(input.Password))
			ap.log.WithField("user", input.Username).WithError(err).Warning("failed to parse user")
			return status.Wrap(err, status.Internal)
		}
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
			ap.log.WithField("user", input.Username).Warning("invalid credentials")
			return status.Unauthenticated
		}
		session := ctx.Value(types.RequestSession).(*sessions.Session)
		session.Values[types.RequestKeyUser] = user
		output.Successful = true
		return nil
	})
	u.SetName("api.users_login")
	u.SetTitle("API Users")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	u.SetExpectedErrors(status.Unauthenticated)
	return u
}

func (ap *AuthProvider) apiHandlerAuthUserPut() usecase.Interactor {
	type authUsersPut struct {
		Username string `query:"username"`

		Password string `json:"password"`
	}
	u := usecase.NewIOI(new(authUsersPut), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*authUsersPut)
		)
		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		user := &User{
			Username: in.Username,
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

func (ap *AuthProvider) apiHandlerAuthUserRead() usecase.Interactor {
	type user struct {
		Username string `json:"username"`
	}
	type authUsersOutput struct {
		Users []user `json:"users"`
	}
	u := usecase.NewIOI(new(struct{}), new(authUsersOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*authUsersOutput)
		)
		rawUsers, err := ap.inst.KV().Get(
			ctx,
			ap.inst.KV().Key(
				types.KeyRole,
				types.KeyUsers,
			).Prefix(true).String(),
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
			out.Users = append(out.Users, user{
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

func (ap *AuthProvider) apiHandlerAuthUserDelete() usecase.Interactor {
	type authUserDeleteInput struct {
		Username string `query:"username"`
	}
	u := usecase.NewIOI(new(authUserDeleteInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*authUserDeleteInput)
		)
		_, err := ap.inst.KV().Delete(ctx, ap.inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
			in.Username,
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
