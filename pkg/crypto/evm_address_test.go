package crypto

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/constant"
)

func TestIsEmptyAddress(t *testing.T) {
	tests := []struct {
		name    string
		address common.Address
		want    bool
	}{
		{
			name:    "empty address",
			address: common.Address{},
			want:    true,
		},
		{
			name:    "zero address",
			address: common.HexToAddress(constant.EVMZeroAddress),
			want:    true,
		},
		{
			name:    "non empty address",
			address: common.HexToAddress("0x1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.want, IsEmptyAddress(tt.address))
		})
	}
}

func TestIsEVMAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "EVM address",
			address: "0x5a4f260A7D716c859A2736151cB38b9c58C32c64",
			want:    true,
		},
		{
			name:    "EVM address with invalid checksum",
			address: "0x5a4f260a7D716c859A2736151CB38b9c58C32c64",
			want:    true,
		},
		{
			name:    "EVM address all lowercase",
			address: "0x5a4f260a7d716c859a2736151cb38b9c58c32c64",
			want:    true,
		},
		{
			name:    "EVM address all uppercase",
			address: "0x5A4F260A7D716C859A2736151CB38B9C58C32C64",
			want:    true,
		},
		{
			name:    "invalid EVM address",
			address: "5a4f260A7D716c859A2736151cB38b9c58C32c64",
		},
		{
			name:    "empty address",
			address: "",
		},
		{
			name:    "non EVM address",
			address: "Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.want, IsEVMAddress(tt.address))
		})
	}
}

func TestIsChecksumAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "checksum address",
			address: "0x5a4f260A7D716c859A2736151cB38b9c58C32c64",
			want:    true,
		},
		{
			name:    "invalid checksum address",
			address: "0x5a4f260a7D716c859A2736151CB38b9c58C32c64",
		},
		{
			name:    "all lowercase",
			address: "0x5a4f260a7d716c859a2736151cb38b9c58c32c64",
		},
		{
			name:    "all uppercase",
			address: "0x5A4F260A7D716C859A2736151CB38B9C58C32C64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.want, IsChecksumAddress(tt.address))
		})
	}
}

func TestToChecksumAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    string
	}{
		{
			name:    "checksum address",
			address: "0x5a4f260A7D716c859A2736151cB38b9c58C32c64",
			want:    "0x5a4f260A7D716c859A2736151cB38b9c58C32c64",
		},
		{
			name:    "all lowercase",
			address: "0x5a4f260a7d716c859a2736151cb38b9c58c32c64",
			want:    "0x5a4f260A7D716c859A2736151cB38b9c58C32c64",
		},
		{
			name:    "all uppercase",
			address: "0x5A4F260A7D716C859A2736151CB38B9C58C32C64",
			want:    "0x5a4f260A7D716c859A2736151cB38b9c58C32c64",
		},
		{
			name:    "empty address returns null address",
			address: "",
			want:    "0x0000000000000000000000000000000000000000",
		},
		{
			name:    "non evm address returns null address",
			address: "Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr",
			want:    "0x0000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.want, ToChecksumAddress(tt.address))
		})
	}
}
