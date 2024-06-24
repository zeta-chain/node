package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/constant"
)

// TestDonationEther tests donation of ether to the tss address
func TestDonationEther(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestDonationEther.")

	txDonation, err := r.SendEther(r.TSSAddress, amount, []byte(constant.DonationMessage))
	require.NoError(r, err)

	r.Logger.EVMTransaction(*txDonation, "donation")

	// check contract deployment receipt
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txDonation, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "donation")
	utils.RequireTxSuccessful(r, receipt)
}
