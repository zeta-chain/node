package utils

import (
	"math/big"

	"github.com/stretchr/testify/require"
)

// BalanceChange contains details about the balance change
type BalanceChange struct {
	// If provided, when set to
	//  - true: the balance change must be positive
	//  - false: the balance change must be negative
	positive *bool

	// If provided, it's the exact value of the balance change
	// delta value can be positive or negative
	delta *big.Int
}

// NewBalanceChange returns a new BalanceChange with given positive flag
func NewBalanceChange(positive bool) BalanceChange {
	return BalanceChange{
		positive: &[]bool{positive}[0],
	}
}

// NewExactChange returns a new BalanceChange with the Delta field set to the exact value
func NewExactChange(delta *big.Int) BalanceChange {
	return BalanceChange{
		delta: delta,
	}
}

// Verify verifies the balance change
func (c BalanceChange) Verify(t require.TestingT, oldBalance *big.Int, newBalance *big.Int) {
	// check exact amount change if provided
	if c.delta != nil {
		require.Equal(t, new(big.Int).Add(oldBalance, c.delta), newBalance)
		return
	}

	// otherwise, check positive/negative only
	if c.positive != nil {
		if *c.positive {
			require.True(t, newBalance.Cmp(oldBalance) > 0, "balance should be increased")
		} else {
			require.True(t, newBalance.Cmp(oldBalance) < 0, "balance should be decreased")
		}
	}
}
