package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMongoDBFromEnv_MissingURI(t *testing.T) {
	t.Setenv("MONGODB_URI", "")
	os.Unsetenv("MONGODB_URI")

	cfg, err := NewMongoDBFromEnv()
	require.Nil(t, cfg)
	require.Error(t, err)
}
