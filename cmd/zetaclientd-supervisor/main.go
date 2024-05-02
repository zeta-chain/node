package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type serializedWriter struct {
	upstream io.Writer
	lock     sync.Mutex
}

func (w *serializedWriter) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.upstream.Write(p)
}

func getLogger(cfg config.Config, out io.Writer) zerolog.Logger {
	var logger zerolog.Logger
	switch cfg.LogFormat {
	case "json":
		logger = zerolog.New(out).Level(zerolog.Level(cfg.LogLevel)).With().Timestamp().Logger()
	case "text":
		logger = zerolog.New(zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}).Level(zerolog.Level(cfg.LogLevel)).With().Timestamp().Logger()
	default:
		logger = zerolog.New(zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339})
	}

	return logger
}

func watchForVersionChanges(ctx context.Context, zetaCoreUrl string, logger zerolog.Logger) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:9090", zetaCoreUrl),
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Warn().Err(err).Msg("grpc dial fail")
		return
	}
	defer conn.Close()
	client := tmservice.NewServiceClient(conn)
	prevVersion := ""
	for {
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}
		res, err := client.GetNodeInfo(ctx, &tmservice.GetNodeInfoRequest{})
		if err != nil {
			logger.Warn().Err(err).Msg("get node info")
			continue
		}
		currentVersion := res.ApplicationVersion.Version
		if prevVersion == "" {
			prevVersion = currentVersion
		} else if prevVersion != currentVersion {
			logger.Warn().Msgf("core version change (%s -> %s), signaling for restart", prevVersion, currentVersion)
			return
		}
	}
}

func main() {
	cfg, err := config.Load(app.DefaultNodeHome)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	// log outputs must be serialized since we are writing log messages in this process and
	// also directly from the zetaclient process
	serializedStdout := &serializedWriter{upstream: os.Stdout}
	logger := getLogger(cfg, serializedStdout)
	logger = logger.With().Str("module", "zetaclientd-supervisor").Logger()

	ctx := context.Background()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	hotkeyPassword, tssPassword, err := promptPasswords()
	if err != nil {
		panic(fmt.Errorf("unable to get passwords: %w", err))
	}

	for {
		cmd := exec.Command("zetaclientd", "start")
		cmd.Stdout = serializedStdout
		cmd.Stderr = os.Stderr
		// must reset the passwordInputBuffer every iteration because reads are stateful (seek to end)
		passwordInputBuffer := bytes.Buffer{}
		passwordInputBuffer.Write([]byte(hotkeyPassword + "\n" + tssPassword + "\n"))
		cmd.Stdin = &passwordInputBuffer

		ctx, cancel := context.WithCancel(ctx)
		eg, ctx := errgroup.WithContext(ctx)
		eg.Go(cmd.Run)
		eg.Go(func() error {
			watchForVersionChanges(ctx, cfg.ZetaCoreURL, logger)
			cancel()
			return nil
		})
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case sig := <-signalChan:
					logger.Info().Msgf("got signal %d, forwarding to zetaclientd", sig)
					_ = cmd.Process.Signal(sig)
				}
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
