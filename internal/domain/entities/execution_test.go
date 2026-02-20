package entities

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecution_JSONTags(t *testing.T) {
	now := time.Now()
	exe := Execution{ID: "1", ServiceOrderID: "so", Status: "started", StartedAt: &now}
	b, err := json.Marshal(exe)
	require.NoError(t, err)
	require.Contains(t, string(b), "service_order_id")
}
