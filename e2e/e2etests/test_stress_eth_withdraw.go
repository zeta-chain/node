package e2etests

import (
	"fmt"
	"strconv"

	"math/big"
	"time"

	"golang.org/x/sync/errgroup"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// TestStressEtherWithdraw tests the stressing withdraw of ether
func TestStressEtherWithdraw(r *runner.E2ERunner, args []string) {
	if len(args) != 2 {
		panic("TestStressEtherWithdraw requires exactly two arguments: the withdrawal amount and the number of withdrawals.")
	}

	withdrawalAmount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid withdrawal amount specified for TestStressEtherWithdraw.")
	}

	numWithdraws, err := strconv.Atoi(args[1])
	if err != nil || numWithdraws < 1 {
		panic("Invalid number of withdrawals specified for TestStressEtherWithdraw.")
	}

	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ETHZRC20Addr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	r.WaitForTxReceiptOnZEVM(tx)

	r.Logger.Print("starting stress test of %d withdraws", numWithdraws)

	// create a wait group to wait for all the withdraws to complete
	var eg errgroup.Group

	// send the withdraws
	for i := 0; i < numWithdraws; i++ {
		i := i
		tx, err := r.ETHZRC20.Withdraw(r.ZEVMAuth, r.DeployerAddress.Bytes(), withdrawalAmount)
		if err != nil {
			panic(err)
		}
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		if receipt.Status == 0 {
			//r.Logger.Info("index %d: withdraw evm tx failed", index)
			panic(fmt.Sprintf("index %d: withdraw evm tx %s failed", i, tx.Hash().Hex()))
		}
		r.Logger.Print("index %d: starting withdraw, tx hash: %s", i, tx.Hash().Hex())

		eg.Go(func() error {
			return monitorEtherWithdraw(r, tx, i, time.Now())
		})
	}

	// wait for all the withdraws to complete
	if err := eg.Wait(); err != nil {
		panic(err)
	}

	r.Logger.Print("all withdraws completed")
}

// monitorEtherWithdraw monitors the withdraw of ether, returns once the withdraw is complete
func monitorEtherWithdraw(r *runner.E2ERunner, tx *ethtypes.Transaction, index int, startTime time.Time) error {
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.ReceiptTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return fmt.Errorf(
			"index %d: withdraw cctx failed with status %s, message %s, cctx index %s",
			index,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
			cctx.Index,
		)
	}
	timeToComplete := time.Now().Sub(startTime)
	r.Logger.Print("index %d: withdraw cctx success in %s", index, timeToComplete.String())

	return nil
}
