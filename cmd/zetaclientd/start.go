package main

import (
	"context"
	"net/http"
	_ "net/http/pprof" // #nosec G108 -- runPprof enablement is intentional
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/authz"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	zetaos "github.com/zeta-chain/node/pkg/os"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/maintenance"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/orchestrator"
	zetatss "github.com/zeta-chain/node/zetaclient/tss"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const (
	// enables posting blame data to core for failed TSS signatures
	envFlagPostBlame = "POST_BLAME"
	envPprofAddr     = "PPROF_ADDR"
)

// Start starts zetaclientd process todo revamp
// https://github.com/zeta-chain/node/issues/3112
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

	// Load Config file given path
	cfg, err := config.Load(globalOpts.ZetacoreHome)
	if err != nil {
		return err
	}

	logger, err := base.InitLogger(cfg)
	if err != nil {
		return errors.Wrap(err, "initLogger failed")
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

	go runPprof(startLogger)

	// CreateZetacoreClient:  zetacore client is used for all communication to zetacore , which this client connects to.
	// Zetacore accumulates votes , and provides a centralized source of truth for all clients
	zetacoreClient, err := createZetacoreClient(cfg, hotkeyPass, masterLogger)
	if err != nil {
		return errors.Wrap(err, "unable to create zetacore client")
	}

	// Wait until zetacore is ready to create blocks
	if err = waitForZetacoreToCreateBlocks(ctx, zetacoreClient, startLogger); err != nil {
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
	if err = orchestrator.UpdateAppContext(ctx, appContext, zetacoreClient, startLogger); err != nil {
		return errors.Wrap(err, "unable to update app context")
	}

	startLogger.Info().Msgf("Config is updated from zetacore\n %s", cfg.StringMasked())

	m, err := metrics.NewMetrics()
	if err != nil {
		return errors.Wrap(err, "unable to create metrics")
	}
	m.Start()

	metrics.Info.WithLabelValues(constant.Version).Set(1)
	metrics.LastStartTime.SetToCurrentTime()

	telemetryServer.SetIPAddress(cfg.PublicIP)

	granteePubKeyBech32, err := resolveObserverPubKeyBech32(cfg, hotkeyPass)
	if err != nil {
		return errors.Wrap(err, "unable to resolve observer pub key bech32")
	}

	tssSetupProps := zetatss.SetupProps{
		Config:              cfg,
		Zetacore:            zetacoreClient,
		GranteePubKeyBech32: granteePubKeyBech32,
		HotKeyPassword:      hotkeyPass,
		TSSKeyPassword:      tssKeyPass,
		BitcoinChainIDs:     btcChainIDsFromContext(appContext),
		PostBlame:           isEnvFlagEnabled(envFlagPostBlame),
		Telemetry:           telemetryServer,
	}

	tss, err := zetatss.Setup(ctx, tssSetupProps, startLogger)
	if err != nil {
		return errors.Wrap(err, "unable to setup TSS service")
	}

	// Creating a channel to listen for os signals (or other signals)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Starts various background TSS listeners.
	// Shuts down zetaclientd if any is triggered.
	maintenance.NewTSSListener(zetacoreClient, masterLogger).Listen(ctx, func() {
		masterLogger.Info().Msg("TSS listener received an action to shutdown zetaclientd.")
		signalChannel <- syscall.SIGTERM
	})

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
		return errors.Wrap(err, "unable to create chain observer map")
	}

	// Orchestrator wraps the zetacore client and adds the observers and signer maps to it.
	// This is the high level object used for CCTX interactions
	// It also handles background configuration updates from zetacore
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
		return errors.Wrap(err, "unable to create orchestrator")
	}

	// Start orchestrator with all observers and signers
	if err = maestro.Start(ctx); err != nil {
		return errors.Wrap(err, "unable to start orchestrator")
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

	startLogger.Info().Msg("zetaclientd is running")

	sig := <-signalChannel
	startLogger.Info().Msgf("Stop signal received: %q. Stopping zetaclientd", sig)

	maestro.Stop()

	return nil
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

func resolveObserverPubKeyBech32(cfg config.Config, hotKeyPassword string) (string, error) {
	// Get observer's public key ("grantee pub key")
	_, granteePubKeyBech32, err := keys.GetKeyringKeybase(cfg, hotKeyPassword)
	if err != nil {
		return "", errors.Wrap(err, "unable to get keyring key base")
	}

	return granteePubKeyBech32, nil
}

// runPprof run pprof http server
// zetacored/cometbft is already listening for runPprof on 6060 (by default)
func runPprof(logger zerolog.Logger) {
	addr := os.Getenv(envPprofAddr)
	if addr == "" {
		addr = "localhost:6061"
	}

	logger.Info().Str("addr", addr).Msg("starting pprof http server")

	// #nosec G114 -- timeouts unneeded
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Error().Err(err).Msg("pprof http server error")
	}
}
