package v11

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

type observerKeeper interface {
	GetBallotListForHeight(ctx sdk.Context, height int64) (val types.BallotListForHeight, found bool)
	DeleteBallot(ctx sdk.Context, index string)
	DeleteBallotListForHeight(ctx sdk.Context, height int64)
	GetAllBallots(ctx sdk.Context) (ballots []*types.Ballot)
}

const (
	PendingBallotsDeletionBufferBlocks = int64(144000)
	MaturityBlocks                     = int64(100)
)

// MigrateStore migrates the x/observer module state from consensus version 10 to version 11.
// The migration deletes all ballots and ballot lists for heights less than the maturity blocks on testnet
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	currentHeight := ctx.BlockHeight()
	// TODO evm: ParseChainID is removed, but since this is old migration probably we can remove it completely
	bufferedMaturityBlocks := MaturityBlocks + PendingBallotsDeletionBufferBlocks
	// Maturity blocks is a parameter in the emissions module
	if currentHeight < bufferedMaturityBlocks {
		return nil
	}

	maturedHeight := currentHeight - bufferedMaturityBlocks
	// We cannot use the ballot list for height 0 as the ballot list was not set for stale ballots on testnet
	ballots := observerKeeper.GetAllBallots(ctx)

	for _, ballot := range ballots {
		if ballot.BallotCreationHeight < maturedHeight {
			observerKeeper.DeleteBallot(ctx, ballot.BallotIdentifier)
		}
	}

	// We can attempt to delete the ballot list for height 0 if it exists
	_, found := observerKeeper.GetBallotListForHeight(ctx, 0)
	if !found {
		return nil
	}

	observerKeeper.DeleteBallotListForHeight(ctx, 0)

	return nil
}
