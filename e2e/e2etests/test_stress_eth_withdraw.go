package e2etests

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestStressEtherWithdraw tests the stressing withdraw of ether
func TestStressEtherWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse withdraw amount and number of withdraws
	withdrawalAmount := parseBigInt(r, args[0])

	numWithdraws, err := strconv.Atoi(args[1])
	require.NoError(r, err)
	require.GreaterOrEqual(r, numWithdraws, 1)

	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ETHZRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)

	r.WaitForTxReceiptOnZEVM(tx)

	r.Logger.Print("starting stress test of %d withdraws", numWithdraws)

	// create a wait group to wait for all the withdraws to complete
	var eg errgroup.Group

	// send the withdraws
	for i := 0; i < numWithdraws; i++ {
		i := i

		tx, err := r.ETHZRC20.Withdraw(r.ZEVMAuth, r.EVMAddress().Bytes(), withdrawalAmount)
		require.NoError(r, err)

		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt)

		r.Logger.Print("index %d: starting withdraw, tx hash: %s", i, tx.Hash().Hex())

		eg.Go(func() error {
			return monitorEtherWithdraw(r, tx, i, time.Now())
		})
	}

	require.NoError(r, eg.Wait())

	r.Logger.Print("all withdraws completed")
}

// monitorEtherWithdraw monitors the withdraw of ether, returns once the withdraw is complete
func monitorEtherWithdraw(r *runner.E2ERunner, tx *ethtypes.Transaction, index int, startTime time.Time) error {
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.ReceiptTimeout)
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
