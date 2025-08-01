package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/pkg/constant"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

// BurnFungibleModuleAsset burns the zrc20 balance on the fungible module
// If the zero address is provided, it burns the native ZETA held from the fungible module
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

	// if the zero address is provided, burn the ZETA from the fungible module
	if msg.Zrc20Address == constant.EVMZeroAddress {
		// get the balance of the fungible module
		balance := k.bankKeeper.SpendableCoin(ctx, types.ModuleAddress, config.BaseDenom)
		if balance.IsZero() {
			return nil, cosmoserrors.Wrapf(
				types.ErrZeroBalance,
				"no balance found for fungible module",
			)
		}
		if err := k.bankKeeper.BurnCoins(
			ctx,
			types.ModuleName,
			sdk.NewCoins(sdk.NewCoin(config.BaseDenom, balance.Amount)),
		); err != nil {
			return nil, cosmoserrors.Wrapf(types.ErrFailedToBurn, "failed to burn zeta balance (%s)", err.Error())
		}
	} else {
		// process the zrc20 burn
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
		if balance == nil {
			return nil, cosmoserrors.Wrapf(
				types.ErrZeroBalance,
				"balance is nil for zrc20 address (%s)",
				msg.Zrc20Address,
			)
		}
		if balance.Sign() == 0 {
			return nil, cosmoserrors.Wrapf(
				types.ErrZeroBalance,
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
			return nil, cosmoserrors.Wrapf(types.ErrFailedToBurn, "failed to burn zrc20 balance (%s)", err.Error())
		}
	}

	return &types.MsgBurnFungibleModuleAssetResponse{}, nil
}
