package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSetTrustedProxies_DoesNotPanic(t *testing.T) {
	r := gin.New()
	require.NotPanics(t, func() { SetTrustedProxies(r) })
}
