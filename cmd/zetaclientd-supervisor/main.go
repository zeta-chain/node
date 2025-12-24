package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cosmossdk.io/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/app"
	zetaos "github.com/zeta-chain/node/pkg/os"
	_ "github.com/zeta-chain/node/pkg/sdkconfig/default"
	"github.com/zeta-chain/node/zetaclient/config"
)

func main() {
	// load zetaclient config
	cfg, err := config.Load(app.DefaultNodeHome)
	if err != nil {
		fmt.Println("failed to load config: ", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// log outputs must be serialized since we are writing log messages in this process and
	// also directly from the zetaclient process
	syncWriter := zerolog.SyncWriter(os.Stdout)
	logger := getLogger(cfg, syncWriter)
	logger = logger.With().Str("process", "zetaclientd-supervisor").Logger()

	// these signals will result in the supervisor process shutting down
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// prompt for all necessary passwords
	titles := []string{"HotKey", "TSS", "Solana Relayer Key"}
	passwords, err := zetaos.PromptPasswords(titles)
	if err != nil {
		logger.Error().Err(err).Msg("unable to get passwords")
		os.Exit(1)
	}

	clientCfg := cfg.GetZetacoreClientConfig()
	_, enableAutoDownload := os.LookupEnv("ZETACLIENTD_SUPERVISOR_ENABLE_AUTO_DOWNLOAD")
	supervisor, err := newZetaclientdSupervisor(clientCfg.GRPCURL, logger, enableAutoDownload)
	if err != nil {
		logger.Error().Err(err).Msg("unable to get supervisor")
		os.Exit(1)
	}
	supervisor.Start(ctx)

	shouldRestart := true
	for shouldRestart {
		ctx, cancel := context.WithCancel(ctx)
		// pass args from supervisor directly to zetaclientd
		cmd := exec.CommandContext(ctx, zetaclientdBinaryName, os.Args[1:]...) // #nosec G204
		cmd.Stdout = syncWriter
		cmd.Stderr = os.Stderr

		// by default, CommandContext sends SIGKILL. we want more graceful shutdown.
		cmd.Cancel = func() error {
			return cmd.Process.Signal(syscall.SIGINT)
		}

		// must reset the passwordInputBuffer every iteration because reads are stateful (seek to end)
		passwordInputBuffer := bytes.Buffer{}
		passwordInputBuffer.Write([]byte(strings.Join(passwords, "\n") + "\n"))
		cmd.Stdin = &passwordInputBuffer

		eg, ctx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			defer cancel()
			if err := cmd.Run(); err != nil {
				return errors.Wrap(err, "zetaclient process failed")
			}

			logger.Info().Msg("zetaclient process exited")
			return nil
		})
		eg.Go(func() error {
			supervisor.WaitForReloadSignal(ctx)
			cancel()
			return nil
		})
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case sig := <-shutdownChan:
					logger.Info().Msgf("got signal %d, shutting down", sig)
					shouldRestart = false
				}
				cancel()
			}
		})
		err := eg.Wait()
		if err != nil {
			logger.Error().Err(err).Msg("error while waiting")
		}
		// prevent fast spin
		time.Sleep(time.Second)
	}
}
