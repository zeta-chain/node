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

func TestMsgDisableFastConfirmation_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgDisableFastConfirmation
		err  require.ErrorAssertionFunc
	}{
		{
			name: "invalid address",
			msg:  types.NewMsgDisableFastConfirmation("invalid", []int64{}),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
			},
		},
		{
			name: "chain id list is too long",
			msg: types.NewMsgDisableFastConfirmation(
				sample.AccAddress(),
				make([]int64, types.MaxChainIDListLength+1), // 201 is the max length
			),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
				require.ErrorContains(t, err, "chain id list too long")
			},
		},
		{
			name: "valid",
			msg: types.NewMsgDisableFastConfirmation(
				sample.AccAddress(),
				[]int64{chains.Ethereum.ChainId},
			),
			err: require.NoError,
		},
		{
			name: "valid with empty chain id list",
			msg: types.NewMsgDisableFastConfirmation(
				sample.AccAddress(),
				nil,
			),
			err: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			tt.err(t, err)
		})
	}
}

func TestMsgDisableFastConfirmation_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgDisableFastConfirmation
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgDisableFastConfirmation(signer, []int64{}),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgDisableFastConfirmation("invalid", []int64{}),
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

func TestMsgDisableFastConfirmation_Type(t *testing.T) {
	msg := types.MsgDisableFastConfirmation{}
	require.Equal(t, types.TypeMsgDisableFastConfirmation, msg.Type())
}

func TestMsgDisableFastConfirmation_Route(t *testing.T) {
	msg := types.MsgDisableFastConfirmation{}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgDisableFastConfirmation_GetSignBytes(t *testing.T) {
	msg := types.NewMsgDisableFastConfirmation(sample.AccAddress(), []int64{chains.Ethereum.ChainId})
	require.NotPanics(t, func() {
		bytes := msg.GetSignBytes()
		require.NotEmpty(t, bytes)
	})
}
