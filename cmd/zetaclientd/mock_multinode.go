package main

import (
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	mcconfig "github.com/zeta-chain/zetacore/zetaclient/config"
	"os"
)

func GetZetaTestSignature() mc.TestSigner {
	pkstring := os.Getenv("PRIVKEY")
	if pkstring == "" {
		log.Info().Msgf("missing env variable PRIVKEY; using default with address %s", mcconfig.TSS_TEST_ADDRESS)
		pkstring = mcconfig.TSS_TEST_PRIVKEY
	}
	privateKey, err := crypto.HexToECDSA(pkstring)
	if err != nil {
		log.Err(err).Msg("TEST private key error")
		os.Exit(1)
	}
	tss := mc.TestSigner{
		PrivKey: privateKey,
	}
	log.Debug().Msg(fmt.Sprintf("tss key address: %s", tss.Address()))

	return tss
}

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
	ethMPIAddress := ethcommon.HexToAddress(mcconfig.Chains["ETH"].MPIContractAddress)
	ethSigner, err := mc.NewSigner(common.ETHChain, mcconfig.ETH_ENDPOINT, tss.Address(), tss, mcconfig.MPI_ABI_STRING, ethMPIAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner Ethereum error ")
		return nil, err
	}
	bscMPIAddress := ethcommon.HexToAddress(mcconfig.Chains["BSC"].MPIContractAddress)
	bscSigner, err := mc.NewSigner(common.BSCChain, mcconfig.BSC_ENDPOINT, tss.Address(), tss, mcconfig.MPI_ABI_STRING, bscMPIAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner BSC error")
		return nil, err
	}
	polygonMPIAddress := ethcommon.HexToAddress(mcconfig.Chains["POLYGON"].MPIContractAddress)
	polygonSigner, err := mc.NewSigner(common.POLYGONChain, mcconfig.POLY_ENDPOINT, tss.Address(), tss, mcconfig.MPI_ABI_STRING, polygonMPIAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner POLYGON error")
		return nil, err
	}
	signerMap := map[common.Chain]*mc.Signer{
		common.ETHChain:     ethSigner,
		common.BSCChain:     bscSigner,
		common.POLYGONChain: polygonSigner,
	}

	return signerMap, nil
}

func CreateChainClientMap(bridge *mc.MetachainBridge, tss mc.TSSSigner, dbpath string) (*map[common.Chain]*mc.ChainObserver, error) {
	log.Info().Msg("starting eth observer...")
	clientMap := make(map[common.Chain]*mc.ChainObserver)
	eth1, err := mc.NewChainObserver(common.ETHChain, bridge, tss, dbpath)
	if err != nil {
		log.Err(err).Msg("ETH NewChainObserver")
		return nil, err
	}
	clientMap[common.ETHChain] = eth1
	go eth1.WatchRouter()
	go eth1.WatchGasPrice()

	log.Info().Msg("starting bsc observer...")
	bsc1, err := mc.NewChainObserver(common.BSCChain, bridge, tss, dbpath)
	if err != nil {
		log.Err(err).Msg("BSC NewChainObserver")
		return nil, err
	}
	clientMap[common.BSCChain] = bsc1
	go bsc1.WatchRouter()
	go bsc1.WatchGasPrice()

	log.Info().Msg("starting polygon observer...")
	poly1, err := mc.NewChainObserver(common.POLYGONChain, bridge, tss, dbpath)
	if err != nil {
		log.Err(err).Msg("POLYGON NewChainObserver")
		return nil, err
	}
	clientMap[common.POLYGONChain] = poly1
	go poly1.WatchRouter()
	go poly1.WatchGasPrice()

	return &clientMap, nil
}
