package types_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
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
				ChainIdList: make([]int64, 200),
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
				require.ErrorContains(t, err, "chain id list cannot be greater than supported chains")
			},
		},
		{
			name: "invalid chain id",
			msg: types.MsgDisableHeaderVerification{
				Creator:     sample.AccAddress(),
				ChainIdList: []int64{chains.ZetaPrivnetChain.ChainId},
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
				require.ErrorContains(t, err, fmt.Sprintf("invalid chain id header not supported (%d)", chains.ZetaPrivnetChain.ChainId))
			},
		},
		{
			name: "valid address",
			msg: types.MsgDisableHeaderVerification{
				Creator:     sample.AccAddress(),
				ChainIdList: []int64{chains.EthChain.ChainId},
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
		msg    *types.MsgEnableHeaderVerification
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.NewMsgEnableHeaderVerification(
				signer,
				[]int64{chains.EthChain.ChainId, chains.BtcMainnetChain.ChainId},
			),
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.NewMsgEnableHeaderVerification(
				"invalid",
				[]int64{chains.EthChain.ChainId, chains.BtcMainnetChain.ChainId},
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
	msg := types.MsgEnableHeaderVerification{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgDisableHeaderVerification_GetSignBytes(t *testing.T) {
	msg := types.MsgEnableHeaderVerification{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
