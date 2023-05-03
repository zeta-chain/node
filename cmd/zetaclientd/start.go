package main

import (
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p/core"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/zeta-chain/zetacore/common"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metrics2 "github.com/zeta-chain/zetacore/zetaclient/metrics"
	"gitlab.com/thorchain/tss/go-tss/p2p"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Multiaddr = core.Multiaddr

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start ZetaClient Observer",
	RunE:  start,
}

func init() {
	RootCmd.AddCommand(StartCmd)
}

func start(_ *cobra.Command, _ []string) error {
	setHomeDir()
	SetupConfigForTest()
	//Load Config file given path
	cfg, err := config.Load(rootArgs.zetaCoreHome)
	if err != nil {
		return err
	}

	log.Logger = InitLogger(cfg.LogLevel)
	//Wait until zetacore has started
	if len(cfg.Peer) != 0 {
		err := validatePeer(cfg.Peer)
		if err != nil {
			return err
		}
	}
	waitForZetaCore(cfg)
	masterLogger := log.Logger
	startLogger := masterLogger.With().Str("module", "startup").Logger()
	startLogger.Info().Msgf("ZetaCore is ready")
	// CreateZetaBridge:  Zetabridge is used for all communication to zetacore , which this client connects to.
	// Zetacore accumulates votes , and provides a centralized source of truth for all clients
	zetaBridge, err := CreateZetaBridge(rootArgs.zetaCoreHome, cfg)
	if err != nil {
		panic(err)
	}
	zetaBridge.WaitForCoreToCreateBlocks()
	startLogger.Info().Msgf("ZetaBridge is ready")
	zetaBridge.SetAccountNumber(common.ZetaClientGranteeKey)

	// CreateAuthzSigner : which is used to sign all authz messages . All votes broadcast to zetacore are wrapped in authz exec .
	// This is to ensure that the user does not need to keep their operator key online , and can use a cold key to sign votes
	CreateAuthzSigner(zetaBridge.GetKeys().GetOperatorAddress().String(), zetaBridge.GetKeys().GetAddress())
	startLogger.Debug().Msgf("CreateAuthzSigner is ready")

	// ConfigUpdater : This runs at every tick to check for configuration from zetacore and updates config if there are any changes . Zetacore stores configuration information which is common across all clients
	go zetaBridge.ConfigUpdater(cfg)
	time.Sleep((time.Duration(cfg.ConfigUpdateTicker) + 1) * time.Second)
	startLogger.Info().Msgf("Config is updated from ZetaCore %s", cfg.String())
	if len(cfg.ChainsEnabled) == 0 {
		startLogger.Info().Msgf("No chains enabled, exiting")
		return nil
	}

	// Generate TSS address . The Tss address is generated through Keygen ceremony. The TSS key is used to sign all outbound transactions .
	// Each node processes a portion of the key stored in ~/.tss by default . Custom location can be specified in config file during init.
	// After generating the key , the address is set on the zetacore
	bridgePk, err := zetaBridge.GetKeys().GetPrivateKey()
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

	// Generate pre Params if not present already
	peers, err := initPeers(cfg.Peer)
	if err != nil {
		log.Error().Err(err).Msg("peer address error")
	}
	initPreParams(cfg.PreParamsPath)
	if cfg.P2PDiagnostic {
		err := RunDiagnostics(startLogger, peers, bridgePk, cfg)
		if err != nil {
			startLogger.Error().Err(err).Msg("RunDiagnostics error")
			return err
		}
	}
	tss, err := mc.NewTSS(peers, priKey, preParams, cfg)
	if err != nil {
		startLogger.Error().Err(err).Msg("NewTSS error")
		return err
	}
	// If Keygen block is set it will try to generate new TSS at the block
	// This is a blocking thread and will wait until the ceremony is complete , and report weather it's a success or failure
	// Set TSS block to 0 using genesis file to disable this feature
	if cfg.KeygenBlock > 0 {
		err = genNewTSSAtBlock(cfg, zetaBridge, tss)
		if err != nil {
			startLogger.Error().Err(err).Msg("genNewTSSAtBlock error")
			return err
		}
	}
	startLogger.Info().Msgf("TSS address \n ETH : %s \n BTC : %s \n PubKey : %s ", tss.EVMAddress(), tss.BTCAddress(), tss.CurrentPubkey)

	// Vote Keygen success
	for _, chain := range cfg.ChainsEnabled {
		var tssAddr string
		if common.IsEVMChain(chain.ChainId) {
			tssAddr = tss.EVMAddress().Hex()
		} else if common.IsBitcoinChain(chain.ChainId) {
			tssAddr = tss.BTCAddress()
		}
		zetaTx, err := zetaBridge.SetTSS(chain, tssAddr, tss.CurrentPubkey)
		if err != nil {
			startLogger.Error().Err(err).Msgf("SetTSS fail %s", chain.String())
		}
		startLogger.Info().Msgf("chain %s set TSS to %s, zeta tx hash %s", chain.String(), tssAddr, zetaTx)
	}

	// CreateSignerMap : This creates a map of all signers for each chain . Each signer is responsible for signing transactions for a particular chain
	signerMap1, err := CreateSignerMap(tss, masterLogger, cfg)
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

	// CreateChainClientMap : This creates a map of all chain clients . Each chain client is responsible for listening to events on the chain and processing them
	chainClientMap, err := CreateChainClientMap(zetaBridge, tss, dbpath, metrics, masterLogger, cfg)
	if err != nil {
		startLogger.Err(err).Msg("CreateSignerMap")
		return err
	}
	for _, v := range chainClientMap {
		v.Start()
	}

	// CreateCoreObserver : Core observer wraps the zetacore bridge and adds the client and signer maps to it . This is the high level object used for CCTX interactions
	mo1 := mc.NewCoreObserver(zetaBridge, signerMap1, chainClientMap, metrics, tss, masterLogger, cfg)
	mo1.MonitorCore()

	// report TSS address nonce on all chains except zeta
	for _, chain := range cfg.ChainsEnabled {
		if chain.IsExternalChain() {
			err = (chainClientMap)[chain].PostNonceIfNotRecorded(startLogger)
			if err != nil {
				startLogger.Fatal().Err(err).Msgf("PostNonceIfNotRecorded fail %s", chain.String())
			}
		}
	}

	startLogger.Info().Msgf("awaiting the os.Interrupt, syscall.SIGTERM signals...")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	startLogger.Info().Msgf("stop signal received: %s", sig)

	// stop zetacore observer
	for _, chain := range cfg.ChainsEnabled {
		(chainClientMap)[chain].Stop()
	}
	zetaBridge.Stop()

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

func validatePeer(seedPeer string) error {
	parsedPeer := strings.Split(seedPeer, "/")

	if len(parsedPeer) < 7 {
		log.Error().Msgf("seed peer is malformed: %s", seedPeer)
		return errors.New("seed peer missing IP or ID")
	}

	seedIP := parsedPeer[2]
	seedID := parsedPeer[6]

	if net.ParseIP(seedIP) == nil {
		log.Error().Msgf("invalid seed IP address: %s", seedIP)
		return errors.New("invalid seed IP address")
	}

	if len(seedID) == 0 {
		log.Error().Msgf("seed id is empty")
		return errors.New("seed id is empty")
	}

	return nil
}
