package keeper_test

import (
	"fmt"
	"testing"

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
	}{
		{
			name:             "successfull votes only",
			ballotIdentifier: sample.ZetaIndex(t),
			ballotType:       types.ObservationType_InboundTx,
			voters:           []string{"voter1", "voter2"},
			voteType:         []types.VoteType{types.VoteType_SuccessObservation, types.VoteType_SuccessObservation},
		},

		{
			name:             "failed votes only",
			ballotIdentifier: sample.ZetaIndex(t),
			ballotType:       types.ObservationType_InboundTx,
			voters:           []string{"voter1", "voter2"},
			voteType:         []types.VoteType{types.VoteType_FailureObservation, types.VoteType_FailureObservation},
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
			for _, event := range ctx.EventManager().Events() {
				for _, attr := range event.Attributes {
					if attr.Key == "ballot_identifier" {
						require.Equal(t, tc.ballotIdentifier, RemoveQuotes(attr.Value))
					}
					if attr.Key == "ballot_type" {
						require.Equal(t, tc.ballotType.String(), RemoveQuotes(attr.Value))
					}
					if attr.Key == "voters" {
						expectedString := ""
						for _, voter := range ballot.GenerateVoterList() {
							st := fmt.Sprintf(
								"{\"voter_address\":\"%s\",\"vote_type\":\"%s\"}",
								voter.VoterAddress,
								voter.VoteType,
							)
							expectedString += st
							expectedString += ","
						}
						expectedString = expectedString[:len(expectedString)-1]
						require.Equal(t, expectedString, RemoveQuotes(attr.Value))
					}
				}
			}

		})
	}
}

func RemoveQuotes(s string) string {
	return s[1 : len(s)-1]
}
