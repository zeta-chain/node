package main

import (
	"fmt"
	"github.com/libp2p/go-libp2p-peerstore/addr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/zeta-chain/zetacore/common"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metrics2 "github.com/zeta-chain/zetacore/zetaclient/metrics"
	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func startTestMode(validatorName string, peers addr.AddrList, zetacoreHome string) {
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
	chainHomeFoler := filepath.Join(homeDir, zetacoreHome)

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
	time.Sleep(12 * time.Second)

	if keygenBlock > 0 {
		log.Info().Msgf("Keygen at blocknum %d", keygenBlock)
		bn, err := bridge1.GetZetaBlockHeight()
		if err != nil {
			log.Error().Err(err).Msg("GetZetaBlockHeight error")
			return
		}
		if int64(bn)+3 > keygenBlock {
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
		err = tss.InsertPubKey(res.PubKey)
		tss.CurrentPubkey = res.PubKey
		if err != nil {
			log.Error().Msgf("SetPubKey fail")
			return
		}
		log.Info().Msgf("TSS address in hex: %s", tss.Address().Hex())
		return
	}

	for _, chain := range config.ChainsEnabled {
		zetaTx, err := bridge1.SetTSS(chain, tss.Address().Hex(), tss.CurrentPubkey)
		if err != nil {
			log.Error().Err(err).Msgf("SetTSS fail %s", chain)
		}
		log.Info().Msgf("chain %s set TSS to %s, zeta tx hash %s", chain, tss.Address().Hex(), zetaTx)

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
