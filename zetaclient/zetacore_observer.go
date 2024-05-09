package zetaclient

import (
	"fmt"
	"math"
	"time"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/pkg/chains"
	zetamath "github.com/zeta-chain/zetacore/pkg/math"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	appcontext "github.com/zeta-chain/zetacore/zetaclient/app_context"
	"github.com/zeta-chain/zetacore/zetaclient/bitcoin"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outboundprocessor"
	"github.com/zeta-chain/zetacore/zetaclient/ratelimiter"
)

const (
	// evmOutboundTxLookbackFactor is the factor to determine how many nonces to look back for pending cctxs
	// For example, give OutboundTxScheduleLookahead of 120, pending NonceLow of 1000 and factor of 1.0,
	// the scheduler need to be able to pick up and schedule any pending cctx with nonce < 880 (1000 - 120 * 1.0)
	// NOTE: 1.0 means look back the same number of cctxs as we look ahead
	evmOutboundTxLookbackFactor = 1.0

	// sampling rate for sampled observer logger
	loggerSamplingRate = 10
)

type ZetaCoreLog struct {
	Observer          zerolog.Logger
	ObserverSampled   zerolog.Logger
	OutboundProcessor zerolog.Logger
}

// CoreObserver wraps the zetabridge, chain clients and signers. This is the high level object used for CCTX scheduling
type CoreObserver struct {
	bridge              interfaces.ZetaCoreBridger
	signerMap           map[int64]interfaces.ChainSigner
	clientMap           map[int64]interfaces.ChainClient
	logger              ZetaCoreLog
	ts                  *metrics.TelemetryServer
	stop                chan struct{}
	lastOperatorBalance sdkmath.Int
}

// NewCoreObserver creates a new CoreObserver
func NewCoreObserver(
	bridge interfaces.ZetaCoreBridger,
	signerMap map[int64]interfaces.ChainSigner,
	clientMap map[int64]interfaces.ChainClient,
	logger zerolog.Logger,
	ts *metrics.TelemetryServer,
) *CoreObserver {
	co := CoreObserver{
		ts:   ts,
		stop: make(chan struct{}),
	}

	// create loggers
	chainLogger := logger.With().Str("chain", "ZetaChain").Logger()
	co.logger = ZetaCoreLog{
		Observer:          chainLogger.With().Str("module", "Observer").Logger(),
		OutboundProcessor: chainLogger.With().Str("module", "OutboundProcessor").Logger(),
	}
	co.logger.ObserverSampled = co.logger.Observer.Sample(&zerolog.BasicSampler{N: loggerSamplingRate})

	// set bridge, signers and clients
	co.bridge = bridge
	co.signerMap = signerMap
	co.clientMap = clientMap
	co.logger.Observer.Info().Msg("starting core observer")
	balance, err := bridge.GetZetaHotKeyBalance()
	if err != nil {
		co.logger.Observer.Error().Err(err).Msg("error getting last balance of the hot key")
	}
	co.lastOperatorBalance = balance

	return &co
}

func (co *CoreObserver) MonitorCore(appContext *appcontext.AppContext) {
	myid := co.bridge.GetKeys().GetAddress()
	co.logger.Observer.Info().Msgf("Starting cctx scheduler for %s", myid)
	go co.StartCctxScheduler(appContext)

	go func() {
		// bridge queries UpgradePlan from zetabridge and send to its pause channel if upgrade height is reached
		co.bridge.Pause()
		// now stop everything
		close(co.stop) // this stops the startSendScheduler() loop
		for _, c := range co.clientMap {
			c.Stop()
		}
	}()
}

// GetUpdatedSigner returns signer with updated chain parameters
func (co *CoreObserver) GetUpdatedSigner(coreContext *corecontext.ZetaCoreContext, chainID int64) (interfaces.ChainSigner, error) {
	signer, found := co.signerMap[chainID]
	if !found {
		return nil, fmt.Errorf("signer not found for chainID %d", chainID)
	}
	// update EVM signer parameters only. BTC signer doesn't use chain parameters for now.
	if chains.IsEVMChain(chainID) {
		evmParams, found := coreContext.GetEVMChainParams(chainID)
		if found {
			// update zeta connector and ERC20 custody addresses
			zetaConnectorAddress := ethcommon.HexToAddress(evmParams.GetConnectorContractAddress())
			erc20CustodyAddress := ethcommon.HexToAddress(evmParams.GetErc20CustodyContractAddress())
			if zetaConnectorAddress != signer.GetZetaConnectorAddress() {
				signer.SetZetaConnectorAddress(zetaConnectorAddress)
				co.logger.Observer.Info().Msgf(
					"updated zeta connector address for chainID %d, new address: %s", chainID, zetaConnectorAddress)
			}
			if erc20CustodyAddress != signer.GetERC20CustodyAddress() {
				signer.SetERC20CustodyAddress(erc20CustodyAddress)
				co.logger.Observer.Info().Msgf(
					"updated ERC20 custody address for chainID %d, new address: %s", chainID, erc20CustodyAddress)
			}
		}
	}
	return signer, nil
}

