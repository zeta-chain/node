package local

import (
	"time"

	"github.com/spf13/cobra"
)

// Default timeout constants
const (
	DefaultTestTimeout    = 20 * time.Minute
	DefaultReceiptTimeout = 8 * time.Minute
	DefaultCctxTimeout    = 8 * time.Minute
	StressReceiptTimeout  = 15 * time.Minute
	StressCctxTimeout     = 15 * time.Minute
)

// TestTimeouts holds the various timeout values for E2E tests
type TestTimeouts struct {
	// TestTimeout is the overall timeout for the entire test suite
	TestTimeout time.Duration

	// ReceiptTimeout is the timeout for waiting for transaction receipts
	ReceiptTimeout time.Duration

	// CctxTimeout is the timeout for waiting for cross-chain transactions to be mined
	CctxTimeout time.Duration
}

// RegularTestTimeouts returns timeout values for regular tests from command flags
func RegularTestTimeouts(cmd *cobra.Command) TestTimeouts {
	return TestTimeouts{
		TestTimeout:    must(cmd.Flags().GetDuration(flagTestTimeout)),
		ReceiptTimeout: must(cmd.Flags().GetDuration(flagReceiptTimeout)),
		CctxTimeout:    must(cmd.Flags().GetDuration(flagCctxTimeout)),
	}
}

// StressTestTimeouts returns timeout values for stress tests.
// It adjusts the overall TestTimeout based on the number of iterations.
// Can be overridden by command flags.
func StressTestTimeouts(cmd *cobra.Command, iterations int) TestTimeouts {
	timeouts := TestTimeouts{
		TestTimeout:    DefaultTestTimeout,
		ReceiptTimeout: StressReceiptTimeout,
		CctxTimeout:    StressCctxTimeout,
	}

	if iterations > 100 {
		timeouts.TestTimeout = time.Hour
	}

	if cmd.Flags().Changed(flagTestTimeout) {
		timeouts.TestTimeout = must(cmd.Flags().GetDuration(flagTestTimeout))
	}
	if cmd.Flags().Changed(flagReceiptTimeout) {
		timeouts.ReceiptTimeout = must(cmd.Flags().GetDuration(flagReceiptTimeout))
	}
	if cmd.Flags().Changed(flagCctxTimeout) {
		timeouts.CctxTimeout = must(cmd.Flags().GetDuration(flagCctxTimeout))
	}

	return timeouts
}
