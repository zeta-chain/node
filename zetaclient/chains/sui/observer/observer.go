package observer

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// Observer Sui observer
type Observer struct {
	*base.Observer

	suiClient SuiClient

	gateway *sui.Gateway

	// nonce -> sui outbound tx
	txMap map[uint64]models.SuiTransactionBlockResponse
	txMu  sync.RWMutex

	latestGasPrice uint64
	gasPriceMu     sync.RWMutex
}

// SuiClient represents subset of Sui SuiClient methods.
type SuiClient interface {
	HealthCheck(ctx context.Context) (time.Time, error)
	GetLatestCheckpoint(ctx context.Context) (models.CheckpointResponse, error)
	QueryModuleEvents(ctx context.Context, q client.EventQuery) ([]models.SuiEventResponse, string, error)

	SuiXGetReferenceGasPrice(ctx context.Context) (uint64, error)
	SuiGetTransactionBlock(
		ctx context.Context,
		req models.SuiGetTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
}

// New Observer constructor.
func New(baseObserver *base.Observer, suiClient SuiClient, gateway *sui.Gateway) *Observer {
	ob := &Observer{
		Observer:  baseObserver,
		suiClient: suiClient,
		gateway:   gateway,
		txMap:     make(map[uint64]models.SuiTransactionBlockResponse),
	}

	ob.LoadLastTxScanned()

	return ob
}

// Gateway returns Sui gateway.
func (ob *Observer) Gateway() *sui.Gateway { return ob.gateway }

// CheckRPCStatus checks the RPC status of the chain.
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	blockTime, err := ob.suiClient.HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to check Sui client health")
	}

	// It's not a "real" block latency as Sui uses concept of "checkpoints"
	metrics.ReportBlockLatency(ob.Chain().Name, blockTime)

	return nil
}

// ObserveGasPrice posts Sui gas price to zetacore.
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
func (ob *Observer) ObserveGasPrice(ctx context.Context) error {
	checkpoint, err := ob.suiClient.GetLatestCheckpoint(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get latest checkpoint")
	}

	epoch, err := uint64FromStr(checkpoint.Epoch)
	if err != nil {
		return errors.Wrap(err, "unable to parse epoch number")
	}

	// gas price in MIST. 1 SUI = 10^9 MIST (a billion)
	// e.g. { "jsonrpc": "2.0", "id": 1, "result": "750" }
	gasPrice, err := ob.suiClient.SuiXGetReferenceGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get ref gas price")
	}

	// There's no concept of priority fee in Sui.
	const priorityFee = 0

	logger := ob.Logger().Chain
	_, err = ob.ZetaRepo().VoteGasPrice(ctx, logger, gasPrice, priorityFee, epoch)
	if err != nil {
		return err
	}

	ob.setLatestGasPrice(gasPrice)

	return nil
}

func (ob *Observer) getLatestGasPrice() uint64 {
	ob.gasPriceMu.RLock()
	defer ob.gasPriceMu.RUnlock()

	return ob.latestGasPrice
}

func (ob *Observer) setLatestGasPrice(price uint64) {
	ob.gasPriceMu.Lock()
	defer ob.gasPriceMu.Unlock()
	ob.latestGasPrice = price
}

func (ob *Observer) getCursor() string { return ob.LastTxScanned() }

func (ob *Observer) setCursor(eventID models.EventId) error {
	cursor := client.EncodeCursor(eventID)

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
