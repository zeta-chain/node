package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestNewMsgUpdateZRC20LiquidityCap_ValidateBasics(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgUpdateZRC20LiquidityCap
		err  error
	}{
		{
			name: "valid message",
			msg: types.NewMsgUpdateZRC20LiquidityCap(
				sample.AccAddress(),
				sample.EthAddress().String(),
				math.NewUint(1000),
			),
		},
		{
			name: "valid message with liquidity cap 0",
			msg: types.NewMsgUpdateZRC20LiquidityCap(
				sample.AccAddress(),
				sample.EthAddress().String(),
				math.ZeroUint(),
			),
		},
		{
			name: "valid message with liquidity cap nil",
			msg: types.NewMsgUpdateZRC20LiquidityCap(
				sample.AccAddress(),
				sample.EthAddress().String(),
				math.NewUint(1000),
			),
		},
		{
			name: "invalid address",
			msg: types.NewMsgUpdateZRC20LiquidityCap(
				"invalid_address",
				sample.EthAddress().String(),
				math.NewUint(1000),
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid contract address",
			msg: types.NewMsgUpdateZRC20LiquidityCap(
				sample.AccAddress(),
				"invalid_address",
				math.NewUint(1000),
			),
			err: sdkerrors.ErrInvalidAddress,
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

func TestNewMsgUpdateZRC20LiquidityCap_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateZRC20LiquidityCap
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateZRC20LiquidityCap{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateZRC20LiquidityCap{
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

func TestNewMsgUpdateZRC20LiquidityCap_Type(t *testing.T) {
	msg := types.MsgUpdateZRC20LiquidityCap{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateZRC20LiquidityCap, msg.Type())
}

func TestNewMsgUpdateZRC20LiquidityCap_Route(t *testing.T) {
	msg := types.MsgUpdateZRC20LiquidityCap{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestNewMsgUpdateZRC20LiquidityCap_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateZRC20LiquidityCap{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
