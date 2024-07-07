package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/hashicorp/go-getter"
	"github.com/rs/zerolog"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc"

	"github.com/zeta-chain/zetacore/zetaclient/config"
)

const zetaclientdBinaryName = "zetaclientd"

var defaultUpgradesDir = os.ExpandEnv("$HOME/.zetaclientd/upgrades")

// serializedWriter wraps an io.Writer and ensures that writes to it from multiple goroutines
// are serialized
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
		logger = zerolog.New(zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}).
			Level(zerolog.Level(cfg.LogLevel)).
			With().
			Timestamp().
			Logger()
	default:
		logger = zerolog.New(zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339})
	}

	return logger
}

type zetaclientdSupervisor struct {
	zetacoredConn      *grpc.ClientConn
	reloadSignals      chan bool
	logger             zerolog.Logger
	upgradesDir        string
	upgradePlanName    string
	enableAutoDownload bool
	restartChan        chan os.Signal
}

func newZetaclientdSupervisor(
	zetaCoreURL string,
	logger zerolog.Logger,
	enableAutoDownload bool,
) (*zetaclientdSupervisor, error) {
	logger = logger.With().Str("module", "zetaclientdSupervisor").Logger()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:9090", zetaCoreURL),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}
	// these signals will result in the supervisor process only restarting zetaclientd
	restartChan := make(chan os.Signal, 1)
	return &zetaclientdSupervisor{
		zetacoredConn:      conn,
		logger:             logger,
		reloadSignals:      make(chan bool, 1),
		upgradesDir:        defaultUpgradesDir,
		enableAutoDownload: enableAutoDownload,
		restartChan:        restartChan,
	}, nil
}

func (s *zetaclientdSupervisor) Start(ctx context.Context) {
	go s.watchForVersionChanges(ctx)
	go s.handleCoreUpgradePlan(ctx)
	go s.handleNewKeygen(ctx)
	go s.handleNewTssKeyGeneration(ctx)
	go s.handleTssUpdate(ctx)
}

func (s *zetaclientdSupervisor) WaitForReloadSignal(ctx context.Context) {
	select {
	case <-s.reloadSignals:
	case <-ctx.Done():
	}
}

func (s *zetaclientdSupervisor) dirForVersion(version string) string {
	return path.Join(s.upgradesDir, version)
}

func atomicSymlink(target, linkName string) error {
	linkNameTmp := linkName + ".tmp"
	_, err := os.Stat(target)
	if err != nil {
		return fmt.Errorf("stat target: %w", err)
	}
	err = os.Remove(linkNameTmp)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove old current tmp: %w", err)
	}
	err = os.Symlink(target, linkNameTmp)
	if err != nil {
		return fmt.Errorf("new symlink: %w", err)
	}
	err = os.Rename(linkNameTmp, linkName)
	if err != nil {
		return fmt.Errorf("rename symlink: %w", err)
	}
	return nil
}

func (s *zetaclientdSupervisor) watchForVersionChanges(ctx context.Context) {
	client := tmservice.NewServiceClient(s.zetacoredConn)
	prevVersion := ""
	for {
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}
		res, err := client.GetNodeInfo(ctx, &tmservice.GetNodeInfoRequest{})
		if err != nil {
			s.logger.Warn().Err(err).Msg("get node info")
			continue
		}
		newVersion := res.ApplicationVersion.Version
		if prevVersion == "" {
			prevVersion = newVersion
		}
		if prevVersion == newVersion {
			continue
		}
		s.logger.Warn().Msgf("core version change (%s -> %s)", prevVersion, newVersion)

		prevVersion = newVersion

		// TODO: just use newVersion when #2135 is merged
		// even without #2135, the version will still change and trigger the update
		newVersionDir := s.dirForVersion(s.upgradePlanName)
		currentLinkPath := s.dirForVersion("current")

		err = atomicSymlink(newVersionDir, currentLinkPath)
		if err != nil {
			s.logger.Error().
				Err(err).
				Msgf("unable to update current symlink (%s -> %s)", newVersionDir, currentLinkPath)
			return
		}
		s.reloadSignals <- true
	}
}

func (s *zetaclientdSupervisor) handleTssUpdate(ctx context.Context) {
	maxRetries := 11
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		client := observertypes.NewQueryClient(s.zetacoredConn)
		tss, err := client.TSS(ctx, &observertypes.QueryGetTSSRequest{})
		if err != nil {
			s.logger.Warn().Err(err).Msg("unable to get original tss")
			time.Sleep(retryInterval)
			continue
		}
		i = 0
		for {
			select {
			case <-time.After(time.Second):
			case <-ctx.Done():
				return
			}
			tssNew, err := client.TSS(ctx, &observertypes.QueryGetTSSRequest{})
			if err != nil {
				s.logger.Warn().Err(err).Msg("unable to get tss")
				continue
			}

			if tssNew.TSS.TssPubkey == tss.TSS.TssPubkey {
				continue
			}

			tss = tssNew
			s.logger.Warn().Msg(fmt.Sprintf("tss address is updated from %s to %s", tss.TSS.TssPubkey, tssNew.TSS.TssPubkey))
			time.Sleep(6 * time.Second)
			s.logger.Info().Msg("restarting zetaclientd to update tss address")
			s.restartChan <- syscall.SIGHUP
		}
	}
	return
}

