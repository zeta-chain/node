package emissions_test

import (
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
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
	t.Run("no distribution happens if emissions module account is empty", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		emissionsModule.BeginBlocker(ctx, *k)
	})
}
func TestObserverRewards(t *testing.T) {
	// setup the test
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
	assert.NoError(t, err)
	// Fund the emission pool to start the emission process
	err = sk.BankKeeper.MintCoins(ctx, emissionstypes.ModuleName, totalRewardCoins)
	assert.NoError(t, err)

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
	assert.True(t, blockRewards.TruncateInt().GT(distributedRewards))

	for i := 0; i < 100; i++ {
		emissionPoolBeforeBlockDistribution := sk.BankKeeper.GetBalance(ctx, emissionPool, config.BaseDenom).Amount
		// produce a block
		emissionsModule.BeginBlocker(ctx, *k)

		// Assert distribution amount
		emissionPoolBalanceAfterBlockDistribution := sk.BankKeeper.GetBalance(ctx, emissionPool, config.BaseDenom).Amount
		assert.True(t, emissionPoolBeforeBlockDistribution.Sub(emissionPoolBalanceAfterBlockDistribution).Equal(distributedRewards))

		// totalDistributedTillCurrentBlock is the net amount of rewards distributed till the current block, this works in a unit test as the fees are not being collected by validators
		totalDistributedTillCurrentBlock := sk.BankKeeper.GetBalance(ctx, feeCollecterAddress, config.BaseDenom).Amount.
			Add(sk.BankKeeper.GetBalance(ctx, undistributedObserverPoolAddress, config.BaseDenom).Amount).
			Add(sk.BankKeeper.GetBalance(ctx, undistributedTssPoolAddress, config.BaseDenom).Amount)
		// Assert we are always under the max limit of block rewards
		assert.True(t, totalRewardCoins.AmountOf(config.BaseDenom).
			Sub(totalDistributedTillCurrentBlock).GTE(sdk.ZeroInt()))

		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	}

	// We can simplify the calculation as the rewards are distributed equally among all the observers
	rewardPerUnit := observerRewardsForABlock.Quo(sdk.NewInt(int64(len(ballotList) * len(observerSet.ObserverList))))
	emissionAmount := rewardPerUnit.Mul(sdk.NewInt(int64(len(ballotList))))

	for _, observer := range observerSet.ObserverList {
		observerEmission, found := k.GetWithdrawableEmission(ctx, observer)
		require.True(t, found)
		require.Equal(t, emissionAmount, observerEmission.Amount)
	}
}

func TestValidatorRewards(t *testing.T) {
	k, ctx, sk, zk := keepertest.EmissionsKeeper(t)
	k.SetParams(ctx, emissionstypes.DefaultParams())
	observerSet := make([]string, 10)
	for i := 0; i < 10; i++ {
		validator := sample.Validator(t, rand.New(rand.NewSource(int64(i))))
		validator.Tokens = sample.IntInRange(100000, 100000000)
		sk.StakingKeeper.SetValidator(ctx, validator)
		observerSet[i] = validator.OperatorAddress
	}
	zk.ObserverKeeper.SetObserverSet(ctx, observerTypes.ObserverSet{
		ObserverList: observerSet,
	})

	// Total block rewards is the fixed amount of rewards that are distributed
	totalBlockRewards, err := common.GetAzetaDecFromAmountInZeta(emissionstypes.BlockRewardsInZeta)
	rewardCoins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, totalBlockRewards.TruncateInt()))
	assert.NoError(t, err)
	// Fund the emission pool to start the emission process
	err = sk.BankKeeper.MintCoins(ctx, emissionstypes.ModuleName, rewardCoins)
	assert.NoError(t, err)

	// Setup module accounts for emission pools
	_ = sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedObserverRewardsPool).GetAddress()
	_ = sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedTssRewardsPool).GetAddress()
	feeCollectorAddress := sk.AuthKeeper.GetModuleAccount(ctx, types.FeeCollectorName).GetAddress()
	_ = sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.ModuleName).GetAddress()
	blockRewards := emissionstypes.BlockReward
	// Produce blocks and distribute rewards
	validatorRewards := sdk.MustNewDecFromStr(k.GetParams(ctx).ValidatorEmissionPercentage).Mul(blockRewards).TruncateInt()
	for i := 0; i < 100; i++ {
		// produce a block
		emissionsModule.BeginBlocker(ctx, *k)
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	}
	feeCollectorBalance := sk.BankKeeper.GetBalance(ctx, feeCollectorAddress, config.BaseDenom).Amount
	assert.Equal(t, feeCollectorBalance, validatorRewards.Mul(sdk.NewInt(100)))
}
