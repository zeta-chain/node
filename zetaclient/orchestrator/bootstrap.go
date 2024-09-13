package orchestrator

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	ethrpc2 "github.com/onrik/ethrpc"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/base"
	btcobserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
	btcsigner "github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	evmobserver "github.com/zeta-chain/node/zetaclient/chains/evm/observer"
	evmsigner "github.com/zeta-chain/node/zetaclient/chains/evm/signer"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	solbserver "github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	solanasigner "github.com/zeta-chain/node/zetaclient/chains/solana/signer"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// btcDatabaseFilename is the Bitcoin database file name now used in mainnet,
// so we keep using it here for backward compatibility
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
// Note that data race handling is the responsibility of the caller.
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

		onAfterAdd = func(chainID int64, _ interfaces.ChainSigner) {
			logger.Std.Info().Msgf("Added signer for chain %d", chainID)
			added++
		}

		addSigner = func(chainID int64, signer interfaces.ChainSigner) {
			mapSet[int64, interfaces.ChainSigner](signers, chainID, signer, onAfterAdd)
		}

		onBeforeRemove = func(chainID int64, _ interfaces.ChainSigner) {
			logger.Std.Info().Msgf("Removing signer for chain %d", chainID)
			removed++
		}
	)

	for _, chain := range app.ListChains() {
		// skip ZetaChain
		if chain.IsZeta() {
			continue
		}

		chainID := chain.ID()

		presentChainIDs = append(presentChainIDs, chainID)

		// noop for existing signers
		if mapHas(signers, chainID) {
			continue
		}

		var (
			params   = chain.Params()
			rawChain = chain.RawChain()
		)

		switch {
		case chain.IsEVM():
			var (
				zetaConnectorAddress = ethcommon.HexToAddress(chain.Params().ConnectorContractAddress)
				erc20CustodyAddress  = ethcommon.HexToAddress(chain.Params().Erc20CustodyContractAddress)
				gatewayAddress       = ethcommon.HexToAddress(chain.Params().GatewayAddress)
			)

			cfg, found := app.Config().GetEVMConfig(chainID)
			if !found || cfg.Empty() {
				logger.Std.Warn().Msgf("Unable to find EVM config for chain %d", chainID)
				continue
			}

			signer, err := evmsigner.NewSigner(
				ctx,
				*rawChain,
				tss,
				ts,
				logger,
				cfg.Endpoint,
				zetaConnectorAddress,
				erc20CustodyAddress,
				gatewayAddress,
			)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("Unable to construct signer for EVM chain %d", chainID)
				continue
			}

			addSigner(chainID, signer)
		case chain.IsUTXO():
			cfg, found := app.Config().GetBTCConfig()
			if !found {
				logger.Std.Warn().Msgf("Unable to find UTXO config for chain %d", chainID)
				continue
			}

			signer, err := btcsigner.NewSigner(*rawChain, tss, ts, logger, cfg)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("Unable to construct signer for UTXO chain %d", chainID)
				continue
			}

			addSigner(chainID, signer)
		case chain.IsSolana():
			cfg, found := app.Config().GetSolanaConfig()
			if !found {
				logger.Std.Warn().Msgf("Unable to find SOL config for chain %d", chainID)
				continue
			}

			// create Solana client
			rpcClient := solrpc.New(cfg.Endpoint)
			if rpcClient == nil {
				// should never happen
				logger.Std.Error().Msgf("Unable to create SOL client from endpoint %s", cfg.Endpoint)
				continue
			}

			// try loading Solana relayer key if present
			password := chain.RelayerKeyPassword()
			relayerKey, err := keys.LoadRelayerKey(app.Config().GetRelayerKeyPath(), rawChain.Network, password)
			if err != nil {
				logger.Std.Error().Err(err).Msg("Unable to load Solana relayer key")
				continue
			}

			// create Solana signer
			signer, err := solanasigner.NewSigner(*rawChain, *params, rpcClient, tss, relayerKey, ts, logger)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("Unable to construct signer for SOL chain %d", chainID)
				continue
			}

			addSigner(chainID, signer)
		default:
			logger.Std.Warn().
				Int64("signer.chain_id", chain.ID()).
				Str("signer.chain_name", chain.RawChain().Name).
				Msgf("Unable to create a signer")
		}
	}

	// Remove all disabled signers
	mapDeleteMissingKeys(signers, presentChainIDs, onBeforeRemove)

	return added, removed, nil
}

