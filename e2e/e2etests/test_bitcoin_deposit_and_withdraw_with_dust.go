package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/withdrawer"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestBitcoinDepositAndWithdrawWithDust deposits Bitcoin and call a smart contract that withdraw dust amount
// It tests the edge case where during a cross-chain call, a invaild withdraw is initiated (processLogs fails)
func TestBitcoinDepositAndWithdrawWithDust(r *runner.E2ERunner, args []string) {
	// Given "Live" BTC network
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	require.Len(r, args, 0)

	// ARRANGE

	// Deploy the withdrawer contract on ZetaChain with a withdraw amount of 100 satoshis (dust amount is 1000 satoshis)
	withdrawerAddr, tx, _, err := withdrawer.DeployWithdrawer(r.ZEVMAuth, r.ZEVMClient, big.NewInt(100))
	require.NoError(r, err)

	// Wait for the transaction to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.Equal(r, receipt.Status, uint64(1))

	// Given a list of UTXOs
	utxos, err := r.ListDeployerUTXOs()
	require.NoError(r, err)
	require.NotEmpty(r, utxos)

	// ACT
	// Deposit 0.01 BTC to withdrawer, this is an arbitrary amount, must be greater than dust amount
	txHash, err := r.SendToTSSFromDeployerWithMemo(0.01, utxos, withdrawerAddr.Bytes())
	require.NoError(r, err)
	require.NotEmpty(r, txHash)

	// ASSERT
	// Now we want to make sure the cctx is aborted with expected error message

	// cctx status would be pending revert if using v21 or before
	cctx := utils.WaitCctxAbortedByInboundHash(r.Ctx, r, txHash.String(), r.CctxClient)

	require.Equal(r, crosschaintypes.CctxStatus_Reverted, cctx.CctxStatus.Status)
	//require.True(r, cctx.GetCurrentOutboundParam().Amount.Uint64() < constant.BTCWithdrawalDustAmount)
	//require.True(r, strings.Contains(cctx.CctxStatus.ErrorMessage, crosschaintypes.ErrInvalidWithdrawalAmount.Error()))
}
