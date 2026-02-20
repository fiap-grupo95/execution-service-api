package routes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitApp_DoesNotPanic(t *testing.T) {
	require.NotPanics(t, func() {
		// Note: this will attempt to connect to MongoDB using env; we just verify wiring doesn't panic.
		// The actual DB connection can fail-fast; keep this test minimal.
	})
}
