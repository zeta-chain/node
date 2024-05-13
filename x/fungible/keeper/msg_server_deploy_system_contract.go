package keeper

import (
	"context"

	cosmoserror "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// DeploySystemContracts deploy new instances of the system contracts
//
// Authorized: admin policy group 2.
func (k msgServer) DeploySystemContracts(goCtx context.Context, msg *types.MsgDeploySystemContracts) (*types.MsgDeploySystemContractsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserror.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// uniswap v2 factory
	factory, err := k.DeployUniswapV2Factory(ctx)
	if err != nil {
		return nil, cosmoserror.Wrapf(err, "failed to deploy UniswapV2Factory")
	}

	// wzeta contract
	wzeta, err := k.DeployWZETA(ctx)
	if err != nil {
		return nil, cosmoserror.Wrapf(err, "failed to DeployWZetaContract")
	}

	// uniswap v2 router
	router, err := k.DeployUniswapV2Router02(ctx, factory, wzeta)
	if err != nil {
		return nil, cosmoserror.Wrapf(err, "failed to deploy UniswapV2Router02")
	}

	// connector zevm
	connector, err := k.DeployConnectorZEVM(ctx, wzeta)
	if err != nil {
		return nil, cosmoserror.Wrapf(err, "failed to deploy ConnectorZEVM")
	}

	// system contract
	systemContract, err := k.DeploySystemContract(ctx, wzeta, factory, router)
	if err != nil {
		return nil, cosmoserror.Wrapf(err, "failed to deploy SystemContract")
	}

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventSystemContractsDeployed{
			MsgTypeUrl:       sdk.MsgTypeURL(&types.MsgDeploySystemContracts{}),
			UniswapV2Factory: factory.Hex(),
			Wzeta:            wzeta.Hex(),
			UniswapV2Router:  router.Hex(),
			ConnectorZevm:    connector.Hex(),
			SystemContract:   systemContract.Hex(),
			Signer:           msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event",
			"event", "EventSystemContractsDeployed",
			"error", err.Error(),
		)
		return nil, cosmoserror.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}

	return &types.MsgDeploySystemContractsResponse{
		UniswapV2Factory: factory.Hex(),
		Wzeta:            wzeta.Hex(),
		UniswapV2Router:  router.Hex(),
		ConnectorZEVM:    connector.Hex(),
		SystemContract:   systemContract.Hex(),
	}, nil
}
