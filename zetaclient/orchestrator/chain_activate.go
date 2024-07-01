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

	ticker := time.NewTicker(common.ZetaBlockTime)
	for {
		select {
		case <-ticker.C:
			oc.ActivateDeactivateChains()
		case <-oc.stop:
			oc.logger.Std.Info().Msg("WatchChainActivation stopped")
			return
		}
	}
}

// ActivateDeactivateChains activates or deactivates chain observers and signers
func (oc *Orchestrator) ActivateDeactivateChains() {
	// create new signer and observer maps
	// Note: the keys of the two maps are chain IDs and they are always exactly matched
	newSignerMap := make(map[int64]interfaces.ChainSigner)
	newObserverMap := make(map[int64]interfaces.ChainObserver)

	// create new signers and observers
	oc.CreateObserversEVM(newSignerMap, newObserverMap)
	oc.CreateObserversBTC(newSignerMap, newObserverMap)

	// loop through existing observer map to deactivate chains that are not in new observer map
	for chainID, observer := range oc.observerMap {
		_, found := newObserverMap[chainID]
		if !found {
			oc.logger.Std.Info().Msgf("orchestrator deactivating chain %d", chainID)

			observer.Stop()
			delete(oc.signerMap, chainID)
			delete(oc.observerMap, chainID)
		}
	}

	// loop through new observer map to activate chains that are not in existing observer map
	for chainID, observer := range newObserverMap {
		_, found := oc.observerMap[chainID]
		if !found {
			oc.logger.Std.Info().Msgf("orchestrator activating chain %d", chainID)

			observer.Start()
			oc.signerMap[chainID] = newSignerMap[chainID]
			oc.observerMap[chainID] = observer
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
			oc.dbPath,
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
			oc.dbPath,
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
