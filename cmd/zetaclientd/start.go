package main

import (
	"encoding/json"
	"fmt"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/zeta-chain/zetacore/common"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metrics2 "github.com/zeta-chain/zetacore/zetaclient/metrics"
	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"
	"gitlab.com/thorchain/tss/go-tss/p2p"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start ZetaClient Observer",
	RunE:  start,
}

const maxRetryCount = 10

func init() {
	RootCmd.AddCommand(StartCmd)
}

func start(_ *cobra.Command, _ []string) error {
	setHomeDir()
	SetupConfigForTest()
	//Load Config file given path
	configData, err := config.Load(rootArgs.zetaCoreHome)
	if err != nil {
		return err
	}
	log.Logger = InitLogger(configData.LogLevel)
	//Wait until zetacore has started
	waitForZetaCore(configData)
	masterLogger := log.Logger
	startLogger := masterLogger.With().Str("module", "startup").Logger()
	startLogger.Info().Msgf("ZetaCore is ready")
	// first signer & bridge

	bridge1, err := CreateZetaBridge(rootArgs.zetaCoreHome, configData)
	if err != nil {
		panic(err)
	}
	startLogger.Info().Msgf("ZetaBridge is ready")

	bridge1.SetAccountNumber(common.ZetaClientGranteeKey)
	CreateAuthzSigner(bridge1.GetKeys().GetOperatorAddress().String(),
		bridge1.GetKeys().GetAddress())

	startLogger.Debug().Msgf("CreateAuthzSigner is ready")

	bridgePk, err := bridge1.GetKeys().GetPrivateKey()
	if err != nil {
		startLogger.Error().Err(err).Msg("GetKeys GetPrivateKey error:")
	}

	startLogger.Debug().Msgf("bridgePk %s", bridgePk.String())
	if len(bridgePk.Bytes()) != 32 {
		errMsg := fmt.Sprintf("key bytes len %d != 32", len(bridgePk.Bytes()))
		log.Error().Msgf(errMsg)
		return errors.New(errMsg)
	}
	var priKey secp256k1.PrivKey
	priKey = bridgePk.Bytes()[:32]

	startLogger.Debug().Msgf("NewTSS: with peer pubkey %s", bridgePk.PubKey())

	peers, err := initPeers(configData.Peer)
	if err != nil {
		log.Error().Err(err).Msg("peer address error")
	}
	initPreParams(configData.PreParamsPath)
	tss, err := mc.NewTSS(peers, priKey, preParams)
	if err != nil {
		startLogger.Error().Err(err).Msg("NewTSS error")
		return err
	}
	retryCount := 0
	for {
		block, err := bridge1.GetLatestZetaBlock()
		if err == nil && block.Header.Height > 1 {
			startLogger.Info().Msgf("Zeta-core height: %d", block.Header.Height)
			break
		}
		retryCount++
		startLogger.Debug().Msgf("Failed to get latest Block , Retry : %d/%d", retryCount, maxRetryCount)
		if retryCount > 10 {
			panic("ZetaCore is not ready , Waited for 60s")
		}
		time.Sleep(6 * time.Second)

	}

	//Check if keygen block is set and generate new keys at specified height
	genNewKeysAtBlock(configData.KeygenBlock, bridge1, tss)

	for _, chain := range config.ChainsEnabled {
		var tssAddr string
		if chain.IsEVMChain() {
			tssAddr = tss.EVMAddress().Hex()
		} else {
			tssAddr = tss.BTCAddress()
		}
		zetaTx, err := bridge1.SetTSS(chain, tssAddr, tss.CurrentPubkey)
		if err != nil {
			startLogger.Error().Err(err).Msgf("SetTSS fail %s", chain.String())
		}
		startLogger.Info().Msgf("chain %s set TSS to %s, zeta tx hash %s", chain.String(), tssAddr, zetaTx)

	}
	signerMap1, err := CreateSignerMap(tss, masterLogger)
	if err != nil {
		log.Error().Err(err).Msg("CreateSignerMap")
		return err
	}

	metrics, err := metrics2.NewMetrics()
	if err != nil {
		log.Error().Err(err).Msg("NewMetrics")
		return err
	}
	metrics.Start()

	userDir, _ := os.UserHomeDir()
	dbpath := filepath.Join(userDir, ".zetaclient/chainobserver")
	chainClientMap1, err := CreateChainClientMap(bridge1, tss, dbpath, metrics, masterLogger)
	if err != nil {
		startLogger.Err(err).Msg("CreateSignerMap")
		return err
	}
	for _, v := range chainClientMap1 {
		v.Start()
	}

	mo1 := mc.NewCoreObserver(bridge1, signerMap1, chainClientMap1, metrics, tss, masterLogger)

	mo1.MonitorCore()

	// report TSS address nonce on ETHish chains
	for _, chain := range config.ChainsEnabled {
		err = (chainClientMap1)[chain].PostNonceIfNotRecorded(startLogger)
		if err != nil {
			startLogger.Fatal().Err(err).Msgf("PostNonceIfNotRecorded fail %s", chain.String())
		}
	}

	startLogger.Info().Msgf("TSS address \n ETH : %s \n BTC : %s \n PubKey : %s ", tss.EVMAddress(), tss.BTCAddress(), tss.CurrentPubkey)
	startLogger.Info().Msgf("awaiting the os.Interrupt, syscall.SIGTERM signals...")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	startLogger.Info().Msgf("stop signal received: %s", sig)

	// stop zetacore observer
	for _, chain := range config.ChainsEnabled {
		(chainClientMap1)[chain].Stop()
	}

	return nil
}

