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
	zetaBlock, _, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get block event from context")
	}

	// get real-time zeta height
	zetaHeight, err := e.observer.ZetaRepo().GetBlockHeight(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get zeta height")
	}

	logger := e.signer.Logger().Std.With().Int64("zeta_height", zetaHeight).Logger()

	// real-time zeta height are the signals to trigger TSS keysign on the exact same time,
	// so we need to ensure the block event is up to date (not a stale one accumulated in the channel)
	if zetaBlock.Block.Height != zetaHeight {
		logger.Info().Int64("event_height", zetaBlock.Block.Height).Msg("skip stale block event")
		return nil
	}

	// keysign happens only when zeta height is a multiple of the schedule interval
	scheduleInterval := e.observer.ChainParams().OutboundScheduleInterval
	if zetaHeight%scheduleInterval != 0 {
		return nil
	}

	// query pending nonces and see if there is any tx to sign
	nonceLow := uint64(0)
	p, err := e.observer.ZetaRepo().GetPendingNonces(ctx)
	switch {
	case err != nil:
		return errors.Wrapf(err, "unable to get pending nonces for chain %d", e.Chain().ChainId)
	case p.NonceLow < 0: // never happens
		return fmt.Errorf("negative pending nonce %d for chain %d", p.NonceLow, e.Chain().ChainId)
	default:
		// remove stale keysign info
		// #nosec G115 - checked positive
		nonceLow = uint64(p.NonceLow)
		e.signer.RemoveKeysignInfo(nonceLow)

		// return if nothing to sign
		if p.NonceLow >= p.NonceHigh {
			return nil
		}
	}

	// use the first pending nonce's batch number as the starting batch for keysign, please note that:
	// 1. the starting batch is deterministic across all TSS signers
	// 2. for any signed batch, there should always be >= 2/3 of signers have cached signatures for it
	// 3. so far as any one of the signers has cached signatures for a batch, txs can be processed, no worries
	batchNumber := base.NonceToBatchNumber(nonceLow)
	batch := e.signer.GetKeysignBatch(batchNumber)
	if batch == nil {
		e.signer.Logger().Std.Info().Msg("waiting for pending cctxs to finalize")
	}

	// sign from the first to the last batch sequentially. In sequential mode,
	// next keysign happens immediately after completing the previous one, so
	// there is no need to wait for another interval to align the signers.
	for batch != nil {
		if err := e.signer.SignBatch(ctx, *batch, zetaHeight); err != nil {
			return errors.Wrapf(err, "break keysign loop: %d", batch.BatchNumber())
		}

		batchNumber++
		batch = e.signer.GetKeysignBatch(batchNumber)
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
		var (
			params     = cctx.GetCurrentOutboundParam()
			outboundID = base.OutboundIDFromCCTX(cctx)
		)

		switch {
		case params.ReceiverChainId != chainID:
			e.outboundLogger(outboundID).Error().Msg("chain id mismatch")
			continue
		case params.TssNonce > cctxList[0].GetCurrentOutboundParam().TssNonce+outboundScheduleLookBack:
			return fmt.Errorf(
				"nonce %d is too high (%s). Earliest nonce %d",
				params.TssNonce,
				outboundID,
				cctxList[0].GetCurrentOutboundParam().TssNonce,
			)
		}

		// if next nonce is higher than cctx's nonce, it means the outbound is already processed
		// vote outbound if it's already confirmed
		if params.TssNonce < nextNonce {
			// vote outbound if it's already confirmed
			continueKeysign, err := e.observer.VoteOutboundIfConfirmed(ctx, cctx)
			switch {
			case err != nil:
				e.outboundLogger(outboundID).Error().
					Err(err).
					Msg("schedule CCTX: call to VoteOutboundIfConfirmed failed")
				continue
			case !continueKeysign:
				e.outboundLogger(outboundID).Debug().Msg("schedule CCTX: outbound already processed")
				continue
			case e.signer.IsOutboundActive(outboundID):
				// outbound is already being processed
				continue
			}
		}

		// process this CCTX
		go e.signer.TryProcessOutbound(
			ctx,
			cctx,
			e.observer.ZetaRepo(),
			zetaHeight,
		)

		// only look at 'lookahead' cctxs per chain
		if int64(idx) >= lookahead-1 {
			break
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
