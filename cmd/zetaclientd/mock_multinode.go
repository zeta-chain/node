package main

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	mcconfig "github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"os"
)

func CreateMetaBridge(chainHomeFoler string, signerName string, signerPass string) (*mc.MetachainBridge, bool) {
	kb, _, err := mc.GetKeyringKeybase(chainHomeFoler, signerName, signerPass)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to get keyring keybase")
		return nil, true
	}

	k := mc.NewKeysWithKeybase(kb, signerName, signerPass)

	chainIP := os.Getenv("CHAIN_IP")
	if chainIP == "" {
		chainIP = "127.0.0.1"
	}

	bridge, err := mc.NewMetachainBridge(k, chainIP, signerName)
	if err != nil {
		log.Fatal().Err(err).Msg("NewMetachainBridge")
		return nil, true
	}
	return bridge, false
}

func CreateSignerMap(tss mc.TSSSigner) (map[common.Chain]*mc.Signer, error) {
	ethMPIAddress := ethcommon.HexToAddress(mcconfig.Chains[common.GoerliChain.String()].ConnectorContractAddress)
	ethSigner, err := mc.NewSigner(common.GoerliChain, mcconfig.GOERLI_ENDPOINT, tss, mcconfig.CONNECTOR_ABI_STRING, ethMPIAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner Ethereum error ")
		return nil, err
	}
	bscMPIAddress := ethcommon.HexToAddress(mcconfig.Chains[common.BSCTestnetChain.String()].ConnectorContractAddress)
	bscSigner, err := mc.NewSigner(common.BSCTestnetChain, mcconfig.BSCTESTNET_ENDPOINT, tss, mcconfig.CONNECTOR_ABI_STRING, bscMPIAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner BSC error")
		return nil, err
	}
	polygonMPIAddress := ethcommon.HexToAddress(mcconfig.Chains[common.MumbaiChain.String()].ConnectorContractAddress)
	polygonSigner, err := mc.NewSigner(common.MumbaiChain, mcconfig.MUMBAI_ENDPOINT, tss, mcconfig.CONNECTOR_ABI_STRING, polygonMPIAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner POLYGON error")
		return nil, err
	}
	ropstenMPIAddress := ethcommon.HexToAddress(mcconfig.Chains[common.RopstenChain.String()].ConnectorContractAddress)
	ropstenSigner, err := mc.NewSigner(common.RopstenChain, mcconfig.ROPSTEN_ENDPOINT, tss, mcconfig.CONNECTOR_ABI_STRING, ropstenMPIAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner ROPSTEN error")
		return nil, err
	}
	signerMap := map[common.Chain]*mc.Signer{
		common.GoerliChain:     ethSigner,
		common.BSCTestnetChain: bscSigner,
		common.MumbaiChain:     polygonSigner,
		common.RopstenChain:    ropstenSigner,
	}

	return signerMap, nil
}

func CreateChainClientMap(bridge *mc.MetachainBridge, tss mc.TSSSigner, dbpath string, metrics *metrics.Metrics) (*map[common.Chain]*mc.ChainObserver, error) {
	log.Info().Msg("starting eth observer...")
	clientMap := make(map[common.Chain]*mc.ChainObserver)
	eth1, err := mc.NewChainObserver(common.GoerliChain, bridge, tss, dbpath, metrics)
	if err != nil {
		log.Err(err).Msg("ETH NewChainObserver")
		return nil, err
	}
	clientMap[common.GoerliChain] = eth1
	eth1.Start()

	log.Info().Msg("starting bsc observer...")
	bsc1, err := mc.NewChainObserver(common.BSCTestnetChain, bridge, tss, dbpath, metrics)
	if err != nil {
		log.Err(err).Msg("BSC NewChainObserver")
		return nil, err
	}
	clientMap[common.BSCTestnetChain] = bsc1
	bsc1.Start()

	log.Info().Msg("starting polygon observer...")
	poly1, err := mc.NewChainObserver(common.MumbaiChain, bridge, tss, dbpath, metrics)
	if err != nil {
		log.Err(err).Msg("POLYGON NewChainObserver")
		return nil, err
	}
	clientMap[common.MumbaiChain] = poly1
	poly1.Start()

	log.Info().Msg("starting ropsten observer...")
	ropsten1, err := mc.NewChainObserver(common.RopstenChain, bridge, tss, dbpath, metrics)
	if err != nil {
		log.Err(err).Msg("ROPSTEN NewChainObserver")
		return nil, err
	}
	clientMap[common.RopstenChain] = ropsten1
	ropsten1.Start()

	return &clientMap, nil
}
