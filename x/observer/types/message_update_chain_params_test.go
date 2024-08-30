package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgUpdateChainParams_ValidateBasic(t *testing.T) {
	chainList := chains.ExternalChainList([]chains.Chain{})

	tests := []struct {
		name string
		msg  *types.MsgUpdateChainParams
		err  error
	}{
		{
			name: "valid message",
			msg: types.NewMsgUpdateChainParams(
				sample.AccAddress(),
				sample.ChainParams(chainList[0].ChainId),
			),
		},
		{
			name: "invalid address",
			msg: types.NewMsgUpdateChainParams(
				"invalid_address",
				sample.ChainParams(chainList[0].ChainId),
			),
			err: sdkerrors.ErrInvalidAddress,
		},

		{
			name: "invalid chain params (nil)",
			msg: types.NewMsgUpdateChainParams(
				sample.AccAddress(),
				nil,
			),
			err: types.ErrInvalidChainParams,
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

func TestMsgUpdateChainParams_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateChainParams
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateChainParams{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateChainParams{
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

func TestMsgUpdateChainParams_Type(t *testing.T) {
	msg := types.MsgUpdateChainParams{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateChainParams, msg.Type())
}

func TestMsgUpdateChainParams_Route(t *testing.T) {
	msg := types.MsgUpdateChainParams{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateChainParams_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateChainParams{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
