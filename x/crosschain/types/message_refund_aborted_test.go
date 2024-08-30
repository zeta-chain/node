package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestMsgRefundAbortedCCTX_ValidateBasic(t *testing.T) {
	t.Run("successfully validate message", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		msg := types.NewMsgRefundAbortedCCTX(sample.AccAddress(), cctx.Index, "")
		require.NoError(t, msg.ValidateBasic())
	})
	t.Run("invalid creator address", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		msg := types.NewMsgRefundAbortedCCTX("invalid", cctx.Index, "")
		require.ErrorContains(t, msg.ValidateBasic(), "invalid creator address")
	})
	t.Run("invalid cctx index", func(t *testing.T) {
		msg := types.NewMsgRefundAbortedCCTX(sample.AccAddress(), "invalid", "")
		require.ErrorContains(t, msg.ValidateBasic(), "invalid index hash")
	})
	t.Run("invalid refund address", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		msg := types.NewMsgRefundAbortedCCTX(sample.AccAddress(), cctx.Index, "invalid")
		require.ErrorContains(t, msg.ValidateBasic(), "invalid address")
	})
	t.Run("invalid refund address 2", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		msg := types.NewMsgRefundAbortedCCTX(
			sample.AccAddress(),
			cctx.Index,
			"0x91da5bf3F8Eb72724E6f50Ec6C3D199C6355c59",
		)
		require.ErrorContains(t, msg.ValidateBasic(), "invalid address")
	})
}

func TestMsgRefundAbortedCCTX_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgRefundAbortedCCTX
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgRefundAbortedCCTX(signer, "test", ""),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgRefundAbortedCCTX("invalid", "invalid", ""),
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				require.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				require.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestMsgRefundAbortedCCTX_Type(t *testing.T) {
	msg := types.NewMsgRefundAbortedCCTX(sample.AccAddress(), "test", "")
	require.Equal(t, types.RefundAborted, msg.Type())
}

func TestMsgRefundAbortedCCTX_Route(t *testing.T) {
	msg := types.NewMsgRefundAbortedCCTX(sample.AccAddress(), "test", "")
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgRefundAbortedCCTX_GetSignBytes(t *testing.T) {
	msg := types.NewMsgRefundAbortedCCTX(sample.AccAddress(), "test", "")
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
