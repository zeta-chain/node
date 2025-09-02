package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestMsgUpdateGatewayGasLimit_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgUpdateGatewayGasLimit
		err  error
	}{
		{
			name: "invalid address",
			msg:  types.NewMsgUpdateGatewayGasLimit("invalid_address", sdkmath.NewInt(1000000)),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid gas limit - zero",
			msg:  types.NewMsgUpdateGatewayGasLimit(sample.AccAddress(), sdkmath.ZeroInt()),
			err:  sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid gas limit - negative",
			msg:  types.NewMsgUpdateGatewayGasLimit(sample.AccAddress(), sdkmath.NewInt(-1)),
			err:  sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid message",
			msg:  types.NewMsgUpdateGatewayGasLimit(sample.AccAddress(), sdkmath.NewInt(1000000)),
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

func TestMsgUpdateGatewayGasLimit_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateGatewayGasLimit
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateGatewayGasLimit{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateGatewayGasLimit{
				Creator: "invalid",
			},
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

func TestMsgUpdateGatewayGasLimit_Type(t *testing.T) {
	msg := types.MsgUpdateGatewayGasLimit{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateGatewayGasLimit, msg.Type())
}

func TestMsgUpdateGatewayGasLimit_Route(t *testing.T) {
	msg := types.MsgUpdateGatewayGasLimit{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateGatewayGasLimit_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateGatewayGasLimit{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
