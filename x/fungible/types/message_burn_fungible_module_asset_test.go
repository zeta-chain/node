package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestMsgBurnFungibleModuleAsset_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgBurnFungibleModuleAsset
		err  error
	}{
		{
			name: "invalid address",
			msg: types.NewMsgBurnFungibleModuleAsset(
				"invalid_address",
				sample.EthAddress().String(),
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid zrc20 address",
			msg: types.NewMsgBurnFungibleModuleAsset(
				sample.AccAddress(),
				"invalid_address",
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid message",
			msg: types.NewMsgBurnFungibleModuleAsset(
				sample.AccAddress(),
				sample.EthAddress().String(),
			),
		},
		{
			name: "zero address is valid",
			msg: types.NewMsgBurnFungibleModuleAsset(
				sample.AccAddress(),
				constant.EVMZeroAddress,
			),
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

func TestMsgBurnFungibleModuleAsset_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgBurnFungibleModuleAsset
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgBurnFungibleModuleAsset{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgBurnFungibleModuleAsset{
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

func TestMsgBurnFungibleModuleAsset_Type(t *testing.T) {
	msg := types.MsgBurnFungibleModuleAsset{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgBurnFungibleModuleAsset, msg.Type())
}

func TestMsgBurnFungibleModuleAsset_Route(t *testing.T) {
	msg := types.MsgBurnFungibleModuleAsset{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgBurnFungibleModuleAsset_GetSignBytes(t *testing.T) {
	msg := types.MsgBurnFungibleModuleAsset{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
