package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"testing"
)

func TestKeeper_GetBallot(t *testing.T) {
	tt := []struct {
		name               string
		identifierInserted string
		identifierQueried  string
		assert             assert.BoolAssertionFunc
	}{
		{"test1", "identifier1", "identifier1", assert.True},
		{"test2", "identifier2", "identifier1", assert.False},
	}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			k, ctx := SetupKeeper(t)
			k.SetBallot(ctx, &types.Ballot{
				Index:            "",
				BallotIdentifier: test.identifierInserted,
				VoterList:        nil,
				Votes:            nil,
				ObservationType:  0,
				BallotThreshold:  sdk.Dec{},
				BallotStatus:     0,
				CreationHeight:   10,
			})
			_, found := k.GetBallot(ctx, test.identifierQueried)
			test.assert(t, found)
		})
	}
}
