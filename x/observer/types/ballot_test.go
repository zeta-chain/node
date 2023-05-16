package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVoter_IsBallotFinalized(t *testing.T) {
	type votes struct {
		address string
		vote    VoteType
	}

	tt := []struct {
		name        string
		threshold   sdk.Dec
		voterList   []string
		votes       []votes
		finalVotes  []VoteType
		finalStatus BallotStatus
		isFinalized bool
	}{
		{
			name:      "All success",
			threshold: sdk.MustNewDecFromStr("0.66"),
			voterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer2", VoteType_SuccessObservation},
				{"Observer3", VoteType_SuccessObservation},
				{"Observer4", VoteType_SuccessObservation},
			},
			finalVotes:  []VoteType{VoteType_SuccessObservation, VoteType_SuccessObservation, VoteType_SuccessObservation, VoteType_SuccessObservation},
			finalStatus: BallotStatus_BallotFinalized_SuccessObservation,
			isFinalized: true,
		},
		{
			name:      "Multiple votes by a observer , Ballot success",
			threshold: sdk.MustNewDecFromStr("0.66"),
			voterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer1", VoteType_FailureObservation},
				{"Observer2", VoteType_SuccessObservation},
				{"Observer3", VoteType_SuccessObservation},
				{"Observer4", VoteType_SuccessObservation},
			},
			finalVotes:  []VoteType{VoteType_SuccessObservation, VoteType_SuccessObservation, VoteType_SuccessObservation, VoteType_SuccessObservation},
			finalStatus: BallotStatus_BallotFinalized_SuccessObservation,
			isFinalized: true,
		},
		{
			name:      "Multiple votes by a observer , Ballot in progress",
			threshold: sdk.MustNewDecFromStr("0.66"),
			voterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer1", VoteType_FailureObservation},
				{"Observer1", VoteType_SuccessObservation},
				{"Observer1", VoteType_SuccessObservation},
				{"Observer1", VoteType_SuccessObservation},
			},
			finalVotes:  []VoteType{VoteType_SuccessObservation, VoteType_NotYetVoted, VoteType_NotYetVoted, VoteType_NotYetVoted},
			finalStatus: BallotStatus_BallotInProgress,
			isFinalized: false,
		},

		{
			name:      "Two observers ",
			threshold: sdk.MustNewDecFromStr("1.00"),
			voterList: []string{"Observer1", "Observer2"},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer2", VoteType_SuccessObservation},
			},
			finalVotes:  []VoteType{VoteType_SuccessObservation, VoteType_SuccessObservation},
			finalStatus: BallotStatus_BallotFinalized_SuccessObservation,
			isFinalized: true,
		},
	}
	for _, test := range tt {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ballot := Ballot{
				BallotIdentifier: "identifier",
				VoterList:        test.voterList,
				Votes:            CreateVotes(len(test.voterList)),
				ObservationType:  ObservationType_InBoundTx,
				BallotThreshold:  test.threshold,
				BallotStatus:     BallotStatus_BallotInProgress,
			}
			for _, vote := range test.votes {
				ballot, _ = ballot.AddVote(vote.address, vote.vote)
			}

			finalBallot, isFinalized := ballot.IsBallotFinalized()
			assert.Equal(t, test.finalStatus, finalBallot.BallotStatus)
			assert.Equal(t, test.finalVotes, finalBallot.Votes)
			assert.Equal(t, test.isFinalized, isFinalized)
		})
	}
}
