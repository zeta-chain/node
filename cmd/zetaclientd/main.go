package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/tendermint/tendermint/crypto/secp256k1"
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
	var validatorName = flag.String("val", "alice", "validator name")
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

	fmt.Println("multi-node client")
	integration_test(*validatorName, peers)
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

	ethEndPoint := os.Getenv("ETH_ENDPOINT")
	if ethEndPoint != "" {
		config.ETH_ENDPOINT = ethEndPoint
		log.Info().Msgf("ETH_ENDPOINT: %s", ethEndPoint)
	}
	bscEndPoint := os.Getenv("BSC_ENDPOINT")
	if bscEndPoint != "" {
		config.BSC_ENDPOINT = bscEndPoint
		log.Info().Msgf("BSC_ENDPOINT: %s", bscEndPoint)
	}
	polygonEndPoint := os.Getenv("POLYGON_ENDPOINT")
	if polygonEndPoint != "" {
		config.POLY_ENDPOINT = polygonEndPoint
		log.Info().Msgf("POLYGON_ENDPOINT: %s", polygonEndPoint)
	}

	ethMpiAddress := os.Getenv("ETH_MPI_ADDRESS")
	if ethMpiAddress != "" {
		config.Chains["ETH"].MPIContractAddress = ethMpiAddress
		log.Info().Msgf("ETH_MPI_ADDRESS: %s", ethMpiAddress)
	}
	bscMpiAddress := os.Getenv("BSC_MPI_ADDRESS")
	if bscMpiAddress != "" {
		config.Chains["BSC"].MPIContractAddress = bscMpiAddress
		log.Info().Msgf("BSC_MPI_ADDRESS: %s", bscMpiAddress)
	}
	polygonMpiAddress := os.Getenv("POLYGON_MPI_ADDRESS")
	if polygonMpiAddress != "" {
		config.Chains["POLYGON"].MPIContractAddress = polygonMpiAddress
		log.Info().Msgf("polygonMpiAddress: %s", polygonMpiAddress)
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
			time.Sleep(5 * time.Second)
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

	key, err := bridge1.GetKeys().GetPrivateKey()
	if err != nil {
		log.Error().Err(err).Msg("GetKeys GetPrivateKey error:")
	}
	if len(key.Bytes()) != 32 {
		log.Error().Msgf("key bytes len %d != 32", len(key.Bytes()))
		return
	}
	var priKey secp256k1.PrivKey
	priKey = key.Bytes()[:32]

	log.Info().Msgf("NewTSS: with peer pubkey %s", key.PubKey())
	tss, err := mc.NewTSS(peers, priKey)
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
	mo1 := mc.NewCoreObserver(bridge1, signerMap1, *chainClientMap1, metrics, tss)

	mo1.MonitorCore()

	// report node key
	// convert key.PubKey() [cosmos-sdk/crypto/PubKey] to bech32?
	s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, key.PubKey())
	if err != nil {
		log.Error().Err(err).Msgf("Bech32ifyPubKey fail in main")
	}
	log.Info().Msgf("GetPrivateKey for pubkey bech32 %s", s)

	pubkey, err := common.NewPubKey(s)
	if err != nil {
		log.Error().Err(err).Msgf("NewPubKey error from string %s:", key.PubKey().String())
	}
	pubkeyset := common.PubKeySet{
		Secp256k1: pubkey,
		Ed25519:   "",
	}
	conskey := ""
	ztx, err := bridge1.SetNodeKey(pubkeyset, conskey)
	log.Info().Msgf("SetNodeKey: %s by node %s zeta tx %s", pubkeyset.Secp256k1.String(), conskey, ztx)
	if err != nil {
		log.Error().Err(err).Msgf("SetNodeKey error")
	}

	// report TSS address nonce on ETHish chains
	err = (*chainClientMap1)[common.ETHChain].PostNonceIfNotRecorded()
	if err != nil {
		log.Error().Err(err).Msgf("PostNonceIfNotRecorded fail %s", common.ETHChain)
	}
	err = (*chainClientMap1)[common.BSCChain].PostNonceIfNotRecorded()
	if err != nil {
		log.Error().Err(err).Msgf("PostNonceIfNotRecorded fail %s", common.BSCChain)
	}
	err = (*chainClientMap1)[common.POLYGONChain].PostNonceIfNotRecorded()
	if err != nil {
		log.Error().Err(err).Msgf("PostNonceIfNotRecorded fail %s", common.POLYGONChain)
	}

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
