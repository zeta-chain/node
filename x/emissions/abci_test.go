package emissions_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	emissionsModule "github.com/zeta-chain/zetacore/x/emissions"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestBeginBlocker(t *testing.T) {
	t.Run("no observer distribution happens if emissions module account is empty", func(t *testing.T) {
		k, ctx, _, zk := keepertest.EmissionsKeeper(t)
		var ballotIdentifiers []string

		observerSet := sample.ObserverSet(10)
		zk.ObserverKeeper.SetObserverSet(ctx, observerSet)

		ballotList := sample.BallotList(10, observerSet.ObserverList)
		for _, ballot := range ballotList {
			zk.ObserverKeeper.SetBallot(ctx, &ballot)
			ballotIdentifiers = append(ballotIdentifiers, ballot.BallotIdentifier)
		}
		zk.ObserverKeeper.SetBallotList(ctx, &observerTypes.BallotListForHeight{
			Height:           0,
			BallotsIndexList: ballotIdentifiers,
		})
		for i := 0; i < 100; i++ {
			emissionsModule.BeginBlocker(ctx, *k)
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		}
		for _, observer := range observerSet.ObserverList {
			_, found := k.GetWithdrawableEmission(ctx, observer)
			require.False(t, found)
		}
	})
	t.Run("no validator distribution happens if emissions module account is empty", func(t *testing.T) {
		k, ctx, sk, _ := keepertest.EmissionsKeeper(t)
		feeCollectorAddress := sk.AuthKeeper.GetModuleAccount(ctx, types.FeeCollectorName).GetAddress()
		for i := 0; i < 100; i++ {
			emissionsModule.BeginBlocker(ctx, *k)
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		}
		require.True(t, sk.BankKeeper.GetBalance(ctx, feeCollectorAddress, config.BaseDenom).Amount.IsZero())
	})
	t.Run("tmp ctx is not committed if any of the distribution fails", func(t *testing.T) {
		k, ctx, sk, _ := keepertest.EmissionsKeeper(t)
		// Fund the emission pool to start the emission process
		err := sk.BankKeeper.MintCoins(ctx, emissionstypes.ModuleName, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewInt(1000000000000))))
		require.NoError(t, err)
		// Setup module accounts for emission pools except for observer pool , so that the observer distribution fails
		_ = sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedTssRewardsPool).GetAddress()
		feeCollectorAddress := sk.AuthKeeper.GetModuleAccount(ctx, types.FeeCollectorName).GetAddress()
		_ = sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.ModuleName).GetAddress()

		for i := 0; i < 100; i++ {
			// produce a block
			emissionsModule.BeginBlocker(ctx, *k)
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		}
		require.True(t, sk.BankKeeper.GetBalance(ctx, feeCollectorAddress, config.BaseDenom).Amount.IsZero())
		require.True(t, sk.BankKeeper.GetBalance(ctx, emissionstypes.EmissionsModuleAddress, config.BaseDenom).Amount.Equal(sdk.NewInt(1000000000000)))
	})
	t.Run("successfully distribute rewards", func(t *testing.T) {
		numberOfTestBlocks := 100
		k, ctx, sk, zk := keepertest.EmissionsKeeper(t)
		observerSet := sample.ObserverSet(10)
		zk.ObserverKeeper.SetObserverSet(ctx, observerSet)
		ballotList := sample.BallotList(10, observerSet.ObserverList)

		// set the ballot list
		ballotIdentifiers := []string{}
		for _, ballot := range ballotList {
			zk.ObserverKeeper.SetBallot(ctx, &ballot)
			ballotIdentifiers = append(ballotIdentifiers, ballot.BallotIdentifier)
		}
		zk.ObserverKeeper.SetBallotList(ctx, &observerTypes.BallotListForHeight{
			Height:           0,
			BallotsIndexList: ballotIdentifiers,
		})

		// Total block rewards is the fixed amount of rewards that are distributed
		totalBlockRewards, err := common.GetAzetaDecFromAmountInZeta(emissionstypes.BlockRewardsInZeta)
		totalRewardCoins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, totalBlockRewards.TruncateInt()))
		require.NoError(t, err)
		// Fund the emission pool to start the emission process
		err = sk.BankKeeper.MintCoins(ctx, emissionstypes.ModuleName, totalRewardCoins)
		require.NoError(t, err)

		// Setup module accounts for emission pools
		undistributedObserverPoolAddress := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedObserverRewardsPool).GetAddress()
		undistributedTssPoolAddress := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedTssRewardsPool).GetAddress()
		feeCollecterAddress := sk.AuthKeeper.GetModuleAccount(ctx, types.FeeCollectorName).GetAddress()
		emissionPool := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.ModuleName).GetAddress()

		blockRewards := emissionstypes.BlockReward
		observerRewardsForABlock := blockRewards.Mul(sdk.MustNewDecFromStr(k.GetParams(ctx).ObserverEmissionPercentage)).TruncateInt()
		validatorRewardsForABlock := blockRewards.Mul(sdk.MustNewDecFromStr(k.GetParams(ctx).ValidatorEmissionPercentage)).TruncateInt()
		tssSignerRewardsForABlock := blockRewards.Mul(sdk.MustNewDecFromStr(k.GetParams(ctx).TssSignerEmissionPercentage)).TruncateInt()
		distributedRewards := observerRewardsForABlock.Add(validatorRewardsForABlock).Add(tssSignerRewardsForABlock)

		require.True(t, blockRewards.TruncateInt().GT(distributedRewards))

		for i := 0; i < numberOfTestBlocks; i++ {
			emissionPoolBeforeBlockDistribution := sk.BankKeeper.GetBalance(ctx, emissionPool, config.BaseDenom).Amount
			// produce a block
			emissionsModule.BeginBlocker(ctx, *k)

			// require distribution amount
			emissionPoolBalanceAfterBlockDistribution := sk.BankKeeper.GetBalance(ctx, emissionPool, config.BaseDenom).Amount
			require.True(t, emissionPoolBeforeBlockDistribution.Sub(emissionPoolBalanceAfterBlockDistribution).Equal(distributedRewards))

			// totalDistributedTillCurrentBlock is the net amount of rewards distributed till the current block, this works in a unit test as the fees are not being collected by validators
			totalDistributedTillCurrentBlock := sk.BankKeeper.GetBalance(ctx, feeCollecterAddress, config.BaseDenom).Amount.
				Add(sk.BankKeeper.GetBalance(ctx, undistributedObserverPoolAddress, config.BaseDenom).Amount).
				Add(sk.BankKeeper.GetBalance(ctx, undistributedTssPoolAddress, config.BaseDenom).Amount)
			// require we are always under the max limit of block rewards
			require.True(t, totalRewardCoins.AmountOf(config.BaseDenom).
				Sub(totalDistributedTillCurrentBlock).GTE(sdk.ZeroInt()))

			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		}

		// We can simplify the calculation as the rewards are distributed equally among all the observers
		rewardPerUnit := observerRewardsForABlock.Quo(sdk.NewInt(int64(len(ballotList) * len(observerSet.ObserverList))))
		emissionAmount := rewardPerUnit.Mul(sdk.NewInt(int64(len(ballotList))))

		// Check if the rewards are distributed equally among all the observers
		for _, observer := range observerSet.ObserverList {
			observerEmission, found := k.GetWithdrawableEmission(ctx, observer)
			require.True(t, found)
			require.Equal(t, emissionAmount, observerEmission.Amount)
		}

		// Check pool balances after the distribution
		feeCollectorBalance := sk.BankKeeper.GetBalance(ctx, feeCollecterAddress, config.BaseDenom).Amount
		require.Equal(t, feeCollectorBalance, validatorRewardsForABlock.Mul(sdk.NewInt(int64(numberOfTestBlocks))))

		tssPoolBalances := sk.BankKeeper.GetBalance(ctx, undistributedTssPoolAddress, config.BaseDenom).Amount
		require.Equal(t, tssSignerRewardsForABlock.Mul(sdk.NewInt(int64(numberOfTestBlocks))).String(), tssPoolBalances.String())

		observerPoolBalances := sk.BankKeeper.GetBalance(ctx, undistributedObserverPoolAddress, config.BaseDenom).Amount
		require.Equal(t, observerRewardsForABlock.Mul(sdk.NewInt(int64(numberOfTestBlocks))).String(), observerPoolBalances.String())
	})
}

