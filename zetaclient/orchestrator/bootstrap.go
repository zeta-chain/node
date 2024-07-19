package orchestrator

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/rpc"
	btcsigner "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/signer"
	evmobserver "github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	evmsigner "github.com/zeta-chain/zetacore/zetaclient/chains/evm/signer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	solbserver "github.com/zeta-chain/zetacore/zetaclient/chains/solana/observer"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/db"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

// backwards compatibility
const btcDatabaseFilename = "btc_chain_client"

// CreateSignerMap creates a map of interfaces.ChainSigner (by chainID) for all chains in the config.
// Note that signer construction failure for a chain does not prevent the creation of signers for other chains.
func CreateSignerMap(
	ctx context.Context,
	tss interfaces.TSSSigner,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (map[int64]interfaces.ChainSigner, error) {
	signers := make(map[int64]interfaces.ChainSigner)
	_, _, err := syncSignerMap(ctx, tss, logger, ts, &signers)
	if err != nil {
		return nil, err
	}

	return signers, nil
}

// syncSignerMap synchronizes the given signers map with the signers for all chains in the config.
// This semantic is used to allow dynamic updates to the signers map.
func syncSignerMap(
	ctx context.Context,
	tss interfaces.TSSSigner,
	logger base.Logger,
	ts *metrics.TelemetryServer,
	signers *map[int64]interfaces.ChainSigner,
) (int, int, error) {
	if signers == nil {
		return 0, 0, errors.New("signers map is nil")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to get app context")
	}

	var (
		added, removed int

		presentChainIDs = make([]int64, 0)

		onAfterSet = func(chainID int64, _ interfaces.ChainSigner) {
			logger.Std.Info().Msgf("Added signer for chain %d", chainID)
			added++
		}

		onBeforeUnset = func(chainID int64, _ interfaces.ChainSigner) {
			logger.Std.Info().Msgf("Removing signer for chain %d", chainID)
			removed++
		}
	)

	// EVM signers
	for _, evmConfig := range app.Config().GetAllEVMConfigs() {
		chainID := evmConfig.Chain.ChainId

		if evmConfig.Chain.IsZetaChain() {
			continue
		}

		evmChainParams, found := app.GetEVMChainParams(chainID)
		switch {
		case !found:
			logger.Std.Warn().Msgf("Unable to find chain params for EVM chain %d", chainID)
			continue
		case !evmChainParams.IsSupported:
			logger.Std.Warn().Msgf("EVM chain %d is not supported", chainID)
			continue
		}

		presentChainIDs = append(presentChainIDs, chainID)

		// noop for existing signers
		if mapHas(signers, chainID) {
			continue
		}

		var (
			mpiAddress          = ethcommon.HexToAddress(evmChainParams.ConnectorContractAddress)
			erc20CustodyAddress = ethcommon.HexToAddress(evmChainParams.Erc20CustodyContractAddress)
		)

		signer, err := evmsigner.NewSigner(
			ctx,
			evmConfig.Chain,
			tss,
			ts,
			logger,
			evmConfig.Endpoint,
			config.GetConnectorABI(),
			config.GetERC20CustodyABI(),
			mpiAddress,
			erc20CustodyAddress,
		)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("Unable to construct signer for EVM chain %d", chainID)
			continue
		}

		mapSet[int64, interfaces.ChainSigner](signers, chainID, signer, onAfterSet)
	}

	// BTC signer
	// Emulate same loop semantics as for EVM chains
	for i := 0; i < 1; i++ {
		btcChain, btcChainParams, btcChainParamsFound := app.GetBTCChainParams()
		switch {
		case !btcChainParamsFound:
			logger.Std.Warn().Msgf("Unable to find chain params for BTC chain")
			continue
		case !btcChainParams.IsSupported:
			logger.Std.Warn().Msgf("BTC chain is not supported")
			continue
		}

		chainID := btcChainParams.ChainId

		presentChainIDs = append(presentChainIDs, chainID)

		// noop
		if mapHas(signers, chainID) {
			continue
		}

		cfg, _ := app.Config().GetBTCConfig()

		utxoSigner, err := btcsigner.NewSigner(btcChain, tss, ts, logger, cfg)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("Unable to construct signer for UTXO chain %d", chainID)
			continue
		}

		mapSet[int64, interfaces.ChainSigner](signers, chainID, utxoSigner, onAfterSet)
	}

	// Remove all disabled signers
	mapDeleteMissingKeys(signers, presentChainIDs, onBeforeUnset)

	return added, removed, nil
}

