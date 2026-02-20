package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoopCollector_ReturnsNil(t *testing.T) {
	c := &noopMetricsCollector{}
	require.NoError(t, c.IncrementCounter(context.Background(), "m", nil))
}
