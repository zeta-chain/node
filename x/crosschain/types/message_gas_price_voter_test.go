package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/authz"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgGasPriceVoter_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgGasPriceVoter
		err  error
	}{
		{
			name: "invalid address",
			msg: types.NewMsgGasPriceVoter(
				"invalid",
				1,
				1,
				"1000",
				1,
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid chain id",
			msg: types.NewMsgGasPriceVoter(
				sample.AccAddress(),
				-1,
				1,
				"1000",
				1,
			),
			err: sdkerrors.ErrInvalidChainID,
		},
		{
			name: "valid address",
			msg: types.NewMsgGasPriceVoter(
				sample.AccAddress(),
				1,
				1,
				"1000",
				1,
			),
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

func TestMsgGasPriceVoter_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgGasPriceVoter
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgGasPriceVoter{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgGasPriceVoter{
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

func TestMsgGasPriceVoter_Type(t *testing.T) {
	msg := types.MsgGasPriceVoter{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, authz.GasPriceVoter.String(), msg.Type())
}

func TestMsgGasPriceVoter_Route(t *testing.T) {
	msg := types.MsgGasPriceVoter{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgGasPriceVoter_GetSignBytes(t *testing.T) {
	msg := types.MsgGasPriceVoter{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
