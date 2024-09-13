package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestBallot_AddVote(t *testing.T) {
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
		wantErr     bool
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
			threshold: sdk.MustNewDecFromStr("0.66"),
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
			threshold: sdk.MustNewDecFromStr("0.66"),
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
			threshold: sdk.MustNewDecFromStr("0.66"),
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
			threshold: sdk.MustNewDecFromStr("0.66"),
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
		{
			name:      "Low threshold 1 always fails as Failure is checked first",
			threshold: sdk.MustNewDecFromStr("0.01"),
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
			threshold: sdk.MustNewDecFromStr("0.01"),
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
			threshold: sdk.MustNewDecFromStr("1"),
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
			threshold: sdk.MustNewDecFromStr("0.66"),
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
		test := test
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
		BallotThreshold sdk.Dec
		Votes           []VoteType
		finalizingVote  int
		finalStatus     BallotStatus
	}{
		{
			name:            "finalized to success",
			BallotThreshold: sdk.MustNewDecFromStr("0.66"),
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
			BallotThreshold: sdk.MustNewDecFromStr("0.66"),
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
			BallotThreshold: sdk.MustNewDecFromStr("0.01"),
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
			BallotThreshold: sdk.MustNewDecFromStr("1"),
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
			BallotThreshold: sdk.MustNewDecFromStr("1"),
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
		test := test
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
		name         string
		voterList    []string
		votes        []VoteType
		ballotStatus BallotStatus
		expectedMap  map[string]int64
	}{
		{
			name:      "BallotFinalized_SuccessObservation",
			voterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []VoteType{
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_FailureObservation,
			},
			ballotStatus: BallotStatus_BallotFinalized_SuccessObservation,
			expectedMap: map[string]int64{
				"Observer1": 1,
				"Observer2": 1,
				"Observer3": 1,
				"Observer4": -1,
			},
		},
		{
			name:      "BallotFinalized_FailureObservation",
			voterList: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []VoteType{
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
			},
			ballotStatus: BallotStatus_BallotFinalized_FailureObservation,
			expectedMap: map[string]int64{
				"Observer1": -1,
				"Observer2": -1,
				"Observer3": 1,
				"Observer4": 1,
			},
		},
	}
	for _, test := range tt {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ballot := Ballot{
				Index:            "",
				BallotIdentifier: "",
				VoterList:        test.voterList,
				Votes:            test.votes,
				ObservationType:  0,
				BallotThreshold:  sdk.Dec{},
				BallotStatus:     test.ballotStatus,
			}
			rewardsMap := map[string]int64{}
			ballot.BuildRewardsDistribution(rewardsMap)
			require.Equal(t, test.expectedMap, rewardsMap)
		})
	}
}

func TestBallot_GenerateVoterList(t *testing.T) {
	tt := []struct {
		name              string
		voters            []string
		votes             []VoteType
		expectedVoterList []VoterList
		err               require.ErrorAssertionFunc
	}{
		{
			name:   "Success observation",
			voters: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []VoteType{
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
			},
			expectedVoterList: []VoterList{
				{
					VoterAddress: "Observer1",
					VoteType:     VoteType_SuccessObservation,
				},
				{
					VoterAddress: "Observer2",
					VoteType:     VoteType_SuccessObservation,
				},
				{
					VoterAddress: "Observer3",
					VoteType:     VoteType_SuccessObservation,
				},
				{
					VoterAddress: "Observer4",
					VoteType:     VoteType_SuccessObservation,
				},
			},
			err: require.NoError,
		},

		{
			name:   "Failure observation",
			voters: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []VoteType{
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_FailureObservation,
			},
			expectedVoterList: []VoterList{
				{
					VoterAddress: "Observer1",
					VoteType:     VoteType_FailureObservation,
				},
				{
					VoterAddress: "Observer2",
					VoteType:     VoteType_FailureObservation,
				},
				{
					VoterAddress: "Observer3",
					VoteType:     VoteType_FailureObservation,
				},
				{
					VoterAddress: "Observer4",
					VoteType:     VoteType_FailureObservation,
				},
			},
			err: require.NoError,
		},

		{
			name:   "mixed observation",
			voters: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []VoteType{
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_SuccessObservation,
				VoteType_SuccessObservation,
			},
			expectedVoterList: []VoterList{
				{
					VoterAddress: "Observer1",
					VoteType:     VoteType_FailureObservation,
				},
				{
					VoterAddress: "Observer2",
					VoteType:     VoteType_FailureObservation,
				},
				{
					VoterAddress: "Observer3",
					VoteType:     VoteType_SuccessObservation,
				},
				{
					VoterAddress: "Observer4",
					VoteType:     VoteType_SuccessObservation,
				},
			},
			err: require.NoError,
		},

		{
			name:   "voterList and votes length mismatch",
			voters: []string{"Observer1", "Observer2", "Observer3", "Observer4"},
			votes: []VoteType{
				VoteType_FailureObservation,
				VoteType_FailureObservation,
				VoteType_SuccessObservation,
			},
			expectedVoterList: nil,
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.Error(t, err)
				require.Equal(t, err, ErrInvalidVoterList)
			},
		},
	}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			ballot := Ballot{
				VoterList: test.voters,
				Votes:     test.votes,
			}
			voterList, err := ballot.GenerateVoterList()
			test.err(t, err)
			require.Equal(t, test.expectedVoterList, voterList)
		})
	}
}
