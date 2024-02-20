package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
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
	tss interfaces.TSSSigner,
	logger zerolog.Logger,
	cfg *config.Config,
	ts *metrics.TelemetryServer,
) (map[common.Chain]interfaces.ChainSigner, error) {
	signerMap := make(map[common.Chain]interfaces.ChainSigner)
	// EVM signers
	for _, evmConfig := range cfg.GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		mpiAddress := ethcommon.HexToAddress(evmConfig.ChainParams.ConnectorContractAddress)
		erc20CustodyAddress := ethcommon.HexToAddress(evmConfig.ChainParams.Erc20CustodyContractAddress)
		signer, err := evm.NewEVMSigner(evmConfig.Chain, evmConfig.Endpoint, tss, config.GetConnectorABI(), config.GetERC20CustodyABI(), mpiAddress, erc20CustodyAddress, logger, ts)
		if err != nil {
			logger.Error().Err(err).Msgf("NewEVMSigner error for chain %s", evmConfig.Chain.String())
			continue
		}
		signerMap[evmConfig.Chain] = signer
	}
	// BTC signer
	btcChain, btcConfig, enabled := cfg.GetBTCConfig()
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
	bridge *zetabridge.ZetaCoreBridge,
	tss interfaces.TSSSigner,
	dbpath string,
	logger zerolog.Logger,
	cfg *config.Config,
	ts *metrics.TelemetryServer,
) (map[common.Chain]interfaces.ChainClient, error) {
	clientMap := make(map[common.Chain]interfaces.ChainClient)
	// EVM clients
	for _, evmConfig := range cfg.GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		co, err := evm.NewEVMChainClient(bridge, tss, dbpath, logger, cfg, *evmConfig, ts)
		if err != nil {
			logger.Error().Err(err).Msgf("NewEVMChainClient error for chain %s", evmConfig.Chain.String())
			continue
		}
		clientMap[evmConfig.Chain] = co
	}
	// BTC client
	btcChain, btcConfig, enabled := cfg.GetBTCConfig()
	if enabled {
		co, err := bitcoin.NewBitcoinClient(btcChain, bridge, tss, dbpath, logger, btcConfig, ts)
		if err != nil {
			logger.Error().Err(err).Msgf("NewBitcoinClient error for chain %s", btcChain.String())

		} else {
			clientMap[btcChain] = co
		}
	}

	return clientMap, nil
}
