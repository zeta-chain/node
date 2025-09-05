package e2etests

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// randomPayload generates a random payload to be used in gateway calls for testing purposes
func randomPayload(r *runner.E2ERunner) string {
	return randomPayloadWithSize(r, 100)
}

// randomPayloadWithSize generates a random payload string containing exactly 'size' characters
func randomPayloadWithSize(r *runner.E2ERunner, size int) string {
	// return empty string for invalid size
	if size <= 0 {
		return ""
	}

	var (
		bytes []byte
		even  = size%2 == 0
	)

	// the size will double when converting to hex string, so we just need half the size
	if even {
		bytes = make([]byte, size/2)
	} else {
		// add 1 more byte if size is odd, we trim it later
		bytes = make([]byte, (size/2)+1)
	}

	// generate random bytes
	_, err := rand.Read(bytes)
	require.NoError(r, err)

	// hex encode the bytes, the size is double
	if even {
		return hex.EncodeToString(bytes)
	}

	// trim the last letter to fit the odd size
	return hex.EncodeToString(bytes)[:size]
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
