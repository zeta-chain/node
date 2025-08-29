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
	client  RPC
	gateway *sui.Gateway

	// nonce -> sui outbound tx
	txMap map[uint64]models.SuiTransactionBlockResponse
	txMu  sync.RWMutex

	latestGasPrice uint64
	gasPriceMu     sync.RWMutex
}

// RPC represents subset of Sui RPC methods.
type RPC interface {
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
func New(baseObserver *base.Observer, client RPC, gateway *sui.Gateway) *Observer {
	ob := &Observer{
		Observer: baseObserver,
		client:   client,
		gateway:  gateway,
		txMap:    make(map[uint64]models.SuiTransactionBlockResponse),
	}

	ob.LoadLastTxScanned()

	return ob
}

// Gateway returns Sui gateway.
func (ob *Observer) Gateway() *sui.Gateway { return ob.gateway }

// CheckRPCStatus checks the RPC status of the chain.
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	blockTime, err := ob.client.HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to check rpc health")
	}

	// It's not a "real" block latency as Sui uses concept of "checkpoints"
	metrics.ReportBlockLatency(ob.Chain().Name, blockTime)

	return nil
}

// MigrateInboundCursorV35 migrates old inbound cursor in the database
// zetaclient v35 needs two separate inbound cursors for old and new gateway packages,
// so the cursors have to be stored under separate keys - the package IDs
func (ob *Observer) MigrateInboundCursorV35() error {
	// all we need to do is to migrate the cursor for original package
	// the old cursor is stored as 'LastTransactionSQLType'
	oldCursor := ob.LastTxScanned()
	if oldCursor == "" {
		// nothing to migrate
		return nil
	}

	// Sui chain params may or may not contain new gateway package ID
	// the 'originalPackageID' should be used as the DB key for old cursor
	originalPackageID := ob.gateway.Original().PackageID()
	if err := ob.WriteAnyStringToDB(originalPackageID, oldCursor); err != nil {
		return errors.Wrapf(err, "unable to write old cursor to db for package %s", originalPackageID)
	}
	ob.WithAnyString(originalPackageID, oldCursor)

	ob.Logger().
		Inbound.Info().
		Str("package", originalPackageID).
		Str("cursor", oldCursor).
		Msgf("Migrated Sui inbound cursor")

	// clean up old cursor from DB
	if err := ob.WriteLastTxScannedToDB(""); err != nil {
		return errors.Wrap(err, "unable to clean last tx scanned from db")
	}
	ob.WithLastTxScanned("")

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

// getCursor retrieves the inbound cursor for a given packageID
func (ob *Observer) getCursor(packageID string) string {
	return ob.GetAnyString(packageID)
}

// setCursor saves the inbound cursor for a given packageID
func (ob *Observer) setCursor(packageID string, eventID models.EventId) error {
	cursor := client.EncodeCursor(eventID)

	if err := ob.WriteAnyStringToDB(packageID, cursor); err != nil {
		return errors.Wrapf(err, "unable to write cursor to db for package %s", packageID)
	}
	ob.WithAnyString(packageID, cursor)

	return nil
}

func uint64FromStr(raw string) (uint64, error) {
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to parse uint64 from %s", raw)
	}

	return v, nil
}
