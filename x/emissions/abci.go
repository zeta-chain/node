package emissions

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/x/emissions/keeper"
	"github.com/zeta-chain/zetacore/x/emissions/types"
	"sort"
)

func BeginBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	reservesFactor, bondFactor, durationFactor := keeper.GetBlockRewardComponents(ctx)
	blockRewards := reservesFactor.Mul(bondFactor).Mul(durationFactor)
	if blockRewards.IsZero() {
		return
	}
	validatorRewards := sdk.MustNewDecFromStr(keeper.GetParams(ctx).ValidatorEmissionPercentage).Mul(blockRewards).TruncateInt()
	observerRewards := sdk.MustNewDecFromStr(keeper.GetParams(ctx).ObserverEmissionPercentage).Mul(blockRewards).TruncateInt()
	tssSignerRewards := sdk.MustNewDecFromStr(keeper.GetParams(ctx).TssSignerEmissionPercentage).Mul(blockRewards).TruncateInt()
	err := DistributeValidatorRewards(ctx, validatorRewards, keeper.GetBankKeeper(), keeper.GetFeeCollector())
	if err != nil {
		panic(err)
	}
	err = DistributeObserverRewards(ctx, observerRewards, keeper.GetBankKeeper(), keeper.GetObserverKeeper())
	if err != nil {
		panic(err)
	}
	err = DistributeTssRewards(ctx, tssSignerRewards, keeper.GetBankKeeper())
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
	//fmt.Println("amount", amount, ctx.BlockHeight())
	//for _, ballot := range ballots {
	//	fmt.Println("ballot : ", ballot.String())
	//}
	rewardsDistributer := map[string]int64{}
	totalRewardsUnits := int64(0)
	if len(ballots) == 0 {
		return nil
	}
	for _, ballot := range ballots {
		totalRewardsUnits = totalRewardsUnits + ballot.BuildRewardsDistribution(rewardsDistributer)
	}
	if totalRewardsUnits == 0 || amount.IsZero() {
		return nil
	}
	rewardPerUnit := amount.Quo(sdk.NewInt(totalRewardsUnits))
	sortedKeys := make([]string, 0, len(rewardsDistributer))
	for k := range rewardsDistributer {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	finalDistributionList := []string{}
	for _, key := range sortedKeys {
		observerAddress := sdk.AccAddress(key)
		rewardUnits := rewardsDistributer[key]
		if rewardUnits <= 0 {
			continue
		}
		rewardAmount := rewardPerUnit.Mul(sdk.NewInt(rewardUnits))
		finalDistributionList = append(finalDistributionList, fmt.Sprintf("%s:%s", observerAddress.String(), rewardAmount.String()))
		err := bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, observerAddress, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, rewardAmount)))
		if err != nil {
			// Check : Should we panic here, or add undistributed rewards to a keeper?
			ctx.Logger().Error("Error while distributing observer rewards", "observer", observerAddress.String(), "amount", rewardAmount, "error", err.Error())
			continue
		}
	}
	for _, d := range finalDistributionList {
		fmt.Println(d)
	}
	fmt.Println("-----------------------------------")
	obsKeeper.DeleteFinalizedBallots(ctx)
	return nil
}

func DistributeTssRewards(ctx sdk.Context, amount sdk.Int, bankKeeper types.BankKeeper) error {
	coin := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount))
	return bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.UndistributedTssRewardsPool, coin)
}
