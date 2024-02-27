package e2etests

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

func TestBitcoinWithdrawToInvalidAddress(r *runner.E2ERunner) {
	// try to withdraw 0.00001 BTC from ZRC20 to BTC legacy address
	// first, approve the ZRC20 contract to spend 1 BTC from the deployer address
	WithdrawToInvalidAddress(r)
}

func WithdrawToInvalidAddress(r *runner.E2ERunner) {

	amount := big.NewInt(0.00001 * btcutil.SatoshiPerBitcoin)

	// approve the ZRC20 contract to spend 1 BTC from the deployer address
	tx, err := r.BTCZRC20.Approve(r.ZevmAuth, r.BTCZRC20Addr, big.NewInt(amount.Int64()*2)) // approve more to cover withdraw fee
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic(fmt.Errorf("approve receipt status is not 1"))
	}

	// mine blocks
	stop := r.MineBlocks()

	// withdraw 0.00001 BTC from ZRC20 to BTC legacy address
	tx, err = r.BTCZRC20.Withdraw(r.ZevmAuth, []byte("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3"), amount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 1 {
		panic(fmt.Errorf("withdraw receipt status is successful for an invalid BTC address"))
	}
	// stop mining
	stop <- struct{}{}
}
