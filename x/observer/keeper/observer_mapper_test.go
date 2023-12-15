package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_GetObserverSet(t *testing.T) {
	t.Run("get observer set", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		os := sample.ObserverSet(10)
		k.SetObservers(ctx, os)
		tfm, found := k.GetObserverSet(ctx)
		assert.True(t, found)
		assert.Equal(t, os, tfm)
	})
}
