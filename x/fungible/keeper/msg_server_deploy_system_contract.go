package keeper

import (
	"context"

	cosmoserror "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// DeploySystemContracts deploy new instances of the system contracts
func (k msgServer) DeploySystemContracts(goCtx context.Context, _ *types.MsgDeploySystemContracts) (*types.MsgDeploySystemContractsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// uniswap v2 factory
	factory, err := k.DeployUniswapV2Factory(ctx)
	if err != nil {
		return nil, cosmoserror.Wrapf(err, "failed to DeployUniswapV2Factory")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("UniswapV2Factory", factory.String()),
		),
	)

	// wzeta contract
	wzeta, err := k.DeployWZETA(ctx)
	if err != nil {
		return nil, cosmoserror.Wrapf(err, "failed to DeployWZetaContract")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployWZetaContract", wzeta.String()),
		),
	)

	// uniswap v2 router
	router, err := k.DeployUniswapV2Router02(ctx, factory, wzeta)
	if err != nil {
		return nil, cosmoserror.Wrapf(err, "failed to DeployUniswapV2Router02")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployUniswapV2Router02", router.String()),
		),
	)

	// connector zevm
	connector, err := k.DeployConnectorZEVM(ctx, wzeta)
	if err != nil {
		return nil, cosmoserror.Wrapf(err, "failed to DeployConnectorZEVM")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployConnectorZEVM", connector.String()),
		),
	)

	// system contract
	systemContract, err := k.DeploySystemContract(ctx, wzeta, factory, router)
	if err != nil {
		return nil, cosmoserror.Wrapf(err, "failed to SystemContractAddress")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("SystemContractAddress", systemContract.String()),
		),
	)

	return &types.MsgDeploySystemContractsResponse{
		UniswapV2Factory: factory.Hex(),
		Wzeta:            wzeta.Hex(),
		UniswapV2Router:  router.Hex(),
		ConnectorZEVM:    connector.Hex(),
		SystemContract:   systemContract.Hex(),
	}, nil
}
