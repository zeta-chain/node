package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

// BurnFungibleModuleAsset burns the zrc20 balance on the fungible module
func (k msgServer) BurnFungibleModuleAsset(
	goCtx context.Context,
	msg *types.MsgBurnFungibleModuleAsset,
) (*types.MsgBurnFungibleModuleAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check signer permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// check the zrc20 exists
	zrc20Addr := ethcommon.HexToAddress(msg.Zrc20Address)
	if zrc20Addr == (ethcommon.Address{}) {
		return nil, cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid zrc20 contract address (%s)",
			msg.Zrc20Address,
		)
	}
	_, found := k.GetForeignCoins(ctx, msg.Zrc20Address)
	if !found {
		return nil, cosmoserrors.Wrapf(
			types.ErrForeignCoinNotFound,
			"no foreign coin match requested zrc20 address (%s)",
			msg.Zrc20Address,
		)
	}

	// get the balance of the fungible module
	balance, err := k.ZRC20BalanceOf(ctx, zrc20Addr, types.ModuleAddressEVM)
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to query zrc20 balance (%s)", err.Error())
	}
	if balance.Uint64() == 0 {
		return nil, cosmoserrors.Wrapf(
			types.ErrForeignCoinNotFound,
			"no balance found for zrc20 address (%s)",
			msg.Zrc20Address,
		)
	}

	// burn the zrc20 balance
	if err := k.CallZRC20Burn(
		ctx,
		types.ModuleAddressEVM,
		zrc20Addr,
		balance,
		false,
	); err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to burn zrc20 balance (%s)", err.Error())
	}

	return &types.MsgBurnFungibleModuleAssetResponse{}, nil
}
