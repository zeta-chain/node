package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestBallot_AddVote(t *testing.T) {
	type votes struct {
		address string
		vote    VoteType
	}

	tt := []struct {
		name        string
		threshold   sdkmath.LegacyDec
		voterList   []string
		votes       []votes
		finalVotes  []VoteType
		finalStatus BallotStatus
		isFinalized bool
		wantErr     bool
	}{
		{
			name:      "All success",
			threshold: sdkmath.LegacyMustNewDecFromStr("0.66"),
			voterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer2", VoteType_SuccessObservation},
				{"Observer3", VoteType_SuccessObservation},
				{"Observer4", VoteType_SuccessObservation},
			},
			finalVotes: []VoteType{
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
			},
			finalStatus: BallotStatus_BallotFinalized_SuccessObservation,
			isFinalized: true,
		},
		{
			name:      "Multiple votes by a observer , Ballot success",
			threshold: sdkmath.LegacyMustNewDecFromStr("0.66"),
			voterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer1", VoteType_FailureObservation},
				{"Observer2", VoteType_SuccessObservation},
				{"Observer3", VoteType_SuccessObservation},
				{"Observer4", VoteType_SuccessObservation},
			},
			finalVotes: []VoteType{
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
			},
			finalStatus: BallotStatus_BallotFinalized_SuccessObservation,
			isFinalized: true,
		},
		{
			name:      "Multiple votes by a observer , Ballot in progress",
			threshold: sdkmath.LegacyMustNewDecFromStr("0.66"),
			voterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer1", VoteType_FailureObservation},
				{"Observer1", VoteType_SuccessObservation},
				{"Observer1", VoteType_SuccessObservation},
				{"Observer1", VoteType_SuccessObservation},
			},
			finalVotes: []VoteType{
				VoteType_SuccessObservation,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
			},
			finalStatus: BallotStatus_BallotInProgress,
			isFinalized: false,
		},
		{
			name:      "Ballot finalized at threshold",
			threshold: sdkmath.LegacyMustNewDecFromStr("0.66"),
			voterList: []string{
				"Observer1",
				"Observer2",
				"Observer3",
				"Observer4",
				"Observer5",
				"Observer6",
				"Observer7",
				"Observer8",
				"Observer9",
				"Observer10",
				"Observer11",
				"Observer12",
			},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer2", VoteType_SuccessObservation},
				{"Observer3", VoteType_SuccessObservation},
				{"Observer4", VoteType_SuccessObservation},
				{"Observer5", VoteType_SuccessObservation},
				{"Observer6", VoteType_SuccessObservation},
				{"Observer7", VoteType_SuccessObservation},
				{"Observer8", VoteType_SuccessObservation},
				{"Observer9", VoteType_NotYetVoted},
				{"Observer10", VoteType_NotYetVoted},
				{"Observer11", VoteType_NotYetVoted},
				{"Observer12", VoteType_NotYetVoted},
			},
			finalVotes: []VoteType{VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
			},
			finalStatus: BallotStatus_BallotFinalized_SuccessObservation,
			isFinalized: true,
		},
		{
			name:      "Ballot finalized at threshold but more votes added after",
			threshold: sdkmath.LegacyMustNewDecFromStr("0.66"),
			voterList: []string{
				"Observer1",
				"Observer2",
				"Observer3",
				"Observer4",
				"Observer5",
				"Observer6",
				"Observer7",
				"Observer8",
				"Observer9",
				"Observer10",
				"Observer11",
				"Observer12",
			},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer2", VoteType_SuccessObservation},
				{"Observer3", VoteType_SuccessObservation},
				{"Observer4", VoteType_SuccessObservation},
				{"Observer5", VoteType_SuccessObservation},
				{"Observer6", VoteType_SuccessObservation},
				{"Observer7", VoteType_SuccessObservation},
				{"Observer8", VoteType_SuccessObservation},
				{"Observer9", VoteType_SuccessObservation},
				{"Observer10", VoteType_SuccessObservation},
				{"Observer11", VoteType_SuccessObservation},
				{"Observer12", VoteType_SuccessObservation},
			},
			finalVotes: []VoteType{VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
			},
			finalStatus: BallotStatus_BallotFinalized_SuccessObservation,
			isFinalized: true,
		},
		{
			name:      "Two observers",
			threshold: sdkmath.LegacyMustNewDecFromStr("1.00"),
			voterList: []string{"Observer1", "Observer2"},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer2", VoteType_SuccessObservation},
			},
			finalVotes:  []VoteType{VoteType_SuccessObservation, VoteType_SuccessObservation},
			finalStatus: BallotStatus_BallotFinalized_SuccessObservation,
			isFinalized: true,
		},
		{
			name:      "Low threshold 1 always fails as Failure is checked first",
			threshold: sdkmath.LegacyMustNewDecFromStr("0.01"),
			voterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer2", VoteType_SuccessObservation},
				{"Observer3", VoteType_SuccessObservation},
				{"Observer4", VoteType_FailureObservation},
			},
			finalVotes: []VoteType{
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_FailureObservation,
			},
			finalStatus: BallotStatus_BallotFinalized_FailureObservation,
			isFinalized: true,
		},
		{
			name:      "Low threshold 2 always fails as Failure is checked first",
			threshold: sdkmath.LegacyMustNewDecFromStr("0.01"),
			voterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer2", VoteType_FailureObservation},
				{"Observer3", VoteType_SuccessObservation},
				{"Observer4", VoteType_SuccessObservation},
			},
			finalVotes: []VoteType{
				VoteType_SuccessObservation,
				VoteType_FailureObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
			},
			finalStatus: BallotStatus_BallotFinalized_FailureObservation,
			isFinalized: true,
		},
		{
			name:      "100 percent threshold cannot finalze with less than 100 percent votes",
			threshold: sdkmath.LegacyMustNewDecFromStr("1"),
			voterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []votes{
				{"Observer1", VoteType_SuccessObservation},
				{"Observer2", VoteType_FailureObservation},
				{"Observer3", VoteType_SuccessObservation},
				{"Observer4", VoteType_SuccessObservation},
			},
			finalVotes: []VoteType{
				VoteType_SuccessObservation,
				VoteType_FailureObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
			},
			finalStatus: BallotStatus_BallotInProgress,
			isFinalized: false,
		},
		{
			name:      "Voter not in voter list",
			threshold: sdkmath.LegacyMustNewDecFromStr("0.66"),
			voterList: []string{},
			votes: []votes{
				{"Observer5", VoteType_SuccessObservation},
			},
			wantErr:     true,
			finalVotes:  []VoteType{},
			finalStatus: BallotStatus_BallotInProgress,
			isFinalized: false,
		},
	}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			ballot := Ballot{
				Index:            "index",
				BallotIdentifier: "identifier",
				VoterList:        test.voterList,
				Votes:            CreateVotes(len(test.voterList)),
				ObservationType:  ObservationType_InboundTx,
				BallotThreshold:  test.threshold,
				BallotStatus:     BallotStatus_BallotInProgress,
			}
			for _, vote := range test.votes {
				b, err := ballot.AddVote(vote.address, vote.vote)
				if test.wantErr {
					require.Error(t, err)
				}
				ballot = b
			}

			finalBallot, isFinalized := ballot.IsFinalizingVote()
			require.Equal(t, test.finalStatus, finalBallot.BallotStatus)
			require.Equal(t, test.finalVotes, finalBallot.Votes)
			require.Equal(t, test.isFinalized, isFinalized)
		})
	}
}

