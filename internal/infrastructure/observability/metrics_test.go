package observability

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type errCollector struct{}

func (e *errCollector) IncrementCounter(ctx context.Context, name string, labels map[string]string) error {
	return errors.New("fail")
}

func TestSetMetricsCollector_NilSetsNoop(t *testing.T) {
	SetMetricsCollector(nil)
	require.NotNil(t, defaultMetricsCollector)
}

func TestIncrementCounter_NoCollectorNoPanic(t *testing.T) {
	defaultMetricsCollector = nil
	require.NotPanics(t, func() {
		IncrementCounter(context.Background(), "m", nil)
	})
}

func TestIncrementCounter_ErrorPathDoesNotPanic(t *testing.T) {
	defaultMetricsCollector = &errCollector{}
	require.NotPanics(t, func() {
		IncrementCounter(context.Background(), "m", map[string]string{"a": "1"})
	})
}
