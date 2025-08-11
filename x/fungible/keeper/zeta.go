package keeper

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/x/fungible/types"
)

// ZETAMaxSupplyStr is the maximum mintable ZETA in the fungible module
// 1.85 billion ZETA
const ZETAMaxSupplyStr = "1850000000000000000000000000"

// MintZetaToEVMAccount mints ZETA (gas token) to the given address
// NOTE: this method should be used with a temporary context, and it should not be committed if the method returns an error
func (k *Keeper) MintZetaToEVMAccount(ctx sdk.Context, to sdk.AccAddress, amount *big.Int) error {
	if err := k.validateZetaSupply(ctx, amount); err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewIntFromBigInt(amount)))
	// Mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	// Send minted coins to the receiver
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, coins)
}

// MintZetaToFungibleModule mints ZETA (gas token) to the fungible module account
// Fungible module account is the protocol address used in the smart contracts
func (k *Keeper) MintZetaToFungibleModule(ctx sdk.Context, amount *big.Int) error {
	if err := k.validateZetaSupply(ctx, amount); err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewIntFromBigInt(amount)))
	// Mint coins
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
}

// validateZetaSupply checks if the minted ZETA amount exceeds the maximum supply
func (k *Keeper) validateZetaSupply(ctx sdk.Context, amount *big.Int) error {
	zetaMaxSupply, ok := sdkmath.NewIntFromString(ZETAMaxSupplyStr)
	if !ok {
		return fmt.Errorf("failed to parse ZETA max supply: %s", ZETAMaxSupplyStr)
	}

	supply := k.bankKeeper.GetSupply(ctx, config.BaseDenom)
	if supply.Amount.Add(sdkmath.NewIntFromBigInt(amount)).GT(zetaMaxSupply) {
		return types.ErrMaxSupplyReached
	}
	return nil
}

// ExecuteWithMintedZeta is a helper function that mints ZETA to the fungible module account
// and executes the provided operation within a temporary context.
// If the operation is successful, it commits the temporary context.
// If it is not successful, it rolls back the temporary context thus preventing surplus ZETA from being minted.
func (k Keeper) ExecuteWithMintedZeta(
	ctx sdk.Context,
	amount *big.Int,
	operation func(sdk.Context) (*evmtypes.MsgEthereumTxResponse, bool, error),
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	tmpCtx, commit := ctx.CacheContext()

	if err := k.MintZetaToFungibleModule(tmpCtx, amount); err != nil {
		return nil, false, err
	}

	res, isCrossChain, err := operation(tmpCtx)
	if err == nil {
		commit()
	}

	return res, isCrossChain, err
}
