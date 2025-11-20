package evm

import (
	"context"
	"fmt"
	"time"

	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/node/zetaclient/chains/evm/signer"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
)

type EVM struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
}

const (
	// outboundLookBackFactor is the factor to determine how many nonces to look back for pending cctxs
	// For example, give OutboundScheduleLookahead of 120, pending NonceLow of 1000 and factor of 1.0,
	// the scheduler need to be able to pick up and schedule any pending cctx with nonce < 880 (1000 - 120 * 1.0)
	// NOTE: 1.0 means look back the same number of cctxs as we look ahead
	outboundLookBackFactor = 1.0
)

func New(scheduler *scheduler.Scheduler, observer *observer.Observer, signer *signer.Signer) *EVM {
	return &EVM{
		scheduler: scheduler,
		observer:  observer,
		signer:    signer,
	}
}

func (e *EVM) Chain() chains.Chain {
	return e.observer.Chain()
}

func (e *EVM) Start(ctx context.Context) error {
	if ok := e.observer.Observer.Start(); !ok {
		return errors.New("observer is already started")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get app from context")
	}

	cctxBlockChan, err := e.observer.ZetaRepo().WatchNewBlocks(ctx)
	if err != nil {
		return err
	}

	keysignBlockChan, err := e.observer.ZetaRepo().WatchNewBlocks(ctx)
	if err != nil {
		return err
	}

	optInboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(e.observer.ChainParams().InboundTicker)
	})

	optGasInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(e.observer.ChainParams().GasPriceTicker)
	})

	optOutboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(e.observer.ChainParams().OutboundTicker)
	})

	optInboundSkipper := scheduler.Skipper(func() bool { return base.CheckSkipInbound(e.observer.Observer, app) })
	optOutboundSkipper := scheduler.Skipper(func() bool { return base.CheckSkipOutbound(e.observer.Observer, app) })
	optGasPriceSkipper := scheduler.Skipper(func() bool { return base.CheckSkipGasPrice(e.observer.Observer, app) })

	register := func(exec scheduler.Executable, name string, opts ...scheduler.Opt) {
		opts = append([]scheduler.Opt{
			scheduler.GroupName(e.group()),
			scheduler.Name(name),
		}, opts...)

		e.scheduler.Register(ctx, exec, opts...)
	}

	// Observers
	register(e.observer.ObserveInbound, "observe_inbound", optInboundInterval, optInboundSkipper)
	register(e.observer.ProcessInboundTrackers, "process_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(e.observer.ProcessInternalTrackers, "process_internal_trackers", optInboundInterval, optInboundSkipper)
	register(e.observer.PostGasPrice, "post_gas_price", optGasInterval, optGasPriceSkipper)
	register(e.observer.CheckRPCStatus, "check_rpc_status")
	register(e.observer.ProcessOutboundTrackers, "process_outbound_trackers", optOutboundInterval, optOutboundSkipper)

	// CCTX Scheduler
	register(e.scheduleCCTX, "schedule_cctx", scheduler.BlockTicker(cctxBlockChan), optOutboundSkipper)

	// TSS keysign scheduler
	register(e.scheduleKeysign, "schedule_keysign", scheduler.BlockTicker(keysignBlockChan), optOutboundSkipper)

	return nil
}

func (e *EVM) Stop() {
	e.observer.Logger().Chain.Info().Msg("stopping observer-signer")
	e.scheduler.StopGroup(e.group())
}

func (e *EVM) group() scheduler.Group {
	return scheduler.Group(
		fmt.Sprintf("evm:%d", e.observer.Chain().ChainId),
	)
}

