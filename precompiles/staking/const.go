package staking

const (
	// State changing methods.
	ClaimRewardsMethodName = "claimRewards"
	ClaimRewardsEventName  = "ClaimedRewards"
	ClaimRewardsMethodGas  = 10000

	DistributeMethodName = "distribute"
	DistributeEventName  = "Distributed"
	DistributeMethodGas  = 10000

	StakeMethodName = "stake"
	StakeEventName  = "Stake"
	StakeMethodGas  = 10000

	UnstakeMethodName = "unstake"
	UnstakeEventName  = "Unstake"
	UnstakeMethodGas  = 1000

	MoveStakeMethodName = "moveStake"
	MoveStakeEventName  = "MoveStake"
	MoveStakeMethodGas  = 10000

	// Query methods.
	GetAllValidatorsMethodName = "getAllValidators"
	GetSharesMethodName        = "getShares"
	GetRewardsMethodName       = "getRewards"
	GetValidatorsMethodName    = "getDelegatorValidators"
)
