package simulation_test

import (
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func Test_Matrix(t *testing.T) {
	numVotesTransitionMatrix, _ := simulation.CreateTransitionMatrix([][]int{
		{20, 10, 0, 0, 0, 0},
		{55, 50, 20, 10, 0, 0},
		{25, 25, 30, 25, 30, 15},
		{0, 15, 30, 25, 30, 30},
		{0, 0, 20, 30, 30, 30},
		{0, 0, 0, 10, 10, 25},
	})

	statePercentageArray := []float64{1, .9, .75, .4, .15, 0}
	yesVoteArray := []float64{1, .5, 0}
	curNumVotesState := 1

	ballotTransitionMatrix, _ := simulation.CreateTransitionMatrix([][]int{
		{20, 50, 50},
		{55, 50, 20},
		{25, 0, 30},
	})
	ballotVotesState := 1

	r := rand.New(rand.NewSource(1))

	for i := 0; i < 10; i++ {
		curNumVotesState = numVotesTransitionMatrix.NextState(r, curNumVotesState)
		ballotVotesState = ballotTransitionMatrix.NextState(r, ballotVotesState)
		percentageVote := statePercentageArray[curNumVotesState]
		percentageBallotYest := yesVoteArray[ballotVotesState]

		t.Logf("iteration %d ,curNumVotesState: %d,"+
			" ballotVotesState: %d, "+
			"percentageVote: %f, "+
			"percentageBallotYest: %f", i, curNumVotesState, ballotVotesState, percentageVote, percentageBallotYest)

	}
}
