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
	"github.com/zeta-chain/zetacore/zetaclient/outtxprocessor"
	"github.com/zeta-chain/zetacore/zetaclient/ratelimiter"
)

const (
	// EVMOutboundTxLookbackFactor is the factor to determine how many nonces to look back for pending cctxs
	// For example, give OutboundTxScheduleLookahead of 120, pending NonceLow of 1000 and factor of 1.0,
	// the scheduler need to be able to pick up and schedule any pending cctx with nonce < 880 (1000 - 120 * 1.0)
	// NOTE: 1.0 means look back the same number of cctxs as we look ahead
	EVMOutboundTxLookbackFactor = 1.0
)

type ZetaCoreLog struct {
	ChainLogger      zerolog.Logger
	ZetaChainWatcher zerolog.Logger
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
	chainLogger := logger.With().
		Str("chain", "ZetaChain").
		Logger()
	co.logger = ZetaCoreLog{
		ChainLogger:      chainLogger,
		ZetaChainWatcher: chainLogger.With().Str("module", "ZetaChainWatcher").Logger(),
	}

	co.bridge = bridge
	co.signerMap = signerMap

	co.clientMap = clientMap
	co.logger.ChainLogger.Info().Msg("starting core observer")
	balance, err := bridge.GetZetaHotKeyBalance()
	if err != nil {
		co.logger.ChainLogger.Error().Err(err).Msg("error getting last balance of the hot key")
	}
	co.lastOperatorBalance = balance

	return &co
}

func (co *CoreObserver) MonitorCore(appContext *appcontext.AppContext) {
	myid := co.bridge.GetKeys().GetAddress()
	co.logger.ZetaChainWatcher.Info().Msgf("Starting Send Scheduler for %s", myid)
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
				co.logger.ZetaChainWatcher.Info().Msgf(
					"updated zeta connector address for chainID %d, new address: %s", chainID, zetaConnectorAddress)
			}
			if erc20CustodyAddress != signer.GetERC20CustodyAddress() {
				signer.SetERC20CustodyAddress(erc20CustodyAddress)
				co.logger.ZetaChainWatcher.Info().Msgf(
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
			co.logger.ZetaChainWatcher.Info().Msgf(
				"updated chain params for chainID %d, new params: %v", chainID, *evmParams)
		}
	} else if chains.IsBitcoinChain(chainID) {
		_, btcParams, found := coreContext.GetBTCChainParams()

		if found && !observertypes.ChainParamsEqual(curParams, *btcParams) {
			chainOb.SetChainParams(*btcParams)
			co.logger.ZetaChainWatcher.Info().Msgf(
				"updated chain params for Bitcoin, new params: %v", *btcParams)
		}
	}
	return chainOb, nil
}

// GetPendingCctxsWithinRatelimit get pending cctxs across foreign chains within rate limit
func (co *CoreObserver) GetPendingCctxsWithinRatelimit(foreignChains []chains.Chain, logger zerolog.Logger) (map[int64][]*types.CrossChainTx, error) {
	// get rate limiter flags
	rateLimitFlags, err := co.bridge.GetRateLimiterFlags()
	if err != nil {
		return nil, err
	}

	// apply rate limiter or not according to the flags
	applyLimit := ratelimiter.IsRateLimiterUsable(rateLimitFlags)

	// fallback to non-rate-limited query if rate limiter is not usable
	cctxsMap := make(map[int64][]*types.CrossChainTx)
	if !applyLimit {
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
		logger.Info().Msgf("current rate limiter window: %d rate: %s, percentage: %f",
			output.CurrentWithdrawWindow, output.CurrentWithdrawRate.String(), percentageFloat)
	}

	return output.CctxsMap, nil
}

