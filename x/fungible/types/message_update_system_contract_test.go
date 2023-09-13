package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestMsgUpdateSystemContract_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateSystemContract
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateSystemContract{
				Creator:                  "invalid_address",
				NewSystemContractAddress: sample.EthAddress().String(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid new system contract address",
			msg: types.MsgUpdateSystemContract{
				Creator:                  sample.AccAddress(),
				NewSystemContractAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid message",
			msg: types.MsgUpdateSystemContract{
				Creator:                  sample.AccAddress(),
				NewSystemContractAddress: sample.EthAddress().String(),
			},
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
