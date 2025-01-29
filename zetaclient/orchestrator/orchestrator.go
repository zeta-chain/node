// Package orchestrator provides the orchestrator for orchestrating cross-chain transactions
package orchestrator

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	sdkmath "cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/samber/lo"

	"github.com/zeta-chain/node/pkg/bg"
	"github.com/zeta-chain/node/pkg/constant"
	zetamath "github.com/zeta-chain/node/pkg/math"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	solanaobserver "github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	solanasigner "github.com/zeta-chain/node/zetaclient/chains/solana/signer"
	tonobserver "github.com/zeta-chain/node/zetaclient/chains/ton/observer"
	tonsigner "github.com/zeta-chain/node/zetaclient/chains/ton/signer"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/ratelimiter"
)

const (
	// outboundLookbackFactor is the factor to determine how many nonces to look back for pending cctxs
	// For example, give OutboundScheduleLookahead of 120, pending NonceLow of 1000 and factor of 1.0,
	// the scheduler need to be able to pick up and schedule any pending cctx with nonce < 880 (1000 - 120 * 1.0)
	// NOTE: 1.0 means look back the same number of cctxs as we look ahead
	outboundLookbackFactor = 1.0

	// sampling rate for sampled orchestrator logger
	loggerSamplingRate = 10
)

var defaultLogSampler = &zerolog.BasicSampler{N: loggerSamplingRate}

// Orchestrator wraps the zetacore client, chain observers and signers. This is the high level object used for CCTX scheduling
type Orchestrator struct {
	// zetacore client
	zetacoreClient interfaces.ZetacoreClient

	// signerMap contains the chain signers indexed by chainID
	signerMap map[int64]interfaces.ChainSigner

	// observerMap contains the chain observers indexed by chainID
	observerMap map[int64]interfaces.ChainObserver

	// last operator balance
	lastOperatorBalance sdkmath.Int

	// observer & signer props
	tss         interfaces.TSSSigner
	dbDirectory string
	baseLogger  base.Logger

	// signerBlockTimeOffset
	signerBlockTimeOffset time.Duration

	// misc
	logger  multiLogger
	ts      *metrics.TelemetryServer
	stop    chan struct{}
	stopped atomic.Bool
	mu      sync.RWMutex
}

type multiLogger struct {
	zerolog.Logger
	Sampled zerolog.Logger
}

// New creates a new Orchestrator
func New(
	ctx context.Context,
	client interfaces.ZetacoreClient,
	signerMap map[int64]interfaces.ChainSigner,
	observerMap map[int64]interfaces.ChainObserver,
	tss interfaces.TSSSigner,
	dbDirectory string,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (*Orchestrator, error) {
	if signerMap == nil || observerMap == nil {
		return nil, errors.New("signerMap or observerMap is nil")
	}

	logging := logger.Std.With().Str(logs.FieldModule, "orchestrator").Logger()
	multiLog := multiLogger{Logger: logging, Sampled: logging.Sample(defaultLogSampler)}

	balance, err := client.GetZetaHotKeyBalance(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get last balance of the hot key")
	}

	return &Orchestrator{
		zetacoreClient: client,

		signerMap:   signerMap,
		observerMap: observerMap,

		lastOperatorBalance: balance,

		// observer & signer props
		tss:         tss,
		dbDirectory: dbDirectory,
		baseLogger:  logger,

		logger: multiLog,
		ts:     ts,
		stop:   make(chan struct{}),
	}, nil
}

// Start starts the orchestrator for CCTXs.
func (oc *Orchestrator) Start(ctx context.Context) error {
	signerAddress, err := oc.zetacoreClient.GetKeys().GetAddress()
	if err != nil {
		return errors.Wrap(err, "unable to get signer address")
	}

	oc.logger.Info().Str("signer", signerAddress.String()).Msg("Starting orchestrator")

	bg.Work(ctx, oc.runScheduler, bg.WithName("runScheduler"), bg.WithLogger(oc.logger.Logger))
	bg.Work(ctx, oc.runObserverSignerSync, bg.WithName("runObserverSignerSync"), bg.WithLogger(oc.logger.Logger))
	bg.Work(
		ctx,
		oc.runSyncObserverOperationalFlags,
		bg.WithName("runSyncObserverOperationalFlags"),
		bg.WithLogger(oc.logger.Logger),
	)

	return nil
}

func (oc *Orchestrator) Stop() {
	// noop
	if oc.stopped.Load() {
		oc.logger.Warn().Msg("Already stopped")
		return
	}

	close(oc.stop)

	oc.stopped.Store(true)
}

// returns signer with updated chain parameters.
func (oc *Orchestrator) resolveSigner(app *zctx.AppContext, chainID int64) (interfaces.ChainSigner, error) {
	signer, err := oc.getSigner(chainID)
	if err != nil {
		return nil, err
	}

	chain, err := app.GetChain(chainID)
	switch {
	case err != nil:
		return nil, err
	case chain.IsZeta():
		// should not happen
		return nil, fmt.Errorf("unable to resolve signer for zeta chain %d", chainID)
	case chain.IsEVM():
		params := chain.Params()

		// update zeta connector, ERC20 custody, and gateway addresses
		signer.SetZetaConnectorAddress(eth.HexToAddress(params.ConnectorContractAddress))
		signer.SetERC20CustodyAddress(eth.HexToAddress(params.Erc20CustodyContractAddress))
		signer.SetGatewayAddress(params.GatewayAddress)
	case chain.IsSolana():
		signer.SetGatewayAddress(chain.Params().GatewayAddress)
	case chain.IsTON():
		signer.SetGatewayAddress(chain.Params().GatewayAddress)
	}

	return signer, nil
}

func (oc *Orchestrator) getSigner(chainID int64) (interfaces.ChainSigner, error) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	s, found := oc.signerMap[chainID]
	if !found {
		return nil, fmt.Errorf("signer not found for chainID %d", chainID)
	}

	return s, nil
}

