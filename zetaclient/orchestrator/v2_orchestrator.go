package orchestrator

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// V2 represents the orchestrator V2 while they co-exist with Orchestrator.
type V2 struct {
	deps      *Dependencies
	scheduler *scheduler.Scheduler

	chains map[int64]ObserverSigner
	mu     sync.RWMutex

	logger loggers
}

type loggers struct {
	zerolog.Logger
	sampled zerolog.Logger
	base    base.Logger
}

const schedulerGroup = scheduler.Group("orchestrator")

type ObserverSigner interface {
	Chain() chains.Chain
	Start(ctx context.Context) error
	Stop()
}

type Dependencies struct {
	Zetacore  interfaces.ZetacoreClient
	TSS       interfaces.TSSSigner
	DBPath    string
	Telemetry *metrics.TelemetryServer
}

func NewV2(scheduler *scheduler.Scheduler, deps *Dependencies, logger base.Logger) (*V2, error) {
	if err := validateConstructor(scheduler, deps); err != nil {
		return nil, errors.Wrap(err, "invalid args")
	}

	return &V2{
		scheduler: scheduler,
		deps:      deps,
		chains:    make(map[int64]ObserverSigner),
		logger:    newLoggers(logger),
	}, nil
}

func (oc *V2) Start(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	// syntax sugar
	opts := func(name string, opts ...scheduler.Opt) []scheduler.Opt {
		return append(opts, scheduler.GroupName(schedulerGroup), scheduler.Name(name))
	}

	contextInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(app.Config().ConfigUpdateTicker)
	})

	// every other block
	syncInterval := scheduler.Interval(2 * constant.ZetaBlockTime)

	oc.scheduler.Register(ctx, oc.UpdateContext, opts("update_context", contextInterval)...)
	oc.scheduler.Register(ctx, oc.SyncChains, opts("sync_chains", syncInterval)...)

	return nil
}

func (oc *V2) Stop() {
	oc.logger.Info().Msg("Stopping orchestrator")

	// stops *all* scheduler tasks
	oc.scheduler.Stop()
}

func (oc *V2) UpdateContext(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	err = UpdateAppContext(ctx, app, oc.deps.Zetacore, oc.logger.Logger)

	switch {
	case errors.Is(err, ErrUpgradeRequired):
		const msg = "Upgrade detected. Kill the process, " +
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

func (oc *V2) SyncChains(ctx context.Context) error {
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
			// todo
			// https://github.com/zeta-chain/node/issues/3302
			continue
		case chain.IsSolana():
			// todo
			// https://github.com/zeta-chain/node/issues/3301
			continue
		case chain.IsTON():
			// todo
			// https://github.com/zeta-chain/node/issues/3300
			continue
		}

		switch {
		case errors.Is(errSkipChain, err):
			oc.logger.sampled.Warn().Err(err).Fields(chain.LogFields()).Msg("Skipping observer-signer")
			continue
		case err != nil:
			oc.logger.Error().Err(err).Fields(chain.LogFields()).Msg("Failed to bootstrap observer-signer")
			continue
		case observerSigner == nil:
			// should not happen
			oc.logger.Error().Fields(chain.LogFields()).Msg("Nil observer-signer")
			continue
		}

		if err = observerSigner.Start(ctx); err != nil {
			oc.logger.Error().Err(err).Fields(chain.LogFields()).Msg("Failed to start observer-signer")
			continue
		}

		oc.addChain(observerSigner)
		added++
	}

	removed = oc.removeMissingChains(presentChainIDs)

	if (added + removed) > 0 {
		oc.logger.Info().
			Int("chains.added", added).
			Int("chains.removed", removed).
			Msg("Synced observer-signers")
	}

	return nil
}

func (oc *V2) hasChain(chainID int64) bool {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	_, ok := oc.chains[chainID]
	return ok
}

func (oc *V2) chainIDs() []int64 {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	ids := make([]int64, 0, len(oc.chains))
	for chainID := range oc.chains {
		ids = append(ids, chainID)
	}

	return ids
}

func (oc *V2) addChain(observerSigner ObserverSigner) {
	chain := observerSigner.Chain()

	oc.mu.Lock()
	oc.chains[chain.ChainId] = observerSigner
	oc.mu.Unlock()

	oc.logger.Info().Fields(chain.LogFields()).Msg("Added observer-signer")
}

func (oc *V2) removeChain(chainID int64) {
	// noop, should not happen
	if !oc.hasChain(chainID) {
		return
	}

	// blocking call
	oc.chains[chainID].Stop()

	oc.mu.Lock()
	delete(oc.chains, chainID)
	oc.mu.Unlock()

	oc.logger.Info().Int64(logs.FieldChain, chainID).Msg("Removed observer-signer")
}

// removeMissingChains stops and deletes chains
// that are not present in the list of chainIDs (e.g. after governance proposal)
func (oc *V2) removeMissingChains(presentChainIDs []int64) int {
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

func validateConstructor(s *scheduler.Scheduler, dep *Dependencies) error {
	switch {
	case s == nil:
		return errors.New("scheduler is nil")
	case dep == nil:
		return errors.New("dependencies are nil")
	case dep.Zetacore == nil:
		return errors.New("zetacore is nil")
	case dep.TSS == nil:
		return errors.New("tss is nil")
	case dep.Telemetry == nil:
		return errors.New("telemetry is nil")
	case dep.DBPath == "":
		return errors.New("db path is empty")
	}

	return nil
}

func newLoggers(baseLogger base.Logger) loggers {
	std := baseLogger.Std.With().Str(logs.FieldModule, "orchestrator").Logger()

	return loggers{
		Logger:  std,
		sampled: std.Sample(&zerolog.BasicSampler{N: 10}),
		base:    baseLogger,
	}
}
