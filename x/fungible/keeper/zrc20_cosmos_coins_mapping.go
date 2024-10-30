package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/pkg/crypto"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

// LockZRC20 locks ZRC20 tokens in the specified address
// The caller must have approved the locker contract to spend the amount of ZRC20 tokens.
// Warning: This function does not mint cosmos coins, if the depositor needs to be rewarded
// it has to be implemented by the caller of this function.
func (k Keeper) LockZRC20(
	ctx sdk.Context,
	zrc20Address, spender, owner, locker common.Address,
	amount *big.Int,
) error {
	// owner is the EOA owner of the ZRC20 tokens.
	// spender is the EOA allowed to spend ZRC20 on owner's behalf.
	// locker is the address that will lock the ZRC20 tokens, i.e: bank precompile.
	if err := k.CheckZRC20Allowance(ctx, owner, spender, zrc20Address, amount); err != nil {
		return errors.Wrap(err, "failed allowance check")
	}

	// Check amount_to_be_locked <= total_erc20_balance - already_locked
	// Max amount of ZRC20 tokens that exists in zEVM are the total supply.
	totalSupply, err := k.ZRC20TotalSupply(ctx, zrc20Address)
	if err != nil {
		return errors.Wrap(err, "failed totalSupply check")
	}

	// The alreadyLocked amount is the amount of ZRC20 tokens that have been locked by the locker.
	// TODO: Implement list of whitelisted locker addresses (https://github.com/zeta-chain/node/issues/2991)
	alreadyLocked, err := k.ZRC20BalanceOf(ctx, zrc20Address, locker)
	if err != nil {
		return errors.Wrap(err, "failed getting the ZRC20 already locked amount")
	}

	if !k.IsValidDepositAmount(totalSupply, alreadyLocked, amount) {
		return errors.Wrap(fungibletypes.ErrInvalidAmount, "amount to be locked is not valid")
	}

	// Initiate a transferFrom the owner to the locker. This will lock the ZRC20 tokens.
	// locker has to initiate the transaction and have enough allowance from owner.
	transferred, err := k.ZRC20TransferFrom(ctx, zrc20Address, spender, owner, locker, amount)
	if err != nil {
		return errors.Wrap(err, "failed executing transferFrom")
	}

	if !transferred {
		return fmt.Errorf("transferFrom returned false (no success)")
	}

	return nil
}

// UnlockZRC20 unlocks ZRC20 tokens and sends them to the owner.
// Warning: Before unlocking ZRC20 tokens, the caller must check if
// the owner has enough collateral (cosmos coins) to be exchanged (burnt) for the ZRC20 tokens.
func (k Keeper) UnlockZRC20(
	ctx sdk.Context,
	zrc20Address, owner, locker common.Address,
	amount *big.Int,
) error {
	// Check if the account locking the ZRC20 tokens has enough balance.
	if err := k.CheckZRC20Balance(ctx, zrc20Address, locker, amount); err != nil {
		return errors.Wrap(err, "failed balance check")
	}

	// transfer from the EOA locking the assets to the owner.
	transferred, err := k.ZRC20Transfer(ctx, zrc20Address, locker, owner, amount)
	if err != nil {
		return errors.Wrap(err, "failed executing transfer")
	}

	if !transferred {
		return fmt.Errorf("transfer returned false (no success)")
	}

	return nil
}

// CheckZRC20Allowance checks if the allowance of ZRC20 tokens,
// is equal or greater than the provided amount.
func (k Keeper) CheckZRC20Allowance(
	ctx sdk.Context,
	owner, spender, zrc20Address common.Address,
	amount *big.Int,
) error {
	if amount.Sign() <= 0 || amount == nil {
		return fungibletypes.ErrInvalidAmount
	}

	if crypto.IsEmptyAddress(owner) || crypto.IsEmptyAddress(spender) {
		return fungibletypes.ErrZeroAddress
	}

	if err := k.IsValidZRC20(ctx, zrc20Address); err != nil {
		return errors.Wrap(err, "ZRC20 is not valid")
	}

	allowanceValue, err := k.ZRC20Allowance(ctx, zrc20Address, owner, spender)
	if err != nil {
		return errors.Wrap(err, "failed while checking spender's allowance")
	}

	if allowanceValue.Cmp(amount) < 0 || allowanceValue.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("invalid allowance, got %s, wanted %s", allowanceValue.String(), amount.String())
	}

	return nil
}

// CheckZRC20Balance checks if the balance of ZRC20 tokens,
// is equal or greater than the provided amount.
func (k Keeper) CheckZRC20Balance(
	ctx sdk.Context,
	zrc20Address, owner common.Address,
	amount *big.Int,
) error {
	if amount.Sign() <= 0 || amount == nil {
		return fungibletypes.ErrInvalidAmount
	}

	if err := k.IsValidZRC20(ctx, zrc20Address); err != nil {
		return errors.Wrap(err, "ZRC20 is not valid")
	}

	if crypto.IsEmptyAddress(owner) {
		return fungibletypes.ErrZeroAddress
	}

	// Check the ZRC20 balance of a given account.
	// function balanceOf(address account)
	balance, err := k.ZRC20BalanceOf(ctx, zrc20Address, owner)
	if err != nil {
		return errors.Wrap(err, "failed getting owner's ZRC20 balance")
	}

	if balance.Cmp(amount) < 0 {
		return fmt.Errorf("invalid balance, got %s, wanted %s", balance.String(), amount.String())
	}

	return nil
}

// IsValidZRC20 returns an error whenever a ZRC20 is not whitelisted or paused.
func (k Keeper) IsValidZRC20(ctx sdk.Context, zrc20Address common.Address) error {
	if crypto.IsEmptyAddress(zrc20Address) {
		return fungibletypes.ErrZRC20ZeroAddress
	}

	t, found := k.GetForeignCoins(ctx, zrc20Address.String())
	if !found {
		return fungibletypes.ErrZRC20NotWhiteListed
	}

	if t.Paused {
		return fungibletypes.ErrPausedZRC20
	}

	return nil
}

// IsValidDepositAmount checks "totalSupply >= amount_to_be_locked + amount_already_locked".
// A failure here means the user is trying to lock more than the available ZRC20 supply.
// This suggests that an actor is minting ZRC20 tokens out of thin air.
func (k Keeper) IsValidDepositAmount(totalSupply, alreadyLocked, amountToDeposit *big.Int) bool {
	if totalSupply == nil || alreadyLocked == nil || amountToDeposit == nil {
		return false
	}

	return totalSupply.Cmp(alreadyLocked.Add(alreadyLocked, amountToDeposit)) >= 0
}
