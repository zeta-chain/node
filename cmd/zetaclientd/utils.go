package main

import (
	gocontext "context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/zetaclient/authz"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	btcrpc "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/rpc"
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

// CreateSignerMap creates a map of ChainSigners for all chains in the config
func CreateSignerMap(
	ctx gocontext.Context,
	appContext *context.AppContext,
	tss interfaces.TSSSigner,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (map[int64]interfaces.ChainSigner, error) {
	signerMap := make(map[int64]interfaces.ChainSigner)

	// EVM signers
	for _, evmConfig := range appContext.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		evmChainParams, found := appContext.GetEVMChainParams(evmConfig.Chain.ChainId)
		if !found {
			logger.Std.Error().Msgf("ChainParam not found for chain %s", evmConfig.Chain.String())
			continue
		}

		chainName := evmConfig.Chain.ChainName.String()
		mpiAddress := ethcommon.HexToAddress(evmChainParams.ConnectorContractAddress)
		erc20CustodyAddress := ethcommon.HexToAddress(evmChainParams.Erc20CustodyContractAddress)

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
			logger.Std.Error().Err(err).Msgf("NewSigner error for EVM chain %q", chainName)
			continue
		}

		signerMap[evmConfig.Chain.ChainId] = signer
		logger.Std.Info().Msgf("NewSigner succeeded for EVM chain %q", chainName)
	}

	// BTC signer
	btcChain, btcConfig, btcEnabled := appContext.GetBTCChainAndConfig()
	if btcEnabled {
		chainName := btcChain.ChainName.String()

		signer, err := btcsigner.NewSigner(btcChain, tss, ts, logger, btcConfig)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewSigner error for BTC chain %q", chainName)
		} else {
			signerMap[btcChain.ChainId] = signer
			logger.Std.Info().Msgf("NewSigner succeeded for BTC chain %q", chainName)
		}
	}

	return signerMap, nil
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

		// create EVM chain observer
		observer, err := evmobserver.NewObserver(
			ctx,
			evmConfig,
			evmClient,
			*chainParams,
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
	_, chainParams, found := appContext.GetBTCChainParams()
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
