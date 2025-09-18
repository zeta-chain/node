package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/fungible/types"
)

// EnsureGasStabilityPoolAccountCreated ensures the gas stability pool account exists
func (k Keeper) EnsureGasStabilityPoolAccountCreated(ctx sdk.Context) {
	address := types.GasStabilityPoolAddress()

	ak := k.GetAuthKeeper()
	accExists := ak.HasAccount(ctx, address)
	if !accExists {
		ak.SetAccount(ctx, ak.NewAccountWithAddress(ctx, address))
	}
}

// GetGasStabilityPoolBalance returns the balance of the gas stability pool
func (k Keeper) GetGasStabilityPoolBalance(
	ctx sdk.Context,
	chainID int64,
) (*big.Int, error) {
	// get the gas zrc20 contract from the chain
	gasZRC20, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return nil, fmt.Errorf(
			"error getting gas zrc20 contract for chain ID %d: %s",
			chainID,
			err.Error(),
		)
	}

	return k.BalanceOfZRC4(ctx, gasZRC20, types.GasStabilityPoolAddressEVM())
}

// FundGasStabilityPool mints the ZRC20 into a special address called gas stability pool for the chain
func (k Keeper) FundGasStabilityPool(
	ctx sdk.Context,
	chainID int64,
	amount *big.Int,
) error {
	k.EnsureGasStabilityPoolAccountCreated(ctx)

	// get the gas zrc20 contract from the chain
	gasZRC20, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return err
	}

	// call deposit ZRC20 method
	return k.CallZRC20Deposit(
		ctx,
		types.ModuleAddressEVM,
		gasZRC20,
		types.GasStabilityPoolAddressEVM(),
		amount,
	)
}

// WithdrawFromGasStabilityPool burns the ZRC20 from the gas stability pool
func (k Keeper) WithdrawFromGasStabilityPool(
	ctx sdk.Context,
	chainID int64,
	amount *big.Int,
) error {
	k.EnsureGasStabilityPoolAccountCreated(ctx)

	// get the gas zrc20 contract from the chain
	gasZRC20, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return err
	}

	// call burn ZRC20 method
	return k.CallZRC20Burn(
		ctx,
		types.GasStabilityPoolAddressEVM(),
		gasZRC20,
		amount,
		false,
	)
}
