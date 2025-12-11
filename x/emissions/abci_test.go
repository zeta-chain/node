package emissions_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/emissions"
	emissionskeeper "github.com/zeta-chain/node/x/emissions/keeper"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestBeginBlocker(t *testing.T) {
	t.Run("no distribution happens if params are not found", func(t *testing.T) {
		//Arrange
		k, ctx, _, zk := keepertest.EmissionsKeeper(t)
		_, found := k.GetParams(ctx)
		require.True(t, found)
		store := ctx.KVStore(k.GetStoreKey())
		store.Delete(emissionstypes.KeyPrefix(emissionstypes.ParamsKey))

		var ballotIdentifiers []string
		observerSet := sample.ObserverSet(10)
		zk.ObserverKeeper.SetObserverSet(ctx, observerSet)
		ballotList := sample.BallotList(10, observerSet.ObserverList)
		for _, ballot := range ballotList {
			zk.ObserverKeeper.SetBallot(ctx, &ballot)
			ballotIdentifiers = append(ballotIdentifiers, ballot.BallotIdentifier)
		}
		zk.ObserverKeeper.SetBallotList(ctx, &observertypes.BallotListForHeight{
			Height:           0,
			BallotsIndexList: ballotIdentifiers,
		})

		//Act
		for i := 0; i < 100; i++ {
			emissions.BeginBlocker(ctx, *k)
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		}

		//Assert
		for _, observer := range observerSet.ObserverList {
			_, found := k.GetWithdrawableEmission(ctx, observer)
			require.False(t, found)
		}
	})
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
		zk.ObserverKeeper.SetBallotList(ctx, &observertypes.BallotListForHeight{
			Height:           0,
			BallotsIndexList: ballotIdentifiers,
		})
		for i := 0; i < 100; i++ {
			emissions.BeginBlocker(ctx, *k)
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
			emissions.BeginBlocker(ctx, *k)
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		}
		require.True(t, sk.BankKeeper.GetBalance(ctx, feeCollectorAddress, config.BaseDenom).Amount.IsZero())
	})

	t.Run("tmp ctx is not committed if any of the distribution fails", func(t *testing.T) {
		k, ctx, sk, _ := keepertest.EmissionsKeeper(t)
		// Fund the emission pool to start the emission process
		blockRewards := emissionstypes.BlockReward
		err := sk.BankKeeper.MintCoins(
			ctx,
			emissionstypes.ModuleName,
			sdk.NewCoins(sdk.NewCoin(config.BaseDenom, blockRewards.TruncateInt())),
		)
		require.NoError(t, err)
		// Setup module accounts for emission pools except for observer pool , so that the observer distribution fails
		_ = sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedTSSRewardsPool).GetAddress()
		feeCollectorAddress := sk.AuthKeeper.GetModuleAccount(ctx, types.FeeCollectorName).GetAddress()
		_ = sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.ModuleName).GetAddress()

		for i := 0; i < 100; i++ {
			// produce a block
			emissions.BeginBlocker(ctx, *k)
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		}
		require.True(t, sk.BankKeeper.GetBalance(ctx, feeCollectorAddress, config.BaseDenom).Amount.IsZero())
		require.True(
			t,
			sk.BankKeeper.GetBalance(
				ctx,
				emissionstypes.EmissionsModuleAddress,
				config.BaseDenom,
			).Amount.Equal(
				blockRewards.TruncateInt(),
			),
		)
	})

	t.Run("begin blocker returns early if validator distribution fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionKeeperWithMockOptions(t, keepertest.EmissionMockOptions{
			UseBankMock: true,
		})
		// Over funding the emission pool to avoid any errors due to truncated values
		blockRewards := emissionstypes.BlockReward
		totalRewardAmount := blockRewards.TruncateInt().Mul(sdkmath.NewInt(2))
		totalRewardCoins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, totalRewardAmount))
		bankMock := keepertest.GetEmissionsBankMock(t, k)
		bankMock.On("GetBalance",
			ctx, mock.Anything, config.BaseDenom).
			Return(totalRewardCoins[0], nil).Once()

		// fail first distribution
		bankMock.On("SendCoinsFromModuleToModule",
			mock.Anything, emissionstypes.ModuleName, k.GetFeeCollector(), mock.Anything).
			Return(emissionstypes.ErrUnableToWithdrawEmissions).Once()
		emissions.BeginBlocker(ctx, *k)

		bankMock.AssertNumberOfCalls(t, "SendCoinsFromModuleToModule", 1)
	})

	t.Run("begin blocker returns early if observer distribution fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionKeeperWithMockOptions(t, keepertest.EmissionMockOptions{
			UseBankMock: true,
		})
		// Over funding the emission pool to avoid any errors due to truncated values
		blockRewards := emissionstypes.BlockReward
		totalRewardAmount := blockRewards.TruncateInt().Mul(sdkmath.NewInt(2))
		totalRewardCoins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, totalRewardAmount))
		bankMock := keepertest.GetEmissionsBankMock(t, k)
		bankMock.On("GetBalance",
			ctx, mock.Anything, config.BaseDenom).
			Return(totalRewardCoins[0], nil).Once()

		// allow first distribution
		bankMock.On("SendCoinsFromModuleToModule",
			mock.Anything, emissionstypes.ModuleName, k.GetFeeCollector(), mock.Anything).
			Return(nil).Once()

		// fail second distribution
		bankMock.On("SendCoinsFromModuleToModule",
			mock.Anything, emissionstypes.ModuleName, emissionstypes.UndistributedObserverRewardsPool, mock.Anything).
			Return(emissionstypes.ErrUnableToWithdrawEmissions).Once()
		emissions.BeginBlocker(ctx, *k)

		bankMock.AssertNumberOfCalls(t, "SendCoinsFromModuleToModule", 2)
	})

	t.Run("begin blocker returns early if tss distribution fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionKeeperWithMockOptions(t, keepertest.EmissionMockOptions{
			UseBankMock: true,
		})

		// Over funding the emission pool to avoid any errors due to truncated values
		blockRewards := emissionstypes.BlockReward
		totalRewardAmount := blockRewards.TruncateInt().Mul(sdkmath.NewInt(2))
		totalRewardCoins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, totalRewardAmount))
		bankMock := keepertest.GetEmissionsBankMock(t, k)
		bankMock.On("GetBalance",
			ctx, mock.Anything, config.BaseDenom).
			Return(totalRewardCoins[0], nil).Once()

		// allow first distribution
		bankMock.On("SendCoinsFromModuleToModule",
			mock.Anything, emissionstypes.ModuleName, k.GetFeeCollector(), mock.Anything).
			Return(nil).Once()

		// allow second distribution
		bankMock.On("SendCoinsFromModuleToModule",
			mock.Anything, emissionstypes.ModuleName, emissionstypes.UndistributedObserverRewardsPool, mock.Anything).
			Return(nil).Once()

		// fail third distribution
		bankMock.On("SendCoinsFromModuleToModule",
			mock.Anything, emissionstypes.ModuleName, emissionstypes.UndistributedTSSRewardsPool, mock.Anything).
			Return(emissionstypes.ErrUnableToWithdrawEmissions).Once()
		emissions.BeginBlocker(ctx, *k)

		bankMock.AssertNumberOfCalls(t, "SendCoinsFromModuleToModule", 3)
	})

	t.Run("successfully distribute rewards", func(t *testing.T) {
		//Arrange
		numberOfTestBlocks := 10
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
		zk.ObserverKeeper.SetBallotList(ctx, &observertypes.BallotListForHeight{
			Height:           0,
			BallotsIndexList: ballotIdentifiers,
		})

		// Fund the emission pool to start the emission process
		// Use Ceil() to ensure there's enough balance for the
		// blockRewards.GT(emissionPoolBalance) check in BeginBlocker, which compares the full decimal value
		blockRewards := emissionstypes.BlockReward
		totalRewardAmount := blockRewards.Ceil().TruncateInt().Mul(sdkmath.NewInt(int64(numberOfTestBlocks)))
		totalRewardCoins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, totalRewardAmount))

		err := sk.BankKeeper.MintCoins(ctx, emissionstypes.ModuleName, totalRewardCoins)
		require.NoError(t, err)

		// Setup module accounts for emission pools
		undistributedObserverPoolAddress := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedObserverRewardsPool).
			GetAddress()
		undistributedTssPoolAddress := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedTSSRewardsPool).
			GetAddress()
		feeCollecterAddress := sk.AuthKeeper.GetModuleAccount(ctx, types.FeeCollectorName).GetAddress()
		emissionPool := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.ModuleName).GetAddress()

		params, found := k.GetParams(ctx)
		require.True(t, found)
		// Set the ballot maturity blocks to numberOfTestBlocks so that the ballot mature at the end of the for loop which produces blocks
		params.BallotMaturityBlocks = int64(numberOfTestBlocks)
		err = k.SetParams(ctx, params)

		// Get the rewards distribution, this is a fixed amount based on total block rewards and distribution percentages
		validatorRewardsForABlock, observerRewardsForABlock, tssSignerRewardsForABlock := emissionstypes.GetRewardsDistributions(
			params,
		)

		distributedRewards := observerRewardsForABlock.Add(validatorRewardsForABlock).Add(tssSignerRewardsForABlock)
		// Block rewards should be >= distributed rewards.
		// They can be equal if the truncated block reward is perfectly divisible by 4 (for 50%+25%+25% split),
		// or slightly greater if truncation causes rounding losses in the distribution.
		require.True(t, blockRewards.TruncateInt().GTE(distributedRewards))

		require.Len(t, zk.ObserverKeeper.GetAllBallots(ctx), len(ballotList))
		_, found = zk.ObserverKeeper.GetBallotListForHeight(ctx, 0)
		require.True(t, found)

		// Act
		for i := 0; i < numberOfTestBlocks; i++ {
			emissionPoolBeforeBlockDistribution := sk.BankKeeper.GetBalance(ctx, emissionPool, config.BaseDenom).Amount
			// produce a block
			emissions.BeginBlocker(ctx, *k)

			// require distribution amount
			emissionPoolBalanceAfterBlockDistribution := sk.BankKeeper.GetBalance(
				ctx,
				emissionPool,
				config.BaseDenom,
			).Amount

			require.True(
				t,
				emissionPoolBeforeBlockDistribution.Sub(emissionPoolBalanceAfterBlockDistribution).
					Equal(distributedRewards),
			)

			// totalDistributedTillCurrentBlock is the net amount of rewards distributed till the current block, this works in a unit test as the fees are not being collected by validators
			totalDistributedTillCurrentBlock := sk.BankKeeper.GetBalance(ctx, feeCollecterAddress, config.BaseDenom).
				Amount.
				Add(sk.BankKeeper.GetBalance(ctx, undistributedObserverPoolAddress, config.BaseDenom).Amount).
				Add(sk.BankKeeper.GetBalance(ctx, undistributedTssPoolAddress, config.BaseDenom).Amount)
			// require we are always under the max limit of block rewards
			require.True(t, totalRewardCoins.AmountOf(config.BaseDenom).
				Sub(totalDistributedTillCurrentBlock).GTE(sdkmath.ZeroInt()))

			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		}

		// Assert

		// 1. Assert Observer rewards, these are distributed at the block in which the ballots mature.
		// numberOfTestBlocks is the same maturity blocks for the ballots

		// We can simplify the calculation as the rewards are distributed equally among all the observers
		rewardPerUnit := observerRewardsForABlock.Quo(
			sdkmath.NewInt(int64(len(ballotList) * len(observerSet.ObserverList))),
		)
		emissionAmount := rewardPerUnit.Mul(sdkmath.NewInt(int64(len(ballotList))))

		// 2 . Assert ballots and ballot list are deleted on maturity
		require.Len(t, zk.ObserverKeeper.GetAllBallots(ctx), 0)
		_, found = zk.ObserverKeeper.GetBallotListForHeight(ctx, 0)
		require.False(t, found)

		//3. Assert amounts in undistributed pools
		// Check if the rewards are distributed equally among all the observers
		for _, observer := range observerSet.ObserverList {
			observerEmission, found := k.GetWithdrawableEmission(ctx, observer)
			require.True(t, found)
			require.Equal(t, emissionAmount, observerEmission.Amount)
		}

		// Check pool balances after the distribution
		feeCollectorBalance := sk.BankKeeper.GetBalance(ctx, feeCollecterAddress, config.BaseDenom).Amount
		require.Equal(t, feeCollectorBalance, validatorRewardsForABlock.Mul(sdkmath.NewInt(int64(numberOfTestBlocks))))

		tssPoolBalances := sk.BankKeeper.GetBalance(ctx, undistributedTssPoolAddress, config.BaseDenom).Amount
		require.Equal(
			t,
			tssSignerRewardsForABlock.Mul(sdkmath.NewInt(int64(numberOfTestBlocks))).String(),
			tssPoolBalances.String(),
		)

		observerPoolBalances := sk.BankKeeper.GetBalance(ctx, undistributedObserverPoolAddress, config.BaseDenom).Amount
		require.Equal(
			t,
			observerRewardsForABlock.Mul(sdkmath.NewInt(int64(numberOfTestBlocks))).String(),
			observerPoolBalances.String(),
		)
	})
}

