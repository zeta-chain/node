package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/cmd"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metrics2 "github.com/zeta-chain/zetacore/zetaclient/metrics"

	//mcconfig "github.com/Meta-Protocol/zetacore/metaclient/config"
	"github.com/cosmos/cosmos-sdk/types"
	//"github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p-peerstore/addr"
	maddr "github.com/multiformats/go-multiaddr"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const ()

func main() {
	var notss = flag.Bool("notss", false, "use fake TSS")
	var validatorName = flag.String("val", "alice", "validator name")
	var tssTestFlag = flag.Bool("tss", false, "2 node TSS test mode")
	var peer = flag.String("peer", "", "peer address, e.g. /dns/tss1/tcp/6668/ipfs/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp")
	flag.Parse()
	//BOOTSTRAP_PEER := "/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ"

	var peers addr.AddrList
	fmt.Println("peer", *peer)
	if *peer != "" {
		address, err := maddr.NewMultiaddr(*peer)
		if err != nil {
			log.Error().Err(err).Msg("NewMultiaddr error")
			return
		}
		peers = append(peers, address)
	}

	if *tssTestFlag {
		SetupConfigForTest()
		fmt.Println("testing TSS signing")

		tssServer, _, err := mc.SetupTSSServer(peers, "")
		if err != nil {
			log.Error().Err(err).Msg("setup TSS server error")
			return
		}

		time.Sleep(5 * time.Second)
		kgRes := mc.TestKeygen(tssServer)
		log.Debug().Msgf("keygen succeeds! TSS pubkey: %s", kgRes.PubKey)

		log.Debug().Msgf("Keysign test begins...")

		mc.TestKeysign(kgRes.PubKey, tssServer)

		// wait....
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Info().Msg("stop signal received")
		return
	}

	if *notss {
		fmt.Println("fake TSS mode")
		integration_test_notss(*validatorName, peers)
		return
	} else {
		fmt.Println("multi-node client")
		integration_test(*validatorName, peers)
		return
	}

}

