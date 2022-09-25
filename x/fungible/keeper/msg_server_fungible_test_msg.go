package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"
)

func (k msgServer) FungibleTestMsg(goCtx context.Context, msg *types.MsgFungibleTestMsg) (*types.MsgFungibleTestMsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ZDCAddree, err := k.DeployZetaDepositAndCall(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to DeployZetaDepositAndCall")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("action", "FungibleTestMsg"),
			sdk.NewAttribute("ZDCAddree", ZDCAddree.String()),
		),
	)

	gasPriceOracle, err := k.DeployGasPriceOracleContract(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to DeployZetaDepositAndCall")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("gasPriceOracle", gasPriceOracle.String()),
		),
	)

	system, _ := k.GetSystemContract(ctx)
	system.ZetaDepositAndCallContract = ZDCAddree.String()
	system.GasPriceOracleContract = gasPriceOracle.String()
	k.SetSystemContract(ctx, system)

	transferGasLimit := big.NewInt(21_000)
	addr, err := k.DeployZRC4Contract(ctx, "ETH", "zETH", 18, "GOERLI", common.CoinType_Gas, "", transferGasLimit)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("ETH_GOERLI_ZRC4", addr.String()),
		),
	)

	addr, err = k.DeployZRC4Contract(ctx, "BNB", "zBNB", 18, "BSCTESTNET", common.CoinType_Gas, "", transferGasLimit)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("BNB_BSCTESTNET_ZRC4", addr.String()),
		),
	)

	addr, err = k.DeployZRC4Contract(ctx, "MATIC", "zMATIC", 18, "MUMBAI", common.CoinType_Gas, "", transferGasLimit)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("MATIC_BSCTESTNET_ZRC4", addr.String()),
		),
	)

	return &types.MsgFungibleTestMsgResponse{}, nil
}
