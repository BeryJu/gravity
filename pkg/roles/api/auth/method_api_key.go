package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/sessions"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

const (
	AuthorizationHeader = "Authorization"
	BearerType          = "bearer"
)

type APIKey struct {
	Username string `json:"username"`
}

func (ap *AuthProvider) apiKeyFromKV(raw *mvccpb.KeyValue) (*APIKey, error) {
	apiKey := &APIKey{}
	err := json.Unmarshal(raw.Value, &apiKey)
	if err != nil {
		return apiKey, err
	}
	return apiKey, nil
}

func (ap *AuthProvider) checkAPIKey(r *http.Request) bool {
	header := r.Header.Get(AuthorizationHeader)
	if header == "" {
		return false
	}
	parts := strings.SplitN(header, " ", 1)
	if len(parts) < 2 {
		return false
	}
	if !strings.EqualFold(parts[0], BearerType) {
		return false
	}
	rawKeys, err := ap.inst.KV().Get(
		r.Context(),
		ap.inst.KV().Key(
			types.KeyRole,
			types.KeyAPIKeys,
			parts[1],
		).String(),
	)
	if err != nil {
		ap.log.WithError(err).Warning("failed to check API keys")
		return false
	}
	if len(rawKeys.Kvs) < 1 {
		return false
	}
	key, err := ap.apiKeyFromKV(rawKeys.Kvs[0])
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
