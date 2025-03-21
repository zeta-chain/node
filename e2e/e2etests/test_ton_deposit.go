package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestTONDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	ctx := r.Ctx

	// Given gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given amount
	amount := utils.ParseUint(r, args[0])

	// Given approx deposit fee
	depositFee, err := gw.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDeposit)
	if err != nil {
		r.Logger.Print("Failed to retrieve deposit fee: %v", err)
		require.NoError(r, err)
	}

	// Debugging: Log deposit fee
	r.Logger.Print("Deposit fee: %s", depositFee.String())

	// Given a sender
	r.Logger.Print("Preparing to call AsTONWallet...")
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	if err != nil {
		r.Logger.Print("Failed to retrieve TON Wallet: %v", err)
	}
	require.NoError(r, err)

	// Debugging: Log sender address
	r.Logger.Print("Sender TON address: %s", sender.GetAddress().ToRaw())

	// Given sample EVM address
	recipient := sample.EthAddress()

	// ACT
	cctx, err := r.TONDeposit(gw, sender, amount, recipient)

	// ASSERT
	require.NoError(r, err)

	// Check CCTX
	expectedDeposit := amount.Sub(depositFee)

	require.Equal(r, sender.GetAddress().ToRaw(), cctx.InboundParams.Sender)
	require.Equal(r, expectedDeposit.Uint64(), cctx.InboundParams.Amount.Uint64())

	// Check receiver's balance
	balance, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, recipient)
	require.NoError(r, err)

	r.Logger.Info("Recipient's zEVM TON balance after deposit: %d", balance.Uint64())

	require.Equal(r, expectedDeposit.Uint64(), balance.Uint64())
}
