package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgDisableCCTXFlags_ValidateBasic(t *testing.T) {
	tt := []struct {
		name string
		msg  *types.MsgDisableCCTXFlags
		err  require.ErrorAssertionFunc
	}{
		{
			name: "invalid creator address",
			msg:  types.NewMsgDisableCCTXFlags("invalid", true, true),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.Contains(t, err.Error(), "invalid creator address")
			},
		},
		{
			name: "valid",
			msg:  types.NewMsgDisableCCTXFlags(sample.AccAddress(), true, true),
			err:  require.NoError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.err(t, tc.msg.ValidateBasic())
		})
	}
}

func TestMsgDisableCCTXFlags_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgDisableCCTXFlags
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgDisableCCTXFlags{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgDisableCCTXFlags{
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

func TestMsgDisableCCTXFlags_Type(t *testing.T) {
	msg := types.MsgDisableCCTXFlags{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgDisableCCTXFlags, msg.Type())
}

func TestMsgDisableCCTXFlags_Route(t *testing.T) {
	msg := types.MsgDisableCCTXFlags{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgDisableCCTXFlags_GetSignBytes(t *testing.T) {
	msg := types.MsgDisableCCTXFlags{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
