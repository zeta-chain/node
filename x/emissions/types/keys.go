package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	// ModuleName defines the module name
	ModuleName                       = "emissions"
	UndistributedObserverRewardsPool = ModuleName + "Observers"
	UndistributedTssRewardsPool      = ModuleName + "Tss"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey              = "mem_emissions"
	WithdrawableEmissionsKey = "WithdrawableEmissions-value-"

	SecsInMonth        = 30 * 24 * 60 * 60
	BlockRewardsInZeta = "210000000"

	EmissionScheduledYears = 4
	AvgBlockTime           = "5.7"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	EmissionsTrackerKey              = "EmissionsTracker-value-"
	ParamMaxBondFactor               = "MaxBondFactor"
	ParamMinBondFactor               = "MinBondFactor"
	ParamAvgBlockTime                = "AvgBlockTime"
	ParamTargetBondRatio             = "TargetBondRation"
	ParamValidatorEmissionPercentage = "ValidatorEmissionPercentage"
	ParamObserverEmissionPercentage  = "ObserverEmissionPercentage"
	ParamTssSignerEmissionPercentage = "SignerEmissionPercentage"
	ParamDurationFactorConstant      = "DurationFactorConstant"
)

var (
	EmissionsModuleAddress                  = authtypes.NewModuleAddress(ModuleName)
	UndistributedObserverRewardsPoolAddress = authtypes.NewModuleAddress(UndistributedObserverRewardsPool)
	UndistributedTssRewardsPoolAddress      = authtypes.NewModuleAddress(UndistributedTssRewardsPool)
)
