package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgAbortStuckCCTX_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgAbortStuckCCTX
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgAbortStuckCCTX{
				Creator:   "invalid_address",
				CctxIndex: "cctx_index",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: types.MsgAbortStuckCCTX{
				Creator:   sample.AccAddress(),
				CctxIndex: "cctx_index",
			},
			err: nil,
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
