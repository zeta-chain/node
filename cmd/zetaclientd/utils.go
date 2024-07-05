package main

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/pkg/chains"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"

	"github.com/zeta-chain/zetacore/zetaclient/authz"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	btcrpc "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/rpc"
	btcsigner "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/signer"
	evmobserver "github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	evmsigner "github.com/zeta-chain/zetacore/zetaclient/chains/evm/signer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	solanaobserver "github.com/zeta-chain/zetacore/zetaclient/chains/solana/observer"
	solanasigner "github.com/zeta-chain/zetacore/zetaclient/chains/solana/signer"

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

// CreateSignerMap creates a map of ChainSigners for all chains in the config
func CreateSignerMap(
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
		mpiAddress := ethcommon.HexToAddress(evmChainParams.ConnectorContractAddress)
		erc20CustodyAddress := ethcommon.HexToAddress(evmChainParams.Erc20CustodyContractAddress)
		signer, err := evmsigner.NewSigner(
			evmConfig.Chain,
			appContext,
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
		signer, err := btcsigner.NewSigner(btcChain, appContext, tss, ts, logger, btcConfig)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewBTCSigner error for chain %s", btcChain.String())
		} else {
			signerMap[btcChain.ChainId] = signer
		}
	}

	// FIXME: config this
	solChain := chains.SolanaLocalnet
	{
		signer, err := solanasigner.NewSigner(solChain, appContext, tss, ts, logger)
		if err != nil {
			logger.Std.Error().Err(err).Msgf("NewSolanaSigner error for chain %s", solChain.String())
		} else {
			logger.Std.Info().Msgf("NewSolanaSigner for chain %s", solChain.String())
			signerMap[solChain.ChainId] = signer
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
		chainParams, found := appContext.GetEVMChainParams(evmConfig.Chain.ChainId)
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
			appContext,
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
				appContext,
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

	// FIXME: config this
	solChainParams := observertypes.ChainParams{
		GatewayAddress: "94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d",
		IsSupported:    true,
		ChainId:        chains.SolanaLocalnet.ChainId,
	}
	co, err := solanaobserver.NewObserver(appContext, zetacoreClient, solChainParams, tss, dbpath, ts)
	if err != nil {
		logger.Std.Error().Err(err).Msg("NewObserver error for solana chain")
	} else {
		// TODO: config this
		observerMap[solChainParams.ChainId] = co
	}

	return observerMap, nil
}
