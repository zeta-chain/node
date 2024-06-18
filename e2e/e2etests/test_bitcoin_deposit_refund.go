package e2etests

import (
	"strconv"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

func TestDepositBTCRefund(r *runner.E2ERunner, args []string) {
	// ARRANGE
	// Given amount to send
	require.Len(r, args, 1)
	amount := parseFloat(r, args[0])

	// Given BTC address
	r.SetBtcAddress(r.Name, false)

	// Given a list of UTXOs
	utxos, err := r.BtcRPCClient.ListUnspent()
	require.NoError(r, err)
	require.NotEmpty(r, utxos)

	// ACT
	// Send a single UTXO to TSS address
	txHash, err := r.SendToTSSFromDeployerWithMemo(amount, utxos, []byte("gibberish-memo"))
	require.NotEmpty(r, err)

	// Wait for processing in zetaclient
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)

	// ASSERT
	// todo...
}

func parseFloat(t require.TestingT, s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	require.NoError(t, err, "unable to parse float %q", s)
	return f
}
