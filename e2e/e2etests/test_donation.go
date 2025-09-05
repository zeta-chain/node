package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/constant"
)

// TestDonationEther tests donation of ether to the tss address
func TestDonationEther(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the donation amount
	amount := utils.ParseBigInt(r, args[0])

	txDonation, err := r.LegacySendEther(r.TSSAddress, amount, []byte(constant.DonationMessage))
	require.NoError(r, err)

	r.Logger.EVMTransaction(txDonation, "donation")

	// check contract deployment receipt
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txDonation, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "donation")
	utils.RequireTxSuccessful(r, receipt)
}
