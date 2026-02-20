package observability

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRelicApp_MissingLicense(t *testing.T) {
	os.Unsetenv("NEW_RELIC_LICENSE_KEY")
	app, err := NewRelicApp()
	require.Nil(t, app)
	require.Error(t, err)
}
