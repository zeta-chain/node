package types

import (
	"math"

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

// BlockTimeSeconds is the estimated block time in seconds
// TargetBallotMaturitySeconds is the real-world target duration for ballot maturity
const (
	BlockTimeSeconds            = 4.5
	TargetBallotMaturitySeconds = 600
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
	// TODO: how is this value calculated?
	BlockReward = sdkmath.LegacyMustNewDecFromStr("9620949074074074074.074070733466756687")
	// ObserverSlashAmount is the amount of tokens to be slashed from observer in case of incorrect vote
	// by default it is set to 0.1 ZETA
	ObserverSlashAmount = sdkmath.NewInt(100000000000000000)

	// BallotMaturityBlocks is amount of blocks needed for ballot to mature
	BallotMaturityBlocks = int64(
		math.Round(TargetBallotMaturitySeconds / BlockTimeSeconds),
	) // approximately 9-10 minutes
	// PendingBallotsBufferBlocks is a buffer number of blocks
	//(in addition to BallotMaturityBlocks)
	// that we use only for pending ballots before deleting them
	PendingBallotsBufferBlocks = int64(144000) // 10 days(60 * 60 * 24 * 10)
)
