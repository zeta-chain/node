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
	mappers := k.GetAllObserverMappersForAddress(ctx, accAddress.String())
	if len(mappers) == 0 {
		return nil
	}

	tokensToBurn := sdk.NewDecFromInt(validator.Tokens).Mul(fraction)
	resultingTokens := validator.Tokens.Sub(tokensToBurn.Ceil().TruncateInt())
	for _, mapper := range mappers {

		cp, found := k.GetCoreParamsByChainID(ctx, mapper.ObserverChain.ChainId)
		if !found || cp == nil || !cp.IsSupported {
			return types.ErrSupportedChains
		}

		if sdk.NewDecFromInt(resultingTokens).LT(cp.MinObserverDelegation) {
			mapper.ObserverList = CleanAddressList(mapper.ObserverList, accAddress.String())
			k.SetObserverMapper(ctx, mapper)
		}
	}
	return nil
}

// CleanObservers cleans a observer Mapper without checking delegation amount
func (k Keeper) CleanObservers(ctx sdk.Context, valAddress sdk.ValAddress) error {
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	mappers := k.GetAllObserverMappersForAddress(ctx, accAddress.String())
	for _, mapper := range mappers {
		mapper.ObserverList = CleanAddressList(mapper.ObserverList, accAddress.String())
		k.SetObserverMapper(ctx, mapper)
	}
	return nil
}

// CleanObservers cleans a observer Mapper checking delegation amount
func (k Keeper) CheckAndCleanObserver(ctx sdk.Context, valAddress sdk.ValAddress) error {
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	k.CleanMapper(ctx, accAddress)
	return nil
}

// CleanObservers cleans a observer Mapper checking delegation amount for a speficific delagator. It is used when delgator is the validator .
// That is when when the validator is trying to remove self delgation
func (k Keeper) CheckAndCleanObserverDelegator(ctx sdk.Context, valAddress sdk.ValAddress, delAddress sdk.AccAddress) error {
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	if !(accAddress.String() == delAddress.String()) {
		return nil
	}
	k.CleanMapper(ctx, accAddress)
	return nil
}

func CleanAddressList(addresslist []string, address string) []string {
	index := -1
	for i, addr := range addresslist {
		if addr == address {
			index = i
		}
	}
	if index != -1 {
		addresslist = RemoveIndex(addresslist, index)
	}
	return addresslist
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func (k Keeper) CleanMapper(ctx sdk.Context, accAddress sdk.AccAddress) {
	mappers := k.GetAllObserverMappersForAddress(ctx, accAddress.String())
	for _, mapper := range mappers {
		err := k.CheckObserverDelegation(ctx, accAddress.String(), mapper.ObserverChain)
		if err != nil {
			mapper.ObserverList = CleanAddressList(mapper.ObserverList, accAddress.String())
			k.SetObserverMapper(ctx, mapper)
		}
	}
}
