package observability

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartSegment_NoTxnNoPanic(t *testing.T) {
	end := StartSegment(nil, "seg")
	require.NotNil(t, end)
	require.NotPanics(t, func() { end() })
}
