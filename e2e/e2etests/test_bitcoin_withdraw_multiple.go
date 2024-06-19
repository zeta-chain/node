package e2etests

import (
	"math/big"
	"strconv"

	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

const defaultReceiver = "mxpYha3UJKUgSwsAz2qYRqaDSwAkKZ3YEY"

func WithdrawBitcoinMultipleTimes(r *runner.E2ERunner, args []string) {
	// ARRANGE
	// Given amount and repeat count
	require.Len(r, args, 2)
	var (
		amount = btcAmountFromFloat64(r, parseFloat(r, args[0]))
		times  = parseInt(r, args[1])
	)

	// Given BTC address set
	r.SetBtcAddress(r.Name, false)

	// Given a receiver
	receiver, err := chains.DecodeBtcAddress(defaultReceiver, r.GetBitcoinChainID())
	require.NoError(r, err)

	// ACT
	for i := 0; i < times; i++ {
		withdrawBTCZRC20(r, receiver, amount)
	}
}

func parseInt(t require.TestingT, s string) int {
	v, err := strconv.Atoi(s)
	require.NoError(t, err, "unable to parse int from %q", s)

	return v
}

// bigIntFromFloat64 takes float64 (e.g. 0.001) that represents btc amount
// and converts it to big.Int for downstream usage.
func btcAmountFromFloat64(t require.TestingT, amount float64) *big.Int {
	satoshi, err := btcutil.NewAmount(amount)
	require.NoError(t, err)

	return big.NewInt(int64(satoshi))
}
