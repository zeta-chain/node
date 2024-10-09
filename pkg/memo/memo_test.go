package memo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func Test_EncodeToBytes(t *testing.T) {
	// create sample fields
	fAddress := sample.EthAddress()
	fBytes := []byte("here_s_some_bytes_field")
	fString := "this_is_a_string_field"

	tests := []struct {
		name         string
		memo         *memo.InboundMemo
		expectedHead []byte
		expectedData []byte
		errMsg       string
	}{
		{
			name: "encode memo with ABI encoding",
			memo: &memo.InboundMemo{
				Version:        0,
				EncodingFormat: memo.EncodingFmtABI,
				OpCode:         memo.OpCodeDepositAndCall,
				FieldsV0: memo.FieldsV0{
					Receiver: fAddress,
					Payload:  fBytes,
					RevertOptions: crosschaintypes.RevertOptions{
						RevertAddress: fString,
						CallOnRevert:  true,
						AbortAddress:  fAddress.String(), // it's a ZEVM address
						RevertMessage: fBytes,
					},
				},
			},
			expectedHead: sample.MemoHead(
				0,
				memo.EncodingFmtABI,
				memo.OpCodeDepositAndCall,
				0,
				0b00001111,
			), // all fields are set
			expectedData: sample.ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
		},
		{
			name: "encode memo with compact encoding",
			memo: &memo.InboundMemo{
				Version:        0,
				EncodingFormat: memo.EncodingFmtCompactShort,
				OpCode:         memo.OpCodeDepositAndCall,
				FieldsV0: memo.FieldsV0{
					Receiver: fAddress,
					Payload:  fBytes,
					RevertOptions: crosschaintypes.RevertOptions{
						RevertAddress: fString,
						CallOnRevert:  true,
						AbortAddress:  fAddress.String(), // it's a ZEVM address
						RevertMessage: fBytes,
					},
				},
			},
			expectedHead: sample.MemoHead(
				0,
				memo.EncodingFmtCompactShort,
				memo.OpCodeDepositAndCall,
				0,
				0b00001111,
			), // all fields are set
			expectedData: sample.CompactPack(
				memo.EncodingFmtCompactShort,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
		},
		{
			name: "failed to encode if basic validation fails",
			memo: &memo.InboundMemo{
				EncodingFormat: memo.EncodingFmtMax,
			},
			errMsg: "invalid encoding format",
		},
		{
			name: "failed to encode if version is invalid",
			memo: &memo.InboundMemo{
				Version: 1,
			},
			errMsg: "invalid memo version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.memo.EncodeToBytes()
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}
			require.NoError(t, err)
			require.Equal(t, append(tt.expectedHead, tt.expectedData...), data)
		})
	}
}

