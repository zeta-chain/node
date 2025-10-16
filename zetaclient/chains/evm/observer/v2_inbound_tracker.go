package observer

import (
	"context"
	"fmt"
	"strconv"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/evm/client"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrGatewayNotSet = errors.New("gateway contract not set")
)

// ProcessInboundTrackerV2 processes inbound tracker events from the gateway
// TODO: add test coverage
// https://github.com/zeta-chain/node/issues/2669
func (ob *Observer) ProcessInboundTrackerV2(
	ctx context.Context,
	tx *client.Transaction,
	receipt *ethtypes.Receipt,
	isInternalTracker bool,
) error {
	gatewayAddr, gateway, err := ob.getGatewayContract()
	if err != nil {
		ob.Logger().Inbound.Debug().
			Err(err).
			Msg("error getting gateway contract for processing inbound tracker")
		return ErrGatewayNotSet
	}

	// check confirmations
	if !ob.IsBlockConfirmedForInboundSafe(receipt.BlockNumber.Uint64()) {
		return fmt.Errorf(
			"inbound %s has not been confirmed yet: receipt block %d",
			tx.Hash,
			receipt.BlockNumber.Uint64(),
		)
	}

	// Check if multiple calls are enabled
	allowMultipleCalls := zctx.EnableMultipleCallsFeatureFlag(ctx)
	eventFound := false

	for _, log := range receipt.Logs {
		if log == nil || log.Address != gatewayAddr {
			continue
		}

		// try parsing deposit
		eventDeposit, err := gateway.ParseDeposited(*log)
		if err == nil {
			eventFound = true

			// check if the event is processable
			if !ob.isEventProcessable(
				eventDeposit.Sender,
				eventDeposit.Receiver,
				eventDeposit.Raw.TxHash,
				eventDeposit.Payload,
			) {
				return fmt.Errorf("event from inbound tracker %s is not processable", tx.Hash)
			}

			metrics.InboundObservationsTrackerTotal.WithLabelValues(ob.Chain().Name, strconv.FormatBool(isInternalTracker)).
				Inc()
			msg := ob.newDepositInboundVote(eventDeposit)
			_, err = ob.ZetaRepo().VoteInbound(ctx,
				ob.Logger().Inbound,
				&msg,
				zetacore.PostVoteInboundExecutionGasLimit,
				ob.WatchMonitoringError,
			)
			if err != nil || !allowMultipleCalls {
				return err
			}
		}

		// try parsing deposit and call
		eventDepositAndCall, err := gateway.ParseDepositedAndCalled(*log)
		if err == nil {
			eventFound = true

			// check if the event is processable
			if !ob.isEventProcessable(
				eventDepositAndCall.Sender,
				eventDepositAndCall.Receiver,
				eventDepositAndCall.Raw.TxHash,
				eventDepositAndCall.Payload,
			) {
				return fmt.Errorf("event from inbound tracker %s is not processable", tx.Hash)
			}
			metrics.InboundObservationsTrackerTotal.WithLabelValues(ob.Chain().Name, strconv.FormatBool(isInternalTracker)).
				Inc()
			msg := ob.newDepositAndCallInboundVote(eventDepositAndCall)
			_, err = ob.ZetaRepo().VoteInbound(ctx,
				ob.Logger().Inbound,
				&msg,
				zetacore.PostVoteInboundExecutionGasLimit,
				ob.WatchMonitoringError,
			)
			if err != nil || !allowMultipleCalls {
				return err
			}
		}

		// try parsing call
		eventCall, err := gateway.ParseCalled(*log)
		if err == nil {
			eventFound = true

			// check if the event is processable
			if !ob.isEventProcessable(
				eventCall.Sender,
				eventCall.Receiver,
				eventCall.Raw.TxHash,
				eventCall.Payload,
			) {
				return fmt.Errorf("event from inbound tracker %s is not processable", tx.Hash)
			}
			metrics.InboundObservationsTrackerTotal.WithLabelValues(ob.Chain().Name, strconv.FormatBool(isInternalTracker)).
				Inc()
			msg := ob.newCallInboundVote(eventCall)
			_, err = ob.ZetaRepo().VoteInbound(ctx,
				ob.Logger().Inbound,
				&msg,
				zetacore.PostVoteInboundExecutionGasLimit,
				ob.WatchMonitoringError,
			)
			if err != nil || !allowMultipleCalls {
				return err
			}
		}
	}

	if eventFound {
		return nil
	}

	return errors.Wrapf(ErrEventNotFound, "inbound tracker %s", tx.Hash)
}
