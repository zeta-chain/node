package e2etests

import (
	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

func TestBitcoinWithdrawP2SH(r *runner.E2ERunner, args []string) {
	// check length of arguments
	if len(args) != 2 {
		panic("TestBitcoinWithdrawP2SH requires two arguments: [receiver, amount]")
	}
	r.SetBtcAddress(r.Name, false)

	// parse arguments and withdraw BTC
	defaultReceiver := "2N6AoUj3KPS7wNGZXuCckh8YEWcSYNsGbqd"
	receiver, amount := parseBitcoinWithdrawArgs(args, defaultReceiver)
	_, ok := receiver.(*btcutil.AddressScriptHash)
	if !ok {
		panic("Invalid receiver address specified for TestBitcoinWithdrawP2SH.")
	}

	withdrawBTCZRC20(r, receiver, amount)
}
