package externaldns

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/externaldns/generated/externaldnsapi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAdjustRecords(t *testing.T) {
	for name, endpoints := range map[string]string{
		"records": `[
			{
				"dnsName": "app.example.com",
				"targets": ["192.0.2.1", "192.0.2.2"],
				"recordType": "A",
				"recordTTL": 300,
				"setIdentifier": "app",
				"labels": {"owner": "external-dns"},
				"providerSpecific": [{"name": "custom", "value": "preserved"}]
			},
			{
				"dnsName": "app.example.com",
				"targets": ["2001:db8::1"],
				"recordType": "AAAA",
				"recordTTL": 300
			},
			{
				"dnsName": "txt.app.example.com",
				"targets": ["heritage=external-dns,external-dns/owner=test"],
				"recordType": "TXT"
			}
		]`,
		"empty": `[]`,
	} {
		t.Run(name, func(t *testing.T) {
			var requests atomic.Int32
			backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{}`))
			}))
			defer backend.Close()

			config := api.NewConfiguration()
			config.Servers[0].URL = backend.URL
			server := &Server{
				api:   api.NewAPIClient(config),
				log:   zap.NewNop(),
				zones: []api.DnsAPIZone{{Name: "example.com."}},
			}
			router := externaldnsapi.NewRouter(externaldnsapi.NewUpdateAPIController(server))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/adjustendpoints", strings.NewReader(endpoints)))

			assert.Equal(t, http.StatusOK, response.Code)
			assert.JSONEq(t, endpoints, response.Body.String())
			assert.Zero(t, requests.Load(), "adjustment must not call the Gravity API")
		})
	}
}

func TestSetRecords(t *testing.T) {
	for _, test := range []struct {
		name          string
		backendStatus int
		wantStatus    int
	}{
		{"success", http.StatusOK, http.StatusNoContent},
		{"backend error", http.StatusInternalServerError, http.StatusInternalServerError},
	} {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/api/v1/dns/zones/records", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(test.backendStatus)
				_, _ = w.Write([]byte(`{}`))
			}))
			defer backend.Close()

			config := api.NewConfiguration()
			config.Servers[0].URL = backend.URL
			server := &Server{
				api:   api.NewAPIClient(config),
				log:   zap.NewNop(),
				zones: []api.DnsAPIZone{{Name: "example.com."}},
			}
			router := externaldnsapi.NewRouter(externaldnsapi.NewUpdateAPIController(server))
			response := httptest.NewRecorder()
			changes := `{"create":[{"dnsName":"app.example.com","targets":["192.0.2.1"],"recordType":"A"}]}`
			router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/records", strings.NewReader(changes)))

			assert.Equal(t, test.wantStatus, response.Code)
			assert.EqualValues(t, 1, requests.Load())
			if test.wantStatus == http.StatusNoContent {
				assert.Empty(t, response.Body.String())
			}
		})
	}
}
