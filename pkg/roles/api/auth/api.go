package auth

import (
	"context"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/sessions"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"golang.org/x/crypto/bcrypt"
)

func (ap *AuthProvider) apiHandlerAuthConfig() usecase.Interactor {
	type authConfigOutput struct {
		Local bool `json:"bool"`
		OIDC  bool `json:"oidc"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *authConfigOutput) error {
		if ap.oidc != nil {
			output.OIDC = true
		}
		output.Local = true
		return nil
	})
	u.SetName("api.auth_config")
	u.SetTitle("API Users")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}

func (ap *AuthProvider) hasUsers(ctx context.Context) bool {
	rawUsers, err := ap.inst.KV().Get(
		ctx,
		ap.inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
		).Prefix(true).String(),
	)
	// Fallback to true to not give access when etcd request fails
	if err != nil {
		return true
	}
	return len(rawUsers.Kvs) > 0
}

func (ap *AuthProvider) apiHandlerAuthMe() usecase.Interactor {
	type userMeOutput struct {
		Authenticated bool   `json:"authenticated" required:"true"`
		Username      string `json:"username" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *userMeOutput) error {
		session := ctx.Value(types.RequestSession).(*sessions.Session)
		u, ok := session.Values[types.SessionKeyUser]
		if u == nil || !ok {
			if !ap.hasUsers(ctx) {
				session.Values[types.SessionKeyUser] = User{
					Username: "default-user",
					Password: "",
				}
				u = session.Values[types.SessionKeyUser]
			} else {
				output.Authenticated = false
				return nil
			}
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

func (ap *AuthProvider) apiHandlerAuthUserPut() usecase.Interactor {
	type authUsersPut struct {
		Username string `query:"username"`

		Password string `json:"password"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input authUsersPut, output *interface{}) error {
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

func (ap *AuthProvider) apiHandlerAuthUserRead() usecase.Interactor {
	type user struct {
		Username string `json:"username"`
	}
	type authUsersOutput struct {
		Users []user `json:"users"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *authUsersOutput) error {
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

func (ap *AuthProvider) apiHandlerAuthUserDelete() usecase.Interactor {
	type authUserDeleteInput struct {
		Username string `query:"username"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input authUserDeleteInput, output *interface{}) error {
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
