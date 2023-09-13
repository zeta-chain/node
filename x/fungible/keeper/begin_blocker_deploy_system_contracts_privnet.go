//go:build PRIVNET

package keeper

import (
	"context"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// This is for privnet/testnet only
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
	// FIXME: remove unnecessary SetGasPrice and setupChainGasCoinAndPool
	k.SetSystemContract(ctx, system)
	//err = k.SetGasPrice(ctx, big.NewInt(1337), big.NewInt(1))
	if err != nil {
		return err
	}

	_, err = k.SetupChainGasCoinAndPool(ctx, common.GoerliChain().ChainId, "ETH", "gETH", 18)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to setupChainGasCoinAndPool")
	}

	_, err = k.SetupChainGasCoinAndPool(ctx, common.BtcRegtestChain().ChainId, "BTC", "tBTC", 8)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to setupChainGasCoinAndPool")
	}

	//FIXME: clean up and config the above based on localnet/testnet/mainnet

	// for localnet only: USDT ZRC20
	USDTAddr := "0xff3135df4F2775f4091b81f4c7B6359CfA07862a"
	USDTZRC20Addr, err := k.DeployZRC20Contract(ctx, "USDT", "USDT", uint8(6), common.GoerliChain().ChainId, common.CoinType_ERC20, USDTAddr, big.NewInt(90_000))
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeployZRC20Contract USDT")
	}
	ctx.Logger().Info("Deployed USDT ZRC20 at " + USDTZRC20Addr.String())
	// for localnet only: ZEVM Swap App
	//ZEVMSwapAddress, err := k.DeployZEVMSwapApp(ctx, router, SystemContractAddress)
	//if err != nil {
	//	return sdkerrors.Wrapf(err, "failed to deploy ZEVMSwapApp")
	//}
	//ctx.Logger().Info("Deployed ZEVM Swap App at " + ZEVMSwapAddress.String())
	fmt.Println("Successfully deployed contracts")
	return nil
}

func (k Keeper) TestUpdateSystemContractAddress(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	wzeta, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to GetWZetaContractAddress")
	}
	uniswapV2Factory, err := k.GetUniswapV2FactoryAddress(ctx)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to GetUniswapv2FacotryAddress")
	}
	router, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to GetUniswapV2Router02Address")
	}

	SystemContractAddress, err := k.DeploySystemContract(ctx, wzeta, uniswapV2Factory, router)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to DeploySystemContract")
	}
	creator := k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_deploy_fungible_coin)
	msg := types.NewMsgUpdateSystemContract(creator, SystemContractAddress.Hex())
	_, err = k.UpdateSystemContract(ctx, msg)
	return err
}

func (k Keeper) TestUpdateZRC20WithdrawFee(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	foreignCoins := k.GetAllForeignCoins(ctx)
	creator := k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_deploy_fungible_coin)

	for _, foreignCoin := range foreignCoins {
		msg := types.NewMsgUpdateZRC20WithdrawFee(creator, foreignCoin.Zrc20ContractAddress, sdk.NewUint(uint64(foreignCoin.ForeignChainId)))
		_, err := k.UpdateZRC20WithdrawFee(ctx, msg)
		if err != nil {
			return err
		}
	}

	return nil
}
