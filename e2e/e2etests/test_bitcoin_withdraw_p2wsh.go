package e2etests

import (
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
)

func TestBitcoinWithdrawP2WSH(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	r.SetBtcAddress(r.Name, false)

	// parse arguments and withdraw BTC
	defaultReceiver := "bcrt1qm9mzhyky4w853ft2ms6dtqdyyu3z2tmrq8jg8xglhyuv0dsxzmgs2f0sqy"
	receiver, amount := parseBitcoinWithdrawArgs(r, args, defaultReceiver)
	_, ok := receiver.(*btcutil.AddressWitnessScriptHash)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawP2WSH.")

	withdrawBTCZRC20(r, receiver, amount)
}
