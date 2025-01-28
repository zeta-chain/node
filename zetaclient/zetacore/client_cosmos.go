package zetacore

import (
	"context"

	sdkmath "cosmossdk.io/math"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/cmd/zetacored/config"
)

// GetZetaHotKeyBalance returns the zeta hot key balance
func (c *Client) GetZetaHotKeyBalance(ctx context.Context) (sdkmath.Int, error) {
	address, err := c.keys.GetAddress()
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "failed to get address")
	}

	in := &banktypes.QueryBalanceRequest{
		Address: address.String(),
		Denom:   config.BaseDenom,
	}

	resp, err := c.Clients.Bank.Balance(ctx, in)
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "failed to get zeta hot key balance")
	}

	return resp.Balance.Amount, nil
}
