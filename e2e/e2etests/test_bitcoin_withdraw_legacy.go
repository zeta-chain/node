package e2etests

import (
	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

func TestBitcoinWithdrawLegacy(r *runner.E2ERunner, args []string) {
	// check length of arguments
	if len(args) != 2 {
		panic("TestBitcoinWithdrawLegacy requires two arguments: [receiver, amount]")
	}
	r.SetBtcAddress(r.Name, false)

	// parse arguments and withdraw BTC
	defaultReceiver := "mxpYha3UJKUgSwsAz2qYRqaDSwAkKZ3YEY"
	receiver, amount := parseBitcoinWithdrawArgs(args, defaultReceiver)
	_, ok := receiver.(*btcutil.AddressPubKeyHash)
	if !ok {
		panic("Invalid receiver address specified for TestBitcoinWithdrawLegacy.")
	}

	withdrawBTCZRC20(r, receiver, amount)
}
