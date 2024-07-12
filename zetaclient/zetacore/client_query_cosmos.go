package zetacore

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	tmhttp "github.com/cometbft/cometbft/rpc/client/http"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/pkg/errors"

	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
)

// GetGenesisSupply returns the genesis supply.
// NOTE that this method is brittle as it uses STATEFUL connection
func (c *Client) GetGenesisSupply(ctx context.Context) (sdkmath.Int, error) {
	tmURL := fmt.Sprintf("http://%s", c.config.ChainRPC)

	s, err := tmhttp.New(tmURL, "/websocket")
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "failed to create tm client")
	}

	// nolint:errcheck
	defer s.Stop()

	res, err := s.Genesis(ctx)
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "failed to get genesis")
	}

	appState, err := genutiltypes.GenesisStateFromGenDoc(*res.Genesis)
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "failed to get app state")
	}

	bankstate := banktypes.GetGenesisStateFromAppState(c.encodingCfg.Codec, appState)

	return bankstate.Supply.AmountOf(config.BaseDenom), nil
}

// GetUpgradePlan returns the current upgrade plan.
// if there is no active upgrade plan, plan will be nil, err will be nil as well.
func (c *Client) GetUpgradePlan(ctx context.Context) (*upgradetypes.Plan, error) {
	in := &upgradetypes.QueryCurrentPlanRequest{}

	resp, err := c.client.upgrade.CurrentPlan(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current upgrade plan")
	}

	return resp.Plan, nil
}

// GetZetaTokenSupplyOnNode returns the zeta token supply on the node
func (c *Client) GetZetaTokenSupplyOnNode(ctx context.Context) (sdkmath.Int, error) {
	in := &banktypes.QuerySupplyOfRequest{Denom: config.BaseDenom}

	resp, err := c.client.bank.SupplyOf(ctx, in)
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "failed to get zeta token supply")
	}

	return resp.GetAmount().Amount, nil
}

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

	resp, err := c.client.bank.Balance(ctx, in)
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "failed to get zeta hot key balance")
	}

	return resp.Balance.Amount, nil
}