func waitForZetaCore(configData *config.Config) {
	// wait until zetacore is up
	log.Debug().Msg("Waiting for ZetaCore to open 9090 port...")
	for {
		_, err := grpc.Dial(
			fmt.Sprintf("%s:9090", configData.ZetaCoreURL),
			grpc.WithInsecure(),
		)
		if err != nil {
			log.Warn().Err(err).Msg("grpc dial fail")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}

func initPeers(peer string) (p2p.AddrList, error) {
	var peers p2p.AddrList

	if peer != "" {
		address, err := maddr.NewMultiaddr(peer)
		if err != nil {
			log.Error().Err(err).Msg("NewMultiaddr error")
			return p2p.AddrList{}, err
		}
		peers = append(peers, address)
	}
	return peers, nil
}

func initPreParams(path string) {
	if path != "" {
		path = filepath.Clean(path)
		log.Info().Msgf("pre-params file path %s", path)
		preParamsFile, err := os.Open(path)
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
}

func genNewKeysAtBlock(height int64, bridge *mc.ZetaCoreBridge, tss *mc.TSS) {
	if height > 0 {
		log.Info().Msgf("Keygen at blocknum %d", height)
		bn, err := bridge.GetZetaBlockHeight()
		if err != nil {
			log.Error().Err(err).Msg("GetZetaBlockHeight error")
			return
		}
		if bn+3 > height {
			log.Fatal().Msgf("Keygen at Blocknum %d, but current blocknum %d , Too late to take part in this keygen. Try again at a later block", height, bn)
			return
		}
		nodeAccounts, err := bridge.GetAllNodeAccounts()
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
			bn, err := bridge.GetZetaBlockHeight()
			if err != nil {
				log.Error().Err(err).Msg("GetZetaBlockHeight error")
				return
			}
			if bn == height {
				break
			}
			log.Debug().Msgf("Waiting for KeygenBlock %d, Current blocknum %d", height, bn)
		}
		log.Info().Msgf("Keygen with %d TSS signers", len(nodeAccounts))
		log.Info().Msgf("%s", pubkeys)
		var req keygen.Request
		req = keygen.NewRequest(pubkeys, height, "0.14.0")
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
		log.Info().Msgf("TSS address in hex: %s", tss.EVMAddress().Hex())
		return
	}
}
