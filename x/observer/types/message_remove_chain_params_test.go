package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgRemoveChainParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgRemoveChainParams
		err  error
	}{
		{
			name: "valid message",
			msg: types.MsgRemoveChainParams{
				Creator: sample.AccAddress(),
				ChainId: common.ExternalChainList()[0].ChainId,
			},
		},
		{
			name: "invalid address",
			msg: types.MsgRemoveChainParams{
				Creator: "invalid_address",
				ChainId: common.ExternalChainList()[0].ChainId,
			},
			err: sdkerrors.ErrInvalidAddress,
		},

		{
			name: "invalid chain ID",
			msg: types.MsgRemoveChainParams{
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
