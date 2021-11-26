package main

import (
	"flag"
	"fmt"
	"github.com/Meta-Protocol/metacore/cmd"
	"github.com/Meta-Protocol/metacore/common/cosmos"
	mc "github.com/Meta-Protocol/metacore/metaclient"
	mcconfig "github.com/Meta-Protocol/metacore/metaclient/config"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	var mockFlag = flag.Bool("mock", false, "mock 2 nodes environment")
	var validatorName = flag.String("val", "alice", "validator name")
	flag.Parse()

	if *mockFlag {
		fmt.Println("single node multiple clients tests")
		mock_integration_test() // single node testing environment; mocking multiple clients
		return
	} else {
		fmt.Println("multi-node client")
		integration_test(*validatorName)
		return
	}

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

	rand.Seed(time.Now().UnixNano())

}

func integration_test(validatorName string) {
	SetupConfigForTest() // setup meta-prefix

	// wait until metacore is up
	log.Info().Msg("Waiting for ZetaCore to open 9090 port...")
	for {
		_, err := grpc.Dial(
			fmt.Sprintf("127.0.0.1:9090"),
			grpc.WithInsecure(),
		)
		if err != nil {
			log.Warn().Err(err).Msg("grpc dial fail")
			time.Sleep(3*time.Second)
		} else {
			break
		}
	}
	log.Info().Msgf("ZetaCore to open 9090 port...")


	// setup 2 metabridges
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Err(err).Msg("UserHomeDir error")
		return
	}
	chainHomeFoler := filepath.Join(homeDir, ".metacore")

	// first signer & bridge
	signerName := validatorName
	signerPass := "password"
	bridge1, done := CreateMetaBridge(chainHomeFoler, signerName, signerPass)
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


	chainClientMap1, err := CreateChainClientMap(bridge1, tss)
	if err != nil {
		log.Err(err).Msg("CreateSignerMap")
		return
	}


	log.Info().Msg("starting metacore observer...")
	mo1 := mc.NewCoreObserver(bridge1, signerMap1, *chainClientMap1)
	mo1.MonitorCore()


	// wait....
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("stop signal received")
}
