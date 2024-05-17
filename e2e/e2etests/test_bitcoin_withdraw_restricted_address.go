package e2etests

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestBitcoinWithdrawRestricted(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestBitcoinWithdrawRestricted requires exactly one argument for the amount.")
	}

	withdrawalAmount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		panic("Invalid withdrawal amount specified for TestBitcoinWithdrawRestricted.")
	}

	withdrawalAmountSat, err := btcutil.NewAmount(withdrawalAmount)
	if err != nil {
		panic(err)
	}
	amount := big.NewInt(int64(withdrawalAmountSat))

	r.SetBtcAddress(r.Name, false)

	withdrawBitcoinRestricted(r, amount)
}

func withdrawBitcoinRestricted(r *runner.E2ERunner, amount *big.Int) {
	// use restricted BTC P2WPKH address
	addressRestricted, err := chains.DecodeBtcAddress(testutils.RestrictedBtcAddressTest, chains.BtcRegtestChain.ChainId)
	if err != nil {
		panic(err)
	}

	// the cctx should be cancelled
	rawTx := withdrawBTCZRC20(r, addressRestricted, amount)
	if len(rawTx.Vout) != 2 {
		panic(fmt.Errorf("BTC cancelled outtx rawTx.Vout should have 2 outputs"))
	}
}
