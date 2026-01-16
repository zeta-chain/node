package memo_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/memo"
)

// MakeHead is a helper function to create a memo head
// Note: all arguments are assumed to be <= 0b1111 for simplicity.
func MakeHead(version, encodingFmt, opCode, reserved, flags uint8) []byte {
	head := make([]byte, memo.HeaderSize)
	head[0] = memo.Identifier
	head[1] = version<<4 | encodingFmt
	head[2] = opCode<<4 | reserved
	head[3] = flags
	return head
}

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
				Version:     0,
				EncodingFmt: memo.EncodingFmtABI,
				OpCode:      memo.OpCodeCall,
				DataFlags:   0b00011111,
			},
			expected: []byte{memo.Identifier, 0b00000000, 0b00100000, 0b00011111},
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
			data: append(
				MakeHead(0, uint8(memo.EncodingFmtABI), uint8(memo.OpCodeCall), 0, 0),
				[]byte{0x01, 0x02}...),
			expected: memo.Header{
				Version:     0,
				EncodingFmt: memo.EncodingFmtABI,
				OpCode:      memo.OpCodeCall,
				Reserved:    0,
			},
		},
		{
			name:   "memo is too short",
			data:   []byte{0x01, 0x02, 0x03},
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
				MakeHead(0, uint8(memo.EncodingFmtInvalid), uint8(memo.OpCodeCall), 0, 0),
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
				Version:     0,
				EncodingFmt: memo.EncodingFmtCompactShort,
				OpCode:      memo.OpCodeDepositAndCall,
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
				EncodingFmt: memo.EncodingFmtInvalid,
			},
			errMsg: "invalid encoding format",
		},
		{
			name: "invalid operation code",
			header: memo.Header{
				OpCode: memo.OpCodeInvalid,
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
