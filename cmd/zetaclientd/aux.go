package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
	"github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

func CreateAuthzSigner(granter string, grantee sdk.AccAddress) {
	zetaclient.SetupAuthZSignerList(granter, grantee)
}

func CreateZetaBridge(cfg *config.Config) (*zetaclient.ZetaCoreBridge, error) {
	chainIP := cfg.ZetaCoreURL

	kb, _, err := zetaclient.GetKeyringKeybase(cfg)
	if err != nil {
		return nil, err
	}

	granterAddreess, err := cosmos.AccAddressFromBech32(cfg.AuthzGranter)
	if err != nil {
		return nil, err
	}

	k := zetaclient.NewKeysWithKeybase(kb, granterAddreess, cfg.AuthzHotkey, "")
	bridge, err := zetaclient.NewZetaCoreBridge(k, chainIP, cfg.AuthzHotkey, cfg.ChainID)
	if err != nil {
		return nil, err
	}

	return bridge, nil
}

func CreateSignerMap(
	tss zetaclient.TSSSigner,
	logger zerolog.Logger,
	cfg *config.Config,
	ts *zetaclient.TelemetryServer,
) (map[common.Chain]zetaclient.ChainSigner, error) {
	signerMap := make(map[common.Chain]zetaclient.ChainSigner)
	// EVM signers
	for _, evmConfig := range cfg.GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		mpiAddress := ethcommon.HexToAddress(evmConfig.CoreParams.ConnectorContractAddress)
		erc20CustodyAddress := ethcommon.HexToAddress(evmConfig.CoreParams.Erc20CustodyContractAddress)
		signer, err := zetaclient.NewEVMSigner(evmConfig.Chain, evmConfig.Endpoint, tss, config.GetConnectorABI(), config.GetERC20CustodyABI(), mpiAddress, erc20CustodyAddress, logger, ts)
		if err != nil {
			logger.Error().Err(err).Msgf("NewEVMSigner error for chain %s", evmConfig.Chain.String())
			continue
		}
		signerMap[evmConfig.Chain] = signer
	}
	// BTC signer
	btcChain, btcConfig, enabled := cfg.GetBTCConfig()
	if enabled {
		signer, err := zetaclient.NewBTCSigner(btcConfig, tss, logger, ts)
		if err != nil {
			logger.Error().Err(err).Msgf("NewBTCSigner error for chain %s", btcChain.String())
		} else {
			signerMap[btcChain] = signer
		}
	}

	return signerMap, nil
}

func CreateChainClientMap(
	bridge *zetaclient.ZetaCoreBridge,
	tss zetaclient.TSSSigner,
	dbpath string,
	metrics *metrics.Metrics,
	logger zerolog.Logger,
	cfg *config.Config,
	ts *zetaclient.TelemetryServer,
) (map[common.Chain]zetaclient.ChainClient, error) {
	clientMap := make(map[common.Chain]zetaclient.ChainClient)
	// EVM clients
	for _, evmConfig := range cfg.GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		co, err := zetaclient.NewEVMChainClient(bridge, tss, dbpath, metrics, logger, cfg, *evmConfig, ts)
		if err != nil {
			logger.Error().Err(err).Msgf("NewEVMChainClient error for chain %s", evmConfig.Chain.String())
			continue
		}
		clientMap[evmConfig.Chain] = co
	}
	// BTC client
	btcChain, btcConfig, enabled := cfg.GetBTCConfig()
	if enabled {
		co, err := zetaclient.NewBitcoinClient(btcChain, bridge, tss, dbpath, metrics, logger, btcConfig, ts)
		if err != nil {
			logger.Error().Err(err).Msgf("NewBitcoinClient error for chain %s", btcChain.String())

		} else {
			clientMap[btcChain] = co
		}
	}

	return clientMap, nil
}
