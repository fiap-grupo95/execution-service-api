package logs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger_NotNil(t *testing.T) {
	require.NotNil(t, Logger())
}
