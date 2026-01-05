package evm

import (
	"context"
	"fmt"
	"time"

	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
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

// scheduleCCTX schedules outbound transactions on each ZetaChain block.
func (e *EVM) scheduleCCTX(ctx context.Context) error {
	// skip stale block event if any
	if _, stale, err := e.signer.IsStaleBlockEvent(ctx, e.observer.ZetaRepo()); err != nil {
		return errors.Wrap(err, "unable to check stale block event")
	} else if stale {
		return nil
	}

	if err := e.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "unable to update chain params")
	}

	var (
		chainID   = e.observer.Chain().ChainId
		lookahead = e.observer.ChainParams().OutboundScheduleLookahead
		// #nosec G115 always in range
		maxNonceOffset = uint64(float64(lookahead) * constant.MaxNonceOffsetFactor)
	)

	cctxList, err := e.observer.ZetaRepo().GetPendingCCTXs(ctx)
	if err != nil {
		return err
	}

	nextTSSNonce, err := e.signer.NextTSSNonce(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get next tss nonce")
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
		case params.TssNonce > cctxList[0].GetCurrentOutboundParam().TssNonce+maxNonceOffset:
			return fmt.Errorf(
				"nonce %d is too high (%s). Earliest nonce %d",
				params.TssNonce,
				outboundID,
				cctxList[0].GetCurrentOutboundParam().TssNonce,
			)
		}

		// if cctx's nonce is lower than next nonce, it means the outbound was already processed;
		// in this case, we just need to check its confirmations and post vote if it's confirmed.
		if params.TssNonce < nextTSSNonce {
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
				0,
			)
		}
	}

	return nil
}

// scheduleKeysign schedules keysign for outbound transactions
func (e *EVM) scheduleKeysign(ctx context.Context) error {
	s := e.signer
	zetaRepo := e.observer.ZetaRepo()

	// skip stale block event if any
	zetaHeight, stale, err := s.IsStaleBlockEvent(ctx, zetaRepo)
	if err != nil {
		return errors.Wrap(err, "unable to check stale block event")
	} else if stale {
		return nil
	}

	// next tss nonce is the starting point to start keysign
	nextTSSNonce, err := s.NextTSSNonce(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get next tss nonce")
	}

	// query pending nonces and see if there is any tx to sign
	p, err := zetaRepo.GetPendingNonces(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to get pending nonces for chain %d", s.Chain().ChainId)
	}

	// remove stale keysign info to release memory
	// #nosec G115 - always positive
	s.RemoveKeysignInfo(uint64(p.NonceLow))

	// check if it's the time to perform keysign
	interval := e.observer.ChainParams().OutboundScheduleInterval
	shouldSign := e.signer.IsTimeToKeysign(*p, nextTSSNonce, zetaHeight, interval)
	if !shouldSign {
		return nil
	}

	var (
		batchNumber = base.NonceToBatchNumber(nextTSSNonce)
		batch       = s.GetKeysignBatch(ctx, zetaRepo, batchNumber)
	)

	// The following are basic prerequisites for TSS batch keysign to work in a sequential way:
	// 1. the starting batch is determined by 'nextTSSNonce', which is deterministic across all TSS signers.
	// 2. the TSS signers ALWAYS fully sign batch N (e.g. nonce 0~9) before signing batch N+1 (e.g. nonce 10~19).
	// 3. for any signed tx, there should be always >= 2/3 of signers have cached signature for it.
	//    so far as any one of the signers has signature for a tx, the tx can be processed, no worries;
	//    the TSS signers will ALWAYS be in sync again on what to sign next once signed txs get processed.
	// 4. any signed batch is skipped when looping forward, as we know they will be processed soon; the TSS
	//    signers can continue signing the batches ahead without waiting.
	//
	// Why signing in batches?
	// - Batching multiple digests in one request is efficient.
	// - Reduce parallel keysign requests and avoid spamming the TSS service.
	//
	// Why signing sequentially?
	// - EVM chain processes txs by nonce sequentially; if tx with nonce N is not processed, it blocks tx with nonce N+1.
	// - For each EVM chain at any time, the TSS signers only need to schedule one single batch keysign request, and this
	//   request will contain the digests of the txs that we immediately need to send out.
	//
	// Why signing adjacent batches in loop without waiting for another interval?
	// - the keysign interval in chain params is used as a timing signal to trigger first keysign handshake.
	// - at the moment when the first keysign is completed, it means the TSS signers are strictly in sync on:
	//   1. the timestamp, regardless of timezones, system time, or Zeta height.
	//   2. which batch number to sign next.
	for batch != nil {
		if !s.IsBatchSigned(batch) {
			if err := s.SignBatch(ctx, *batch, zetaHeight); err != nil {
				return errors.Wrapf(err, "break keysign batch: %d", batch.BatchNumber())
			}
		}

		// now that signed, move to next batch only if this batch is the end
		// we don't want to leave a batch partially signed and jump to next
		if !batch.IsEnd() {
			break
		}

		batchNumber++
		batch = s.GetKeysignBatch(ctx, zetaRepo, batchNumber)
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
