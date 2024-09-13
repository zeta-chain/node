package emissions

import (
	"fmt"
	"sort"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/x/emissions/keeper"
	"github.com/zeta-chain/node/x/emissions/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func BeginBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	emissionPoolBalance := keeper.GetReservesFactor(ctx)

	// Get the block rewards from the params
	params, found := keeper.GetParams(ctx)
	if !found {
		ctx.Logger().Error("Params not found")
		return
	}
	blockRewards := params.BlockRewardAmount
	if blockRewards.GT(emissionPoolBalance) {
		ctx.Logger().
			Info(fmt.Sprintf("Block rewards %s are greater than emission pool balance %s", blockRewards.String(), emissionPoolBalance.String()))
		return
	}

	// Get the distribution of rewards
	validatorRewards, observerRewards, tssSignerRewards := types.GetRewardsDistributions(params)

	// Use a tmpCtx, which is a cache-wrapped context to avoid writing to the store
	// We commit only if all three distributions are successful, if not the funds stay in the emission pool
	tmpCtx, commit := ctx.CacheContext()
	err := DistributeValidatorRewards(tmpCtx, validatorRewards, keeper.GetBankKeeper(), keeper.GetFeeCollector())
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error while distributing validator rewards %s", err))
		return
	}
	err = DistributeObserverRewards(tmpCtx, observerRewards, keeper, params)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error while distributing observer rewards %s", err))
		return
	}
	err = DistributeTssRewards(tmpCtx, tssSignerRewards, keeper.GetBankKeeper())
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error while distributing tss signer rewards %s", err))
		return
	}
	commit()

	types.EmitValidatorEmissions(ctx, "", "",
		"",
		validatorRewards.String(),
		observerRewards.String(),
		tssSignerRewards.String())
}

// DistributeValidatorRewards distributes the rewards to validators who signed the block .
// The block proposer gets a bonus reward
// This function uses the distribution module of cosmos-sdk , by directly sending funds to the feecollector.
func DistributeValidatorRewards(
	ctx sdk.Context,
	amount sdkmath.Int,
	bankKeeper types.BankKeeper,
	feeCollector string,
) error {
	coin := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount))
	ctx.Logger().
		Info(fmt.Sprintf("Distributing Validator Rewards Total:%s To FeeCollector : %s", amount.String(), feeCollector))
	return bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, feeCollector, coin)
}

// DistributeObserverRewards distributes the rewards to all observers who voted in any of the matured ballots
// The total rewards are distributed equally among all Successful votes
// NotVoted or Unsuccessful votes are slashed
// rewards given or slashed amounts are in azeta
func DistributeObserverRewards(
	ctx sdk.Context,
	amount sdkmath.Int,
	keeper keeper.Keeper,
	params types.Params,
) error {
	slashAmount := params.ObserverSlashAmount
	rewardsDistributer := map[string]int64{}
	totalRewardsUnits := int64(0)
	err := keeper.GetBankKeeper().
		SendCoinsFromModuleToModule(ctx, types.ModuleName, types.UndistributedObserverRewardsPool, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)))
	if err != nil {
		return err
	}

	list, found := keeper.GetObserverKeeper().GetMaturedBallots(ctx, params.BallotMaturityBlocks)
	ballotIdentifiers := []string{}
	if found {
		ballotIdentifiers = list.BallotsIndexList
	}

	// do not distribute rewards if no ballots are matured, the rewards can accumulate in the undistributed pool
	if len(ballotIdentifiers) == 0 {
		return nil
	}
	ballots := make([]observertypes.Ballot, 0, len(ballotIdentifiers))
	for _, ballotIdentifier := range ballotIdentifiers {
		ballot, found := keeper.GetObserverKeeper().GetBallot(ctx, ballotIdentifier)
		if !found {
			continue
		}
		ballots = append(ballots, ballot)
		totalRewardsUnits += ballot.BuildRewardsDistribution(rewardsDistributer)
	}
	rewardPerUnit := sdkmath.ZeroInt()
	if totalRewardsUnits > 0 && amount.IsPositive() {
		rewardPerUnit = amount.Quo(sdk.NewInt(totalRewardsUnits))
	}
	ctx.Logger().
		Debug(fmt.Sprintf("Total Rewards Units : %d , rewards per Unit %s ,number of ballots :%d", totalRewardsUnits, rewardPerUnit.String(), len(ballotIdentifiers)))
	sortedKeys := make([]string, 0, len(rewardsDistributer))
	for k := range rewardsDistributer {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	var finalDistributionList []*types.ObserverEmission
	for _, key := range sortedKeys {
		observerAddress, err := sdk.AccAddressFromBech32(key)
		if err != nil {
			ctx.Logger().Error("Error while parsing observer address ", "error", err, "address", key)
			continue
		}
		// observerRewardUnits can be negative if the observer has been slashed
		// an observers earn 1 unit for a correct vote, and -1 unit for an incorrect vote
		observerRewardUnits := rewardsDistributer[key]

		if observerRewardUnits == 0 {
			finalDistributionList = append(finalDistributionList, &types.ObserverEmission{
				EmissionType:    types.EmissionType_Slash,
				ObserverAddress: observerAddress.String(),
				Amount:          sdkmath.ZeroInt(),
			})
			continue
		}
		if observerRewardUnits < 0 {
			keeper.SlashObserverEmission(ctx, observerAddress.String(), slashAmount)
			finalDistributionList = append(finalDistributionList, &types.ObserverEmission{
				EmissionType:    types.EmissionType_Slash,
				ObserverAddress: observerAddress.String(),
				Amount:          slashAmount,
			})
			continue
		}

		// Defensive check
		if rewardPerUnit.GT(sdk.ZeroInt()) {
			rewardAmount := rewardPerUnit.Mul(sdkmath.NewInt(observerRewardUnits))
			keeper.AddObserverEmission(ctx, observerAddress.String(), rewardAmount)
			finalDistributionList = append(finalDistributionList, &types.ObserverEmission{
				EmissionType:    types.EmissionType_Rewards,
				ObserverAddress: observerAddress.String(),
				Amount:          rewardAmount,
			})
		}
	}
	types.EmitObserverEmissions(ctx, finalDistributionList)

	// Clear the matured ballots
	keeper.GetObserverKeeper().ClearMaturedBallots(ctx, ballots, params.BallotMaturityBlocks)
	return nil
}

// DistributeTssRewards trasferes the allocated rewards to the Undistributed Tss Rewards Pool.
// This is done so that the reserves factor is properly calculated in the next block
func DistributeTssRewards(ctx sdk.Context, amount sdk.Int, bankKeeper types.BankKeeper) error {
	coin := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount))
	return bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.UndistributedTssRewardsPool, coin)
}
