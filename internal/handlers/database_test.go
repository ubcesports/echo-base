package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ubcesports/echo-base/internal/tests"
)

func TestDatabasePing(t *testing.T) {
	tests.SetupTestDBForTest(t)
	handler := http.HandlerFunc(DatabasePing)

	rr := tests.ExecuteTestRequest(t, handler, "GET", "/db/ping", nil)

	var response DatabasePingResponse
	tests.AssertResponse(t, rr, http.StatusOK, &response)
	require.Equal(t, "ok", response.Status, "Expected status 'ok'")
	require.NotEmpty(t, response.ResponseTime, "Response time should not be empty")
}
