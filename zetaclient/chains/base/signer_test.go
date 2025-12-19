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

// createSigner creates a new signer for testing
func createSigner(t *testing.T) *Signer {
	// constructor parameters
	chain := chains.Ethereum
	tss := mocks.NewTSS(t)
	logger := DefaultLogger()

	// create signer
	return NewSigner(chain, tss, logger, mode.StandardMode)
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

func Test_PassesCompliance(t *testing.T) {
	signer := createSigner(t)

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
