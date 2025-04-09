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
	r.Logger.Print("🔍 Gateway address match: %v", gw.AccountID().ToRaw() == r.TONGateway.ToRaw())

	// Try to verify chain parameters
	r.Logger.Print("🔍 Checking chain parameters...")
	chainID := chains.TONTestnet.ChainId
	r.Logger.Print("🔍 Using TON testnet chain ID: %d", chainID)

	chainParams, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, &types.QueryGetChainParamsForChainRequest{
		ChainId: chainID,
	})

	if err != nil {
		r.Logger.Print("🔍 Failed to get chain params: %v", err)

		// Force getting the chain parameters from the observer module
		r.Logger.Print("🔍 Trying to get all chain params...")
		allParams, paramErr := r.ObserverClient.GetChainParams(r.Ctx, &types.QueryGetChainParamsRequest{})
		if paramErr != nil {
			r.Logger.Print("❌ Failed to get any chain params: %v", paramErr)
		} else if allParams != nil && allParams.ChainParams != nil {
			r.Logger.Print("✅ Found chain params")
			r.Logger.Print("🔍 Chain params: %+v", allParams.ChainParams)
		} else {
			r.Logger.Print("⚠️ No chain params found")
		}
	} else {
		r.Logger.Print("✅ Successfully retrieved chain parameters")
		r.Logger.Print("🔍 ZetaCore has TON Gateway address: %s", chainParams.ChainParams.GatewayAddress)
		r.Logger.Print("🔍 Gateway matches test gateway: %v", chainParams.ChainParams.GatewayAddress == gw.AccountID().ToRaw())

		if chainParams.ChainParams.GatewayAddress != gw.AccountID().ToRaw() {
			r.Logger.Print("⚠️ Gateway address mismatch, this may cause test failure!")
			r.Logger.Print("🔍 Expected: %s, Got: %s", gw.AccountID().ToRaw(), chainParams.ChainParams.GatewayAddress)
		}
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

	// Verify chain parameters one more time before deposit
	chainParams, err = r.ObserverClient.GetChainParamsForChain(r.Ctx, &types.QueryGetChainParamsForChainRequest{
		ChainId: chainID,
	})
	if err != nil {
		r.Logger.Print("⚠️ Final check: Chain parameters still not set, test will likely fail")
	} else {
		r.Logger.Print("✅ Final check: Chain parameters are set with gateway: %s", chainParams.ChainParams.GatewayAddress)
		r.Logger.Print("🔍 Final check: Gateway match: %v", chainParams.ChainParams.GatewayAddress == gw.AccountID().ToRaw())
	}

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