func TestDistributeObserverRewards(t *testing.T) {

	k, ctx, sk, zk := keepertest.EmissionsKeeper(t)
	observerSet := sample.ObserverSet(4)
	zk.ObserverKeeper.SetObserverSet(ctx, observerSet)
	// Total block rewards is the fixed amount of rewards that are distributed
	totalBlockRewards, err := common.GetAzetaDecFromAmountInZeta(emissionstypes.BlockRewardsInZeta)
	totalRewardCoins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, totalBlockRewards.TruncateInt()))
	require.NoError(t, err)
	// Fund the emission pool to start the emission process
	err = sk.BankKeeper.MintCoins(ctx, emissionstypes.ModuleName, totalRewardCoins)
	require.NoError(t, err)
	// Set starting emission for all observers to 100 so that we can calculate the rewards and slashes
	for _, observer := range observerSet.ObserverList {
		k.SetWithdrawableEmission(ctx, emissionstypes.WithdrawableEmissions{
			Address: observer,
			Amount:  sdkmath.NewInt(100),
		})
	}

	tt := []struct {
		name                 string
		votes                [][]observerTypes.VoteType
		totalRewardsForBlock sdkmath.Int
		expectedRewards      map[string]int64
		ballotStatus         observerTypes.BallotStatus
		slashAmount          sdkmath.Int
	}{
		{
			name:  "all observers rewarded correctly",
			votes: [][]observerTypes.VoteType{{observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation}},
			// total reward units would be 4 as all votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(100),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 125,
				observerSet.ObserverList[1]: 125,
				observerSet.ObserverList[2]: 125,
				observerSet.ObserverList[3]: 125,
			},
			ballotStatus: observerTypes.BallotStatus_BallotFinalized_SuccessObservation,
			slashAmount:  sdkmath.NewInt(25),
		},
		{
			name:  "one observer slashed",
			votes: [][]observerTypes.VoteType{{observerTypes.VoteType_FailureObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation}},
			// total reward units would be 3 as 3 votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(75),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 75,
				observerSet.ObserverList[1]: 125,
				observerSet.ObserverList[2]: 125,
				observerSet.ObserverList[3]: 125,
			},
			ballotStatus: observerTypes.BallotStatus_BallotFinalized_SuccessObservation,
			slashAmount:  sdkmath.NewInt(25),
		},
		{
			name:  "all observer slashed",
			votes: [][]observerTypes.VoteType{{observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation}},
			// total reward units would be 0 as no votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(100),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 75,
				observerSet.ObserverList[1]: 75,
				observerSet.ObserverList[2]: 75,
				observerSet.ObserverList[3]: 75,
			},
			ballotStatus: observerTypes.BallotStatus_BallotFinalized_FailureObservation,
			slashAmount:  sdkmath.NewInt(25),
		},
		{
			name:  "slashed to zero if slash amount is greater than available emissions",
			votes: [][]observerTypes.VoteType{{observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation}},
			// total reward units would be 0 as no votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(100),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 0,
				observerSet.ObserverList[1]: 0,
				observerSet.ObserverList[2]: 0,
				observerSet.ObserverList[3]: 0,
			},
			ballotStatus: observerTypes.BallotStatus_BallotFinalized_FailureObservation,
			slashAmount:  sdkmath.NewInt(2500),
		},
		{
			name: "withdraw able emissions unchanged if rewards and slashes are equal",
			votes: [][]observerTypes.VoteType{
				{observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation},
				{observerTypes.VoteType_FailureObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_SuccessObservation},
			},
			// total reward units would be 7 as 7 votes match the ballot status, including both ballots
			totalRewardsForBlock: sdkmath.NewInt(70),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 100,
				observerSet.ObserverList[1]: 120,
				observerSet.ObserverList[2]: 120,
				observerSet.ObserverList[3]: 120,
			},
			ballotStatus: observerTypes.BallotStatus_BallotFinalized_SuccessObservation,
			slashAmount:  sdkmath.NewInt(25),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			params := emissionstypes.DefaultParams()
			params.ObserverSlashAmount = tc.slashAmount
			k.SetParams(ctx, params)
			ballotIdentifiers := []string{}
			for i, votes := range tc.votes {
				ballot := observerTypes.Ballot{
					BallotIdentifier: "ballot" + string(rune(i)),
					BallotStatus:     tc.ballotStatus,
					VoterList:        observerSet.ObserverList,
					Votes:            votes,
				}
				zk.ObserverKeeper.SetBallot(ctx, &ballot)
				ballotIdentifiers = append(ballotIdentifiers, ballot.BallotIdentifier)
			}
			zk.ObserverKeeper.SetBallotList(ctx, &observerTypes.BallotListForHeight{
				Height:           0,
				BallotsIndexList: ballotIdentifiers,
			})

			ctx = ctx.WithBlockHeight(100)
			err := emissionsModule.DistributeObserverRewards(ctx, tc.totalRewardsForBlock, *k)
			require.NoError(t, err)
			for _, observer := range observerSet.ObserverList {
				observerEmission, found := k.GetWithdrawableEmission(ctx, observer)
				require.True(t, found)
				require.Equal(t, tc.expectedRewards[observer], observerEmission.Amount.Int64())
			}
		})
	}
}
