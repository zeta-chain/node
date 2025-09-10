package observer

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

var (
	errTxNotFound = errors.New("no tx found")
	errCompliance = errors.New("compliance check failed")
)

// ObserveInbound processes inbound deposit cross-chain transactions.
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	// always query inbound events from original gateway package
	// querying events from upgraded packageID will get nothing
	packageID := ob.gateway.Original().PackageID()
	if err := ob.observeGatewayInbound(ctx, packageID); err != nil {
		return errors.Wrap(err, "unable to observe gateway inbound")
	}

	return nil
}

// observeGatewayInbound observes inbound deposits for the given gateway packageID
// The last processed event will be used as the next cursor
func (ob *Observer) observeGatewayInbound(ctx context.Context, packageID string) error {
	cursor := ob.getCursor(packageID)
	query := client.EventQuery{
		PackageID: packageID,
		Module:    sui.GatewayModule,
		Cursor:    cursor,
		Limit:     client.DefaultEventsLimit,
	}

	// Sui has a nice access-pattern of scrolling through contract events
	events, _, err := ob.client.QueryModuleEvents(ctx, query)
	if err != nil {
		return errors.Wrap(err, "unable to query module events")
	}

	if len(events) == 0 {
		ob.Logger().Inbound.Debug().Msg("found no inbound events")
		return nil
	}

	ob.Logger().
		Inbound.Info().
		Str("package", packageID).
		Str("cursor", cursor).
		Int("events", len(events)).
		Msg("Processing inbound events")

	for _, event := range events {
		// Note: we can make this concurrent if needed.
		// Let's revisit later
		err := ob.processInboundEvent(ctx, event, nil)

		switch {
		case errors.Is(err, errTxNotFound):
			// try again later
			ob.Logger().Inbound.Warn().Err(err).
				Str(logs.FieldTx, event.Id.TxDigest).
				Msg("tx not found or not finalized; pausing")
			return nil
		case errors.Is(err, errCompliance):
			// skip restricted tx and update the cursor
			ob.Logger().Inbound.Warn().Err(err).
				Str(logs.FieldTx, event.Id.TxDigest).
				Msg("tx contains restricted address; skipping")
		case err != nil:
			// failed processing also updates the cursor
			ob.Logger().Inbound.Err(err).
				Str(logs.FieldTx, event.Id.TxDigest).
				Msg("unable to process inbound event")
		}

		// update the cursor
		if err := ob.setCursor(packageID, event.Id); err != nil {
			return errors.Wrapf(err, "unable to set cursor %+v", event.Id)
		}
	}

	return nil
}

// ProcessInboundTrackers processes trackers for inbound transactions.
func (ob *Observer) ProcessInboundTrackers(ctx context.Context) error {
	chainID := ob.Chain().ChainId

	trackers, err := ob.ZetacoreClient().GetInboundTrackersForChain(ctx, chainID)
	if err != nil {
		return errors.Wrap(err, "unable to get inbound trackers")
	}

	for _, tracker := range trackers {
		if err := ob.processInboundTracker(ctx, tracker); err != nil {
			ob.Logger().Inbound.Err(err).
				Str(logs.FieldTx, tracker.TxHash).
				Msg("unable to process inbound tracker")
		}
	}

	return nil
}

// processInboundEvent parses raw event into Inbound, augments it with origin tx and votes on the inbound.
// - Invalid/Non-inbound txs are skipped.
// - Unconfirmed txs pause the whole tail sequence.
// - If tx is empty, it fetches the tx from RPC.
// - Sui tx is finalized if it's returned from RPC
//
// See https://docs.sui.io/concepts/sui-architecture/transaction-lifecycle#verifying-finality
func (ob *Observer) processInboundEvent(
	ctx context.Context,
	raw models.SuiEventResponse,
	tx *models.SuiTransactionBlockResponse,
) error {
	event, err := ob.gateway.ParseEvent(raw)
	switch {
	case errors.Is(err, sui.ErrParseEvent):
		ob.Logger().Inbound.Err(err).Msg("unable to parse event; skipping")
		return nil
	case err != nil:
		return errors.Wrap(err, "unable to parse event")
	case event.EventIndex != 0:
		// Is it possible to have multiple events per tx?
		// e.g. contract "A" calls Gateway multiple times in a single tx (deposit to multiple accounts)
		// most likely not, so let's explicitly fail to prevent undefined behavior.
		return errors.Errorf("unexpected event index %d for tx %s", event.EventIndex, event.TxHash)
	}

	if tx == nil {
		txReq := models.SuiGetTransactionBlockRequest{
			Digest:  event.TxHash,
			Options: models.SuiTransactionBlockOptions{ShowEffects: true},
		}
		txFresh, err := ob.client.SuiGetTransactionBlock(ctx, txReq)
		if err != nil {
			return errors.Wrap(errTxNotFound, err.Error())
		}

		tx = &txFresh
	}

	msg, err := ob.constructInboundVote(event, *tx)
	if err != nil {
		return errors.Wrap(err, "unable to construct inbound vote")
	}

	_, err = ob.PostVoteInbound(ctx, msg, zetacore.PostVoteInboundExecutionGasLimit)
	if err != nil {
		return errors.Wrap(err, "unable to post vote inbound")
	}

	return nil
}

