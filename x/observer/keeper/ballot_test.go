package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"testing"
)

func TestKeeper_GetBallot(t *testing.T) {
	k, ctx := SetupKeeper(t)
	identifier := "0x9ea007f0f60e32d58577a8cf25678942d2b10791c2a34f48e237b76a7e998e4d"
	k.SetBallot(ctx, &types.Ballot{
		BallotIdentifier: identifier,
		VoterList:        nil,
		ObservationType:  0,
		BallotThreshold:  sdk.Dec{},
		BallotStatus:     0,
	})

	k.GetBallot(ctx, identifier)
}
