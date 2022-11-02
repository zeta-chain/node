package keeper

import (
	"fmt"
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

func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec)                {}
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

// CleanObservers cleans a observer Mapper checking delegation amount for a speficific delagator
func (k Keeper) CheckAndCleanObserverDelegator(ctx sdk.Context, valAddress sdk.ValAddress, delAddress sdk.AccAddress) error {
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddress.String())
	if err != nil {
		return err
	}
	fmt.Println(accAddress.String(), delAddress.String())
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
		err := k.CheckObserverDelegation(ctx, accAddress.String(), mapper.ObserverChain, mapper.ObservationType)
		if err != nil {
			mapper.ObserverList = CleanAddressList(mapper.ObserverList, accAddress.String())
			k.SetObserverMapper(ctx, mapper)
		}
	}
}
