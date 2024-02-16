package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
	appcontext "github.com/zeta-chain/zetacore/zetaclient/app_context"
	"github.com/zeta-chain/zetacore/zetaclient/authz"
	"github.com/zeta-chain/zetacore/zetaclient/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"

	"github.com/zeta-chain/zetacore/zetaclient/evm"
)

func CreateAuthzSigner(granter string, grantee sdk.AccAddress) {
	authz.SetupAuthZSignerList(granter, grantee)
}

func CreateZetaBridge(cfg *config.Config, telemetry *metrics.TelemetryServer, hotkeyPassword string) (*zetabridge.ZetaCoreBridge, error) {
	hotKey := cfg.AuthzHotkey
	if cfg.HsmMode {
		hotKey = cfg.HsmHotKey
	}

	chainIP := cfg.ZetaCoreURL

	kb, _, err := keys.GetKeyringKeybase(cfg, hotkeyPassword)
	if err != nil {
		return nil, err
	}

	granterAddreess, err := cosmos.AccAddressFromBech32(cfg.AuthzGranter)
	if err != nil {
		return nil, err
	}

	k := keys.NewKeysWithKeybase(kb, granterAddreess, cfg.AuthzHotkey, hotkeyPassword)

	bridge, err := zetabridge.NewZetaCoreBridge(k, chainIP, hotKey, cfg.ChainID, cfg.HsmMode, telemetry)
	if err != nil {
		return nil, err
	}

	return bridge, nil
}

func CreateSignerMap(
	appContext *appcontext.AppContext,
	tss interfaces.TSSSigner,
	ts *metrics.TelemetryServer,
) (map[common.Chain]interfaces.ChainSigner, error) {
	signerMap := make(map[common.Chain]interfaces.ChainSigner)
	logger := appContext.MasterLogger()
	// EVM signers
	for _, evmConfig := range appContext.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		evmChainParams, found := appContext.ZetaCoreContext().GetEVMChainParams(evmConfig.Chain.ChainId)
		if !found {
			logger.Error().Msgf("ChainParam not found for chain %s", evmConfig.Chain.String())
			continue
		}
		mpiAddress := ethcommon.HexToAddress(evmChainParams.ConnectorContractAddress)
		erc20CustodyAddress := ethcommon.HexToAddress(evmChainParams.Erc20CustodyContractAddress)
		signer, err := evm.NewEVMSigner(evmConfig.Chain, evmConfig.Endpoint, tss, config.GetConnectorABI(), config.GetERC20CustodyABI(), mpiAddress, erc20CustodyAddress, logger, ts)
		if err != nil {
			logger.Error().Err(err).Msgf("NewEVMSigner error for chain %s", evmConfig.Chain.String())
			continue
		}
		signerMap[evmConfig.Chain] = signer
	}
	// BTC signer
	btcChain, btcConfig, enabled := appContext.Config().GetBTCConfig()
	if enabled {
		signer, err := bitcoin.NewBTCSigner(btcConfig, tss, logger, ts)
		if err != nil {
			logger.Error().Err(err).Msgf("NewBTCSigner error for chain %s", btcChain.String())
		} else {
			signerMap[btcChain] = signer
		}
	}

	return signerMap, nil
}

func CreateChainClientMap(
	appContext *appcontext.AppContext,
	bridge *zetabridge.ZetaCoreBridge,
	tss interfaces.TSSSigner,
	dbpath string,
	metrics *metrics.Metrics,
	ts *metrics.TelemetryServer,
) (map[common.Chain]interfaces.ChainClient, error) {
	clientMap := make(map[common.Chain]interfaces.ChainClient)
	logger := appContext.MasterLogger()
	// EVM clients
	for _, evmConfig := range appContext.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		co, err := evm.NewEVMChainClient(appContext, bridge, tss, dbpath, metrics, *evmConfig, ts)
		if err != nil {
			logger.Error().Err(err).Msgf("NewEVMChainClient error for chain %s", evmConfig.Chain.String())
			continue
		}
		clientMap[evmConfig.Chain] = co
	}
	// BTC client
	btcChain, btcConfig, enabled := appContext.Config().GetBTCConfig()
	if enabled {
		co, err := bitcoin.NewBitcoinClient(appContext, btcChain, bridge, tss, dbpath, metrics, logger, btcConfig, ts)
		if err != nil {
			logger.Error().Err(err).Msgf("NewBitcoinClient error for chain %s", btcChain.String())

		} else {
			clientMap[btcChain] = co
		}
	}

	return clientMap, nil
}
