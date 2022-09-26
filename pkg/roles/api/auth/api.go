package auth

import (
	"context"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/sessions"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
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

func (ap *AuthProvider) apiHandlerAuthMe() usecase.Interactor {
	type userMeOutput struct {
		Authenticated bool   `json:"authenticated" required:"true"`
		Username      string `json:"username" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *userMeOutput) error {
		session := ctx.Value(types.RequestSession).(*sessions.Session)
		u, ok := session.Values[types.SessionKeyUser]
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
