// Package orchestrator provides the orchestrator for orchestrating cross-chain transactions
package orchestrator

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	zetamath "github.com/zeta-chain/zetacore/pkg/math"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outboundprocessor"
	"github.com/zeta-chain/zetacore/zetaclient/ratelimiter"
)

const (
	// evmOutboundLookbackFactor is the factor to determine how many nonces to look back for pending cctxs
	// For example, give OutboundScheduleLookahead of 120, pending NonceLow of 1000 and factor of 1.0,
	// the scheduler need to be able to pick up and schedule any pending cctx with nonce < 880 (1000 - 120 * 1.0)
	// NOTE: 1.0 means look back the same number of cctxs as we look ahead
	evmOutboundLookbackFactor = 1.0

	// sampling rate for sampled orchestrator logger
	loggerSamplingRate = 10
)

// Log is a struct that contains the logger
// TODO(revamp): rename to logger
type Log struct {
	// Base are the original base loggers used by orchestrator to create observers
	Base base.Logger

	// Std is the standard logger for orchestrator module
	Std zerolog.Logger

	// Sampled is the sampled logger for orchestrator module
	Sampled zerolog.Logger
}

// Orchestrator wraps the zetacore client, chain observers and signers. This is the high level object used for CCTX scheduling
type Orchestrator struct {
	// appContext contains the zetaclient application context
	appContext *context.AppContext

	// zetacore client
	zetacoreClient interfaces.ZetacoreClient

	// tss is the TSS signer
	tss interfaces.TSSSigner

	// signerMap contains the chain signers indexed by chainID
	signerMap map[int64]interfaces.ChainSigner

	// observerMap contains the chain observers indexed by chainID
	observerMap map[int64]interfaces.ChainObserver

	// outbound processor
	outboundProc *outboundprocessor.Processor

	// last operator balance
	lastOperatorBalance sdkmath.Int

	// logger contains the loggers used by the orchestrator
	logger Log

	// dbPath is the path observer database
	dbPath string

	// ts is the telemetry server for metrics
	ts *metrics.TelemetryServer

	// mu protects fields from concurrent access
	mu sync.Mutex

	// stop channel and flag to avoid closing twice
	stop    chan struct{}
	stopped bool
}

// NewOrchestrator creates a new orchestrator
func NewOrchestrator(
	appContext *context.AppContext,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	logger base.Logger,
	dbPath string,
	ts *metrics.TelemetryServer,
) *Orchestrator {
	oc := Orchestrator{
		appContext:     appContext,
		zetacoreClient: zetacoreClient,
		tss:            tss,
		signerMap:      make(map[int64]interfaces.ChainSigner),
		observerMap:    make(map[int64]interfaces.ChainObserver),
		dbPath:         dbPath,
		ts:             ts,
		mu:             sync.Mutex{},
		stop:           make(chan struct{}),
		stopped:        false,
	}

	// create loggers
	oc.logger = Log{
		Base: logger,
		Std:  logger.Std.With().Str("module", "orchestrator").Logger(),
	}
	oc.logger.Sampled = oc.logger.Std.Sample(&zerolog.BasicSampler{N: loggerSamplingRate})

	// create outbound processor
	oc.outboundProc = outboundprocessor.NewProcessor(logger.Std)

	// initialize hot key balance
	balance, err := oc.zetacoreClient.GetZetaHotKeyBalance()
	if err != nil {
		oc.logger.Std.Error().Err(err).Msg("error getting last balance of the hot key")
	}
	oc.lastOperatorBalance = balance

	return &oc
}

// Start all orchestrator routines
func (oc *Orchestrator) Start() {
	// watch for zetaclient app context changes
	go oc.WatchAppContext()

	// watch for upgrade plan in zetacore
	go oc.WatchUpgradePlan()

	// watch for enabling/disabling chains
	go oc.WatchEnabledChains()

	// schedule pending cctxs across all enabled chains
	go oc.SchedulePendingCctxs()

	// watch for stop signals
	oc.AwaitStopSignals()
}