func SetupConfigForTest() {
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount(cmd.Bech32PrefixAccAddr, cmd.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(cmd.Bech32PrefixValAddr, cmd.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(cmd.Bech32PrefixConsAddr, cmd.Bech32PrefixConsPub)
	//config.SetCoinType(cmd.MetaChainCoinType)
	config.SetFullFundraiserPath(cmd.ZetaChainHDPath)
	types.SetCoinDenomRegex(func() string {
		return cmd.DenomRegex
	})

	rand.Seed(time.Now().UnixNano())

}

func integration_test(validatorName string, peers addr.AddrList) {
	SetupConfigForTest() // setup meta-prefix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	chainIP := os.Getenv("CHAIN_IP")
	if chainIP == "" {
		chainIP = "127.0.0.1"
	}

	// wait until zetacore is up
	log.Info().Msg("Waiting for ZetaCore to open 9090 port...")
	for {
		_, err := grpc.Dial(
			fmt.Sprintf("%s:9090", chainIP),
			grpc.WithInsecure(),
		)
		if err != nil {
			log.Warn().Err(err).Msg("grpc dial fail")
			time.Sleep(3 * time.Second)
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
	chainHomeFoler := filepath.Join(homeDir, ".zetacore")

	// first signer & bridge
	signerName := validatorName
	signerPass := "password"
	bridge1, done := CreateMetaBridge(chainHomeFoler, signerName, signerPass)
	if done {
		return
	}

	tss, err := mc.NewTSS(peers)
	if err != nil {
		log.Error().Err(err).Msg("NewTSS error")
		return
	}

	signerMap1, err := CreateSignerMap(tss)
	if err != nil {
		log.Err(err).Msg("CreateSignerMap")
		return
	}

	userDir, _ := os.UserHomeDir()
	dbpath := filepath.Join(userDir, ".zetaclient/chainobserver")
	chainClientMap1, err := CreateChainClientMap(bridge1, tss, dbpath)
	if err != nil {
		log.Err(err).Msg("CreateSignerMap")
		return
	}

	metrics, err := metrics2.NewMetrics()
	if err != nil {
		log.Error().Err(err).Msg("NewMetric")
		return
	}
	metrics.Start()

	log.Info().Msg("starting zetacore observer...")
	mo1 := mc.NewCoreObserver(bridge1, signerMap1, *chainClientMap1, metrics)

	mo1.MonitorCore()

	// do a test keysign when block reaches 30
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			bn, _ := bridge1.GetMetaBlockHeight()
			if bn == 30 {
				mc.TestKeysign(tss.PubkeyInBech32, tss.Server)
			}
		}
	}()

	// printout debug info from SIGUSR1
	// trigger by $ kill -SIGUSR1 <PID of zetaclient>
	usr := make(chan os.Signal, 1)
	signal.Notify(usr, syscall.SIGUSR1)
	go func() {
		for {
			<-usr
			fmt.Printf("Last blocks:\n")
			fmt.Printf("ETH     %d:\n", (*chainClientMap1)[common.ETHChain].LastBlock)
			fmt.Printf("BSC     %d:\n", (*chainClientMap1)[common.BSCChain].LastBlock)
			fmt.Printf("POLYGON %d:\n", (*chainClientMap1)[common.POLYGONChain].LastBlock)

		}
	}()

	// wait....
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("stop signal received")

	_ = (*chainClientMap1)[common.ETHChain].Stop()
	_ = (*chainClientMap1)[common.BSCChain].Stop()
	_ = (*chainClientMap1)[common.POLYGONChain].Stop()
}

func integration_test_notss(validatorName string, peers addr.AddrList) {
	SetupConfigForTest() // setup meta-prefix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("Fake TSS")

	// Configure MPI contract addresses
	if ethMPIAddress := os.Getenv("ETH_MPI_ADDRESS"); ethMPIAddress != "" {
		config.ETH_MPI_ADDRESS = ethMPIAddress
	}
	if bscMPIAddress := os.Getenv("BSC_MPI_ADDRESS"); bscMPIAddress != "" {
		config.BSC_MPI_ADDRESS = bscMPIAddress
	}
	if polygonMPIAddress := os.Getenv("POLYGON_MPI_ADDRESS"); polygonMPIAddress != "" {
		config.POLYGON_MPI_ADDRESS = polygonMPIAddress
	}

	chainIP := os.Getenv("CHAIN_IP")
	if chainIP == "" {
		chainIP = "127.0.0.1"
	}
	// wait until zetacore is up
	log.Info().Msg("Waiting for ZetaCore to open 9090 port...")
	for {
		_, err := grpc.Dial(
			fmt.Sprintf("%s:9090", chainIP),
			grpc.WithInsecure(),
		)
		if err != nil {
			log.Warn().Err(err).Msg("grpc dial fail")
			time.Sleep(3 * time.Second)
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
	chainHomeFoler := filepath.Join(homeDir, ".zetacore")

	// first signer & bridge
	signerName := validatorName
	signerPass := "password"
	bridge1, done := CreateMetaBridge(chainHomeFoler, signerName, signerPass)
	if done {
		return
	}

	// setup mock TSS signers:
	tss := GetZetaTestSignature()

	signerMap1, err := CreateSignerMap(tss)
	if err != nil {
		log.Err(err).Msg("CreateSignerMap")
		return
	}

	userDir, _ := os.UserHomeDir()
	dbpath := filepath.Join(userDir, ".zetaclient/chainobserver")
	chainClientMap1, err := CreateChainClientMap(bridge1, tss, dbpath)
	if err != nil {
		log.Err(err).Msg("CreateSignerMap")
		return
	}

	metrics, err := metrics2.NewMetrics()
	if err != nil {
		log.Error().Err(err).Msg("NewMetric")
		return
	}
	metrics.Start()

	log.Info().Msg("starting zetacore observer...")
	mo1 := mc.NewCoreObserver(bridge1, signerMap1, *chainClientMap1, metrics)

	mo1.MonitorCore()

	// printout debug info from SIGUSR1
	// trigger by $ kill -SIGUSR1 <PID of zetaclient>
	usr := make(chan os.Signal, 1)
	signal.Notify(usr, syscall.SIGUSR1)
	go func() {
		for {
			<-usr
			fmt.Printf("Last blocks:\n")
			fmt.Printf("ETH     %d:\n", (*chainClientMap1)[common.ETHChain].LastBlock)
			fmt.Printf("BSC     %d:\n", (*chainClientMap1)[common.BSCChain].LastBlock)
			fmt.Printf("POLYGON %d:\n", (*chainClientMap1)[common.POLYGONChain].LastBlock)

		}
	}()

	// wait....
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("stop signal received")

	_ = (*chainClientMap1)[common.ETHChain].Stop()
	_ = (*chainClientMap1)[common.BSCChain].Stop()
	_ = (*chainClientMap1)[common.POLYGONChain].Stop()
}
