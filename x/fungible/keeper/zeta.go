package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// MintZetaToEVMAccount mints ZETA (gas token) to the given address
func (k *Keeper) MintZetaToEVMAccount(ctx sdk.Context, to sdk.AccAddress, amount *big.Int) error {
	balanceCoin := k.bankKeeper.GetBalance(ctx, to, config.BaseDenom)
	coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewIntFromBigInt(amount)))
	// Mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	// Send minted coins to the receiver
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, coins)

	if err == nil {
		// Check expected receiver balance after transfer
		balanceCoinAfter := k.bankKeeper.GetBalance(ctx, to, config.BaseDenom)
		expCoin := balanceCoin.Add(coins[0])

		if ok := balanceCoinAfter.IsEqual(expCoin); !ok {
			err = sdkerrors.Wrapf(
				types.ErrBalanceInvariance,
				"invalid coin balance - expected: %v, actual: %v",
				expCoin, balanceCoinAfter,
			)
		}
	}

	if err != nil {
		// Revert minting if an error is found.
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
			return err
		}
		return err
	}

	return nil
}
