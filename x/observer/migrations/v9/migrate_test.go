package v9_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	v9 "github.com/zeta-chain/node/x/observer/migrations/v9"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("delete all matured ballots", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		// Set current height to 1000
		currentHeight := int64(1000)
		// 100 is the maturity blocks parameter defined in emissions module and used by the migrator script
		lastMaturedHeight := currentHeight - v9.MaturityBlocks
		// The first block height is 1 for zeta chain
		firstBlockHeight := int64(1)
		// Account for the first block height to be 1
		numberOfActualBlocks := currentHeight - firstBlockHeight
		blocksWithUnMaturedBallots := numberOfActualBlocks - lastMaturedHeight
		// Use this constant to add ballot to each block
		numberOfBallotsPerBlock := int64(10)

		ctx = ctx.WithBlockHeight(currentHeight)

		for i := firstBlockHeight; i < currentHeight; i++ {
			for j := int64(0); j < numberOfBallotsPerBlock; j++ {
				b := types.Ballot{
					BallotIdentifier:     sample.ZetaIndex(t),
					BallotCreationHeight: i,
				}
				k.AddBallotToList(ctx, b)
				k.SetBallot(ctx, &b)
			}
		}

		allBallots := k.GetAllBallots(ctx)
		require.Equal(t, numberOfActualBlocks*numberOfBallotsPerBlock, int64(len(allBallots)))

		//Act
		err := v9.MigrateStore(ctx, k)
		require.NoError(t, err)

		//Assert
		// We have 10 ballots per block for the last 999 blocks.
		// However, since the maturity blocks are 100, the last 99 blocks will have unmatured ballots.
		// 99*10 = 990
		remainingBallotsAfterMigration := k.GetAllBallots(ctx)
		require.Equal(t, blocksWithUnMaturedBallots*numberOfBallotsPerBlock, int64(len(remainingBallotsAfterMigration)))
		require.Equal(t, int64(990), int64(len(remainingBallotsAfterMigration)))

		// remaining ballots should have creation height greater than last matured height
		for _, b := range remainingBallotsAfterMigration {
			require.Greater(t, b.BallotCreationHeight, lastMaturedHeight)
		}

		// all ballots lists before last matured height should be deleted
		for i := firstBlockHeight; i < lastMaturedHeight; i++ {
			_, found := k.GetBallotListForHeight(ctx, i)
			require.False(t, found)
		}

		// all ballots lists after last matured height should be present
		for i := lastMaturedHeight + 1; i < currentHeight; i++ {
			_, found := k.GetBallotListForHeight(ctx, i)
			require.True(t, found)
		}
	})

	t.Run("do not thing if ballot list for height is not found", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		// Set current height to 1000
		currentHeight := int64(1000)
		// The first block height is 1 for zeta chain
		firstBlockHeight := int64(1)
		// Account for the first block height to be 1
		numberOfActualBlocks := currentHeight - firstBlockHeight
		// Use this constant to add ballot to each block
		numberOfBallotsPerBlock := int64(10)

		ctx = ctx.WithBlockHeight(currentHeight)

		for i := firstBlockHeight; i < currentHeight; i++ {
			for j := int64(0); j < numberOfBallotsPerBlock; j++ {
				b := types.Ballot{
					BallotIdentifier:     sample.ZetaIndex(t),
					BallotCreationHeight: i,
				}
				k.SetBallot(ctx, &b)
			}
		}

		allBallots := k.GetAllBallots(ctx)
		require.Equal(t, numberOfActualBlocks*numberOfBallotsPerBlock, int64(len(allBallots)))

		//Act
		err := v9.MigrateStore(ctx, k)
		require.NoError(t, err)

		//Assert
		allBallotsAfterMigration := k.GetAllBallots(ctx)
		require.Equal(t, numberOfActualBlocks*numberOfBallotsPerBlock, int64(len(allBallotsAfterMigration)))
	})

	t.Run("do not thing if current height is less than maturity blocks", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		// Set current height to 100
		currentHeight := int64(10)
		// The first block height is 1 for zeta chain
		firstBlockHeight := int64(1)
		// Account for the first block height to be 1
		numberOfActualBlocks := currentHeight - firstBlockHeight
		// Use this constant to add ballot to each block
		numberOfBallotsPerBlock := int64(10)

		ctx = ctx.WithBlockHeight(currentHeight)

		for i := firstBlockHeight; i < currentHeight; i++ {
			for j := int64(0); j < numberOfBallotsPerBlock; j++ {
				b := types.Ballot{
					BallotIdentifier:     sample.ZetaIndex(t),
					BallotCreationHeight: i,
				}
				k.AddBallotToList(ctx, b)
				k.SetBallot(ctx, &b)
			}
		}

		allBallots := k.GetAllBallots(ctx)
		require.Equal(t, numberOfActualBlocks*numberOfBallotsPerBlock, int64(len(allBallots)))

		//Act
		err := v9.MigrateStore(ctx, k)
		require.NoError(t, err)

		//Assert
		allBallotsAfterMigration := k.GetAllBallots(ctx)
		require.Equal(t, numberOfActualBlocks*numberOfBallotsPerBlock, int64(len(allBallotsAfterMigration)))
	})
}
