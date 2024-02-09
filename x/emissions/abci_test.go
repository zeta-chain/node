package emissions_test

import (
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	emissionsModule "github.com/zeta-chain/zetacore/x/emissions"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestBeginBlocker(t *testing.T) {
	// setup the test
	k, ctx, sk, zk := keepertest.EmisionKeeper(t)
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
	rewardCoins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, totalBlockRewards.TruncateInt()))
	assert.NoError(t, err)
	// Fund the emission pool to start the emission process
	err = sk.BankKeeper.MintCoins(ctx, emissionstypes.ModuleName, rewardCoins)
	assert.NoError(t, err)

	// Setup module accounts for emission pools
	undistributedObserverPoolAddress := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedObserverRewardsPool).GetAddress()
	undistributedTssPoolAddress := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedTssRewardsPool).GetAddress()
	feeCollecterAddress := sk.AuthKeeper.GetModuleAccount(ctx, types.FeeCollectorName).GetAddress()
	emissionPool := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.ModuleName).GetAddress()
	totalDistributedTillLastBlock := sdk.ZeroInt()
	blockRewards := emissionstypes.BlockReward
	for i := 0; i < 100; i++ {
		balanceEmissionPoolBeforeBlockDistribution := sk.BankKeeper.GetBalance(ctx, emissionPool, config.BaseDenom).Amount
		// produce a block
		emissionsModule.BeginBlocker(ctx, *k)

		feeCollecterBalance := sk.BankKeeper.GetBalance(ctx, feeCollecterAddress, config.BaseDenom).Amount
		observerPoolBalance := sk.BankKeeper.GetBalance(ctx, undistributedObserverPoolAddress, config.BaseDenom).Amount
		tssPoolBalance := sk.BankKeeper.GetBalance(ctx, undistributedTssPoolAddress, config.BaseDenom).Amount
		emissionPoolBalanceAfterBlockDistribution := sk.BankKeeper.GetBalance(ctx, emissionPool, config.BaseDenom).Amount

		// Assert the rewards pool has enough balance to distribute rewards
		rewardsDistributedAndLeftInPool := emissionPoolBalanceAfterBlockDistribution.Sub(rewardCoins.AmountOf(config.BaseDenom))
		assert.True(t, balanceEmissionPoolBeforeBlockDistribution.Sub(rewardsDistributedAndLeftInPool).GTE(sdk.ZeroInt()))

		// totalDistributedTillCurrentBlock is the net amount of rewards distributed till the current block , this works in a unit test as the fees are not being collected by validators
		totalDistributedTillCurrentBlock := feeCollecterBalance.Add(observerPoolBalance).Add(tssPoolBalance)

		// Assert that a maximum of value of block rewards is distributed in each block
		assert.True(t, totalDistributedTillCurrentBlock.Sub(totalDistributedTillLastBlock).LTE(blockRewards.TruncateInt()))

		// Assert we are always under the max limit of block rewards
		assert.True(t, rewardCoins.AmountOf(config.BaseDenom).Sub(totalDistributedTillCurrentBlock).GTE(sdk.ZeroInt()))

		totalDistributedTillLastBlock = totalDistributedTillCurrentBlock
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	}
}

func TestDistributeValidatorRewards(t *testing.T) {
	k, ctx, sk, zk := keepertest.EmisionKeeper(t)
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
	_ = sk.StakingKeeper.GetAllValidators(ctx)

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
	feeCollecterAddress := sk.AuthKeeper.GetModuleAccount(ctx, types.FeeCollectorName).GetAddress()
	_ = sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.ModuleName).GetAddress()
	blockRewards := emissionstypes.BlockReward
	// Produce blocks and distribute rewards
	for i := 0; i < 100; i++ {
		validatorRewards := sdk.MustNewDecFromStr(k.GetParams(ctx).ValidatorEmissionPercentage).Mul(blockRewards).TruncateInt()
		// produce a block
		emissionsModule.BeginBlocker(ctx, *k)
		feeCollecterBalance := sk.BankKeeper.GetBalance(ctx, feeCollecterAddress, config.BaseDenom).Amount
		assert.Equal(t, validatorRewards, feeCollecterBalance)

		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	}

}
