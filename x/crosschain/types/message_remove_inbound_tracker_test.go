package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestMsgRemoveInboundTracker_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgRemoveInboundTracker
		err  error
	}{
		{
			name: "invalid address",
			msg: types.NewMsgRemoveInboundTracker(
				"invalid_address",
				chains.Goerli.ChainId,
				sample.ZetaIndex(t),
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid if chain id is negative",
			msg: types.NewMsgRemoveInboundTracker(
				sample.AccAddress(),
				-1,
				sample.ZetaIndex(t),
			),
			err: sdkerrors.ErrInvalidChainID,
		},
		{
			name: "valid",
			msg: types.NewMsgRemoveInboundTracker(
				sample.AccAddress(),
				chains.Goerli.ChainId,
				sample.ZetaIndex(t),
			),
			err: nil,
		},
		{
			name: "valid even if chain id is not supported",
			msg: types.NewMsgRemoveInboundTracker(
				sample.AccAddress(),
				999,
				sample.ZetaIndex(t),
			),
			err: nil,
		},
		{
			name: "valid even if tx hash is not supported",
			msg: types.NewMsgRemoveInboundTracker(
				sample.AccAddress(),
				chains.Goerli.ChainId,
				"invalid",
			),
			err: nil,
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

func TestMsgRemoveInboundTracker_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgRemoveInboundTracker
		panics bool
	}{
		{
			name: "valid",
			msg: types.NewMsgRemoveInboundTracker(
				signer,
				chains.Goerli.ChainId,
				sample.ZetaIndex(t),
			),
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.NewMsgRemoveInboundTracker(
				"invalid_address",
				chains.Goerli.ChainId,
				"hash",
			),
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

func TestMsgRemoveInboundTracker_Type(t *testing.T) {
	msg := types.NewMsgRemoveInboundTracker(
		sample.AccAddress(),
		chains.Goerli.ChainId,
		sample.ZetaIndex(t),
	)
	require.Equal(t, types.TypeMsgRemoveInboundTracker, msg.Type())
}

func TestMsgRemoveInboundTracker_Route(t *testing.T) {
	msg := types.NewMsgRemoveInboundTracker(
		sample.AccAddress(),
		chains.Goerli.ChainId,
		sample.ZetaIndex(t),
	)
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgRemoveInboundTracker_GetSignBytes(t *testing.T) {
	msg := types.NewMsgRemoveInboundTracker(
		sample.AccAddress(),
		chains.Goerli.ChainId,
		sample.ZetaIndex(t),
	)

	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
	require.NotEmpty(t, msg.GetSignBytes())
	require.Equal(t, string(msg.GetSignBytes()), string(msg.GetSignBytes()), "sign bytes should be deterministic")
}
