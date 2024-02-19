package zetaclient

import (
	"fmt"
	"math"
	"time"

	appcontext "github.com/zeta-chain/zetacore/zetaclient/app_context"
	"github.com/zeta-chain/zetacore/zetaclient/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/outtxprocessor"

	observertypes "github.com/zeta-chain/zetacore/x/observer/types"

	sdkmath "cosmossdk.io/math"

	"github.com/pkg/errors"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

const (
	MaxLookaheadNonce   = 120
	OutboundTxSignCount = "zetaclient_Outbound_tx_sign_count"
	HotKeyBurnRate      = "zetaclient_hotkey_burn_rate"
)

type ZetaCoreLog struct {
	ChainLogger      zerolog.Logger
	ZetaChainWatcher zerolog.Logger
}

// CoreObserver wraps the zetabridge bridge and adds the client and signer maps to it . This is the high level object used for CCTX interactions
type CoreObserver struct {
	bridge              interfaces.ZetaCoreBridger
	signerMap           map[common.Chain]interfaces.ChainSigner
	clientMap           map[common.Chain]interfaces.ChainClient
	metrics             *metrics.Metrics
	logger              ZetaCoreLog
	ts                  *metrics.TelemetryServer
	stop                chan struct{}
	lastOperatorBalance sdkmath.Int
}

// NewCoreObserver creates a new CoreObserver
func NewCoreObserver(
	appContext *appcontext.AppContext,
	bridge interfaces.ZetaCoreBridger,
	signerMap map[common.Chain]interfaces.ChainSigner,
	clientMap map[common.Chain]interfaces.ChainClient,
	metrics *metrics.Metrics,
	ts *metrics.TelemetryServer,
) *CoreObserver {
	co := CoreObserver{
		ts:   ts,
		stop: make(chan struct{}),
	}
	chainLogger := appContext.Logger().With().
		Str("chain", "ZetaChain").
		Logger()
	co.logger = ZetaCoreLog{
		ChainLogger:      chainLogger,
		ZetaChainWatcher: chainLogger.With().Str("module", "ZetaChainWatcher").Logger(),
	}

	co.bridge = bridge
	co.signerMap = signerMap

	co.clientMap = clientMap
	co.metrics = metrics
	co.logger.ChainLogger.Info().Msg("starting core observer")
	err := metrics.RegisterCounter(OutboundTxSignCount, "number of Outbound tx signed")
	if err != nil {
		co.logger.ChainLogger.Error().Err(err).Msg("error registering counter")
	}
	err = metrics.RegisterGauge(HotKeyBurnRate, "Fee burn rate of the hotkey")
	if err != nil {
		co.logger.ChainLogger.Error().Err(err).Msg("error registering gauge")
	}
	balance, err := bridge.GetZetaHotKeyBalance()
	if err != nil {
		co.logger.ChainLogger.Error().Err(err).Msg("error getting last balance of the hot key")
	}
	co.lastOperatorBalance = balance

	return &co
}

func (co *CoreObserver) GetPromCounter(name string) (prom.Counter, error) {
	cnt, found := metrics.Counters[name]
	if !found {
		return nil, errors.New("counter not found")
	}
	return cnt, nil
}

func (co *CoreObserver) GetPromGauge(name string) (prom.Gauge, error) {
	gauge, found := metrics.Gauges[name]
	if !found {
		return nil, errors.New("gauge not found")
	}
	return gauge, nil
}

func (co *CoreObserver) MonitorCore(appContext *appcontext.AppContext) {
	myid := co.bridge.GetKeys().GetAddress()
	co.logger.ZetaChainWatcher.Info().Msgf("Starting Send Scheduler for %s", myid)
	go co.startCctxScheduler(appContext)

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

// startCctxScheduler schedules keysigns for cctxs on each ZetaChain block (the ticker)
func (co *CoreObserver) startCctxScheduler(appContext *appcontext.AppContext) {
	outTxMan := outtxprocessor.NewOutTxProcessorManager(co.logger.ChainLogger)
	observeTicker := time.NewTicker(3 * time.Second)
	var lastBlockNum int64
	for {
		select {
		case <-co.stop:
			co.logger.ZetaChainWatcher.Warn().Msg("startCctxScheduler: stopped")
			return
		case <-observeTicker.C:
			{
				bn, err := co.bridge.GetZetaBlockHeight()
				if err != nil {
					co.logger.ZetaChainWatcher.Error().Err(err).Msg("startCctxScheduler: GetZetaBlockHeight fail")
					continue
				}
				if bn < 0 {
					co.logger.ZetaChainWatcher.Error().Msg("startCctxScheduler: GetZetaBlockHeight returned negative height")
					continue
				}
				if lastBlockNum == 0 {
					lastBlockNum = bn - 1
				}
				if bn > lastBlockNum { // we have a new block
					bn = lastBlockNum + 1
					if bn%10 == 0 {
						co.logger.ZetaChainWatcher.Debug().Msgf("startCctxScheduler: ZetaCore heart beat: %d", bn)
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

					// Set Current Hot key burn rate
					gauge, err := co.GetPromGauge(HotKeyBurnRate)
					if err != nil {
						co.logger.ZetaChainWatcher.Error().Err(err).Msgf("scheduleCctxEVM: failed to get prometheus gauge: %s for observer", metrics.PendingTxs)
						continue
					} // Gauge only takes float values
					gauge.Set(float64(co.ts.HotKeyBurnRate.GetBurnRate().Int64()))

					// schedule keysign for pending cctxs on each chain
					supportedChains := appContext.ZetaCoreContext().GetEnabledChains()
					for _, c := range supportedChains {
						if c.ChainId == co.bridge.ZetaChain().ChainId {
							continue
						}
						signer := co.signerMap[c]

						cctxList, totalPending, err := co.bridge.ListPendingCctx(c.ChainId)
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("startCctxScheduler: ListPendingCctx failed for chain %d", c.ChainId)
							continue
						}
						ob, err := co.getUpdatedChainOb(appContext, c.ChainId)
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("startCctxScheduler: getTargetChainOb failed for chain %d", c.ChainId)
							continue
						}
						// Set Pending transactions prometheus gauge
						gauge, err := ob.GetPromGauge(metrics.PendingTxs)
						if err != nil {
							co.logger.ZetaChainWatcher.Error().Err(err).Msgf("scheduleCctxEVM: failed to get prometheus gauge: %s for chain %d", metrics.PendingTxs, c.ChainId)
							continue
						}
						gauge.Set(float64(totalPending))

						// #nosec G701 range is verified
						zetaHeight := uint64(bn)
						if common.IsEVMChain(c.ChainId) {
							co.scheduleCctxEVM(outTxMan, zetaHeight, c.ChainId, cctxList, ob, signer)
						} else if common.IsBitcoinChain(c.ChainId) {
							co.scheduleCctxBTC(outTxMan, zetaHeight, c.ChainId, cctxList, ob, signer)
						} else {
							co.logger.ZetaChainWatcher.Error().Msgf("startCctxScheduler: unsupported chain %d", c.ChainId)
							continue
						}
					}
					// update last processed block number
					lastBlockNum = bn
					co.ts.SetCoreBlockNumber(lastBlockNum)
				}
			}
		}
	}
}

// scheduleCctxEVM schedules evm outtx keysign on each ZetaChain block (the ticker)
func (co *CoreObserver) scheduleCctxEVM(
	outTxMan *outtxprocessor.Processor,
	zetaHeight uint64,
	chainID int64,
	cctxList []*types.CrossChainTx,
	ob interfaces.ChainClient,
	signer interfaces.ChainSigner) {
	res, err := co.bridge.GetAllOutTxTrackerByChain(chainID, interfaces.Ascending)
	if err != nil {
		co.logger.ZetaChainWatcher.Warn().Err(err).Msgf("scheduleCctxEVM: GetAllOutTxTrackerByChain failed for chain %d", chainID)
		return
	}
	trackerMap := make(map[uint64]bool)
	for _, v := range res {
		trackerMap[v.Nonce] = true
	}

	for idx, cctx := range cctxList {
		params := cctx.GetCurrentOutTxParam()
		nonce := params.OutboundTxTssNonce
		outTxID := outtxprocessor.ToOutTxID(cctx.Index, params.ReceiverChainId, nonce)

		if params.ReceiverChainId != chainID {
			co.logger.ZetaChainWatcher.Error().Msgf("scheduleCctxEVM: outtx %s chainid mismatch: want %d, got %d", outTxID, chainID, params.ReceiverChainId)
			continue
		}
		if params.OutboundTxTssNonce > cctxList[0].GetCurrentOutTxParam().OutboundTxTssNonce+MaxLookaheadNonce {
			co.logger.ZetaChainWatcher.Error().Msgf("scheduleCctxEVM: nonce too high: signing %d, earliest pending %d, chain %d",
				params.OutboundTxTssNonce, cctxList[0].GetCurrentOutTxParam().OutboundTxTssNonce, chainID)
			break
		}

		// try confirming the outtx
		included, _, err := ob.IsSendOutTxProcessed(cctx.Index, params.OutboundTxTssNonce, params.CoinType, co.logger.ZetaChainWatcher)
		if err != nil {
			co.logger.ZetaChainWatcher.Error().Err(err).Msgf("scheduleCctxEVM: IsSendOutTxProcessed faild for chain %d nonce %d", chainID, nonce)
			continue
		}
		if included {
			co.logger.ZetaChainWatcher.Info().Msgf("scheduleCctxEVM: outtx %s already included; do not schedule keysign", outTxID)
			continue
		}

		// #nosec G701 positive
		interval := uint64(ob.GetChainParams().OutboundTxScheduleInterval)
		lookahead := ob.GetChainParams().OutboundTxScheduleLookahead

		// determining critical outtx; if it satisfies following criteria
		// 1. it's the first pending outtx for this chain
		// 2. the following 5 nonces have been in tracker
		criticalInterval := uint64(10)      // for critical pending outTx we reduce re-try interval
		nonCriticalInterval := interval * 2 // for non-critical pending outTx we increase re-try interval
		if nonce%criticalInterval == zetaHeight%criticalInterval {
			count := 0
			for i := nonce + 1; i <= nonce+10; i++ {
				if _, found := trackerMap[i]; found {
					count++
				}
			}
			if count >= 5 {
				interval = criticalInterval
			}
		}
		// if it's already in tracker, we increase re-try interval
		if _, ok := trackerMap[nonce]; ok {
			interval = nonCriticalInterval
		}

		// otherwise, the normal interval is used
		if nonce%interval == zetaHeight%interval && !outTxMan.IsOutTxActive(outTxID) {
			outTxMan.StartTryProcess(outTxID)
			co.logger.ZetaChainWatcher.Debug().Msgf("scheduleCctxEVM: sign outtx %s with value %d\n", outTxID, cctx.GetCurrentOutTxParam().Amount)
			go signer.TryProcessOutTx(cctx, outTxMan, outTxID, ob, co.bridge, zetaHeight)
		}

		// #nosec G701 always in range
		if int64(idx) >= lookahead-1 { // only look at 'lookahead' cctxs per chain
			break
		}
	}
}

// scheduleCctxBTC schedules bitcoin outtx keysign on each ZetaChain block (the ticker)
// 1. schedule at most one keysign per ticker
// 2. schedule keysign only when nonce-mark UTXO is available
// 3. stop keysign when lookahead is reached
func (co *CoreObserver) scheduleCctxBTC(
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
		included, confirmed, err := btcClient.IsSendOutTxProcessed(cctx.Index, nonce, params.CoinType, co.logger.ZetaChainWatcher)
		if err != nil {
			co.logger.ZetaChainWatcher.Error().Err(err).Msgf("scheduleCctxBTC: IsSendOutTxProcessed faild for chain %d nonce %d", chainID, nonce)
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

func (co *CoreObserver) getUpdatedChainOb(appContext *appcontext.AppContext, chainID int64) (interfaces.ChainClient, error) {
	chainOb, err := co.getTargetChainOb(chainID)
	if err != nil {
		return nil, err
	}
	// update chain client core parameters
	curParams := chainOb.GetChainParams()
	if common.IsEVMChain(chainID) {
		evmParams, found := appContext.ZetaCoreContext().GetEVMChainParams(chainID)
		if found && !observertypes.ChainParamsEqual(curParams, *evmParams) {
			chainOb.SetChainParams(*evmParams)
			co.logger.ZetaChainWatcher.Info().Msgf(
				"updated chain params for chainID %d, new params: %v",
				chainID,
				*evmParams,
			)
		}
	} else if common.IsBitcoinChain(chainID) {
		_, btcParams, found := appContext.ZetaCoreContext().GetBTCChainParams()

		if found && !observertypes.ChainParamsEqual(curParams, *btcParams) {
			chainOb.SetChainParams(*btcParams)
			co.logger.ZetaChainWatcher.Info().Msgf(
				"updated chain params for Bitcoin, new params: %v",
				*btcParams,
			)
		}
	}
	return chainOb, nil
}

func (co *CoreObserver) getTargetChainOb(chainID int64) (interfaces.ChainClient, error) {
	c := common.GetChainFromChainID(chainID)
	if c == nil {
		return nil, fmt.Errorf("chain not found for chainID %d", chainID)
	}
	chainOb, found := co.clientMap[*c]
	if !found {
		return nil, fmt.Errorf("chain client not found for chainID %d", chainID)
	}
	return chainOb, nil
}
