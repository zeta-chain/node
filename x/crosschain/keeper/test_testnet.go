//go:build !PRIVNET
// +build !PRIVNET

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) TestWhitelistERC20(ctx sdk.Context) error {
	return nil
}