// CreateChainObserverMap creates a map of interfaces.ChainObserver (by chainID) for all chains in the config.
// Note (!) that it calls observer.Start() on creation
func CreateChainObserverMap(
	ctx context.Context,
	client interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	dbpath string,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (map[int64]interfaces.ChainObserver, error) {
	observerMap := make(map[int64]interfaces.ChainObserver)

	_, _, err := syncObserverMap(ctx, client, tss, dbpath, logger, ts, &observerMap)
	if err != nil {
		return nil, err
	}

	return observerMap, nil
}

// syncObserverMap synchronizes the given observer map with the observers for all chains in the config.
// This semantic is used to allow dynamic updates to the map.
// Note (!) that it calls observer.Start() on creation and observer.Stop() on deletion.
func syncObserverMap(
	ctx context.Context,
	client interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	dbpath string,
	logger base.Logger,
	ts *metrics.TelemetryServer,
	observerMap *map[int64]interfaces.ChainObserver,
) (int, int, error) {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to get app context")
	}

	var (
		added, removed int

		presentChainIDs = make([]int64, 0)

		onAfterSet = func(_ int64, ob interfaces.ChainObserver) {
			ob.Start(ctx)
			added++
		}

		onBeforeUnset = func(_ int64, ob interfaces.ChainObserver) {
			ob.Stop()
			removed++
		}
	)

	// EVM observers
	for _, evmConfig := range app.Config().GetAllEVMConfigs() {
		var (
			chainID   = evmConfig.Chain.ChainId
			chainName = evmConfig.Chain.ChainName.String()
		)

		if evmConfig.Chain.IsZetaChain() {
			continue
		}

		chainParams, found := app.GetEVMChainParams(evmConfig.Chain.ChainId)
		switch {
		case !found:
			logger.Std.Error().Msgf("Unable to find chain params for EVM chain %d", chainID)
			continue
		case !chainParams.IsSupported:
			logger.Std.Error().Msgf("EVM chain %d is not supported", chainID)
			continue
		}

		presentChainIDs = append(presentChainIDs, chainID)

		// noop
		if mapHas(observerMap, chainID) {
			continue
		}

		// create EVM client
		evmClient, err := ethclient.DialContext(ctx, evmConfig.Endpoint)
		if err != nil {
			logger.Std.Error().Err(err).Str("rpc.endpoint", evmConfig.Endpoint).Msgf("Unable to dial EVM RPC")
			continue
		}

		database, err := db.NewFromSqlite(dbpath, chainName, true)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("Unable to open a database for EVM chain %q", chainName)
			continue
		}

		// create EVM chain observer
		observer, err := evmobserver.NewObserver(
			ctx,
			evmConfig,
			evmClient,
			*chainParams,
			client,
			tss,
			database,
			logger,
			ts,
		)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewObserver error for EVM chain %s", evmConfig.Chain.String())
			continue
		}

		mapSet[int64, interfaces.ChainObserver](observerMap, chainID, observer, onAfterSet)
	}

	// Emulate same loop semantics as for EVM chains
	// create BTC chain observer
	for i := 0; i < 1; i++ {
		btcChain, btcConfig, btcEnabled := app.GetBTCChainAndConfig()
		if !btcEnabled {
			continue
		}

		chainID := btcChain.ChainId

		_, btcChainParams, found := app.GetBTCChainParams()
		switch {
		case !found:
			logger.Std.Warn().Msgf("Unable to find chain params for BTC chain %d", chainID)
			continue
		case !btcChainParams.IsSupported:
			logger.Std.Warn().Msgf("BTC chain %d is not supported", chainID)
			continue
		}

		presentChainIDs = append(presentChainIDs, chainID)

		// noop
		if mapHas(observerMap, chainID) {
			continue
		}

		btcRPC, err := rpc.NewRPCClient(btcConfig)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("unable to create rpc client for BTC chain %d", chainID)
			continue
		}

		database, err := db.NewFromSqlite(dbpath, btcDatabaseFilename, true)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("unable to open database for BTC chain %d", chainID)
			continue
		}

		btcObserver, err := btcobserver.NewObserver(
			btcChain,
			btcRPC,
			*btcChainParams,
			client,
			tss,
			database,
			logger,
			ts,
		)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewObserver error for BTC chain %d", chainID)
			continue
		}

		mapSet[int64, interfaces.ChainObserver](observerMap, btcChain.ChainId, btcObserver, onAfterSet)
	}

	// Emulate same loop semantics as for EVM chains
	// create SOL chain observer
	for i := 0; i < 1; i++ {
		solChain, solConfig, solEnabled := app.GetSolanaChainAndConfig()
		if !solEnabled {
			continue
		}

		var (
			chainID   = solChain.ChainId
			chainName = solChain.ChainName.String()
		)

		_, solanaChainParams, found := app.GetSolanaChainParams()
		switch {
		case !found:
			logger.Std.Warn().Msgf("Unable to find chain params for SOL chain %d", chainID)
			continue
		case !solanaChainParams.IsSupported:
			logger.Std.Warn().Msgf("SOL chain %d is not supported", chainID)
			continue
		}

		presentChainIDs = append(presentChainIDs, chainID)

		// noop
		if mapHas(observerMap, chainID) {
			continue
		}

		rpcClient := solrpc.New(solConfig.Endpoint)
		if rpcClient == nil {
			// should never happen
			logger.Std.Error().Msg("solana create Solana client error")
			continue
		}

		database, err := db.NewFromSqlite(dbpath, chainName, true)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("unable to open database for SOL chain %d", chainID)
			continue
		}

		solObserver, err := solbserver.NewObserver(
			solChain,
			rpcClient,
			*solanaChainParams,
			client,
			tss,
			database,
			logger,
			ts,
		)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewObserver error for SOL chain %d", chainID)
			continue
		}

		mapSet[int64, interfaces.ChainObserver](observerMap, chainID, solObserver, onAfterSet)
	}

	// Remove all disabled observers
	mapDeleteMissingKeys(observerMap, presentChainIDs, onBeforeUnset)

	return added, removed, nil
}
