// Package orchestrator is responsible for (de)provisioning, running, and monitoring various
// observer-signer instances.
//
// It also updates app context with data from zetacore (eg chain parameters).
package orchestrator

import (
	"context"
	"math"
	"sync"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/tssrepo"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/mode/chaos"
)

// Orchestrator chain orchestrator.
type Orchestrator struct {
	scheduler *scheduler.Scheduler

	zetacoreClient zetacoreClient
	tssClient      tssrepo.TSSClient
	telemetry      *metrics.TelemetryServer
	dbPath         string

	chains map[int64]ObserverSigner
	mu     sync.RWMutex

	operatorBalance sdkmath.Int

	chaosSource *chaos.Source

	logger loggers
}

// zetacoreClient aggregates the orchestrator and zrepo interfaces for ZetacoreClient.
type zetacoreClient interface {
	ZetacoreClient
	zrepo.ZetacoreClient
}

type loggers struct {
	zerolog.Logger
	sampled zerolog.Logger
	base    base.Logger
}

type ObserverSigner interface {
	Chain() chains.Chain
	Start(ctx context.Context) error
	Stop()
}

const schedulerGroup = scheduler.Group("orchestrator")

func New(
	scheduler *scheduler.Scheduler,
	zetacoreClient zetacoreClient,
	tssClient tssrepo.TSSClient,
	telemetry *metrics.TelemetryServer,
	dbPath string,
	config config.Config,
	logger base.Logger,
) (*Orchestrator, error) {
	switch {
	case scheduler == nil:
		return nil, errors.New("invalid scheduler")
	case zetacoreClient == nil:
		return nil, errors.New("invalid zetacore client")
	case tssClient == nil:
		return nil, errors.New("invalid TSS client")
	case telemetry == nil:
		return nil, errors.New("invalid telemetry server")
	case dbPath == "":
		return nil, errors.New("invalid database path")
	}

	var chaosSource *chaos.Source
	if config.ClientMode.IsChaosMode() {
		source, err := chaos.NewSource(logger.Std, config)
		if err != nil {
			return nil, err
		}
		chaosSource = source
	}

	return &Orchestrator{
		scheduler:      scheduler,
		zetacoreClient: zetacoreClient,
		tssClient:      tssClient,
		telemetry:      telemetry,
		dbPath:         dbPath,
		chains:         make(map[int64]ObserverSigner),
		chaosSource:    chaosSource,
		logger:         newLoggers(logger),
	}, nil
}

func (oc *Orchestrator) Start(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	newBlocksChan, err := oc.zetacoreClient.NewBlockSubscriber(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to subscribe to new block")
	}

	// syntax sugar
	opts := func(name string, opts ...scheduler.Opt) []scheduler.Opt {
		return append(opts, scheduler.GroupName(schedulerGroup), scheduler.Name(name))
	}

	contextInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(app.Config().ConfigUpdateTicker)
	})

	// every other block, regardless of block events from zetacore
	syncInterval := scheduler.Interval(2 * constant.ZetaBlockTime)

	blocksTicker := scheduler.BlockTicker(newBlocksChan)

	// refresh preflight metrics in a lazy manner
	preflightTicker := scheduler.Interval(1 * time.Minute)

	// check feature flags and log their status
	oc.logFeatureFlags(app.Config())

	oc.scheduler.Register(ctx, oc.UpdateContext, opts("update_context", contextInterval)...)
	oc.scheduler.Register(ctx, oc.SyncChains, opts("sync_chains", syncInterval)...)
	oc.scheduler.Register(ctx, oc.updateMetrics, opts("update_metrics", blocksTicker)...)
	oc.scheduler.Register(ctx, oc.reportPreflightMetrics, opts("report_preflight_metrics", preflightTicker)...)

	return nil
}

func (oc *Orchestrator) Stop() {
	oc.logger.Info().Msg("stopping the orchestrator")

	// stops *all* scheduler tasks
	oc.scheduler.Stop()
}

func (oc *Orchestrator) UpdateContext(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	err = UpdateAppContext(ctx, app, oc.zetacoreClient, oc.logger.Logger)

	switch {
	case errors.Is(err, ErrUpgradeRequired):
		const msg = "upgrade detected; kill the process, " +
			"replace the binary with upgraded version, and restart zetaclientd"
		oc.logger.Warn().Str("upgrade", err.Error()).Msg(msg)

		// stop the orchestrator
		go oc.Stop()

		return nil
	case err != nil:
		return errors.Wrap(err, "unable to update app context")
	default:
		return nil
	}
}

var errSkipChain = errors.New("skip chain")

