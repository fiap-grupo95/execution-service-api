package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type dummyHandler struct{}

func (d *dummyHandler) CreateExecution(c *gin.Context) { c.Status(http.StatusCreated) }
func (d *dummyHandler) FinishExecution(c *gin.Context) { c.Status(http.StatusOK) }
func (d *dummyHandler) CancelExecution(c *gin.Context) { c.Status(http.StatusOK) }

func TestAddExecutionRoutes_Registers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// We're in the same package, so we can call the unexported function.
	addExecutionRoutes(r, &dummyHandler{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/execution", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/v1/execution/finish/exe_1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/v1/execution/cancel/exe_1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
