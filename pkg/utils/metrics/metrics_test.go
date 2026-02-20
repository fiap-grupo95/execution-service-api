package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildMetricName_NoLabels(t *testing.T) {
	require.Equal(t, "base", BuildMetricName("base", nil))
	require.Equal(t, "base", BuildMetricName("base", map[string]string{}))
}

func TestBuildMetricName_SortsLabels(t *testing.T) {
	name := BuildMetricName("m", map[string]string{"b": "2", "a": "1"})
	require.Equal(t, "m[a=1,b=2]", name)
}
