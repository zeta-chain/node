package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestBlockHeaderVerification_EnableChain(t *testing.T) {
	t.Run("should enable chain", func(t *testing.T) {
		bhv := sample.BlockHeaderVerification()
		bhv.EnableChain(chains.BscMainnetChain.ChainId)
		require.True(t, bhv.IsChainEnabled(chains.BscMainnetChain.ChainId))
	})
}
