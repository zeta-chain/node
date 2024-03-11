package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestMsgUpdateSystemContract_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgUpdateSystemContract
		err  error
	}{
		{
			name: "invalid address",
			msg:  types.NewMsgUpdateSystemContract("invalid_address", sample.EthAddress().String()),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid new system contract address",
			msg:  types.NewMsgUpdateSystemContract(sample.AccAddress(), "invalid_address"),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid message",
			msg:  types.NewMsgUpdateSystemContract(sample.AccAddress(), sample.EthAddress().String()),
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

func TestMsgUpdateSystemContract_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateSystemContract
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateSystemContract{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateSystemContract{
				Creator: "invalid",
			},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				assert.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				assert.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestMsgUpdateSystemContract_Type(t *testing.T) {
	msg := types.MsgUpdateSystemContract{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.TypeMsgUpdateSystemContract, msg.Type())
}

func TestMsgUpdateSystemContract_Route(t *testing.T) {
	msg := types.MsgUpdateSystemContract{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateSystemContract_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateSystemContract{
		Creator: sample.AccAddress(),
	}
	assert.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
