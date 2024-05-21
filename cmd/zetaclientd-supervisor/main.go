package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	// load zetaclient config
	cfg, err := config.Load(app.DefaultNodeHome)
	if err != nil {
		fmt.Println("failed to load config: ", err)
		os.Exit(1)
	}

	// log outputs must be serialized since we are writing log messages in this process and
	// also directly from the zetaclient process
	serializedStdout := &serializedWriter{upstream: os.Stdout}
	logger := getLogger(cfg, serializedStdout)
	logger = logger.With().Str("process", "zetaclientd-supervisor").Logger()

	ctx := context.Background()

	// these signals will result in the supervisor process shutting down
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// these signals will result in the supervisor process only restarting zetaclientd
	restartChan := make(chan os.Signal, 1)
	signal.Notify(restartChan, syscall.SIGHUP)

	hotkeyPassword, tssPassword, err := promptPasswords()
	if err != nil {
		logger.Error().Err(err).Msg("unable to get passwords")
		os.Exit(1)
	}

	_, enableAutoDownload := os.LookupEnv("ZETACLIENTD_SUPERVISOR_ENABLE_AUTO_DOWNLOAD")
	supervisor, err := newZetaclientdSupervisor(cfg.ZetaCoreURL, logger, enableAutoDownload)
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
		// by default, CommandContext sends SIGKILL. we want more graceful shutdown.
		cmd.Cancel = func() error {
			return cmd.Process.Signal(syscall.SIGINT)
		}
		cmd.Stdout = serializedStdout
		cmd.Stderr = os.Stderr
		// must reset the passwordInputBuffer every iteration because reads are stateful (seek to end)
		passwordInputBuffer := bytes.Buffer{}
		passwordInputBuffer.Write([]byte(hotkeyPassword + "\n" + tssPassword + "\n"))
		cmd.Stdin = &passwordInputBuffer

		eg, ctx := errgroup.WithContext(ctx)
		eg.Go(cmd.Run)
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
				case sig := <-restartChan:
					logger.Info().Msgf("got signal %d, sending SIGINT to zetaclientd", sig)
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
