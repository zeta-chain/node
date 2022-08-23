package main

import (
	"encoding/json"
	"flag"
	"fmt"
	ecdsakeygen "github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/rs/zerolog"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/zeta-chain/zetacore/cmd"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metrics2 "github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"

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

var (
	preParams   *ecdsakeygen.LocalPreParams
	keygenBlock int64
)

func main() {
	fmt.Printf("zeta-node commit hash %s version %s build time %s \n", common.CommitHash, common.Version, common.BuildTime)
	enabledChains := flag.String("enable-chains", "GOERLI,BSCTESTNET,MUMBAI,ROPSTEN", "enable chains, comma separated list")
	valKeyName := flag.String("val", "alice", "validator name")
	peer := flag.String("peer", "", "peer address, e.g. /dns/tss1/tcp/6668/ipfs/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp")
	logConsole := flag.Bool("log-console", false, "log to console (pretty print)")
	preParamsPath := flag.String("pre-params", "", "pre-params file path")
	keygen := flag.Int64("keygen-block", 0, "keygen at block height (default: 0 means no keygen)")

	flag.Parse()
	keygenBlock = *keygen
	if *logConsole {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

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

	updateEndpoint(common.GoerliChain, "GOERLI_ENDPOINT")
	updateEndpoint(common.BSCTestnetChain, "BSCTESTNET_ENDPOINT")
	updateEndpoint(common.MumbaiChain, "MUMBAI_ENDPOINT")
	updateEndpoint(common.RopstenChain, "ROPSTEN_ENDPOINT")

	updateMPIAddress(common.GoerliChain, "GOERLI_MPI_ADDRESS")
	updateMPIAddress(common.BSCTestnetChain, "BSCTESTNET_MPI_ADDRESS")
	updateMPIAddress(common.MumbaiChain, "MUMBAI_MPI_ADDRESS")
	updateMPIAddress(common.RopstenChain, "ROPSTEN_MPI_ADDRESS")

	// pools
	updatePoolAddress("GOERLI_POOL_ADDRESS", common.GoerliChain)
	updatePoolAddress("MUMBAI_POOL_ADDRESS", common.MumbaiChain)
	updatePoolAddress("BSCTESTNET_POOL_ADDRESS", common.BSCTestnetChain)
	updatePoolAddress("ROPSTEN_POOL_ADDRESS", common.RopstenChain)

	updateTokenAddress(common.GoerliChain, "GOERLI_ZETA_ADDRESS")
	updateTokenAddress(common.BSCTestnetChain, "BSCTESTNET_ZETA_ADDRESS")
	updateTokenAddress(common.MumbaiChain, "MUMBAI_ZETA_ADDRESS")
	updateTokenAddress(common.RopstenChain, "ROPSTEN_ZETA_ADDRESS")

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

	log.Info().Msg("wait for 20s for all node to SetNodeKey")
	time.Sleep(20 * time.Second)

	if keygenBlock > 0 {
		log.Info().Msgf("Keygen at blocknum %d", keygenBlock)
		bn, err := bridge1.GetZetaBlockHeight()
		if err != nil {
			log.Error().Err(err).Msg("GetZetaBlockHeight error")
			return
		}
		if int64(bn)+10 > keygenBlock {
			log.Warn().Msgf("Keygen at blocknum %d, but current blocknum %d", keygenBlock, bn)
			return
		}
		nodeAccounts, err := bridge1.GetAllNodeAccounts()
		if err != nil {
			log.Error().Err(err).Msg("GetAllNodeAccounts error")
			return
		}
		pubkeys := make([]string, 0)
		for _, na := range nodeAccounts {
			pubkeys = append(pubkeys, na.PubkeySet.Secp256k1.String())
		}
		ticker := time.NewTicker(time.Second * 2)
		for range ticker.C {
			bn, err := bridge1.GetZetaBlockHeight()
			if err != nil {
				log.Error().Err(err).Msg("GetZetaBlockHeight error")
				return
			}
			if int64(bn) == keygenBlock {
				break
			}
		}
		log.Info().Msgf("Keygen with %d TSS signers", len(nodeAccounts))
		log.Info().Msgf("%s", pubkeys)
		var req keygen.Request
		req = keygen.NewRequest(pubkeys, keygenBlock, "0.14.0")
		res, err := tss.Server.Keygen(req)
		if err != nil || res.Status != tsscommon.Success {
			log.Error().Msgf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)
			return
		}
		// Keygen succeed! Report TSS address
		log.Info().Msgf("Keygen success! keygen response: %v...", res)

		log.Info().Msgf("doing a keysign test...")
		err = mc.TestKeysign(res.PubKey, tss.Server)
		if err != nil {
			log.Error().Err(err).Msg("TestKeysign error")
			return
		}

		log.Info().Msgf("setting TSS pubkey: %s", res.PubKey)
		err = tss.SetPubKey(res.PubKey)
		if err != nil {
			log.Error().Msgf("SetPubKey fail")
			return
		}
		log.Info().Msgf("TSS address in hex: %s", tss.AddressInHex)
		return
	}

	kg, err := bridge1.GetKeyGen()
	if err != nil {
		log.Error().Err(err).Msg("GetKeyGen error")
		return
	}
	log.Info().Msgf("Setting TSS pubkeys: %s", kg.Pubkeys)
	tss.Pubkeys = kg.Pubkeys

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

func updatePoolAddress(envvar string, chain common.Chain) {
	pool := os.Getenv(envvar)
	parts := strings.Split(pool, ":")
	if len(parts) != 3 {
		log.Error().Msgf("%s is not a valid type:address", pool)
		return
	}
	if strings.EqualFold(parts[0], "v2") {
		config.Chains[chain.String()].PoolContract = clienttypes.UniswapV2
	} else if strings.EqualFold(parts[0], "v3") {
		config.Chains[chain.String()].PoolContract = clienttypes.UniswapV3
	} else {
		log.Error().Msgf("%s is not a valid type:address", pool)
		return
	}
	if parts[1] != "" {
		config.Chains[chain.String()].PoolContractAddress = parts[1]
		log.Info().Msgf("Pool address ENVVAR: %s: %s", envvar, parts[1])
	} else {
		log.Info().Msgf("Pool address DEFAULT: %s", config.Chains[chain.String()].PoolContractAddress)
	}
	if strings.EqualFold(parts[2], "ZETAETH") {
		config.Chains[chain.String()].PoolTokenOrder = clienttypes.ZETAETH
	} else if strings.EqualFold(parts[2], "ETHZETA") {
		config.Chains[chain.String()].PoolTokenOrder = clienttypes.ETHZETA
	} else {
		log.Error().Msgf("%s is not a valid type:address:order", pool)
		return
	}
}

func updateMPIAddress(chain common.Chain, envvar string) {
	mpi := os.Getenv(envvar)
	if mpi != "" {
		config.Chains[chain.String()].ConnectorContractAddress = mpi
		log.Info().Msgf("MPI: %s", mpi)
	}
}

func updateEndpoint(chain common.Chain, envvar string) {
	endpoint := os.Getenv(envvar)
	if endpoint != "" {
		config.Chains[chain.String()].Endpoint = endpoint
		log.Info().Msgf("ENDPOINT: %s", endpoint)
	}
}

func updateTokenAddress(chain common.Chain, envvar string) {
	token := os.Getenv(envvar)
	if token != "" {
		config.Chains[chain.String()].ZETATokenContractAddress = token
		log.Info().Msgf("TOKEN: %s", token)
	}
}
