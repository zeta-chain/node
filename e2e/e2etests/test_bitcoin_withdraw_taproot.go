package e2etests

import (
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

func TestBitcoinWithdrawTaproot(r *runner.E2ERunner, args []string) {
	// check length of arguments
	if len(args) != 2 {
		panic("TestBitcoinWithdrawTaproot requires two arguments: [receiver, amount]")
	}
	r.SetBtcAddress(r.Name, false)

	// parse arguments and withdraw BTC
	defaultReceiver := "bcrt1pqqqsyqcyq5rqwzqfpg9scrgwpugpzysnzs23v9ccrydpk8qarc0sj9hjuh"
	receiver, amount := parseBitcoinWithdrawArgs(args, defaultReceiver)
	_, ok := receiver.(*chains.AddressTaproot)
	if !ok {
		panic("Invalid receiver address specified for TestBitcoinWithdrawTaproot.")
	}

	withdrawBTCZRC20(r, receiver, amount)
}
