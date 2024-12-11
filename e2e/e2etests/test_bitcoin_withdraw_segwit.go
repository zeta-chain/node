package e2etests

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestBitcoinWithdrawSegWit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse arguments
	defaultReceiver := r.BTCDeployerAddress.EncodeAddress()
	receiver, amount := utils.ParseBitcoinWithdrawArgs(r, args, defaultReceiver, r.GetBitcoinChainID())
	_, ok := receiver.(*btcutil.AddressWitnessPubKeyHash)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawSegWit.")

	withdrawBTCZRC20(r, receiver, amount)
}
