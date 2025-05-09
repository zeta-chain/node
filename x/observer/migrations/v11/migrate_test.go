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

func SaveBallotsToState(t *testing.T, ctx sdk.Context, k *keeper.Keeper, readFile bool, startHeight int64, endHeight int64) {
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
			if height > 0 {
				k.SetBallotList(ctx, &types.BallotListForHeight{
					Height:           height,
					BallotsIndexList: ballotList,
				})
			}
		}
	} else {

		ballotPerBlock := 10

		for i := startHeight; i >= endHeight; i-- {
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
			// The Ballot list was not set for older ballots on testnet
			if i > 0 {
				k.SetBallotList(ctx, &types.BallotListForHeight{
					Height:           i,
					BallotsIndexList: ballotList,
				})
			}
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
			currentHeight = 144500
			testnetBallotsAfterV9Migration = 10
		}
		// Only stale ballots
		startHeight := currentHeight - v11.BufferedMaturityBlocks - 1
		endHeight := int64(0)

		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		SaveBallotsToState(t, ctx, k, useStateExport, startHeight, endHeight)

		ctx = ctx.WithBlockHeight(currentHeight)
		err := v9MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Len(t, k.GetAllBallots(ctx), testnetBallotsAfterV9Migration)

		err = v11.MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Len(t, k.GetAllBallots(ctx), 0)
	})

	t.Run("do not nothing if no ballots are present", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		ballotsBeforeMigrations := k.GetAllBallots(ctx)

		err := v11.MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Len(t, k.GetAllBallots(ctx), len(ballotsBeforeMigrations))
	})

	t.Run("do not nothing if no stale ballots are present(Simulate Mainnet)", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		currentHeight := int64(144500)

		// No stale ballots
		startHeight := currentHeight
		endHeight := currentHeight - v11.BufferedMaturityBlocks

		SaveBallotsToState(t, ctx, k, false, startHeight, endHeight)
		ballotsBeforeMigrations := k.GetAllBallots(ctx)

		err := v11.MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Len(t, k.GetAllBallots(ctx), len(ballotsBeforeMigrations))
	})

}

func v9MigrateStore(ctx sdk.Context, observerKeeper keeper.Keeper) error {
	currentHeight := ctx.BlockHeight()
	// Maturity blocks is a parameter in the emissions module
	if currentHeight < v11.BufferedMaturityBlocks {
		return nil
	}
	maturedHeight := currentHeight - v11.BufferedMaturityBlocks
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
