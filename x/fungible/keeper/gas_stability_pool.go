package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// EnsureGasStabilityPoolAccountCreated ensures the gas stability pool account exists
func (k Keeper) EnsureGasStabilityPoolAccountCreated(ctx sdk.Context, chainID int64) {
	address := types.GasStabilityPoolAddress(chainID)

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
		return nil, err
	}

	address := types.GasStabilityPoolAddressEVM(chainID)
	return k.BalanceOfZRC4(ctx, gasZRC20, address)
}

// FundGasStabilityPool mints the ZRC20 into a special address called gas stability pool for the chain
func (k Keeper) FundGasStabilityPool(
	ctx sdk.Context,
	chainID int64,
	amount *big.Int,
) error {
	k.EnsureGasStabilityPoolAccountCreated(ctx, chainID)

	// get the gas zrc20 contract from the chain
	gasZRC20, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return err
	}

	// send to the gas stability pool address
	to := types.GasStabilityPoolAddressEVM(chainID)

	// call deposit ZRC20 method
	if err := k.CallZRC20Deposit(
		ctx,
		types.ModuleAddressEVM,
		gasZRC20,
		to,
		amount,
	); err != nil {
		return err
	}

	return nil
}

// WithdrawFromGasStabilityPool burns the ZRC20 from the gas stability pool
func (k Keeper) WithdrawFromGasStabilityPool(
	ctx sdk.Context,
	chainID int64,
	amount *big.Int,
) error {
	k.EnsureGasStabilityPoolAccountCreated(ctx, chainID)

	// get the gas zrc20 contract from the chain
	gasZRC20, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return err
	}

	// Ensure the account exists SDK level
	_ = k.GetAuthKeeper().GetAccount(ctx, types.GasStabilityPoolAddress(chainID))

	// Send from the gas stability pool address
	from := types.GasStabilityPoolAddressEVM(chainID)

	// call withdraw ZRC20 method
	if err := k.CallZRC20Burn(
		ctx,
		from,
		gasZRC20,
		amount,
		false,
	); err != nil {
		return err
	}

	return nil
}