// processInboundTracker queries tx with its events by tracker and then votes.
func (ob *Observer) processInboundTracker(ctx context.Context, tracker cctypes.InboundTracker) error {
	req := models.SuiGetTransactionBlockRequest{
		Digest: tracker.TxHash,
		Options: models.SuiTransactionBlockOptions{
			ShowEffects: true,
			ShowEvents:  true,
		},
	}

	tx, err := ob.client.SuiGetTransactionBlock(ctx, req)
	if err != nil {
		return errors.Wrapf(err, "unable to get transaction block")
	}

	for _, event := range tx.Events {
		if err := ob.processInboundEvent(ctx, event, &tx); err != nil {
			return errors.Wrapf(err, "unable to process inbound event %s", event.Id.EventSeq)
		}
	}

	return nil
}

// constructInboundVote creates a vote message for inbound deposit
func (ob *Observer) constructInboundVote(
	event sui.Event,
	tx models.SuiTransactionBlockResponse,
) (*cctypes.MsgVoteInbound, error) {
	deposit, err := event.Deposit()
	if err != nil {
		return nil, errors.Wrap(err, "unable to extract inbound")
	}

	coinType := coin.CoinType_Gas
	asset := ""
	if !deposit.IsGas() {
		coinType = coin.CoinType_ERC20
		asset = string(deposit.CoinType)
	}

	// compliance check, skip restricted tx by returning nil msg
	if config.ContainRestrictedAddress(deposit.Sender, deposit.Receiver.String()) {
		compliance.PrintComplianceLog(
			ob.Logger().Inbound,
			ob.Logger().Compliance,
			false,
			ob.Chain().ChainId,
			event.TxHash,
			deposit.Sender,
			deposit.Receiver.String(),
			asset,
		)
		return nil, errCompliance
	}

	// a valid inbound should be successful
	// in theory, Sui protocol should erase emitted events if tx failed, just in case
	if tx.Effects.Status.Status != client.TxStatusSuccess {
		return nil, errors.Errorf("inbound is failed: %s", tx.Effects.Status.Error)
	}

	// Sui uses checkpoint seq num instead of block height.
	// If checkpoint is invalid (e.g. 0), the tx status remains unclear (e.g. maybe pending).
	// In this case, we should signal the caller to stop scanning further by returning errTxNotFound.
	checkpointSeqNum, err := uint64FromStr(tx.Checkpoint)
	if err != nil || checkpointSeqNum == 0 {
		return nil, errors.Wrap(errTxNotFound, fmt.Sprintf("invalid checkpoint: %s", tx.Checkpoint))
	}

	return cctypes.NewMsgVoteInbound(
		ob.ZetacoreClient().GetKeys().GetOperatorAddress().String(),
		deposit.Sender,
		ob.Chain().ChainId,
		deposit.Sender,
		deposit.Receiver.String(),
		ob.ZetacoreClient().Chain().ChainId,
		deposit.Amount,
		hex.EncodeToString(deposit.Payload),
		event.TxHash,
		checkpointSeqNum,
		zetacore.PostVoteInboundCallOptionsGasLimit,
		coinType,
		asset,
		event.EventIndex,
		cctypes.ProtocolContractVersion_V2,
		false,
		cctypes.InboundStatus_SUCCESS,
		cctypes.ConfirmationMode_SAFE,
		cctypes.WithCrossChainCall(deposit.IsCrossChainCall),
	), nil
}
