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

	r.Logger.Info("Starting TON Deposit and Call Refund test with amount: %s nano TON", amount.String())

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

	// ACT
	// Send a deposit and call transaction from the deployer (faucet)
	// to the reverter contract
	r.Logger.Info("Sending deposit and call with %s nano TON and expecting Reverted status", amount.String())
	cctx, err := r.TONDepositAndCall(
		gw,
		sender,
		amount,
		reverterAddr,
		data,
		runner.TONExpectStatus(cctypes.CctxStatus_Reverted),
		// No revert gas limit for now - let's first get it working with defaults
	)

	// ASSERT
	require.NoError(r, err)
	r.Logger.Info("Received CCTX with status: %s and error message: %s",
		cctx.CctxStatus.Status.String(),
		cctx.CctxStatus.ErrorMessage)
	r.Logger.CCTX(*cctx, "ton_deposit_and_refund")

	// Check for the Foo() error hash
	require.Contains(r, cctx.CctxStatus.ErrorMessage, utils.ErrHashRevertFoo,
		"Error message should contain the Foo() error hash")
}
