package types_test

import (
	"testing"

	"cosmossdk.io/math"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestNewMsgUpdateZRC20LiquidityCap_ValidateBasics(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateZRC20LiquidityCap
		err  error
	}{
		{
			name: "valid message",
			msg: types.MsgUpdateZRC20LiquidityCap{
				Creator:      sample.AccAddress(),
				Zrc20Address: sample.EthAddress().String(),
				LiquidityCap: math.NewUint(1000),
			},
		},
		{
			name: "valid message with liquidity cap 0",
			msg: types.MsgUpdateZRC20LiquidityCap{
				Creator:      sample.AccAddress(),
				Zrc20Address: sample.EthAddress().String(),
				LiquidityCap: math.ZeroUint(),
			},
		},
		{
			name: "valid message with liquidity cap nil",
			msg: types.MsgUpdateZRC20LiquidityCap{
				Creator:      sample.AccAddress(),
				Zrc20Address: sample.EthAddress().String(),
			},
		},
		{
			name: "invalid address",
			msg: types.MsgUpdateZRC20LiquidityCap{
				Creator:      "invalid_address",
				Zrc20Address: sample.EthAddress().String(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid contract address",
			msg: types.MsgUpdateZRC20LiquidityCap{
				Creator:      sample.AccAddress(),
				Zrc20Address: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
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
