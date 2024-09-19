package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/emissions/types"
)

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	t.Run("invalid authority address", func(t *testing.T) {
		msg := types.MsgUpdateParams{
			Authority: "invalid",
			Params:    types.NewParams(),
		}
		err := msg.ValidateBasic()
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("invalid if params are invalid", func(t *testing.T) {
		params := types.NewParams()
		params.BlockRewardAmount = sdk.MustNewDecFromStr("-10.0")
		msg := types.MsgUpdateParams{
			Authority: sample.AccAddress(),
			Params:    params,
		}
		err := msg.ValidateBasic()
		require.ErrorContains(t, err, "block reward amount cannot be less than 0")
	})

	t.Run("valid", func(t *testing.T) {
		msg := types.MsgUpdateParams{
			Authority: sample.AccAddress(),
			Params:    types.NewParams(),
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
			msg:    &types.MsgUpdateParams{signer, types.NewParams()},
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    &types.MsgUpdateParams{"invalid", types.NewParams()},
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
	msg := types.MsgUpdateParams{"invalid", types.NewParams()}
	require.Equal(t, types.MsgUpdateParamsType, msg.Type())
}

func TestMsgUpdateParamsRoute(t *testing.T) {
	msg := types.MsgUpdateParams{"invalid", types.NewParams()}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateParamsGetSignBytes(t *testing.T) {
	msg := types.MsgUpdateParams{"invalid", types.NewParams()}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
