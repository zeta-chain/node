package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func (m Ballot) AddVote(address string, vote VoteType) (Ballot, error) {
	if m.HasVoted(address) {
		return m, errors.Wrap(ErrUnableToAddVote, fmt.Sprintf(" Voter : %s | Ballot :%s | Already Voted", address, m.String()))
	}
	index := m.GetVoterIndex(address)
	m.Votes[index] = vote
	return m, nil
}

func (m Ballot) HasVoted(address string) bool {
	index := m.GetVoterIndex(address)
	return m.Votes[index] != VoteType_NotYetVoted
}

func (m Ballot) GetVoterIndex(address string) int {
	index := -1
	for i, addr := range m.VoterList {
		if addr == address {
			return i
		}
	}
	return index
}

func (m Ballot) IsBallotFinalized() (Ballot, bool) {
	if m.BallotStatus != BallotStatus_BallotInProgress {
		return m, false
	}
	success, failure := sdk.ZeroDec(), sdk.ZeroDec()
	total := sdk.NewDec(int64(len(m.VoterList)))
	for _, vote := range m.Votes {
		if vote == VoteType_SuccessObservation {
			success = success.Add(sdk.OneDec())
		}
		if vote == VoteType_FailureObservation {
			failure = failure.Add(sdk.OneDec())
		}

	}
	if failure.IsPositive() {
		if failure.Quo(total).GTE(m.BallotThreshold) {
			m.BallotStatus = BallotStatus_BallotFinalized_FailureObservation
			return m, true
		}
	}
	if success.IsPositive() {
		if success.Quo(total).GTE(m.BallotThreshold) {
			m.BallotStatus = BallotStatus_BallotFinalized_SuccessObservation
			return m, true
		}
	}
	return m, false
}

func CreateVotes(len int) []VoteType {
	voterList := make([]VoteType, len)
	for i := range voterList {
		voterList[i] = VoteType_NotYetVoted
	}
	return voterList
}
