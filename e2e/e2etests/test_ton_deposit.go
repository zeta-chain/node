package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
)

func TestTONDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	ctx := r.Ctx
	recipient := r.Account.EVMAddress()

	// Get TON ZRC20 balance before deposit
	balanceBefore, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, recipient)
	require.NoError(r, err)
	r.Logger.Info("Recipient's zEVM TON balance before deposit: %d", balanceBefore.Uint64())

	// Given gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given amount
	amount := utils.ParseUint(r, args[0])

	// Given approx deposit fee
	depositFee, err := gw.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDeposit)
	require.NoError(r, err)

	// Given a sender
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	// ACT
	cctx, err := r.TONDeposit(gw, sender, amount, recipient)
	require.NoError(r, err)

	// ASSERT
	// Check CCTX
	expectedDeposit := amount.Sub(depositFee)
	require.Equal(r, sender.GetAddress().ToRaw(), cctx.InboundParams.Sender)
	require.Equal(r, expectedDeposit.Uint64(), cctx.InboundParams.Amount.Uint64())

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(expectedDeposit.BigInt())
	utils.WaitAndVerifyZRC20BalanceChange(r, r.TONZRC20, recipient, balanceBefore, change, r.Logger)
}
