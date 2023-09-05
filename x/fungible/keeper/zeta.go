package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// MintZetaToEVMAccount mints ZETA (gas token) to the given address
// NOTE: this method should be used with a temporary context, and it should not be committed if the method returns an erorr
func (k *Keeper) MintZetaToEVMAccount(ctx sdk.Context, to sdk.AccAddress, amount *big.Int) error {
	coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewIntFromBigInt(amount)))
	// Mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	// Send minted coins to the receiver
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, coins)
}
