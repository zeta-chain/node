package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	appcontext "github.com/zeta-chain/zetacore/zetaclient/app_context"
	"github.com/zeta-chain/zetacore/zetaclient/authz"
	"github.com/zeta-chain/zetacore/zetaclient/bitcoin"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
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

func CreateZetaBridge(cfg config.Config, telemetry *metrics.TelemetryServer, hotkeyPassword string) (*zetabridge.ZetaCoreBridge, error) {
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

	bridge, err := zetabridge.NewZetaCoreBridge(k, chainIP, hotKey, cfg.ChainID, cfg.HsmMode, telemetry)
	if err != nil {
		return nil, err
	}

	return bridge, nil
}

func CreateSignerMap(
	appContext *appcontext.AppContext,
	tss interfaces.TSSSigner,
	loggers clientcommon.ClientLogger,
	ts *metrics.TelemetryServer,
) (map[int64]interfaces.ChainSigner, error) {
	coreContext := appContext.ZetaCoreContext()
	signerMap := make(map[int64]interfaces.ChainSigner)

	// EVM signers
	for _, evmConfig := range appContext.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		evmChainParams, found := coreContext.GetEVMChainParams(evmConfig.Chain.ChainId)
		if !found {
			loggers.Std.Error().Msgf("ChainParam not found for chain %s", evmConfig.Chain.String())
			continue
		}
		mpiAddress := ethcommon.HexToAddress(evmChainParams.ConnectorContractAddress)
		erc20CustodyAddress := ethcommon.HexToAddress(evmChainParams.Erc20CustodyContractAddress)
		signer, err := evm.NewEVMSigner(
			evmConfig,
			tss,
			config.GetConnectorABI(),
			config.GetERC20CustodyABI(),
			mpiAddress,
			erc20CustodyAddress,
			coreContext,
			loggers,
			ts)
		if err != nil {
			loggers.Std.Error().Err(err).Msgf("NewEVMSigner error for chain %s", evmConfig.Chain.String())
			continue
		}
		signerMap[evmConfig.Chain.ChainId] = signer
	}
	// BTC signer
	btcChain, btcConfig, enabled := appContext.GetBTCChainAndConfig()
	if enabled {
		signer, err := bitcoin.NewBTCSigner(btcConfig, tss, loggers, ts, coreContext)
		if err != nil {
			loggers.Std.Error().Err(err).Msgf("NewBTCSigner error for chain %s", btcChain.String())
		} else {
			signerMap[btcChain.ChainId] = signer
		}
	}

	return signerMap, nil
}

func CreateChainClientMap(
	appContext *appcontext.AppContext,
	bridge *zetabridge.ZetaCoreBridge,
	tss interfaces.TSSSigner,
	dbpath string,
	loggers clientcommon.ClientLogger,
	ts *metrics.TelemetryServer,
) (map[int64]interfaces.ChainClient, error) {
	clientMap := make(map[int64]interfaces.ChainClient)
	// EVM clients
	for _, evmConfig := range appContext.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		_, found := appContext.ZetaCoreContext().GetEVMChainParams(evmConfig.Chain.ChainId)
		if !found {
			loggers.Std.Error().Msgf("ChainParam not found for chain %s", evmConfig.Chain.String())
			continue
		}
		co, err := evm.NewEVMChainClient(appContext, bridge, tss, dbpath, loggers, evmConfig, ts)
		if err != nil {
			loggers.Std.Error().Err(err).Msgf("NewEVMChainClient error for chain %s", evmConfig.Chain.String())
			continue
		}
		clientMap[evmConfig.Chain.ChainId] = co
	}
	// BTC client
	btcChain, btcConfig, enabled := appContext.GetBTCChainAndConfig()
	if enabled {
		co, err := bitcoin.NewBitcoinClient(appContext, btcChain, bridge, tss, dbpath, loggers, btcConfig, ts)
		if err != nil {
			loggers.Std.Error().Err(err).Msgf("NewBitcoinClient error for chain %s", btcChain.String())

		} else {
			clientMap[btcChain.ChainId] = co
		}
	}

	return clientMap, nil
}