func TestBallot_IsFinalizingVote(t *testing.T) {
	tt := []struct {
		name            string
		BallotThreshold sdkmath.LegacyDec
		Votes           []VoteType
		finalizingVote  int
		finalStatus     BallotStatus
	}{
		{
			name:            "finalized to success",
			BallotThreshold: sdkmath.LegacyMustNewDecFromStr("0.66"),
			Votes: []VoteType{
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
			},
			finalizingVote: 7,
			finalStatus:    BallotStatus_BallotFinalized_SuccessObservation,
		},
		{
			name:            "finalized to failure",
			BallotThreshold: sdkmath.LegacyMustNewDecFromStr("0.66"),
			Votes: []VoteType{
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
			},
			finalizingVote: 7,
			finalStatus:    BallotStatus_BallotFinalized_FailureObservation,
		},
		{
			name:            "low threshold finalized early to success",
			BallotThreshold: sdkmath.LegacyMustNewDecFromStr("0.01"),
			Votes: []VoteType{
				VoteType_SuccessObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
				VoteType_NotYetVoted,
			},
			finalizingVote: 0,
			finalStatus:    BallotStatus_BallotFinalized_SuccessObservation,
		},
		{
			name:            "100 percent threshold cannot finalize with less than 100 percent votes",
			BallotThreshold: sdkmath.LegacyMustNewDecFromStr("1"),
			Votes: []VoteType{
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_SuccessObservation,
			},
			finalizingVote: 0,
			finalStatus:    BallotStatus_BallotInProgress,
		},
		{
			name:            "100 percent threshold can finalize with 100 percent votes",
			BallotThreshold: sdkmath.LegacyMustNewDecFromStr("1"),
			Votes: []VoteType{
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
			},
			finalizingVote: 11,
			finalStatus:    BallotStatus_BallotFinalized_FailureObservation,
		},
	}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {

			ballot := Ballot{
				BallotStatus:    BallotStatus_BallotInProgress,
				BallotThreshold: test.BallotThreshold,
				VoterList:       make([]string, len(test.Votes)),
			}
			isFinalizingVote := false
			for index, vote := range test.Votes {
				ballot.Votes = append(ballot.Votes, vote)
				ballot, isFinalizingVote = ballot.IsFinalizingVote()
				if isFinalizingVote {
					require.Equal(t, test.finalizingVote, index)
				}
			}
			require.Equal(t, test.finalStatus, ballot.BallotStatus)
		})
	}
}

