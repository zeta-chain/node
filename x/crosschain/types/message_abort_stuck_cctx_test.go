package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestMsgAbortStuckCCTX_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgAbortStuckCCTX
		err  error
	}{
		{
			name: "invalid address",
			msg:  types.NewMsgAbortStuckCCTX("invalid_address", "cctx_index"),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid cctx index",
			msg:  types.NewMsgAbortStuckCCTX(sample.AccAddress(), "cctx_index"),
			err:  types.ErrInvalidIndexValue,
		},
		{
			name: "valid",
			msg:  types.NewMsgAbortStuckCCTX(sample.AccAddress(), sample.GetCctxIndexFromString("test")),
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgAbortStuckCCTX_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgAbortStuckCCTX
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgAbortStuckCCTX(signer, "cctx_index"),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgAbortStuckCCTX("invalid", "cctx_index"),
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

func TestMsgAbortStuckCCTX_Type(t *testing.T) {
	msg := types.NewMsgAbortStuckCCTX(sample.AccAddress(), "cctx_index")
	require.Equal(t, types.TypeMsgAbortStuckCCTX, msg.Type())
}

func TestMsgAbortStuckCCTX_Route(t *testing.T) {
	msg := types.NewMsgAbortStuckCCTX(sample.AccAddress(), "cctx_index")
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgAbortStuckCCTX_GetSignBytes(t *testing.T) {
	msg := types.NewMsgAbortStuckCCTX(sample.AccAddress(), "cctx_index")
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
