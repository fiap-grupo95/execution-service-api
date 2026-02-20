package logs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggerWithContext_NotNil(t *testing.T) {
	require.NotNil(t, LoggerWithContext(context.Background()))
}
