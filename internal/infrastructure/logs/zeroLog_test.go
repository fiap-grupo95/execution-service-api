package logs

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitAndHelpers_DoNotPanic(t *testing.T) {
	Init(nil)
	require.NotNil(t, Logger())

	require.NotPanics(t, func() { Info("i") })
	require.NotPanics(t, func() { Debug("d") })
	require.NotPanics(t, func() { Warn("w") })
	require.NotPanics(t, func() { Error("e", errors.New("x")) })
	require.NotPanics(t, func() { LoggerWithContext(nil) })
	require.NotNil(t, LoggerWithContext(context.Background()))
}
