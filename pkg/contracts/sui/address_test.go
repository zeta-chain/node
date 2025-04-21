package sui

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncodeAddress(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantHex   string
		shouldErr bool
	}{
		{
			name:    "valid short address",
			input:   "0xabc",
			wantHex: "0000000000000000000000000000000000000000000000000000000000000abc",
		},
		{
			name:    "valid full 64 char address",
			input:   "0x" + "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899",
			wantHex: "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899",
		},
		{
			name:      "missing 0x prefix",
			input:     "abcdef",
			shouldErr: true,
		},
		{
			name:      "empty hex part",
			input:     "0x",
			shouldErr: true,
		},
		{
			name:      "too long address",
			input:     "0x" + "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aa",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeAddress(tt.input)
			if tt.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantHex, hex.EncodeToString(got))
			}
		})
	}
}

func TestDecodeAddress(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		want      string
		shouldErr bool
	}{
		{
			name:  "empty input",
			input: []byte{},
			want:  "0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:  "short input",
			input: []byte{0x1, 0x2, 0x3},
			want:  "0x0000000000000000000000000000000000000000000000000000000000010203",
		},
		{
			name: "exact 32 bytes",
			input: func() []byte {
				b, _ := hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
				return b
			}(),
			want: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
		{
			name: "too long (33 bytes)",
			input: func() []byte {
				b := make([]byte, 33)
				return b
			}(),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeAddress(tt.input)
			if tt.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCheckValidSuiAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{
			name:    "Valid full-length address",
			address: "0x2a4c5a97b561ac5b38edc4b4e9b2c183c57b56df5b1ea2f1c6f2e4a44b92d59f",
			wantErr: false,
		},
		{
			name:    "Valid short address",
			address: "0x1a",
			wantErr: false,
		},
		{
			name:    "Missing 0x prefix",
			address: "2a4c5a97b561ac5b38edc4b4e9b2c183c57b56df5b1ea2f1c6f2e4a44b92d59f",
			wantErr: true,
		},
		{
			name:    "Too long address",
			address: "0x" + "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899a1",
			wantErr: true,
		},
		{
			name:    "Invalid hex characters",
			address: "0xZZZZZZ",
			wantErr: true,
		},
		{
			name:    "Empty string",
			address: "",
			wantErr: true,
		},
		{
			name:    "Only 0x",
			address: "0x",
			wantErr: true,
		},
		{
			name:    "Minimal valid single-byte address",
			address: "0x0",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidAddress(tt.address)
			if tt.wantErr {
				require.Error(t, err, "expected error for address: %s", tt.address)
			} else {
				require.NoError(t, err, "unexpected error for address: %s", tt.address)
			}
		})
	}
}
