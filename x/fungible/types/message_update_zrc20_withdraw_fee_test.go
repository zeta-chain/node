package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestMsgUpdateZRC20WithdrawFee_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateZRC20WithdrawFee
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateZRC20WithdrawFee{
				Creator:        "invalid_address",
				Zrc20Address:   sample.EthAddress().String(),
				NewWithdrawFee: sdk.NewUint(1),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid new system contract address",
			msg: MsgUpdateZRC20WithdrawFee{
				Creator:        sample.AccAddress(),
				Zrc20Address:   "invalid_address",
				NewWithdrawFee: sdk.NewUint(1),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid new withdraw fee",
			msg: MsgUpdateZRC20WithdrawFee{
				Creator:      sample.AccAddress(),
				Zrc20Address: sample.EthAddress().String(),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid message",
			msg: MsgUpdateZRC20WithdrawFee{
				Creator:        sample.AccAddress(),
				Zrc20Address:   sample.EthAddress().String(),
				NewWithdrawFee: sdk.NewUint(1),
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
