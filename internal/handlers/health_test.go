package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ubcesports/echo-base/internal/tests"
)

func TestHealthCheck(t *testing.T) {
	tests.SetupTestDBForTest(t)
	handler := http.HandlerFunc(HealthCheck)

	rr := tests.ExecuteTestRequest(t, handler, http.MethodGet, "/health", nil)

	var response HealthResponse
	tests.AssertResponse(t, rr, http.StatusOK, &response)

	require.Equal(t, "ok", response.Status)
	require.Equal(t, "ok", response.Database)
}
