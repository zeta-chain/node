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

type Client struct {
	sui.ISuiAPI
}

// NewFromEndpoint Client constructor based on endpoint string.
func NewFromEndpoint(endpoint string) *Client {
	return New(sui.NewSuiClient(endpoint))
}

// New Client constructor.
func New(client sui.ISuiAPI) *Client {
	return &Client{ISuiAPI: client}
}

// Queries latest seq no and returns its timestamp.
func (c *Client) HealthCheck(ctx context.Context) (time.Time, error) {
	seqNum, err := c.SuiGetLatestCheckpointSequenceNumber(ctx)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to get latest seq num")
	}

	req := models.SuiGetCheckpointRequest{
		CheckpointID: fmt.Sprintf("%d", seqNum),
	}

	checkpoint, err := c.SuiGetCheckpoint(ctx, req)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "unable to get checkpoint %d", seqNum)
	}

	ts, err := strconv.ParseInt(checkpoint.TimestampMs, 10, 64)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to parse checkpoint timestamp")
	}

	return time.UnixMilli(ts).UTC(), nil
}
