// Package orchestrator provides the orchestrator for orchestrating cross-chain transactions
package orchestrator

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/bg"
	"github.com/zeta-chain/zetacore/pkg/chains"
	zetamath "github.com/zeta-chain/zetacore/pkg/math"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	solanaobserver "github.com/zeta-chain/zetacore/zetaclient/chains/solana/observer"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
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

var defaultLogSampler = &zerolog.BasicSampler{N: loggerSamplingRate}

// Orchestrator wraps the zetacore client, chain observers and signers. This is the high level object used for CCTX scheduling
type Orchestrator struct {
	// zetacore client
	zetacoreClient interfaces.ZetacoreClient

	// signerMap contains the chain signers indexed by chainID
	signerMap map[int64]interfaces.ChainSigner

	// observerMap contains the chain observers indexed by chainID
	observerMap map[int64]interfaces.ChainObserver

	// outbound processor
	outboundProc *outboundprocessor.Processor

	// last operator balance
	lastOperatorBalance sdkmath.Int

	// observer & signer props
	tss         interfaces.TSSSigner
	dbDirectory string
	baseLogger  base.Logger

	// misc
	logger multiLogger
	ts     *metrics.TelemetryServer
	stop   chan struct{}
	mu     sync.RWMutex
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

	log := multiLogger{
		Logger:  logger.Std.With().Str("module", "orchestrator").Logger(),
		Sampled: logger.Std.With().Str("module", "orchestrator").Logger().Sample(defaultLogSampler),
	}

	balance, err := client.GetZetaHotKeyBalance(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get last balance of the hot key")
	}

	return &Orchestrator{
		zetacoreClient: client,

		signerMap:   signerMap,
		observerMap: observerMap,

		outboundProc:        outboundprocessor.NewProcessor(logger.Std),
		lastOperatorBalance: balance,

		// observer & signer props
		tss:         tss,
		dbDirectory: dbDirectory,
		baseLogger:  logger,

		logger: log,
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

	// start cctx scheduler
	bg.Work(ctx, oc.runScheduler, bg.WithName("runScheduler"), bg.WithLogger(oc.logger.Logger))
	bg.Work(ctx, oc.runObserverSignerSync, bg.WithName("runObserverSignerSync"), bg.WithLogger(oc.logger.Logger))

	shutdownOrchestrator := func() {
		// now stop orchestrator and all observers
		close(oc.stop)
		for _, c := range oc.observerMap {
			c.Stop()
		}
	}

	oc.zetacoreClient.OnBeforeStop(shutdownOrchestrator)

	return nil
}

// returns signer with updated chain parameters.
func (oc *Orchestrator) resolveSigner(app *zctx.AppContext, chainID int64) (interfaces.ChainSigner, error) {
	signer, err := oc.getSigner(chainID)
	if err != nil {
		return nil, err
	}

	// update signer chain parameters
	if chains.IsEVMChain(chainID, app.GetAdditionalChains()) {
		evmParams, found := app.GetEVMChainParams(chainID)
		if found {
			// update zeta connector and ERC20 custody addresses
			zetaConnectorAddress := ethcommon.HexToAddress(evmParams.GetConnectorContractAddress())
			if zetaConnectorAddress != signer.GetZetaConnectorAddress() {
				signer.SetZetaConnectorAddress(zetaConnectorAddress)
				oc.logger.Info().
					Str("signer.connector_address", zetaConnectorAddress.String()).
					Msgf("updated zeta connector address for chain %d", chainID)
			}

			erc20CustodyAddress := ethcommon.HexToAddress(evmParams.GetErc20CustodyContractAddress())
			if erc20CustodyAddress != signer.GetERC20CustodyAddress() {
				signer.SetERC20CustodyAddress(erc20CustodyAddress)
				oc.logger.Info().
					Str("signer.erc20_custody", erc20CustodyAddress.String()).
					Msgf("updated zeta connector address for chain %d", chainID)
			}
		}
	} else if chains.IsSolanaChain(chainID, app.GetAdditionalChains()) {
		_, solParams, found := app.GetSolanaChainParams()
		if found {
			// update solana gateway address
			if solParams.GatewayAddress != signer.GetGatewayAddress() {
				signer.SetGatewayAddress(solParams.GatewayAddress)
				oc.logger.Info().
					Str("signer.gateway_address", solParams.GatewayAddress).
					Msgf("updated gateway address for chain %d", chainID)
			}
		}
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

	// update chain observer chain parameters
	curParams := observer.GetChainParams()
	if chains.IsEVMChain(chainID, app.GetAdditionalChains()) {
		evmParams, found := app.GetEVMChainParams(chainID)
		if found && !observertypes.ChainParamsEqual(curParams, *evmParams) {
			observer.SetChainParams(*evmParams)
			oc.logger.Info().
				Interface("observer.chain_params", *evmParams).
				Msgf("updated chain params for EVM chainID %d", chainID)
		}
	} else if chains.IsBitcoinChain(chainID, app.GetAdditionalChains()) {
		_, btcParams, found := app.GetBTCChainParams()
		if found && !observertypes.ChainParamsEqual(curParams, *btcParams) {
			observer.SetChainParams(*btcParams)
			oc.logger.Info().
				Interface("observer.chain_params", *btcParams).
				Msgf("updated chain params for UTXO chainID %d", btcParams.ChainId)
		}
	} else if chains.IsSolanaChain(chainID, app.GetAdditionalChains()) {
		_, solParams, found := app.GetSolanaChainParams()
		if found && !observertypes.ChainParamsEqual(curParams, *solParams) {
			observer.SetChainParams(*solParams)
			oc.logger.Info().Msgf(
				"updated chain params for Solana, new params: %v", *solParams)
		}
	}

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
func (oc *Orchestrator) GetPendingCctxsWithinRateLimit(
	ctx context.Context,
	foreignChains []chains.Chain,
) (map[int64][]*types.CrossChainTx, error) {
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
		for _, chain := range foreignChains {
			resp, _, err := oc.zetacoreClient.ListPendingCCTX(ctx, chain.ChainId)
			if err == nil && resp != nil {
				cctxsMap[chain.ChainId] = resp
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

	observeTicker := time.NewTicker(3 * time.Second)
	var lastBlockNum int64
	for {
		select {
		case <-oc.stop:
			oc.logger.Warn().Msg("runScheduler: stopped")
			return nil
		case <-observeTicker.C:
			{
				bn, err := oc.zetacoreClient.GetBlockHeight(ctx)
				if err != nil {
					oc.logger.Error().Err(err).Msg("StartCctxScheduler: GetBlockHeight fail")
					continue
				}
				if bn < 0 {
					oc.logger.Error().Msg("runScheduler: GetBlockHeight returned negative height")
					continue
				}
				if lastBlockNum == 0 {
					lastBlockNum = bn - 1
				}
				if bn > lastBlockNum { // we have a new block
					bn = lastBlockNum + 1
					if bn%10 == 0 {
						oc.logger.Debug().Msgf("runScheduler: zetacore heart beat: %d", bn)
					}

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

					// get supported external chains
					externalChains := app.GetEnabledExternalChains()

					// query pending cctxs across all external chains within rate limit
					cctxMap, err := oc.GetPendingCctxsWithinRateLimit(ctx, externalChains)
					if err != nil {
						oc.logger.Error().Err(err).Msgf("runScheduler: GetPendingCctxsWithinRatelimit failed")
					}

					// schedule keysign for pending cctxs on each chain
					for _, c := range externalChains {
						// get cctxs from map and set pending transactions prometheus gauge
						cctxList := cctxMap[c.ChainId]
						metrics.PendingTxsPerChain.WithLabelValues(c.Name).Set(float64(len(cctxList)))
						if len(cctxList) == 0 {
							continue
						}

						// update chain parameters for signer and chain observer
						signer, err := oc.resolveSigner(app, c.ChainId)
						if err != nil {
							oc.logger.Error().Err(err).
								Msgf("runScheduler: unable to resolve signer for chain %d", c.ChainId)
							continue
						}
						ob, err := oc.resolveObserver(app, c.ChainId)
						if err != nil {
							oc.logger.Error().Err(err).
								Msgf("runScheduler: resolveObserver failed for chain %d", c.ChainId)
							continue
						}
						if !app.IsOutboundObservationEnabled(ob.GetChainParams()) {
							continue
						}

						// #nosec G115 range is verified
						zetaHeight := uint64(bn)
						if chains.IsEVMChain(c.ChainId, app.GetAdditionalChains()) {
							oc.ScheduleCctxEVM(ctx, zetaHeight, c.ChainId, cctxList, ob, signer)
						} else if chains.IsBitcoinChain(c.ChainId, app.GetAdditionalChains()) {
							oc.ScheduleCctxBTC(ctx, zetaHeight, c.ChainId, cctxList, ob, signer)
						} else if chains.IsSolanaChain(c.ChainId, app.GetAdditionalChains()) {
							oc.ScheduleCctxSolana(ctx, zetaHeight, c.ChainId, cctxList, ob, signer)
						} else {
							oc.logger.Error().Msgf("runScheduler: unsupported chain %d", c.ChainId)
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
	ctx context.Context,
	zetaHeight uint64,
	chainID int64,
	cctxList []*types.CrossChainTx,
	observer interfaces.ChainObserver,
	signer interfaces.ChainSigner,
) {
	res, err := oc.zetacoreClient.GetAllOutboundTrackerByChain(ctx, chainID, interfaces.Ascending)
	if err != nil {
		oc.logger.Warn().Err(err).Msgf("ScheduleCctxEVM: GetAllOutboundTrackerByChain failed for chain %d", chainID)
		return
	}
	trackerMap := make(map[uint64]bool)
	for _, v := range res {
		trackerMap[v.Nonce] = true
	}
	outboundScheduleLookahead := observer.GetChainParams().OutboundScheduleLookahead
	// #nosec G115 always in range
	outboundScheduleLookback := uint64(float64(outboundScheduleLookahead) * evmOutboundLookbackFactor)
	// #nosec G115 positive
	outboundScheduleInterval := uint64(observer.GetChainParams().OutboundScheduleInterval)

	for idx, cctx := range cctxList {
		params := cctx.GetCurrentOutboundParam()
		nonce := params.TssNonce
		outboundID := outboundprocessor.ToOutboundID(cctx.Index, params.ReceiverChainId, nonce)

		if params.ReceiverChainId != chainID {
			oc.logger.Error().
				Msgf("ScheduleCctxEVM: outbound %s chainid mismatch: want %d, got %d", outboundID, chainID, params.ReceiverChainId)
			continue
		}
		if params.TssNonce > cctxList[0].GetCurrentOutboundParam().TssNonce+outboundScheduleLookback {
			oc.logger.Error().Msgf("ScheduleCctxEVM: nonce too high: signing %d, earliest pending %d, chain %d",
				params.TssNonce, cctxList[0].GetCurrentOutboundParam().TssNonce, chainID)
			break
		}

		// try confirming the outbound
		_, confirmed, err := observer.IsOutboundProcessed(ctx, cctx)
		if err != nil {
			oc.logger.Error().
				Err(err).
				Msgf("ScheduleCctxEVM: IsOutboundProcessed failed for chain %d nonce %d", chainID, nonce)
			continue
		}
		if confirmed {
			oc.logger.Info().
				Msgf("ScheduleCctxEVM: outbound %s already confirmed; do not schedule keysign", outboundID)
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
			oc.logger.Debug().
				Msgf("ScheduleCctxEVM: sign outbound %s with value %d", outboundID, cctx.GetCurrentOutboundParam().Amount)
			go signer.TryProcessOutbound(
				ctx,
				cctx,
				oc.outboundProc,
				outboundID,
				observer,
				oc.zetacoreClient,
				zetaHeight,
			)
		}

		// #nosec G115 always in range
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
	ctx context.Context,
	zetaHeight uint64,
	chainID int64,
	cctxList []*types.CrossChainTx,
	observer interfaces.ChainObserver,
	signer interfaces.ChainSigner,
) {
	btcObserver, ok := observer.(*btcobserver.Observer)
	if !ok { // should never happen
		oc.logger.Error().Msgf("ScheduleCctxBTC: chain observer is not a bitcoin observer")
		return
	}
	// #nosec G115 positive
	interval := uint64(observer.GetChainParams().OutboundScheduleInterval)
	lookahead := observer.GetChainParams().OutboundScheduleLookahead

	// schedule at most one keysign per ticker
	for idx, cctx := range cctxList {
		params := cctx.GetCurrentOutboundParam()
		nonce := params.TssNonce
		outboundID := outboundprocessor.ToOutboundID(cctx.Index, params.ReceiverChainId, nonce)

		if params.ReceiverChainId != chainID {
			oc.logger.Error().
				Msgf("ScheduleCctxBTC: outbound %s chainid mismatch: want %d, got %d", outboundID, chainID, params.ReceiverChainId)
			continue
		}
		// try confirming the outbound
		included, confirmed, err := btcObserver.IsOutboundProcessed(ctx, cctx)
		if err != nil {
			oc.logger.Error().
				Err(err).
				Msgf("ScheduleCctxBTC: IsOutboundProcessed failed for chain %d nonce %d", chainID, nonce)
			continue
		}
		if included || confirmed {
			oc.logger.Info().
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
			oc.logger.Warn().
				Msgf("ScheduleCctxBTC: lookahead reached, signing %d, earliest pending %d", nonce, cctxList[0].GetCurrentOutboundParam().TssNonce)
			break
		}
		// schedule a TSS keysign
		if nonce%interval == zetaHeight%interval && !oc.outboundProc.IsOutboundActive(outboundID) {
			oc.outboundProc.StartTryProcess(outboundID)
			oc.logger.Debug().Msgf("ScheduleCctxBTC: sign outbound %s with value %d", outboundID, params.Amount)
			go signer.TryProcessOutbound(
				ctx,
				cctx,
				oc.outboundProc,
				outboundID,
				observer,
				oc.zetacoreClient,
				zetaHeight,
			)
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
	if !ok { // should never happen
		oc.logger.Error().Msgf("ScheduleCctxSolana: chain observer is not a solana observer")
		return
	}
	// #nosec G701 positive
	interval := uint64(observer.GetChainParams().OutboundScheduleInterval)

	// schedule keysign for each pending cctx
	for _, cctx := range cctxList {
		params := cctx.GetCurrentOutboundParam()
		nonce := params.TssNonce
		outboundID := outboundprocessor.ToOutboundID(cctx.Index, params.ReceiverChainId, nonce)

		if params.ReceiverChainId != chainID {
			oc.logger.Error().
				Msgf("ScheduleCctxSolana: outbound %s chainid mismatch: want %d, got %d", outboundID, chainID, params.ReceiverChainId)
			continue
		}

		// try confirming the outbound
		_, finalized, err := solObserver.IsOutboundProcessed(ctx, cctx)
		if err != nil {
			oc.logger.Error().
				Err(err).
				Msgf("ScheduleCctxSolana: IsOutboundProcessed failed for chain %d nonce %d", chainID, nonce)
			continue
		}
		if finalized {
			oc.logger.Info().
				Msgf("ScheduleCctxSolana: outbound %s already finalized; do not schedule keysign", outboundID)
			continue
		}

		// schedule a TSS keysign
		if nonce%interval == zetaHeight%interval && !oc.outboundProc.IsOutboundActive(outboundID) {
			oc.outboundProc.StartTryProcess(outboundID)
			oc.logger.Debug().Msgf("ScheduleCctxSolana: sign outbound %s with value %d", outboundID, params.Amount)
			go signer.TryProcessOutbound(
				ctx,
				cctx,
				oc.outboundProc,
				outboundID,
				observer,
				oc.zetacoreClient,
				zetaHeight,
			)
		}
	}
}

// runObserverSignerSync runs a blocking ticker that observes chain changes from zetacore
// and optionally (de)provisions respective observers and signers.
func (oc *Orchestrator) runObserverSignerSync(ctx context.Context) error {
	// check every other zeta block
	const cadence = 2 * evm.ZetaBlockTime

	ticker := time.NewTicker(cadence)
	defer ticker.Stop()

	for {
		select {
		case <-oc.stop:
			oc.logger.Warn().Msg("runObserverSignerSync: stopped")
			return nil
		case <-ticker.C:
			if err := oc.syncObserverSigner(ctx); err != nil {
				oc.logger.Error().Err(err).Msg("runObserverSignerSync: syncObserverSigner failed")
			}
		}
	}
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

	added, removed, err = syncSignerMap(ctx, oc.tss, oc.baseLogger, oc.ts, &oc.signerMap)
	if err != nil {
		return errors.Wrap(err, "syncSignerMap failed")
	}

	if added+removed > 0 {
		oc.logger.Info().
			Int("signers.added", added).
			Int("signers.removed", removed).
			Msg("synced signers")
	}

	return nil
}
