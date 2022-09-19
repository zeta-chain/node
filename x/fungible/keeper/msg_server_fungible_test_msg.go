package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	contracts "github.com/zeta-chain/zetacore/contracts/evm"
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
			sdk.NewAttribute("action", "DeployZetaDepositAndCall"),
			sdk.NewAttribute("ZDCAddree", ZDCAddree.String()),
		),
	)

	// TODO: Handling the message
	addr, err := k.DeployZRC4Contract(ctx, "ETH", "zETH", 18, "GOERLI", common.CoinType_Gas, "")
	if err != nil {
		return nil, err
	}

	alice := ethcommon.HexToAddress("0x9B4547Dd93e93c526a11c5123Ab74D42aFb41B10")
	bal1 := k.BalanceOfZRC4(ctx, addr, alice)
	if bal1 == nil {
		return nil, sdkerrors.Wrap(types.ErrBlanceQuery, fmt.Sprintf("zrc4 balance of %s", alice.String()))
	}
	_, err = k.DepositZRC4(ctx, addr, alice, big.NewInt(1000))
	if err != nil {
		return nil, err
	}
	bal2 := k.BalanceOfZRC4(ctx, addr, alice)
	if bal2 == nil {
		return nil, sdkerrors.Wrap(types.ErrBlanceQuery, fmt.Sprintf("zrc4 balance of %s", alice.String()))
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("action", "FungibleTestMsg"),
			sdk.NewAttribute("chain", "GOERLI"),
			sdk.NewAttribute("contract", addr.String()),
			sdk.NewAttribute("balance1", bal1.String()),
			sdk.NewAttribute("balance2", bal2.String()),
		),
	)

	addr, err = k.DeployZRC4Contract(ctx, "BNB", "zBNB", 18, "BSCTESTNET", common.CoinType_Gas, "")
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("chain", "BSCTESTNET"),
			sdk.NewAttribute("contract", addr.String()),
		),
	)
	_, err = k.DepositZRC4(ctx, addr, alice, big.NewInt(1000))

	// test withdraw
	abi, err := contracts.ZRC4MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	_, err = k.CallEVM(ctx, *abi, alice, addr, true, "withdraw", alice.Bytes(), big.NewInt(17))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("action", "withdraw"),
		),
	)

	return &types.MsgFungibleTestMsgResponse{}, nil
}
