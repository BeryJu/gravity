package auth

import (
	"net/http"
	"testing"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/stretchr/testify/assert"
)

func mustRequest(meth string, url string) *http.Request {
	req, err := http.NewRequest(meth, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestPermission_Fixed(t *testing.T) {
	ap := AuthProvider{}
	assert.True(t, ap.checkPermission(mustRequest("get", "/foo/bar"), &types.User{
		Permissions: []*types.Permission{
			{
				Path:    "/foo/bar",
				Methods: []string{"get", "post"},
			},
			{
				Path:    "/foo/ba",
				Methods: []string{"post"},
			},
			{
				Path:    "/foo",
				Methods: []string{"head"},
			},
		},
	}))
}

func TestPermission_Wildcard(t *testing.T) {
	ap := AuthProvider{}
	assert.True(t, ap.checkPermission(mustRequest("get", "/foo/bar"), &types.User{
		Permissions: []*types.Permission{
			{
				Path:    "/foo/*",
				Methods: []string{"get"},
			},
		},
	}))
	assert.True(t, ap.checkPermission(mustRequest("get", "/foo/bar"), &types.User{
		Permissions: []*types.Permission{
			{
				Path:    "/foo/*",
				Methods: []string{"*"},
			},
		},
	}))
}
