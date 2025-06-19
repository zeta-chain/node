package e2etests

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// randomPayload generates a random payload to be used in gateway calls for testing purposes
func randomPayload(r *runner.E2ERunner) string {
	bytes := make([]byte, 50)
	_, err := rand.Read(bytes)
	require.NoError(r, err)

	return hex.EncodeToString(bytes)
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

func requireCCTXStatus(
	r *runner.E2ERunner,
	expectedStatus crosschaintypes.CctxStatus,
	cctx *crosschaintypes.CrossChainTx,
) {
	if expectedStatus == cctx.CctxStatus.Status {
		return
	}
	require.Failf(
		r,
		"cctx status mismatch",
		"cctx index %s ,expected status: %s, got status: %s",
		cctx.Index,
		expectedStatus.String(),
		cctx.CctxStatus.Status.String(),
	)
}
