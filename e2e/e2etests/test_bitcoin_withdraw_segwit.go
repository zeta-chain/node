package e2etests

import (
	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

func TestBitcoinWithdrawSegWit(r *runner.E2ERunner, args []string) {
	// check length of arguments
	if len(args) != 2 {
		panic("TestBitcoinWithdrawSegWit requires two arguments: [receiver, amount]")
	}
	r.SetBtcAddress(r.Name, false)

	// parse arguments
	defaultReceiver := r.BTCDeployerAddress.EncodeAddress()
	receiver, amount := parseBitcoinWithdrawArgs(args, defaultReceiver)
	_, ok := receiver.(*btcutil.AddressWitnessPubKeyHash)
	if !ok {
		panic("Invalid receiver address specified for TestBitcoinWithdrawSegWit.")
	}

	withdrawBTCZRC20(r, receiver, amount)
}
