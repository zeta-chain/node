package zetacore

import (
	"context"

	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func (c *Client) GetOperationalFlags(ctx context.Context) (observertypes.OperationalFlags, error) {
	res, err := c.Observer.OperationalFlags(ctx, &observertypes.QueryOperationalFlagsRequest{})
	return res.OperationalFlags, err
}
