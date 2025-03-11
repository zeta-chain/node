package base_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

// createSigner creates a new signer for testing
func createSigner(t *testing.T) *base.Signer {
	// constructor parameters
	chain := chains.Ethereum
	tss := mocks.NewTSS(t)
	logger := base.DefaultLogger()

	// create signer
	signer, err := base.NewSigner(chain, tss, 1, 0, logger)
	require.NoError(t, err)

	return signer
}

func TestNewSigner(t *testing.T) {
	signer := createSigner(t)
	require.NotNil(t, signer)
}

func Test_BeingReportedFlag(t *testing.T) {
	signer := createSigner(t)

	// hash to be reported
	hash := "0x1234"
	alreadySet := signer.SetBeingReportedFlag(hash)
	require.False(t, alreadySet)

	// set reported outbound again and check
	alreadySet = signer.SetBeingReportedFlag(hash)
	require.True(t, alreadySet)

	// clear reported outbound and check again
	signer.ClearBeingReportedFlag(hash)
	alreadySet = signer.SetBeingReportedFlag(hash)
	require.False(t, alreadySet)
}
