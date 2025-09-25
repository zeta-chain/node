package e2etests

import (
	"fmt"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestStressEtherDeposit tests the stressing deposit of ether
func TestStressEtherDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse deposit amount and number of deposits
	depositAmount := utils.ParseBigInt(r, args[0])
	numDeposits := utils.ParseInt(r, args[1])

	r.Logger.Print("starting stress test of %d deposits", numDeposits)

	// create a wait group to wait for all the deposits to complete
	var eg errgroup.Group

	revertOptions := gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)}

	// send the deposits
	for i := 0; i < numDeposits; i++ {
		i := i
		tx := r.ETHDeposit(
			r.EVMAddress(),
			depositAmount,
			revertOptions,
			false,
		)
		hash := tx.Hash()
		r.Logger.Print("index %d: starting deposit, tx hash: %s", i, hash.Hex())

		// slow down submitting transactions a bit.
		// submitting them as fast as possible does actually work.
		// but we want to ensure the workload is a bit more representative.
		time.Sleep(time.Millisecond * 10)

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
	timeToComplete := time.Since(startTime)
	r.Logger.Print("index %d: deposit cctx success in %s", index, timeToComplete.String())

	return nil
}
