package sample

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func WithdrawableEmissions(t *testing.T) types.WithdrawableEmissions {
	addr := AccAddress()
	r := newRandFromStringSeed(t, addr)

	return types.WithdrawableEmissions{
		Address: AccAddress(),
		Amount:  math.NewInt(r.Int63()),
	}
}

func WithdrawEmission(t *testing.T) types.WithdrawEmission {
	addr := AccAddress()
	r := newRandFromStringSeed(t, addr)

	return types.WithdrawEmission{
		Address: AccAddress(),
		Amount:  math.NewInt(r.Int63()),
	}
}
