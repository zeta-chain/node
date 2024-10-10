package memo_test

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func Test_V0_Pack(t *testing.T) {
	// create sample fields
	fAddress := sample.EthAddress()
	fBytes := []byte("here_s_some_bytes_field")
	fString := "this_is_a_string_field"

	tests := []struct {
		name           string
		opCode         uint8
		encodingFormat uint8
		fields         memo.FieldsV0
		expectedFlags  byte
		expectedData   []byte
		errMsg         string
	}{
		{
			name:           "pack all fields with ABI encoding",
			opCode:         memo.OpCodeDepositAndCall,
			encodingFormat: memo.EncodingFmtABI,
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
			expectedFlags: 0b00001111, // all fields are set
			expectedData: sample.ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
		},
		{
			name:           "pack all fields with compact encoding",
			opCode:         memo.OpCodeDepositAndCall,
			encodingFormat: memo.EncodingFmtCompactShort,
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
			expectedFlags: 0b00001111, // all fields are set
			expectedData: sample.CompactPack(
				memo.EncodingFmtCompactShort,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes),
				memo.ArgRevertAddress(fString),
				memo.ArgAbortAddress(fAddress),
				memo.ArgRevertMessage(fBytes)),
		},
		{
			name:           "should not pack invalid abort address",
			opCode:         memo.OpCodeDepositAndCall,
			encodingFormat: memo.EncodingFmtABI,
			fields: memo.FieldsV0{
				Receiver: fAddress,
				RevertOptions: crosschaintypes.RevertOptions{
					AbortAddress: "invalid_address",
				},
			},
			expectedFlags: 0b00000000, // no flag is set
			expectedData:  sample.ABIPack(t, memo.ArgReceiver(fAddress)),
		},
		{
			name:           "fields validation failed due to empty receiver address",
			opCode:         memo.OpCodeDepositAndCall,
			encodingFormat: memo.EncodingFmtABI,
			fields: memo.FieldsV0{
				Receiver: common.Address{},
			},
			errMsg: "receiver address is empty",
		},
		{
			name:   "unable to get codec on invalid encoding format",
			opCode: memo.OpCodeDepositAndCall,
			fields: memo.FieldsV0{
				Receiver: fAddress,
			},
			encodingFormat: 0x0F,
			errMsg:         "unable to get codec",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// pack the fields
			flags, data, err := tc.fields.Pack(tc.opCode, tc.encodingFormat)

			// validate the error message
			if tc.errMsg != "" {
				require.ErrorContains(t, err, tc.errMsg)
				require.Zero(t, flags)
				require.Nil(t, data)
				return
			}

			// compare the fields
			require.NoError(t, err)
			require.Equal(t, tc.expectedFlags, flags)
			require.True(t, bytes.Equal(tc.expectedData, data))
		})
	}
}

func Test_V0_Unpack(t *testing.T) {
	// create sample fields
	fAddress := sample.EthAddress()
	fBytes := []byte("here_s_some_bytes_field")
	fString := "this_is_a_string_field"

	tests := []struct {
		name           string
		opCode         uint8
		encodingFormat uint8
		flags          byte
		data           []byte
		expected       memo.FieldsV0
		errMsg         string
	}{
		{
			name:           "unpack all fields with ABI encoding",
			opCode:         memo.OpCodeDepositAndCall,
			encodingFormat: memo.EncodingFmtABI,
			flags:          0b00001111, // all fields are set
			data: sample.ABIPack(t,
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
			name:           "unpack all fields with compact encoding",
			opCode:         memo.OpCodeDepositAndCall,
			encodingFormat: memo.EncodingFmtCompactShort,
			flags:          0b00001111, // all fields are set
			data: sample.CompactPack(
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
					CallOnRevert:  true,
					AbortAddress:  fAddress.String(),
					RevertMessage: fBytes,
				},
			},
		},
		{
			name:           "unpack empty ABI encoded payload if flag is set",
			opCode:         memo.OpCodeDepositAndCall,
			encodingFormat: memo.EncodingFmtABI,
			flags:          0b00000001, // payload flag is set
			data: sample.ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload([]byte{})), // empty payload
			expected: memo.FieldsV0{
				Receiver: fAddress,
				Payload:  []byte{},
			},
		},
		{
			name:           "unpack empty compact encoded payload if flag is not set",
			opCode:         memo.OpCodeDepositAndCall,
			encodingFormat: memo.EncodingFmtCompactShort,
			flags:          0b00000001, // payload flag is set
			data: sample.CompactPack(
				memo.EncodingFmtCompactShort,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload([]byte{})), // empty payload
			expected: memo.FieldsV0{
				Receiver: fAddress,
				Payload:  []byte{},
			},
		},
		{
			name:           "unable to get codec on invalid encoding format",
			opCode:         memo.OpCodeDepositAndCall,
			encodingFormat: 0x0F,
			flags:          0b00000001,
			data:           []byte{},
			errMsg:         "unable to get codec",
		},
		{
			name:           "failed to unpack ABI encoded data with compact encoding format",
			opCode:         memo.OpCodeDepositAndCall,
			encodingFormat: memo.EncodingFmtCompactShort,
			flags:          0b00000001,
			data: sample.ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes)),
			errMsg: "failed to unpack arguments",
		},
		{
			name:           "fields validation failed due to empty receiver address",
			opCode:         memo.OpCodeDepositAndCall,
			encodingFormat: memo.EncodingFmtABI,
			flags:          0b00000001,
			data: sample.ABIPack(t,
				memo.ArgReceiver(common.Address{}),
				memo.ArgPayload(fBytes)),
			errMsg: "receiver address is empty",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// unpack the fields
			fields := memo.FieldsV0{}
			err := fields.Unpack(tc.opCode, tc.encodingFormat, tc.flags, tc.data)

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
	fAddress := sample.EthAddress()
	fBytes := []byte("here_s_some_bytes_field")
	fString := "this_is_a_string_field"

	tests := []struct {
		name   string
		opCode uint8
		fields memo.FieldsV0
		errMsg string
	}{
		{
			name:   "valid fields",
			opCode: memo.OpCodeDepositAndCall,
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
			name:   "invalid receiver address",
			opCode: memo.OpCodeCall,
			fields: memo.FieldsV0{
				Receiver: common.Address{}, // empty receiver address
			},
			errMsg: "receiver address is empty",
		},
		{
			name:   "payload is not allowed when opCode is deposit",
			opCode: memo.OpCodeDeposit,
			fields: memo.FieldsV0{
				Receiver: fAddress,
				Payload:  fBytes, // payload is mistakenly set
			},
			errMsg: "payload is not allowed for deposit operation",
		},
		{
			name:   "revert message is not allowed when CallOnRevert is false",
			opCode: memo.OpCodeDeposit,
			fields: memo.FieldsV0{
				Receiver: fAddress,
				RevertOptions: crosschaintypes.RevertOptions{
					RevertAddress: fString,
					CallOnRevert:  false,                    // CallOnRevert is false
					RevertMessage: []byte("revert message"), // revert message is mistakenly set
				},
			},
			errMsg: "revert message is not allowed when CallOnRevert is false",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// validate the fields
			err := tc.fields.Validate(tc.opCode)

			// validate the error message
			if tc.errMsg != "" {
				require.ErrorContains(t, err, tc.errMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}
