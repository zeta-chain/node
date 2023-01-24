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
	blockRewards := GetReserservesFactor(ctx, keeper).Mul(GetBondFactor(ctx, stakingKeeper)).Mul(GetDurationFactor(ctx))

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

func GetBondFactor(ctx sdk.Context, stakingKeeper types.StakingKeeper) sdk.Dec {

	targetBondRatio, err := sdk.NewDecFromStr("67.00")
	if err != nil {
		fmt.Println(err)
	}
	maxBondFactor, _ := sdk.NewDecFromStr("1.25")
	minBondFactor, _ := sdk.NewDecFromStr("0.75")

	currentBondedRatio := stakingKeeper.BondedRatio(ctx)
	fmt.Println("currentBondedRatio : ", currentBondedRatio)
	bondFactor := targetBondRatio.Quo(currentBondedRatio)
	if bondFactor.GT(maxBondFactor) {
		bondFactor = maxBondFactor
	}
	if bondFactor.LT(minBondFactor) {
		bondFactor = minBondFactor
	}
	return bondFactor
}

func GetDurationFactor(ctx sdk.Context) sdk.Dec {
	avgBlockTime, _ := sdk.NewDecFromStr("6.0")

	NumberOfBlocksInAMonth := sdk.NewDec(30 * 24 * 60 * 60).Quo(avgBlockTime)
	months := sdk.NewDec(ctx.BlockHeight()).Quo(NumberOfBlocksInAMonth)
	numConstant := 0.02
	denomConstant := 12.00
	fractionConstant := numConstant / denomConstant

	logValue := math.Log(1.0 + fractionConstant)
	logValueDec, _ := sdk.NewDecFromStr(big.NewFloat(logValue).String())

	fractionNumerator := months.Mul(logValueDec)
	fractionDenominator := fractionNumerator.Add(sdk.OneDec())

	durationFactor := fractionNumerator.Quo(fractionDenominator)

	return durationFactor
}

func GetReserservesFactor(ctx sdk.Context, keeper keeper.Keeper) sdk.Dec {
	reserveAmount, _ := keeper.GetEmissionTracker(ctx, types.EmissionCategory_ValidatorEmission)
	return reserveAmount.AmountLeft.ToDec()
}
