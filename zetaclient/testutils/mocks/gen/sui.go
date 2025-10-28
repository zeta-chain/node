package mocks

import (
	context "context"
	time "time"

	models "github.com/block-vision/sui-go-sdk/models"
	suiptb "github.com/pattonkan/sui-go/sui"

	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
)

//go:generate mockery --name suiClient --structname SuiClient --filename sui_client.go --output ../

// suiClient represents interface version of Client.
//
// The interface is unexported on purpose, since we ONLY use it for mock generation.
//
//nolint:unused // used for code gen
type suiClient interface {
	HealthCheck(ctx context.Context) (time.Time, error)
	GetLatestCheckpoint(ctx context.Context) (models.CheckpointResponse, error)
	QueryModuleEvents(ctx context.Context, q client.EventQuery) ([]models.SuiEventResponse, string, error)
	GetOwnedObjectID(ctx context.Context, ownerAddress, structType string) (string, error)
	GetObjectParsedData(ctx context.Context, objectID string) (models.SuiParsedData, error)
	GetSuiCoinObjectRefs(ctx context.Context, owner string, minBalanceMist uint64) ([]*suiptb.ObjectRef, error)

	SuiXGetLatestSuiSystemState(ctx context.Context) (models.SuiSystemStateSummary, error)
	SuiXGetDynamicFieldObject(
		ctx context.Context,
		req models.SuiXGetDynamicFieldObjectRequest,
	) (models.SuiObjectResponse, error)
	SuiXGetReferenceGasPrice(ctx context.Context) (uint64, error)
	SuiXQueryEvents(ctx context.Context, req models.SuiXQueryEventsRequest) (models.PaginatedEventsResponse, error)
	SuiMultiGetObjects(ctx context.Context, req models.SuiMultiGetObjectsRequest) ([]*models.SuiObjectResponse, error)
	SuiGetTransactionBlock(
		ctx context.Context,
		req models.SuiGetTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
	MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error)
	InspectTransactionBlock(
		ctx context.Context,
		req models.SuiDevInspectTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
	SuiExecuteTransactionBlock(
		ctx context.Context,
		req models.SuiExecuteTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
}
