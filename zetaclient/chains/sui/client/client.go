package client

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/pkg/errors"
)

// Client SUI client.
type Client struct {
	sui.ISuiAPI
}

var _ client = (*Client)(nil)

// NewFromEndpoint Client constructor based on endpoint string.
func NewFromEndpoint(endpoint string) *Client {
	return New(sui.NewSuiClient(endpoint))
}

// New Client constructor.
func New(client sui.ISuiAPI) *Client {
	return &Client{ISuiAPI: client}
}

// HealthCheck queries latest checkpoint and returns its timestamp.
func (c *Client) HealthCheck(ctx context.Context) (time.Time, error) {
	checkpoint, err := c.GetLatestCheckpoint(ctx)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "unable to get latest checkpoint")
	}

	ts, err := strconv.ParseInt(checkpoint.TimestampMs, 10, 64)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to parse checkpoint timestamp")
	}

	return time.UnixMilli(ts).UTC(), nil
}

// GetLatestCheckpoint returns the latest checkpoint.
// See https://docs.sui.io/concepts/cryptography/system/checkpoint-verification
func (c *Client) GetLatestCheckpoint(ctx context.Context) (models.CheckpointResponse, error) {
	seqNum, err := c.SuiGetLatestCheckpointSequenceNumber(ctx)
	if err != nil {
		return models.CheckpointResponse{}, errors.Wrap(err, "unable to get latest seq num")
	}

	return c.SuiGetCheckpoint(ctx, models.SuiGetCheckpointRequest{
		CheckpointID: fmt.Sprintf("%d", seqNum),
	})
}