func TestDistributeObserverRewards(t *testing.T) {
	observerSet := sample.ObserverSet(4)

	tt := []struct {
		name                      string
		votes                     [][]observertypes.VoteType
		observerStartingEmissions sdkmath.Int
		totalRewardsForBlock      sdkmath.Int
		expectedRewards           map[string]int64
		ballotStatus              observertypes.BallotStatus
		slashAmount               sdkmath.Int
		rewardsPerBlock           sdkmath.LegacyDec
	}{
		{
			name: "all observers rewarded correctly",
			votes: [][]observertypes.VoteType{
				{
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
				},
			},
			observerStartingEmissions: sdkmath.NewInt(100),
			// total reward units would be 4 as all votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(100),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 125,
				observerSet.ObserverList[1]: 125,
				observerSet.ObserverList[2]: 125,
				observerSet.ObserverList[3]: 125,
			},
			ballotStatus:    observertypes.BallotStatus_BallotFinalized_SuccessObservation,
			slashAmount:     sdkmath.NewInt(25),
			rewardsPerBlock: emissionstypes.BlockReward,
		},
		{
			name: "no rewards if ballot is not finalized,irrespective of votes",
			votes: [][]observertypes.VoteType{
				{
					observertypes.VoteType_NotYetVoted,
					observertypes.VoteType_NotYetVoted,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_FailureObservation,
				},
			},
			observerStartingEmissions: sdkmath.NewInt(100),
			// total reward units would be 4 as all votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(0),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 100,
				observerSet.ObserverList[1]: 100,
				observerSet.ObserverList[2]: 100,
				observerSet.ObserverList[3]: 100,
			},
			ballotStatus:    observertypes.BallotStatus_BallotInProgress,
			slashAmount:     sdkmath.NewInt(25),
			rewardsPerBlock: emissionstypes.BlockReward,
		},
		{
			name: "one observer slashed",
			votes: [][]observertypes.VoteType{
				{
					observertypes.VoteType_FailureObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
				},
			},
			observerStartingEmissions: sdkmath.NewInt(100),
			// total reward units would be 3 as 3 votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(75),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 75,
				observerSet.ObserverList[1]: 125,
				observerSet.ObserverList[2]: 125,
				observerSet.ObserverList[3]: 125,
			},
			ballotStatus:    observertypes.BallotStatus_BallotFinalized_SuccessObservation,
			slashAmount:     sdkmath.NewInt(25),
			rewardsPerBlock: emissionstypes.BlockReward,
		},
		{
			name: "all observer slashed",
			votes: [][]observertypes.VoteType{
				{
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
				},
			},
			observerStartingEmissions: sdkmath.NewInt(100),
			// total reward units would be 0 as no votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(100),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 75,
				observerSet.ObserverList[1]: 75,
				observerSet.ObserverList[2]: 75,
				observerSet.ObserverList[3]: 75,
			},
			ballotStatus:    observertypes.BallotStatus_BallotFinalized_FailureObservation,
			slashAmount:     sdkmath.NewInt(25),
			rewardsPerBlock: emissionstypes.BlockReward,
		},
		{
			name: "slashed to zero if slash amount is greater than available emissions",
			votes: [][]observertypes.VoteType{
				{
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
				},
			},
			observerStartingEmissions: sdkmath.NewInt(100),
			// total reward units would be 0 as no votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(100),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 0,
				observerSet.ObserverList[1]: 0,
				observerSet.ObserverList[2]: 0,
				observerSet.ObserverList[3]: 0,
			},
			ballotStatus:    observertypes.BallotStatus_BallotFinalized_FailureObservation,
			slashAmount:     sdkmath.NewInt(2500),
			rewardsPerBlock: emissionstypes.BlockReward,
		},
		{
			name: "withdraw able emissions added only for correct votes",
			votes: [][]observertypes.VoteType{
				{
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
				},
				{
					observertypes.VoteType_FailureObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
				},
			},
			observerStartingEmissions: sdkmath.NewInt(100),
			// total reward units would be 7 as 7 votes match the ballot status, including both ballots
			totalRewardsForBlock: sdkmath.NewInt(70),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 100,
				observerSet.ObserverList[1]: 122,
				observerSet.ObserverList[2]: 122,
				observerSet.ObserverList[3]: 122,
			},
			ballotStatus:    observertypes.BallotStatus_BallotFinalized_SuccessObservation,
			slashAmount:     sdkmath.NewInt(25),
			rewardsPerBlock: emissionstypes.BlockReward,
		},
		{
			name: "no rewards if block reward is nil",
			votes: [][]observertypes.VoteType{
				{
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
				},
			},
			observerStartingEmissions: sdkmath.NewInt(0),
			// total reward units would be 4 as all votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(0),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 0,
				observerSet.ObserverList[1]: 0,
				observerSet.ObserverList[2]: 0,
				observerSet.ObserverList[3]: 0,
			},
			ballotStatus: observertypes.BallotStatus_BallotFinalized_SuccessObservation,
			slashAmount:  sdkmath.NewInt(25),
			//rewardsPerBlock: nil,
		},
		{
			name: "no rewards if block reward is negative",
			votes: [][]observertypes.VoteType{
				{
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
				},
			},
			observerStartingEmissions: sdkmath.NewInt(0),
			// total reward units would be 4 as all votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(0),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 0,
				observerSet.ObserverList[1]: 0,
				observerSet.ObserverList[2]: 0,
				observerSet.ObserverList[3]: 0,
			},
			ballotStatus:    observertypes.BallotStatus_BallotFinalized_SuccessObservation,
			slashAmount:     sdkmath.NewInt(25),
			rewardsPerBlock: sdkmath.LegacyNewDec(1).NegMut(),
		},
		{
			name: "no rewards if block reward is zero",
			votes: [][]observertypes.VoteType{
				{
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
					observertypes.VoteType_SuccessObservation,
				},
			},
			observerStartingEmissions: sdkmath.NewInt(0),
			// total reward units would be 4 as all votes match the ballot status
			totalRewardsForBlock: sdkmath.NewInt(0),
			expectedRewards: map[string]int64{
				observerSet.ObserverList[0]: 0,
				observerSet.ObserverList[1]: 0,
				observerSet.ObserverList[2]: 0,
				observerSet.ObserverList[3]: 0,
			},
			ballotStatus:    observertypes.BallotStatus_BallotFinalized_SuccessObservation,
			slashAmount:     sdkmath.NewInt(25),
			rewardsPerBlock: sdkmath.LegacyZeroDec(),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			// Keeper initialization
			k, ctx, sk, zk := keepertest.EmissionsKeeper(t)
			zk.ObserverKeeper.SetObserverSet(ctx, observerSet)

			// Total block rewards is the fixed amount of rewards that are distributed
			totalBlockRewards := emissionstypes.BlockReward.TruncateInt()
			totalRewardCoins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, totalBlockRewards))

			// Fund the emission pool to start the emission process
			err := sk.BankKeeper.MintCoins(ctx, emissionstypes.ModuleName, totalRewardCoins)
			require.NoError(t, err)

			// Set starting emission for all observers to a specified value so that we can calculate the rewards and slashes
			for _, observer := range observerSet.ObserverList {
				k.SetWithdrawableEmission(ctx, emissionstypes.WithdrawableEmissions{
					Address: observer,
					Amount:  tc.observerStartingEmissions,
				})
			}

			// Set the params
			params := emissionstypes.DefaultParams()
			params.ObserverSlashAmount = tc.slashAmount
			params.BlockRewardAmount = tc.rewardsPerBlock
			setEmissionsParams(t, ctx, *k, params)

			// Set the ballot list
			ballotIdentifiers := []string{}
			for i, votes := range tc.votes {
				ballot := observertypes.Ballot{
					BallotIdentifier: "ballot" + string(rune(i)),
					BallotStatus:     tc.ballotStatus,
					VoterList:        observerSet.ObserverList,
					Votes:            votes,
				}
				zk.ObserverKeeper.SetBallot(ctx, &ballot)
				ballotIdentifiers = append(ballotIdentifiers, ballot.BallotIdentifier)
			}
			zk.ObserverKeeper.SetBallotList(ctx, &observertypes.BallotListForHeight{
				Height:           0,
				BallotsIndexList: ballotIdentifiers,
			})
			ctx = ctx.WithBlockHeight(300)

			// Act
			// Distribute the rewards and check if the rewards are distributed correctly
			err = emissions.DistributeObserverRewards(ctx, tc.totalRewardsForBlock, *k, params)

			// Assert
			require.NoError(t, err)
			for i, observer := range observerSet.ObserverList {
				observerEmission, found := k.GetWithdrawableEmission(ctx, observer)
				require.True(t, found, "withdrawable emission not found for observer %d", i)
				require.Equal(
					t,
					tc.expectedRewards[observer],
					observerEmission.Amount.Int64(),
					"invalid withdrawable emission for observer %d",
					i,
				)
			}
			if tc.ballotStatus != observertypes.BallotStatus_BallotInProgress {
				require.Len(t, zk.ObserverKeeper.GetAllBallots(ctx), 0)
				_, found := zk.ObserverKeeper.GetBallotListForHeight(ctx, 0)
				require.False(t, found)
			}
		})
	}
}

// setEmissionsParams sets the emissions params in the store without validation
func setEmissionsParams(t *testing.T, ctx sdk.Context, k emissionskeeper.Keeper, params emissionstypes.Params) {
	store := ctx.KVStore(k.GetStoreKey())
	bz, err := k.GetCodec().Marshal(&params)
	if err != nil {
		require.NoError(t, err)
	}

	store.Set(emissionstypes.KeyPrefix(emissionstypes.ParamsKey), bz)
}
