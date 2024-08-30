package e2etests

import (
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
)

func TestBitcoinWithdrawLegacy(r *runner.E2ERunner, args []string) {
	// check length of arguments
	require.Len(r, args, 2)

	r.SetBtcAddress(r.Name, false)

	// parse arguments and withdraw BTC
	receiver, amount := parseBitcoinWithdrawArgs(r, args, defaultReceiver)

	_, ok := receiver.(*btcutil.AddressPubKeyHash)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawLegacy.")

	withdrawBTCZRC20(r, receiver, amount)
}
