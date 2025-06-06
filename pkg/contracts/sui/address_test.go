package sui

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
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
		name  string
		input []byte
		want  string
	}{
		{
			name: "sample bytes",
			input: func() []byte {
				b, _ := hex.DecodeString("1234567890abcdef")
				return b
			}(),
			want: "0x1234567890abcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecodeAddress(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestValidateAddress(t *testing.T) {
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
			name:    "Uppercase addresses are explicitly rejected",
			address: "0X2A4C5A97B561AC5B38EDC4B4E9B2C183C57B56DF5B1EA2F1C6F2E4A44B92D59F",
			wantErr: true,
		},
		{
			name:    "Short addresses are explicitly rejected",
			address: "0x1a",
			wantErr: true,
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
			address: "0xZZZZZZ0000000000000000000000000000000000000000000000000000000000",
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddress(tt.address)
			if tt.wantErr {
				require.Error(t, err, tt.address)
				return
			}

			require.NoError(t, err, tt.address)
		})
	}
}
