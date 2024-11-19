package simulation_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/zeta-chain/node/x/crosschain/simulation"
)

func Test_Matrix(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	numVotes := 10
	for i := 0; i < 10; i++ {

		ballotVotesTransitionMatrix, yesVotePercentageArray, ballotVotesState := simulation.BallotVoteSimulationMatrix()
		ballotVotesState = ballotVotesTransitionMatrix.NextState(r, ballotVotesState)
		yesVotePercentage := yesVotePercentageArray[ballotVotesState]
		numberOfYesVotes := int(math.Ceil(float64(numVotes) * yesVotePercentage))

		t.Logf("Yes Vote Percentage: %v, Number of Yes votes %d", yesVotePercentage, numberOfYesVotes)
	}
}
