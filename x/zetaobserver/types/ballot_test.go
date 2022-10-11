package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVoter_IsBallotFinalized(t *testing.T) {
	tt := []struct {
		name             string
		threshold        sdk.Dec
		orginalVoterList map[string]VoteType
		finalVoterList   map[string]VoteType
		finalStatus      BallotStatus
		isFinalized      bool
	}{
		{
			name:      "All success",
			threshold: sdk.MustNewDecFromStr("0.66"),
			orginalVoterList: map[string]VoteType{
				"Observer1": VoteType_NotYetVoted,
				"Observer2": VoteType_NotYetVoted,
				"Observer3": VoteType_NotYetVoted,
				"Observer4": VoteType_NotYetVoted,
			},
			finalVoterList: map[string]VoteType{
				"Observer1": VoteType_SuccessObservation,
				"Observer2": VoteType_SuccessObservation,
				"Observer3": VoteType_SuccessObservation,
				"Observer4": VoteType_SuccessObservation,
			},
			finalStatus: BallotStatus_BallotFinalized_SuccessObservation,
			isFinalized: true,
		},
		{
			name:      "Unable to finalize",
			threshold: sdk.MustNewDecFromStr("0.66"),
			orginalVoterList: map[string]VoteType{
				"Observer1": VoteType_NotYetVoted,
				"Observer2": VoteType_NotYetVoted,
				"Observer3": VoteType_NotYetVoted,
				"Observer4": VoteType_NotYetVoted,
			},
			finalVoterList: map[string]VoteType{
				"Observer1": VoteType_SuccessObservation,
				"Observer2": VoteType_SuccessObservation,
				"Observer3": VoteType_FailureObservation,
				"Observer4": VoteType_FailureObservation,
			},
			finalStatus: BallotStatus_BallotInProgress,
			isFinalized: false,
		},
		{
			name:      "Low Threshold Failure first",
			threshold: sdk.MustNewDecFromStr("0.33"),
			orginalVoterList: map[string]VoteType{
				"Observer1": VoteType_NotYetVoted,
				"Observer2": VoteType_NotYetVoted,
				"Observer3": VoteType_NotYetVoted,
				"Observer4": VoteType_NotYetVoted,
			},
			finalVoterList: map[string]VoteType{
				"Observer1": VoteType_SuccessObservation,
				"Observer2": VoteType_SuccessObservation,
				"Observer3": VoteType_FailureObservation,
				"Observer4": VoteType_FailureObservation,
			},
			finalStatus: BallotStatus_BallotFinalized_FailureObservation,
			isFinalized: true,
		},
		{
			name:      "High threshold",
			threshold: sdk.MustNewDecFromStr("0.90"),
			orginalVoterList: map[string]VoteType{
				"Observer1": VoteType_NotYetVoted,
				"Observer2": VoteType_NotYetVoted,
				"Observer3": VoteType_NotYetVoted,
				"Observer4": VoteType_NotYetVoted,
			},
			finalVoterList: map[string]VoteType{
				"Observer1": VoteType_SuccessObservation,
				"Observer2": VoteType_FailureObservation,
				"Observer3": VoteType_FailureObservation,
				"Observer4": VoteType_FailureObservation,
			},
			finalStatus: BallotStatus_BallotInProgress,
			isFinalized: false,
		},
		{
			name:      "Two observers ",
			threshold: sdk.MustNewDecFromStr("1.00"),
			orginalVoterList: map[string]VoteType{
				"Observer1": VoteType_NotYetVoted,
				"Observer2": VoteType_NotYetVoted,
			},
			finalVoterList: map[string]VoteType{
				"Observer1": VoteType_SuccessObservation,
				"Observer2": VoteType_NotYetVoted,
			},
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
				VoterList:        test.orginalVoterList,
				ObservationType:  ObservationType_InBoundTx,
				BallotThreshold:  test.threshold,
				BallotStatus:     BallotStatus_BallotInProgress,
			}
			ballot.VoterList = test.finalVoterList
			finalBallot, isFinalized := ballot.IsBallotFinalized()
			assert.Equal(t, test.finalStatus, finalBallot.BallotStatus)
			assert.Equal(t, test.isFinalized, isFinalized)
		})
	}
}
