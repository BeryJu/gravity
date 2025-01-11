package api

import (
	"context"
	"net/http"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/roles/api/utils"
)

// This middleware injects an API config...in the request context...to call the API on behalf
// of the request we just got
func NewAPIConfigMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := api.NewConfiguration()
			c.AddDefaultHeader("Authorization", r.Header.Get("Authorization"))
			c.AddDefaultHeader("Cookie", r.Header.Get("Cookie"))
			c.Host = r.URL.Host
			c.Scheme = r.URL.Scheme
			c.UserAgent = r.Header.Get("User-Agent")
			nr := r.WithContext(context.WithValue(r.Context(), utils.ContextKeyAPIConfig, c))
			h.ServeHTTP(w, nr)
		})
	}
}
