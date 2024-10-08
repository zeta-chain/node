package memo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_NewArg(t *testing.T) {
	argAddress := sample.EthAddress()
	argString := sample.String()
	argBytes := sample.Bytes()

	tests := []struct {
		name    string
		argType string
		arg     interface{}
	}{
		{
			name:    "receiver",
			argType: "address",
			arg:     &argAddress,
		},
		{
			name:    "payload",
			argType: "bytes",
			arg:     &argBytes,
		},
		{
			name:    "revertAddress",
			argType: "string",
			arg:     &argString,
		},
		{
			name:    "abortAddress",
			argType: "address",
			arg:     &argAddress,
		},
		{
			name:    "revertMessage",
			argType: "bytes",
			arg:     &argBytes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg := memo.NewArg(tt.name, memo.ArgType(tt.argType), tt.arg)

			require.Equal(t, tt.name, arg.Name)
			require.Equal(t, memo.ArgType(tt.argType), arg.Type)
			require.Equal(t, tt.arg, arg.Arg)

			switch tt.name {
			case "receiver":
				require.Equal(t, arg, memo.ArgReceiver(&argAddress))
			case "payload":
				require.Equal(t, arg, memo.ArgPayload(&argBytes))
			case "revertAddress":
				require.Equal(t, arg, memo.ArgRevertAddress(&argString))
			case "abortAddress":
				require.Equal(t, arg, memo.ArgAbortAddress(&argAddress))
			case "revertMessage":
				require.Equal(t, arg, memo.ArgRevertMessage(&argBytes))
			}
		})
	}
}
