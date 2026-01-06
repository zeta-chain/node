package base

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

// signerTestSuite is a test suite for testing the signer
type signerTestSuite struct {
	*Signer
	tss *mocks.TSS
}

// newTestSuite creates a new test suite for testing
func newSignerTestSuite(t *testing.T) *signerTestSuite {
	// constructor parameters
	chain := chains.Ethereum
	tss := mocks.NewTSS(t)

	//logger := DefaultLogger()
	logger := Logger{}
	signer := NewSigner(chain, tss, logger, mode.StandardMode)

	suite := &signerTestSuite{
		Signer: signer,
		tss:    tss,
	}

	return suite
}

func TestNewSigner(t *testing.T) {
	signer := newSignerTestSuite(t)
	require.NotNil(t, signer)
}

func Test_BeingReportedFlag(t *testing.T) {
	signer := newSignerTestSuite(t)

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

func Test_PassesCompliance(t *testing.T) {
	signer := newSignerTestSuite(t)

	// create config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return false for restricted CCTX", func(t *testing.T) {
		cctx := sample.CrossChainTxV2(t, "abcd")
		cfg.ComplianceConfig.RestrictedAddresses = []string{cctx.InboundParams.Sender}
		config.SetRestrictedAddressesFromConfig(cfg)

		require.False(t, signer.PassesCompliance(cctx))
	})
	t.Run("should return true for non restricted CCTX", func(t *testing.T) {
		cctx := sample.CrossChainTxV2(t, "abcd")
		cfg.ComplianceConfig.RestrictedAddresses = []string{sample.EthAddress().Hex()}
		config.SetRestrictedAddressesFromConfig(cfg)

		require.True(t, signer.PassesCompliance(cctx))
	})
}
