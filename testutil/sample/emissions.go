package sample

import (
	"cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/x/emissions/types"
	"testing"
)

func WithdrawableEmissions(t *testing.T) types.WithdrawableEmissions {
	addr := AccAddress()
	r := newRandFromStringSeed(t, addr)

	return types.WithdrawableEmissions{
		Address: AccAddress(),
		Amount:  math.NewInt(r.Int63()),
	}
}
