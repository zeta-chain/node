package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestTONDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// Given deployer
	ctx, deployer := r.Ctx, r.TONDeployer

	// Given amount
	amount := parseUint(r, args[0])

	// Given approx deposit fee
	depositFee, err := r.TONGateway.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDeposit)
	require.NoError(r, err)

	// Given sample wallet with a balance of 50 TON
	sender, err := deployer.CreateWallet(ctx, toncontracts.Coins(50))
	require.NoError(r, err)

	// Given sample EVM address
	recipient := sample.EthAddress()

	// ACT
	cctx, err := r.TONDeposit(sender, amount, recipient)

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
