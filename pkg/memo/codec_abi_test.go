package memo_test

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/testutil/sample"
)

// newArgInstance creates a new instance of the given argument type
func newArgInstance(v interface{}) interface{} {
	switch v.(type) {
	case common.Address:
		return new(common.Address)
	case []byte:
		return &[]byte{}
	case string:
		return new(string)
	}
	return nil
}

// ensureArgEquality ensures the expected argument and actual value are equal
func ensureArgEquality(t *testing.T, expected, actual interface{}) {
	switch v := expected.(type) {
	case common.Address:
		require.Equal(t, v.Hex(), actual.(*common.Address).Hex())
	case []byte:
		require.True(t, bytes.Equal(v, *actual.(*[]byte)))
	case string:
		require.Equal(t, v, *actual.(*string))
	default:
		t.Fatalf("unexpected argument type: %T", v)
	}
}

func Test_NewCodecABI(t *testing.T) {
	c := memo.NewCodecABI()
	require.NotNil(t, c)
}

func Test_CodecABI_AddArguments(t *testing.T) {
	codec := memo.NewCodecABI()
	require.NotNil(t, codec)

	address := sample.EthAddress()
	codec.AddArguments(memo.ArgReceiver(&address))
}

func Test_CodecABI_PackArgument(t *testing.T) {
	// create sample arguments
	argAddress := sample.EthAddress()
	argBytes := []byte("some test bytes argument")
	argString := "some test string argument"

	// test cases
	tests := []struct {
		name   string
		args   []memo.CodecArg
		errMsg string
	}{
		{
			name: "pack in the order of [address, bytes, string]",
			args: []memo.CodecArg{
				memo.ArgReceiver(argAddress),
				memo.ArgPayload(argBytes),
				memo.ArgRevertAddress(argString),
			},
		},
		{
			name: "pack in the order of [string, address, bytes]",
			args: []memo.CodecArg{
				memo.ArgRevertAddress(argString),
				memo.ArgReceiver(argAddress),
				memo.ArgPayload(argBytes),
			},
		},
		{
			name: "pack empty bytes array and string",
			args: []memo.CodecArg{
				memo.ArgPayload([]byte{}),
				memo.ArgRevertAddress(""),
			},
		},
		{
			name: "unable to parse unsupported ABI type",
			args: []memo.CodecArg{
				memo.ArgReceiver(&argAddress),
				memo.NewArg("payload", memo.ArgType("unknown"), nil),
			},
			errMsg: "failed to parse ABI string",
		},
		{
			name: "packing should fail on argument type mismatch",
			args: []memo.CodecArg{
				memo.ArgReceiver(argBytes), // expect address type, but passed bytes
			},
			errMsg: "failed to pack ABI arguments",
		},
	}

	// loop through each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// create a new ABI codec and add arguments
			codec := memo.NewCodecABI()
			codec.AddArguments(tc.args...)

			// pack arguments into ABI-encoded packedData
			packedData, err := codec.PackArguments()
			if tc.errMsg != "" {
				require.ErrorContains(t, err, tc.errMsg)
				require.Nil(t, packedData)
				return
			}
			require.NoError(t, err)

			// calc expected data for comparison
			expectedData := sample.ABIPack(t, tc.args...)

			// validate the packed data
			require.True(t, bytes.Equal(expectedData, packedData), "ABI encoded data mismatch")
		})
	}
}

func Test_CodecABI_UnpackArguments(t *testing.T) {
	// create sample arguments
	argAddress := sample.EthAddress()
	argBytes := []byte("some test bytes argument")
	argString := "some test string argument"

	// test cases
	tests := []struct {
		name     string
		data     []byte
		expected []memo.CodecArg
		errMsg   string
	}{
		{
			name: "unpack in the order of [address, bytes, string]",
			data: sample.ABIPack(t,
				memo.ArgReceiver(argAddress),
				memo.ArgPayload(argBytes),
				memo.ArgRevertAddress(argString)),
			expected: []memo.CodecArg{
				memo.ArgReceiver(argAddress),
				memo.ArgPayload(argBytes),
				memo.ArgRevertAddress(argString),
			},
		},
		{
			name: "unpack in the order of [string, address, bytes]",
			data: sample.ABIPack(t,
				memo.ArgRevertAddress(argString),
				memo.ArgReceiver(argAddress),
				memo.ArgPayload(argBytes)),
			expected: []memo.CodecArg{
				memo.ArgRevertAddress(argString),
				memo.ArgReceiver(argAddress),
				memo.ArgPayload(argBytes),
			},
		},
		{
			name: "unpack empty bytes array and string",
			data: sample.ABIPack(t,
				memo.ArgPayload([]byte{}),
				memo.ArgRevertAddress("")),
			expected: []memo.CodecArg{
				memo.ArgPayload([]byte{}),
				memo.ArgRevertAddress(""),
			},
		},
		{
			name: "unable to parse unsupported ABI type",
			data: []byte{},
			expected: []memo.CodecArg{
				memo.ArgReceiver(argAddress),
				memo.NewArg("payload", memo.ArgType("unknown"), nil),
			},
			errMsg: "failed to parse ABI string",
		},
		{
			name: "unpacking should fail on argument type mismatch",
			data: sample.ABIPack(t,
				memo.ArgReceiver(argAddress),
			),
			expected: []memo.CodecArg{
				memo.ArgReceiver(argBytes), // expect address type, but passed bytes
			},
			errMsg: "failed to unpack ABI encoded data",
		},
	}

	// loop through each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// create a new ABI codec
			codec := memo.NewCodecABI()

			// add output arguments
			output := make([]memo.CodecArg, len(tc.expected))
			for i, arg := range tc.expected {
				output[i] = memo.NewArg(arg.Name, arg.Type, newArgInstance(arg.Arg))
			}
			codec.AddArguments(output...)

			// unpack arguments from ABI-encoded data
			err := codec.UnpackArguments(tc.data)

			// validate the error message
			if tc.errMsg != "" {
				require.ErrorContains(t, err, tc.errMsg)
				return
			}

			// validate the unpacked arguments values
			require.NoError(t, err)
			for i, arg := range tc.expected {
				ensureArgEquality(t, arg.Arg, output[i].Arg)
			}
		})
	}
}
