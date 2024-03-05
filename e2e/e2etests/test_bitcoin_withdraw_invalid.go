package e2etests

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

func TestBitcoinWithdrawToInvalidAddress(r *runner.E2ERunner) {
	WithdrawToInvalidAddress(r)
}

func WithdrawToInvalidAddress(r *runner.E2ERunner) {

	amount := big.NewInt(0.00001 * btcutil.SatoshiPerBitcoin)
	approvalAmount := 1000000000000000000
	// approve the ZRC20 contract to spend approvalAmount BTC from the deployer address.
	// the actual amount transferred is 0.00001 BTC, but we approve more to cover withdraw fee
	tx, err := r.BTCZRC20.Approve(r.ZevmAuth, r.BTCZRC20Addr, big.NewInt(int64(approvalAmount)))
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
