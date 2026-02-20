package request

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecutionRequest_JSONTags(t *testing.T) {
	r := ExecutionRequest{ID: "1", ServiceOrderID: "so"}
	b, err := json.Marshal(r)
	require.NoError(t, err)
	require.Contains(t, string(b), "service_order_id")
}
