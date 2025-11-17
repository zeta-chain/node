package simulation

import (
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// TSSVoteSimulationMatrix returns a transition matrix and a state array for the TSS ballot vote simulation
// This simulates the vote cast for the TSS creation ballot
func TSSVoteSimulationMatrix() (simtypes.TransitionMatrix, []float64, int) {
	ballotTransitionMatrix, _ := simulation.CreateTransitionMatrix([][]int{
		{70, 10},
		{30, 10},
	})
	// The states are:
	// column 1: 100% vote yes
	// column 2: 0% vote yes
	// For all conditions we assume if the vote is not a yes
	// then it is a no .
	yesVoteArray := []float64{1, 0}
	ballotVotesState := 1
	return ballotTransitionMatrix, yesVoteArray, ballotVotesState
}

// OutboundVoteStatusSimulationMatrix returns a transition matrix and a state array for the outbound vote simulation
// This simulates the ReceiveStatus field of the outbound vote
func OutboundVoteStatusSimulationMatrix() (simtypes.TransitionMatrix, []float64, int) {
	ballotTransitionMatrix, _ := simulation.CreateTransitionMatrix([][]int{
		{70, 10, 20},
		{20, 30, 30},
		{10, 60, 50},
	})
	// The states are:
	// column 1: 100% vote yes
	// column 2: 50% vote yes
	// column 3: 0% vote yes
	// For all conditions we assume if the vote is not a yes.
	// then it is a no .Not voting condition is handled by the ObserverVotesSimulationMatrix matrix
	yesVoteArray := []float64{1, .5, 0}
	ballotVotesState := 1
	return ballotTransitionMatrix, yesVoteArray, ballotVotesState
}

// ObserverVotesSimulationMatrix returns a transition matrix and a state array for the observer votes simulation.
// This is used to simulate the number of votes that will be cast by the observers for both inbound and outbound votes
func ObserverVotesSimulationMatrix() (simtypes.TransitionMatrix, []float64, int) {
	observerVotesTransitionMatrix, _ := simulation.CreateTransitionMatrix([][]int{
		{20, 10, 0, 0, 0, 0},
		{55, 50, 20, 10, 0, 0},
		{25, 25, 30, 25, 30, 15},
		{0, 15, 30, 25, 30, 30},
		{0, 0, 20, 30, 30, 30},
		{0, 0, 0, 10, 10, 25},
	})
	// The states are:
	// column 1: All observers vote
	// column 2: 90% vote
	// column 3: 75% vote
	// column 4: 40% vote
	// column 5: 15% vote
	// column 6: no one votes
	// All columns sum to 100 for simplicity, but this is arbitrary and can be changed
	statePercentageArray := []float64{1, .9, .75, .4, .15, 0}
	curNumVotesState := 1
	return observerVotesTransitionMatrix, statePercentageArray, curNumVotesState
}
