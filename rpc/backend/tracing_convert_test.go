package backend

// import (
// 	"encoding/json"
// 	"testing"

// 	"github.com/stretchr/testify/require"
// 	rpctypes "github.com/zeta-chain/node/rpc/types"
// )

// const expectedConvertedValue = `{"onlyTopCall":false}`

// const brokenStyleRaw = `
// {
//     "tracer": "callTracer",
//     "tracerConfig": "{\"onlyTopCall\":false}"
// }
// `

// const compliantStyleRaw = `
// {
//     "tracer": "callTracer",
//     "tracerConfig": {"onlyTopCall":false}
// }
// `

// const emptyStyleRaw = `
// {
//     "tracer": "callTracer"
// }
// `

// const invalidStyleRaw = `
// {
//     "tracer": "callTracer",
// 	"tracerConfig": []
// }
// `

// func TestConvertConfig(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		input       string
// 		expectError bool
// 		expected    string
// 	}{
// 		{
// 			name:        "broken style",
// 			input:       brokenStyleRaw,
// 			expectError: false,
// 			expected:    expectedConvertedValue,
// 		},
// 		{
// 			name:        "compliant style",
// 			input:       compliantStyleRaw,
// 			expectError: false,
// 			expected:    expectedConvertedValue,
// 		},
// 		{
// 			name:        "empty style",
// 			input:       emptyStyleRaw,
// 			expectError: false,
// 			expected:    "",
// 		},
// 		{
// 			name:        "invalid style",
// 			input:       invalidStyleRaw,
// 			expectError: true,
// 			expected:    "",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rpcConfig := &rpctypes.TraceConfig{}
// 			err := json.Unmarshal([]byte(tt.input), rpcConfig)
// 			require.NoError(t, err)

// 			ethermintConfig, err := convertConfig(rpcConfig)
// 			if tt.expectError {
// 				require.Error(t, err)
// 				return
// 			}
// 			require.NoError(t, err)
// 			require.Equal(t, tt.expected, ethermintConfig.TracerJsonConfig)
// 		})
// 	}
// }
