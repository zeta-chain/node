package e2etests

import (
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
)

func TestBitcoinWithdrawSegWit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	r.SetBtcAddress(r.Name, false)

	// parse arguments
	defaultReceiver := r.BTCDeployerAddress.EncodeAddress()
	receiver, amount := parseBitcoinWithdrawArgs(r, args, defaultReceiver)
	_, ok := receiver.(*btcutil.AddressWitnessPubKeyHash)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawSegWit.")

	withdrawBTCZRC20(r, receiver, amount)
}
