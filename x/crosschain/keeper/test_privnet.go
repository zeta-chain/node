//go:build PRIVNET
// +build PRIVNET

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) TestWhitelistERC20(ctx sdk.Context) error {
	goCtx := sdk.UnwrapSDKContext(ctx)
	creator := k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_group1)
	msg := types.NewMsgWhitelistERC20(creator, types.ModuleAddressEVM.Hex(), common.GoerliChain().ChainId, "test", "testerc20", 17, 90_000)

	_, err := k.WhitelistERC20(goCtx, msg)
	if err != nil {
		panic(err)
	}
	return nil
}
