package memo_test

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

const (
	// flagsAllFieldsSet sets all fields: [receiver, payload, revert address, abort address, CallOnRevert, revert message]
	flagsAllFieldsSet = 0b00111111
)

func Test_V0_Pack(t *testing.T) {
	// create sample fields
	fAddress := sample.EthAddress()
	fBytes := []byte("here_s_some_bytes_field")
	fString := "this_is_a_string_field"

	tests := []struct {
		name          string
		opCode        memo.OpCode
		encodeFmt     memo.EncodingFormat
		dataFlags     uint8
		fields        memo.FieldsV0
		expectedFlags byte
		expectedData  []byte
		errMsg        string
	}{
		{
			name:      "pack all fields with ABI encoding",
			opCode:    memo.OpCodeDepositAndCall,
			encodeFmt: memo.EncodingFmtABI,
			dataFlags: flagsAllFieldsSet, // all fields are set
			fields: memo.FieldsV0{
				Receiver: fAddress,
				Payload:  fBytes,
				RevertOptions: crosschaintypes.RevertOptions{
					RevertAddress: fString,
					CallOnRevert:  true,
					AbortAddress:  fAddress.String(), // it's a ZEVM address
					RevertMessage: fBytes,
				},
			},
			expectedData: ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
		},
		{
			name:      "pack all fields with compact encoding",
			opCode:    memo.OpCodeDepositAndCall,
			encodeFmt: memo.EncodingFmtCompactShort,
			dataFlags: 0b00101111, // all fields are set except callOnRevert flag
			fields: memo.FieldsV0{
				Receiver: fAddress,
				Payload:  fBytes,
				RevertOptions: crosschaintypes.RevertOptions{
					RevertAddress: fString,
					CallOnRevert:  false,             // CallOnRevert is irrelevant to RevertMessage
					AbortAddress:  fAddress.String(), // it's a ZEVM address
					RevertMessage: fBytes,
				},
			},
			expectedData: CompactPack(
				memo.EncodingFmtCompactShort,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
		},
		{
			name:      "fields validation failed due to empty receiver address",
			opCode:    memo.OpCodeDepositAndCall,
			encodeFmt: memo.EncodingFmtABI,
			dataFlags: 0b00000001, // receiver flag is set
			fields: memo.FieldsV0{
				Receiver: common.Address{},
			},
			errMsg: "receiver address is empty",
		},
		{
			name:      "unable to get codec on invalid encoding format",
			opCode:    memo.OpCodeDepositAndCall,
			dataFlags: 0b00000001,
			fields: memo.FieldsV0{
				Receiver: fAddress,
			},
			encodeFmt: 0x0F,
			errMsg:    "unable to get codec",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// pack the fields
			data, err := tc.fields.Pack(tc.opCode, tc.encodeFmt, tc.dataFlags)

			// validate the error message
			if tc.errMsg != "" {
				require.ErrorContains(t, err, tc.errMsg)
				require.Nil(t, data)
				return
			}

			// compare the fields
			require.NoError(t, err)
			require.True(t, bytes.Equal(tc.expectedData, data))
		})
	}
}

func Test_V0_Unpack(t *testing.T) {
	// create sample fields
	fAddress := common.HexToAddress("0xA029D053E13223E2442E28be80b3CeDA27ecbE31")
	fBytes := []byte("here_s_some_bytes_field")
	fString := "this_is_a_string_field"

	tests := []struct {
		name      string
		encodeFmt memo.EncodingFormat
		dataFlags byte
		data      []byte
		expected  memo.FieldsV0
		errMsg    string
	}{
		{
			name:      "unpack all fields with ABI encoding",
			encodeFmt: memo.EncodingFmtABI,
			dataFlags: flagsAllFieldsSet, // all fields are set
			data: ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
			expected: memo.FieldsV0{
				Receiver: fAddress,
				Payload:  fBytes,
				RevertOptions: crosschaintypes.RevertOptions{
					RevertAddress: fString,
					CallOnRevert:  true,
					AbortAddress:  fAddress.String(),
					RevertMessage: fBytes,
				},
			},
		},
		{
			name:      "unpack all fields with compact encoding",
			encodeFmt: memo.EncodingFmtCompactShort,
			dataFlags: 0b00101111, // all fields are set except callOnRevert flag
			data: CompactPack(
				memo.EncodingFmtCompactShort,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
			expected: memo.FieldsV0{
				Receiver: fAddress,
				Payload:  fBytes,
				RevertOptions: crosschaintypes.RevertOptions{
					RevertAddress: fString,
					CallOnRevert:  false, // CallOnRevert is irrelevant to RevertMessage
					AbortAddress:  fAddress.String(),
					RevertMessage: fBytes,
				},
			},
		},
		{
			name:      "unpack empty ABI encoded payload if flag is set",
			encodeFmt: memo.EncodingFmtABI,
			dataFlags: 0b00000010, // payload flags are set
			data: ABIPack(t,
				memo.ArgPayload([]byte{})), // empty payload
			expected: memo.FieldsV0{},
		},
		{
			name:      "unpack empty compact encoded payload if flag is set",
			encodeFmt: memo.EncodingFmtCompactShort,
			dataFlags: 0b00000010, // payload flag is set
			data: CompactPack(
				memo.EncodingFmtCompactShort,
				memo.ArgPayload([]byte{})), // empty payload
			expected: memo.FieldsV0{},
		},
		{
			name:      "unable to get codec on invalid encoding format",
			encodeFmt: 0x0F,
			dataFlags: 0b00000001,
			data:      []byte{},
			errMsg:    "unable to get codec",
		},
		{
			name:      "failed to unpack ABI encoded data with compact encoding format",
			encodeFmt: memo.EncodingFmtCompactShort,
			dataFlags: 0b00000011, // receiver and payload flags are set
			data: ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes)),
			errMsg: "failed to unpack arguments",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// unpack the fields
			fields := memo.FieldsV0{}
			err := fields.Unpack(tc.encodeFmt, tc.dataFlags, tc.data)

			// validate the error message
			if tc.errMsg != "" {
				require.ErrorContains(t, err, tc.errMsg)
				return
			}

			// compare the fields
			require.NoError(t, err)
			require.Equal(t, tc.expected, fields)
		})
	}
}

