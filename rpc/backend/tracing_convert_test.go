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

const emptyStyleRaw = `
{
    "tracer": "callTracer"
}
`

func TestConvertConfig(t *testing.T) {
	t.Run("broken style", func(t *testing.T) {
		brokenStyle := &rpctypes.TraceConfig{}
		err := json.Unmarshal([]byte(brokenStyleRaw), brokenStyle)
		require.NoError(t, err)
		brokenStyleConverted := convertConfig(brokenStyle)
		require.Equal(t, expectedConvertedValue, brokenStyleConverted.TracerJsonConfig)
	})

	t.Run("compliant style", func(t *testing.T) {
		compliantStyle := &rpctypes.TraceConfig{}
		err := json.Unmarshal([]byte(compliantStyleRaw), compliantStyle)
		require.NoError(t, err)
		compliantStyleConverted := convertConfig(compliantStyle)
		require.Equal(t, expectedConvertedValue, compliantStyleConverted.TracerJsonConfig)
	})

	t.Run("empty style", func(t *testing.T) {
		emptyStyle := &rpctypes.TraceConfig{}
		err := json.Unmarshal([]byte(emptyStyleRaw), emptyStyle)
		require.NoError(t, err)
		emptyStyleConverted := convertConfig(emptyStyle)
		require.Equal(t, "", emptyStyleConverted.TracerJsonConfig)
	})
}
