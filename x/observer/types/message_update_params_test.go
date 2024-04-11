package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	t.Run("invalid authority address", func(t *testing.T) {
		msg := types.MsgUpdateParams{
			Authority: "invalid",
			Params:    types.DefaultParams(),
		}
		err := msg.ValidateBasic()
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("valid", func(t *testing.T) {
		msg := types.MsgUpdateParams{
			Authority: sample.AccAddress(),
			Params:    types.DefaultParams(),
		}
		err := msg.ValidateBasic()
		require.NoError(t, err)
	})
}

func TestMsgUpdateParamsGetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgUpdateParams
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    &types.MsgUpdateParams{signer, types.DefaultParams()},
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    &types.MsgUpdateParams{"invalid", types.DefaultParams()},
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

func TestMsgUpdateParamsType(t *testing.T) {
	msg := types.MsgUpdateParams{"invalid", types.DefaultParams()}
	require.Equal(t, types.MsgUpdateParamsType, msg.Type())
}

func TestMsgUpdateParamsRoute(t *testing.T) {
	msg := types.MsgUpdateParams{"invalid", types.DefaultParams()}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateParamsGetSignBytes(t *testing.T) {
	msg := types.MsgUpdateParams{"invalid", types.DefaultParams()}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
