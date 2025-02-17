package observer

import (
	"context"
	"encoding/hex"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

var errTxNotFound = errors.New("no tx found")

// ObserveInbound processes inbound deposit cross-chain transactions.
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	if err := ob.ensureCursor(ctx); err != nil {
		return errors.Wrap(err, "unable to ensure inbound cursor")
	}

	query := client.EventQuery{
		PackageID: ob.gateway.PackageID(),
		Module:    ob.gateway.Module(),
		Cursor:    ob.getCursor(),
		Limit:     client.DefaultEventsLimit,
	}

	// Sui has a nice access-pattern of scrolling through contract events
	events, _, err := ob.client.QueryModuleEvents(ctx, query)
	if err != nil {
		return errors.Wrap(err, "unable to query module events")
	}

	for _, event := range events {
		// Note: we can make this concurrent if needed.
		// Let's revisit later
		err := ob.processInboundEvent(ctx, event)

		switch {
		case errors.Is(err, errTxNotFound):
			// try again later
			ob.Logger().Inbound.Warn().Err(err).
				Str(logs.FieldTx, event.Id.TxDigest).
				Msg("TX not found or unfinalized. Pausing")
			return nil
		case err != nil:
			// failed processing also updates the cursor
			ob.Logger().Inbound.Err(err).
				Str(logs.FieldTx, event.Id.TxDigest).
				Msg("Unable to process inbound event")
		}

		// update the cursor
		if err := ob.setCursor(client.EncodeCursor(event.Id)); err != nil {
			return errors.Wrapf(err, "unable to set cursor %+v", event.Id)
		}
	}

	return nil
}

// processInboundEvent parses raw event into Inbound,
// augment it with origin tx and vote on the inbound.
// Invalid/Non-inbound txs are skipped. Unconfirmed txs pause the whole tail sequence.
func (ob *Observer) processInboundEvent(ctx context.Context, raw models.SuiEventResponse) error {
	event, err := ob.gateway.ParseEvent(raw)
	switch {
	case errors.Is(err, sui.ErrParseEvent):
		ob.Logger().Inbound.Err(err).Msg("Unable to parse event. Skipping")
		return nil
	case err != nil:
		return errors.Wrap(err, "unable to parse event")
	case !event.IsInbound():
		ob.Logger().Inbound.Info().Msg("Not an inbound event. Skipping")
	case event.EventIndex != 0:
		// Is it possible to have multiple events per tx?
		// e.g. contract "A" calls Gateway multiple times in a single tx (deposit to multiple accounts)
		// most likely not, so let's explicitly fail to prevent undefined behavior.
		return errors.Errorf("unexpected event index %d for tx %s", event.EventIndex, event.TxHash)
	}

	txReq := models.SuiGetTransactionBlockRequest{Digest: event.TxHash}

	tx, err := ob.client.SuiGetTransactionBlock(ctx, txReq)
	if err != nil {
		return errors.Wrap(errTxNotFound, err.Error())
	}

	msg, err := ob.constructInboundVote(event, tx)
	if err != nil {
		return errors.Wrap(err, "unable to construct inbound vote")
	}

	_, err = ob.PostVoteInbound(ctx, msg, zetacore.PostVoteInboundExecutionGasLimit)
	if err != nil {
		return errors.Wrap(err, "unable to post vote inbound")
	}

	return nil
}

// constructInboundVote creates a vote message for inbound deposit
func (ob *Observer) constructInboundVote(
	event sui.Event,
	tx models.SuiTransactionBlockResponse,
) (*cctypes.MsgVoteInbound, error) {
	inbound, err := event.Inbound()
	if err != nil {
		return nil, errors.Wrap(err, "unable to extract inbound")
	}

	coinType := coin.CoinType_Gas
	if !inbound.IsGasDeposit() {
		coinType = coin.CoinType_ERC20
	}

	// Sui uses checkpoint seq num instead of block height
	checkpointSeqNum, err := uint64FromStr(tx.Checkpoint)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse checkpoint")
	}

	return cctypes.NewMsgVoteInbound(
		ob.ZetacoreClient().GetKeys().GetOperatorAddress().String(),
		inbound.Sender,
		ob.Chain().ChainId,
		inbound.Sender,
		inbound.Receiver.String(),
		ob.ZetacoreClient().Chain().ChainId,
		inbound.Amount,
		hex.EncodeToString(inbound.Memo()),
		event.TxHash,
		checkpointSeqNum,
		0,
		coinType,
		string(inbound.CoinType),
		uint(event.EventIndex),
		cctypes.ProtocolContractVersion_V2,
		false,
		cctypes.InboundStatus_SUCCESS,
		cctypes.ConfirmationMode_SAFE,
		cctypes.WithCrossChainCall(inbound.IsCrossChainCall),
	), nil
}

// ensureCursor ensures tx scroll cursor for inbound observations
func (ob *Observer) ensureCursor(ctx context.Context) error {
	if ob.LastTxScanned() == "" {
		return nil
	}

	// Note that this would only work for the empty chain database
	envValue := base.EnvVarLatestTxByChain(ob.Chain())
	if envValue != "" {
		ob.WithLastTxScanned(envValue)
		return nil
	}

	// let's take the first tx that was ever registered for the Gateway (deployment tx)
	// Note that this might have for a non-archival node
	req := models.SuiGetObjectRequest{
		ObjectId: ob.gateway.PackageID(),
		Options: models.SuiObjectDataOptions{
			ShowPreviousTransaction: true,
		},
	}

	res, err := ob.client.SuiGetObject(ctx, req)
	switch {
	case err != nil:
		return errors.Wrap(err, "unable to get object")
	case res.Error != nil:
		return errors.Errorf("get object error: %s (code %s)", res.Error.Error, res.Error.Code)
	case res.Data == nil:
		return errors.New("object data is empty")
	case res.Data.PreviousTransaction == "":
		return errors.New("previous transaction is empty")
	}

	cursor := client.EncodeCursor(models.EventId{
		TxDigest: res.Data.PreviousTransaction,
		EventSeq: "0",
	})

	return ob.setCursor(cursor)
}

func (ob *Observer) getCursor() string { return ob.LastTxScanned() }

func (ob *Observer) setCursor(cursor string) error {
	if err := ob.WriteLastTxScannedToDB(cursor); err != nil {
		return errors.Wrap(err, "unable to write last tx scanned to db")
	}

	ob.WithLastTxScanned(cursor)

	return nil
}