// returns chain observer with updated chain parameters
func (oc *Orchestrator) resolveObserver(app *zctx.AppContext, chainID int64) (interfaces.ChainObserver, error) {
	observer, err := oc.getObserver(chainID)
	if err != nil {
		return nil, err
	}

	chain, err := app.GetChain(chainID)
	switch {
	case err != nil:
		return nil, errors.Wrapf(err, "unable to get chain %d", chainID)
	case chain.IsZeta():
		// should not happen
		return nil, fmt.Errorf("unable to resolve observer for zeta chain %d", chainID)
	}

	// update chain observer chain parameters
	observer.SetChainParams(*chain.Params())

	return observer, nil
}

func (oc *Orchestrator) getObserver(chainID int64) (interfaces.ChainObserver, error) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	ob, found := oc.observerMap[chainID]
	if !found {
		return nil, fmt.Errorf("observer not found for chainID %d", chainID)
	}

	return ob, nil
}

// GetPendingCctxsWithinRateLimit get pending cctxs across foreign chains within rate limit
func (oc *Orchestrator) GetPendingCctxsWithinRateLimit(ctx context.Context, chainIDs []int64) (
	map[int64][]*types.CrossChainTx,
	error,
) {
	// get rate limiter flags
	rateLimitFlags, err := oc.zetacoreClient.GetRateLimiterFlags(ctx)
	if err != nil {
		return nil, err
	}

	// apply rate limiter or not according to the flags
	rateLimiterUsable := ratelimiter.IsRateLimiterUsable(rateLimitFlags)

	// fallback to non-rate-limited query if rate limiter is not usable
	cctxsMap := make(map[int64][]*types.CrossChainTx)
	if !rateLimiterUsable {
		for _, chainID := range chainIDs {
			resp, _, err := oc.zetacoreClient.ListPendingCCTX(ctx, chainID)
			if err == nil && resp != nil {
				cctxsMap[chainID] = resp
			}
		}
		return cctxsMap, nil
	}

	// query rate limiter input
	resp, err := oc.zetacoreClient.GetRateLimiterInput(ctx, rateLimitFlags.Window)
	if err != nil {
		return nil, err
	}
	input, ok := ratelimiter.NewInput(*resp)
	if !ok {
		return nil, fmt.Errorf("failed to create rate limiter input")
	}

	// apply rate limiter
	output := ratelimiter.ApplyRateLimiter(input, rateLimitFlags.Window, rateLimitFlags.Rate)

	// set metrics
	percentage := zetamath.Percentage(output.CurrentWithdrawRate.BigInt(), rateLimitFlags.Rate.BigInt())
	if percentage != nil {
		percentageFloat, _ := percentage.Float64()
		metrics.PercentageOfRateReached.Set(percentageFloat)
		oc.logger.Sampled.Info().Msgf("current rate limiter window: %d rate: %s, percentage: %f",
			output.CurrentWithdrawWindow, output.CurrentWithdrawRate.String(), percentageFloat)
	}

	return output.CctxsMap, nil
}

