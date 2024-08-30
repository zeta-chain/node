package rpc

import (
	"context"

	sdkmath "cosmossdk.io/math"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/cmd/zetacored/config"
)

// GetUpgradePlan returns the current upgrade plan.
// if there is no active upgrade plan, plan will be nil, err will be nil as well.
func (c *Clients) GetUpgradePlan(ctx context.Context) (*upgradetypes.Plan, error) {
	in := &upgradetypes.QueryCurrentPlanRequest{}

	resp, err := c.Upgrade.CurrentPlan(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current upgrade plan")
	}

	return resp.Plan, nil
}

// GetZetaTokenSupplyOnNode returns the zeta token supply on the node
func (c *Clients) GetZetaTokenSupplyOnNode(ctx context.Context) (sdkmath.Int, error) {
	in := &banktypes.QuerySupplyOfRequest{Denom: config.BaseDenom}

	resp, err := c.Bank.SupplyOf(ctx, in)
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "failed to get zeta token supply")
	}

	return resp.GetAmount().Amount, nil
}
