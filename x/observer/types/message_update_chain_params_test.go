package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgUpdateChainParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateChainParams
		err  error
	}{
		{
			name: "valid message",
			msg: types.MsgUpdateChainParams{
				Creator:     sample.AccAddress(),
				ChainParams: sample.ChainParams(common.ExternalChainList()[0].ChainId),
			},
		},
		{
			name: "invalid address",
			msg: types.MsgUpdateChainParams{
				Creator:     "invalid_address",
				ChainParams: sample.ChainParams(common.ExternalChainList()[0].ChainId),
			},
			err: sdkerrors.ErrInvalidAddress,
		},

		{
			name: "invalid chain params (nil)",
			msg: types.MsgUpdateChainParams{
				Creator: sample.AccAddress(),
			},
			err: types.ErrInvalidChainParams,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
