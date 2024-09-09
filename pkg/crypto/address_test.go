package crypto

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/constant"
	"testing"
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
