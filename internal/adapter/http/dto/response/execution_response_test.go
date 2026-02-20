package response

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecutionResponse_JSONTags(t *testing.T) {
	now := time.Now()
	r := ExecutionResponse{ID: "1", ServiceOrderID: "so", Status: "started", StartedAt: &now}
	b, err := json.Marshal(r)
	require.NoError(t, err)
	require.Contains(t, string(b), "service_order_id")
}
