package e2etests

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestStressSolanaWithdraw tests the stressing withdrawal of SOL/SPL
func TestStressSolanaWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 4)

	withdrawSOLAmount := utils.ParseBigInt(r, args[0])
	numWithdrawalsSOL := utils.ParseInt(r, args[1])
	withdrawSPLAmount := utils.ParseBigInt(r, args[2])
	numWithdrawalsSPL := utils.ParseInt(r, args[3])

	balanceBefore, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL before withdraw: %s", balanceBefore.String())

	balanceBefore, err = r.SPLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SPL before withdraw: %s", balanceBefore.String())

	// load deployer private key
	privKey := r.GetSolanaPrivKey()

	r.Logger.Print("starting stress test of %d SOL and %d SPL withdrawals", numWithdrawalsSOL, numWithdrawalsSPL)

	tx, err := r.SOLZRC20.Approve(r.ZEVMAuth, r.SOLZRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve_sol")

	tx, err = r.SPLZRC20.Approve(r.ZEVMAuth, r.SPLZRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve_spl")

	tx, err = r.SOLZRC20.Approve(r.ZEVMAuth, r.SPLZRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve_spl_sol")

	// create a wait group to wait for all the withdrawals to complete
	var eg errgroup.Group

	// send the withdrawals SOL
	for i := 0; i < numWithdrawalsSOL; i++ {
		i := i

		// execute the withdraw SOL transaction
		tx, err = r.SOLZRC20.Withdraw(r.ZEVMAuth, []byte(privKey.PublicKey().String()), withdrawSOLAmount)
		require.NoError(r, err)

		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt)

		r.Logger.Print("index %d: starting SOL withdraw, tx hash: %s", i, tx.Hash().Hex())

		eg.Go(func() error { return monitorWithdrawal(r, tx.Hash(), i, time.Now()) })
	}

	// send the withdrawals SPL
	for i := 0; i < numWithdrawalsSPL; i++ {
		i := i

		// execute the withdraw SPL transaction
		tx, err = r.SPLZRC20.Withdraw(r.ZEVMAuth, []byte(privKey.PublicKey().String()), withdrawSPLAmount)
		require.NoError(r, err)

		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt)

		r.Logger.Print("index %d: starting SPL withdraw, tx hash: %s", i, tx.Hash().Hex())

		eg.Go(func() error { return monitorWithdrawal(r, tx.Hash(), i, time.Now()) })
	}

	require.NoError(r, eg.Wait())

	r.Logger.Print("all withdrawals completed")
}

// monitorWithdrawal monitors the withdrawal of SOL/SPL, returns once the withdrawal is complete
func monitorWithdrawal(r *runner.E2ERunner, hash ethcommon.Hash, index int, startTime time.Time) error {
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, hash.String(), r.CctxClient, r.Logger, r.ReceiptTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return fmt.Errorf(
			"index %d: withdraw cctx failed with status %s, message %s, cctx index %s",
			index,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
			cctx.Index,
		)
	}
	timeToComplete := time.Since(startTime)
	r.Logger.Print("index %d: withdraw cctx success in %s", index, timeToComplete.String())

	return nil
}