// GetUpdatedChainClient returns chain client object with updated chain parameters
func (co *CoreObserver) GetUpdatedChainClient(coreContext *corecontext.ZetaCoreContext, chainID int64) (interfaces.ChainClient, error) {
	chainOb, found := co.clientMap[chainID]
	if !found {
		return nil, fmt.Errorf("chain client not found for chainID %d", chainID)
	}
	// update chain client chain parameters
	curParams := chainOb.GetChainParams()
	if chains.IsEVMChain(chainID) {
		evmParams, found := coreContext.GetEVMChainParams(chainID)
		if found && !observertypes.ChainParamsEqual(curParams, *evmParams) {
			chainOb.SetChainParams(*evmParams)
			co.logger.Observer.Info().Msgf(
				"updated chain params for chainID %d, new params: %v", chainID, *evmParams)
		}
	} else if chains.IsBitcoinChain(chainID) {
		_, btcParams, found := coreContext.GetBTCChainParams()

		if found && !observertypes.ChainParamsEqual(curParams, *btcParams) {
			chainOb.SetChainParams(*btcParams)
			co.logger.Observer.Info().Msgf(
				"updated chain params for Bitcoin, new params: %v", *btcParams)
		}
	}
	return chainOb, nil
}

// GetPendingCctxsWithinRatelimit get pending cctxs across foreign chains within rate limit
func (co *CoreObserver) GetPendingCctxsWithinRatelimit(foreignChains []chains.Chain) (map[int64][]*types.CrossChainTx, error) {
	// get rate limiter flags
	rateLimitFlags, err := co.bridge.GetRateLimiterFlags()
	if err != nil {
		return nil, err
	}

	// apply rate limiter or not according to the flags
	rateLimiterUsable := ratelimiter.IsRateLimiterUsable(rateLimitFlags)

	// fallback to non-rate-limited query if rate limiter is not usable
	cctxsMap := make(map[int64][]*types.CrossChainTx)
	if !rateLimiterUsable {
		for _, chain := range foreignChains {
			resp, _, err := co.bridge.ListPendingCctx(chain.ChainId)
			if err == nil && resp != nil {
				cctxsMap[chain.ChainId] = resp
			}
		}
		return cctxsMap, nil
	}

	// query rate limiter input
	resp, err := co.bridge.GetRateLimiterInput(rateLimitFlags.Window)
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
		co.logger.ObserverSampled.Info().Msgf("current rate limiter window: %d rate: %s, percentage: %f",
			output.CurrentWithdrawWindow, output.CurrentWithdrawRate.String(), percentageFloat)
	}

	return output.CctxsMap, nil
}

