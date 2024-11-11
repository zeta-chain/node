package e2etests

import (
	"github.com/stretchr/testify/require"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	zetabitcoin "github.com/zeta-chain/node/zetaclient/chains/bitcoin"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

// TestBitcoinDepositAndCallRevertWithDust sends a Bitcoin deposit that reverts with a dust amount in the revert outbound.
// Given the dust is too smart, the CCTX should revert
func TestBitcoinDepositAndCallRevertWithDust(r *runner.E2ERunner, args []string) {
	// ARRANGE
	// Given BTC address
	r.SetBtcAddress(r.Name, false)

	require.Len(r, args, 0)

	// Given "Live" BTC network
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// 0.002 BTC is consumed in a deposit and revert, the dust is set to 1000 satoshis in the protocol
	// Therefore the deposit amount should be 0.002 + 0.000001 = 0.00200100 should trigger the condition
	// As only 100 satoshis are left after the deposit

	amount := 0.00200100
	amount += zetabitcoin.DefaultDepositorFee

	// Given a list of UTXOs
	utxos, err := r.ListDeployerUTXOs()
	require.NoError(r, err)
	require.NotEmpty(r, utxos)

	// ACT
	// Send BTC to TSS address with a dummy memo
	// zetacore should revert cctx if call is made on a non-existing address
	nonExistReceiver := sample.EthAddress()
	badMemo := append(nonExistReceiver.Bytes(), []byte("gibberish-memo")...)
	txHash, err := r.SendToTSSFromDeployerWithMemo(amount, utxos, badMemo)
	require.NoError(r, err)
	require.NotEmpty(r, txHash)
	r.Logger.Print("BITCOIN tx hash: %s", txHash.String())

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Aborted)
}
