package runner

import (
	"context"
	"time"

	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	tontypes "github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"

	"github.com/zeta-chain/node/e2e/config"
	ton "github.com/zeta-chain/node/e2e/runner/ton"
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

	// 2. Deploy Gateway
	gwAccount, err := ton.ConstructGatewayAccount(deployerID, r.TSSAddress)
	require.NoError(r, err, "unable to initialize TON gateway")

	err = deployer.Deploy(ctx, gwAccount, toncontracts.Coins(1))
	require.NoError(r, err, "unable to deploy TON gateway")

	// 3. Check that the gateway indeed was deployed and has desired TON balance.
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
	err = r.ensureTONChainParams(gwAccount)
	require.NoError(r, err, "unable to ensure TON chain params")

	gw := toncontracts.NewGateway(gwAccount.ID)

	// 5. Deposit TON to userTON
	zevmRecipient := userTON.EVMAddress()

	cctx, err := r.TONDeposit(gw, &deployer.Wallet, amount, zevmRecipient)
	require.NoError(r, err, "unable to deposit TON to userTON (additional account)")
	require.Equal(r, cctxtypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// Set runner field
	r.TONGateway = gw.AccountID()
}

func (r *E2ERunner) ensureTONChainParams(gw *ton.AccountInit) error {
	if r.ZetaTxServer == nil {
		return errors.New("ZetaTxServer is not initialized")
	}

	creator := r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName)

	chainID := chains.TONLocalnet.ChainId

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
		GasPriceMultiplier: observertypes.DefaultGasPriceMultiplier,
	}

	if err := r.ZetaTxServer.UpdateChainParams(chainParams); err != nil {
		return errors.Wrap(err, "unable to broadcast TON chain params tx")
	}

	resetMsg := observertypes.NewMsgResetChainNonces(creator, chainID, 0, 0)
	if _, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, resetMsg); err != nil {
		return errors.Wrap(err, "unable to broadcast TON chain nonce reset tx")
	}

	r.Logger.Print("💎 Voted for adding TON chain params (localnet). Waiting for confirmation")

	query := &observertypes.QueryGetChainParamsForChainRequest{ChainId: chainID}

	const duration = 2 * time.Second

	for i := 0; i < 10; i++ {
		_, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, query)
		if err == nil {
			r.Logger.Print("💎 TON chain params are set")
			return nil
		}

		time.Sleep(duration)
	}

	return errors.New("unable to set TON chain params")
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

// UpdateTSSAddressTON updates the TSS address on the TON gateway contract and resets the seqno.
func (r *E2ERunner) UpdateTSSAddressTON(gatewayAccountID, faucetURL string) {
	require.NotNil(r, r.Clients.TON, "TON client is not initialized")
	require.NotEmpty(r, faucetURL, "TON faucet url is empty")

	ctx := r.Ctx
	faucetConfig, err := ton.GetFaucet(ctx, faucetURL)
	require.NoError(r, err, "unable to get faucet config")
	deployer, err := ton.NewDeployer(r.Clients.TON, faucetConfig)
	require.NoError(r, err, "unable to create TON deployer")

	var gwIDStr string
	if gatewayAccountID != "" {
		gwIDStr = gatewayAccountID
	} else if !r.TONGateway.IsZero() {
		gwIDStr = r.TONGateway.ToRaw()
	}
	require.NotEmpty(r, gwIDStr, "TON gateway account id is empty")
	gw := toncontracts.NewGateway(tontypes.MustParseAccountID(gwIDStr))

	deployerSeqno, err := r.Clients.TON.GetSeqno(ctx, deployer.GetAddress())
	require.NoError(r, err, "unable to get deployer seqno")

	// Send update_tss admin message
	// Use 0.1 TON for gas (authority operations require minimal gas)
	amount := math.NewUint(100_000_000)
	sendMode := toncontracts.SendFlagSeparateFees + toncontracts.SendFlagIgnoreErrors

	err = gw.SendUpdateTSS(ctx, &deployer.Wallet, amount, r.TSSAddress, sendMode)
	require.NoError(r, err, "unable to send update_tss to TON gateway")

	const maxWaitingTime = 30 * time.Second
	err = r.Clients.TON.WaitForNextSeqno(ctx, deployer.GetAddress(), deployerSeqno, maxWaitingTime)
	require.NoError(r, err, "unable to confirm update_tss transaction")
	err = r.Clients.TON.WaitForBlocks(ctx)
	require.NoError(r, err, "unable to wait for blocks after update_tss")

	r.Logger.Print("💎 successfully updated TSS address on TON gateway to %s", r.TSSAddress.Hex())
	deployerSeqno, err = r.Clients.TON.GetSeqno(ctx, deployer.GetAddress())
	require.NoError(r, err, "unable to get deployer seqno before reset_seqno")

	err = gw.SendResetSeqno(ctx, &deployer.Wallet, amount, 0, sendMode)
	require.NoError(r, err, "unable to send reset_seqno to TON gateway")

	err = r.Clients.TON.WaitForNextSeqno(ctx, deployer.GetAddress(), deployerSeqno, maxWaitingTime)
	require.NoError(r, err, "unable to confirm reset_seqno transaction")

	err = r.Clients.TON.WaitForBlocks(ctx)
	require.NoError(r, err, "unable to wait for blocks after reset_seqno")
}
