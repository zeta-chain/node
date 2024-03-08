package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

// GetDistributions returns the current distribution of rewards
// for validators, observers and TSS signers
// If the percentage is not set, it returns 0
func (k Keeper) GetDistributions(ctx sdk.Context) (sdkmath.Int, sdkmath.Int, sdkmath.Int) {
	// Fetch the validator rewards, use 0 if the percentage is not set
	validatorRewards := sdk.NewInt(0)
	validatorRewardsDec, err := sdk.NewDecFromStr(k.GetParamSetIfExists(ctx).ValidatorEmissionPercentage)
	if err == nil {
		validatorRewards = validatorRewardsDec.Mul(types.BlockReward).TruncateInt()
	}

	// Fetch the observer rewards, use 0 if the percentage is not set
	observerRewards := sdk.NewInt(0)
	observerRewardsDec, err := sdk.NewDecFromStr(k.GetParamSetIfExists(ctx).ObserverEmissionPercentage)
	if err == nil {
		observerRewards = observerRewardsDec.Mul(types.BlockReward).TruncateInt()
	}

	// Fetch the TSS signer rewards, use 0 if the percentage is not set
	tssSignerRewards := sdk.NewInt(0)
	tssSignerRewardsDec, err := sdk.NewDecFromStr(k.GetParamSetIfExists(ctx).TssSignerEmissionPercentage)
	if err == nil {
		tssSignerRewards = tssSignerRewardsDec.Mul(types.BlockReward).TruncateInt()
	}

	return validatorRewards, observerRewards, tssSignerRewards
}
