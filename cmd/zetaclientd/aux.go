package main

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	mcconfig "github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"os"
)

func CreateZetaBridge(chainHomeFoler string, signerName string, signerPass string) (*mc.ZetaCoreBridge, bool) {
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

	bridge, err := mc.NewZetaCoreBridge(k, chainIP, signerName)
	if err != nil {
		log.Fatal().Err(err).Msg("NewZetaCoreBridge")
		return nil, true
	}
	return bridge, false
}

func CreateSignerMap(tss mc.TSSSigner) (map[zetaObserverTypes.Chain]*mc.Signer, error) {
	signerMap := make(map[zetaObserverTypes.Chain]*mc.Signer)
	supportedChains := mc.GetSupportedChains()
	for _, chain := range supportedChains {
		mpiAddress := ethcommon.HexToAddress(mcconfig.Chains[chain.String()].ConnectorContractAddress)
		signer, err := mc.NewSigner(*chain, mcconfig.Chains[chain.String()].Endpoint, tss, mcconfig.ConnectorAbiString, mpiAddress)
		if err != nil {
			log.Fatal().Err(err).Msg("NewSigner Ethereum error ")
			return nil, err
		}
		signerMap[*chain] = signer
	}

	return signerMap, nil
}

func CreateChainClientMap(bridge *mc.ZetaCoreBridge, tss mc.TSSSigner, dbpath string, metrics *metrics.Metrics) (*map[zetaObserverTypes.Chain]*mc.ChainObserver, error) {
	clientMap := make(map[zetaObserverTypes.Chain]*mc.ChainObserver)
	supportedChains := mc.GetSupportedChains()
	for _, chain := range supportedChains {
		log.Info().Msgf("starting %s observer...", chain)
		co, err := mc.NewChainObserver(*chain, bridge, tss, dbpath, metrics)
		if err != nil {
			log.Err(err).Msgf("%s NewChainObserver", chain)
			return nil, err
		}
		clientMap[*chain] = co
	}

	return &clientMap, nil
}
