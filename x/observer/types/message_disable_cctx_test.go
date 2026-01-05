package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgDisableCCTX_ValidateBasic(t *testing.T) {
	tt := []struct {
		name string
		msg  *types.MsgDisableCCTX
		err  require.ErrorAssertionFunc
	}{
		{
			name: "invalid creator address",
			msg:  types.NewMsgDisableCCTX("invalid", true, true),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "invalid creator address")
			},
		},
		{
			name: "invalid flags",
			msg:  types.NewMsgDisableCCTX(sample.AccAddress(), false, false),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "at least one of DisableInbound or DisableOutbound must be true")
			},
		},
		{
			name: "valid",
			msg:  types.NewMsgDisableCCTX(sample.AccAddress(), true, true),
			err:  require.NoError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.err(t, tc.msg.ValidateBasic())
		})
	}
}

func TestMsgDisableCCTX_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgDisableCCTX
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgDisableCCTX{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgDisableCCTX{
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

func TestMsgDisableCCTX_Type(t *testing.T) {
	msg := types.MsgDisableCCTX{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgDisableCCTX, msg.Type())
}

func TestMsgDisableCCTX_Route(t *testing.T) {
	msg := types.MsgDisableCCTX{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgDisableCCTX_GetSignBytes(t *testing.T) {
	msg := types.MsgDisableCCTX{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
