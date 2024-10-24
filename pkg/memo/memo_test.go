package memo_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func Test_Memo_EncodeToBytes(t *testing.T) {
	// create sample fields
	fAddress := common.HexToAddress("0xEA9808f0Ac504d1F521B5BbdfC33e6f1953757a7")
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
				Header: memo.Header{
					Version:     0,
					EncodingFmt: memo.EncodingFmtABI,
					OpCode:      memo.OpCodeDepositAndCall,
				},
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
			expectedHead: MakeHead(
				0,
				uint8(memo.EncodingFmtABI),
				uint8(memo.OpCodeDepositAndCall),
				0,
				flagsAllFieldsSet, // all fields are set
			),
			expectedData: ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
		},
		{
			name: "encode memo with compact encoding",
			memo: &memo.InboundMemo{
				Header: memo.Header{
					Version:     0,
					EncodingFmt: memo.EncodingFmtCompactShort,
					OpCode:      memo.OpCodeDepositAndCall,
				},
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
			expectedHead: MakeHead(
				0,
				uint8(memo.EncodingFmtCompactShort),
				uint8(memo.OpCodeDepositAndCall),
				0,
				flagsAllFieldsSet, // all fields are set
			),
			expectedData: CompactPack(
				memo.EncodingFmtCompactShort,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
		},
		{
			name: "failed to encode memo header",
			memo: &memo.InboundMemo{
				Header: memo.Header{
					OpCode: memo.OpCodeInvalid, // invalid operation code
				},
			},
			errMsg: "failed to encode memo header",
		},
		{
			name: "failed to encode if version is invalid",
			memo: &memo.InboundMemo{
				Header: memo.Header{
					Version: 1,
				},
			},
			errMsg: "invalid memo version",
		},
		{
			name: "failed to pack memo fields",
			memo: &memo.InboundMemo{
				Header: memo.Header{
					Version:     0,
					EncodingFmt: memo.EncodingFmtABI,
					OpCode:      memo.OpCodeDeposit,
				},
				FieldsV0: memo.FieldsV0{
					Receiver: fAddress,
					Payload:  fBytes, // payload is not allowed for deposit
				},
			},
			errMsg: "failed to pack memo fields version: 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.memo.EncodeToBytes()
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				require.Nil(t, data)
				return
			}
			require.NoError(t, err)
			require.Equal(t, append(tt.expectedHead, tt.expectedData...), data)

			// decode the memo and compare with the original
			decodedMemo, isMemo, err := memo.DecodeFromBytes(data)
			require.True(t, isMemo)
			require.NoError(t, err)
			require.Equal(t, tt.memo, decodedMemo)
		})
	}
}

func Test_Memo_DecodeFromBytes(t *testing.T) {
	// create sample fields
	fAddress := common.HexToAddress("0xEA9808f0Ac504d1F521B5BbdfC33e6f1953757a7")
	fBytes := []byte("here_s_some_bytes_field")
	fString := "this_is_a_string_field"

	tests := []struct {
		name         string
		head         []byte
		data         []byte
		isMemo       bool
		expectedMemo memo.InboundMemo
		errMsg       string
	}{
		{
			name: "decode memo with ABI encoding",
			head: MakeHead(
				0,
				uint8(memo.EncodingFmtABI),
				uint8(memo.OpCodeDepositAndCall),
				0,
				flagsAllFieldsSet, // all fields are set
			),
			data: ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
			isMemo: true,
			expectedMemo: memo.InboundMemo{
				Header: memo.Header{
					Version:     0,
					EncodingFmt: memo.EncodingFmtABI,
					OpCode:      memo.OpCodeDepositAndCall,
					DataFlags:   0b00011111,
				},
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
			head: MakeHead(
				0,
				uint8(memo.EncodingFmtCompactLong),
				uint8(memo.OpCodeDepositAndCall),
				0,
				flagsAllFieldsSet, // all fields are set
			),
			data: CompactPack(
				memo.EncodingFmtCompactLong,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
			isMemo: true,
			expectedMemo: memo.InboundMemo{
				Header: memo.Header{
					Version:     0,
					EncodingFmt: memo.EncodingFmtCompactLong,
					OpCode:      memo.OpCodeDepositAndCall,
					DataFlags:   0b00011111,
				},
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
			name:   "failed to decode memo header",
			head:   MakeHead(0, uint8(memo.EncodingFmtABI), uint8(memo.OpCodeInvalid), 0, 0),
			data:   ABIPack(t, memo.ArgReceiver(fAddress)),
			errMsg: "failed to decode memo header",
		},
		{
			name:   "failed to decode if version is invalid",
			head:   MakeHead(1, uint8(memo.EncodingFmtABI), uint8(memo.OpCodeDeposit), 0, 0),
			data:   ABIPack(t, memo.ArgReceiver(fAddress)),
			errMsg: "invalid memo version",
		},
		{
			name: "failed to decode compact encoded data with ABI encoding format",
			head: MakeHead(
				0,
				uint8(memo.EncodingFmtABI),
				uint8(memo.OpCodeDepositAndCall),
				0,
				0,
			), // header says ABI encoding
			data: CompactPack(
				memo.EncodingFmtCompactShort,
				memo.ArgReceiver(fAddress),
			), // but data is compact encoded
			errMsg: "failed to unpack memo FieldsV0",
		},
		{
			name: "should return [nil,  true, err] if fields validation fails",
			head: MakeHead(
				0,
				uint8(memo.EncodingFmtABI),
				uint8(memo.OpCodeDepositAndCall),
				0,
				0b00000011, // receiver flag is set
			),
			data: ABIPack(t,
				memo.ArgReceiver(common.Address{}), // empty receiver address provided
				memo.ArgPayload(fBytes)),
			isMemo: true, // it's still a memo, but with invalid field values
			errMsg: "failed to validate memo FieldsV0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := append(tt.head, tt.data...)
			memo, isMemo, err := memo.DecodeFromBytes(data)

			// check error message
			if tt.errMsg != "" {
				require.Nil(t, memo)
				require.ErrorContains(t, err, tt.errMsg)
				return
			}

			// a standard memo or not
			require.Equal(t, tt.isMemo, isMemo)
			if !isMemo {
				require.Nil(t, memo)
				return
			}

			// if it's a standard memo, depending on validation result
			if err != nil {
				require.Nil(t, memo)
			} else {
				require.NotNil(t, memo)
				require.Equal(t, tt.expectedMemo, *memo)
			}
		})
	}
}

func Test_DecodeLegacyMemoHex(t *testing.T) {
	expectedShortMsgResult, err := hex.DecodeString("1a2b3c4d5e6f708192a3b4c5d6e7f808")
	require.NoError(t, err)
	tests := []struct {
		name       string
		message    string
		expectAddr common.Address
		expectData []byte
		wantErr    bool
	}{
		{
			"valid msg",
			"95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5",
			common.HexToAddress("95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5"),
			[]byte{},
			false,
		},
		{"empty msg", "", common.Address{}, nil, false},
		{"invalid hex", "invalidHex", common.Address{}, nil, true},
		{"short msg", "1a2b3c4d5e6f708192a3b4c5d6e7f808", common.Address{}, expectedShortMsgResult, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, data, err := memo.DecodeLegacyMemoHex(tt.message)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectAddr, addr)
				require.Equal(t, tt.expectData, data)
			}
		})
	}
}
