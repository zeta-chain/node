package e2etests

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestBitcoinWithdrawTaproot(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse arguments and withdraw BTC
	defaultReceiver := "bcrt1pqqqsyqcyq5rqwzqfpg9scrgwpugpzysnzs23v9ccrydpk8qarc0sj9hjuh"
	receiver, amount := utils.ParseBitcoinWithdrawArgs(r, args, defaultReceiver, r.GetBitcoinChainID())
	_, ok := receiver.(*btcutil.AddressTaproot)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawTaproot.")

	withdrawBTCZRC20(r, receiver, amount)
}
