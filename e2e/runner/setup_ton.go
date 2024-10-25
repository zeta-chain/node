package runner

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/e2e/runner/ton"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// SetupTON setups TON deployer and deploys Gateway contract
func (r *E2ERunner) SetupTON() error {
	if r.Clients.TON == nil {
		return fmt.Errorf("TON clients are not initialized")
	}

	ctx := r.Ctx

	// 1. Setup Deployer (acts as a faucet as well)
	faucetConfig, err := r.Clients.TON.GetFaucet(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get faucet config")
	}

	deployer, err := ton.NewDeployer(r.Clients.TON, faucetConfig)
	if err != nil {
		return errors.Wrap(err, "unable to create TON deployer")
	}

	depAddr := deployer.GetAddress()
	r.Logger.Print("💎TON Deployer %s (%s)", depAddr.ToRaw(), depAddr.ToHuman(false, true))

	// 2. Deploy Gateway
	gwAccount, err := ton.ConstructGatewayAccount(depAddr, r.TSSAddress)
	if err != nil {
		return errors.Wrap(err, "unable to initialize TON gateway")
	}

	if err = deployer.Deploy(ctx, gwAccount, toncontracts.Coins(1)); err != nil {
		return errors.Wrapf(err, "unable to deploy TON gateway")
	}

	r.Logger.Print(
		"💎TON Gateway deployed %s (%s) with TSS address %s",
		gwAccount.ID.ToRaw(),
		gwAccount.ID.ToHuman(false, true),
		r.TSSAddress.Hex(),
	)

	// 3. Check that the gateway indeed was deployed and has desired TON balance.
	gwBalance, err := deployer.GetBalanceOf(ctx, gwAccount.ID)
	if err != nil {
		return errors.Wrap(err, "unable to get balance of TON gateway")
	}

	if gwBalance.IsZero() {
		return fmt.Errorf("TON gateway balance is zero")
	}

	r.TONDeployer = deployer
	r.TONGateway = toncontracts.NewGateway(gwAccount.ID)

	return r.ensureTONChainParams(gwAccount)
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
	}

	msg := observertypes.NewMsgUpdateChainParams(creator, chainParams)

	if _, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msg); err != nil {
		return errors.Wrap(err, "unable to broadcast TON chain params tx")
	}

	resetMsg := observertypes.NewMsgResetChainNonces(creator, chainID, 0, 0)
	if _, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, resetMsg); err != nil {
		return errors.Wrap(err, "unable to broadcast TON chain nonce reset tx")
	}

	r.Logger.Print("💎Voted for adding TON chain params (localnet). Waiting for confirmation")

	query := &observertypes.QueryGetChainParamsForChainRequest{ChainId: chainID}

	const duration = 2 * time.Second

	for i := 0; i < 10; i++ {
		_, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, query)
		if err == nil {
			r.Logger.Print("💎TON chain params are set")
			return nil
		}

		time.Sleep(duration)
	}

	return errors.New("unable to set TON chain params")
}
