package v11_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/observer/keeper"
	v11 "github.com/zeta-chain/node/x/observer/migrations/v11"
	"github.com/zeta-chain/node/x/observer/types"
)

const MaturityBlocks = int64(100)

func SaveBallotsToState(t *testing.T, ctx sdk.Context, k *keeper.Keeper, currentHeight int64, readFile bool) {
	type Ballots struct {
		Ballots []types.Ballot `json:"Ballots"`
	}

	b := Ballots{}

	if readFile {
		file := os.Getenv("STATE_EXPORT_PATH")
		data, err := os.ReadFile(file)
		require.NoError(t, err)
		err = json.Unmarshal(data, &b)
		require.NoError(t, err)
		ballotsList := map[int64][]string{}
		for _, ballot := range b.Ballots {
			k.SetBallot(ctx, &ballot)
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
	} else {
		startHeight := currentHeight - MaturityBlocks - 1
		ballotPerBlock := 10

		for i := startHeight; i >= 0; i-- {
			ballotList := make([]string, ballotPerBlock)
			for j := 0; j < ballotPerBlock; j++ {
				index := fmt.Sprintf("ballot-%d-%d", i, j)
				ballot := &types.Ballot{
					Index:                index,
					BallotIdentifier:     index,
					BallotCreationHeight: i,
				}
				k.SetBallot(ctx, ballot)
				ballotList[j] = index
			}
			k.SetBallotList(ctx, &types.BallotListForHeight{
				Height:           i,
				BallotsIndexList: ballotList,
			})
		}
	}

}

func Test_MigrateStore(t *testing.T) {
	t.Run("delete ballots with creation height 0", func(t *testing.T) {
		useStateExport := false

		currentHeight := int64(0)
		testnetBallotsAfterV9Migration := 0
		if useStateExport {
			currentHeight = 10232850
			testnetBallotsAfterV9Migration = 270431
		} else {
			currentHeight = 500
			testnetBallotsAfterV9Migration = 10
		}

		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		SaveBallotsToState(t, ctx, k, currentHeight, useStateExport)

		ctx = ctx.WithBlockHeight(currentHeight)
		err := v9MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Len(t, k.GetAllBallots(ctx), testnetBallotsAfterV9Migration)

		err = v11.MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Len(t, k.GetAllBallots(ctx), 0)
	})

	t.Run("do not nothing if no ballots are present at height 0", func(t *testing.T) {
		useStateExport := false

		currentHeight := int64(0)
		if useStateExport {
			currentHeight = 10232850
		} else {
			currentHeight = 500
		}

		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		SaveBallotsToState(t, ctx, k, currentHeight, useStateExport)
		ballotList, found := k.GetBallotListForHeight(ctx, 0)
		require.True(t, found)
		for _, ballotIndex := range ballotList.BallotsIndexList {
			k.DeleteBallot(ctx, ballotIndex)
		}

		ballotsBeforeMigrations := k.GetAllBallots(ctx)

		err := v11.MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Len(t, k.GetAllBallots(ctx), len(ballotsBeforeMigrations))
	})

}

func v9MigrateStore(ctx sdk.Context, observerKeeper keeper.Keeper) error {
	currentHeight := ctx.BlockHeight()
	// Maturity blocks is a parameter in the emissions module
	if currentHeight < MaturityBlocks {
		return nil
	}
	maturedHeight := currentHeight - MaturityBlocks
	for i := maturedHeight; i > 0; i-- {
		ballotList, found := observerKeeper.GetBallotListForHeight(ctx, i)
		if !found {
			continue
		}
		for _, ballotIndex := range ballotList.BallotsIndexList {
			observerKeeper.DeleteBallot(ctx, ballotIndex)
		}
		observerKeeper.DeleteBallotListForHeight(ctx, i)
	}
	return nil
}
