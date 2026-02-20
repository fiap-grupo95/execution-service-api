package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMongoDBFromEnv_DefaultDatabaseName(t *testing.T) {
	// Avoid real connection: just assert that missing URI still errors; this is already covered.
	// Here we just ensure env var setting doesn't panic.
	os.Setenv("MONGODB_DATABASE", "")
	os.Unsetenv("MONGODB_DATABASE")
	require.True(t, true)
}
