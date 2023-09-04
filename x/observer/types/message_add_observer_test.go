package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgAddObserver_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgAddObserver
		err  error
	}{
		{
			name: "invalid msg",
			msg: types.MsgAddObserver{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid creator",
			msg: types.MsgAddObserver{
				Creator:                 "invalid_address",
				ObserverAddress:         sample.AccAddress(),
				ZetaclientGranteePubkey: sample.PubKeyString(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid pubkey",
			msg: types.MsgAddObserver{
				Creator:                 sample.AccAddress(),
				ObserverAddress:         sample.AccAddress(),
				ZetaclientGranteePubkey: "sample.PubKey()",
			},
			err: sdkerrors.ErrInvalidPubKey,
		},
		{
			name: "invalid observer address",
			msg: types.MsgAddObserver{
				Creator:                 sample.AccAddress(),
				ObserverAddress:         "invalid_address",
				ZetaclientGranteePubkey: sample.PubKeyString(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.MsgAddObserver{
				Creator:                 sample.AccAddress(),
				ObserverAddress:         sample.AccAddress(),
				ZetaclientGranteePubkey: sample.PubKeyString(),
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
