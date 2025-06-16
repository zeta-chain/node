package e2etests

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
)

// randomPayload generates a random payload to be used in gateway calls for testing purposes
func randomPayload(r *runner.E2ERunner) string {
	return hex.EncodeToString(randomPayloadBytes(r))
}

func randomPayloadBytes(r *runner.E2ERunner) []byte {
	bytes := make([]byte, 50)
	_, err := rand.Read(bytes)
	require.NoError(r, err)

	return bytes
}

// bigAdd is shorthand for new(big.Int).Add(x, y)
func bigAdd(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Add(x, y)
}

// bigSub is shorthand for new(big.Int).Sub(x, y)
func bigSub(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Sub(x, y)
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := d.Seconds() - float64(minutes*60)
	return fmt.Sprintf("%dm%.1fs", minutes, seconds)
}
