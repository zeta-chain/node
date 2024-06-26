package orchestrator

import (
	"fmt"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	btcrpc "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/rpc"
	btcsigner "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/signer"
	evmobserver "github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	evmsigner "github.com/zeta-chain/zetacore/zetaclient/chains/evm/signer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// WatchUpgradePlan watches for upgrade plan and stops orchestrator if upgrade height is reached
func (oc *Orchestrator) WatchUpgradePlan() {
	oc.logger.Std.Info().Msg("WatchUpgradePlan started")

	// detect upgrade plan every half Zeta block
	ticker := time.NewTicker(common.ZetaBlockTime / 2)
	for range ticker.C {
		reached, err := oc.UpgradeHeightReached()
		if err != nil {
			oc.logger.Sampled.Error().Err(err).Msg("error detecting upgrade plan")
		} else if reached {
			oc.Stop()
			oc.logger.Std.Info().Msg("WatchUpgradePlan stopped")

			return
		}
	}
}

// UpdateAppContext is a polling goroutine that checks and updates app context periodically
func (oc *Orchestrator) UpdateAppContext() {
	oc.logger.Std.Info().Msg("UpdateAppContext started")

	ticker := time.NewTicker(time.Duration(oc.appContext.Config().ConfigUpdateTicker) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := UpdateZetacoreContext(oc.zetacoreClient, oc.appContext.ZetacoreContext(), false, oc.logger.Std)
			if err != nil {
				oc.logger.Std.Err(err).Msg("error updating zetaclient app context")
			}
		case <-oc.stop:
			oc.logger.Std.Info().Msg("UpdateAppContext stopped")
			return
		}
	}
}

// CreateAppContext creates new app context from config and zetacore client
func CreateAppContext(
	cfg config.Config,
	zetacoreClient interfaces.ZetacoreClient,
	logger zerolog.Logger,
) (*context.AppContext, error) {
	// create app context from config
	appContext := context.NewAppContext(context.NewZetacoreContext(cfg), cfg)

	// update zetacore context from zetacore
	err := UpdateZetacoreContext(zetacoreClient, appContext.ZetacoreContext(), true, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error updating zetacore context")
	}

	return appContext, nil
}

// UpdateZetacoreContext updates zetacore context
// zetacore stores zetacore context for all clients
func UpdateZetacoreContext(
	zetacoreClient interfaces.ZetacoreClient,
	coreContext *context.ZetacoreContext,
	init bool,
	logger zerolog.Logger,
) error {
	// create latest zetacore context from zetacore
	zetacoreContext, err := zetacoreClient.GetLatestZetacoreContext()
	if err != nil {
		return errors.Wrap(err, "error getting latest zetacore context")
	}

	// get keygen
	keygen := zetacoreContext.GetKeygen()

	// get btc chain params
	_, newBTCParams, _ := zetacoreContext.GetBTCChainParams()

	coreContext.Update(
		&keygen,
		zetacoreContext.GetEnabledExternalChains(),
		zetacoreContext.GetAllEVMChainParams(),
		newBTCParams,
		zetacoreContext.GetCurrentTssPubkey(),
		zetacoreContext.GetCrossChainFlags(),
		zetacoreContext.GetAllHeaderEnabledChains(),
		init,
		logger,
	)

	return nil
}

// UpgradeHeightReached returns true if upgrade height is reached
func (oc *Orchestrator) UpgradeHeightReached() (bool, error) {
	// query for active upgrade plan
	plan, err := oc.zetacoreClient.GetUpgradePlan()
	if err != nil {
		return false, fmt.Errorf("failed to get upgrade plan: %w", err)
	}

	// if there is no active upgrade plan, plan will be nil.
	if plan == nil {
		return false, nil
	}

	// get ZetaChain block height
	height, err := oc.zetacoreClient.GetBlockHeight()
	if err != nil {
		return false, fmt.Errorf("failed to get block height: %w", err)
	}

	// if upgrade height is not reached, do nothing
	if height != plan.Height-1 {
		return false, nil
	}

	// stop zetaclients if upgrade height is reached; notify operator to upgrade and restart
	oc.logger.Std.Warn().
		Msgf("Active upgrade plan detected and upgrade height reached: %s at height %d; ZetaClient is stopped;"+
			"please kill this process, replace zetaclientd binary with upgraded version, and restart zetaclientd", plan.Name, plan.Height)

	return true, nil
}

