package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareChain_CanAttach(t *testing.T) {
	r := gin.New()
	require.NotPanics(t, func() {
		r.Use(gin.Logger())
	})
}
