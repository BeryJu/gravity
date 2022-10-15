package auth

import (
	"net/http"
	"strings"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

func (ap *AuthProvider) checkToken(r *http.Request) bool {
	header := r.Header.Get(AuthorizationHeader)
	if header == "" {
		return false
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) < 2 {
		return false
	}
	if !strings.EqualFold(parts[0], BearerType) {
		return false
	}
	rawTokens, err := ap.inst.KV().Get(
		r.Context(),
		ap.inst.KV().Key(
			types.KeyRole,
			types.KeyTokens,
			parts[1],
		).String(),
	)
	if err != nil {
		ap.log.Warn("failed to check token", zap.Error(err))
		return false
	}
	if len(rawTokens.Kvs) < 1 {
		return false
	}
	key, err := ap.tokenFromKV(rawTokens.Kvs[0])
	if err != nil {
		return false
	}
	session := r.Context().Value(types.RequestSession).(*sessions.Session)
	session.Values[types.SessionKeyUser] = User{
		Username: key.Username,
	}
	session.Values[types.SessionKeyDirty] = true
	return false
}
