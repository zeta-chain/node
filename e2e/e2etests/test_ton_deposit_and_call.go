package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
)

func TestTONDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	ctx := r.Ctx

	// Using gateway address from config.yml (specified as gateway_account_id)
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given amount
	amount := utils.ParseUint(r, args[0])

	// Given approx depositAndCall fee
	depositFee, err := gw.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDepositAndCall)
	require.NoError(r, err)

	// Given a sender
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	// Get zEVM address as the recipient
	recipientAddr := r.EVMAddress()

	// Using the TONZRC20 address from config.yml (specified as ton_zrc20)
	contractAddr := r.TONZRC20Addr
	r.Logger.Info("Using TON ZRC20 contract address %s and TON gateway address %s", contractAddr.Hex(), r.TONGateway.ToRaw())
	r.Logger.Info("Recipient address (user's EVM address): %s", recipientAddr.Hex())

	balanceBefore, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, contractAddr)
	require.NoError(r, err)
	r.Logger.Info("Recipient's zEVM TON balance before deposit (contract address: %s): %d (0x%x)", contractAddr.Hex(), balanceBefore.Uint64(), balanceBefore.Uint64())

	// Given call data
	callData := []byte("hello from TON!")

	// Call TONDepositAndCall with the contract address
	_, err = r.TONDepositAndCall(gw, sender, amount, contractAddr, callData)

	// ASSERT
	require.NoError(r, err)

	expectedDeposit := amount.Sub(depositFee)

	// Check the balance after deposit
	balanceAfter, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, contractAddr)
	require.NoError(r, err)
	r.Logger.Info("Recipient's zEVM TON balance after deposit: %d (0x%x)", balanceAfter.Uint64(), balanceAfter.Uint64())

	// Calculate and log expected deposit amount (amount minus fee)

	r.Logger.Info("Expected deposit amount: %d (0x%x)", expectedDeposit.Uint64(), expectedDeposit.Uint64())

	require.Equal(r, expectedDeposit.Uint64(), balanceAfter.Uint64())

}
