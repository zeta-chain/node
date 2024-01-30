package emissions_test

import (
	"fmt"
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

func TestDistributeObserverRewards(t *testing.T) {
	k, ctx, sk, zk := keepertest.EmisionKeeper(t)
	observerSet := sample.ObserverSet(10)
	zk.ObserverKeeper.SetObserverSet(ctx, observerSet)
	ballotList := sample.BallotList(10, observerSet.ObserverList)
	ballotIdentifiers := []string{}
	for _, ballot := range ballotList {
		zk.ObserverKeeper.SetBallot(ctx, &ballot)
		ballotIdentifiers = append(ballotIdentifiers, ballot.BallotIdentifier)
	}
	zk.ObserverKeeper.SetBallotList(ctx, &observerTypes.BallotListForHeight{
		Height:           0,
		BallotsIndexList: ballotIdentifiers,
	})
	totalBlockRewards, err := common.GetAzetaDecFromAmountInZeta(emissionstypes.BlockRewardsInZeta)
	rewardCoins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, totalBlockRewards.TruncateInt()))
	assert.NoError(t, err)
	err = sk.BankKeeper.MintCoins(ctx, emissionstypes.ModuleName, rewardCoins)
	assert.NoError(t, err)
	undistributedObserverPoolAddress := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedObserverRewardsPool).GetAddress()
	undistributedTssPoolAddress := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.UndistributedTssRewardsPool).GetAddress()
	feeCollecterAddress := sk.AuthKeeper.GetModuleAccount(ctx, types.FeeCollectorName).GetAddress()
	emissionPool := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.ModuleName).GetAddress()

	for i := 0; i < 20736000; i++ {
		balanceEmissonPoolBeforeBlockDistribution := sk.BankKeeper.GetBalance(ctx, emissionPool, config.BaseDenom).Amount

		emissionsModule.BeginBlocker(ctx, *k)

		feeCollecterBalance := sk.BankKeeper.GetBalance(ctx, feeCollecterAddress, config.BaseDenom).Amount
		observerPoolBalance := sk.BankKeeper.GetBalance(ctx, undistributedObserverPoolAddress, config.BaseDenom).Amount
		tssPoolBalance := sk.BankKeeper.GetBalance(ctx, undistributedTssPoolAddress, config.BaseDenom).Amount
		emissionPoolBalanceAfterBlockDistribution := sk.BankKeeper.GetBalance(ctx, emissionPool, config.BaseDenom).Amount
		rewardsDistributedAndLeftInPool := emissionPoolBalanceAfterBlockDistribution.Sub(rewardCoins.AmountOf(config.BaseDenom))
		assert.True(t, balanceEmissonPoolBeforeBlockDistribution.Sub(rewardsDistributedAndLeftInPool).GTE(sdk.ZeroInt()))
		totalPoolBalance := feeCollecterBalance.Add(observerPoolBalance).Add(tssPoolBalance)
		assert.True(t, rewardCoins.AmountOf(config.BaseDenom).Sub(totalPoolBalance).GTE(sdk.ZeroInt()))
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	}

}

func TestKeeper_CalculateFixedObserverRewards(t *testing.T) {
	SecsInAHour := float64(60 * 60)
	BlockTime := 6.0
	BlocksInAHour := SecsInAHour / BlockTime
	NoOfHours := 2
	StartingBlock := 1940897
	fmt.Println("Proposal Block :", StartingBlock+int(BlocksInAHour*float64(NoOfHours)))
}
