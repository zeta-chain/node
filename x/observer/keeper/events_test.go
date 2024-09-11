package keeper_test

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestEmitEventBallotDeleted(t *testing.T) {
	tt := []struct {
		name             string
		ballotIdentifier string
		ballotType       types.ObservationType
		voters           []string
		voteType         []types.VoteType
		expectErr        bool
	}{
		{
			name:             "successful votes only",
			ballotIdentifier: sample.ZetaIndex(t),
			ballotType:       types.ObservationType_InboundTx,
			voters:           []string{"voter1", "voter2"},
			voteType:         []types.VoteType{types.VoteType_SuccessObservation, types.VoteType_SuccessObservation},
			expectErr:        false,
		},

		{
			name:             "failed votes only",
			ballotIdentifier: sample.ZetaIndex(t),
			ballotType:       types.ObservationType_InboundTx,
			voters:           []string{"voter1", "voter2"},
			voteType:         []types.VoteType{types.VoteType_FailureObservation, types.VoteType_FailureObservation},
			expectErr:        false,
		},

		{
			name:             "invalid voter list",
			ballotIdentifier: sample.ZetaIndex(t),
			ballotType:       types.ObservationType_InboundTx,
			voters:           []string{"voter1", "voter2"},
			voteType:         []types.VoteType{types.VoteType_FailureObservation},
			expectErr:        true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, ctx, _, _ := keepertest.ObserverKeeper(t)
			ballot := types.Ballot{
				BallotIdentifier: tc.ballotIdentifier,
				ObservationType:  tc.ballotType,
				VoterList:        tc.voters,
				Votes:            tc.voteType,
			}
			keeper.EmitEventBallotDeleted(ctx, ballot)
			checkEventAttributes(t, ctx, ballot, tc.expectErr)
		})
	}
}

func checkEventAttributes(t *testing.T, ctx sdk.Context, ballot types.Ballot, expectErr bool) {
	for _, event := range ctx.EventManager().Events() {
		for _, attr := range event.Attributes {
			if attr.Key == "ballot_identifier" {
				require.Equal(t, ballot.BallotIdentifier, RemoveQuotes(attr.Value))
			}
			if attr.Key == "ballot_type" {
				require.Equal(t, ballot.ObservationType.String(), RemoveQuotes(attr.Value))
			}
			if attr.Key == "voters" {
				expectedString := ""
				list, err := ballot.GenerateVoterList()
				if !expectErr {
					require.NoError(t, err)
					var voterStrings []string
					for _, voter := range list {
						voterStrings = append(voterStrings, fmt.Sprintf(
							"{\"voter_address\":\"%s\",\"vote_type\":\"%s\"}",
							voter.VoterAddress,
							voter.VoteType,
						))
					}
					expectedString = strings.Join(voterStrings, ",")
				} else {
					require.ErrorIs(t, err, types.ErrInvalidVoterList)
				}
				require.Equal(t, expectedString, RemoveQuotes(attr.Value))
			}
		}
	}
}

func RemoveQuotes(s string) string {
	return s[1 : len(s)-1]
}
