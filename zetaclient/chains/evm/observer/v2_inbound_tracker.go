package observer

import (
	"context"
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onrik/ethrpc"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/zetacore"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrGatewayNotSet = errors.New("gateway contract not set")
)

// ProcessInboundTrackerV2 processes inbound tracker events from the gateway
func (ob *Observer) ProcessInboundTrackerV2(
	ctx context.Context,
	tx *ethrpc.Transaction,
	receipt *ethtypes.Receipt,
) error {
	gatewayAddr, gateway, err := ob.GetGatewayContract()
	if err != nil {
		ob.Logger().Inbound.Debug().Err(err).Msg("error getting gateway contract for processing inbound tracker")
		return ErrGatewayNotSet
	}

	// check confirmations
	if confirmed := ob.HasEnoughConfirmations(receipt, ob.LastBlock()); !confirmed {
		return fmt.Errorf(
			"inbound %s has not been confirmed yet: receipt block %d",
			tx.Hash,
			receipt.BlockNumber.Uint64(),
		)
	}

	for _, log := range receipt.Logs {
		if log == nil || log.Address != gatewayAddr {
			continue
		}

		// try parsing deposit
		eventDeposit, err := gateway.ParseDeposited(*log)
		if err == nil {
			// check if the event is processable
			if !ob.checkEventProcessability(
				eventDeposit.Sender,
				eventDeposit.Receiver,
				eventDeposit.Raw.TxHash,
				eventDeposit.Payload,
			) {
				return fmt.Errorf("event from inbound tracker %s is not processable", tx.Hash)
			}
			msg := ob.newDepositInboundVote(eventDeposit)
			_, err = ob.PostVoteInbound(ctx, &msg, zetacore.PostVoteInboundExecutionGasLimit)
			return err
		}

		// try parsing deposit and call
		eventDepositAndCall, err := gateway.ParseDepositedAndCalled(*log)
		if err == nil {
			// check if the event is processable
			if !ob.checkEventProcessability(
				eventDepositAndCall.Sender,
				eventDepositAndCall.Receiver,
				eventDepositAndCall.Raw.TxHash,
				eventDepositAndCall.Payload,
			) {
				return fmt.Errorf("event from inbound tracker %s is not processable", tx.Hash)
			}
			msg := ob.newDepositAndCallInboundVote(eventDepositAndCall)
			_, err = ob.PostVoteInbound(ctx, &msg, zetacore.PostVoteInboundExecutionGasLimit)
			return err
		}

		// try parsing call
		eventCall, err := gateway.ParseCalled(*log)
		if err == nil {
			// check if the event is processable
			if !ob.checkEventProcessability(
				eventCall.Sender,
				eventCall.Receiver,
				eventCall.Raw.TxHash,
				eventCall.Payload,
			) {
				return fmt.Errorf("event from inbound tracker %s is not processable", tx.Hash)
			}
			msg := ob.newCallInboundVote(eventCall)
			_, err = ob.PostVoteInbound(ctx, &msg, zetacore.PostVoteInboundExecutionGasLimit)
			return err
		}
	}

	return errors.Wrapf(ErrEventNotFound, "inbound tracker %s", tx.Hash)
}
