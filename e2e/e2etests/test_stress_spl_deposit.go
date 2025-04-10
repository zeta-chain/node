package e2etests

import (
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

// TestStressSPLDeposit tests the stressing deposit of SPL
func TestStressSPLDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	depositSPLAmount := utils.ParseBigInt(r, args[0])
	numDepositsSPL := utils.ParseInt(r, args[1])

	// load deployer private key
	privKey := r.GetSolanaPrivKey()

	r.Logger.Print("starting stress test of %d SPL deposits", numDepositsSPL)

	// create a wait group to wait for all the deposits to complete
	var eg errgroup.Group

	// send the deposits SPL
	for i := 0; i < numDepositsSPL; i++ {
		i := i

		// execute the deposit SPL transaction
		sig := r.SPLDepositAndCall(&privKey, depositSPLAmount.Uint64(), r.SPLAddr, r.EVMAddress(), nil, nil)
		r.Logger.Print("index %d: starting SPL deposit, sig: %s", i, sig.String())

		eg.Go(func() error { return monitorDeposit(r, sig, i, time.Now()) })
	}

	require.NoError(r, eg.Wait())

	r.Logger.Print("all SPL deposits completed")
}