// schedules keysigns for cctxs on each ZetaChain block (the ticker)
// TODO(revamp): make this function simpler
func (oc *Orchestrator) runScheduler(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	newBlockChan, err := oc.zetacoreClient.NewBlockSubscriber(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-oc.stop:
			oc.logger.Warn().Msg("runScheduler: stopped")
			return nil
		case newBlock := <-newBlockChan:
			bn := newBlock.Block.Height

			blockTimeLatency := time.Since(newBlock.Block.Time)
			blockTimeLatencySeconds := blockTimeLatency.Seconds()
			metrics.CoreBlockLatency.Set(blockTimeLatencySeconds)

			if blockTimeLatencySeconds > 15 {
				oc.logger.Warn().
					Float64("latency", blockTimeLatencySeconds).
					Msgf("runScheduler: core block latency too high")
				continue
			}

			sleepDuration := time.Until(newBlock.Block.Time.Add(oc.signerBlockTimeOffset))
			if sleepDuration < 0 {
				sleepDuration = 0
			}
			metrics.CoreBlockLatencySleep.Set(sleepDuration.Seconds())
			time.Sleep(sleepDuration)

			balance, err := oc.zetacoreClient.GetZetaHotKeyBalance(ctx)
			if err != nil {
				oc.logger.Error().Err(err).Msgf("couldn't get operator balance")
			} else {
				diff := oc.lastOperatorBalance.Sub(balance)
				if diff.GT(sdkmath.NewInt(0)) && diff.LT(sdkmath.NewInt(math.MaxInt64)) {
					oc.ts.AddFeeEntry(bn, diff.Int64())
					oc.lastOperatorBalance = balance
				}
			}

			// set current hot key burn rate
			metrics.HotKeyBurnRate.Set(float64(oc.ts.HotKeyBurnRate.GetBurnRate().Int64()))

			// get chain ids without zeta chain
			chainIDs := lo.FilterMap(app.ListChains(), func(c zctx.Chain, _ int) (int64, bool) {
				return c.ID(), !c.IsZeta()
			})

			// query pending cctxs across all external chains within rate limit
			cctxMap, err := oc.GetPendingCctxsWithinRateLimit(ctx, chainIDs)
			if err != nil {
				oc.logger.Error().Err(err).Msgf("runScheduler: GetPendingCctxsWithinRatelimit failed")
			}

			// schedule keysign for pending cctxs on each chain
			for _, chain := range app.ListChains() {
				// skip zeta chain
				if chain.IsZeta() {
					continue
				}

				// managed by V2
				if chain.IsBitcoin() || chain.IsEVM() {
					continue
				}

				// todo move metrics to v2

				chainID := chain.ID()

				// update chain parameters for signer and chain observer
				signer, err := oc.resolveSigner(app, chainID)
				if err != nil {
					oc.logger.Error().Err(err).
						Int64(logs.FieldChain, chainID).
						Msg("runScheduler: unable to resolve signer")
					continue
				}

				ob, err := oc.resolveObserver(app, chainID)
				if err != nil {
					oc.logger.Error().Err(err).
						Int64(logs.FieldChain, chainID).
						Msg("runScheduler: unable to resolve observer")
					continue
				}

				// get cctxs from map and set pending transactions prometheus gauge
				cctxList := cctxMap[chainID]

				if len(cctxList) == 0 {
					continue
				}

				if !app.IsOutboundObservationEnabled() {
					continue
				}

				// #nosec G115 range is verified
				zetaHeight := uint64(bn)

				switch {
				case chain.IsEVM():
					// Managed by orchestrator V2
					continue
				case chain.IsBitcoin():
					// Managed by orchestrator V2
					continue
				case chain.IsSolana():
					oc.ScheduleCctxSolana(ctx, zetaHeight, chainID, cctxList, ob, signer)
				case chain.IsTON():
					oc.ScheduleCCTXTON(ctx, zetaHeight, chainID, cctxList, ob, signer)
				default:
					oc.logger.Error().Msgf("runScheduler: no scheduler found chain %d", chainID)
					continue
				}
			}

			// update last processed block number
			oc.ts.SetCoreBlockNumber(bn)
		}
	}
}

