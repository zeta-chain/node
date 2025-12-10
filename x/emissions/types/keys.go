package types

import (
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	// ModuleName defines the module name
	ModuleName                       = "emissions"
	UndistributedObserverRewardsPool = ModuleName + "Observers"
	UndistributedTSSRewardsPool      = ModuleName + "Tss"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey              = "mem_emissions"
	WithdrawableEmissionsKey = "WithdrawableEmissions-value-"
	ParamsKey                = "Params-value-"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	EmissionsTrackerKey              = "EmissionsTracker-value-"
	ParamValidatorEmissionPercentage = "ValidatorEmissionPercentage"
	ParamObserverEmissionPercentage  = "ObserverEmissionPercentage"
	ParamTssSignerEmissionPercentage = "SignerEmissionPercentage"
	ParamObserverSlashAmount         = "ObserverSlashAmount"
)

var (
	EmissionsModuleAddress                  = authtypes.NewModuleAddress(ModuleName)
	UndistributedObserverRewardsPoolAddress = authtypes.NewModuleAddress(UndistributedObserverRewardsPool)
	UndistributedTssRewardsPoolAddress      = authtypes.NewModuleAddress(UndistributedTSSRewardsPool)
	// BlockReward is an initial block reward amount when emissions module was initialized.
	// The current value can be obtained from by querying the params
	BlockReward = sdkmath.LegacyMustNewDecFromStr("3375771604938271604.938271604938271605")
	// ObserverSlashAmount is the amount of tokens to be slashed from observer in case of incorrect vote
	// by default it is set to 0.1 ZETA
	ObserverSlashAmount = sdkmath.NewInt(100000000000000000)

	// BallotMaturityBlocks is amount of blocks needed for ballot to mature
	// by default is set to 300
	BallotMaturityBlocks = 300 // approximately 9-10 minutes
	// PendingBallotsBufferBlocks is a buffer number of blocks
	//(in addition to BallotMaturityBlocks)
	// that we use only for pending ballots before deleting them
	PendingBallotsBufferBlocks = int64(432000) // 10 days(60 * 60 * 24 * 10)
)
