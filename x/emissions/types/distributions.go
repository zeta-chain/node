package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetRewardsDistributions returns the current distribution of rewards
// for validators, observers and TSS signers
// If the percentage is not set, it returns 0
func GetRewardsDistributions(params Params) (sdkmath.Int, sdkmath.Int, sdkmath.Int) {
	// Fetch the validator rewards, use 0 if the percentage is not set
	validatorRewards := sdk.NewInt(0)
	validatorRewardsDec, err := sdk.NewDecFromStr(params.ValidatorEmissionPercentage)
	if err == nil {
		validatorRewards = validatorRewardsDec.Mul(BlockReward).TruncateInt()
	}

	// Fetch the observer rewards, use 0 if the percentage is not set
	observerRewards := sdk.NewInt(0)
	observerRewardsDec, err := sdk.NewDecFromStr(params.ObserverEmissionPercentage)
	if err == nil {
		observerRewards = observerRewardsDec.Mul(BlockReward).TruncateInt()
	}

	// Fetch the TSS signer rewards, use 0 if the percentage is not set
	tssSignerRewards := sdk.NewInt(0)
	tssSignerRewardsDec, err := sdk.NewDecFromStr(params.TssSignerEmissionPercentage)
	if err == nil {
		tssSignerRewards = tssSignerRewardsDec.Mul(BlockReward).TruncateInt()
	}

	return validatorRewards, observerRewards, tssSignerRewards
}
