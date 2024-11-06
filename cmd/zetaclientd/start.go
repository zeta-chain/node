package main

import (
	"context"
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
	"github.com/libp2p/go-libp2p/core/peer"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gitlab.com/thorchain/tss/go-tss/conversion"

	"github.com/zeta-chain/node/pkg/authz"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	zetaos "github.com/zeta-chain/node/pkg/os"
	observerTypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/maintenance"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/orchestrator"
	mc "github.com/zeta-chain/node/zetaclient/tss"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

func Start(_ *cobra.Command, _ []string) error {
	// Prompt for Hotkey, TSS key-share and relayer key passwords
	titles := []string{"HotKey", "TSS", "Solana Relayer Key"}
	passwords, err := zetaos.PromptPasswords(titles)
	if err != nil {
		return errors.Wrap(err, "unable to get passwords")
	}
	hotkeyPass, tssKeyPass, solanaKeyPass := passwords[0], passwords[1], passwords[2]
	relayerKeyPasswords := map[string]string{
		chains.Network_solana.String(): solanaKeyPass,
	}

	//Load Config file given path
	cfg, err := config.Load(globalOpts.ZetacoreHome)
	if err != nil {
		return err
	}

	logger, err := base.InitLogger(cfg)
	if err != nil {
		return errors.Wrap(err, "initLogger failed")
	}

	// Wait until zetacore has started
	if len(cfg.Peer) != 0 {
		if err := validatePeer(cfg.Peer); err != nil {
			return errors.Wrap(err, "unable to validate peer")
		}
	}

	masterLogger := logger.Std
	startLogger := logger.Std.With().Str("module", "startup").Logger()

	appContext := zctx.New(cfg, relayerKeyPasswords, masterLogger)
	ctx := zctx.WithAppContext(context.Background(), appContext)

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

	// CreateZetacoreClient:  zetacore client is used for all communication to zetacore , which this client connects to.
	// Zetacore accumulates votes , and provides a centralized source of truth for all clients
	zetacoreClient, err := createZetacoreClient(cfg, hotkeyPass, masterLogger)
	if err != nil {
		return errors.Wrap(err, "unable to create zetacore client")
	}

	// Wait until zetacore is ready to create blocks
	if err = zetacoreClient.WaitForZetacoreToCreateBlocks(ctx); err != nil {
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
	res, err := zetacoreClient.GetNodeInfo(ctx)
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

	// CreateAuthzSigner : which is used to sign all authz messages . All votes broadcast to zetacore are wrapped in authz exec .
	// This is to ensure that the user does not need to keep their operator key online , and can use a cold key to sign votes
	signerAddress, err := zetacoreClient.GetKeys().GetAddress()
	if err != nil {
		return errors.Wrap(err, "error getting signer address")
	}

	createAuthzSigner(zetacoreClient.GetKeys().GetOperatorAddress().String(), signerAddress)
	startLogger.Debug().Msgf("createAuthzSigner is ready")

	// Initialize core parameters from zetacore
	if err = zetacoreClient.UpdateAppContext(ctx, appContext, startLogger); err != nil {
		return errors.Wrap(err, "unable to update app context")
	}

	startLogger.Info().Msgf("Config is updated from zetacore\n %s", cfg.StringMasked())

	go zetacoreClient.UpdateAppContextWorker(ctx, appContext)

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
		log.Error().Msg(errMsg)
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
	tssHistoricalList, err = zetacoreClient.GetTSSHistory(ctx)
	if err != nil {
		startLogger.Error().Err(err).Msg("GetTssHistory error")
	}

	telemetryServer.SetIPAddress(cfg.PublicIP)

	keygen := appContext.GetKeygen()
	whitelistedPeers := []peer.ID{}
	for _, pk := range keygen.GranteePubkeys {
		pid, err := conversion.Bech32PubkeyToPeerID(pk)
		if err != nil {
			return err
		}
		whitelistedPeers = append(whitelistedPeers, pid)
	}

	// Create TSS server
	server, err := mc.SetupTSSServer(peers, priKey, preParams, appContext.Config(), tssKeyPass, true, whitelistedPeers)
	if err != nil {
		return fmt.Errorf("SetupTSSServer error: %w", err)
	}
	// Set P2P ID for telemetry
	telemetryServer.SetP2PID(server.GetLocalPeerID())

	// Creating a channel to listen for os signals (or other signals)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Maintenance workers ============
	maintenance.NewTSSListener(zetacoreClient, masterLogger).Listen(ctx, func() {
		masterLogger.Info().Msg("TSS listener received an action to shutdown zetaclientd.")
		signalChannel <- syscall.SIGTERM
	})

	// Generate a new TSS if keygen is set and add it into the tss server
	// If TSS has already been generated, and keygen was successful ; we use the existing TSS
	err = GenerateTSS(ctx, masterLogger, zetacoreClient, server)
	if err != nil {
		return err
	}

	tss, err := mc.NewTSS(
		ctx,
		zetacoreClient,
		tssHistoricalList,
		hotkeyPass,
		server,
	)
	if err != nil {
		startLogger.Error().Err(err).Msg("NewTSS error")
		return err
	}
	if cfg.TestTssKeysign {
		err = TestTSS(tss.CurrentPubkey, *tss.Server, masterLogger)
		if err != nil {
			startLogger.Error().Err(err).Msgf("TestTSS error : %s", tss.CurrentPubkey)
		}
	}

	// Wait for TSS keygen to be successful before proceeding, This is a blocking thread only for a new keygen.
	// For existing keygen, this should directly proceed to the next step
	ticker := time.NewTicker(time.Second * 1)
	for range ticker.C {
		keyGen := appContext.GetKeygen()
		if keyGen.Status != observerTypes.KeygenStatus_KeyGenSuccess {
			startLogger.Info().Msgf("Waiting for TSS Keygen to be a success, current status %s", keyGen.Status)
			continue
		}
		break
	}

	// Update Current TSS value from zetacore, if TSS keygen is successful, the TSS address is set on zeta-core
	// Returns err if the RPC call fails as zeta client needs the current TSS address to be set
	// This is only needed in case of a new Keygen , as the TSS address is set on zetacore only after the keygen is successful i.e enough votes have been broadcast
	currentTss, err := zetacoreClient.GetTSS(ctx)
	if err != nil {
		startLogger.Error().Err(err).Msg("GetCurrentTSS error")
		return err
	}

	// Filter supported BTC chain IDs
	btcChains := appContext.FilterChains(zctx.Chain.IsBitcoin)
	btcChainIDs := make([]int64, len(btcChains))
	for i, chain := range btcChains {
		btcChainIDs[i] = chain.ID()
	}

	// Make sure the TSS EVM/BTC addresses are well formed.
	// Zetaclient should not start if TSS addresses cannot be properly derived.
	tss.CurrentPubkey = currentTss.TssPubkey
	err = tss.ValidateAddresses(btcChainIDs)
	if err != nil {
		startLogger.Error().Err(err).Msg("TSS address validation failed")
		return err
	}

	if len(appContext.ListChainIDs()) == 0 {
		startLogger.Error().Interface("config", cfg).Msgf("No chains in updated config")
	}

	isObserver, err := isObserverNode(ctx, zetacoreClient)
	switch {
	case err != nil:
		startLogger.Error().Msgf("Unable to determine if node is an observer")
		return err
	case !isObserver:
		addr := zetacoreClient.GetKeys().GetOperatorAddress().String()
		startLogger.Info().Str("operator_address", addr).Msg("This node is not an observer. Exit 0")
		return nil
	}

	// CreateSignerMap: This creates a map of all signers for each chain.
	// Each signer is responsible for signing transactions for a particular chain
	signerMap, err := orchestrator.CreateSignerMap(ctx, tss, logger, telemetryServer)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create signer map")
		return err
	}

	userDir, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Msg("os.UserHomeDir")
		return err
	}
	dbpath := filepath.Join(userDir, ".zetaclient/chainobserver")

	// Creates a map of all chain observers for each chain.
	// Each chain observer is responsible for observing events on the chain and processing them.
	observerMap, err := orchestrator.CreateChainObserverMap(ctx, zetacoreClient, tss, dbpath, logger, telemetryServer)
	if err != nil {
		startLogger.Err(err).Msg("CreateChainObserverMap")
		return err
	}

	// Orchestrator wraps the zetacore client and adds the observers and signer maps to it.
	// This is the high level object used for CCTX interactions
	maestro, err := orchestrator.New(
		ctx,
		zetacoreClient,
		signerMap,
		observerMap,
		tss,
		dbpath,
		logger,
		telemetryServer,
	)
	if err != nil {
		startLogger.Error().Err(err).Msg("Unable to create orchestrator")
		return err
	}

	// Start orchestrator with all observers and signers
	if err := maestro.Start(ctx); err != nil {
		startLogger.Error().Err(err).Msg("Unable to start orchestrator")
		return err
	}

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

	startLogger.Info().Msgf("Zetaclientd is running")

	sig := <-signalChannel
	startLogger.Info().Msgf("Stop signal received: %q", sig)

	zetacoreClient.Stop()

	return nil
}

func initPeers(peer string) ([]maddr.Multiaddr, error) {
	var peers []maddr.Multiaddr

	if peer != "" {
		address, err := maddr.NewMultiaddr(peer)
		if err != nil {
			log.Error().Err(err).Msg("NewMultiaddr error")
			return []maddr.Multiaddr{}, err
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

// isObserverNode checks whether THIS node is an observer node.
func isObserverNode(ctx context.Context, client *zetacore.Client) (bool, error) {
	observers, err := client.GetObserverList(ctx)
	if err != nil {
		return false, errors.Wrap(err, "unable to get observers list")
	}

	operatorAddress := client.GetKeys().GetOperatorAddress().String()

	for _, observer := range observers {
		if observer == operatorAddress {
			return true, nil
		}
	}

	return false, nil
}
