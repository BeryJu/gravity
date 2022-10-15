package auth

import (
	"context"
	"net/http"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/sessions"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type APILoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type APILoginOutput struct {
	Successful bool `json:"successful"`
}

func (ap *AuthProvider) APILogin() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input *APILoginInput, output *APILoginOutput) error {
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
			ap.log.Warn("failed to get users", zap.Error(err))
			return status.Wrap(err, status.Internal)
		}
		if len(rawUsers.Kvs) < 1 {
			bcrypt.CompareHashAndPassword([]byte{}, []byte(input.Password))
			return status.Unauthenticated
		}
		user, err := ap.userFromKV(rawUsers.Kvs[0])
		if err != nil {
			bcrypt.CompareHashAndPassword([]byte{}, []byte(input.Password))
			ap.log.Warn("failed to parse user", zap.Error(err), zap.String("user", input.Username))
			return status.Wrap(err, status.Internal)
		}
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
			ap.log.Warn("invalid credentials", zap.String("user", input.Username))
			return status.Unauthenticated
		}
		session := ctx.Value(types.RequestSession).(*sessions.Session)
		session.Values[types.SessionKeyUser] = user
		session.Values[types.SessionKeyDirty] = true
		output.Successful = true
		return nil
	})
	u.SetName("api.login_user")
	u.SetTitle("API Users")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal, status.Unauthenticated)
	return u
}

func (ap *AuthProvider) APILogout(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(types.RequestSession).(*sessions.Session)
	session.Values[types.SessionKeyUser] = nil
	session.Values[types.SessionKeyDirty] = true
	http.Redirect(w, r, "/", http.StatusFound)
}
