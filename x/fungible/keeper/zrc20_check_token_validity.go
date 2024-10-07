package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// IsValidZRC20 returns an error whenever a ZRC20 is not whitelisted or paused.
func (k Keeper) IsValidZRC20(ctx sdk.Context, zrc20Address common.Address) error {
	if zrc20Address == zeroAddress {
		return fmt.Errorf("zrc20 address cannot be zero")
	}

	t, found := k.GetForeignCoins(ctx, zrc20Address.String())
	if !found {
		return fmt.Errorf("ZRC20 is not whitelisted, address: %s", zrc20Address.String())
	}

	if t.Paused {
		return fmt.Errorf("ZRC20 is paused, address: %s", zrc20Address.String())
	}

	return nil
}
