package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

var _ types.StakingHooks = Hooks{}

type Hooks struct {
	k Keeper
}

func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {
	err := h.k.CleanObservers(ctx, valAddr)
	if err != nil {
		panic(err)
	}
}

func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {
	err := h.k.CheckAndCleanObserver(ctx, valAddr)
	if err != nil {
		panic(err)
	}
}

func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	err := h.k.CheckAndCleanObserverDelegator(ctx, valAddr, delAddr)
	if err != nil {
		panic(err)
	}
}

func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	err := h.k.CleanSlashedValidator(ctx, valAddr, fraction)
	if err != nil {
		panic(err)
	}
}
func (h Hooks) AfterValidatorCreated(_ sdk.Context, _ sdk.ValAddress)                            {}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                          {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)          {}
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}

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
	tokensToBurn := validator.Tokens.ToDec().Mul(fraction)
	resultingTokens := validator.Tokens.Sub(tokensToBurn.Ceil().TruncateInt())
	for _, mapper := range mappers {
		obsParams, supported := k.GetParams(ctx).GetParamsForChainAndType(mapper.ObserverChain, mapper.ObservationType)
		if !supported {
			return types.ErrSupportedChains
		}
		if resultingTokens.ToDec().LT(obsParams.MinObserverDelegation) {
			mapper.ObserverList = CleanAddressList(mapper.ObserverList, accAddress.String())
			k.SetObserverMapper(ctx, mapper)
			RemoveObserverEvent(ctx, *mapper, accAddress.String(), "validator slashed below minimum observer delegation")
		}
	}
	return nil
}

// CleanObservers cleans an observer Mapper without checking delegation amount
func (k Keeper) CleanObservers(ctx sdk.Context, valAddress sdk.ValAddress) error {
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	mappers := k.GetAllObserverMappersForAddress(ctx, accAddress.String())
	for _, mapper := range mappers {
		mapper.ObserverList = CleanAddressList(mapper.ObserverList, accAddress.String())
		k.SetObserverMapper(ctx, mapper)
		RemoveObserverEvent(ctx, *mapper, accAddress.String(), "validator unbonded")
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

// CleanObservers cleans a observer Mapper checking delegation amount for a speficific delagator
func (k Keeper) CheckAndCleanObserverDelegator(ctx sdk.Context, valAddress sdk.ValAddress, delAddress sdk.AccAddress) error {
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	// verify this is a self delegation
	if !(accAddress.String() == delAddress.String()) {
		return nil
	}
	k.CleanMapper(ctx, accAddress)
	return nil
}

func (k Keeper) CleanMapper(ctx sdk.Context, accAddress sdk.AccAddress) {
	mappers := k.GetAllObserverMappersForAddress(ctx, accAddress.String())
	for _, mapper := range mappers {
		err := k.CheckObserverDelegation(ctx, accAddress.String(), mapper.ObserverChain, mapper.ObservationType)
		if err != nil {
			mapper.ObserverList = CleanAddressList(mapper.ObserverList, accAddress.String())
			k.SetObserverMapper(ctx, mapper)
			RemoveObserverEvent(ctx, *mapper, accAddress.String(), "validators self delegation is below minimum observer delegation")
		}
	}
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
