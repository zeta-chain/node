package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) BlockZeroDeploySystemContracts(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// setup uniswap v2 factory
	uniswapV2Factory, err := k.DeployUniswapV2Factory(ctx)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeployUniswapV2Factory")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("UniswapV2Factory", uniswapV2Factory.String()),
		),
	)

	// setup WZETA contract
	wzeta, err := k.DeployWZETA(ctx)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeployWZetaContract")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployWZetaContract", wzeta.String()),
		),
	)

	router, err := k.DeployUniswapV2Router02(ctx, uniswapV2Factory, wzeta)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeployUniswapV2Router02")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployUniswapV2Router02", router.String()),
		),
	)

	SystemContractAddress, err := k.DeploySystemContract(ctx, wzeta, uniswapV2Factory, router)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to SystemContractAddress")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("SystemContractAddress", SystemContractAddress.String()),
		),
	)

	// set the system contract
	system, _ := k.GetSystemContract(ctx)
	system.SystemContract = SystemContractAddress.String()
	k.SetSystemContract(ctx, system)
	return nil
}
