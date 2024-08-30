package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/lightclient/types"
)

func TestMsgDisableHeaderVerification_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgDisableHeaderVerification
		err  require.ErrorAssertionFunc
	}{
		{
			name: "invalid address",
			msg: types.MsgDisableHeaderVerification{
				Creator: "invalid_address",
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
			},
		},
		{
			name: "empty chain id list",
			msg: types.MsgDisableHeaderVerification{
				Creator:     sample.AccAddress(),
				ChainIdList: []int64{},
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
				require.ErrorContains(t, err, "chain id list cannot be empty")
			},
		},
		{
			name: "chain id list is too long",
			msg: types.MsgDisableHeaderVerification{
				Creator:     sample.AccAddress(),
				ChainIdList: make([]int64, types.MaxChainIDListLength+1),
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
				require.ErrorContains(t, err, "chain id list too long")
			},
		},
		{
			name: "valid address",
			msg: types.MsgDisableHeaderVerification{
				Creator:     sample.AccAddress(),
				ChainIdList: []int64{chains.Ethereum.ChainId},
			},
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

func TestMsgDisableHeaderVerification_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgDisableHeaderVerification
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.NewMsgDisableHeaderVerification(
				signer,
				[]int64{chains.Ethereum.ChainId, chains.BitcoinMainnet.ChainId},
			),
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.NewMsgDisableHeaderVerification(
				"invalid",
				[]int64{chains.Ethereum.ChainId, chains.BitcoinMainnet.ChainId},
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

func TestMsgDisableHeaderVerification_Type(t *testing.T) {
	msg := types.MsgDisableHeaderVerification{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgDisableHeaderVerification, msg.Type())
}

func TestMsgDisableHeaderVerification_Route(t *testing.T) {
	msg := types.MsgDisableHeaderVerification{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgDisableHeaderVerification_GetSignBytes(t *testing.T) {
	msg := types.MsgDisableHeaderVerification{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
