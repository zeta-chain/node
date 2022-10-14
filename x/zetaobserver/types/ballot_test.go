package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVoter_IsBallotFinalized(t *testing.T) {
	tt := []struct {
		name           string
		threshold      sdk.Dec
		finalVoterList map[string]VoteType
		finalStatus    BallotStatus
		isFinalized    bool
	}{
		{
			name:      "All success",
			threshold: sdk.MustNewDecFromStr("0.66"),
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
			finalVoterList: map[string]VoteType{
				"Observer1": VoteType_SuccessObservation,
				"Observer2": VoteType_FailureObservation,
				"Observer3": VoteType_FailureObservation,
				"Observer4": VoteType_FailureObservation,
			},
			finalStatus: BallotStatus_BallotInProgress,
			isFinalized: false,
		},
	}
	for _, test := range tt {
		test := test
		t.Run(test.name, func(t *testing.T) {
			voterList := make(map[string]VoteType)
			voterList["Observer1"] = VoteType_NotYetVoted
			voterList["Observer2"] = VoteType_NotYetVoted
			voterList["Observer3"] = VoteType_NotYetVoted
			voterList["Observer4"] = VoteType_NotYetVoted
			ballot := Ballot{
				Index:            "index",
				BallotIdentifier: "identifier",
				VoterList:        voterList,
				ObservationType:  ObservationType_InBoundTx,
				BallotThreshold:  test.threshold,
				BallotStatus:     BallotStatus_BallotInProgress,
			}
			ballot.VoterList = test.finalVoterList
			isFinalized := ballot.IsBallotFinalized()
			assert.Equal(t, test.finalStatus, ballot.BallotStatus)
			assert.Equal(t, test.isFinalized, isFinalized)
		})
	}
}
