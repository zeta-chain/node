package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

var _ types.StakingHooks = Hooks{}

type Hooks struct {
	k Keeper
}

func (h Hooks) AfterValidatorRemoved(c context.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	ctx := sdk.UnwrapSDKContext(c)
	err := h.k.CleanObservers(ctx, valAddr)
	if err != nil {
		ctx.Logger().Error("Error cleaning observer set", "error", err)
	}
	return nil
}

func (h Hooks) AfterValidatorBeginUnbonding(c context.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	ctx := sdk.UnwrapSDKContext(c)
	err := h.k.CheckAndCleanObserver(ctx, valAddr)
	if err != nil {
		ctx.Logger().Error("Error cleaning observer set", "error", err)
	}
	return nil
}

func (h Hooks) AfterDelegationModified(c context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	ctx := sdk.UnwrapSDKContext(c)
	err := h.k.CheckAndCleanObserverDelegator(ctx, valAddr, delAddr)
	if err != nil {
		ctx.Logger().Error("Error cleaning observer set", "error", err)
	}
	return nil
}

func (h Hooks) BeforeValidatorSlashed(c context.Context, valAddr sdk.ValAddress, fraction sdkmath.LegacyDec) error {
	ctx := sdk.UnwrapSDKContext(c)
	err := h.k.CleanSlashedValidator(ctx, valAddr, fraction)
	if err != nil {
		ctx.Logger().Error("Error cleaning observer set", "error", err)
	}
	return nil
}
func (h Hooks) AfterValidatorCreated(_ context.Context, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeValidatorModified(_ context.Context, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorBonded(_ context.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationCreated(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationSharesModified(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationRemoved(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (k Keeper) CleanSlashedValidator(
	ctx context.Context,
	valAddress sdk.ValAddress,
	fraction sdkmath.LegacyDec,
) error {
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	observerSet, found := k.GetObserverSet(sdkCtx)
	if !found || observerSet.Len() == 0 {
		return nil
	}

	tokensToBurn := sdkmath.LegacyNewDecFromInt(validator.Tokens).Mul(fraction)
	resultingTokens := validator.Tokens.Sub(tokensToBurn.Ceil().TruncateInt())

	mindelegation, found := types.GetMinObserverDelegation()
	if !found {
		return types.ErrMinDelegationNotFound
	}
	if resultingTokens.LT(mindelegation) {
		k.RemoveObserverFromSet(sdkCtx, accAddress.String())
	}
	return nil
}

// CleanObservers cleans a observer Mapper without checking delegation amount
func (k Keeper) CleanObservers(ctx context.Context, valAddress sdk.ValAddress) error {
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.RemoveObserverFromSet(sdkCtx, accAddress.String())
	return nil
}

// CheckAndCleanObserver checks if the observer self-delegation is sufficient,
// if not it removes the observer from the set
func (k Keeper) CheckAndCleanObserver(c context.Context, valAddress sdk.ValAddress) error {
	ctx := sdk.UnwrapSDKContext(c)
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
func (k Keeper) CheckAndCleanObserverDelegator(
	c context.Context,
	valAddress sdk.ValAddress,
	delAddress sdk.AccAddress,
) error {
	ctx := sdk.UnwrapSDKContext(c)
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}

	// Check if this is a self delegation, if it's not then return, we only check self-delegation for cleaning observer set
	if accAddress.String() != delAddress.String() {
		return nil
	}

	err = k.CheckObserverSelfDelegation(ctx, accAddress.String())
	if err != nil {
		return err
	}
	return nil
}

func (h Hooks) AfterUnbondingInitiated(_ context.Context, _ uint64) error {
	return nil
}
