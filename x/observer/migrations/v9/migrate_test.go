package v9_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/testdata"
	observerkeeper "github.com/zeta-chain/node/x/observer/keeper"
	v9 "github.com/zeta-chain/node/x/observer/migrations/v9"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMigrateStore(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	SetupMainnetData(t, k, ctx)
	ctx = ctx.WithBlockHeight(5073825)
	err := v9.MigrateStore(ctx, k)
	require.NoError(t, err)

	start := time.Now()
	ballots := k.GetAllBallots(ctx)
	elapsed := time.Since(start)

	// Elapsed should be less than 1 second
	require.Less(t, elapsed.Seconds(), 1)
	require.Len(t, ballots, 0)
}

func SetupMainnetData(t *testing.T, k *observerkeeper.Keeper, ctx sdk.Context) {
	b := testdata.ReadMainnetBallots(t)
	ballotsList := map[int64][]string{}

	for _, ballot := range b {
		k.SetBallot(ctx, ballot)
		ballotsList[ballot.BallotCreationHeight] = append(
			ballotsList[ballot.BallotCreationHeight],
			ballot.BallotIdentifier,
		)
	}

	for height, ballotList := range ballotsList {
		k.SetBallotList(ctx, &types.BallotListForHeight{
			Height:           height,
			BallotsIndexList: ballotList,
		})
	}
}