func (s *zetaclientdSupervisor) handleNewTssKeyGeneration(ctx context.Context) {
	maxRetries := 11
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		client := observertypes.NewQueryClient(s.zetacoredConn)
		alltss, err := client.TssHistory(ctx, &observertypes.QueryTssHistoryRequest{})
		if err != nil {
			s.logger.Warn().Err(err).Msg("unable to get tss original history")
			time.Sleep(retryInterval)
			continue
		}
		i = 0
		tssLenCurrent := len(alltss.TssList)
		for {
			select {
			case <-time.After(time.Second):
			case <-ctx.Done():
				return
			}
			tssListNew, err := client.TssHistory(ctx, &observertypes.QueryTssHistoryRequest{})
			if err != nil {
				s.logger.Warn().Err(err).Msg("unable to get tss new history")
				continue
			}
			tssLenUpdated := len(tssListNew.TssList)

			if tssLenUpdated == tssLenCurrent {
				continue
			}
			if tssLenUpdated < tssLenCurrent {
				tssLenCurrent = len(tssListNew.TssList)
				continue
			}

			tssLenCurrent = tssLenUpdated
			s.logger.Warn().Msg(fmt.Sprintf("tss list updated from %d to %d", tssLenCurrent, tssLenUpdated))
			time.Sleep(5 * time.Second)
			s.logger.Info().Msg("restarting zetaclientd to update tss list")
			s.restartChan <- syscall.SIGHUP
		}
	}
	return
}

func (s *zetaclientdSupervisor) handleNewKeygen(ctx context.Context) {
	client := observertypes.NewQueryClient(s.zetacoredConn)
	prevKeygenBlock := int64(0)
	for {
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}
		resp, err := client.Keygen(ctx, &observertypes.QueryGetKeygenRequest{})
		if err != nil {
			s.logger.Warn().Err(err).Msg("unable to get keygen")
			continue
		}
		if resp.Keygen == nil {
			s.logger.Warn().Err(err).Msg("keygen is nil")
			continue
		}

		if resp.Keygen.Status != observertypes.KeygenStatus_PendingKeygen {
			continue
		}
		keygenBlock := resp.Keygen.BlockNumber
		if prevKeygenBlock == keygenBlock {
			continue
		}
		prevKeygenBlock = keygenBlock
		s.logger.Info().Msgf("got new keygen at block %d", keygenBlock)
		s.restartChan <- syscall.SIGHUP
	}
}
func (s *zetaclientdSupervisor) handleCoreUpgradePlan(ctx context.Context) {
	client := upgradetypes.NewQueryClient(s.zetacoredConn)

	prevPlanName := ""
	for {
		// wait for either a second or context cancel
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}

		resp, err := client.CurrentPlan(ctx, &upgradetypes.QueryCurrentPlanRequest{})
		if err != nil {
			s.logger.Warn().Err(err).Msg("get current upgrade plan")
			continue
		}
		if resp.Plan == nil {
			continue
		}
		plan := resp.Plan
		if prevPlanName == plan.Name {
			continue
		}
		s.logger.Warn().Msgf("got new upgrade plan (%s)", plan.Name)
		prevPlanName = plan.Name
		s.upgradePlanName = plan.Name

		if !s.enableAutoDownload {
			s.logger.Warn().Msg("skipping autodownload because of configuration")
			continue
		}
		err = s.downloadZetaclientd(ctx, plan)
		if err != nil {
			s.logger.Error().Err(err).Msg("downloadZetaclientd failed")
		}
	}
}

// UpgradeConfig is expected format for the info field to allow auto-download
// this structure is copied from cosmosvisor
type upgradeConfig struct {
	Binaries map[string]string `json:"binaries"`
}

func (s *zetaclientdSupervisor) downloadZetaclientd(ctx context.Context, plan *upgradetypes.Plan) error {
	if plan.Info == "" {
		return errors.New("upgrade info empty")
	}
	var config upgradeConfig
	err := json.Unmarshal([]byte(plan.Info), &config)
	if err != nil {
		return fmt.Errorf("unmarshal upgrade config: %w", err)
	}

	s.logger.Info().Msg("downloading zetaclientd")

	binKey := fmt.Sprintf("%s-%s/%s", zetaclientdBinaryName, runtime.GOOS, runtime.GOARCH)
	binURL, ok := config.Binaries[binKey]
	if !ok {
		return fmt.Errorf("no binary found for: %s", binKey)
	}
	upgradeDir := s.dirForVersion(plan.Name)
	err = os.MkdirAll(upgradeDir, 0o750)
	if err != nil {
		return fmt.Errorf("mkdir %s: %w", upgradeDir, err)
	}
	upgradePath := path.Join(upgradeDir, zetaclientdBinaryName)
	// TODO: retry?
	// GetFile should validate checksum so long as it was provided in the url
	err = getter.GetFile(upgradePath, binURL, getter.WithContext(ctx), getter.WithUmask(0o750))
	if err != nil {
		return fmt.Errorf("get file %s: %w", binURL, err)
	}

	// ensure binary is executable
	info, err := os.Stat(upgradePath)
	if err != nil {
		return fmt.Errorf("stat binary: %w", err)
	}
	newMode := info.Mode().Perm() | 0o111
	err = os.Chmod(upgradePath, newMode)
	if err != nil {
		return fmt.Errorf("chmod %s: %w", upgradePath, err)
	}
	return nil
}

func promptPasswords() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("HotKey Password: ")
	hotKeyPass, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	fmt.Print("TSS Password: ")
	tssKeyPass, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	//trim delimiters
	hotKeyPass = strings.TrimSuffix(hotKeyPass, "\n")
	tssKeyPass = strings.TrimSuffix(tssKeyPass, "\n")

	return hotKeyPass, tssKeyPass, err
}
