package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestNewMsgRefundAbortedCCTX(t *testing.T) {
	t.Run("successfully validate message", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		msg := types.NewMsgRefundAbortedCCTX(sample.AccAddress(), cctx.Index)
		assert.NoError(t, msg.ValidateBasic())
	})
	t.Run("invalid creator address", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		msg := types.NewMsgRefundAbortedCCTX("invalid", cctx.Index)
		assert.ErrorContains(t, msg.ValidateBasic(), "invalid creator address")
	})
	t.Run("invalid cctx index", func(t *testing.T) {
		msg := types.NewMsgRefundAbortedCCTX(sample.AccAddress(), "invalid")
		assert.ErrorContains(t, msg.ValidateBasic(), "invalid cctx index")
	})
}
