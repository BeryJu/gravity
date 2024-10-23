package auth

import (
	"net/http"
	"strings"
)

const wildcard = "*"

func (ap *AuthProvider) checkPermission(req *http.Request, u User) bool {
	var longestMatch *Permission
	for _, perm := range u.Permissions() {
		if strings.HasSuffix(perm.Path, wildcard) && strings.HasPrefix(req.URL.Path, strings.TrimSuffix(perm.Path, wildcard)) {
			if longestMatch == nil || len(perm.Path) > len(longestMatch.Path) {
				longestMatch = &perm
			}
		} else if perm.Path == req.URL.Path {
			if longestMatch == nil || len(perm.Path) > len(longestMatch.Path) {
				longestMatch = &perm
			}
		}
	}
	if longestMatch == nil {
		return false
	}
	for _, meth := range longestMatch.Methods {
		if strings.EqualFold(meth, req.Method) {
			return true
		}
	}
	return false
}