// ScheduleCctxSolana schedules solana outbound keysign on each ZetaChain block (the ticker)
func (oc *Orchestrator) ScheduleCctxSolana(
	ctx context.Context,
	zetaHeight uint64,
	chainID int64,
	cctxList []*types.CrossChainTx,
	observer interfaces.ChainObserver,
	signer interfaces.ChainSigner,
) {
	solObserver, ok := observer.(*solanaobserver.Observer)

	// should never happen
	if !ok {
		oc.logger.Error().Msgf("ScheduleCctxSolana: chain observer is not a solana observer")
		return
	}

	solSigner, ok := signer.(*solanasigner.Signer)

	// should never happen
	if !ok {
		oc.logger.Error().Msgf("ScheduleCctxSolana: chain signer is not a solana signer")
		return
	}

	// outbound keysign scheduler parameters
	// #nosec G115 positive
	interval := uint64(observer.ChainParams().OutboundScheduleInterval)
	outboundScheduleLookahead := observer.ChainParams().OutboundScheduleLookahead
	outboundScheduleLookback := uint64(float64(outboundScheduleLookahead) * outboundLookbackFactor)

	// schedule keysign for each pending cctx
	for _, cctx := range cctxList {
		params := cctx.GetCurrentOutboundParam()
		nonce := params.TssNonce
		outboundID := base.OutboundIDFromCCTX(cctx)

		if params.ReceiverChainId != chainID {
			oc.logger.Error().
				Msgf("ScheduleCctxSolana: outbound %s chainid mismatch: want %d, got %d", outboundID, chainID, params.ReceiverChainId)
			continue
		}
		if params.TssNonce > cctxList[0].GetCurrentOutboundParam().TssNonce+outboundScheduleLookback {
			oc.logger.Warn().Msgf("ScheduleCctxSolana: nonce too high: signing %d, earliest pending %d",
				params.TssNonce, cctxList[0].GetCurrentOutboundParam().TssNonce)
			break
		}

		// vote outbound if it's already confirmed
		continueKeysign, err := solObserver.VoteOutboundIfConfirmed(ctx, cctx)
		if err != nil {
			oc.logger.Error().
				Err(err).
				Msgf("ScheduleCctxSolana: VoteOutboundIfConfirmed failed for chain %d nonce %d", chainID, nonce)
			continue
		}
		if !continueKeysign {
			oc.logger.Info().
				Msgf("ScheduleCctxSolana: outbound %s already processed; do not schedule keysign", outboundID)
			continue
		}

		// noop
		if solSigner.IsOutboundActive(outboundID) {
			continue
		}

		// schedule a TSS keysign
		if nonce%interval == zetaHeight%interval {
			go signer.TryProcessOutbound(
				ctx,
				cctx,
				observer,
				oc.zetacoreClient,
				zetaHeight,
			)
		}
	}
}

// ScheduleCCTXTON schedules TON outbound keySign on each ZetaChain block
func (oc *Orchestrator) ScheduleCCTXTON(
	ctx context.Context,
	zetaHeight uint64,
	chainID int64,
	cctxList []*types.CrossChainTx,
	observer interfaces.ChainObserver,
	signer interfaces.ChainSigner,
) {
	// should never happen
	if _, ok := observer.(*tonobserver.Observer); !ok {
		oc.logger.Error().Msg("ScheduleCCTXTON: observer is not TON")
		return
	}

	tonSigner, ok := signer.(*tonsigner.Signer)
	if !ok {
		oc.logger.Error().Msg("ScheduleCCTXTON: signer is not TON")
		return
	}

	// Scheduler interval measured in zeta blocks.
	// runScheduler() guarantees that this function is called every zeta block.
	// Note that TON blockchain is async and doesn't have a concept of confirmations
	// i.e. tx is finalized as soon as it's included in the next block (less than 6 seconds)
	// #nosec G115 positive
	interval := uint64(observer.ChainParams().OutboundScheduleInterval)

	shouldProcessOutbounds := zetaHeight%interval == 0

	for i := range cctxList {
		var (
			cctx       = cctxList[i]
			params     = cctx.GetCurrentOutboundParam()
			nonce      = params.TssNonce
			outboundID = base.OutboundIDFromCCTX(cctx)
		)

		if params.ReceiverChainId != chainID {
			// should not happen
			oc.logger.Error().Msgf("ScheduleCCTXTON: outbound chain id mismatch (got %d)", params.ReceiverChainId)
			continue
		}

		// outbound is already being processed
		if tonSigner.IsOutboundActive(outboundID) {
			continue
		}

		// vote outbound if it's already confirmed
		continueKeySign, err := observer.VoteOutboundIfConfirmed(ctx, cctx)

		switch {
		case err != nil:
			oc.logger.Error().Err(err).Uint64("outbound.nonce", nonce).
				Msg("ScheduleCCTXTON: VoteOutboundIfConfirmed failed")
			continue
		case !continueKeySign:
			oc.logger.Info().Uint64("outbound.nonce", nonce).
				Msg("ScheduleCCTXTON: outbound already processed")
			continue
		case !shouldProcessOutbounds:
			// well, let's wait for another block to (probably) trigger the processing
			continue
		}

		// try to sign and broadcast cctx to TON
		task := func(ctx context.Context) error {
			signer.TryProcessOutbound(
				ctx,
				cctx,
				observer,
				oc.zetacoreClient,
				zetaHeight,
			)

			return nil
		}

		// fire async task
		taskLogger := oc.logger.Logger.With().
			Int64(logs.FieldChain, chainID).
			Str("outbound.id", outboundID).
			Uint64("outbound.nonce", cctx.GetCurrentOutboundParam().TssNonce).
			Logger()

		bg.Work(ctx, task, bg.WithName("TryProcessOutbound"), bg.WithLogger(taskLogger))
	}
}