// StartCctxScheduler schedules keysigns for cctxs on each ZetaChain block (the ticker)
func (co *CoreObserver) StartCctxScheduler(appContext *appcontext.AppContext) {
	outTxMan := outtxprocessor.NewOutTxProcessorManager(co.logger.ChainLogger)
	observeTicker := time.NewTicker(3 * time.Second)
	sampledLogger := co.logger.ZetaChainWatcher.Sample(&zerolog.BasicSampler{N: 10})
	var lastBlockNum int64
	for {
		select {
		case <-co.stop:
			co.logger.ZetaChainWatcher.Warn().Msg("StartCctxScheduler: stopped")
			return
		case <-observeTicker.C:
			{
				bn, err := co.bridge.GetZetaBlockHeight()
				if err != nil {
					co.logger.ZetaChainWatcher.Error().Err(err).Msg("StartCctxScheduler: GetZetaBlockHeight fail")
					continue
				}
				if bn < 0 {
					co.logger.ZetaChainWatcher.Error().Msg("StartCctxScheduler: GetZetaBlockHeight returned negative height")
					continue
				}
				if lastBlockNum == 0 {
					lastBlockNum = bn - 1
				}
				if bn > lastBlockNum { // we have a new block
					bn = lastBlockNum + 1
					if bn%10 == 0 {
						co.logger.ZetaChainWatcher.Debug().Msgf("StartCctxScheduler: ZetaCore heart beat: %d", bn)
					}

					balance, err := co.bridge.GetZetaHotKeyBalance()
					if err != nil {
						co.logger.ZetaChainWatcher.Error().Err(err).Msgf("couldn't get operator balance")
					} else {
						diff := co.lastOperatorBalance.Sub(balance)
						if diff.GT(sdkmath.NewInt(0)) && diff.LT(sdkmath.NewInt(math.MaxInt64)) {
							co.ts.AddFeeEntry(bn, diff.Int64())
							co.lastOperatorBalance = balance
						}
					}

					// set current hot key burn rate
					metrics.HotKeyBurnRate.Set(float64(co.ts.HotKeyBurnRate.GetBurnRate().Int64()))

					// get supported foreign chains
					coreContext := appContext.ZetaCoreContext()
					foreignChains := coreContext.GetEnabledForeignChains()

					// query pending cctxs across all foreign chains within rate limit
					cctxMap, err := co.GetPendingCctxsWithinRatelimit(foreignChains, sampledLogger)
					if err != nil {
						co.logger.ZetaChainWatcher.Error().Err(err).Msgf("StartCctxScheduler: GetPendingCctxsWithinRatelimit failed")
					}

					// schedule keysign for pending cctxs on each chain
					for _, c := range foreignChains {
						// get cctxs from map and set pending transactions prometheus gauge
						cctxList := cctxMap[c.ChainId]
						metrics.PendingTxsPerChain.WithLabelValues(c.ChainName.String()).Set(float64(len(cctxList)))
						if len(cctxList) == 0 {
							continue
						}

						// update chain parameters for signer and chain client
						signer, err := co.GetUpdatedSigner(coreContext, c.ChainId)
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("StartCctxScheduler: GetUpdatedSigner failed for chain %d", c.ChainId)
							continue
						}
						ob, err := co.GetUpdatedChainClient(coreContext, c.ChainId)
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("StartCctxScheduler: GetUpdatedChainClient failed for chain %d", c.ChainId)
							continue
						}
						if !corecontext.IsOutboundObservationEnabled(coreContext, ob.GetChainParams()) {
							continue
						}

						// #nosec G701 range is verified
						zetaHeight := uint64(bn)
						if chains.IsEVMChain(c.ChainId) {
							co.ScheduleCctxEVM(outTxMan, zetaHeight, c.ChainId, cctxList, ob, signer)
						} else if chains.IsBitcoinChain(c.ChainId) {
							co.ScheduleCctxBTC(outTxMan, zetaHeight, c.ChainId, cctxList, ob, signer)
						} else {
							co.logger.ZetaChainWatcher.Error().Msgf("StartCctxScheduler: unsupported chain %d", c.ChainId)
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
	outTxMan *outtxprocessor.Processor,
	zetaHeight uint64,
	chainID int64,
	cctxList []*types.CrossChainTx,
	ob interfaces.ChainClient,
	signer interfaces.ChainSigner) {
	res, err := co.bridge.GetAllOutTxTrackerByChain(chainID, interfaces.Ascending)
	if err != nil {
		co.logger.ZetaChainWatcher.Warn().Err(err).Msgf("ScheduleCctxEVM: GetAllOutTxTrackerByChain failed for chain %d", chainID)
		return
	}
	trackerMap := make(map[uint64]bool)
	for _, v := range res {
		trackerMap[v.Nonce] = true
	}
	outboundScheduleLookahead := ob.GetChainParams().OutboundTxScheduleLookahead
	// #nosec G701 always in range
	outboundScheduleLookback := uint64(float64(outboundScheduleLookahead) * EVMOutboundTxLookbackFactor)
	// #nosec G701 positive
	outboundScheduleInterval := uint64(ob.GetChainParams().OutboundTxScheduleInterval)

	for idx, cctx := range cctxList {
		params := cctx.GetCurrentOutTxParam()
		nonce := params.OutboundTxTssNonce
		outTxID := outtxprocessor.ToOutTxID(cctx.Index, params.ReceiverChainId, nonce)

		if params.ReceiverChainId != chainID {
			co.logger.ZetaChainWatcher.Error().Msgf("ScheduleCctxEVM: outtx %s chainid mismatch: want %d, got %d", outTxID, chainID, params.ReceiverChainId)
			continue
		}
		if params.OutboundTxTssNonce > cctxList[0].GetCurrentOutTxParam().OutboundTxTssNonce+outboundScheduleLookback {
			co.logger.ZetaChainWatcher.Error().Msgf("ScheduleCctxEVM: nonce too high: signing %d, earliest pending %d, chain %d",
				params.OutboundTxTssNonce, cctxList[0].GetCurrentOutTxParam().OutboundTxTssNonce, chainID)
			break
		}

		// try confirming the outtx
		included, _, err := ob.IsOutboundProcessed(cctx, co.logger.ZetaChainWatcher)
		if err != nil {
			co.logger.ZetaChainWatcher.Error().Err(err).Msgf("ScheduleCctxEVM: IsOutboundProcessed faild for chain %d nonce %d", chainID, nonce)
			continue
		}
		if included {
			co.logger.ZetaChainWatcher.Info().Msgf("ScheduleCctxEVM: outtx %s already included; do not schedule keysign", outTxID)
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
		if nonce%outboundScheduleInterval == zetaHeight%outboundScheduleInterval && !outTxMan.IsOutTxActive(outTxID) {
			outTxMan.StartTryProcess(outTxID)
			co.logger.ZetaChainWatcher.Debug().Msgf("ScheduleCctxEVM: sign outtx %s with value %d\n", outTxID, cctx.GetCurrentOutTxParam().Amount)
			go signer.TryProcessOutTx(cctx, outTxMan, outTxID, ob, co.bridge, zetaHeight)
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
	outTxMan *outtxprocessor.Processor,
	zetaHeight uint64,
	chainID int64,
	cctxList []*types.CrossChainTx,
	ob interfaces.ChainClient,
	signer interfaces.ChainSigner) {
	btcClient, ok := ob.(*bitcoin.BTCChainClient)
	if !ok { // should never happen
		co.logger.ZetaChainWatcher.Error().Msgf("scheduleCctxBTC: chain client is not a bitcoin client")
		return
	}
	// #nosec G701 positive
	interval := uint64(ob.GetChainParams().OutboundTxScheduleInterval)
	lookahead := ob.GetChainParams().OutboundTxScheduleLookahead

	// schedule at most one keysign per ticker
	for idx, cctx := range cctxList {
		params := cctx.GetCurrentOutTxParam()
		nonce := params.OutboundTxTssNonce
		outTxID := outtxprocessor.ToOutTxID(cctx.Index, params.ReceiverChainId, nonce)

		if params.ReceiverChainId != chainID {
			co.logger.ZetaChainWatcher.Error().Msgf("scheduleCctxBTC: outtx %s chainid mismatch: want %d, got %d", outTxID, chainID, params.ReceiverChainId)
			continue
		}
		// try confirming the outtx
		included, confirmed, err := btcClient.IsOutboundProcessed(cctx, co.logger.ZetaChainWatcher)
		if err != nil {
			co.logger.ZetaChainWatcher.Error().Err(err).Msgf("scheduleCctxBTC: IsOutboundProcessed faild for chain %d nonce %d", chainID, nonce)
			continue
		}
		if included || confirmed {
			co.logger.ZetaChainWatcher.Info().Msgf("scheduleCctxBTC: outtx %s already included; do not schedule keysign", outTxID)
			continue
		}

		// stop if the nonce being processed is higher than the pending nonce
		if nonce > btcClient.GetPendingNonce() {
			break
		}
		// stop if lookahead is reached
		if int64(idx) >= lookahead { // 2 bitcoin confirmations span is 20 minutes on average. We look ahead up to 100 pending cctx to target TPM of 5.
			co.logger.ZetaChainWatcher.Warn().Msgf("scheduleCctxBTC: lookahead reached, signing %d, earliest pending %d", nonce, cctxList[0].GetCurrentOutTxParam().OutboundTxTssNonce)
			break
		}
		// try confirming the outtx or scheduling a keysign
		if nonce%interval == zetaHeight%interval && !outTxMan.IsOutTxActive(outTxID) {
			outTxMan.StartTryProcess(outTxID)
			co.logger.ZetaChainWatcher.Debug().Msgf("scheduleCctxBTC: sign outtx %s with value %d\n", outTxID, params.Amount)
			go signer.TryProcessOutTx(cctx, outTxMan, outTxID, ob, co.bridge, zetaHeight)
		}
	}
}
