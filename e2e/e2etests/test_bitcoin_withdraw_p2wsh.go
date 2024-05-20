package e2etests

import (
	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

func TestBitcoinWithdrawP2WSH(r *runner.E2ERunner, args []string) {
	// check length of arguments
	if len(args) != 2 {
		panic("TestBitcoinWithdrawP2WSH requires two arguments: [receiver, amount]")
	}
	r.SetBtcAddress(r.Name, false)

	// parse arguments and withdraw BTC
	defaultReceiver := "bcrt1qm9mzhyky4w853ft2ms6dtqdyyu3z2tmrq8jg8xglhyuv0dsxzmgs2f0sqy"
	receiver, amount := parseBitcoinWithdrawArgs(args, defaultReceiver)
	_, ok := receiver.(*btcutil.AddressWitnessScriptHash)
	if !ok {
		panic("Invalid receiver address specified for TestBitcoinWithdrawP2WSH.")
	}

	withdrawBTCZRC20(r, receiver, amount)
}
