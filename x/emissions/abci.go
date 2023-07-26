package emissions

import (
	sdkmath "cosmossdk.io/math"
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
	err = DistributeObserverRewards(ctx, observerRewards, keeper)
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

func DistributeObserverRewards(ctx sdk.Context, amount sdkmath.Int, keeper keeper.Keeper) error {
	ballots := keeper.GetObserverKeeper().GetFinalizedBallots(ctx)
	rewardsDistributer := map[string]int64{}
	totalRewardsUnits := int64(0)
	err := keeper.GetBankKeeper().SendCoinsFromModuleToModule(ctx, types.ModuleName, types.UndistributedObserverRewardsPool, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)))
	if err != nil {
		return err
	}

	if len(ballots) == 0 {
		return nil
	}
	for _, ballot := range ballots {
		totalRewardsUnits = totalRewardsUnits + ballot.BuildRewardsDistribution(rewardsDistributer)
	}
	if totalRewardsUnits <= 0 || amount.IsZero() {
		return nil
	}

	rewardPerUnit := amount.Quo(sdk.NewInt(totalRewardsUnits))
	sortedKeys := make([]string, 0, len(rewardsDistributer))
	for k := range rewardsDistributer {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var finalDistributionList []*types.ObserverEmission
	for _, key := range sortedKeys {
		observerAddress, err := sdk.AccAddressFromBech32(key)
		if err != nil {
			continue
		}
		observerRewardsRatio := rewardsDistributer[key]
		if observerRewardsRatio == 0 {
			finalDistributionList = append(finalDistributionList, &types.ObserverEmission{
				EmissionType:    types.EmissionType_Slash,
				ObserverAddress: observerAddress.String(),
				Amount:          sdkmath.ZeroInt().String(),
			})
			continue
		}
		if observerRewardsRatio < 0 {
			keeper.SlashRewards(ctx, observerAddress.String(), sdkmath.NewInt(10000))
			finalDistributionList = append(finalDistributionList, &types.ObserverEmission{
				EmissionType:    types.EmissionType_Slash,
				ObserverAddress: observerAddress.String(),
				Amount:          sdkmath.NewInt(10000).String(),
			})
			continue
		}
		rewardAmount := rewardPerUnit.Mul(sdkmath.NewInt(observerRewardsRatio))
		keeper.AddRewards(ctx, observerAddress.String(), rewardAmount)
		finalDistributionList = append(finalDistributionList, &types.ObserverEmission{
			EmissionType:    types.EmissionType_Rewards,
			ObserverAddress: observerAddress.String(),
			Amount:          rewardAmount.String(),
		})
	}
	types.EmitObserverEmissions(ctx, finalDistributionList)

	keeper.GetObserverKeeper().DeleteFinalizedBallots(ctx)
	return nil
}

func DistributeTssRewards(ctx sdk.Context, amount sdk.Int, bankKeeper types.BankKeeper) error {
	coin := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount))
	return bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.UndistributedTssRewardsPool, coin)
}
