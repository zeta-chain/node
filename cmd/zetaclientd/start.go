package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/go-tss/p2p"
	"github.com/zeta-chain/zetacore/pkg/authz"
	"github.com/zeta-chain/zetacore/pkg/constant"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/orchestrator"
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
	err := setHomeDir()
	if err != nil {
		return err
	}

	SetupConfigForTest()

	//Prompt for Hotkey and TSS key-share passwords
	hotkeyPass, tssKeyPass, err := promptPasswords()
	if err != nil {
		return err
	}

	//Load Config file given path
	cfg, err := config.Load(rootArgs.zetaCoreHome)
	if err != nil {
		return err
	}
	loggers, err := InitLogger(cfg)
	if err != nil {
		log.Error().Err(err).Msg("InitLogger failed")
		return err
	}

	//Wait until zetacore has started
	if len(cfg.Peer) != 0 {
		err := validatePeer(cfg.Peer)
		if err != nil {
			log.Error().Err(err).Msg("invalid peer")
			return err
		}
	}

	masterLogger := loggers.Std
	startLogger := masterLogger.With().Str("module", "startup").Logger()

	// Wait until zetacore is up
	waitForZetaCore(cfg, startLogger)
	startLogger.Info().Msgf("Zetacore is ready, trying to connect to %s", cfg.Peer)

	telemetryServer := metrics.NewTelemetryServer()
	go func() {
		err := telemetryServer.Start()
		if err != nil {
			startLogger.Error().Err(err).Msg("telemetryServer error")
			panic("telemetryServer error")
		}
	}()

	// CreateZetaCoreClient:  zetacore client is used for all communication to zetacore , which this client connects to.
	// Zetacore accumulates votes , and provides a centralized source of truth for all clients
	zetacoreClient, err := CreateZetaCoreClient(cfg, telemetryServer, hotkeyPass)
	if err != nil {
		panic(err)
	}
	zetacoreClient.WaitForCoreToCreateBlocks()
	startLogger.Info().Msgf("Zetacore client is ready")
	zetacoreClient.SetAccountNumber(authz.ZetaClientGranteeKey)

	// cross-check chainid
	res, err := zetacoreClient.GetNodeInfo()
	if err != nil {
		panic(err)
	}

	if strings.Compare(res.GetDefaultNodeInfo().Network, cfg.ChainID) != 0 {
		startLogger.Warn().Msgf("chain id mismatch, zeta-core chain id %s, zeta client chain id %s; reset zeta client chain id", res.GetDefaultNodeInfo().Network, cfg.ChainID)
		cfg.ChainID = res.GetDefaultNodeInfo().Network
		err := zetacoreClient.UpdateChainID(cfg.ChainID)
		if err != nil {
			return err
		}
	}

	// CreateAuthzSigner : which is used to sign all authz messages . All votes broadcast to zetacore are wrapped in authz exec .
	// This is to ensure that the user does not need to keep their operator key online , and can use a cold key to sign votes
	CreateAuthzSigner(zetacoreClient.GetKeys().GetOperatorAddress().String(), zetacoreClient.GetKeys().GetAddress())
	startLogger.Debug().Msgf("CreateAuthzSigner is ready")

	// Initialize core parameters from zetacore
	appContext := context.NewAppContext(context.NewZetaCoreContext(cfg), cfg)
	err = zetacoreClient.UpdateZetaCoreContext(appContext.ZetaCoreContext(), true, startLogger)
	if err != nil {
		startLogger.Error().Err(err).Msg("Error getting core parameters")
		return err
	}
	startLogger.Info().Msgf("Config is updated from zetacore %s", maskCfg(cfg))

	go zetacoreClient.CoreContextUpdater(appContext)

	// Generate TSS address . The Tss address is generated through Keygen ceremony. The TSS key is used to sign all outbound transactions .
	// The hotkeyPk is private key for the Hotkey. The Hotkey is used to sign all inbound transactions
	// Each node processes a portion of the key stored in ~/.tss by default . Custom location can be specified in config file during init.
	// After generating the key , the address is set on the zetacore
	hotkeyPk, err := zetacoreClient.GetKeys().GetPrivateKey(hotkeyPass)
	if err != nil {
		startLogger.Error().Err(err).Msg("zetacore client GetPrivateKey error")
	}
	startLogger.Debug().Msgf("hotkeyPk %s", hotkeyPk.String())
	if len(hotkeyPk.Bytes()) != 32 {
		errMsg := fmt.Sprintf("key bytes len %d != 32", len(hotkeyPk.Bytes()))
		log.Error().Msgf(errMsg)
		return errors.New(errMsg)
	}
	priKey := secp256k1.PrivKey(hotkeyPk.Bytes()[:32])

	// Generate pre Params if not present already
	peers, err := initPeers(cfg.Peer)
	if err != nil {
		log.Error().Err(err).Msg("peer address error")
	}
	initPreParams(cfg.PreParamsPath)
	if cfg.P2PDiagnostic {
		err := RunDiagnostics(startLogger, peers, hotkeyPk, cfg)
		if err != nil {
			startLogger.Error().Err(err).Msg("RunDiagnostics error")
			return err
		}
	}

	m, err := metrics.NewMetrics()
	if err != nil {
		log.Error().Err(err).Msg("NewMetrics")
		return err
	}
	m.Start()

	metrics.Info.WithLabelValues(constant.Version).Set(1)
	metrics.LastStartTime.SetToCurrentTime()

	var tssHistoricalList []observerTypes.TSS
	tssHistoricalList, err = zetacoreClient.GetTssHistory()
	if err != nil {
		startLogger.Error().Err(err).Msg("GetTssHistory error")
	}

	telemetryServer.SetIPAddress(cfg.PublicIP)
	tss, err := GenerateTss(appContext, masterLogger, zetacoreClient, peers, priKey, telemetryServer, tssHistoricalList, tssKeyPass, hotkeyPass)
	if err != nil {
		return err
	}
	if cfg.TestTssKeysign {
		err = TestTSS(tss, masterLogger)
		if err != nil {
			startLogger.Error().Err(err).Msgf("TestTSS error : %s", tss.CurrentPubkey)
		}
	}

	// Wait for TSS keygen to be successful before proceeding, This is a blocking thread only for a new keygen.
	// For existing keygen, this should directly proceed to the next step
	ticker := time.NewTicker(time.Second * 1)
	for range ticker.C {
		keyGen := appContext.ZetaCoreContext().GetKeygen()
		if keyGen.Status != observerTypes.KeygenStatus_KeyGenSuccess {
			startLogger.Info().Msgf("Waiting for TSS Keygen to be a success, current status %s", keyGen.Status)
			continue
		}
		break
	}

	// Update Current TSS value from zetacore, if TSS keygen is successful, the TSS address is set on zeta-core
	// Returns err if the RPC call fails as zeta client needs the current TSS address to be set
	// This is only needed in case of a new Keygen , as the TSS address is set on zetacore only after the keygen is successful i.e enough votes have been broadcast
	currentTss, err := zetacoreClient.GetCurrentTss()
	if err != nil {
		startLogger.Error().Err(err).Msg("GetCurrentTSS error")
		return err
	}

	// Defensive check: Make sure the tss address is set to the current TSS address and not the newly generated one
	tss.CurrentPubkey = currentTss.TssPubkey
	if tss.EVMAddress() == (ethcommon.Address{}) || tss.BTCAddress() == "" {
		startLogger.Error().Msg("TSS address is not set in zetacore")
	}
	startLogger.Info().Msgf("Current TSS address \n ETH : %s \n BTC : %s \n PubKey : %s ", tss.EVMAddress(), tss.BTCAddress(), tss.CurrentPubkey)
	if len(appContext.ZetaCoreContext().GetEnabledChains()) == 0 {
		startLogger.Error().Msgf("No chains enabled in updated config %s ", cfg.String())
	}

	observerList, err := zetacoreClient.GetObserverList()
	if err != nil {
		startLogger.Error().Err(err).Msg("GetObserverList error")
		return err
	}
	isNodeActive := false
	for _, observer := range observerList {
		if observer == zetacoreClient.GetKeys().GetOperatorAddress().String() {
			isNodeActive = true
			break
		}
	}

	// CreateSignerMap: This creates a map of all signers for each chain . Each signer is responsible for signing transactions for a particular chain
	signerMap, err := CreateSignerMap(appContext, tss, loggers, telemetryServer)
	if err != nil {
		log.Error().Err(err).Msg("CreateSignerMap")
		return err
	}

	userDir, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Msg("os.UserHomeDir")
		return err
	}
	dbpath := filepath.Join(userDir, ".zetaclient/chainobserver")

	// Creates a map of all chain observers for each chain. Each chain observer is responsible for observing events on the chain and processing them.
	observerMap, err := CreateChainObserverMap(appContext, zetacoreClient, tss, dbpath, loggers, telemetryServer)
	if err != nil {
		startLogger.Err(err).Msg("CreateChainObserverMap")
		return err
	}

	if !isNodeActive {
		startLogger.Error().Msgf("Node %s is not an active observer external chain observers will not be started", zetacoreClient.GetKeys().GetOperatorAddress().String())
	} else {
		startLogger.Debug().Msgf("Node %s is an active observer starting external chain observers", zetacoreClient.GetKeys().GetOperatorAddress().String())
		for _, observer := range observerMap {
			observer.Start()
		}
	}

	// Orchestrator wraps the zetacore client and adds the observers and signer maps to it . This is the high level object used for CCTX interactions
	mo1 := orchestrator.NewOrchestrator(zetacoreClient, signerMap, observerMap, masterLogger, telemetryServer)
	mo1.MonitorCore(appContext)

	// start zeta supply checker
	// TODO: enable
	// https://github.com/zeta-chain/node/issues/1354
	// NOTE: this is disabled for now because we need to determine the frequency on how to handle invalid check
	// The method uses GRPC query to the node we might need to improve for performance
	//zetaSupplyChecker, err := mc.NewZetaSupplyChecker(cfg, zetacoreClient, masterLogger)
	//if err != nil {
	//	startLogger.Err(err).Msg("NewZetaSupplyChecker")
	//}
	//if err == nil {
	//	zetaSupplyChecker.Start()
	//	defer zetaSupplyChecker.Stop()
	//}

	startLogger.Info().Msgf("awaiting the os.Interrupt, syscall.SIGTERM signals...")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	startLogger.Info().Msgf("stop signal received: %s", sig)

	// stop chain observers
	for _, observer := range observerMap {
		observer.Stop()
	}
	zetacoreClient.Stop()

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
			bz, err := io.ReadAll(preParamsFile)
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

// promptPasswords() This function will prompt for passwords which will be used to decrypt two key files:
// 1. HotKey
// 2. TSS key-share
func promptPasswords() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("HotKey Password: ")
	hotKeyPass, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	fmt.Print("TSS Password: ")
	TSSKeyPass, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	//trim delimiters
	hotKeyPass = strings.TrimSuffix(hotKeyPass, "\n")
	TSSKeyPass = strings.TrimSuffix(TSSKeyPass, "\n")

	return hotKeyPass, TSSKeyPass, err
}
