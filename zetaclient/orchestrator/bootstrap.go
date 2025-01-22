package orchestrator

import (
	"context"
	"fmt"

	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/chains"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	solbserver "github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	solanasigner "github.com/zeta-chain/node/zetaclient/chains/solana/signer"
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
	tonobserver "github.com/zeta-chain/node/zetaclient/chains/ton/observer"
	tonsigner "github.com/zeta-chain/node/zetaclient/chains/ton/signer"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// CreateSignerMap creates a map of interfaces.ChainSigner (by chainID) for all chains in the config.
// Note that signer construction failure for a chain does not prevent the creation of signers for other chains.
func CreateSignerMap(
	ctx context.Context,
	tss interfaces.TSSSigner,
	logger base.Logger,
) (map[int64]interfaces.ChainSigner, error) {
	signers := make(map[int64]interfaces.ChainSigner)
	_, _, err := syncSignerMap(ctx, tss, logger, &signers)
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
			logger.Std.Info().Int64(logs.FieldChain, chainID).Msg("Added signer")
			added++
		}

		addSigner = func(chainID int64, signer interfaces.ChainSigner) {
			mapSet[int64, interfaces.ChainSigner](signers, chainID, signer, onAfterAdd)
		}

		onBeforeRemove = func(chainID int64, _ interfaces.ChainSigner) {
			logger.Std.Info().Int64(logs.FieldChain, chainID).Msg("Removing signer")
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
			// managed by orchestrator V2
			continue

		case chain.IsBitcoin():
			// managed by orchestrator V2
			continue

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
			signer, err := solanasigner.NewSigner(*rawChain, *params, rpcClient, tss, relayerKey, logger)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("Unable to construct signer for SOL chain %d", chainID)
				continue
			}

			addSigner(chainID, signer)
		case chain.IsTON():
			cfg, found := app.Config().GetTONConfig()
			if !found {
				logger.Std.Warn().Msgf("Unable to find TON config for chain %d", chainID)
				continue
			}

			tonClient, gateway, err := makeTONClient(ctx, cfg, chain.Params().GatewayAddress)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("Unable to create TON client for chain %d", chainID)
				continue
			}

			tonSigner := tonsigner.New(
				base.NewSigner(*rawChain, tss, logger),
				tonClient,
				gateway,
			)

			addSigner(chainID, tonSigner)
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

		onAfterAdd = func(chainID int64, ob interfaces.ChainObserver) {
			logger.Std.Info().Int64(logs.FieldChain, chainID).Msg("Added observer")
			ob.Start(ctx)
			added++
		}

		addObserver = func(chainID int64, ob interfaces.ChainObserver) {
			mapSet[int64, interfaces.ChainObserver](observerMap, chainID, ob, onAfterAdd)
		}

		onBeforeRemove = func(chainID int64, ob interfaces.ChainObserver) {
			logger.Std.Info().Int64(logs.FieldChain, chainID).Msg("Removing observer")
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
			// managed by orchestrator V2
			continue

		case chain.IsBitcoin():
			// managed by orchestrator V2
			continue

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
				database,
				logger,
				ts,
			)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("NewObserver error for SOL chain %d", chainID)
				continue
			}

			addObserver(chainID, solObserver)
		case chain.IsTON():
			cfg, found := app.Config().GetTONConfig()
			if !found {
				logger.Std.Warn().Msgf("Unable to find chain params for TON chain %d", chainID)
				continue
			}

			database, err := db.NewFromSqlite(dbpath, chainName, true)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("unable to open database for TON chain %d", chainID)
				continue
			}

			baseObserver, err := base.NewObserver(
				*rawChain,
				*params,
				client,
				tss,
				base.DefaultBlockCacheSize,
				ts,
				database,
				logger,
			)

			if err != nil {
				logger.Std.Error().Err(err).Msgf("Unable to create base observer for TON chain %d", chainID)
				continue
			}

			tonClient, gateway, err := makeTONClient(ctx, cfg, chain.Params().GatewayAddress)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("Unable to create TON client for chain %d", chainID)
				continue
			}

			tonObserver, err := tonobserver.New(baseObserver, tonClient, gateway)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("Unable to create TON observer for chain %d", chainID)
				continue
			}

			addObserver(chainID, tonObserver)
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

func btcDatabaseFileName(chain chains.Chain) string {
	// legacyBTCDatabaseFilename is the Bitcoin database file name now used in mainnet and testnet3
	// so we keep using it here for backward compatibility
	const legacyBTCDatabaseFilename = "btc_chain_client"

	// For additional bitcoin networks, we use the chain name as the database file name
	switch chain.ChainId {
	case chains.BitcoinMainnet.ChainId, chains.BitcoinTestnet.ChainId:
		return legacyBTCDatabaseFilename
	default:
		return fmt.Sprintf("%s_%s", legacyBTCDatabaseFilename, chain.Name)
	}
}

func makeTONClient(
	ctx context.Context,
	cfg config.TONConfig,
	gatewayAddr string,
) (*liteapi.Client, *toncontracts.Gateway, error) {
	client, err := liteapi.NewFromSource(ctx, cfg.LiteClientConfigURL)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to create TON liteapi")
	}

	gatewayID, err := ton.ParseAccountID(gatewayAddr)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Unable to parse gateway address %q", gatewayAddr)
	}

	gw := toncontracts.NewGateway(gatewayID)

	return client, gw, nil
}
