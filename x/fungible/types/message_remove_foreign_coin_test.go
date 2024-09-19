package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestMsgRemoveForeignCoin_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgRemoveForeignCoin
		err  error
	}{
		{
			name: "invalid address",
			msg:  types.NewMsgRemoveForeignCoin("invalid_address", "name"),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg:  types.NewMsgRemoveForeignCoin(sample.AccAddress(), "name"),
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

func TestMsgRemoveForeignCoin_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgRemoveForeignCoin
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgRemoveForeignCoin{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgRemoveForeignCoin{
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

func TestMsgRemoveForeignCoin_Type(t *testing.T) {
	msg := types.MsgRemoveForeignCoin{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgRemoveForeignCoin, msg.Type())
}

func TestMsgRemoveForeignCoin_Route(t *testing.T) {
	msg := types.MsgRemoveForeignCoin{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgRemoveForeignCoin_GetSignBytes(t *testing.T) {
	msg := types.MsgRemoveForeignCoin{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