func Test_V0_Validate(t *testing.T) {
	// create sample fields
	fAddress := common.HexToAddress("0xA029D053E13223E2442E28be80b3CeDA27ecbE31")
	fBytes := []byte("here_s_some_bytes_field")
	fString := "this_is_a_string_field"

	tests := []struct {
		name      string
		opCode    memo.OpCode
		dataFlags uint8
		fields    memo.FieldsV0
		errMsg    string
	}{
		{
			name:      "valid fields",
			opCode:    memo.OpCodeDepositAndCall,
			dataFlags: flagsAllFieldsSet, // all fields are set
			fields: memo.FieldsV0{
				Receiver: fAddress,
				Payload:  fBytes,
				RevertOptions: crosschaintypes.RevertOptions{
					RevertAddress: fString,
					CallOnRevert:  true,
					AbortAddress:  fAddress.String(),
					RevertMessage: fBytes,
				},
			},
		},
		{
			name:      "receiver address flag is not set",
			opCode:    memo.OpCodeDepositAndCall,
			dataFlags: 0b00000000, // receiver flag is not set
			fields: memo.FieldsV0{
				Receiver: fAddress,
			},
			errMsg: "must set receiver address flag",
		},
		{
			name:      "invalid receiver address",
			opCode:    memo.OpCodeCall,
			dataFlags: 0b00000001, // receiver flag is set
			fields: memo.FieldsV0{
				Receiver: common.Address{}, // provide empty receiver address
			},
			errMsg: "receiver address is empty",
		},
		{
			name:      "payload is not allowed when opCode is deposit",
			opCode:    memo.OpCodeDeposit,
			dataFlags: 0b00000001, // receiver flag is set
			fields: memo.FieldsV0{
				Receiver: fAddress,
				Payload:  fBytes, // payload is mistakenly set
			},
			errMsg: "payload is not allowed for deposit operation",
		},
		{
			name:      "revert message is empty",
			opCode:    memo.OpCodeDepositAndCall,
			dataFlags: 0b00000101, // revert message flag is set
			fields: memo.FieldsV0{
				Receiver: fAddress,
				RevertOptions: crosschaintypes.RevertOptions{
					CallOnRevert:  true,
					RevertMessage: []byte("revert message"),
				},
			},
			errMsg: "revert address is empty",
		},
		{
			name:      "abort address is invalid",
			opCode:    memo.OpCodeDeposit,
			dataFlags: 0b00001001, // abort address flag is set
			fields: memo.FieldsV0{
				Receiver: fAddress,
				RevertOptions: crosschaintypes.RevertOptions{
					AbortAddress: "invalid abort address",
				},
			},
			errMsg: "invalid abort address",
		},
		{
			name:      "abort address is empty",
			opCode:    memo.OpCodeDepositAndCall,
			dataFlags: 0b00001001, // abort address flag is set
			fields: memo.FieldsV0{
				Receiver: fAddress,
				RevertOptions: crosschaintypes.RevertOptions{
					AbortAddress: constant.EVMZeroAddress,
				},
			},
			errMsg: "abort address is empty",
		},
		{
			name:      "reserved flags are not zero",
			opCode:    memo.OpCodeDepositAndCall,
			dataFlags: 0b01000001, // reserved flags are not zero
			fields: memo.FieldsV0{
				Receiver: fAddress,
			},
			errMsg: "reserved flags are not zero",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// validate the fields
			err := tc.fields.Validate(tc.opCode, tc.dataFlags)

			// validate the error message
			if tc.errMsg != "" {
				require.ErrorContains(t, err, tc.errMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func Test_V0_DataFlags(t *testing.T) {
	// create sample fields
	fAddress := common.HexToAddress("0xA029D053E13223E2442E28be80b3CeDA27ecbE31")
	fBytes := []byte("here_s_some_bytes_field")
	fString := "this_is_a_string_field"

	tests := []struct {
		name          string
		fields        memo.FieldsV0
		expectedFlags uint8
	}{
		{
			name: "all fields set",
			fields: memo.FieldsV0{
				Receiver: fAddress,
				Payload:  fBytes,
				RevertOptions: crosschaintypes.RevertOptions{
					RevertAddress: fString,
					CallOnRevert:  true,
					AbortAddress:  fAddress.String(),
					RevertMessage: fBytes,
				},
			},
			expectedFlags: flagsAllFieldsSet,
		},
		{
			name:          "no fields set",
			fields:        memo.FieldsV0{},
			expectedFlags: 0b00000000,
		},
		{
			name: "a few fields set",
			fields: memo.FieldsV0{
				Receiver: fAddress,
				RevertOptions: crosschaintypes.RevertOptions{
					RevertAddress: fString,
					CallOnRevert:  false, // CallOnRevert is irrelevant to RevertMessage
					RevertMessage: fBytes,
				},
			},
			expectedFlags: 0b00100101,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// get the data flags
			flags := tc.fields.DataFlags()

			// compare the flags
			require.Equal(t, tc.expectedFlags, flags)
		})
	}
}
