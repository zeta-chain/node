package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	authzclient "github.com/zeta-chain/zetacore/zetaclient/authz"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/orchestrator"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

type Multiaddr = core.Multiaddr

const (
	// ObserverDBPath is the path (relative to user's home) to the observer database.
	ObserverDBPath = ".zetaclient/chainobserver"
)

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

	// Load Config file from given path
	cfg, err := config.Load(rootArgs.zetaCoreHome)
	if err != nil {
		return err
	}
	if len(cfg.Peer) != 0 {
		err := validatePeer(cfg.Peer)
		if err != nil {
			log.Error().Err(err).Msg("invalid peer")
			return err
		}
	}

	// Load compliance config
	compliance.LoadComplianceConfig(cfg)

	// Initialize base logger
	logger, err := base.InitLogger(cfg)
	if err != nil {
		log.Error().Err(err).Msg("InitLogger failed")
		return err
	}
	masterLogger := logger.Std
	startLogger := masterLogger.With().Str("module", "startup").Logger()
	startLogger.Info().Msgf("zetaclient config file: \n%s", maskCfg(cfg))

	// Wait until zetacore is up
	waitForZetaCore(cfg, startLogger)
	startLogger.Info().Msgf("zetacore is ready, trying to connect to %s", cfg.Peer)

	// Start telemetry server
	telemetryServer := metrics.NewTelemetryServer()
	telemetryServer.SetIPAddress(cfg.PublicIP)
	go func() {
		err := telemetryServer.Start()
		if err != nil {
			startLogger.Error().Err(err).Msg("telemetryServer error")
			panic("telemetryServer error")
		}
	}()

	// Start metrics server
	m, err := metrics.NewMetrics()
	if err != nil {
		log.Error().Err(err).Msg("NewMetrics")
		return err
	}
	m.Start()
	metrics.Info.WithLabelValues(constant.Version).Set(1)
	metrics.LastStartTime.SetToCurrentTime()

	// Create zetacore client to communicate with zetacore.
	// Zetacore accumulates votes, and provides a centralized source of truth for all clients
	zetacoreClient, err := zetacore.CreateClient(cfg, telemetryServer, hotkeyPass)
	if err != nil {
		startLogger.Error().Err(err).Msg("Create zetacore client error")
		return err
	}

	// Wait until zetacore is ready to create blocks
	err = zetacoreClient.WaitForZetacoreToCreateBlocks()
	if err != nil {
		startLogger.Error().Err(err).Msg("WaitForZetacoreToCreateBlocks error")
		return err
	}
	startLogger.Info().Msgf("Zetacore client is ready")

	// Set grantee account number and sequence number
	err = zetacoreClient.SetAccountNumber(authz.ZetaClientGranteeKey)
	if err != nil {
		startLogger.Error().Err(err).Msg("SetAccountNumber error")
		return err
	}

	// cross-check chainid
	res, err := zetacoreClient.GetNodeInfo()
	if err != nil {
		startLogger.Error().Err(err).Msg("GetNodeInfo error")
		return err
	}

	if strings.Compare(res.GetDefaultNodeInfo().Network, cfg.ChainID) != 0 {
		startLogger.Warn().
			Msgf("chain id mismatch, zetacore chain id %s, zetaclient configured chain id %s; reset zetaclient chain id", res.GetDefaultNodeInfo().Network, cfg.ChainID)
		cfg.ChainID = res.GetDefaultNodeInfo().Network
		err := zetacoreClient.UpdateChainID(cfg.ChainID)
		if err != nil {
			return err
		}
	}

	// Set up authz signer to sign all authz messages. All votes broadcast to zetacore are wrapped in authz exec.
	// This is to ensure that the user does not need to keep their operator key online, and can use a cold key to sign votes
	granter := zetacoreClient.GetKeys().GetOperatorAddress().String()
	grantee, err := zetacoreClient.GetKeys().GetAddress()
	if err != nil {
		startLogger.Error().Err(err).Msg("error getting signer address")
		return err
	}
	authzclient.SetupAuthZSignerList(granter, grantee)
	startLogger.Info().Msgf("Authz is ready for granter %s grantee %s", granter, grantee)

	// Initialize zetaclient app context
	appContext, err := orchestrator.CreateAppContext(cfg, zetacoreClient, startLogger)
	if err != nil {
		startLogger.Error().Err(err).Msg("error creating app context")
		return err
	}

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

	// Generate a new TSS key if there is a planned keygen ceremony
	tssHistoricalList, err := zetacoreClient.GetTssHistory()
	if err != nil {
		startLogger.Error().Err(err).Msg("GetTssHistory error")
	}

	tss, err := GenerateTss(
		appContext,
		masterLogger,
		zetacoreClient,
		peers,
		priKey,
		telemetryServer,
		tssHistoricalList,
		tssKeyPass,
		hotkeyPass,
	)
	if err != nil {
		return err
	}
	if cfg.TestTssKeysign {
		err = TestTSS(tss, masterLogger)
		if err != nil {
			startLogger.Error().Err(err).Msgf("TestTSS error : %s", tss.CurrentPubkey)
		}
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
	startLogger.Info().
		Msgf("Current TSS address \n ETH : %s \n BTC : %s \n PubKey : %s ", tss.EVMAddress(), tss.BTCAddress(), tss.CurrentPubkey)
	if len(appContext.GetEnabledExternalChains()) == 0 {
		startLogger.Error().Msgf("No external chains enabled in the zetacore %s ", cfg.String())
	}

	// Stop zetaclient if this node is not an active observer
	observerList, err := zetacoreClient.GetObserverList()
	if err != nil {
		startLogger.Error().Err(err).Msg("GetObserverList error")
		return err
	}
	isNodeActive := false
	for _, observer := range observerList {
		if observer == granter {
			isNodeActive = true
			break
		}
	}
	if !isNodeActive {
		startLogger.Error().Msgf("Node %s is not an active observer, zetaclient stopped", granter)
		return nil
	}
	startLogger.Info().Msgf("Node %s is an active observer, starting orchestrator", granter)

	// use the user's home path to store observer database
	userDir, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Msg("os.UserHomeDir")
		return err
	}
	dbPath := filepath.Join(userDir, ObserverDBPath)

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

	// Orchestrator wraps the zetacore client and adds the observers and signer maps to it . This is the high level object used for CCTX interactions
	orch := orchestrator.NewOrchestrator(
		appContext,
		zetacoreClient,
		tss,
		logger,
		dbPath,
		telemetryServer,
	)
	orch.Start()

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
