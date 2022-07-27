package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/rs/zerolog"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/zeta-chain/zetacore/cmd"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metrics2 "github.com/zeta-chain/zetacore/zetaclient/metrics"
	"io/ioutil"
	"strings"
	"syscall"

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
	"time"
)

var preParams *keygen.LocalPreParams

func main() {
	fmt.Printf("zeta-node commit hash %s version %s build time %s \n", common.CommitHash, common.Version, common.BuildTime)
	enabledChains := flag.String("enable-chains", "GOERLI,BSCTESTNET,MUMBAI,ROPSTEN", "enable chains, comma separated list")
	valKeyName := flag.String("val", "alice", "validator name")
	peer := flag.String("peer", "", "peer address, e.g. /dns/tss1/tcp/6668/ipfs/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp")
	logConsole := flag.Bool("log-console", false, "log to console (pretty print)")
	preParamsPath := flag.String("pre-params", "", "pre-params file path")

	flag.Parse()
	chains := strings.Split(*enabledChains, ",")
	for _, chain := range chains {
		if c, err := common.ParseChain(chain); err == nil {
			config.ChainsEnabled = append(config.ChainsEnabled, c)
		} else {
			log.Error().Err(err).Msgf("invalid chain %s", chain)
			return
		}
	}
	log.Info().Msgf("enabled chains %v", config.ChainsEnabled)

	if *logConsole {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if *preParamsPath != "" {
		log.Info().Msgf("pre-params file path %s", *preParamsPath)
		preParamsFile, err := os.Open(*preParamsPath)
		if err != nil {
			log.Error().Err(err).Msg("open pre-params file failed; skip")
		} else {
			bz, err := ioutil.ReadAll(preParamsFile)
			if err != nil {
				log.Error().Err(err).Msg("read pre-params file failed; skip")
			} else {
				err = json.Unmarshal(bz, &preParams)
				if err != nil {
					log.Error().Err(err).Msg("unmarshal pre-params file failed; skip and generate new one")
					preParams = nil // skip reading pre-params; generate new one instead
				}
			}
		}
	}

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
	start(*valKeyName, peers)
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

func start(validatorName string, peers addr.AddrList) {
	SetupConfigForTest() // setup meta-prefix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	chainIP := os.Getenv("CHAIN_IP")
	if chainIP == "" {
		chainIP = "127.0.0.1"
	}

	ethEndPoint := os.Getenv("GOERLI_ENDPOINT")
	if ethEndPoint != "" {
		config.Chains[common.GoerliChain.String()].Endpoint = ethEndPoint
		log.Info().Msgf("GOERLI_ENDPOINT: %s", ethEndPoint)
	}
	bscEndPoint := os.Getenv("BSCTESTNET_ENDPOINT")
	if bscEndPoint != "" {
		config.Chains[common.BSCTestnetChain.String()].Endpoint = bscEndPoint
		log.Info().Msgf("BSCTESTNET_ENDPOINT: %s", bscEndPoint)
	}
	polygonEndPoint := os.Getenv("MUMBAI_ENDPOINT")
	if polygonEndPoint != "" {
		config.Chains[common.MumbaiChain.String()].Endpoint = polygonEndPoint
		log.Info().Msgf("MUMBAI_ENDPOINT: %s", polygonEndPoint)
	}
	ropstenEndPoint := os.Getenv("ROPSTEN_ENDPOINT")
	if ropstenEndPoint != "" {
		config.Chains[common.RopstenChain.String()].Endpoint = ropstenEndPoint
		log.Info().Msgf("ROPSTEN_ENDPOINT: %s", ropstenEndPoint)
	}

	ethMpiAddress := os.Getenv("GOERLI_MPI_ADDRESS")
	if ethMpiAddress != "" {
		config.Chains[common.GoerliChain.String()].ConnectorContractAddress = ethMpiAddress
		log.Info().Msgf("ETH_MPI_ADDRESS: %s", ethMpiAddress)
	}
	bscMpiAddress := os.Getenv("BSCTESTNET_MPI_ADDRESS")
	if bscMpiAddress != "" {
		config.Chains[common.BSCTestnetChain.String()].ConnectorContractAddress = bscMpiAddress
		log.Info().Msgf("BSC_MPI_ADDRESS: %s", bscMpiAddress)
	}
	polygonMpiAddress := os.Getenv("MUMBAI_MPI_ADDRESS")
	if polygonMpiAddress != "" {
		config.Chains[common.MumbaiChain.String()].ConnectorContractAddress = polygonMpiAddress
		log.Info().Msgf("polygonMpiAddress: %s", polygonMpiAddress)
	}
	ropstenMpiAddress := os.Getenv("ROPSTEN_MPI_ADDRESS")
	if ropstenMpiAddress != "" {
		config.Chains[common.RopstenChain.String()].ConnectorContractAddress = ropstenMpiAddress
		log.Info().Msgf("ropstenMpiAddress: %s", ropstenMpiAddress)
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
	bridge1, done := CreateZetaBridge(chainHomeFoler, signerName, signerPass)
	if done {
		return
	}

	bridgePk, err := bridge1.GetKeys().GetPrivateKey()
	if err != nil {
		log.Error().Err(err).Msg("GetKeys GetPrivateKey error:")
	}
	if len(bridgePk.Bytes()) != 32 {
		log.Error().Msgf("key bytes len %d != 32", len(bridgePk.Bytes()))
		return
	}
	var priKey secp256k1.PrivKey
	priKey = bridgePk.Bytes()[:32]

	log.Info().Msgf("NewTSS: with peer pubkey %s", bridgePk.PubKey())
	tss, err := mc.NewTSS(peers, priKey, preParams)
	if err != nil {
		log.Error().Err(err).Msg("NewTSS error")
		return
	}

	for _, chain := range config.ChainsEnabled {
		_, err = bridge1.SetTSS(chain, tss.Address().Hex(), tss.PubkeyInBech32)
		if err != nil {
			log.Error().Err(err).Msgf("SetTSS fail %s", chain)
		}
	}

	signerMap1, err := CreateSignerMap(tss)
	if err != nil {
		log.Error().Err(err).Msg("CreateSignerMap")
		return
	}

	metrics, err := metrics2.NewMetrics()
	if err != nil {
		log.Error().Err(err).Msg("NewMetrics")
		return
	}
	metrics.Start()

	userDir, _ := os.UserHomeDir()
	dbpath := filepath.Join(userDir, ".zetaclient/chainobserver")
	chainClientMap1, err := CreateChainClientMap(bridge1, tss, dbpath, metrics)
	if err != nil {
		log.Err(err).Msg("CreateSignerMap")
		return
	}
	for _, v := range *chainClientMap1 {
		v.Start()
	}

	log.Info().Msg("starting zetacore observer...")
	mo1 := mc.NewCoreObserver(bridge1, signerMap1, *chainClientMap1, metrics, tss)

	mo1.MonitorCore()

	consKey := ""
	pubkeySet, err := bridge1.GetKeys().GetPubKeySet()
	if err != nil {
		log.Error().Err(err).Msgf("Get Pubkey Set Error")
	}
	ztx, err := bridge1.SetNodeKey(pubkeySet, consKey)
	log.Info().Msgf("SetNodeKey: %s by node %s zeta tx %s", pubkeySet.Secp256k1.String(), consKey, ztx)
	if err != nil {
		log.Error().Err(err).Msgf("SetNodeKey error")
	}

	// report TSS address nonce on ETHish chains
	for _, chain := range config.ChainsEnabled {
		err = (*chainClientMap1)[chain].PostNonceIfNotRecorded()
		if err != nil {
			log.Error().Err(err).Msgf("PostNonceIfNotRecorded fail %s", chain)
		}
	}

	// wait....
	log.Info().Msgf("awaiting the os.Interrupt, syscall.SIGTERM signals...")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	log.Info().Msgf("stop signal received: %s", sig)

	// stop zetacore observer
	for _, chain := range config.ChainsEnabled {
		(*chainClientMap1)[chain].Stop()
	}

}
