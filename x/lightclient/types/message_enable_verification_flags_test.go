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

func TestMsgEnableHeaderVerification_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgEnableHeaderVerification
		err  require.ErrorAssertionFunc
	}{
		{
			name: "invalid address",
			msg: types.MsgEnableHeaderVerification{
				Creator: "invalid_address",
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
			},
		},
		{
			name: "empty chain id list",
			msg: types.MsgEnableHeaderVerification{
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
			msg: types.MsgEnableHeaderVerification{
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
			msg: types.MsgEnableHeaderVerification{
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

func TestMsgEnableHeaderVerification_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgEnableHeaderVerification
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.NewMsgEnableHeaderVerification(
				signer,
				[]int64{chains.Ethereum.ChainId, chains.BitcoinMainnet.ChainId},
			),
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.NewMsgEnableHeaderVerification(
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

func TestMsgEnableHeaderVerification_Type(t *testing.T) {
	msg := types.MsgEnableHeaderVerification{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgEnableHeaderVerification, msg.Type())
}

func TestMsgEnableHeaderVerification_Route(t *testing.T) {
	msg := types.MsgEnableHeaderVerification{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgEnableHeaderVerification_GetSignBytes(t *testing.T) {
	msg := types.MsgEnableHeaderVerification{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
