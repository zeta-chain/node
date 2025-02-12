package mocks

import (
	context "context"
	time "time"

	models "github.com/block-vision/sui-go-sdk/models"
)

// client represents interface version of Client.
// It's unexported on purpose ONLY for mock generation.
//
//go:generate mockery --name suiClient --structname SuiClient --filename sui_client.go --output ./
//nolint:unused // used for code gen
type suiClient interface {
	HealthCheck(ctx context.Context) (time.Time, error)
	GetLatestCheckpoint(ctx context.Context) (models.CheckpointResponse, error)

	SuiXGetReferenceGasPrice(ctx context.Context) (uint64, error)
	SuiXQueryEvents(ctx context.Context, req models.SuiXQueryEventsRequest) (models.PaginatedEventsResponse, error)
}