// StartCctxScheduler schedules keysigns for cctxs on each ZetaChain block (the ticker)
func (co *CoreObserver) StartCctxScheduler(appContext *appcontext.AppContext) {
	outboundManager := outboundprocessor.NewOutboundProcessorManager(co.logger.OutboundProcessor)
	observeTicker := time.NewTicker(3 * time.Second)
	var lastBlockNum int64
	for {
		select {
		case <-co.stop:
			co.logger.Observer.Warn().Msg("StartCctxScheduler: stopped")
			return
		case <-observeTicker.C:
			{
				bn, err := co.bridge.GetZetaBlockHeight()
				if err != nil {
					co.logger.Observer.Error().Err(err).Msg("StartCctxScheduler: GetZetaBlockHeight fail")
					continue
				}
				if bn < 0 {
					co.logger.Observer.Error().Msg("StartCctxScheduler: GetZetaBlockHeight returned negative height")
					continue
				}
				if lastBlockNum == 0 {
					lastBlockNum = bn - 1
				}
				if bn > lastBlockNum { // we have a new block
					bn = lastBlockNum + 1
					if bn%10 == 0 {
						co.logger.Observer.Debug().Msgf("StartCctxScheduler: ZetaCore heart beat: %d", bn)
					}

					balance, err := co.bridge.GetZetaHotKeyBalance()
					if err != nil {
						co.logger.Observer.Error().Err(err).Msgf("couldn't get operator balance")
					} else {
						diff := co.lastOperatorBalance.Sub(balance)
						if diff.GT(sdkmath.NewInt(0)) && diff.LT(sdkmath.NewInt(math.MaxInt64)) {
							co.ts.AddFeeEntry(bn, diff.Int64())
							co.lastOperatorBalance = balance
						}
					}

					// set current hot key burn rate
					metrics.HotKeyBurnRate.Set(float64(co.ts.HotKeyBurnRate.GetBurnRate().Int64()))

					// get supported external chains
					coreContext := appContext.ZetaCoreContext()
					externalChains := coreContext.GetEnabledExternalChains()

					// query pending cctxs across all external chains within rate limit
					cctxMap, err := co.GetPendingCctxsWithinRatelimit(externalChains)
					if err != nil {
						co.logger.Observer.Error().Err(err).Msgf("StartCctxScheduler: GetPendingCctxsWithinRatelimit failed")
					}

					// schedule keysign for pending cctxs on each chain
					for _, c := range externalChains {
						// get cctxs from map and set pending transactions prometheus gauge
						cctxList := cctxMap[c.ChainId]
						metrics.PendingTxsPerChain.WithLabelValues(c.ChainName.String()).Set(float64(len(cctxList)))
						if len(cctxList) == 0 {
							continue
						}

						// update chain parameters for signer and chain client
						signer, err := co.GetUpdatedSigner(coreContext, c.ChainId)
						if err != nil {
							co.logger.Observer.Error().Err(err).Msgf("StartCctxScheduler: GetUpdatedSigner failed for chain %d", c.ChainId)
							continue
						}
						ob, err := co.GetUpdatedChainClient(coreContext, c.ChainId)
						if err != nil {
							co.logger.Observer.Error().Err(err).Msgf("StartCctxScheduler: GetUpdatedChainClient failed for chain %d", c.ChainId)
							continue
						}
						if !corecontext.IsOutboundObservationEnabled(coreContext, ob.GetChainParams()) {
							continue
						}

						// #nosec G701 range is verified
						zetaHeight := uint64(bn)
						if chains.IsEVMChain(c.ChainId) {
							co.ScheduleCctxEVM(outboundManager, zetaHeight, c.ChainId, cctxList, ob, signer)
						} else if chains.IsBitcoinChain(c.ChainId) {
							co.ScheduleCctxBTC(outboundManager, zetaHeight, c.ChainId, cctxList, ob, signer)
						} else {
							co.logger.Observer.Error().Msgf("StartCctxScheduler: unsupported chain %d", c.ChainId)
							continue
						}
					}

					// update last processed block number
					lastBlockNum = bn
					metrics.LastCoreBlockNumber.Set(float64(lastBlockNum))
				}
			}
		}
	}
}

// ScheduleCctxEVM schedules evm outtx keysign on each ZetaChain block (the ticker)
func (co *CoreObserver) ScheduleCctxEVM(
	outboundManager *outboundprocessor.Processor,
	zetaHeight uint64,
	chainID int64,
	cctxList []*types.CrossChainTx,
	ob interfaces.ChainClient,
	signer interfaces.ChainSigner,
) {
	res, err := co.bridge.GetAllOutboundTrackerByChainbound(chainID, interfaces.Ascending)
	if err != nil {
		co.logger.Observer.Warn().Err(err).Msgf("ScheduleCctxEVM: GetAllOutboundTrackerByChainbound failed for chain %d", chainID)
		return
	}
	trackerMap := make(map[uint64]bool)
	for _, v := range res {
		trackerMap[v.Nonce] = true
	}
	outboundScheduleLookahead := ob.GetChainParams().OutboundScheduleLookahead
	// #nosec G701 always in range
	outboundScheduleLookback := uint64(float64(outboundScheduleLookahead) * evmOutboundTxLookbackFactor)
	// #nosec G701 positive
	outboundScheduleInterval := uint64(ob.GetChainParams().OutboundScheduleInterval)

	for idx, cctx := range cctxList {
		params := cctx.GetCurrentOutboundParam()
		nonce := params.TssNonce
		outboundID := outboundprocessor.ToOutboundID(cctx.Index, params.ReceiverChainId, nonce)

		if params.ReceiverChainId != chainID {
			co.logger.Observer.Error().Msgf("ScheduleCctxEVM: outtx %s chainid mismatch: want %d, got %d", outboundID, chainID, params.ReceiverChainId)
			continue
		}
		if params.TssNonce > cctxList[0].GetCurrentOutboundParam().TssNonce+outboundScheduleLookback {
			co.logger.Observer.Error().Msgf("ScheduleCctxEVM: nonce too high: signing %d, earliest pending %d, chain %d",
				params.TssNonce, cctxList[0].GetCurrentOutboundParam().TssNonce, chainID)
			break
		}

		// try confirming the outtx
		included, _, err := ob.IsOutboundProcessed(cctx, co.logger.Observer)
		if err != nil {
			co.logger.Observer.Error().Err(err).Msgf("ScheduleCctxEVM: IsOutboundProcessed faild for chain %d nonce %d", chainID, nonce)
			continue
		}
		if included {
			co.logger.Observer.Info().Msgf("ScheduleCctxEVM: outtx %s already included; do not schedule keysign", outboundID)
			continue
		}

		// determining critical outtx; if it satisfies following criteria
		// 1. it's the first pending outtx for this chain
		// 2. the following 5 nonces have been in tracker
		criticalInterval := uint64(10)                      // for critical pending outTx we reduce re-try interval
		nonCriticalInterval := outboundScheduleInterval * 2 // for non-critical pending outTx we increase re-try interval
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
		if nonce%outboundScheduleInterval == zetaHeight%outboundScheduleInterval && !outboundManager.IsOutboundActive(outboundID) {
			outboundManager.StartTryProcess(outboundID)
			co.logger.Observer.Debug().Msgf("ScheduleCctxEVM: sign outbound %s with value %d\n", outboundID, cctx.GetCurrentOutboundParam().Amount)
			go signer.TryProcessOutbound(cctx, outboundManager, outboundID, ob, co.bridge, zetaHeight)
		}

		// #nosec G701 always in range
		if int64(idx) >= outboundScheduleLookahead-1 { // only look at 'lookahead' cctxs per chain
			break
		}
	}
}

