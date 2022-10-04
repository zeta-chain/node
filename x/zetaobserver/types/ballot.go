package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func (m *Ballot) AddVote(address string, vote VoteType) error {
	if !(m.VoterList[address] == VoteType_NotYetVoted) {
		return errors.Wrap(ErrUnableToAddVote, fmt.Sprintf(" Voter : %s | Status : %s", address, m.VoterList[address]))
	}
	m.VoterList[address] = vote
	return nil
}

func (m *Ballot) IsBallotFinalized() bool {
	success, failure, total := sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()
	for _, vote := range m.VoterList {
		if vote != VoteType_NotYetVoted {
			total = total.Add(sdk.OneDec())
		}
		if vote == VoteType_SuccessObservation {
			success = success.Add(sdk.OneDec())
		}
		if vote == VoteType_FailureObservation {
			failure = failure.Add(sdk.OneDec())
		}

	}
	if total.IsZero() {
		m.BallotStatus = BallotStatus_BallotInProgress
		return false
	}
	if failure.IsPositive() {
		if failure.Quo(total).GTE(m.BallotThreshold) {
			m.BallotStatus = BallotStatus_BallotFinalized_FailureObservation
			return true
		}
	}
	if success.IsPositive() {
		if success.Quo(total).GTE(m.BallotThreshold) {
			m.BallotStatus = BallotStatus_BallotFinalized_SuccessObservation
			return true
		}
	}
	return false
}
