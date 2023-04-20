package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestMsgUpdateClientParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateClientParams
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateClientParams{
				Creator: "invalid_address",
				ClientParams: &ClientParams{
					ConfirmationCount: 1,
					GasPriceTicker:    1,
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateClientParams{
				Creator: sample.AccAddress(),
				ClientParams: &ClientParams{
					ConfirmationCount: 1,
					GasPriceTicker:    1,
				},
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