func (oc *Orchestrator) SyncChains(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	var (
		added, removed  int
		presentChainIDs = make([]int64, 0)
	)

	for _, chain := range app.ListChains() {
		// skip zetachain
		if chain.IsZeta() {
			continue
		}

		presentChainIDs = append(presentChainIDs, chain.ID())

		// skip existing chain
		if oc.hasChain(chain.ID()) {
			continue
		}

		var observerSigner ObserverSigner

		switch {
		case chain.IsBitcoin():
			observerSigner, err = oc.bootstrapBitcoin(ctx, chain)
		case chain.IsEVM():
			observerSigner, err = oc.bootstrapEVM(ctx, chain)
		case chain.IsSolana():
			observerSigner, err = oc.bootstrapSolana(ctx, chain)
		case chain.IsSui():
			observerSigner, err = oc.bootstrapSui(ctx, chain)
		case chain.IsTON():
			observerSigner, err = oc.bootstrapTON(ctx, chain)
		}

		switch {
		case errors.Is(err, errSkipChain):
			// TODO use throttled logger instead of sampled one.
			// https://github.com/zeta-chain/node/issues/3336
			oc.logger.sampled.Warn().
				Err(err).
				Fields(chain.LogFields()).
				Msg("skipping observer-signer")
			continue
		case err != nil:
			oc.logger.Error().
				Err(err).
				Fields(chain.LogFields()).
				Msg("failed to bootstrap observer-signer")
			continue
		case observerSigner == nil:
			// should not happen
			oc.logger.Error().
				Fields(chain.LogFields()).
				Msg("nil observer-signer")
			continue
		}

		if err = observerSigner.Start(ctx); err != nil {
			oc.logger.Error().
				Err(err).
				Fields(chain.LogFields()).
				Msg("failed to start observer-signer")
			continue
		}

		oc.addChain(observerSigner)
		added++
	}

	removed = oc.removeMissingChains(presentChainIDs)

	if (added + removed) > 0 {
		oc.logger.Info().
			Int("chains_added", added).
			Int("chains_removed", removed).
			Msg("synced observer-signers")
	}

	return nil
}

var (
	zero   = sdkmath.NewInt(0)
	maxInt = sdkmath.NewInt(math.MaxInt64)
)

func (oc *Orchestrator) updateMetrics(ctx context.Context) error {
	block, sleepDuration, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "unable get block from context")
	}

	zetaBlockHeight := block.Block.Height

	// 0. Set block metrics
	metrics.CoreBlockLatency.Set(time.Since(block.Block.Time).Seconds())
	metrics.CoreBlockLatencySleep.Set(sleepDuration.Seconds())

	oc.telemetry.SetCoreBlockNumber(zetaBlockHeight)

	// 1. Fetch hot key balance
	balance, err := oc.zetacoreClient.GetZetaHotKeyBalance(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get hot key balance")
	}

	// 2. Set it within orchestrator
	oc.operatorBalance = balance

	// 3. Update telemetry
	diff := oc.operatorBalance.Sub(balance)
	if diff.GT(zero) && diff.LT(maxInt) {
		oc.telemetry.AddFeeEntry(zetaBlockHeight, diff.Int64())
	}

	// 4. Update metrics
	burnRate := oc.telemetry.HotKeyBurnRate.GetBurnRate().Int64()
	metrics.HotKeyBurnRate.Set(float64(burnRate))

	return nil
}

func (oc *Orchestrator) reportPreflightMetrics(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	return ReportPreflightMetrics(ctx, app, oc.zetacoreClient, oc.logger.Logger)
}

func (oc *Orchestrator) hasChain(chainID int64) bool {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	_, ok := oc.chains[chainID]
	return ok
}

func (oc *Orchestrator) chainIDs() []int64 {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	ids := make([]int64, 0, len(oc.chains))
	for chainID := range oc.chains {
		ids = append(ids, chainID)
	}

	return ids
}

func (oc *Orchestrator) addChain(observerSigner ObserverSigner) {
	chain := observerSigner.Chain()

	oc.mu.Lock()
	defer oc.mu.Unlock()

	// noop
	if _, ok := oc.chains[chain.ChainId]; ok {
		return
	}

	oc.chains[chain.ChainId] = observerSigner
	oc.logger.Info().Fields(chain.LogFields()).Msg("added observer-signer")
}

func (oc *Orchestrator) removeChain(chainID int64) {
	// noop, should not happen
	if !oc.hasChain(chainID) {
		return
	}

	// blocking call
	oc.chains[chainID].Stop()

	oc.mu.Lock()
	delete(oc.chains, chainID)
	oc.mu.Unlock()

	oc.logger.Info().Int64(logs.FieldChain, chainID).Msg("removed observer-signer")
}

// removeMissingChains stops and deletes chains
// that are not present in the list of chainIDs (e.g. after governance proposal)
func (oc *Orchestrator) removeMissingChains(presentChainIDs []int64) int {
	presentChainsSet := make(map[int64]struct{})
	for _, chainID := range presentChainIDs {
		presentChainsSet[chainID] = struct{}{}
	}

	existingIDs := oc.chainIDs()
	removed := 0

	for _, chainID := range existingIDs {
		if _, ok := presentChainsSet[chainID]; ok {
			// all good, chain is present
			continue
		}

		oc.removeChain(chainID)
		removed++
	}

	return removed
}

func newLoggers(baseLogger base.Logger) loggers {
	std := baseLogger.Std.With().Str(logs.FieldModule, logs.ModNameOrchestrator).Logger()
	return loggers{
		Logger:  std,
		sampled: std.Sample(&zerolog.BasicSampler{N: 10}),
		base:    baseLogger,
	}
}

// logFeatureFlags logs the current status of feature flags
func (oc *Orchestrator) logFeatureFlags(config config.Config) {
	flags := config.GetFeatureFlags()

	oc.logger.Info().
		Bool("enable_multiple_calls", flags.EnableMultipleCalls).
		Bool("enable_solana_address_lookup_table", flags.EnableSolanaAddressLookupTable).
		Msg("feature flags status")
}
