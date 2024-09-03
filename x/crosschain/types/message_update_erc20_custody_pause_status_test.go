package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestNewMsgUpdateERC20CustodyPauseStatus_ValidateBasic(t *testing.T) {
	keeper.SetConfig(false)
	tests := []struct {
		name  string
		msg   *types.MsgUpdateERC20CustodyPauseStatus
		error bool
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgUpdateERC20CustodyPauseStatus(
				"invalid address",
				chains.DefaultChainsList()[0].ChainId,
				true,
			),
			error: true,
		},
		{
			name: "valid msg",
			msg: types.NewMsgUpdateERC20CustodyPauseStatus(
				sample.AccAddress(),
				chains.DefaultChainsList()[0].ChainId,
				true,
			),
		},
		{
			name: "valid msg with pause false",
			msg: types.NewMsgUpdateERC20CustodyPauseStatus(
				sample.AccAddress(),
				chains.DefaultChainsList()[0].ChainId,
				false,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.error {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgUpdateERC20CustodyPauseStatus_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateERC20CustodyPauseStatus
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateERC20CustodyPauseStatus{
				Creator: signer,
				ChainId: chains.DefaultChainsList()[0].ChainId,
				Pause:   true,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateERC20CustodyPauseStatus{
				Creator: "invalid_address",
				ChainId: chains.DefaultChainsList()[0].ChainId,
				Pause:   true,
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

func TestMsgUpdateERC20CustodyPauseStatus_Type(t *testing.T) {
	msg := types.MsgUpdateERC20CustodyPauseStatus{
		Creator: sample.AccAddress(),
		ChainId: chains.DefaultChainsList()[0].ChainId,
		Pause:   true,
	}
	require.Equal(t, types.TypeUpdateERC20CustodyPauseStatus, msg.Type())
}

func TestMsgUpdateERC20CustodyPauseStatus_Route(t *testing.T) {
	msg := types.MsgUpdateERC20CustodyPauseStatus{
		Creator: sample.AccAddress(),
		ChainId: chains.DefaultChainsList()[0].ChainId,
		Pause:   true,
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateERC20CustodyPauseStatus_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateERC20CustodyPauseStatus{
		Creator: sample.AccAddress(),
		ChainId: chains.DefaultChainsList()[0].ChainId,
		Pause:   true,
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
