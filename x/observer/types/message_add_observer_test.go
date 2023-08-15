package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestMsgAddObserver_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddObserver
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAddObserver{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg:  GenerateObserverMsg(),
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

func GenerateObserverMsg() MsgAddObserver {
	address, pubkey := sample.PubKey()
	return MsgAddObserver{
		Creator:                  sample.AccAddress(),
		ObserverAddress:          sample.AccAddress(),
		ZetaclientGranteeAddress: address,
		ZetaclientGranteePubkey:  pubkey,
	}
}
