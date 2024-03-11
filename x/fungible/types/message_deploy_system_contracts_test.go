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

func TestMsgDeploySystemContract_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgDeploySystemContracts
		err  error
	}{
		{
			name: "invalid address",
			msg:  types.NewMsgDeploySystemContracts("invalid"),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid message",
			msg:  types.NewMsgDeploySystemContracts(sample.AccAddress()),
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

func TestMsgDeploySystemContract_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgDeploySystemContracts
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgDeploySystemContracts(signer),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgDeploySystemContracts("invalid"),
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

func TestMsgDeploySystemContract_Type(t *testing.T) {
	msg := types.MsgDeploySystemContracts{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.TypeMsgDeploySystemContracts, msg.Type())
}

func TestMsgDeploySystemContract_Route(t *testing.T) {
	msg := types.MsgDeploySystemContracts{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgDeploySystemContract_GetSignBytes(t *testing.T) {
	msg := types.MsgDeploySystemContracts{
		Creator: sample.AccAddress(),
	}
	assert.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