// AwaitStopSignals waits for stop signals
func (oc *Orchestrator) AwaitStopSignals() {
	oc.logger.Std.Info().Msgf("Orchestrator awaiting the os.Interrupt, syscall.SIGTERM signals...")

	// subscribe to stop signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch

	// stop orchestrator
	oc.Stop()
	oc.logger.Std.Info().Msgf("Orchestrator stopped on signal: %s", sig)
}

// Stop notifies all zetaclient goroutines to stop
func (oc *Orchestrator) Stop() {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	// stop orchestrator only once.
	// both WatchUpgradePlan and system signals can trigger Stop()
	if !oc.stopped {
		// notify app context updater and CCTX scheduler to stop
		close(oc.stop)

		// notify all chain observers to stop
		for _, c := range oc.observerMap {
			c.Stop()
		}
		// set stopped flag
		oc.stopped = true
	}
}

// GetUpdatedSigner returns signer with updated chain parameters
func (oc *Orchestrator) GetUpdatedSigner(chainID int64) (interfaces.ChainSigner, error) {
	oc.mu.Lock()
	signer, found := oc.signerMap[chainID]
	oc.mu.Unlock()
	if !found {
		return nil, fmt.Errorf("signer not found for chainID %d", chainID)
	}

	// update signer parameters for the chain.
	// the logic is consistent for all chains, even if BTC chain doesn't have zetaConnector/erc20Custody.
	evmParams, found := oc.appContext.GetExternalChainParams(chainID)
	if found {
		// update zeta connector and ERC20 custody addresses
		zetaConnectorAddress := ethcommon.HexToAddress(evmParams.GetConnectorContractAddress())
		erc20CustodyAddress := ethcommon.HexToAddress(evmParams.GetErc20CustodyContractAddress())
		if zetaConnectorAddress != signer.GetZetaConnectorAddress() {
			signer.SetZetaConnectorAddress(zetaConnectorAddress)
			oc.logger.Std.Info().Msgf(
				"updated zeta connector address for chainID %d, new address: %s", chainID, zetaConnectorAddress)
		}
		if erc20CustodyAddress != signer.GetERC20CustodyAddress() {
			signer.SetERC20CustodyAddress(erc20CustodyAddress)
			oc.logger.Std.Info().Msgf(
				"updated ERC20 custody address for chainID %d, new address: %s", chainID, erc20CustodyAddress)
		}
	}
	return signer, nil
}

// GetUpdatedChainObserver returns chain observer with updated chain parameters
func (oc *Orchestrator) GetUpdatedChainObserver(chainID int64) (interfaces.ChainObserver, error) {
	oc.mu.Lock()
	observer, found := oc.observerMap[chainID]
	oc.mu.Unlock()
	if !found {
		return nil, fmt.Errorf("chain observer not found for chainID %d", chainID)
	}

	// update chain observer chain parameters
	oldParams := observer.GetChainParams()
	newParams, found := oc.appContext.GetExternalChainParams(chainID)
	if found && !observertypes.ChainParamsEqual(oldParams, *newParams) {
		observer.SetChainParams(*newParams)
		oc.logger.Std.Info().Msgf(
			"updated chain params for chainID %d, new params: %v", chainID, *newParams)
	}
	return observer, nil
}

