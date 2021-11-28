package main

import (
	"github.com/Meta-Protocol/metacore/common"
	mc "github.com/Meta-Protocol/metacore/metaclient"
	mcconfig "github.com/Meta-Protocol/metacore/metaclient/config"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func mock_integration_test() {
	SetupConfigForTest() // setup meta-prefix

	// setup 2 metabridges
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Err(err).Msg("UserHomeDir error")
		return
	}
	chainHomeFoler := filepath.Join(homeDir, ".metacore")

	// first signer & bridge
	signerName := "alice"
	signerPass := "password"
	bridge1, done := CreateMetaBridge(chainHomeFoler, signerName, signerPass)
	if done {
		return
	}

	signerName = "bob"
	signerPass = "password"
	bridge2, done := CreateMetaBridge(chainHomeFoler, signerName, signerPass)
	if done {
		return
	}

	// setup mock TSS signers:
	// The following privkey has address 0xE80B6467863EbF8865092544f441da8fD3cF6074
	privateKey, err := crypto.HexToECDSA(mcconfig.TSS_TEST_PRIVKEY)
	if err != nil {
		log.Err(err).Msg("TEST private key error")
		return
	}
	tss := mc.TestSigner{
		PrivKey: privateKey,
	}

	signerMap1, err := CreateSignerMap(tss)
	if err != nil {
		log.Err(err).Msg("CreateSignerMap")
		return
	}
	signerMap2, err := CreateSignerMap(tss)
	if err != nil {
		log.Err(err).Msg("CreateSignerMap")
		return
	}

	chainClientMap1, err := CreateChainClientMap(bridge1, tss)
	if err != nil {
		log.Err(err).Msg("CreateSignerMap")
		return
	}
	chainClientMap2, err := CreateChainClientMap(bridge2, tss)
	if err != nil {
		log.Err(err).Msg("CreateSignerMap")
		return
	}

	log.Info().Msg("starting metacore observer...")
	mo1 := mc.NewCoreObserver(bridge1, signerMap1, *chainClientMap1)
	mo1.MonitorCore()
	mo2 := mc.NewCoreObserver(bridge2, signerMap2, *chainClientMap2)
	mo2.MonitorCore()

	// wait....
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("stop signal received")

}

func CreateMetaBridge(chainHomeFoler string, signerName string, signerPass string) (*mc.MetachainBridge, bool) {
	kb, _, err := mc.GetKeyringKeybase(chainHomeFoler, signerName, signerPass)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to get keyring keybase")
		return nil, true
	}

	k := mc.NewKeysWithKeybase(kb, signerName, signerPass)

	chainIP := "127.0.0.1"
	bridge, err := mc.NewMetachainBridge(k, chainIP, signerName)
	if err != nil {
		log.Fatal().Err(err).Msg("NewMetachainBridge")
		return nil, true
	}
	return bridge, false
}

func CreateSignerMap(tss mc.TSSSigner) (map[common.Chain]*mc.Signer, error) {
	ethSigner, err := mc.NewSigner(common.ETHChain, mcconfig.ETH_ENDPOINT, tss.Address(), tss, mcconfig.ETH_ZETA_LOCK_ABI, ethcommon.HexToAddress(mcconfig.ETH_METALOCK_ADDRESS))
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner Ethereum error ")
		return nil, err
	}
	bscSigner, err := mc.NewSigner(common.BSCChain, mcconfig.BSC_ENDPOINT, tss.Address(), tss, mcconfig.BSC_ZETA_ABI, ethcommon.HexToAddress(mcconfig.BSC_TOKEN_ADDRESS))
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner BSC error")
		return nil, err
	}
	polygonSigner, err := mc.NewSigner(common.POLYGONChain, mcconfig.POLY_ENDPOINT, tss.Address(), tss, mcconfig.BSC_ZETA_ABI, ethcommon.HexToAddress(mcconfig.POLYGON_TOKEN_ADDRESS))
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner POLYGON error")
		return nil, err
	}
	signerMap := map[common.Chain]*mc.Signer{
		common.ETHChain: ethSigner,
		common.BSCChain: bscSigner,
		common.POLYGONChain: polygonSigner,
	}

	return signerMap, nil
}

func CreateChainClientMap(bridge *mc.MetachainBridge, tss mc.TSSSigner) (*map[common.Chain]*mc.ChainObserver, error){
	log.Info().Msg("starting eth observer...")
	clientMap := make(map[common.Chain]*mc.ChainObserver)
	eth1, err := mc.NewChainObserver(common.ETHChain, bridge, tss.Address())
	if err != nil {
		log.Err(err).Msg("ETH NewChainObserver")
		return nil, err
	}
	clientMap[common.ETHChain] = eth1
	go eth1.WatchRouter()
	go eth1.WatchGasPrice()


	log.Info().Msg("starting bsc observer...")
	bsc1, err := mc.NewChainObserver(common.BSCChain, bridge, tss.Address())
	if err != nil {
		log.Err(err).Msg("BSC NewChainObserver")
		return nil, err
	}
	clientMap[common.BSCChain] = bsc1
	go bsc1.WatchRouter()
	go bsc1.WatchGasPrice()



	log.Info().Msg("starting polygon observer...")
	poly1, err := mc.NewChainObserver(common.POLYGONChain, bridge, tss.Address())
	if err != nil {
		log.Err(err).Msg("POLYGON NewChainObserver")
		return nil, err
	}
	clientMap[common.POLYGONChain] = poly1
	go poly1.WatchRouter()
	go poly1.WatchGasPrice()

	return &clientMap, nil
}