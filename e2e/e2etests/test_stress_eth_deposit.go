package e2etests

import (
	"fmt"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestStressEtherDeposit tests the stressing deposit of ether
func TestStressEtherDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse deposit amount and number of deposits
	depositAmount := parseBigInt(r, args[0])
	numDeposits := parseInt(r, args[1])

	r.Logger.Print("starting stress test of %d deposits", numDeposits)

	// create a wait group to wait for all the deposits to complete
	var eg errgroup.Group

	// send the deposits
	for i := 0; i < numDeposits; i++ {
		i := i
		hash := r.DepositEtherWithAmount(depositAmount)
		r.Logger.Print("index %d: starting deposit, tx hash: %s", i, hash.Hex())

		eg.Go(func() error { return monitorEtherDeposit(r, hash, i, time.Now()) })
	}

	require.NoError(r, eg.Wait())

	r.Logger.Print("all deposits completed")
}

// monitorEtherDeposit monitors the deposit of ether, returns once the deposit is complete
func monitorEtherDeposit(r *runner.E2ERunner, hash ethcommon.Hash, index int, startTime time.Time) error {
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.ReceiptTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return fmt.Errorf(
			"index %d: deposit cctx failed with status %s, message %s, cctx index %s",
			index,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
			cctx.Index,
		)
	}
	timeToComplete := time.Now().Sub(startTime)
	r.Logger.Print("index %d: deposit cctx success in %s", index, timeToComplete.String())

	return nil
}
