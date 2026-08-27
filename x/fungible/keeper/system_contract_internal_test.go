package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestUniswapV2Deadline(t *testing.T) {
	t.Run("uses block time when available", func(t *testing.T) {
		blockTime := time.Unix(1_700_000_000, 0).UTC()
		ctx := sdk.Context{}.WithBlockTime(blockTime)

		require.Equal(t, blockTime.Add(uniswapV2SwapDeadline).Unix(), uniswapV2Deadline(ctx).Int64())
	})

	t.Run("uses far future fallback for zero block time", func(t *testing.T) {
		require.Equal(t, uniswapV2FallbackDeadline, uniswapV2Deadline(sdk.Context{}).Int64())
	})
}