func Test_BuildRewardsDistribution(t *testing.T) {
	tt := []struct {
		name        string
		ballotList  []Ballot
		expectedMap map[string]int64
	}{
		{
			name: "all success votes",
			ballotList: []Ballot{{
				VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
				Votes: []VoteType{
					VoteType_SuccessObservation,
					VoteType_SuccessObservation,
					VoteType_SuccessObservation,
					VoteType_SuccessObservation,
				},
				BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
			}},
			expectedMap: map[string]int64{
				"Observer1": 1,
				"Observer2": 1,
				"Observer3": 1,
				"Observer4": 1,
			},
		},
		{
			name: "all success votes 3 ballots",
			ballotList: []Ballot{{
				VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
				Votes: []VoteType{
					VoteType_SuccessObservation,
					VoteType_SuccessObservation,
					VoteType_SuccessObservation,
					VoteType_SuccessObservation,
				},
				BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
			},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
			},
			expectedMap: map[string]int64{
				"Observer1": 3,
				"Observer2": 3,
				"Observer3": 3,
				"Observer4": 3,
			},
		},
		{
			name: "mixed votes with some NotVoted - 3 ballots",
			ballotList: []Ballot{
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
						VoteType_NotYetVoted,
						VoteType_FailureObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_SuccessObservation,
						VoteType_NotYetVoted,
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_NotYetVoted,
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
			},
			expectedMap: map[string]int64{
				"Observer1": 1,
				"Observer2": 1,
				"Observer3": 1,
				"Observer4": 1,
			},
		},
		{
			name: "all failure ballots with mixed votes - 3 ballots",
			ballotList: []Ballot{
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3"},
					Votes: []VoteType{
						VoteType_FailureObservation,
						VoteType_FailureObservation,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_FailureObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3"},
					Votes: []VoteType{
						VoteType_FailureObservation,
						VoteType_NotYetVoted,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_FailureObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3"},
					Votes: []VoteType{
						VoteType_SuccessObservation,
						VoteType_FailureObservation,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_FailureObservation,
				},
			},
			expectedMap: map[string]int64{
				"Observer1": 1,
				"Observer2": 1,
				"Observer3": -3,
			},
		},
		{
			name: "mixed ballot outcomes with varied votes - 3 ballots",
			ballotList: []Ballot{
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
						VoteType_FailureObservation,
						VoteType_NotYetVoted,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_FailureObservation,
						VoteType_FailureObservation,
						VoteType_SuccessObservation,
						VoteType_FailureObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_FailureObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_SuccessObservation,
						VoteType_NotYetVoted,
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
			},
			expectedMap: map[string]int64{
				"Observer1": 3,
				"Observer2": 1,
				"Observer3": -1,
				"Observer4": 1,
			},
		},
		{
			name: "heavy NotVoted scenario - 3 ballots",
			ballotList: []Ballot{
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_NotYetVoted,
						VoteType_NotYetVoted,
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_SuccessObservation,
						VoteType_NotYetVoted,
						VoteType_NotYetVoted,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_NotYetVoted,
						VoteType_SuccessObservation,
						VoteType_NotYetVoted,
						VoteType_NotYetVoted,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
			},
			expectedMap: map[string]int64{
				"Observer1": -1,
				"Observer2": -1,
				"Observer3": -1,
				"Observer4": 1,
			},
		},
		{
			name: "non finalized ballots not counted - 3 ballots",
			ballotList: []Ballot{
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_SuccessObservation,
						VoteType_NotYetVoted,
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_SuccessObservation,
						VoteType_NotYetVoted,
						VoteType_NotYetVoted,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotInProgress,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
					Votes: []VoteType{
						VoteType_NotYetVoted,
						VoteType_SuccessObservation,
						VoteType_NotYetVoted,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
			},
			expectedMap: map[string]int64{
				"Observer1": 0,
				"Observer2": 0,
				"Observer3": 0,
				"Observer4": 2,
			},
		},
		{
			name: "ballots not finalized",
			ballotList: []Ballot{
				{
					VoterList:    []string{"Observer1", "Observer2", "Observer3"},
					Votes:        []VoteType{VoteType_SuccessObservation, VoteType_NotYetVoted, VoteType_NotYetVoted},
					BallotStatus: BallotStatus_BallotInProgress,
				},
			},
			expectedMap: map[string]int64{},
		},
		{
			name: "observers performing differently across ballots - 3 ballots",
			ballotList: []Ballot{
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3"},
					Votes: []VoteType{
						VoteType_SuccessObservation,
						VoteType_FailureObservation,
						VoteType_NotYetVoted,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3"},
					Votes: []VoteType{
						VoteType_NotYetVoted,
						VoteType_FailureObservation,
						VoteType_FailureObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_FailureObservation,
				},
				{
					VoterList: []string{"Observer1", "Observer2", "Observer3"},
					Votes: []VoteType{
						VoteType_FailureObservation,
						VoteType_SuccessObservation,
						VoteType_SuccessObservation,
					},
					BallotStatus: BallotStatus_BallotFinalized_SuccessObservation,
				},
			},

			expectedMap: map[string]int64{
				"Observer1": -1,
				"Observer2": 1,
				"Observer3": 1,
			},
		}}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			result := BuildRewardsDistribution(test.ballotList)
			require.Equal(t, test.expectedMap, result)
		})
	}
}
