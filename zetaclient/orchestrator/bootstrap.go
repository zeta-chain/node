package orchestrator

import (
	"context"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/rpc"
	btcsigner "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/signer"
	evmobserver "github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	evmsigner "github.com/zeta-chain/zetacore/zetaclient/chains/evm/signer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/db"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
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
	signers := make(signerMap)
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
	signers *signerMap,
) (int, int, error) {
	if signers == nil {
		return 0, 0, errors.New("signers map is nil")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to get app context")
	}

	var added int

	presentChainIDs := make([]int64, 0)

	// EVM signers
	for _, evmConfig := range app.Config().GetAllEVMConfigs() {
		chainID := evmConfig.Chain.ChainId

		if evmConfig.Chain.IsZetaChain() {
			continue
		}

		evmChainParams, found := app.GetEVMChainParams(chainID)
		if !found {
			logger.Std.Warn().Msgf("Unable to find chain params for EVM chain %d", chainID)
			continue
		}

		presentChainIDs = append(presentChainIDs, chainID)

		// noop for existing signers
		if signers.has(chainID) {
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

		signers.set(chainID, signer)
		logger.Std.Info().Msgf("Added signer for EVM chain %d", chainID)
		added++
	}

	// BTC signer
	btcChain, btcConfig, btcEnabled := app.GetBTCChainAndConfig()
	if btcEnabled {
		chainID := btcChain.ChainId

		presentChainIDs = append(presentChainIDs, chainID)

		if !signers.has(chainID) {
			utxoSigner, err := btcsigner.NewSigner(btcChain, tss, ts, logger, btcConfig)
			if err != nil {
				logger.Std.Error().Err(err).Msgf("Unable to construct signer for UTXO chain %d", chainID)
			} else {
				signers.set(chainID, utxoSigner)
				logger.Std.Info().Msgf("Added signer for UTXO chain %d", chainID)
				added++
			}
		}
	}

	// Remove all disabled signers
	removed := signers.unsetMissing(presentChainIDs, logger.Std)

	return added, removed, nil
}

type signerMap map[int64]interfaces.ChainSigner

func (m *signerMap) has(chainID int64) bool {
	_, ok := (*m)[chainID]
	return ok
}

func (m *signerMap) set(chainID int64, signer interfaces.ChainSigner) {
	(*m)[chainID] = signer
}

func (m *signerMap) unset(chainID int64, logger zerolog.Logger) bool {
	if _, ok := (*m)[chainID]; !ok {
		return false
	}

	logger.Info().Msgf("Removing signer for chain %d", chainID)
	delete(*m, chainID)

	return true
}

// unsetMissing removes signers from the map IF they are not in the enabledChains list.
func (m *signerMap) unsetMissing(enabledChains []int64, logger zerolog.Logger) int {
	enabledMap := make(map[int64]struct{}, len(enabledChains))
	for _, id := range enabledChains {
		enabledMap[id] = struct{}{}
	}

	var removed int

	for id := range *m {
		if _, ok := enabledMap[id]; !ok {
			m.unset(id, logger)
			removed++
		}
	}

	return removed
}

// CreateChainObserverMap creates a map of ChainObservers for all chains in the config
func CreateChainObserverMap(
	ctx context.Context,
	client *zetacore.Client,
	tss interfaces.TSSSigner,
	dbpath string,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (map[int64]interfaces.ChainObserver, error) {
	observerMap := make(map[int64]interfaces.ChainObserver)

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get app context")
	}

	// EVM observers
	for _, evmConfig := range app.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}

		chainParams, found := app.GetEVMChainParams(evmConfig.Chain.ChainId)
		if !found {
			logger.Std.Error().Msgf("ChainParam not found for chain %s", evmConfig.Chain.String())
			continue
		}

		// create EVM client
		evmClient, err := ethclient.Dial(evmConfig.Endpoint)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("error dailing endpoint %q", evmConfig.Endpoint)
			continue
		}

		chainName := evmConfig.Chain.ChainName.String()

		database, err := db.NewFromSqlite(dbpath, chainName, true)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("Unable to open a database for EVM chain %q", chainName)
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
		observerMap[evmConfig.Chain.ChainId] = observer
	}

	// create BTC chain observer
	btcChain, btcConfig, btcEnabled := app.GetBTCChainAndConfig()
	if !btcEnabled {
		return observerMap, nil
	}

	_, btcChainParams, found := app.GetBTCChainParams()
	if !found {
		return nil, fmt.Errorf("BTC is enabled, but chains params not found")
	}

	btcClient, err := rpc.NewRPCClient(btcConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create rpc client for BTC chain")
	}

	btcDatabase, err := db.NewFromSqlite(dbpath, btcDatabaseFilename, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open a database for BTC chain")
	}

	btcObserver, err := btcobserver.NewObserver(
		btcChain,
		btcClient,
		*btcChainParams,
		client,
		tss,
		btcDatabase,
		logger,
		ts,
	)
	if err != nil {
		logger.Std.Error().Err(err).Msgf("NewObserver error for BTC chain %s", btcChain.ChainName.String())
	} else {
		observerMap[btcChain.ChainId] = btcObserver
	}

	return observerMap, nil
}
