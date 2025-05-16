package v11_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/observer/keeper"
	v11 "github.com/zeta-chain/node/x/observer/migrations/v11"
	"github.com/zeta-chain/node/x/observer/types"
)

const TestnetStringChainID = "zetachain_7001-1"
const MainnetStringChainID = "zetachain_7000-2"

func SaveBallotsToState(
	t *testing.T,
	ctx sdk.Context,
	k *keeper.Keeper,
	readFile bool,
	startHeight int64,
	endHeight int64,
) {

	if readFile {
		ballots := ImportData(t, k.Codec())
		ballotsList := map[int64][]string{}
		for _, b := range ballots {
			k.SetBallot(ctx, b)
			ballotsList[b.BallotCreationHeight] = append(
				ballotsList[b.BallotCreationHeight],
				b.BallotIdentifier,
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
	t.Run("delete ballots on testnet with creation height 0", func(t *testing.T) {
		useStateExport := false
		bufferedMaturityBlocks := v11.MaturityBlocks + v11.PendingBallotsDeletionBufferBlocks
		currentHeight := int64(0)
		if useStateExport {
			currentHeight = 10309298
		} else {
			currentHeight = 144500
		}
		// Only stale ballots
		startHeight := currentHeight - bufferedMaturityBlocks
		endHeight := int64(0)

		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		SaveBallotsToState(t, ctx, k, useStateExport, startHeight, endHeight)

		ctx = ctx.WithBlockHeight(currentHeight)
		ctx = ctx.WithChainID(TestnetStringChainID)
		err := v9MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Greater(t, len(k.GetAllBallots(ctx)), 0)

		err = v11.MigrateStore(ctx, *k)
		require.NoError(t, err)

		deletionHeight := currentHeight - bufferedMaturityBlocks

		ballots := k.GetAllBallots(ctx)
		for _, ballot := range ballots {
			require.GreaterOrEqual(t, ballot.BallotCreationHeight, deletionHeight)
		}
	})

	t.Run("do not nothing if no ballots are present", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		ballotsBeforeMigrations := k.GetAllBallots(ctx)

		ctx = ctx.WithChainID(TestnetStringChainID)
		err := v11.MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Len(t, k.GetAllBallots(ctx), len(ballotsBeforeMigrations))
	})

	t.Run("do not nothing on mainnet with no stale ballots", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		currentHeight := int64(144500)
		bufferedMaturityBlocks := v11.MaturityBlocks + v11.PendingBallotsDeletionBufferBlocks

		// No stale ballots ( All created ballots are higher than maturity blocks)
		startHeight := currentHeight
		endHeight := currentHeight - bufferedMaturityBlocks

		SaveBallotsToState(t, ctx, k, false, startHeight, endHeight)
		ballotsBeforeMigrations := k.GetAllBallots(ctx)
		ctx = ctx.WithBlockHeight(currentHeight)
		ctx = ctx.WithChainID(MainnetStringChainID)

		err := v11.MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Len(t, k.GetAllBallots(ctx), len(ballotsBeforeMigrations))

		deletionHeight := currentHeight - bufferedMaturityBlocks
		ballots := k.GetAllBallots(ctx)
		for _, ballot := range ballots {
			require.GreaterOrEqual(t, ballot.BallotCreationHeight, deletionHeight)
		}
	})

	t.Run("do nothing on mainnet even if stale ballots are present", func(t *testing.T) {
		useStateExport := false
		bufferedMaturityBlocks := v11.MaturityBlocks + v11.PendingBallotsDeletionBufferBlocks
		currentHeight := int64(0)
		if useStateExport {
			currentHeight = 10309298
		} else {
			currentHeight = 144500
		}
		// Only stale ballots
		startHeight := currentHeight - bufferedMaturityBlocks
		endHeight := int64(0)

		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		SaveBallotsToState(t, ctx, k, useStateExport, startHeight, endHeight)

		ctx = ctx.WithBlockHeight(currentHeight)
		ctx = ctx.WithChainID(MainnetStringChainID)
		err := v9MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Greater(t, len(k.GetAllBallots(ctx)), 0)

		ballotsBeforeMigrations := k.GetAllBallots(ctx)

		err = v11.MigrateStore(ctx, *k)
		require.NoError(t, err)

		require.Len(t, k.GetAllBallots(ctx), len(ballotsBeforeMigrations))
	})

}

func v9MigrateStore(ctx sdk.Context, observerKeeper keeper.Keeper) error {
	currentHeight := ctx.BlockHeight()
	bufferedMaturityBlocks := v11.MaturityBlocks + v11.PendingBallotsDeletionBufferBlocks
	// Maturity blocks is a parameter in the emissions module
	if currentHeight < bufferedMaturityBlocks {
		return nil
	}
	maturedHeight := currentHeight - bufferedMaturityBlocks
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

func ImportData(t *testing.T, cdc codec.JSONCodec) []*types.Ballot {
	file := os.Getenv("STATE_EXPORT_PATH")
	_, genesis, err := genutiltypes.GenesisStateFromGenFile(file)
	require.NoError(t, err)
	appState, err := genutiltypes.GenesisStateFromAppGenesis(genesis)
	require.NoError(t, err)
	importedAppState := types.GetGenesisStateFromAppState(cdc, appState)
	return importedAppState.Ballots
}
