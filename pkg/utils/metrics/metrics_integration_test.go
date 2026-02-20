package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildMetricName_MultipleLabels(t *testing.T) {
	name := BuildMetricName("x", map[string]string{"z": "9", "a": "1", "m": "5"})
	require.Equal(t, "x[a=1,m=5,z=9]", name)
}
