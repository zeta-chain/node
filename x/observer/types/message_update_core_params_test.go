package types_test

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"testing"
)

func TestMsgUpdateCoreParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateCoreParams
		err  error
	}{
		{
			name: "valid message",
			msg: types.MsgUpdateCoreParams{
				Creator:    sample.AccAddress(),
				CoreParams: sample.CoreParams(common.ExternalChainList()[0].ChainId),
			},
		},
		{
			name: "invalid address",
			msg: types.MsgUpdateCoreParams{
				Creator:    "invalid_address",
				CoreParams: sample.CoreParams(common.ExternalChainList()[0].ChainId),
			},
			err: sdkerrors.ErrInvalidAddress,
		},

		{
			name: "invalid core params (nil)",
			msg: types.MsgUpdateCoreParams{
				Creator: sample.AccAddress(),
			},
			err: types.ErrInvalidCoreParams,
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