// ScheduleCctxBTC schedules bitcoin outtx keysign on each ZetaChain block (the ticker)
// 1. schedule at most one keysign per ticker
// 2. schedule keysign only when nonce-mark UTXO is available
// 3. stop keysign when lookahead is reached
func (co *CoreObserver) ScheduleCctxBTC(
	outboundManager *outboundprocessor.Processor,
	zetaHeight uint64,
	chainID int64,
	cctxList []*types.CrossChainTx,
	ob interfaces.ChainClient,
	signer interfaces.ChainSigner,
) {
	btcClient, ok := ob.(*bitcoin.BTCChainClient)
	if !ok { // should never happen
		co.logger.Observer.Error().Msgf("ScheduleCctxBTC: chain client is not a bitcoin client")
		return
	}
	// #nosec G701 positive
	interval := uint64(ob.GetChainParams().OutboundScheduleInterval)
	lookahead := ob.GetChainParams().OutboundScheduleLookahead

	// schedule at most one keysign per ticker
	for idx, cctx := range cctxList {
		params := cctx.GetCurrentOutboundParam()
		nonce := params.TssNonce
		outboundID := outboundprocessor.ToOutboundID(cctx.Index, params.ReceiverChainId, nonce)

		if params.ReceiverChainId != chainID {
			co.logger.Observer.Error().Msgf("ScheduleCctxBTC: outtx %s chainid mismatch: want %d, got %d", outboundID, chainID, params.ReceiverChainId)
			continue
		}
		// try confirming the outtx
		included, confirmed, err := btcClient.IsOutboundProcessed(cctx, co.logger.Observer)
		if err != nil {
			co.logger.Observer.Error().Err(err).Msgf("ScheduleCctxBTC: IsOutboundProcessed faild for chain %d nonce %d", chainID, nonce)
			continue
		}
		if included || confirmed {
			co.logger.Observer.Info().Msgf("ScheduleCctxBTC: outtx %s already included; do not schedule keysign", outboundID)
			continue
		}

		// stop if the nonce being processed is higher than the pending nonce
		if nonce > btcClient.GetPendingNonce() {
			break
		}
		// stop if lookahead is reached
		if int64(idx) >= lookahead { // 2 bitcoin confirmations span is 20 minutes on average. We look ahead up to 100 pending cctx to target TPM of 5.
			co.logger.Observer.Warn().Msgf("ScheduleCctxBTC: lookahead reached, signing %d, earliest pending %d", nonce, cctxList[0].GetCurrentOutboundParam().TssNonce)
			break
		}
		// try confirming the outtx or scheduling a keysign
		if nonce%interval == zetaHeight%interval && !outboundManager.IsOutboundActive(outboundID) {
			outboundManager.StartTryProcess(outboundID)
			co.logger.Observer.Debug().Msgf("ScheduleCctxBTC: sign outtx %s with value %d\n", outboundID, params.Amount)
			go signer.TryProcessOutbound(cctx, outboundManager, outboundID, ob, co.bridge, zetaHeight)
		}
	}
}
