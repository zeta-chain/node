package main

import (
	"github.com/Meta-Protocol/metacore/cmd"
	"github.com/Meta-Protocol/metacore/common"
	"github.com/Meta-Protocol/metacore/common/cosmos"
	mc "github.com/Meta-Protocol/metacore/metaclient"
	mcconfig "github.com/Meta-Protocol/metacore/metaclient/config"
	"github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	mock_integration_test()
}

func SetupConfigForTest() {
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount(cmd.Bech32PrefixAccAddr, cmd.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(cmd.Bech32PrefixValAddr, cmd.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(cmd.Bech32PrefixConsAddr, cmd.Bech32PrefixConsPub)
	//config.SetCoinType(cmd.MetaChainCoinType)
	config.SetFullFundraiserPath(cmd.METAChainHDPath)
	types.SetCoinDenomRegex(func() string {
		return cmd.DenomRegex
	})
}

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

	log.Info().Msg("starting metacore observer...")
	mo1 := mc.NewCoreObserver(bridge1, signerMap1)
	mo1.MonitorCore()
	mo2 := mc.NewCoreObserver(bridge2, signerMap2)
	mo2.MonitorCore()

	//{
	//	metaHash, err := bridge1.PostSend("0xfrom", "Ethereum", "0xto", "BSC", "123456", "23245", "little message",
	//		"0xtxhash", 123123)
	//	log.Info().Msgf("PostSend metaHash %s err %v", metaHash, err)
	//
	//	// wait for the next block
	//	timer1 := time.NewTimer(2 * time.Second)
	//	<-timer1.C
	//
	//	metaHash, err = bridge2.PostSend("0xfrom", "Ethereum", "0xto", "BSC", "123456", "23245", "little message",
	//		"0xtxhash", 123123)
	//	log.Info().Msgf("Second PostSend metaHash %s", metaHash)
	//
	//	// wait for the next block
	//	timer2 := time.NewTimer(2 * time.Second)
	//	<-timer2.C
	//}

	log.Info().Msg("starting eth observer...")
	eth1, _ := mc.NewChainObserver("Ethereum", bridge1)
	go eth1.WatchRouter()
	eth2, _ := mc.NewChainObserver("Ethereum", bridge2)
	go eth2.WatchRouter()

	log.Info().Msg("starting bsc observer...")
	bsc1, _ := mc.NewChainObserver("BSC", bridge1)
	go bsc1.WatchRouter()
	bsc2, _ := mc.NewChainObserver("BSC", bridge2)
	go bsc2.WatchRouter()

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
	metaContractAddress := ethcommon.HexToAddress(mcconfig.META_TEST_GOERLI_ADDRESS)
	ethSigner, err := mc.NewSigner(common.ETHChain, mcconfig.GOERLI_RPC_ENDPOINT, tss.Address(), tss, mcconfig.META_TEST_GOERLI_ABI, metaContractAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner Ethereum error ")
		return nil, err
	}
	bscSigner, err := mc.NewSigner(common.BSCChain, mcconfig.BSC_ENDPOINT, tss.Address(), tss, mcconfig.META_ABI, ethcommon.HexToAddress(mcconfig.BSC_ROUTER))
	if err != nil {
		log.Fatal().Err(err).Msg("NewSigner BSC error")
		return nil, err
	}
	signerMap := map[common.Chain]*mc.Signer{
		common.ETHChain: ethSigner,
		common.BSCChain: bscSigner,
	}

	return signerMap, nil
}
