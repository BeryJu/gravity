package externaldns

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"beryju.io/gravity/pkg/externaldns/generated/externaldnsapi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAdjustRecords(t *testing.T) {
	// Adjustment must work without a configured Gravity API or loaded zones.
	server := &Server{log: zap.NewNop()}
	router := externaldnsapi.NewRouter(externaldnsapi.NewUpdateAPIController(server))
	endpoints := `[{"dnsName":"app.example.com","targets":["192.0.2.1","192.0.2.2"],"recordType":"A","recordTTL":300,"labels":{"owner":"external-dns"}}]`
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/adjustendpoints", strings.NewReader(endpoints)))

	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, endpoints, response.Body.String())
}

func TestSetRecordsResponse(t *testing.T) {
	router := externaldnsapi.NewRouter(externaldnsapi.NewUpdateAPIController(&Server{}))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/records", strings.NewReader(`{}`)))

	assert.Equal(t, http.StatusNoContent, response.Code)
	assert.Empty(t, response.Body.String())
}
