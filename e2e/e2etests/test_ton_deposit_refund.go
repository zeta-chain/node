package e2etests

import (
	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/e2e/contracts/reverter"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestTONDepositAndCallRefund(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// Given amount and arbitrary call data
	var (
		amount = utils.ParseUint(r, args[0])
		data   = []byte("hello reverter")
	)

	// Given gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given deployer mock revert contract
	// deploy a reverter contract in ZEVM
	reverterAddr, _, _, err := testcontract.DeployReverter(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Reverter contract deployed at: %s", reverterAddr.String())

	// Given a sender
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	// Check sender balance
	ctx := r.Ctx
	senderBalance, err := r.Clients.TON.GetBalanceOf(ctx, sender.GetAddress(), false)
	if err != nil {
		r.Logger.Print("Failed to get sender balance: %v", err)
		require.NoError(r, err)
	}

	r.Logger.Print("Sender balance: %s", toncontracts.FormatCoins(senderBalance))

	// Get deposit fee
	depositFee, err := gw.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDepositAndCall)
	require.NoError(r, err)

	// Calculate total required amount (deposit + fee)
	totalRequired := amount.Add(depositFee)

	// Check if sender has enough balance
	if senderBalance.LT(totalRequired) {
		r.Logger.Print("⚠️ WARNING: Sender doesn't have enough TON to complete the deposit and call refund!")
		r.Logger.Print("Required: %s, Available: %s",
			toncontracts.FormatCoins(totalRequired),
			toncontracts.FormatCoins(senderBalance))
		r.Logger.Print("❓ This is expected when running without a faucet URL (ton_faucet: \"\")")
		r.Logger.Print("⏩ SKIPPING TEST: pre-conditions aren't met (insufficient balance).")
		return // Skip test instead of failing
	}

	// ACT
	// Send a deposit and call transaction from the deployer (faucet)
	// to the reverter contract
	cctx, err := r.TONDepositAndCall(
		gw,
		sender,
		amount,
		reverterAddr,
		data,
		runner.TONExpectStatus(cctypes.CctxStatus_Reverted),
	)

	// ASSERT
	require.NoError(r, err)
	r.Logger.CCTX(*cctx, "ton_deposit_and_refund")

	require.Contains(r, cctx.CctxStatus.ErrorMessage, utils.ErrHashRevertFoo)
}
