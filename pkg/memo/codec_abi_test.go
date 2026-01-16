package memo_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/memo"
)

const (
	// abiAlignment is the number of bytes used to align the ABI encoded data
	abiAlignment = 32
)

// ABIPack is a helper function that simulates the abi.Pack function.
// Note: all arguments are assumed to be <= 32 bytes for simplicity.
func ABIPack(t *testing.T, args ...memo.CodecArg) []byte {
	packedData := make([]byte, 0)

	// data offset for 1st dynamic-length field
	offset := abiAlignment * len(args)

	// 1. pack 32-byte offset for each dynamic-length field (bytes, string)
	// 2. pack actual data for each fixed-length field (address)
	for _, arg := range args {
		switch arg.Type {
		case memo.ArgTypeBytes:
			// left-pad length as uint16
			buff := make([]byte, 2)
			binary.BigEndian.PutUint16(buff, uint16(offset))
			offsetData := abiPad32(t, buff, true)
			packedData = append(packedData, offsetData...)

			argLen := len(arg.Arg.([]byte))
			if argLen > 0 {
				offset += abiAlignment * 2 // [length + data]
			} else {
				offset += abiAlignment // only [length]
			}

		case memo.ArgTypeString:
			// left-pad length as uint16
			buff := make([]byte, 2)
			binary.BigEndian.PutUint16(buff, uint16(offset))
			offsetData := abiPad32(t, buff, true)
			packedData = append(packedData, offsetData...)

			argLen := len([]byte(arg.Arg.(string)))
			if argLen > 0 {
				offset += abiAlignment * 2 // [length + data]
			} else {
				offset += abiAlignment // only [length]
			}

		case memo.ArgTypeAddress: // left-pad for address
			data := abiPad32(t, arg.Arg.(common.Address).Bytes(), true)
			packedData = append(packedData, data...)
		}
	}

	// pack dynamic-length fields
	dynamicData := abiPackDynamicData(t, args...)
	packedData = append(packedData, dynamicData...)

	return packedData
}

// abiPad32 is a helper function to pad a byte slice to 32 bytes
func abiPad32(t *testing.T, data []byte, left bool) []byte {
	// nothing needs to be encoded, return empty bytes
	if len(data) == 0 {
		return []byte{}
	}

	require.LessOrEqual(t, len(data), abiAlignment)
	padded := make([]byte, 32)

	if left {
		// left-pad the data for fixed-size types
		copy(padded[32-len(data):], data)
	} else {
		// right-pad the data for dynamic types
		copy(padded, data)
	}
	return padded
}

// abiPackDynamicData is a helper function to pack dynamic-length data
func abiPackDynamicData(t *testing.T, args ...memo.CodecArg) []byte {
	packedData := make([]byte, 0)

	// pack with ABI format: length + data
	for _, arg := range args {
		// get length
		var length int
		switch arg.Type {
		case memo.ArgTypeBytes:
			length = len(arg.Arg.([]byte))
		case memo.ArgTypeString:
			length = len([]byte(arg.Arg.(string)))
		default:
			continue
		}

		// append length in bytes
		lengthData := abiPad32(t, []byte{byte(length)}, true)
		packedData = append(packedData, lengthData...)

		// append actual data in bytes
		switch arg.Type {
		case memo.ArgTypeBytes: // right-pad for bytes
			data := abiPad32(t, arg.Arg.([]byte), false)
			packedData = append(packedData, data...)
		case memo.ArgTypeString: // right-pad for string
			data := abiPad32(t, []byte(arg.Arg.(string)), false)
			packedData = append(packedData, data...)
		}
	}

	return packedData
}

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
		require.FailNow(t, "unexpected argument type", "Type: %T", v)
	}
}

func Test_NewCodecABI(t *testing.T) {
	c := memo.NewCodecABI()
	require.NotNil(t, c)
}

func Test_CodecABI_AddArguments(t *testing.T) {
	codec := memo.NewCodecABI()
	require.NotNil(t, codec)

	address := common.HexToAddress("0xEf221eC80f004E6A2ee4E5F5d800699c1C68cD6F")
	codec.AddArguments(memo.ArgReceiver(&address))

	// attempt to pack the arguments, result should not be nil
	packedData, err := codec.PackArguments()
	require.NoError(t, err)
	require.True(t, len(packedData) > 0)
}

func Test_CodecABI_PackArgument(t *testing.T) {
	// create sample arguments
	argAddress := common.HexToAddress("0xEf221eC80f004E6A2ee4E5F5d800699c1C68cD6F")
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
			expectedData := ABIPack(t, tc.args...)

			// validate the packed data
			require.True(t, bytes.Equal(expectedData, packedData), "ABI encoded data mismatch")
		})
	}
}

func Test_CodecABI_UnpackArguments(t *testing.T) {
	// create sample arguments
	argAddress := common.HexToAddress("0xEf221eC80f004E6A2ee4E5F5d800699c1C68cD6F")
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
			data: ABIPack(t,
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
			data: ABIPack(t,
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
			data: ABIPack(t,
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
			data: ABIPack(t,
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
