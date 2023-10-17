//go:build MOCK_MAINNET
// +build MOCK_MAINNET

package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
)

func (k Keeper) BlockOneDeploySystemContracts(goCtx context.Context) error {
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

	connector, err := k.DeployConnectorZEVM(ctx, wzeta)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeployConnectorZEVM")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployConnectorZEVM", connector.String()),
		),
	)
	ctx.Logger().Info("Deployed Connector ZEVM at " + connector.String())

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
	//err = k.SetGasPrice(ctx, big.NewInt(1337), big.NewInt(1))
	if err != nil {
		return err
	}
	_, err = k.SetupChainGasCoinAndPool(ctx, common.EthChain().ChainId, "ETH", "ETH", 18, nil)
	if err != nil {
		return errorsmod.Wrapf(err, fmt.Sprintf("failed to setupChainGasCoinAndPool for %s", common.EthChain().ChainName))
	}

	_, err = k.SetupChainGasCoinAndPool(ctx, common.BscMainnetChain().ChainId, "BNB", "BNB", 18, nil)
	if err != nil {
		return errorsmod.Wrapf(err, fmt.Sprintf("failed to setupChainGasCoinAndPool for %s", common.BscMainnetChain().ChainName))
	}
	_, err = k.SetupChainGasCoinAndPool(ctx, common.BtcMainnetChain().ChainId, "BTC", "BTC", 8, nil)
	if err != nil {
		return errorsmod.Wrapf(err, fmt.Sprintf("failed to setupChainGasCoinAndPool for %s", common.BtcMainnetChain().ChainName))
	}
	return nil
}
func (k Keeper) TestUpdateSystemContractAddress(_ context.Context) error {
	return nil
}
func (k Keeper) TestUpdateZRC20WithdrawFee(_ context.Context) error {
	return nil
}
