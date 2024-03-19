package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgResetChainNonces_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgResetChainNonces
		err  error
	}{
		{
			name: "valid message",
			msg: types.MsgResetChainNonces{
				Creator:        sample.AccAddress(),
				ChainId:        common.ExternalChainList()[0].ChainId,
				ChainNonceLow:  1,
				ChainNonceHigh: 5,
			},
		},
		{
			name: "invalid address",
			msg: types.MsgResetChainNonces{
				Creator: "invalid_address",
				ChainId: common.ExternalChainList()[0].ChainId,
			},
			err: sdkerrors.ErrInvalidAddress,
		},

		{
			name: "invalid chain ID",
			msg: types.MsgResetChainNonces{
				Creator: sample.AccAddress(),
				ChainId: 999,
			},
			err: sdkerrors.ErrInvalidChainID,
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

func TestMsgResetChainNonces_Type(t *testing.T) {
	msg := types.NewMsgResetChainNonces(sample.AccAddress(), 5, 1, 5)
	require.Equal(t, types.TypeMsgResetChainNonces, msg.Type())
}

func TestMsgResetChainNonces_Route(t *testing.T) {
	msg := types.NewMsgResetChainNonces(sample.AccAddress(), 5, 1, 5)
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgResetChainNonces_GetSignBytes(t *testing.T) {
	msg := types.NewMsgResetChainNonces(sample.AccAddress(), 5, 1, 5)
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
