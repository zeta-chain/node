package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func (k Keeper) GetBlockRewardComponents(ctx sdk.Context) (sdk.Dec, sdk.Dec, sdk.Dec) {
	reservesFactor := k.GetReservesFactor(ctx)
	if reservesFactor.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()
	}
	bondFactor := k.GetBondFactor(ctx, k.GetStakingKeeper())
	durationFactor := k.GetDurationFactor(ctx)
	return reservesFactor, bondFactor, durationFactor
}
func (k Keeper) GetBondFactor(ctx sdk.Context, stakingKeeper types.StakingKeeper) sdk.Dec {
	targetBondRatio := sdk.MustNewDecFromStr(k.GetParams(ctx).TargetBondRatio)
	maxBondFactor := sdk.MustNewDecFromStr(k.GetParams(ctx).MaxBondFactor)
	minBondFactor := sdk.MustNewDecFromStr(k.GetParams(ctx).MinBondFactor)

	currentBondedRatio := stakingKeeper.BondedRatio(ctx)
	// Bond factor ranges between minBondFactor (0.75) to maxBondFactor (1.25)
	if currentBondedRatio.IsZero() {
		return sdk.ZeroDec()
	}
	bondFactor := targetBondRatio.Quo(currentBondedRatio)
	if bondFactor.GT(maxBondFactor) {
		return maxBondFactor
	}
	if bondFactor.LT(minBondFactor) {
		return minBondFactor
	}
	return bondFactor
}

func (k Keeper) GetDurationFactor(ctx sdk.Context) sdk.Dec {
	avgBlockTime := sdk.MustNewDecFromStr(k.GetParams(ctx).AvgBlockTime)
	NumberOfBlocksInAMonth := sdk.NewDec(types.SecsInMonth).Quo(avgBlockTime)
	monthFactor := sdk.NewDec(ctx.BlockHeight()).Quo(NumberOfBlocksInAMonth)
	logValueDec := sdk.MustNewDecFromStr(k.GetParams(ctx).DurationFactorConstant)
	// month * log(1 + 0.02 / 12)
	fractionNumerator := monthFactor.Mul(logValueDec)
	// (month * log(1 + 0.02 / 12) ) + 1
	fractionDenominator := fractionNumerator.Add(sdk.OneDec())

	// (month * log(1 + 0.02 / 12)) / (month * log(1 + 0.02 / 12) ) + 1
	if fractionDenominator.IsZero() {
		return sdk.OneDec()
	}
	if fractionNumerator.IsZero() {
		return sdk.ZeroDec()
	}
	return fractionNumerator.Quo(fractionDenominator)
}

func (k Keeper) GetReservesFactor(ctx sdk.Context) sdk.Dec {
	reserveAmount := k.GetBankKeeper().GetBalance(ctx, types.EmissionsModuleAddress, config.BaseDenom)
	return sdk.NewDecFromInt(reserveAmount.Amount)
}

func (k Keeper) GetFixedBlockRewards() (sdk.Dec, error) {
	return CalculateFixedValidatorRewards(types.AvgBlockTime)
}

func CalculateFixedValidatorRewards(avgBlockTimeString string) (sdk.Dec, error) {
	azetaAmountTotalRewards, err := common.GetAzetaDecFromAmountInZeta(types.BlockRewardsInZeta)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	avgBlockTime, err := sdk.NewDecFromStr(avgBlockTimeString)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	numberOfBlocksInAMonth := sdk.NewDec(types.SecsInMonth).Quo(avgBlockTime)
	numberOfBlocksTotal := numberOfBlocksInAMonth.Mul(sdk.NewDec(12)).Mul(sdk.NewDec(types.EmissionScheduledYears))
	constantRewardPerBlock := azetaAmountTotalRewards.Quo(numberOfBlocksTotal)
	return constantRewardPerBlock, nil
}
