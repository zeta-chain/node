package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestStatus_AbortRefunded(t *testing.T) {
	t.Run("should set status to aborted and message to CCTX aborted and Refunded", func(t *testing.T) {
		status := types.Status{
			Status:              0,
			StatusMessage:       "",
			LastUpdateTimestamp: 0,
			IsAbortRefunded:     false,
		}
		timestamp := time.Now().Unix()
		status.AbortRefunded(timestamp)
		require.Equal(t, status.IsAbortRefunded, true)
		require.Equal(t, status.StatusMessage, "CCTX aborted and Refunded")
		require.Equal(t, status.LastUpdateTimestamp, timestamp)
	})
}
