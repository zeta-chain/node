package e2etests

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestBitcoinWithdrawLegacy(r *runner.E2ERunner, args []string) {
	// check length of arguments
	require.Len(r, args, 2)

	// parse arguments and withdraw BTC
	receiver, amount := utils.ParseBitcoinWithdrawArgs(r, args, defaultReceiver, r.GetBitcoinChainID())

	_, ok := receiver.(*btcutil.AddressPubKeyHash)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawLegacy.")

	withdrawBTCZRC20(r, receiver, amount)
}
