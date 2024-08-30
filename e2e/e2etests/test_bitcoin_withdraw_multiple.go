package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/pkg/chains"
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