// GetLatestZetacoreContext queries zetacore to build the latest zetacore context
func (oc *Orchestrator) GetLatestZetacoreContext(client interfaces.ZetacoreClient) (*context.ZetacoreContext, error) {
	// get latest supported chains
	supportedChains, err := client.GetSupportedChains()
	if err != nil {
		return nil, errors.Wrap(err, "GetSupportedChains failed")
	}
	supportedChainsMap := make(map[int64]chains.Chain)
	for _, chain := range supportedChains {
		supportedChainsMap[chain.ChainId] = *chain
	}

	// get latest chain parameters
	chainParams, err := client.GetChainParams()
	if err != nil {
		return nil, errors.Wrap(err, "GetChainParams failed")
	}

	chainsEnabled := make([]chains.Chain, 0)
	chainParamMap := make(map[int64]*observertypes.ChainParams)

	newEVMParams := make(map[int64]*observertypes.ChainParams)
	var newBTCParams *observertypes.ChainParams

	for _, chainParam := range chainParams {
		// skip unsupported chain
		if !chainParam.IsSupported {
			continue
		}

		// chain should exist in chain list
		chain, found := supportedChainsMap[chainParam.ChainId]
		if !found {
			continue
		}

		// skip ZetaChain
		if !chain.IsExternalChain() {
			continue
		}

		// add chain param to map
		chainParamMap[chainParam.ChainId] = chainParam

		// keep this chain
		chainsEnabled = append(chainsEnabled, chain)
		if chains.IsBitcoinChain(chainParam.ChainId) {
			newBTCParams = chainParam
		} else if chains.IsEVMChain(chainParam.ChainId) {
			newEVMParams[chainParam.ChainId] = chainParam
		}
	}

	// get latest keygen
	keyGen, err := client.GetKeyGen()
	if err != nil {
		return nil, errors.Wrap(err, "GetKeyGen failed")
	}

	// get latest TSS public key
	tss, err := client.GetCurrentTss()
	if err != nil {
		return nil, errors.Wrap(err, "GetCurrentTss failed")
	}
	tssPubKey := tss.GetTssPubkey()

	// get latest crosschain flags
	crosschainFlags, err := client.GetCrosschainFlags()
	if err != nil {
		return nil, errors.Wrap(err, "GetCrosschainFlags failed")
	}

	// get latest block header enabled chains
	blockHeaderEnabledChains, err := client.GetBlockHeaderEnabledChains()
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockHeaderEnabledChains failed")
	}

	return context.CreateZetacoreContext(
		keyGen,
		chainsEnabled,
		chainParamMap,
		newEVMParams,
		newBTCParams,
		tssPubKey,
		crosschainFlags,
		blockHeaderEnabledChains,
	), nil
}

// CreateSignerMap creates a map of ChainSigners for all chains in the config
func CreateSignerMap(
	appContext *context.AppContext,
	tss interfaces.TSSSigner,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (map[int64]interfaces.ChainSigner, error) {
	zetacoreContext := appContext.ZetacoreContext()
	signerMap := make(map[int64]interfaces.ChainSigner)

	// EVM signers
	for _, evmConfig := range appContext.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		evmChainParams, found := zetacoreContext.GetEVMChainParams(evmConfig.Chain.ChainId)
		if !found {
			logger.Std.Error().Msgf("ChainParam not found for chain %s", evmConfig.Chain.String())
			continue
		}
		mpiAddress := ethcommon.HexToAddress(evmChainParams.ConnectorContractAddress)
		erc20CustodyAddress := ethcommon.HexToAddress(evmChainParams.Erc20CustodyContractAddress)
		signer, err := evmsigner.NewSigner(
			evmConfig.Chain,
			zetacoreContext,
			tss,
			ts,
			logger,
			evmConfig.Endpoint,
			config.GetConnectorABI(),
			config.GetERC20CustodyABI(),
			mpiAddress,
			erc20CustodyAddress)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewEVMSigner error for chain %s", evmConfig.Chain.String())
			continue
		}
		signerMap[evmConfig.Chain.ChainId] = signer
	}
	// BTC signer
	btcChain, btcConfig, enabled := appContext.GetBTCChainAndConfig()
	if enabled {
		signer, err := btcsigner.NewSigner(btcChain, zetacoreContext, tss, ts, logger, btcConfig)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewBTCSigner error for chain %s", btcChain.String())
		} else {
			signerMap[btcChain.ChainId] = signer
		}
	}

	return signerMap, nil
}

// CreateChainObserverMap creates a map of ChainObservers for all chains in the config
func CreateChainObserverMap(
	appContext *context.AppContext,
	zetacoreClient *zetacore.Client,
	tss interfaces.TSSSigner,
	dbpath string,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (map[int64]interfaces.ChainObserver, error) {
	zetacoreContext := appContext.ZetacoreContext()
	observerMap := make(map[int64]interfaces.ChainObserver)
	// EVM observers
	for _, evmConfig := range appContext.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		chainParams, found := zetacoreContext.GetEVMChainParams(evmConfig.Chain.ChainId)
		if !found {
			logger.Std.Error().Msgf("ChainParam not found for chain %s", evmConfig.Chain.String())
			continue
		}

		// create EVM client
		evmClient, err := ethclient.Dial(evmConfig.Endpoint)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("error dailing endpoint %s", evmConfig.Endpoint)
			continue
		}

		// create EVM chain observer
		observer, err := evmobserver.NewObserver(
			evmConfig,
			evmClient,
			*chainParams,
			zetacoreContext,
			zetacoreClient,
			tss,
			dbpath,
			logger,
			ts,
		)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewObserver error for evm chain %s", evmConfig.Chain.String())
			continue
		}
		observerMap[evmConfig.Chain.ChainId] = observer
	}

	// BTC observer
	_, chainParams, found := zetacoreContext.GetBTCChainParams()
	if !found {
		return nil, fmt.Errorf("bitcoin chains params not found")
	}

	// create BTC chain observer
	btcChain, btcConfig, enabled := appContext.GetBTCChainAndConfig()
	if enabled {
		btcClient, err := btcrpc.NewRPCClient(btcConfig)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("error creating rpc client for bitcoin chain %s", btcChain.String())
		} else {
			// create BTC chain observer
			observer, err := btcobserver.NewObserver(
				btcChain,
				btcClient,
				*chainParams,
				zetacoreContext,
				zetacoreClient,
				tss,
				dbpath,
				logger,
				ts,
			)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("NewObserver error for bitcoin chain %s", btcChain.String())
			} else {
				observerMap[btcChain.ChainId] = observer
			}
		}
	}

	return observerMap, nil
}
