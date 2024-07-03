package orchestrator

import (
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	btcrpc "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/rpc"
	btcsigner "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/signer"
	evmobserver "github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	evmsigner "github.com/zeta-chain/zetacore/zetaclient/chains/evm/signer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

// WatchEnabledChains watches for run-time chain activation and deactivation
func (oc *Orchestrator) WatchEnabledChains() {
	oc.logger.Std.Info().Msg("WatchChainActivation started")

	ticker := time.NewTicker(common.ZetaBlockTime * 2)
	for {
		select {
		case <-ticker.C:
			oc.ActivateAndDeactivateChains()
		case <-oc.stop:
			oc.logger.Std.Info().Msg("WatchChainActivation stopped")
			return
		}
	}
}

// ActivateAndDeactivateChains activates and deactivates chains according to chain params and config file
//
// The chains to be activated:
//   - chain params flag 'IsSupported' is true AND
//   - chain is configured in config file
//
// The chains to be deactivated:
//   - chain params flag 'IsSupported' is false OR
//   - chain is not configured in config file
//
// Note:
//   - zetaclient will reload config file periodically and update in-memory config accordingly.
//   - As an tss signer, please make sure the config file is always well configured and not missing any chain
func (oc *Orchestrator) ActivateAndDeactivateChains() {
	// create new signer and observer maps
	// Note: the keys of the two maps are chain IDs and they are always exactly matched
	newSignerMap := make(map[int64]interfaces.ChainSigner)
	newObserverMap := make(map[int64]interfaces.ChainObserver)

	// create new signers and observers
	oc.CreateObserversEVM(newSignerMap, newObserverMap)
	oc.CreateObserversBTC(newSignerMap, newObserverMap)

	// activate newly supported chains and deactivate chains that are no longer supported
	oc.DeactivateChains(newObserverMap)
	oc.ActivateChains(newSignerMap, newObserverMap)
}

// DeactivateChains deactivates chains that are no longer supported
func (oc *Orchestrator) DeactivateChains(
	newObserverMap map[int64]interfaces.ChainObserver,
) {
	// loop through existing observer map to deactivate chains that are not in new observer map
	oc.mu.Lock()
	defer oc.mu.Unlock()
	for chainID, observer := range oc.observerMap {
		_, found := newObserverMap[chainID]
		if !found {
			oc.logger.Std.Info().Msgf("DeactivateChains: deactivating chain %d", chainID)
			observer.Stop()

			// remove signer and observer from maps
			delete(oc.signerMap, chainID)
			delete(oc.observerMap, chainID)
			oc.logger.Std.Info().Msgf("DeactivateChains: deactivated chain %d", chainID)
		}
	}
}

// ActivateChains activates newly supported chains
func (oc *Orchestrator) ActivateChains(
	newSignerMap map[int64]interfaces.ChainSigner,
	newObserverMap map[int64]interfaces.ChainObserver,
) {
	// loop through new observer map to activate chains that are not in existing observer map
	for chainID, observer := range newObserverMap {
		_, found := oc.observerMap[chainID]
		if !found {
			oc.logger.Std.Info().Msgf("ActivateChains: activating chain %d", chainID)

			// open database and load data
			err := observer.LoadDB(oc.dbPath)
			if err != nil {
				oc.logger.Std.Error().
					Err(err).
					Msgf("ActivateChains: error LoadDB for chain %d", chainID)
				continue
			}
			observer.Start()

			// add signer and observer to maps
			oc.mu.Lock()
			oc.signerMap[chainID] = newSignerMap[chainID]
			oc.observerMap[chainID] = observer
			oc.mu.Unlock()

			oc.logger.Std.Info().Msgf("ActivateChains: activated chain %d", chainID)
		}
	}
}

// CreateObserversEVM creates signer and observer maps for all enabled EVM chains
func (oc *Orchestrator) CreateObserversEVM(
	resultSignerMap map[int64]interfaces.ChainSigner,
	resultObserverMap map[int64]interfaces.ChainObserver,
) {
	// create EVM-chain signers
	for _, evmConfig := range oc.appContext.Config().GetAllEVMConfigs() {
		chainParams, found := oc.appContext.GetExternalChainParams(evmConfig.Chain.ChainId)
		if !found {
			oc.logger.Sampled.Warn().
				Msgf("CreateObserversEVM: chain parameter not found for chain %d", evmConfig.Chain.ChainId)
			continue
		}
		connectorAddress := ethcommon.HexToAddress(chainParams.ConnectorContractAddress)
		erc20CustodyAddress := ethcommon.HexToAddress(chainParams.Erc20CustodyContractAddress)

		// create RPC client
		evmClient, err := ethclient.Dial(evmConfig.Endpoint)
		if err != nil {
			oc.logger.Std.Error().
				Err(err).
				Msgf("CreateObserversEVM: error dailing endpoint %s for chain %d", evmConfig.Endpoint, evmConfig.Chain.ChainId)
			continue
		}

		// create signer
		signer, err := evmsigner.NewSigner(
			evmConfig.Chain,
			oc.appContext,
			oc.tss,
			oc.ts,
			oc.logger.Base,
			evmConfig.Endpoint,
			config.GetConnectorABI(),
			config.GetERC20CustodyABI(),
			connectorAddress,
			erc20CustodyAddress)
		if err != nil {
			oc.logger.Std.Error().
				Err(err).
				Msgf("CreateObserversEVM: error NewSigner for chain %d", evmConfig.Chain.ChainId)
			continue
		}

		// create observer
		observer, err := evmobserver.NewObserver(
			evmConfig,
			evmClient,
			*chainParams,
			oc.appContext,
			oc.zetacoreClient,
			oc.tss,
			oc.logger.Base,
			oc.ts,
		)
		if err != nil {
			oc.logger.Std.Error().
				Err(err).
				Msgf("CreateObserversEVM: error NewObserver for chain %d", evmConfig.Chain.ChainId)
			continue
		}

		// add signer and observer to result maps
		resultSignerMap[evmConfig.Chain.ChainId] = signer
		resultObserverMap[evmConfig.Chain.ChainId] = observer
	}
}

// CreateObserversBTC creates signer and observer maps for all enabled BTC chains
func (oc *Orchestrator) CreateObserversBTC(
	resultSignerMap map[int64]interfaces.ChainSigner,
	resultObserverMap map[int64]interfaces.ChainObserver,
) {
	// get enabled BTC chains and config
	btcChains := oc.appContext.GetEnabledBTCChains()
	btcConfig, found := oc.appContext.Config().GetBTCConfig()

	// currently only one single BTC chain is supported
	if !found {
		oc.logger.Sampled.Warn().Msg("CreateObserversBTC: BTC config not found")
		return
	}
	if len(btcChains) != 1 {
		oc.logger.Std.Error().Msgf("CreateObserversBTC: want single BTC chain, got %d", len(btcChains))
		return
	}

	// create BTC-chain signers and observers
	// loop is used here in case we have multiple btc chains in the future
	for _, btcChain := range btcChains {
		chainParams, found := oc.appContext.GetExternalChainParams(btcChain.ChainId)
		if !found {
			oc.logger.Sampled.Warn().
				Msgf("CreateObserversBTC: chain parameter not found for chain %d", btcChain.ChainId)
			continue
		}

		// create RPC client
		btcClient, err := btcrpc.NewRPCClient(btcConfig)
		if err != nil {
			oc.logger.Std.Error().
				Err(err).
				Msgf("CreateObserversBTC: error NewRPCClient for chain %s", btcChain.String())
			continue
		}

		// create signer
		signer, err := btcsigner.NewSigner(btcChain, oc.appContext, oc.tss, oc.ts, oc.logger.Base, btcConfig)
		if err != nil {
			oc.logger.Std.Error().Err(err).Msgf("CreateObserversBTC: error NewSigner for chain %d", btcChain.ChainId)
			continue
		}

		// create observer
		observer, err := btcobserver.NewObserver(
			btcChain,
			btcClient,
			*chainParams,
			oc.appContext,
			oc.zetacoreClient,
			oc.tss,
			oc.logger.Base,
			oc.ts,
		)
		if err != nil {
			oc.logger.Std.Error().Err(err).Msgf("NewObserver error for bitcoin chain %s", btcChain.String())
		}

		// add signer and observer to result maps
		resultSignerMap[btcChain.ChainId] = signer
		resultObserverMap[btcChain.ChainId] = observer
	}
}
