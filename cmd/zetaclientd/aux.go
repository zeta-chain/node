package main

import (
	"github.com/btcsuite/btcd/rpcclient"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
	"github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

func CreateAuthzSigner(granter string, grantee sdk.AccAddress) {
	zetaclient.SetupAuthZSignerList(granter, grantee)
}

func CreateZetaBridge(chainHomeFolder string, config *config.Config) (*zetaclient.ZetaCoreBridge, error) {
	signerPass := "password"
	chainIP := config.ZetaCoreURL
	kb, err := zetaclient.GetKeyringKeybase(config.AuthzHotkey, chainHomeFolder, signerPass)
	if err != nil {
		return nil, err
	}
	granterAddreess, err := cosmos.AccAddressFromBech32(config.AuthzGranter)
	if err != nil {
		return nil, err
	}
	k := zetaclient.NewKeysWithKeybase(kb, granterAddreess, config.AuthzHotkey, signerPass)
	bridge, err := zetaclient.NewZetaCoreBridge(k, chainIP, config.AuthzHotkey, config.ChainID)
	if err != nil {
		return nil, err
	}
	return bridge, nil
}

func CreateSignerMap(tss zetaclient.TSSSigner, logger zerolog.Logger, cfg *config.Config, ts *zetaclient.TelemetryServer) map[common.Chain]zetaclient.ChainSigner {
	signerMap := make(map[common.Chain]zetaclient.ChainSigner)
	for _, chain := range cfg.ChainsEnabled {
		if chain.IsZetaChain() {
			continue
		}
		if common.IsEVMChain(chain.ChainId) {
			mpiAddress := ethcommon.HexToAddress(cfg.EVMChainConfigs[chain.ChainId].CoreParams.ConnectorContractAddress)
			erc20CustodyAddress := ethcommon.HexToAddress(cfg.EVMChainConfigs[chain.ChainId].CoreParams.ERC20CustodyContractAddress)
			signer, err := zetaclient.NewEVMSigner(chain, cfg.EVMChainConfigs[chain.ChainId].Endpoint, tss, config.GetConnectorABI(), config.GetERC20CustodyABI(), mpiAddress, erc20CustodyAddress, logger, ts)
			if err != nil {
				logger.Err(err).Msgf("%s: NewEVMSigner Ethereum error ", chain.String())
				continue
			}
			signerMap[chain] = signer
		} else if common.IsBitcoinChain(chain.ChainId) {
			// FIXME: move the construction of rpcclient to somewhere else
			connCfg := &rpcclient.ConnConfig{
				Host:         cfg.BitcoinConfig.RPCHost,
				User:         cfg.BitcoinConfig.RPCUsername,
				Pass:         cfg.BitcoinConfig.RPCPassword,
				HTTPPostMode: true,
				DisableTLS:   true,
				Params:       cfg.BitcoinConfig.RPCParams,
			}
			client, err := rpcclient.New(connCfg, nil)
			if err != nil {
				logger.Err(err).Msgf("error creating rpc client: %s", err)
				continue
			}
			signer, err := zetaclient.NewBTCSigner(tss, client, logger, ts)
			if err != nil {
				logger.Err(err).Msgf("%s: NewBitcoinSigner Bitcoin error ", chain.String())
				continue
			}
			signerMap[chain] = signer
		}
	}

	return signerMap
}

func CreateChainClientMap(bridge *zetaclient.ZetaCoreBridge, tss zetaclient.TSSSigner, dbpath string, metrics *metrics.Metrics, logger zerolog.Logger, cfg *config.Config, ts *zetaclient.TelemetryServer) map[common.Chain]zetaclient.ChainClient {
	clientMap := make(map[common.Chain]zetaclient.ChainClient)
	for _, chain := range cfg.ChainsEnabled {
		if chain.IsZetaChain() {
			continue
		}
		logger.Info().Msgf("starting observer for : %s ", chain.String())
		var co zetaclient.ChainClient
		var err error
		if common.IsEVMChain(chain.ChainId) {
			co, err = zetaclient.NewEVMChainClient(bridge, tss, dbpath, metrics, logger, cfg, chain, ts)
		} else if common.IsBitcoinChain(chain.ChainId) {
			co, err = zetaclient.NewBitcoinClient(chain, bridge, tss, dbpath, metrics, logger, cfg, ts)
		}
		if err != nil {
			log.Err(err).Msgf("%s NewEVMChainClient", chain.String())
			continue
		}
		clientMap[chain] = co
	}

	return clientMap
}
