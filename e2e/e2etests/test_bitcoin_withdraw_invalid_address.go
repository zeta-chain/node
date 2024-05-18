package e2etests

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

func TestBitcoinWithdrawToInvalidAddress(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestBitcoinWithdrawToInvalidAddress requires exactly one argument for the amount.")
	}

	withdrawalAmount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		panic("Invalid withdrawal amount specified for TestBitcoinWithdrawToInvalidAddress.")
	}

	withdrawalAmountSat, err := btcutil.NewAmount(withdrawalAmount)
	if err != nil {
		panic(err)
	}
	amount := big.NewInt(int64(withdrawalAmountSat))

	r.SetBtcAddress(r.Name, false)

	withdrawToInvalidAddress(r, amount)
}

func withdrawToInvalidAddress(r *runner.E2ERunner, amount *big.Int) {
	approvalAmount := 1000000000000000000
	// approve the ZRC20 contract to spend approvalAmount BTC from the deployer address.
	// the actual amount transferred is provided as test arg BTC, but we approve more to cover withdraw fee
	tx, err := r.BTCZRC20.Approve(r.ZEVMAuth, r.BTCZRC20Addr, big.NewInt(int64(approvalAmount)))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic(fmt.Errorf("approve receipt status is not 1"))
	}

	// mine blocks
	stop := r.MineBlocks()

	// withdraw amount provided as test arg BTC from ZRC20 to BTC legacy address
	// the address "1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3" is for mainnet, not regtest
	tx, err = r.BTCZRC20.Withdraw(r.ZEVMAuth, []byte("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3"), amount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 1 {
		panic(fmt.Errorf("withdraw receipt status is successful for an invalid BTC address"))
	}
	// stop mining
	stop <- struct{}{}
}