// CreateChainObserverMap creates a map of interfaces.ChainObserver (by chainID) for all chains in the config.
// - Note (!) that it calls observer.Start() on creation
// - Note that data race handling is the responsibility of the caller.
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

		onAfterAdd = func(_ int64, ob interfaces.ChainObserver) {
			ob.Start(ctx)
			added++
		}

		addObserver = func(chainID int64, ob interfaces.ChainObserver) {
			mapSet[int64, interfaces.ChainObserver](observerMap, chainID, ob, onAfterAdd)
		}

		onBeforeRemove = func(_ int64, ob interfaces.ChainObserver) {
			ob.Stop()
			removed++
		}
	)

	for _, chain := range app.ListChains() {
		// skip ZetaChain
		if chain.IsZeta() {
			continue
		}

		chainID := chain.ID()
		presentChainIDs = append(presentChainIDs, chainID)

		// noop
		if mapHas(observerMap, chainID) {
			continue
		}

		var (
			params    = chain.Params()
			rawChain  = chain.RawChain()
			chainName = rawChain.Name
		)

		switch {
		case chain.IsEVM():
			cfg, found := app.Config().GetEVMConfig(chainID)
			if !found || cfg.Empty() {
				logger.Std.Warn().Msgf("Unable to find EVM config for chain %d", chainID)
				continue
			}

			httpClient, err := metrics.GetInstrumentedHTTPClient(cfg.Endpoint)
			if err != nil {
				logger.Std.Error().Err(err).Str("rpc.endpoint", cfg.Endpoint).Msgf("Unable to create HTTP client")
				continue
			}
			rpcClient, err := ethrpc.DialHTTPWithClient(cfg.Endpoint, httpClient)
			if err != nil {
				logger.Std.Error().Err(err).Str("rpc.endpoint", cfg.Endpoint).Msgf("Unable to dial EVM RPC")
				continue
			}
			evmClient := ethclient.NewClient(rpcClient)

			database, err := db.NewFromSqlite(dbpath, chainName, true)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("Unable to open a database for EVM chain %q", chainName)
				continue
			}

			evmJSONRPCClient := ethrpc2.NewEthRPC(cfg.Endpoint, ethrpc2.WithHttpClient(httpClient))

			// create EVM chain observer
			observer, err := evmobserver.NewObserver(
				ctx,
				*rawChain,
				evmClient,
				evmJSONRPCClient,
				*params,
				client,
				tss,
				cfg.RPCAlertLatency,
				database,
				logger,
				ts,
			)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("NewObserver error for EVM chain %d", chainID)
				continue
			}

			addObserver(chainID, observer)
		case chain.IsUTXO():
			cfg, found := app.Config().GetBTCConfig()
			if !found {
				logger.Std.Warn().Msgf("Unable to find chain params for BTC chain %d", chainID)
				continue
			}

			btcRPC, err := rpc.NewRPCClient(cfg)
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
				*rawChain,
				btcRPC,
				*params,
				client,
				tss,
				cfg.RPCAlertLatency,
				database,
				logger,
				ts,
			)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("NewObserver error for BTC chain %d", chainID)
				continue
			}

			addObserver(chainID, btcObserver)
		case chain.IsSolana():
			cfg, found := app.Config().GetSolanaConfig()
			if !found {
				logger.Std.Warn().Msgf("Unable to find chain params for SOL chain %d", chainID)
				continue
			}

			rpcClient := solrpc.New(cfg.Endpoint)
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
				*rawChain,
				rpcClient,
				*params,
				client,
				tss,
				cfg.RPCAlertLatency,
				database,
				logger,
				ts,
			)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("NewObserver error for SOL chain %d", chainID)
				continue
			}

			addObserver(chainID, solObserver)
		default:
			logger.Std.Warn().
				Int64("observer.chain_id", chain.ID()).
				Str("observer.chain_name", chain.RawChain().Name).
				Msgf("Unable to create an observer")
		}
	}

	// Remove all disabled observers
	mapDeleteMissingKeys(observerMap, presentChainIDs, onBeforeRemove)

	return added, removed, nil
}
