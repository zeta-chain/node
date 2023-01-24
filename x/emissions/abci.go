package emissions

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/x/emissions/keeper"
	"github.com/zeta-chain/zetacore/x/emissions/types"
	"math"
	"math/big"
)

func BeginBlocker(ctx sdk.Context, keeper keeper.Keeper, stakingKeeper types.StakingKeeper, bankKeeper types.BankKeeper) {
	blockRewards := GetReservesFactor(ctx, keeper).
		Mul(GetBondFactor(ctx, stakingKeeper, keeper)).
		Mul(GetDurationFactor(ctx, keeper))

	blockRewardsInt := blockRewards.TruncateInt()
	blockRewardsCoins := sdk.NewCoin(config.BaseDenom, blockRewardsInt)

	err := bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, keeper.GetFeeCollector(), sdk.NewCoins(blockRewardsCoins))
	if err != nil {
		panic(err)
	}

	tracker, _ := keeper.GetEmissionTracker(ctx, types.EmissionCategory_ValidatorEmission)
	tracker.AmountLeft = tracker.AmountLeft.Sub(blockRewardsInt)
	keeper.SetEmissionTracker(ctx, &tracker)

	fmt.Println("blockRewards", blockRewards, blockRewardsInt)
}

func GetBondFactor(ctx sdk.Context, stakingKeeper types.StakingKeeper, keeper keeper.Keeper) sdk.Dec {
	targetBondRatio := sdk.MustNewDecFromStr(keeper.GetParams(ctx).TargetBondRatio)
	maxBondFactor := sdk.MustNewDecFromStr(keeper.GetParams(ctx).MaxBondFactor)
	minBondFactor := sdk.MustNewDecFromStr(keeper.GetParams(ctx).MinBondFactor)

	currentBondedRatio := stakingKeeper.BondedRatio(ctx)
	// Bond factor ranges between minBondFactor (0.75) to maxBondFactor (1.25)
	bondFactor := targetBondRatio.Quo(currentBondedRatio)
	if bondFactor.GT(maxBondFactor) {
		return maxBondFactor
	}
	if bondFactor.LT(minBondFactor) {
		return minBondFactor
	}
	return bondFactor
}

func GetDurationFactor(ctx sdk.Context, keeper keeper.Keeper) sdk.Dec {
	avgBlockTime := sdk.MustNewDecFromStr(keeper.GetParams(ctx).AvgBlockTime)
	NumberOfBlocksInAMonth := sdk.NewDec(types.SecsInMonth).Quo(avgBlockTime)
	monthFactor := sdk.NewDec(ctx.BlockHeight()).Quo(NumberOfBlocksInAMonth)
	//log(1 + 0.02 / 12)
	fractionConstant := 0.02 / 12.00
	logValue := math.Log(1.0 + fractionConstant)
	logValueDec, _ := sdk.NewDecFromStr(big.NewFloat(logValue).String())
	// month * log(1 + 0.02 / 12)
	fractionNumerator := monthFactor.Mul(logValueDec)
	// (month * log(1 + 0.02 / 12) ) + 1
	fractionDenominator := fractionNumerator.Add(sdk.OneDec())
	// (month * log(1 + 0.02 / 12)) / (month * log(1 + 0.02 / 12) ) + 1
	durationFactor := fractionNumerator.Quo(fractionDenominator)

	return durationFactor
}

func GetReservesFactor(ctx sdk.Context, keeper keeper.Keeper) sdk.Dec {
	reserveAmount, _ := keeper.GetEmissionTracker(ctx, types.EmissionCategory_ValidatorEmission)
	return reserveAmount.AmountLeft.ToDec()
}
