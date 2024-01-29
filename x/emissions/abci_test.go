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
	emissionPool := sk.AuthKeeper.GetModuleAccount(ctx, emissionstypes.ModuleName).GetAddress()
	feeCollecterAddress := sk.AuthKeeper.GetModuleAccount(ctx, types.FeeCollectorName).GetAddress()
	for i := 0; i < 90; i++ {
		balanceEmissonPool := sk.BankKeeper.GetBalance(ctx, emissionPool, config.BaseDenom)
		fmt.Println("Emission Pool Balance: ", balanceEmissonPool.String())

		emissionsModule.BeginBlocker(ctx, *k)

		feeCoolecterBalance := sk.BankKeeper.GetBalance(ctx, feeCollecterAddress, config.BaseDenom)
		fmt.Println("Fee Collected : ", feeCoolecterBalance.String())
		balances := sk.BankKeeper.GetBalance(ctx, undistributedObserverPoolAddress, config.BaseDenom)
		fmt.Println("Balance Observer : ", balances.String())
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	}

}
