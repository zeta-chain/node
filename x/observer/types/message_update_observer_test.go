package types_test

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"testing"
)

func TestNewMsgUpdateObserver(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateObserver
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgUpdateObserver{
				Creator:            "invalid_address",
				OldObserverAddress: sample.AccAddress(),
				NewObserverAddress: sample.AccAddress(),
				UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid old observer address",
			msg: types.MsgUpdateObserver{
				Creator:            sample.AccAddress(),
				OldObserverAddress: "invalid_address",
				NewObserverAddress: sample.AccAddress(),
				UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid new observer address",
			msg: types.MsgUpdateObserver{
				Creator:            sample.AccAddress(),
				OldObserverAddress: sample.AccAddress(),
				NewObserverAddress: "invalid_address",
				UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "old observer address is not creator",
			msg: types.MsgUpdateObserver{
				Creator:            sample.AccAddress(),
				OldObserverAddress: sample.AccAddress(),
				NewObserverAddress: sample.AccAddress(),
				UpdateReason:       types.ObserverUpdateReason_Tombstoned,
			},
			err: types.ErrUpdateObserver,
		},
		{
			name: "valid message",
			msg: types.MsgUpdateObserver{
				Creator:            sample.AccAddress(),
				OldObserverAddress: sample.AccAddress(),
				NewObserverAddress: sample.AccAddress(),
				UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
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
