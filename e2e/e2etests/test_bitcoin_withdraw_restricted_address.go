package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestBitcoinWithdrawRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawalAmount := parseFloat(r, args[0])
	amount := btcAmountFromFloat64(r, withdrawalAmount)

	r.SetBtcAddress(r.Name, false)

	withdrawBitcoinRestricted(r, amount)
}

func withdrawBitcoinRestricted(r *runner.E2ERunner, amount *big.Int) {
	// use restricted BTC P2WPKH address
	addressRestricted, err := chains.DecodeBtcAddress(
		sample.RestrictedBtcAddressTest,
		chains.BitcoinRegtest.ChainId,
	)
	require.NoError(r, err)

	// the cctx should be cancelled
	rawTx := withdrawBTCZRC20(r, addressRestricted, amount)
	require.Len(r, rawTx.Vout, 2, "BTC cancelled outtx rawTx.Vout should have 2 outputs")
}
