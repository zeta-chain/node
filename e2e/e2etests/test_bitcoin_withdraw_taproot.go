package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/pkg/chains"
)

func TestBitcoinWithdrawTaproot(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	r.SetBtcAddress(r.Name, false)

	// parse arguments and withdraw BTC
	defaultReceiver := "bcrt1pqqqsyqcyq5rqwzqfpg9scrgwpugpzysnzs23v9ccrydpk8qarc0sj9hjuh"
	receiver, amount := parseBitcoinWithdrawArgs(r, args, defaultReceiver)
	_, ok := receiver.(*chains.AddressTaproot)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawTaproot.")

	withdrawBTCZRC20(r, receiver, amount)
}
