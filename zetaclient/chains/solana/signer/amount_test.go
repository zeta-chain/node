package signer

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	cctypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestOutboundAmountUint64(t *testing.T) {
	t.Run("returns amount within uint64 range", func(t *testing.T) {
		maxUint64 := ^uint64(0)
		amount, err := outboundAmountUint64(sdkmath.NewUint(maxUint64))
		require.NoError(t, err)
		require.Equal(t, maxUint64, amount)
	})

	t.Run("fails if amount exceeds uint64 range", func(t *testing.T) {
		amount := sdkmath.NewUintFromBigInt(new(big.Int).Lsh(big.NewInt(1), 64))

		_, err := outboundAmountUint64(amount)
		require.ErrorIs(t, err, cctypes.ErrInvalidWithdrawalAmount)
	})
}
