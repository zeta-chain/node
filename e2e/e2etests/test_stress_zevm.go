package e2etests

import (
	"github.com/zeta-chain/node/e2e/contracts/testgasconsumer"
	"math/big"
	"strings"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

// TestStressZEVM tests stressing direct interactions with the zEVM using calls that consume a lot of gas
func TestStressZEVM(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse deposit amount and number of deposits
	txNumbers := utils.ParseBigInt(r, args[0])

	r.Logger.Print("starting stress test of %d calls", txNumbers)

	// create a wait group to wait for all the deposits to complete
	//var eg errgroup.Group

	// Deploy the GasConsumer contract
	// the target provided in the contract is not accurate, this value allows to get a gas usage close to the limit of 4M
	gasConsumerAddress, txDeploy, gasConsumer, err := testgasconsumer.DeployTestGasConsumer(
		r.ZEVMAuth,
		r.ZEVMClient,
		big.NewInt(1000000),
	)
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(txDeploy)

	// send the deposits
	for i := uint64(0); i < txNumbers.Uint64(); i++ {
		_, err := gasConsumer.OnCall(
			r.ZEVMAuth,
			testgasconsumer.TestGasConsumerzContext{
				Origin:  []byte{},
				Sender:  gasConsumerAddress,
				ChainID: big.NewInt(0),
			},
			gasConsumerAddress,
			big.NewInt(0),
			[]byte{},
		)
		if err != nil && strings.Contains(err.Error(), "invalid nonce") {
			if strings.Contains(err.Error(), "invalid nonce") {
				// nonce issue happen because of fast submissions, just skip
				r.Logger.Print("index %d: skipped for invalid nonce", i)
				time.Sleep(time.Second)
				continue
			}
			require.Fail(r, "failed to call: %v", err)
		}

		r.Logger.Print("index %d: starting stress zevm call", i)

		// slow down submitting transactions a bit.
		time.Sleep(time.Millisecond * 50)

		//eg.Go(func() error { return monitorEtherDeposit(r, hash, i, time.Now()) })
	}

	//require.NoError(r, eg.Wait())

	r.Logger.Print("all calls completed")
}
