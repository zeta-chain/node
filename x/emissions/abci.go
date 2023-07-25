package emissions

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/x/emissions/keeper"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func BeginBlocker(ctx sdk.Context, keeper keeper.Keeper, stakingKeeper types.StakingKeeper, bankKeeper types.BankKeeper, observerKeeper types.ZetaObserverKeeper) {
	reservesFactor, bondFactor, durationFactor := GetBlockRewardComponents(ctx, bankKeeper, stakingKeeper, keeper)
	blockRewards := reservesFactor.Mul(bondFactor).Mul(durationFactor)
	if blockRewards.IsZero() {
		return
	}
	validatorRewards := sdk.MustNewDecFromStr(keeper.GetParams(ctx).ValidatorEmissionPercentage).Mul(blockRewards).TruncateInt()
	observerRewards := sdk.MustNewDecFromStr(keeper.GetParams(ctx).ObserverEmissionPercentage).Mul(blockRewards).TruncateInt()
	tssSignerRewards := sdk.MustNewDecFromStr(keeper.GetParams(ctx).TssSignerEmissionPercentage).Mul(blockRewards).TruncateInt()
	err := DistributeValidatorRewards(ctx, validatorRewards, bankKeeper, keeper.GetFeeCollector())
	if err != nil {
		panic(err)
	}
	err = DistributeObserverRewards(ctx, observerRewards, bankKeeper, observerKeeper)
	if err != nil {
		panic(err)
	}
	err = DistributeTssRewards(ctx, tssSignerRewards, bankKeeper)
	if err != nil {
		panic(err)
	}
	types.EmitValidatorEmissions(ctx, bondFactor.String(), reservesFactor.String(),
		durationFactor.String(),
		validatorRewards.String(),
		observerRewards.String(),
		tssSignerRewards.String())
}

func DistributeValidatorRewards(ctx sdk.Context, amount sdk.Int, bankKeeper types.BankKeeper, feeCollector string) error {
	coin := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount))
	return bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, feeCollector, coin)
}

func DistributeObserverRewards(ctx sdk.Context, amount sdkmath.Int, bankKeeper types.BankKeeper, obsKeeper types.ZetaObserverKeeper) error {
	ballots := obsKeeper.GetFinalizedBallots(ctx)
	rewardsDistributer := map[string]int64{}
	totalRewardsUnits := int64(0)
	for _, ballot := range ballots {
		totalRewardsUnits = totalRewardsUnits + ballot.BuildRewardsDistribution(rewardsDistributer)
	}
	rewardPerUnit := amount.Quo(sdk.NewInt(totalRewardsUnits))
	for observer, rewardUnits := range rewardsDistributer {
		if rewardUnits <= 0 {
			continue
		}
		rewardAmount := rewardPerUnit.Mul(sdk.NewInt(rewardUnits))
		err := bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(observer), sdk.NewCoins(sdk.NewCoin(config.BaseDenom, rewardAmount)))
		if err != nil {
			ctx.Logger().Error("Error while distributing observer rewards", "observer", observer, "amount", rewardAmount, "error", err.Error())
		}
	}
	return nil
}

func GetBlockRewardComponents(ctx sdk.Context, bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper, emissionKeeper keeper.Keeper) (sdk.Dec, sdk.Dec, sdk.Dec) {
	reservesFactor := GetReservesFactor(ctx, bankKeeper)
	if reservesFactor.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()
	}
	bondFactor := GetBondFactor(ctx, stakingKeeper, emissionKeeper)
	durationFactor := GetDurationFactor(ctx, emissionKeeper)
	return reservesFactor, bondFactor, durationFactor
}

func DistributeTssRewards(ctx sdk.Context, amount sdk.Int, bankKeeper types.BankKeeper) error {
	coin := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount))
	return bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.UndistributedTssRewardsPool, coin)
}
func GetBondFactor(ctx sdk.Context, stakingKeeper types.StakingKeeper, keeper keeper.Keeper) sdk.Dec {
	targetBondRatio := sdk.MustNewDecFromStr(keeper.GetParams(ctx).TargetBondRatio)
	maxBondFactor := sdk.MustNewDecFromStr(keeper.GetParams(ctx).MaxBondFactor)
	minBondFactor := sdk.MustNewDecFromStr(keeper.GetParams(ctx).MinBondFactor)

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

func GetDurationFactor(ctx sdk.Context, keeper keeper.Keeper) sdk.Dec {
	avgBlockTime := sdk.MustNewDecFromStr(keeper.GetParams(ctx).AvgBlockTime)
	NumberOfBlocksInAMonth := sdk.NewDec(types.SecsInMonth).Quo(avgBlockTime)
	monthFactor := sdk.NewDec(ctx.BlockHeight()).Quo(NumberOfBlocksInAMonth)
	logValueDec := sdk.MustNewDecFromStr(keeper.GetParams(ctx).DurationFactorConstant)
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

func GetReservesFactor(ctx sdk.Context, keeper types.BankKeeper) sdk.Dec {
	reserveAmount := keeper.GetBalance(ctx, types.EmissionsModuleAddress, config.BaseDenom)
	return sdk.NewDecFromInt(reserveAmount.Amount)
}