// GetPendingCctxsWithinRatelimit get pending cctxs across foreign chains within rate limit
func (oc *Orchestrator) GetPendingCctxsWithinRatelimit(
	foreignChains []chains.Chain,
) (map[int64][]*types.CrossChainTx, error) {
	// get rate limiter flags
	rateLimitFlags, err := oc.zetacoreClient.GetRateLimiterFlags()
	if err != nil {
		return nil, err
	}

	// apply rate limiter or not according to the flags
	rateLimiterUsable := ratelimiter.IsRateLimiterUsable(rateLimitFlags)

	// fallback to non-rate-limited query if rate limiter is not usable
	cctxsMap := make(map[int64][]*types.CrossChainTx)
	if !rateLimiterUsable {
		for _, chain := range foreignChains {
			resp, _, err := oc.zetacoreClient.ListPendingCctx(chain.ChainId)
			if err == nil && resp != nil {
				cctxsMap[chain.ChainId] = resp
			}
		}
		return cctxsMap, nil
	}

	// query rate limiter input
	resp, err := oc.zetacoreClient.GetRateLimiterInput(rateLimitFlags.Window)
	if err != nil {
		return nil, err
	}
	input, ok := ratelimiter.NewInput(resp)
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

// SchedulePendingCctxs schedules keysigns for pending cctxs across all chains on ticker
// TODO(revamp): make this function simpler
func (oc *Orchestrator) SchedulePendingCctxs() {
	observeTicker := time.NewTicker(3 * time.Second)
	var lastBlockNum int64
	for {
		select {
		case <-oc.stop:
			oc.logger.Std.Warn().Msg("StartCctxScheduler: stopped")
			return
		case <-observeTicker.C:
			{
				bn, err := oc.zetacoreClient.GetBlockHeight()
				if err != nil {
					oc.logger.Std.Error().Err(err).Msg("StartCctxScheduler: GetBlockHeight fail")
					continue
				}
				if bn < 0 {
					oc.logger.Std.Error().Msg("StartCctxScheduler: GetBlockHeight returned negative height")
					continue
				}
				if lastBlockNum == 0 {
					lastBlockNum = bn - 1
				}
				if bn > lastBlockNum { // we have a new block
					bn = lastBlockNum + 1
					if bn%10 == 0 {
						oc.logger.Std.Debug().Msgf("StartCctxScheduler: zetacore heart beat: %d", bn)
					}

					balance, err := oc.zetacoreClient.GetZetaHotKeyBalance()
					if err != nil {
						oc.logger.Std.Error().Err(err).Msgf("couldn't get operator balance")
					} else {
						diff := oc.lastOperatorBalance.Sub(balance)
						if diff.GT(sdkmath.NewInt(0)) && diff.LT(sdkmath.NewInt(math.MaxInt64)) {
							oc.ts.AddFeeEntry(bn, diff.Int64())
							oc.lastOperatorBalance = balance
						}
					}

					// set current hot key burn rate
					metrics.HotKeyBurnRate.Set(float64(oc.ts.HotKeyBurnRate.GetBurnRate().Int64()))

					// get supported external chains
					externalChains := oc.appContext.GetEnabledExternalChains()

					// query pending cctxs across all external chains within rate limit
					cctxMap, err := oc.GetPendingCctxsWithinRatelimit(externalChains)
					if err != nil {
						oc.logger.Std.Error().Err(err).Msgf("StartCctxScheduler: GetPendingCctxsWithinRatelimit failed")
					}

					// schedule keysign for pending cctxs on each chain
					for _, c := range externalChains {
						// get cctxs from map and set pending transactions prometheus gauge
						cctxList := cctxMap[c.ChainId]
						metrics.PendingTxsPerChain.WithLabelValues(c.ChainName.String()).Set(float64(len(cctxList)))
						if len(cctxList) == 0 {
							continue
						}

						// update chain parameters for signer and chain observer
						signer, err := oc.GetUpdatedSigner(c.ChainId)
						if err != nil {
							oc.logger.Std.Error().
								Err(err).
								Msgf("StartCctxScheduler: GetUpdatedSigner failed for chain %d", c.ChainId)
							continue
						}
						ob, err := oc.GetUpdatedChainObserver(c.ChainId)
						if err != nil {
							oc.logger.Std.Error().
								Err(err).
								Msgf("StartCctxScheduler: GetUpdatedChainObserver failed for chain %d", c.ChainId)
							continue
						}
						if !oc.appContext.IsOutboundObservationEnabled(ob.GetChainParams()) {
							continue
						}

						// #nosec G701 range is verified
						zetaHeight := uint64(bn)
						if chains.IsEVMChain(c.ChainId) {
							oc.ScheduleCctxEVM(zetaHeight, c.ChainId, cctxList, ob, signer)
						} else if chains.IsBitcoinChain(c.ChainId) {
							oc.ScheduleCctxBTC(zetaHeight, c.ChainId, cctxList, ob, signer)
						} else {
							oc.logger.Std.Error().Msgf("StartCctxScheduler: unsupported chain %d", c.ChainId)
							continue
						}
					}

					// update last processed block number
					lastBlockNum = bn
					oc.ts.SetCoreBlockNumber(lastBlockNum)
				}
			}
		}
	}
}

// ScheduleCctxEVM schedules evm outbound keysign on each ZetaChain block (the ticker)
func (oc *Orchestrator) ScheduleCctxEVM(
	zetaHeight uint64,
	chainID int64,
	cctxList []*types.CrossChainTx,
	observer interfaces.ChainObserver,
	signer interfaces.ChainSigner,
) {
	res, err := oc.zetacoreClient.GetAllOutboundTrackerByChain(chainID, interfaces.Ascending)
	if err != nil {
		oc.logger.Std.Warn().Err(err).Msgf("ScheduleCctxEVM: GetAllOutboundTrackerByChain failed for chain %d", chainID)
		return
	}
	trackerMap := make(map[uint64]bool)
	for _, v := range res {
		trackerMap[v.Nonce] = true
	}
	outboundScheduleLookahead := observer.GetChainParams().OutboundScheduleLookahead
	// #nosec G701 always in range
	outboundScheduleLookback := uint64(float64(outboundScheduleLookahead) * evmOutboundLookbackFactor)
	// #nosec G701 positive
	outboundScheduleInterval := uint64(observer.GetChainParams().OutboundScheduleInterval)

	for idx, cctx := range cctxList {
		params := cctx.GetCurrentOutboundParam()
		nonce := params.TssNonce
		outboundID := outboundprocessor.ToOutboundID(cctx.Index, params.ReceiverChainId, nonce)

		if params.ReceiverChainId != chainID {
			oc.logger.Std.Error().
				Msgf("ScheduleCctxEVM: outbound %s chainid mismatch: want %d, got %d", outboundID, chainID, params.ReceiverChainId)
			continue
		}
		if params.TssNonce > cctxList[0].GetCurrentOutboundParam().TssNonce+outboundScheduleLookback {
			oc.logger.Std.Error().Msgf("ScheduleCctxEVM: nonce too high: signing %d, earliest pending %d, chain %d",
				params.TssNonce, cctxList[0].GetCurrentOutboundParam().TssNonce, chainID)
			break
		}

		// try confirming the outbound
		included, _, err := observer.IsOutboundProcessed(cctx, oc.logger.Std)
		if err != nil {
			oc.logger.Std.Error().
				Err(err).
				Msgf("ScheduleCctxEVM: IsOutboundProcessed faild for chain %d nonce %d", chainID, nonce)
			continue
		}
		if included {
			oc.logger.Std.Info().
				Msgf("ScheduleCctxEVM: outbound %s already included; do not schedule keysign", outboundID)
			continue
		}

		// determining critical outbound; if it satisfies following criteria
		// 1. it's the first pending outbound for this chain
		// 2. the following 5 nonces have been in tracker
		criticalInterval := uint64(10)                      // for critical pending outbound we reduce re-try interval
		nonCriticalInterval := outboundScheduleInterval * 2 // for non-critical pending outbound we increase re-try interval
		if nonce%criticalInterval == zetaHeight%criticalInterval {
			count := 0
			for i := nonce + 1; i <= nonce+10; i++ {
				if _, found := trackerMap[i]; found {
					count++
				}
			}
			if count >= 5 {
				outboundScheduleInterval = criticalInterval
			}
		}
		// if it's already in tracker, we increase re-try interval
		if _, ok := trackerMap[nonce]; ok {
			outboundScheduleInterval = nonCriticalInterval
		}

		// otherwise, the normal interval is used
		if nonce%outboundScheduleInterval == zetaHeight%outboundScheduleInterval &&
			!oc.outboundProc.IsOutboundActive(outboundID) {
			oc.outboundProc.StartTryProcess(outboundID)
			oc.logger.Std.Debug().
				Msgf("ScheduleCctxEVM: sign outbound %s with value %d\n", outboundID, cctx.GetCurrentOutboundParam().Amount)
			go signer.TryProcessOutbound(cctx, oc.outboundProc, outboundID, observer, oc.zetacoreClient, zetaHeight)
		}

		// #nosec G701 always in range
		if int64(idx) >= outboundScheduleLookahead-1 { // only look at 'lookahead' cctxs per chain
			break
		}
	}
}

// ScheduleCctxBTC schedules bitcoin outbound keysign on each ZetaChain block (the ticker)
// 1. schedule at most one keysign per ticker
// 2. schedule keysign only when nonce-mark UTXO is available
// 3. stop keysign when lookahead is reached
func (oc *Orchestrator) ScheduleCctxBTC(
	zetaHeight uint64,
	chainID int64,
	cctxList []*types.CrossChainTx,
	observer interfaces.ChainObserver,
	signer interfaces.ChainSigner,
) {
	btcObserver, ok := observer.(*btcobserver.Observer)
	if !ok { // should never happen
		oc.logger.Std.Error().Msgf("ScheduleCctxBTC: chain observer is not a bitcoin observer")
		return
	}
	// #nosec G701 positive
	interval := uint64(observer.GetChainParams().OutboundScheduleInterval)
	lookahead := observer.GetChainParams().OutboundScheduleLookahead

	// schedule at most one keysign per ticker
	for idx, cctx := range cctxList {
		params := cctx.GetCurrentOutboundParam()
		nonce := params.TssNonce
		outboundID := outboundprocessor.ToOutboundID(cctx.Index, params.ReceiverChainId, nonce)

		if params.ReceiverChainId != chainID {
			oc.logger.Std.Error().
				Msgf("ScheduleCctxBTC: outbound %s chainid mismatch: want %d, got %d", outboundID, chainID, params.ReceiverChainId)
			continue
		}
		// try confirming the outbound
		included, confirmed, err := btcObserver.IsOutboundProcessed(cctx, oc.logger.Std)
		if err != nil {
			oc.logger.Std.Error().
				Err(err).
				Msgf("ScheduleCctxBTC: IsOutboundProcessed faild for chain %d nonce %d", chainID, nonce)
			continue
		}
		if included || confirmed {
			oc.logger.Std.Info().
				Msgf("ScheduleCctxBTC: outbound %s already included; do not schedule keysign", outboundID)
			continue
		}

		// stop if the nonce being processed is higher than the pending nonce
		if nonce > btcObserver.GetPendingNonce() {
			break
		}
		// stop if lookahead is reached
		if int64(
			idx,
		) >= lookahead { // 2 bitcoin confirmations span is 20 minutes on average. We look ahead up to 100 pending cctx to target TPM of 5.
			oc.logger.Std.Warn().
				Msgf("ScheduleCctxBTC: lookahead reached, signing %d, earliest pending %d", nonce, cctxList[0].GetCurrentOutboundParam().TssNonce)
			break
		}
		// try confirming the outbound or scheduling a keysign
		if nonce%interval == zetaHeight%interval && !oc.outboundProc.IsOutboundActive(outboundID) {
			oc.outboundProc.StartTryProcess(outboundID)
			oc.logger.Std.Debug().Msgf("ScheduleCctxBTC: sign outbound %s with value %d\n", outboundID, params.Amount)
			go signer.TryProcessOutbound(cctx, oc.outboundProc, outboundID, observer, oc.zetacoreClient, zetaHeight)
		}
	}
}
