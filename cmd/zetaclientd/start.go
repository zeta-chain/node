package main

import (
	"context"
	"net/http"
	_ "net/http/pprof" // #nosec G108 -- pprof enablement is intentional
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

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
)

const (
	// enables posting blame data to core for failed TSS signatures
	envFlagPostBlame = "POST_BLAME"
	envPprofAddr     = "PPROF_ADDR"
)

// Start starts zetaclientd process
func Start(_ *cobra.Command, _ []string) error {
	// Load Config file given path
	cfg, err := config.Load(globalOpts.ZetacoreHome)
	if err != nil {
		return errors.Wrap(err, "unable to load config")
	}

	dbPath, err := config.ResolveDBPath()
	if err != nil {
		return errors.Wrap(err, "unable to resolve db path")
	}

	// Configure logger (also overrides the default log level)
	logger, err := base.NewLogger(cfg)
	if err != nil {
		return errors.Wrap(err, "unable to create logger")
	}

	passes, err := promptPasswords()
	if err != nil {
		return errors.Wrap(err, "unable to prompt for passwords")
	}

	appContext := zctx.New(cfg, passes.relayerKeys(), logger.Std)
	ctx := zctx.WithAppContext(context.Background(), appContext)

	// TODO graceful
	telemetryServer := metrics.NewTelemetryServer()
	go func() {
		err := telemetryServer.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("telemetryServer error")
		}
	}()

	m, err := metrics.NewMetrics()
	if err != nil {
		return errors.Wrap(err, "unable to create metrics")
	}
	m.Start()

	metrics.Info.WithLabelValues(constant.Version).Set(1)
	metrics.LastStartTime.SetToCurrentTime()

	telemetryServer.SetIPAddress(cfg.PublicIP)

	// TODO graceful
	go runPprof(logger.Std)

	// zetacore client is used for all communication to zeta node.
	// it accumulates votes, and provides a source of truth for all clients
	zetacoreClient, err := createZetacoreClient(cfg, passes.hotkey, logger.Std)
	if err != nil {
		return errors.Wrap(err, "unable to create zetacore client")
	}

	// Wait until zetacore is ready to produce blocks
	if err = waitForBlocks(ctx, zetacoreClient, logger.Std); err != nil {
		return errors.Wrap(err, "zetacore unavailable")
	}

	if err = prepareZetacoreClient(ctx, zetacoreClient, &cfg, logger.Std); err != nil {
		return errors.Wrap(err, "unable to prepare zetacore client")
	}

	// Initialize core parameters from zetacore
	if err = orchestrator.UpdateAppContext(ctx, appContext, zetacoreClient, logger.Std); err != nil {
		return errors.Wrap(err, "unable to update app context")
	}

	log.Info().Msgf("Config is updated from zetacore\n %s", cfg.StringMasked())

	granteePubKeyBech32, err := resolveObserverPubKeyBech32(cfg, passes.hotkey)
	if err != nil {
		return errors.Wrap(err, "unable to resolve observer pub key bech32")
	}

	tssSetupProps := zetatss.SetupProps{
		Config:              cfg,
		Zetacore:            zetacoreClient,
		GranteePubKeyBech32: granteePubKeyBech32,
		HotKeyPassword:      passes.hotkey,
		TSSKeyPassword:      passes.tss,
		BitcoinChainIDs:     btcChainIDsFromContext(appContext),
		PostBlame:           isEnvFlagEnabled(envFlagPostBlame),
		Telemetry:           telemetryServer,
	}

	tss, err := zetatss.Setup(ctx, tssSetupProps, logger.Std)
	if err != nil {
		return errors.Wrap(err, "unable to setup TSS service")
	}

	// Creating a channel to listen for os signals (or other signals)
	// TODO graceful
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Starts various background TSS listeners.
	// Shuts down zetaclientd if any is triggered.
	maintenance.NewTSSListener(zetacoreClient, logger.Std).Listen(ctx, func() {
		logger.Std.Info().Msg("TSS listener received an action to shutdown zetaclientd.")
		signalChannel <- syscall.SIGTERM
	})

	// CreateSignerMap: This creates a map of all signers for each chain.
	// Each signer is responsible for signing transactions for a particular chain
	signerMap, err := orchestrator.CreateSignerMap(ctx, tss, logger)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create signer map")
		return err
	}

	// Creates a map of all chain observers for each chain.
	// Each chain observer is responsible for observing events on the chain and processing them.
	observerMap, err := orchestrator.CreateChainObserverMap(ctx, zetacoreClient, tss, dbPath, logger, telemetryServer)
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
		dbPath,
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

	log.Info().Msg("zetaclientd is running")

	// todo graceful
	sig := <-signalChannel
	log.Info().Msgf("Stop signal received: %q. Stopping zetaclientd", sig)

	maestro.Stop()

	return nil
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

type passwords struct {
	hotkey           string
	tss              string
	solanaRelayerKey string
}

// promptPasswords prompts for Hotkey, TSS key-share and relayer key passwords
func promptPasswords() (passwords, error) {
	titles := []string{"HotKey", "TSS", "Solana Relayer Key"}

	res, err := zetaos.PromptPasswords(titles)
	if err != nil {
		return passwords{}, errors.Wrap(err, "unable to get passwords")
	}

	return passwords{
		hotkey:           res[0],
		tss:              res[1],
		solanaRelayerKey: res[2],
	}, nil
}

func (p passwords) relayerKeys() map[string]string {
	return map[string]string{
		chains.Network_solana.String(): p.solanaRelayerKey,
	}
}
