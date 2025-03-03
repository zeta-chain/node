package backend

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	rpctypes "github.com/zeta-chain/node/rpc/types"
)

const expectedConvertedValue = `{"onlyTopCall":false}`

const brokenStyleRaw = `
{
    "tracer": "callTracer",
    "tracerConfig": "{\"onlyTopCall\":false}"
}
`

const compliantStyleRaw = `
{
    "tracer": "callTracer",
    "tracerConfig": {"onlyTopCall":false}
}
`

// TestConvertConfig ensures that both the old broken style and the new compliant style
// serialize to the same value
func TestConvertConfig(t *testing.T) {
	brokenStyle := &rpctypes.TraceConfig{}
	err := json.Unmarshal([]byte(brokenStyleRaw), brokenStyle)
	require.NoError(t, err)

	brokenStyleConverted := convertConfig(brokenStyle)
	require.Equal(t, expectedConvertedValue, brokenStyleConverted.TracerJsonConfig)

	compliantStyle := &rpctypes.TraceConfig{}
	err = json.Unmarshal([]byte(compliantStyleRaw), compliantStyle)
	require.NoError(t, err)

	compliantStyleConverted := convertConfig(compliantStyle)
	require.Equal(t, expectedConvertedValue, compliantStyleConverted.TracerJsonConfig)
}
