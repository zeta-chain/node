package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

var _ types.StakingHooks = Hooks{}

type Hooks struct {
	k Keeper
}

func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	err := h.k.CleanObservers(ctx, valAddr)
	if err != nil {
		return err
	}
	return nil
}

func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	err := h.k.CheckAndCleanObserver(ctx, valAddr)
	if err != nil {
		return err
	}
	return nil
}

func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	err := h.k.CheckAndCleanObserverDelegator(ctx, valAddr, delAddr)
	if err != nil {
		return err
	}
	return nil
}

func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {

	err := h.k.CleanSlashedValidator(ctx, valAddr, fraction)
	if err != nil {
		return err
	}
	return nil
}
func (h Hooks) AfterValidatorCreated(_ sdk.Context, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (k Keeper) CleanSlashedValidator(ctx sdk.Context, valAddress sdk.ValAddress, fraction sdk.Dec) error {
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return types.ErrNotValidator
	}
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	observerSet, found := k.GetObserverSet(ctx)
	if !found || observerSet.Len() == 0 {
		return nil
	}

	tokensToBurn := sdk.NewDecFromInt(validator.Tokens).Mul(fraction)
	resultingTokens := validator.Tokens.Sub(tokensToBurn.Ceil().TruncateInt())

	mindelegation, found := types.GetMinObserverDelegation()
	if !found {
		return types.ErrMinDelegationNotFound
	}
	if resultingTokens.LT(mindelegation) {
		k.RemoveObserverFromSet(ctx, accAddress.String())
	}
	return nil
}

// CleanObservers cleans a observer Mapper without checking delegation amount
func (k Keeper) CleanObservers(ctx sdk.Context, valAddress sdk.ValAddress) error {
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	k.RemoveObserverFromSet(ctx, accAddress.String())
	return nil
}

// CheckAndCleanObserver checks if the observer self-delegation is sufficient,
// if not it removes the observer from the set
func (k Keeper) CheckAndCleanObserver(ctx sdk.Context, valAddress sdk.ValAddress) error {
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	err = k.CheckObserverSelfDelegation(ctx, accAddress.String())
	if err != nil {
		return err
	}
	return nil
}

// CheckAndCleanObserverDelegator first checks if the delegation is self delegation,
// if it is, then it checks if the total delegation is sufficient after the delegation is removed,
// if not it removes the observer from the set
func (k Keeper) CheckAndCleanObserverDelegator(ctx sdk.Context, valAddress sdk.ValAddress, delAddress sdk.AccAddress) error {
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	// Check if this is a self delegation, if it's not then return, we only check self-delegation for cleaning observer set
	if !(accAddress.String() == delAddress.String()) {
		return nil
	}
	err = k.CheckObserverSelfDelegation(ctx, accAddress.String())
	if err != nil {
		return err
	}
	return nil
}
