package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/zetacore/zetaclient/authz"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	btcsigner "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/signer"
	evmobserver "github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	evmsigner "github.com/zeta-chain/zetacore/zetaclient/chains/evm/signer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

func CreateAuthzSigner(granter string, grantee sdk.AccAddress) {
	authz.SetupAuthZSignerList(granter, grantee)
}

func CreateZetacoreClient(
	cfg config.Config,
	telemetry *metrics.TelemetryServer,
	hotkeyPassword string,
) (*zetacore.Client, error) {
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

	client, err := zetacore.NewClient(k, chainIP, hotKey, cfg.ChainID, cfg.HsmMode, telemetry)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CreateSignerMap(
	appContext *context.AppContext,
	tss interfaces.TSSSigner,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (map[int64]interfaces.ChainSigner, error) {
	coreContext := appContext.ZetacoreContext()
	signerMap := make(map[int64]interfaces.ChainSigner)

	// EVM signers
	for _, evmConfig := range appContext.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		evmChainParams, found := coreContext.GetEVMChainParams(evmConfig.Chain.ChainId)
		if !found {
			logger.Std.Error().Msgf("ChainParam not found for chain %s", evmConfig.Chain.String())
			continue
		}
		mpiAddress := ethcommon.HexToAddress(evmChainParams.ConnectorContractAddress)
		erc20CustodyAddress := ethcommon.HexToAddress(evmChainParams.Erc20CustodyContractAddress)
		signer, err := evmsigner.NewSigner(
			evmConfig.Chain,
			evmConfig.Endpoint,
			tss,
			config.GetConnectorABI(),
			config.GetERC20CustodyABI(),
			mpiAddress,
			erc20CustodyAddress,
			coreContext,
			logger,
			ts)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewEVMSigner error for chain %s", evmConfig.Chain.String())
			continue
		}
		signerMap[evmConfig.Chain.ChainId] = signer
	}
	// BTC signer
	btcChain, btcConfig, enabled := appContext.GetBTCChainAndConfig()
	if enabled {
		signer, err := btcsigner.NewSigner(btcConfig, tss, logger, ts, coreContext)
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
	observerMap := make(map[int64]interfaces.ChainObserver)
	// EVM observers
	for _, evmConfig := range appContext.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		_, found := appContext.ZetacoreContext().GetEVMChainParams(evmConfig.Chain.ChainId)
		if !found {
			logger.Std.Error().Msgf("ChainParam not found for chain %s", evmConfig.Chain.String())
			continue
		}
		co, err := evmobserver.NewObserver(appContext, zetacoreClient, tss, dbpath, logger, evmConfig, ts)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewObserver error for evm chain %s", evmConfig.Chain.String())
			continue
		}
		observerMap[evmConfig.Chain.ChainId] = co
	}
	// BTC observer
	btcChain, btcConfig, enabled := appContext.GetBTCChainAndConfig()
	if enabled {
		co, err := btcobserver.NewObserver(appContext, btcChain, zetacoreClient, tss, dbpath, logger, btcConfig, ts)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewObserver error for bitcoin chain %s", btcChain.String())

		} else {
			observerMap[btcChain.ChainId] = co
		}
	}

	return observerMap, nil
}
