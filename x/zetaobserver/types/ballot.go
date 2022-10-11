package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func (m Ballot) AddVote(address string, vote VoteType) (Ballot, error) {
	if m.BallotStatus != BallotStatus_BallotInProgress {
		return m, errors.Wrap(ErrUnableToAddVote, fmt.Sprintf(" Voter : %s | Status : %s | Ballot Already Finalized", address, m.VoterList[address]))
	}
	if !(m.VoterList[address] == VoteType_NotYetVoted) {
		return m, errors.Wrap(ErrUnableToAddVote, fmt.Sprintf(" Voter : %s | Status : %s", address, m.VoterList[address]))
	}
	m.VoterList[address] = vote
	return m, nil
}

func (m Ballot) IsBallotFinalized() (Ballot, bool) {
	if m.BallotStatus != BallotStatus_BallotInProgress {
		return m, false
	}
	success, failure := sdk.ZeroDec(), sdk.ZeroDec()
	total := sdk.NewDec(int64(len(m.VoterList)))
	for _, vote := range m.VoterList {
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

func CreateVoterList(addresses []string) map[string]VoteType {
	voterList := make(map[string]VoteType, len(addresses))
	fmt.Println("List ", voterList)
	for _, address := range addresses {
		voterList[address] = VoteType_NotYetVoted
	}
	fmt.Println("List 2", voterList)
	return voterList
}