func Test_DecodeFromBytes(t *testing.T) {
	// create sample fields
	fAddress := sample.EthAddress()
	fBytes := []byte("here_s_some_bytes_field")
	fString := "this_is_a_string_field"

	tests := []struct {
		name         string
		head         []byte
		data         []byte
		expectedMemo memo.InboundMemo
		errMsg       string
	}{
		{
			name: "decode memo with ABI encoding",
			head: sample.MemoHead(
				0,
				memo.EncodingFmtABI,
				memo.OpCodeDepositAndCall,
				0,
				0b00001111,
			), // all fields are set
			data: sample.ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
			expectedMemo: memo.InboundMemo{
				Version:        0,
				EncodingFormat: memo.EncodingFmtABI,
				OpCode:         memo.OpCodeDepositAndCall,
				FieldsV0: memo.FieldsV0{
					Receiver: fAddress,
					Payload:  fBytes,
					RevertOptions: crosschaintypes.RevertOptions{
						RevertAddress: fString,
						CallOnRevert:  true,
						AbortAddress:  fAddress.String(), // it's a ZEVM address
						RevertMessage: fBytes,
					},
				},
			},
		},
		{
			name: "decode memo with compact encoding",
			head: sample.MemoHead(
				0,
				memo.EncodingFmtCompactLong,
				memo.OpCodeDepositAndCall,
				0,
				0b00001111,
			), // all fields are set
			data: sample.CompactPack(
				memo.EncodingFmtCompactLong,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
			expectedMemo: memo.InboundMemo{
				Version:        0,
				EncodingFormat: memo.EncodingFmtCompactLong,
				OpCode:         memo.OpCodeDepositAndCall,
				FieldsV0: memo.FieldsV0{
					Receiver: fAddress,
					Payload:  fBytes,
					RevertOptions: crosschaintypes.RevertOptions{
						RevertAddress: fString,
						CallOnRevert:  true,
						AbortAddress:  fAddress.String(), // it's a ZEVM address
						RevertMessage: fBytes,
					},
				},
			},
		},
		{
			name:   "failed to decode if basic validation fails",
			head:   sample.MemoHead(0, memo.EncodingFmtABI, memo.OpCodeMax, 0, 0),
			data:   sample.ABIPack(t, memo.ArgReceiver(fAddress)),
			errMsg: "invalid operation code",
		},
		{
			name:   "failed to decode if version is invalid",
			head:   sample.MemoHead(1, memo.EncodingFmtABI, memo.OpCodeDeposit, 0, 0),
			data:   sample.ABIPack(t, memo.ArgReceiver(fAddress)),
			errMsg: "invalid memo version",
		},
		{

			name: "failed to decode compact encoded data with ABI encoding format",
			head: sample.MemoHead(
				0,
				memo.EncodingFmtABI,
				memo.OpCodeDepositAndCall,
				0,
				0,
			), // head says ABI encoding
			data: sample.CompactPack(
				memo.EncodingFmtCompactShort,
				memo.ArgReceiver(fAddress),
			), // but data is compact encoded
			errMsg: "failed to unpack memo fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := append(tt.head, tt.data...)
			memo, err := memo.DecodeFromBytes(data)
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedMemo, *memo)
		})
	}
}

func Test_ValidateBasics(t *testing.T) {
	tests := []struct {
		name   string
		memo   *memo.InboundMemo
		errMsg string
	}{
		{
			name: "valid memo",
			memo: &memo.InboundMemo{
				Version:        0,
				EncodingFormat: memo.EncodingFmtCompactShort,
				OpCode:         memo.OpCodeDepositAndCall,
			},
		},
		{
			name: "invalid encoding format",
			memo: &memo.InboundMemo{
				EncodingFormat: memo.EncodingFmtMax,
			},
			errMsg: "invalid encoding format",
		},
		{
			name: "invalid operation code",
			memo: &memo.InboundMemo{
				OpCode: memo.OpCodeMax,
			},
			errMsg: "invalid operation code",
		},
		{
			name: "reserved field is not zero",
			memo: &memo.InboundMemo{
				Reserved: 1,
			},
			errMsg: "reserved control bits are not zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.memo.ValidateBasics()
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func Test_EncodeBasics(t *testing.T) {
	tests := []struct {
		name     string
		memo     *memo.InboundMemo
		expected []byte
		errMsg   string
	}{
		{
			name: "it works",
			memo: &memo.InboundMemo{
				Version:        1,
				EncodingFormat: memo.EncodingFmtABI,
				OpCode:         memo.OpCodeCall,
				Reserved:       15,
			},
			expected: []byte{memo.MemoIdentifier, 0b00010000, 0b00101111},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basics := tt.memo.EncodeBasics()
			require.Equal(t, tt.expected, basics)
		})
	}
}

func Test_DecodeBasics(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected memo.InboundMemo
		errMsg   string
	}{
		{
			name: "it works",
			data: append(sample.MemoHead(1, memo.EncodingFmtABI, memo.OpCodeCall, 15, 0), []byte{0x01, 0x02}...),
			expected: memo.InboundMemo{
				Version:        1,
				EncodingFormat: memo.EncodingFmtABI,
				OpCode:         memo.OpCodeCall,
				Reserved:       15,
			},
		},
		{
			name:   "memo is too short",
			data:   []byte{0x01, 0x02, 0x03, 0x04},
			errMsg: "memo is too short",
		},
		{
			name:   "invalid memo identifier",
			data:   []byte{'M', 0x02, 0x03, 0x04, 0x05},
			errMsg: "invalid memo identifier",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memo := &memo.InboundMemo{}
			err := memo.DecodeBasics(tt.data)
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, *memo)
		})
	}
}
