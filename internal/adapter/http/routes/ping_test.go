package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestPingHandler_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	addPingRoutes(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, PathHealthCheck, nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "pong")
}
