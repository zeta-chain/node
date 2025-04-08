package mocks

import (
	context "context"
	time "time"

	models "github.com/block-vision/sui-go-sdk/models"

	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
)

// client represents interface version of Client.
// It's unexported on purpose ONLY for mock generation.
//
//go:generate mockery --name suiClient --structname SuiClient --filename sui_client.go --output ./
//nolint:unused // used for code gen
type suiClient interface {
	HealthCheck(ctx context.Context) (time.Time, error)
	GetLatestCheckpoint(ctx context.Context) (models.CheckpointResponse, error)
	QueryModuleEvents(ctx context.Context, q client.EventQuery) ([]models.SuiEventResponse, string, error)
	GetOwnedObjectID(ctx context.Context, ownerAddress, structType string) (string, error)

	SuiXGetReferenceGasPrice(ctx context.Context) (uint64, error)
	SuiXQueryEvents(ctx context.Context, req models.SuiXQueryEventsRequest) (models.PaginatedEventsResponse, error)
	SuiGetObject(ctx context.Context, req models.SuiGetObjectRequest) (models.SuiObjectResponse, error)
	SuiGetTransactionBlock(
		ctx context.Context,
		req models.SuiGetTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
	MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error)
	SuiExecuteTransactionBlock(
		ctx context.Context,
		req models.SuiExecuteTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
}
