package e2etests

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

// TestStressSuiDeposit tests the stressing deposit of SUI
func TestStressSuiDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	amount := utils.ParseBigInt(r, args[0])
	numDeposits := utils.ParseInt(r, args[1])

	r.Logger.Print("starting stress test of %d SUI deposits", numDeposits)

	// create a wait group to wait for all the deposits to complete
	var eg errgroup.Group

	// send the SUI deposit transactions
	for i := range numDeposits {
		// each goroutine captures its own copy of i
		resp := r.SuiDepositSUI(r.SuiGateway.PackageID(), r.EVMAddress(), math.NewUintFromBigInt(amount))

		r.Logger.Print("index %d: started with tx hash: %s", i, resp.Digest)

		eg.Go(func() error {
			_, err := r.SuiMonitorCCTXByInboundHash(resp.Digest, i)
			if err != nil {
				return errors.Wrap(err, "failed to monitor deposit")
			}
			return nil
		})
	}

	require.NoError(r, eg.Wait())

	r.Logger.Print("all SUI deposits completed")
}
