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

// GetNumberOfUnconfirmedTxs returns the number of unconfirmed txs in the zetacore mempool
func (c *Client) GetNumberOfUnconfirmedTxs(ctx context.Context) (int, error) {
	resp, err := c.cometBFTClient.NumUnconfirmedTxs(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get number of unconfirmed txs")
	}

	return resp.Count, nil
}

func (c *Client) GetSyncStatus(ctx context.Context) (bool, error) {
	syncing, err := c.Clients.GetSyncing(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get syncing status")
	}
	return syncing, nil
}
