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
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metrics2 "github.com/zeta-chain/zetacore/zetaclient/metrics"
	"gitlab.com/thorchain/tss/go-tss/p2p"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
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
	log.Logger = InitLogger(cfg)
	//Wait until zetacore has started
	if len(cfg.Peer) != 0 {
		err := validatePeer(cfg.Peer)
		if err != nil {
			log.Error().Err(err).Msg("invalid peer")
			return err
		}
	}

	masterLogger := log.Logger
	startLogger := masterLogger.With().Str("module", "startup").Logger()
	setMYIP(cfg, startLogger)
	waitForZetaCore(cfg, startLogger)
	startLogger.Info().Msgf("ZetaCore is ready , Trying to connect to %s", cfg.Peer)

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
	// This is a blocking thread and will wait until the ceremony is complete successfully
	// If the TSS generation is unsuccessful , it will loop indefinitely until a new TSS is generated
	// Set TSS block to 0 using genesis file to disable this feature
	// Note : The TSS generation is done through the "hotkey" or "Zeta-clientGrantee" This key needs to be present on the machine for the TSS signing to happen .
	// "ZetaClientGrantee" key is different from the "operator" key .The "Operator" key gives all zetaclient related permissions such as TSS generation ,reporting and signing, INBOUND and OUTBOUND vote signing, to the "ZetaClientGrantee" key.
	// The votes to signify a successful TSS generation(Or unsuccessful) is signed by the operator key and broadcast to zetacore by the zetcalientGrantee key on behalf of the operator .
	ticker := time.NewTicker(time.Second * 1)
	triedKeygenAtBlock := false
	lastBlock := int64(0)
	for range ticker.C {
		// Break out of loop only when TSS is generated successfully , either at the keygenBlock or if it has been generated already , Block set as zero in genesis file
		// This loop will try keygen at the keygen block and then wait for keygen to be successfully reported by all nodes before breaking out of the loop.
		// If keygen is unsuccessful , it will reset the triedKeygenAtBlock flag and try again at a new keygen block.
		if cfg.KeyGenStatus == crosschaintypes.KeygenStatus_KeyGenSuccess {
			break
		}
		// Arrive at this stage only if keygen is unsuccessfully reported by every node . This will reset the flag and to try again at a new keygen block
		if cfg.KeyGenStatus == crosschaintypes.KeygenStatus_KeyGenFailed {
			triedKeygenAtBlock = false
			continue
		}
		// Try generating TSS at keygen block , only when status is pending keygen and generation has not been tried at the block
		if cfg.KeyGenStatus == crosschaintypes.KeygenStatus_PendingKeygen && !triedKeygenAtBlock {
			// Return error if RPC is not working
			currentBlock, err := zetaBridge.GetZetaBlockHeight()
			if err != nil {
				startLogger.Error().Err(err).Msg("GetZetaBlockHeight RPC  error")
				continue
			}
			// If not at keygen block do not try to generate TSS
			if currentBlock != cfg.KeygenBlock {
				if currentBlock > lastBlock {
					lastBlock = currentBlock
					startLogger.Info().Msgf("Waiting For Keygen Block to arrive or new keygen block to be set. Keygen Block : %d", cfg.KeygenBlock)
				}
				continue
			}
			// Try keygen only once at a particular block, irrespective of whether it is successful or failure
			triedKeygenAtBlock = true
			err = keygenTss(cfg, tss, masterLogger)
			if err != nil {
				startLogger.Error().Err(err).Msg("keygenTss error")
				tssFailedVoteHash, err := zetaBridge.SetTSS("", cfg.KeygenBlock, common.ReceiveStatus_Failed)
				if err != nil {
					startLogger.Error().Err(err).Msg("Failed to broadcast Failed TSS Vote to zetacore")
				}
				startLogger.Info().Msgf("TSS Failed Vote: %s", tssFailedVoteHash)
				continue
			}

			// If TSS is successful , broadcast the vote to zetacore and set Pubkey
			tssSuccessVoteHash, err := zetaBridge.SetTSS(tss.CurrentPubkey, cfg.KeygenBlock, common.ReceiveStatus_Success)
			if err != nil {
				startLogger.Error().Err(err).Msg("TSS successful but unable to broadcast vote to zeta-core")
				return err
			}
			startLogger.Info().Msgf("TSS successful Vote: %s", tssSuccessVoteHash)
			err = SetTSSPubKey(tss, masterLogger)
			if err != nil {
				startLogger.Error().Err(err).Msg("SetTSSPubKey error")
			}
			continue
		}
	}
	err = TestTSS(tss, masterLogger)
	if err != nil {
		startLogger.Error().Err(err).Msg("TestTSS error")
	}

	startLogger.Info().Msgf("TSS address \n ETH : %s \n BTC : %s \n PubKey : %s ", tss.EVMAddress(), tss.BTCAddress(), tss.CurrentPubkey)

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

	//// report TSS address nonce on all chains except zeta
	//for _, chain := range cfg.ChainsEnabled {
	//	if chain.IsExternalChain() {
	//		err = (chainClientMap)[chain].PostNonceIfNotRecorded(startLogger)
	//		if err != nil {
	//			startLogger.Fatal().Err(err).Msgf("PostNonceIfNotRecorded fail %s", chain.String())
	//		}
	//	}
	//}

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
