package memo_test

import (
	"bytes"
	"testing"

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
		encodingFormat uint8
		fields         memo.FieldsV0
		expectedFlags  byte
		expectedData   []byte
		errMsg         string
	}{
		{
			name:           "pack all fields with ABI encoding",
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
			name:           "unable to get codec on invalid encoding format",
			encodingFormat: 0x0F,
			errMsg:         "unable to get codec",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// pack the fields
			data, err := tc.fields.Pack(tc.encodingFormat)

			// validate the error message
			if tc.errMsg != "" {
				require.ErrorContains(t, err, tc.errMsg)
				return
			}

			// compare the fields
			require.NoError(t, err)
			require.Equal(t, tc.expectedFlags, data[0])
			require.True(t, bytes.Equal(tc.expectedData, data[1:]))
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
		encodingFormat uint8
		flags          byte
		data           []byte
		expected       memo.FieldsV0
		errMsg         string
	}{
		{
			name:           "unpack all fields with ABI encoding",
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
			name:           "failed to unpack ABI encoded data with compact encoding format",
			encodingFormat: memo.EncodingFmtCompactShort,
			flags:          0b00000001,
			data: sample.ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes)),
			errMsg: "failed to unpack arguments",
		},
		{
			name:           "failed to unpack data if reserved flag is not zero",
			encodingFormat: memo.EncodingFmtABI,
			flags:          0b00100001, // payload flag and reserved bit5 are set
			data: sample.ABIPack(t,
				memo.ArgReceiver(fAddress),
				memo.ArgPayload(fBytes)),
			errMsg: "reserved flag bits are not zero",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// attach data flags
			tc.data = append([]byte{tc.flags}, tc.data...)

			// unpack the fields
			fields := memo.FieldsV0{}
			err := fields.Unpack(tc.data, tc.encodingFormat)

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