// scheduleKeysign schedules keysign for outbound transactions
func (e *EVM) scheduleKeysign(ctx context.Context) error {
	s := e.signer
	zetaRepo := e.observer.ZetaRepo()

	zetaHeight, stale, err := s.IsStaleBlockEvent(ctx, zetaRepo)
	if err != nil {
		return errors.Wrap(err, "unable to check block event")
	} else if stale {
		return nil
	}

	// check if it's the time to perform keysign
	interval := e.observer.ChainParams().OutboundScheduleInterval
	shouldSign, startNonce, err := e.signer.IsTimeToSign(ctx, zetaRepo, zetaHeight, interval)
	if err != nil {
		return errors.Wrap(err, "unable to determine the timing of TSS keysign")
	} else if !shouldSign {
		return nil
	}

	// use the first pending nonce's batch number as the starting point for keysign, please note that:
	// 1. yhe starting point is deterministic across all TSS signers
	// 2. gor any signed batch, there should always be >= 2/3 of signers have cached signatures for it;
	//    so far as any one of the signers has cached signatures for a batch, txs can be processed, no worries.
	// 3. sny signed batch is skipped during the loop, so the keysign effectively starts from the first batch
	//    that has not been signed by >= 2/3 of signers.
	var (
		batchNumber = base.NonceToBatchNumber(startNonce)
		batch       = s.GetKeysignBatch(ctx, zetaRepo, batchNumber)
	)

	// dign from the first to the last batch sequentially. In sequential mode, the first keysign is a handshake
	// for the signers to get in sync, and the next keysign is scheduled immediately after completing the previous
	// one, so there is no need to wait for another interval of time to do handshake again.
	for batch != nil {
		if !s.IsBatchSigned(batch) {
			if err := s.SignBatch(ctx, *batch, zetaHeight); err != nil {
				return errors.Wrapf(err, "break keysign batch: %d", batch.BatchNumber())
			}
		}

		// now that signed, move to next batch only if this batch is ending
		// we don't want to leave a batch partially signed and jump to next
		if !batch.IsEnding() {
			break
		}

		batchNumber++
		batch = s.GetKeysignBatch(ctx, zetaRepo, batchNumber)
	}

	return nil
}

// scheduleCCTX schedules outbound transactions on each zeta block.
func (e *EVM) scheduleCCTX(ctx context.Context) error {
	if err := e.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "unable to update chain params")
	}

	zetaBlock, delay, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get zeta block from context")
	}

	time.Sleep(delay)

	var (
		chainID = e.observer.Chain().ChainId

		// #nosec G115 always in range
		zetaHeight = uint64(zetaBlock.Block.Height)
		lookahead  = e.observer.ChainParams().OutboundScheduleLookahead

		// #nosec G115 always in range
		outboundScheduleLookBack = uint64(float64(lookahead) * outboundLookBackFactor)
	)

	cctxList, err := e.observer.ZetaRepo().GetPendingCCTXs(ctx)
	if err != nil {
		return err
	}

	nextNonce, err := e.signer.NextNonce(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get next nonce")
	}

	for idx, cctx := range cctxList {
		// only look at 'lookahead' cctxs per chain
		if int64(idx) >= lookahead {
			break
		}

		var (
			params     = cctx.GetCurrentOutboundParam()
			outboundID = base.OutboundIDFromCCTX(cctx)
		)

		switch {
		case params.ReceiverChainId != chainID:
			return fmt.Errorf("chain id mismatch: want %d, got %d", chainID, params.ReceiverChainId)
		case params.TssNonce > cctxList[0].GetCurrentOutboundParam().TssNonce+outboundScheduleLookBack:
			return fmt.Errorf(
				"nonce %d is too high (%s). Earliest nonce %d",
				params.TssNonce,
				outboundID,
				cctxList[0].GetCurrentOutboundParam().TssNonce,
			)
		}

		// if cctx's nonce is lower than next nonce, it means the outbound was already processed;
		// in this case, we just need to check its confirmations and post vote if it's confirmed.
		if params.TssNonce < nextNonce {
			if _, err := e.observer.VoteOutboundIfConfirmed(ctx, cctx); err != nil {
				e.outboundLogger(outboundID).Error().Err(err).Msg("call to VoteOutboundIfConfirmed failed")
			}
			continue
		}

		// process this CCTX
		if !e.signer.IsOutboundActive(outboundID) {
			go e.signer.TryProcessOutbound(
				ctx,
				cctx,
				e.observer.ZetaRepo(),
				zetaHeight,
			)
		}
	}

	return nil
}

func (e *EVM) updateChainParams(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	chain, err := app.GetChain(e.observer.Chain().ChainId)
	if err != nil {
		return errors.Wrap(err, "unable to get chain")
	}

	params := chain.Params()

	// Update chain params
	e.observer.SetChainParams(*params)

	// Update zeta connector, ERC20 custody, and gateway addresses
	e.signer.SetZetaConnectorAddress(eth.HexToAddress(params.ConnectorContractAddress))
	e.signer.SetERC20CustodyAddress(eth.HexToAddress(params.Erc20CustodyContractAddress))
	e.signer.SetGatewayAddress(params.GatewayAddress)

	return nil
}

func (e *EVM) outboundLogger(id string) *zerolog.Logger {
	l := e.observer.Logger().Outbound.With().Str(logs.FieldOutboundID, id).Logger()

	return &l
}
