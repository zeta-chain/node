package emissions

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/x/emissions/keeper"
	"github.com/zeta-chain/zetacore/x/emissions/types"
	"math"
	"math/big"
)

func BeginBlocker(ctx sdk.Context, keeper keeper.Keeper, stakingKeeper types.StakingKeeper, bankKeeper types.BankKeeper) {
	reservesFactor := GetReservesFactor(ctx, bankKeeper)
	if reservesFactor.LTE(sdk.ZeroDec()) {
		return
	}
	bondFactor := GetBondFactor(ctx, stakingKeeper, keeper)
	durationFactor := GetDurationFactor(ctx, keeper)
	blockRewards := reservesFactor.Mul(bondFactor).Mul(durationFactor)

	validatorRewards := sdk.NewCoin(config.BaseDenom, sdk.MustNewDecFromStr(keeper.GetParams(ctx).ValidatorEmissionPercentage).Mul(blockRewards).TruncateInt())
	observerRewards := sdk.NewCoin(config.BaseDenom, sdk.MustNewDecFromStr(keeper.GetParams(ctx).ObserverEmissionPercentage).Mul(blockRewards).TruncateInt())
	tssSignerRewards := sdk.NewCoin(config.BaseDenom, sdk.MustNewDecFromStr(keeper.GetParams(ctx).TssSignerEmissionPercentage).Mul(blockRewards).TruncateInt())

	err := DistributeValidatorRewards(ctx, sdk.NewCoins(validatorRewards), bankKeeper, keeper)
	if err != nil {
		panic(err)
	}
	DistributeObserverRewards(ctx, observerRewards, keeper)
	DistributeTssRewards(ctx, tssSignerRewards, keeper)
	types.EmitValidatorEmissions(ctx, bondFactor.String(), reservesFactor.String(),
		durationFactor.String(),
		validatorRewards.String(),
		observerRewards.String(),
		tssSignerRewards.String())
}

func DistributeValidatorRewards(ctx sdk.Context, amount sdk.Coins, bankKeeper types.BankKeeper, emissionKeeper keeper.Keeper) error {
	return bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, emissionKeeper.GetFeeCollector(), amount)
}

func DistributeObserverRewards(ctx sdk.Context, amount sdk.Coin, emissionKeeper keeper.Keeper) {
	tracker, _ := emissionKeeper.GetEmissionTracker(ctx, types.EmissionCategory_ObserverEmission)
	tracker.UndistributedAmount = tracker.UndistributedAmount.Add(amount)
	emissionKeeper.SetEmissionTracker(ctx, &tracker)
}

func DistributeTssRewards(ctx sdk.Context, amount sdk.Coin, emissionKeeper keeper.Keeper) {
	tracker, _ := emissionKeeper.GetEmissionTracker(ctx, types.EmissionCategory_ObserverEmission)
	tracker.UndistributedAmount = tracker.UndistributedAmount.Add(amount)
	emissionKeeper.SetEmissionTracker(ctx, &tracker)
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

func GetReservesFactor(ctx sdk.Context, keeper types.BankKeeper) sdk.Dec {
	reserveAmount := keeper.GetBalance(ctx, types.EmissionsModuleAddress, config.BaseDenom)
	return sdk.NewDecFromInt(reserveAmount.Amount)
}
