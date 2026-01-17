package coin_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/coin"
)

func Test_AzetaPerZeta(t *testing.T) {
	require.Equal(t, sdkmath.LegacyNewDec(1e18), coin.AzetaPerZeta())
}

func Test_GetAzetaDecFromAmountInZeta(t *testing.T) {
	tt := []struct {
		name        string
		zetaAmount  string
		err         require.ErrorAssertionFunc
		azetaAmount sdkmath.LegacyDec
	}{
		{
			name:        "valid zeta amount",
			zetaAmount:  "210000000",
			err:         require.NoError,
			azetaAmount: sdkmath.LegacyMustNewDecFromStr("210000000000000000000000000"),
		},
		{
			name:        "very high zeta amount",
			zetaAmount:  "21000000000000000000",
			err:         require.NoError,
			azetaAmount: sdkmath.LegacyMustNewDecFromStr("21000000000000000000000000000000000000"),
		},
		{
			name:        "very low zeta amount",
			zetaAmount:  "1",
			err:         require.NoError,
			azetaAmount: sdkmath.LegacyMustNewDecFromStr("1000000000000000000"),
		},
		{
			name:        "zero zeta amount",
			zetaAmount:  "0",
			err:         require.NoError,
			azetaAmount: sdkmath.LegacyMustNewDecFromStr("0"),
		},
		{
			name:        "decimal zeta amount",
			zetaAmount:  "0.1",
			err:         require.NoError,
			azetaAmount: sdkmath.LegacyMustNewDecFromStr("100000000000000000"),
		},
		{
			name:        "invalid zeta amount",
			zetaAmount:  "%%%%%$#",
			err:         require.Error,
			azetaAmount: sdkmath.LegacyMustNewDecFromStr("0"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			azeta, err := coin.GetAzetaDecFromAmountInZeta(tc.zetaAmount)
			tc.err(t, err)
			if err == nil {
				require.Equal(t, tc.azetaAmount, azeta)
			}
		})
	}

}

func TestGetCoinType(t *testing.T) {
	tests := []struct {
		name    string
		coin    string
		want    coin.CoinType
		wantErr bool
	}{
		{
			name:    "valid coin type 0",
			coin:    "0",
			want:    coin.CoinType(0),
			wantErr: false,
		},
		{
			name:    "valid coin type 1",
			coin:    "1",
			want:    coin.CoinType(1),
			wantErr: false,
		},
		{
			name:    "valid coin type 2",
			coin:    "2",
			want:    coin.CoinType(2),
			wantErr: false,
		},
		{
			name:    "valid coin type 3",
			coin:    "3",
			want:    coin.CoinType(3),
			wantErr: false,
		},
		{
			name:    "invalid coin type negative",
			coin:    "-1",
			wantErr: true,
		},
		{
			name: "invalid coin type large number",
			coin: "4",
			want: coin.CoinType(4),
		},
		{
			name:    "invalid coin type non-integer",
			coin:    "abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := coin.GetCoinType(tt.coin)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCoinType_SupportsRefund(t *testing.T) {
	tests := []struct {
		name string
		c    coin.CoinType
		want bool
	}{
		{"should support refund for ERC20", coin.CoinType_ERC20, true},
		{"should support refund forGas", coin.CoinType_Gas, true},
		{"should support refund forZeta", coin.CoinType_Zeta, true},
		{"should not support refund forCmd", coin.CoinType_Cmd, false},
		{"should not support refund forUnknown", coin.CoinType(100), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.SupportsRefund(); got != tt.want {
				t.Errorf("FungibleTokenCoinType.SupportsRefund() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinType_IsAsset(t *testing.T) {
	tests := []struct {
		name string
		c    coin.CoinType
		want bool
	}{
		{"Gas is asset", coin.CoinType_Gas, true},
		{"ERC20 is asset", coin.CoinType_ERC20, true},
		{"Zeta is asset", coin.CoinType_Zeta, true},
		{"Cmd is not asset", coin.CoinType_Cmd, false},
		{"CoinType_NoAssetCall is irrelevant and not asset", coin.CoinType_NoAssetCall, false},
		{"Unknown coin type is not asset", coin.CoinType(100), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsAsset(); got != tt.want {
				t.Errorf("CoinType.IsAsset() = %v, want %v", got, tt.want)
			}
		})
	}
}
