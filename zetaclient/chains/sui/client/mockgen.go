package client

import (
	"context"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
)

// client represents interface version of Client.
// It's unexported on purpose ONLY for mock generation.
//
//go:generate mockery --name client --structname SUIClient --filename sui_client.go --output ../../../testutils/mocks
type client interface {
	HealthCheck(ctx context.Context) (time.Time, error)
	GetLatestCheckpoint(ctx context.Context) (models.CheckpointResponse, error)

	SuiXGetReferenceGasPrice(ctx context.Context) (uint64, error)
}
