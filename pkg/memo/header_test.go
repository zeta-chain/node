package memo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_Header_EncodeToBytes(t *testing.T) {
	tests := []struct {
		name     string
		header   memo.Header
		expected []byte
		errMsg   string
	}{
		{
			name: "it works",
			header: memo.Header{
				Version:        0,
				EncodingFormat: memo.EncodingFmtABI,
				OpCode:         memo.OpCodeCall,
				DataFlags:      0b00001111,
			},
			expected: []byte{memo.Identifier, 0b00000000, 0b00100000, 0b00001111},
		},
		{
			name: "header validation failed",
			header: memo.Header{
				Version: 1, // invalid version
			},
			errMsg: "invalid memo version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header, err := tt.header.EncodeToBytes()
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				require.Nil(t, header)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, header)
		})
	}
}

func Test_Header_DecodeFromBytes(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected memo.Header
		errMsg   string
	}{
		{
			name: "it works",
			data: append(sample.MemoHead(0, memo.EncodingFmtABI, memo.OpCodeCall, 0, 0), []byte{0x01, 0x02}...),
			expected: memo.Header{
				Version:        0,
				EncodingFormat: memo.EncodingFmtABI,
				OpCode:         memo.OpCodeCall,
				Reserved:       0,
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
		{
			name: "header validation failed",
			data: append(
				sample.MemoHead(0, memo.EncodingFmtMax, memo.OpCodeCall, 0, 0),
				[]byte{0x01, 0x02}...), // invalid encoding format
			errMsg: "invalid encoding format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memo := &memo.Header{}
			err := memo.DecodeFromBytes(tt.data)
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, *memo)
		})
	}
}

func Test_Header_Validate(t *testing.T) {
	tests := []struct {
		name   string
		header memo.Header
		errMsg string
	}{
		{
			name: "valid header",
			header: memo.Header{
				Version:        0,
				EncodingFormat: memo.EncodingFmtCompactShort,
				OpCode:         memo.OpCodeDepositAndCall,
			},
		},
		{
			name: "invalid version",
			header: memo.Header{
				Version: 1,
			},
			errMsg: "invalid memo version",
		},
		{
			name: "invalid encoding format",
			header: memo.Header{
				EncodingFormat: memo.EncodingFmtMax,
			},
			errMsg: "invalid encoding format",
		},
		{
			name: "invalid operation code",
			header: memo.Header{
				OpCode: memo.OpCodeMax,
			},
			errMsg: "invalid operation code",
		},
		{
			name: "reserved field is not zero",
			header: memo.Header{
				Reserved: 1,
			},
			errMsg: "reserved control bits are not zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.header.Validate()
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}