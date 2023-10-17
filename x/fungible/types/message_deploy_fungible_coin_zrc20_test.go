package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestMsgDeployFungibleCoinZRC4_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgDeployFungibleCoinZRC20
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgDeployFungibleCoinZRC20{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid gas limit",
			msg: types.MsgDeployFungibleCoinZRC20{
				Creator:  sample.AccAddress(),
				GasLimit: -1,
			},
			err: sdkerrors.ErrInvalidGasLimit,
		},
		{
			name: "invalid decimals",
			msg: types.MsgDeployFungibleCoinZRC20{
				Creator:  sample.AccAddress(),
				Decimals: 256,
			},
			err: sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "decimals must be less than 256"),
		},
		{
			name: "valid message",
			msg: types.MsgDeployFungibleCoinZRC20{
				Creator: sample.AccAddress(),
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
