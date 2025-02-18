package observer

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
)

// Observer Sui observer
type Observer struct {
	*base.Observer
	client  RPC
	gateway *sui.Gateway
}

// RPC represents subset of Sui RPC methods.
type RPC interface {
	HealthCheck(ctx context.Context) (time.Time, error)
	GetLatestCheckpoint(ctx context.Context) (models.CheckpointResponse, error)
	QueryModuleEvents(ctx context.Context, q client.EventQuery) ([]models.SuiEventResponse, string, error)

	SuiXGetReferenceGasPrice(ctx context.Context) (uint64, error)
	SuiGetObject(ctx context.Context, req models.SuiGetObjectRequest) (models.SuiObjectResponse, error)
	SuiGetTransactionBlock(
		ctx context.Context,
		req models.SuiGetTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
}

// New Observer constructor.
func New(baseObserver *base.Observer, client RPC, gateway *sui.Gateway) *Observer {
	return &Observer{
		Observer: baseObserver,
		client:   client,
		gateway:  gateway,
	}
}

// CheckRPCStatus checks the RPC status of the chain.
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	blockTime, err := ob.client.HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to check rpc health")
	}

	// It's not a "real" block latency as Sui uses concept of "checkpoints"
	ob.ReportBlockLatency(blockTime)

	return nil
}

// PostGasPrice posts Sui gas price to zetacore.
// Note (1) that Sui changes gas per EPOCH (not block)
// Note (2) that SuiXGetCurrentEpoch() is deprecated.
//
// See https://docs.sui.io/concepts/tokenomics/gas-pricing
// See https://docs.sui.io/concepts/sui-architecture/transaction-lifecycle#epoch-change
//
// TLDR:
// - GasFees = CompUnits * (ReferencePrice + Tip) + StorageUnits * StoragePrice
// - "During regular network usage, users are NOT expected to pay tips"
// - "Validators update the ReferencePrice every epoch (~24h)"
// - "Storage price is updated infrequently through gov proposals"
func (ob *Observer) PostGasPrice(ctx context.Context) error {
	checkpoint, err := ob.client.GetLatestCheckpoint(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get latest checkpoint")
	}

	epochNum, err := uint64FromStr(checkpoint.Epoch)
	if err != nil {
		return errors.Wrap(err, "unable to parse epoch number")
	}

	// gas price in MIST. 1 SUI = 10^9 MIST (a billion)
	// e.g. { "jsonrpc": "2.0", "id": 1, "result": "750" }
	gasPrice, err := ob.client.SuiXGetReferenceGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get ref gas price")
	}

	// no priority fee for Sui
	const priorityFee = 0

	_, err = ob.ZetacoreClient().PostVoteGasPrice(ctx, ob.Chain(), gasPrice, priorityFee, epochNum)
	if err != nil {
		return errors.Wrap(err, "unable to post vote for gas price")
	}

	return nil
}

// ensureCursor ensures tx scroll cursor for inbound observations
func (ob *Observer) ensureCursor(ctx context.Context) error {
	if ob.LastTxScanned() != "" {
		return nil
	}

	// Note that this would only work for the empty chain database
	envValue := os.Getenv(base.EnvVarLatestTxByChain(ob.Chain()))
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

func uint64FromStr(raw string) (uint64, error) {
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to parse uint64 from %s", raw)
	}

	return v, nil
}
