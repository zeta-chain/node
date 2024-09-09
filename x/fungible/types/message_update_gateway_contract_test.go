package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestMsgUpdateGatewayContract_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgUpdateGatewayContract
		err  error
	}{
		{
			name: "invalid address",
			msg:  types.NewMsgUpdateGatewayContract("invalid_address", sample.EthAddress().String()),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid new gateway contract address",
			msg:  types.NewMsgUpdateGatewayContract(sample.AccAddress(), "invalid_address"),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid message",
			msg:  types.NewMsgUpdateGatewayContract(sample.AccAddress(), sample.EthAddress().String()),
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

func TestMsgUpdateGatewayContract_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateGatewayContract
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateGatewayContract{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateGatewayContract{
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

func TestMsgUpdateGatewayContract_Type(t *testing.T) {
	msg := types.MsgUpdateGatewayContract{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateGatewayContract, msg.Type())
}

func TestMsgUpdateGatewayContract_Route(t *testing.T) {
	msg := types.MsgUpdateGatewayContract{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateGatewayContract_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateGatewayContract{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
