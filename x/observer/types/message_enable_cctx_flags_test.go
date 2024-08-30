package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgEnableCCTX_ValidateBasic(t *testing.T) {
	tt := []struct {
		name string
		msg  *types.MsgEnableCCTX
		err  require.ErrorAssertionFunc
	}{
		{
			name: "invalid creator address",
			msg:  types.NewMsgEnableCCTX("invalid", true, true),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "invalid creator address")
			},
		},
		{
			name: "invalid flags",
			msg:  types.NewMsgEnableCCTX(sample.AccAddress(), false, false),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "at least one of EnableInbound or EnableOutbound must be true")
			},
		},
		{
			name: "valid",
			msg:  types.NewMsgEnableCCTX(sample.AccAddress(), true, true),
			err:  require.NoError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.err(t, tc.msg.ValidateBasic())
		})
	}
}

func TestMsgEnableCCTX_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgEnableCCTX
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgEnableCCTX{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgEnableCCTX{
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

func TestMsgEnableCCTX_Type(t *testing.T) {
	msg := types.MsgEnableCCTX{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgEnableCCTX, msg.Type())
}

func TestMsgEnableCCTX_Route(t *testing.T) {
	msg := types.MsgEnableCCTX{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgEnableCCTX_GetSignBytes(t *testing.T) {
	msg := types.MsgEnableCCTX{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
