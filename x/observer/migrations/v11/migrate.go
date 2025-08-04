package v11

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	pkgchains "github.com/zeta-chain/node/pkg/chains"
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
	zetachain, err := pkgchains.ZetaChainFromCosmosChainID(ctx.ChainID())
	if err != nil {
		// Its fine to return nil here and not try to execute the migration at all if the parsing fails
		ctx.Logger().Error("failed to parse chain ID", "chain_id", ctx.ChainID(), "error", err)
		return err
	}
	if zetachain.ChainId == pkgchains.ZetaChainMainnet.ChainId {
		return nil
	}

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
