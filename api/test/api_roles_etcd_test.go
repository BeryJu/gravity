/*
gravity

Testing RolesEtcdApiService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package api

import (
	"context"
	"testing"

	openapiclient "beryju.io/gravity/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_api_RolesEtcdApiService(t *testing.T) {
	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)

	t.Run("Test RolesEtcdApiService EtcdJoinMember", func(t *testing.T) {
		t.Skip("skip test") // remove to run test

		resp, httpRes, err := apiClient.RolesEtcdApi.EtcdJoinMember(context.Background()).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)
	})
}
