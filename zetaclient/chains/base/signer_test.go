package base_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

// createSigner creates a new signer for testing
func createSigner(_ *testing.T) *base.Signer {
	// constructor parameters
	chain := chains.Ethereum
	zetacoreContext := context.NewZetacoreContext(config.NewConfig())
	tss := mocks.NewTSSMainnet()
	logger := base.DefaultLogger()

	// create signer
	return base.NewSigner(chain, zetacoreContext, tss, nil, logger)
}

func TestNewSigner(t *testing.T) {
	signer := createSigner(t)
	require.NotNil(t, signer)
}

func TestSignerGetterAndSetter(t *testing.T) {
	t.Run("should be able to update chain", func(t *testing.T) {
		signer := createSigner(t)

		// update chain
		newChain := chains.BscMainnet
		signer = signer.WithChain(chains.BscMainnet)
		require.Equal(t, newChain, signer.Chain())
	})
	t.Run("should be able to update zetacore context", func(t *testing.T) {
		signer := createSigner(t)

		// update zetacore context
		newZetacoreContext := context.NewZetacoreContext(config.NewConfig())
		signer = signer.WithZetacoreContext(newZetacoreContext)
		require.Equal(t, newZetacoreContext, signer.ZetacoreContext())
	})
	t.Run("should be able to update tss", func(t *testing.T) {
		signer := createSigner(t)

		// update tss
		newTSS := mocks.NewTSSAthens3()
		signer = signer.WithTSS(newTSS)
		require.Equal(t, newTSS, signer.TSS())
	})
	t.Run("should be able to update telemetry server", func(t *testing.T) {
		signer := createSigner(t)

		// update telemetry server
		newTs := metrics.NewTelemetryServer()
		signer = signer.WithTelemetryServer(newTs)
		require.Equal(t, newTs, signer.TelemetryServer())
	})
	t.Run("should be able to get logger", func(t *testing.T) {
		ob := createSigner(t)
		logger := ob.Logger()

		// should be able to print log
		logger.Std.Info().Msg("print standard log")
		logger.Compliance.Info().Msg("print compliance log")
	})
	t.Run("should be able to get mutex", func(t *testing.T) {
		signer := createSigner(t)
		require.NotNil(t, signer.Mu())
	})
}
