package e2etests

import (
	"math/big"

	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

// TestDonationEther tests donation of ether to the tss address
func TestDonationEther(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestDonationEther requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestDonationEther.")
	}

	txDonation, err := r.SendEther(r.TSSAddress, amount, []byte(common.DonationMessage))
	if err != nil {
		panic(err)
	}
	r.Logger.EVMTransaction(*txDonation, "donation")

	// check contract deployment receipt
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txDonation, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "donation")
	if receipt.Status != 1 {
		panic("donation tx failed")
	}
}