// runObserverSignerSync runs a blocking ticker that observes chain changes from zetacore
// and optionally (de)provisions respective observers and signers.
func (oc *Orchestrator) runObserverSignerSync(ctx context.Context) error {
	// every other block
	const cadence = 2 * constant.ZetaBlockTime

	task := func(ctx context.Context, _ *ticker.Ticker) error {
		if err := oc.syncObserverSigner(ctx); err != nil {
			oc.logger.Error().Err(err).Msg("syncObserverSigner failed")
		}

		return nil
	}

	return ticker.Run(
		ctx,
		cadence,
		task,
		ticker.WithLogger(oc.logger.Logger, "SyncObserverSigner"),
		ticker.WithStopChan(oc.stop),
	)
}

// syncs and provisions observers & signers.
// Note that zctx.AppContext Update is a responsibility of another component
// See zetacore.Client{}.UpdateAppContextWorker
func (oc *Orchestrator) syncObserverSigner(ctx context.Context) error {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	client := oc.zetacoreClient

	added, removed, err := syncObserverMap(ctx, client, oc.tss, oc.dbDirectory, oc.baseLogger, oc.ts, &oc.observerMap)
	if err != nil {
		return errors.Wrap(err, "syncObserverMap failed")
	}

	if added+removed > 0 {
		oc.logger.Info().
			Int("observer.added", added).
			Int("observer.removed", removed).
			Msg("synced observers")
	}

	added, removed, err = syncSignerMap(ctx, oc.tss, oc.baseLogger, &oc.signerMap)
	if err != nil {
		return errors.Wrap(err, "syncSignerMap failed")
	}

	if added+removed > 0 {
		oc.logger.Info().
			Int("signer.added", added).
			Int("signer.removed", removed).
			Msg("synced signers")
	}

	return nil
}

func (oc *Orchestrator) runSyncObserverOperationalFlags(ctx context.Context) error {
	// every other block
	const cadence = 2 * constant.ZetaBlockTime

	task := func(ctx context.Context, _ *ticker.Ticker) error {
		if err := oc.syncObserverOperationalFlags(ctx); err != nil {
			oc.logger.Error().Err(err).Msg("syncObserverOperationalFlags failed")
		}

		return nil
	}

	return ticker.Run(
		ctx,
		cadence,
		task,
		ticker.WithLogger(oc.logger.Logger, "SyncObserverOperationalFlags"),
		ticker.WithStopChan(oc.stop),
	)
}

func (oc *Orchestrator) syncObserverOperationalFlags(ctx context.Context) error {
	client := oc.zetacoreClient
	flags, err := client.GetOperationalFlags(ctx)
	if err != nil {
		return fmt.Errorf("get operational flags: %w", err)
	}

	oc.mu.Lock()
	defer oc.mu.Unlock()
	newSignerBlockTimeOffsetPtr := flags.SignerBlockTimeOffset
	if newSignerBlockTimeOffsetPtr == nil {
		return nil
	}
	newSignerBlockTimeOffset := *newSignerBlockTimeOffsetPtr
	if oc.signerBlockTimeOffset != newSignerBlockTimeOffset {
		oc.logger.Info().
			Dur("offset", newSignerBlockTimeOffset).
			Msg("block time offset updated")
		oc.signerBlockTimeOffset = newSignerBlockTimeOffset
	}

	return nil
}
