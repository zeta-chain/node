package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestTONDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	ctx := r.Ctx

	// Given gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Log important gateway information
	r.Logger.Print("🔍 Test using TON Gateway address: %s", gw.AccountID().ToRaw())
	r.Logger.Print("🔍 Runner's TON Gateway address: %s", r.TONGateway.ToRaw())

	// Verify chain parameters have the correct gateway
	chainParams, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, &types.QueryGetChainParamsForChainRequest{
		ChainId: chains.TONLocalnet.ChainId,
	})
	if err != nil {
		r.Logger.Print("🔍 Failed to get chain params: %v", err)
	} else {
		r.Logger.Print("🔍 ZetaCore has TON Gateway address: %s", chainParams.ChainParams.GatewayAddress)
		r.Logger.Print("🔍 Gateway matches test gateway: %v", chainParams.ChainParams.GatewayAddress == gw.AccountID().ToRaw())
	}

	// Given amount
	amount := utils.ParseUint(r, args[0])

	// Debug messages
	_, s, err := r.Account.AsTONWallet(r.Clients.TON)
	r.Logger.Print("Amount: %s", amount.String())
	r.Logger.Print("Address: %s", s.GetAddress().ToHuman(false, true))
	r.Logger.Print("Gateway Account: %s", gw.AccountID().ToRaw())
	r.Logger.Print("TSS Address: %s", r.TSSAddress.Hex())
	r.Logger.Print("Authority Address: %s", r.Account.EVMAddress().Hex())

	// Verify TSS and authority addresses
	expectedTSS := r.TSSAddress
	expectedAuthority := r.Account.EVMAddress()
	r.Logger.Print("Expected TSS Address: %s", expectedTSS.Hex())
	r.Logger.Print("Expected Authority Address: %s", expectedAuthority.Hex())
	r.Logger.Print("TSS Address Match: %v", r.TSSAddress.Hex() == expectedTSS.Hex())
	r.Logger.Print("Authority Address Match: %v", r.Account.EVMAddress().Hex() == expectedAuthority.Hex())

	// Check Gateway contract state
	state, err := r.Clients.TON.GetAccountState(ctx, gw.AccountID())
	if err != nil {
		r.Logger.Print("Failed to get Gateway state: %v", err)
	} else {
		r.Logger.Print("Gateway state: %+v", state)
	}

	// Given approx deposit fee
	depositFee, err := gw.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDeposit)
	if err != nil {
		r.Logger.Print("Failed to retrieve deposit fee: %v (fee: %s, address: %s, account: %s)", err, depositFee.String(), s.GetAddress().ToHuman(false, true), gw.AccountID().ToRaw())
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
	r.Logger.Print("🔍 Sending TON deposit to gateway: %s", gw.AccountID().ToRaw())
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
