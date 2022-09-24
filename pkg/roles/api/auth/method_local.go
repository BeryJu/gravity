package auth

import (
	"context"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/sessions"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"golang.org/x/crypto/bcrypt"
)

func (ap *AuthProvider) apiHandlerAuthLogin() usecase.Interactor {
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
		session.Values[types.SessionKeyUser] = user
		output.Successful = true
		return nil
	})
	u.SetName("api.login_user")
	u.SetTitle("API Users")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	u.SetExpectedErrors(status.Unauthenticated)
	return u
}
