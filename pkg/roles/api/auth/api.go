package auth

import (
	"context"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/sessions"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type APIConfigOutput struct {
	Local bool `json:"bool"`
	OIDC  bool `json:"oidc"`
}

func (ap *AuthProvider) APIConfig() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIConfigOutput) error {
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

type APIMeOutput struct {
	Username      string              `json:"username" required:"true"`
	Authenticated bool                `json:"authenticated" required:"true"`
	Permissions   []*types.Permission `json:"permissions" required:"true"`
}

func (ap *AuthProvider) APIMe() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIMeOutput) error {
		session := ctx.Value(types.RequestSession).(*sessions.Session)
		u, ok := session.Values[types.SessionKeyUser]
		if u == nil || !ok {
			output.Authenticated = false
			return nil
		}
		user := u.(*types.User)
		output.Authenticated = true
		output.Username = user.Username
		output.Permissions = user.Permissions
		return nil
	})
	u.SetName("api.users_me")
	u.SetTitle("API Users")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}
