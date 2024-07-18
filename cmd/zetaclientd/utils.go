package main

import (
	gocontext "context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/authz"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	btcrpc "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/rpc"
	evmobserver "github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/db"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// backwards compatibility
const btcDatabaseFilename = "btc_chain_client"

func CreateAuthzSigner(granter string, grantee sdk.AccAddress) {
	authz.SetupAuthZSignerList(granter, grantee)
}

func CreateZetacoreClient(cfg config.Config, hotkeyPassword string, logger zerolog.Logger) (*zetacore.Client, error) {
	hotKey := cfg.AuthzHotkey
	if cfg.HsmMode {
		hotKey = cfg.HsmHotKey
	}

	chainIP := cfg.ZetaCoreURL

	kb, _, err := keys.GetKeyringKeybase(cfg, hotkeyPassword)
	if err != nil {
		return nil, err
	}

	granterAddreess, err := sdk.AccAddressFromBech32(cfg.AuthzGranter)
	if err != nil {
		return nil, err
	}

	k := keys.NewKeysWithKeybase(kb, granterAddreess, cfg.AuthzHotkey, hotkeyPassword)

	client, err := zetacore.NewClient(k, chainIP, hotKey, cfg.ChainID, cfg.HsmMode, logger)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// CreateChainObserverMap creates a map of ChainObservers for all chains in the config
func CreateChainObserverMap(
	ctx gocontext.Context,
	appContext *context.AppContext,
	zetacoreClient *zetacore.Client,
	tss interfaces.TSSSigner,
	dbpath string,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (map[int64]interfaces.ChainObserver, error) {
	observerMap := make(map[int64]interfaces.ChainObserver)
	// EVM observers
	for _, evmConfig := range appContext.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		chainParams, found := appContext.GetEVMChainParams(evmConfig.Chain.ChainId)
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
			zetacoreClient,
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
	btcChain, btcConfig, btcEnabled := appContext.GetBTCChainAndConfig()
	if btcEnabled {
		_, chainParams, found := appContext.GetBTCChainParams()
		if !found {
			return nil, fmt.Errorf("BTC is enabled, but chains params not found")
		}

		btcObserver, err := createBTCObserver(
			dbpath,
			btcConfig,
			btcChain,
			*chainParams,
			zetacoreClient,
			tss,
			logger,
			ts,
		)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewObserver error for BTC chain %s", btcChain.ChainName.String())
		} else {
			observerMap[btcChain.ChainId] = btcObserver
		}
	}

	return observerMap, nil
}

func createBTCObserver(
	dbPath string,
	cfg config.BTCConfig,
	chain chains.Chain,
	chainParams observertypes.ChainParams,
	client *zetacore.Client,
	tss interfaces.TSSSigner,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (*btcobserver.Observer, error) {
	btcClient, err := btcrpc.NewRPCClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create rpc client for BTC chain")
	}

	database, err := db.NewFromSqlite(dbPath, btcDatabaseFilename, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open a database for BTC chain")
	}

	// create BTC chain observer
	observer, err := btcobserver.NewObserver(
		chain,
		btcClient,
		chainParams,
		client,
		tss,
		database,
		logger,
		ts,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create observer for BTC chain")
	}

	return observer, nil
}
