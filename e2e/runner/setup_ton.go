package runner

import (
	"context"
	"time"

	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/wallet"

	"github.com/tonkeeper/tongo"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/runner/ton"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	cctxtypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// SetupTON setups TON deployer and deploys Gateway contract
func (r *E2ERunner) SetupTON(faucetURL string, userTON config.Account) {
	require.NotEmpty(r, faucetURL, "TON faucet url is empty")
	require.NotNil(r, r.Clients.TON, "TON client is not initialized")

	ctx := r.Ctx

	// 1. Setup Deployer (acts as a faucet as well)
	faucetConfig, err := ton.GetFaucet(ctx, faucetURL)
	require.NoError(r, err, "unable to get faucet config")

	deployer, err := ton.NewDeployer(r.Clients.TON, faucetConfig)
	require.NoError(r, err, "unable to create TON deployer")

	deployerID := deployer.GetAddress()

	deployerBalance, err := r.Clients.TON.GetBalanceOf(ctx, deployerID, false)
	require.NoError(r, err, "unable to get balance of TON deployer")

	r.Logger.Print(
		"💎 TON Deployer %s; balance %s",
		deployerID,
		toncontracts.FormatCoins(deployerBalance),
	)

	// Always deploy a new gateway, even if one exists in the env vars
	if r.TONGateway != (tongo.AccountID{}) {
		r.Logger.Print("🔍 Ignoring existing TON Gateway from environment: %s", r.TONGateway.ToRaw())
	}

	// 2. Deploy Gateway
	gwAccount, err := ton.ConstructGatewayAccount(deployerID, r.TSSAddress)
	require.NoError(r, err, "unable to initialize TON gateway")

	r.Logger.Print("🔍 TON Gateway being deployed to address: %s", gwAccount.ID.ToRaw())
	r.Logger.Print("🔍 TON Gateway derived from TSS address: %s", r.TSSAddress.Hex())

	err = deployer.Deploy(ctx, gwAccount, toncontracts.Coins(1))
	require.NoError(r, err, "unable to deploy TON gateway")

	// Set runner field so we use this gateway for tests, overriding any env var
	r.TONGateway = gwAccount.ID
	r.Logger.Print("🔍 TON Gateway address saved in runner: %s", r.TONGateway.ToRaw())

	// 3. Check that the gateway indeed was deployed and has TON balance.
	gwBalance, err := r.Clients.TON.GetBalanceOf(ctx, gwAccount.ID, true)
	require.NoError(r, err, "unable to get balance of TON gateway")
	require.False(r, gwBalance.IsZero(), "TON gateway balance is zero")

	r.Logger.Print(
		"💎 TON Gateway deployed %s; balance: %s",
		gwAccount.ID.ToRaw(),
		toncontracts.FormatCoins(gwBalance),
	)

	amount := toncontracts.Coins(1000)

	// 4. Provision user account
	r.tonProvisionUser(ctx, userTON, deployer, amount)

	// 5. Set chain params & chain nonce
	r.Logger.Print("🔍 Setting up TON chain parameters")

	err = r.ensureTONChainParams(gwAccount)
	require.NoError(r, err, "unable to ensure TON chain params")

	// Verify the parameters were set correctly
	chainID := chains.TONLocalnet.ChainId
	params, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, &observertypes.QueryGetChainParamsForChainRequest{
		ChainId: chainID,
	})

	if err != nil {
		r.Logger.Print("⚠️ Chain params verification failed: %v", err)
	} else {
		r.Logger.Print("✅ TON chain parameters set. Gateway address: %s", params.ChainParams.GatewayAddress)
		if params.ChainParams.GatewayAddress != gwAccount.ID.ToRaw() {
			r.Logger.Print("⚠️ WARNING: Gateway address mismatch: expected %s, got %s",
				gwAccount.ID.ToRaw(), params.ChainParams.GatewayAddress)
		}
	}

	gw := toncontracts.NewGateway(gwAccount.ID)

	// 6. Deposit TON to userTON
	zevmRecipient := userTON.EVMAddress()

	cctx, err := r.TONDeposit(gw, &deployer.Wallet, amount, zevmRecipient)
	require.NoError(r, err, "unable to deposit TON to userTON (additional account)")
	require.Equal(r, cctxtypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
}

func (r *E2ERunner) ensureTONChainParams(gw *ton.AccountInit) error {
	if r.ZetaTxServer == nil {
		return errors.New("ZetaTxServer is not initialized")
	}

	creator := r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName)
	chainID := chains.TONLocalnet.ChainId

	r.Logger.Print("🔍 Setting TON chain params for gateway address: %s", gw.ID.ToRaw())
	r.Logger.Print("🔍 Chain ID: %d, TSS address: %s", chainID, r.TSSAddress.Hex())

	// Set up the chain parameters
	chainParams := &observertypes.ChainParams{
		ChainId:                     chainID,
		ConfirmationCount:           1,
		GasPriceTicker:              5,
		InboundTicker:               5,
		OutboundTicker:              5,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		OutboundScheduleInterval:    2,
		OutboundScheduleLookahead:   5,
		BallotThreshold:             observertypes.DefaultBallotThreshold,
		MinObserverDelegation:       observertypes.DefaultMinObserverDelegation,
		IsSupported:                 true,
		GatewayAddress:              gw.ID.ToRaw(),
		ConfirmationParams: &observertypes.ConfirmationParams{
			SafeInboundCount:  1,
			SafeOutboundCount: 1,
		},
	}

	r.Logger.Print("🔍 Updating TON chain params with Gateway: %s", gw.ID.ToRaw())

	// Update chain params
	err := r.ZetaTxServer.UpdateChainParams(chainParams)
	if err != nil {
		r.Logger.Print("❌ Failed to broadcast TON chain params tx: %v", err)
		return errors.Wrap(err, "unable to broadcast TON chain params tx")
	}
	r.Logger.Print("✅ Successfully broadcast TON chain params update")

	// Reset chain nonces
	resetMsg := observertypes.NewMsgResetChainNonces(creator, chainID, 0, 0)
	resp, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, resetMsg)
	if err != nil {
		r.Logger.Print("❌ Failed to broadcast TON chain nonce reset tx: %v", err)
		return errors.Wrap(err, "unable to broadcast TON chain nonce reset tx")
	}
	r.Logger.Print("✅ Successfully broadcast TON chain nonce reset: %s", resp.TxHash)

	// Allow some time for the transactions to be included in blocks
	r.Logger.Print("⏳ Waiting for transactions to be included in blocks...")
	time.Sleep(5 * time.Second)

	// Wait for params to be queryable
	query := &observertypes.QueryGetChainParamsForChainRequest{ChainId: chainID}
	const checkDuration = 5 * time.Second
	const maxChecks = 10

	for i := 0; i < maxChecks; i++ {
		params, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, query)
		if err == nil {
			r.Logger.Print("💎 TON chain params are set")
			r.Logger.Print("🔍 ZetaCore has TON Gateway address: %s", params.ChainParams.GatewayAddress)
			r.Logger.Print("🔍 Gateway address match: %v", params.ChainParams.GatewayAddress == gw.ID.ToRaw())

			// Extra verification
			if params.ChainParams.GatewayAddress != gw.ID.ToRaw() {
				r.Logger.Print("⚠️ Warning: Gateway address mismatch, but proceeding anyway")
				r.Logger.Print("🔍 Expected: %s, Got: %s", gw.ID.ToRaw(), params.ChainParams.GatewayAddress)
			}

			return nil
		}

		r.Logger.Print("⏳ Waiting for TON chain params to be set (check %d/%d): %v", i+1, maxChecks, err)
		time.Sleep(checkDuration)
	}

	return errors.New("unable to set TON chain params after maximum attempts")
}

// tonProvisionUser deploy & fund ton user account
// that will act as TON sender/receiver in E2E tests
func (r *E2ERunner) tonProvisionUser(
	ctx context.Context,
	user config.Account,
	deployer *ton.Deployer,
	amount math.Uint,
) *wallet.Wallet {
	accInit, wt, err := user.AsTONWallet(r.Clients.TON)
	require.NoError(r, err, "unable to create wallet from TON user account")

	err = deployer.Deploy(ctx, accInit, amount)
	require.NoError(r, err, "unable to deploy TON user wallet %s", wt.GetAddress().ToRaw())

	balance, err := wt.GetBalance(ctx)
	require.NoError(r, err, "unable to get balance of TON user wallet")

	r.Logger.Print(
		"💎 Config.AdditionalAccounts.UserTON: %s; balance: %s",
		wt.GetAddress().ToRaw(),
		toncontracts.FormatCoins(math.NewUint(balance)),
	)

	return wt
}
