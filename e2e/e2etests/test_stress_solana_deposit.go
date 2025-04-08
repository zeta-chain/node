package e2etests

import (
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestStressSolanaDeposit tests the stressing deposit of SOL
func TestStressSolanaDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	depositSOLAmount := utils.ParseBigInt(r, args[0])
	numDepositsSOL := utils.ParseInt(r, args[1])

	r.Logger.Print("starting stress test of %d SOL deposits", numDepositsSOL)

	// create a wait group to wait for all the deposits to complete
	var eg errgroup.Group

	// send the deposits SOL
	for i := 0; i < numDepositsSOL; i++ {
		i := i

		// execute the deposit SOL transaction
		sig := r.SOLDepositAndCall(nil, r.EVMAddress(), depositSOLAmount, nil, nil)
		r.Logger.Print("index %d: starting SOL deposit, sig: %s", i, sig.String())

		eg.Go(func() error { return monitorDeposit(r, sig, i, time.Now()) })
	}

	require.NoError(r, eg.Wait())

	r.Logger.Print("all SOL deposits completed")
}

// monitorDeposit monitors the deposit of SOL/SPL, returns once the deposit is complete
func monitorDeposit(r *runner.E2ERunner, sig solana.Signature, index int, startTime time.Time) error {
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.ReceiptTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return fmt.Errorf(
			"index %d: deposit cctx failed with status %s, message %s, cctx index %s",
			index,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
			cctx.Index,
		)
	}
	timeToComplete := time.Since(startTime)
	r.Logger.Print("index %d: deposit cctx success in %s", index, timeToComplete.String())

	return nil
}
