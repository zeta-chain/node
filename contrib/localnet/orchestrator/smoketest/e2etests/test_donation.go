package e2etests

import (
	"math/big"

	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	"github.com/zeta-chain/zetacore/zetaclient"
)

// TestDonationEther tests donation of ether to the tss address
func TestDonationEther(sm *runner.E2ERunner) {
	txDonation, err := sm.SendEther(sm.TSSAddress, big.NewInt(100000000000000000), []byte(zetaclient.DonationMessage))
	if err != nil {
		panic(err)
	}
	sm.Logger.EVMTransaction(*txDonation, "donation")

	// check contract deployment receipt
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, txDonation, sm.Logger, sm.ReceiptTimeout)
	sm.Logger.EVMReceipt(*receipt, "donation")
	if receipt.Status != 1 {
		panic("donation tx failed")
	}
}
