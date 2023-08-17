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
			name: "invalid msg",
			msg: MsgAddObserver{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid creator",
			msg: MsgAddObserver{
				Creator:                 "invalid_address",
				ObserverAddress:         sample.AccAddress(),
				ZetaclientGranteePubkey: sample.PubKey(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid pubkey",
			msg: MsgAddObserver{
				Creator:                 sample.AccAddress(),
				ObserverAddress:         sample.AccAddress(),
				ZetaclientGranteePubkey: "sample.PubKey()",
			},
			err: sdkerrors.ErrInvalidPubKey,
		},
		{
			name: "invalid observer address",
			msg: MsgAddObserver{
				Creator:                 sample.AccAddress(),
				ObserverAddress:         "invalid_address",
				ZetaclientGranteePubkey: sample.PubKey(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgAddObserver{
				Creator:                 sample.AccAddress(),
				ObserverAddress:         sample.AccAddress(),
				ZetaclientGranteePubkey: sample.PubKey(),
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
