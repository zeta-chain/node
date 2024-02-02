package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
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
			name: "invalid Update Reason",
			msg: types.MsgUpdateObserver{
				Creator:            sample.AccAddress(),
				OldObserverAddress: sample.AccAddress(),
				NewObserverAddress: sample.AccAddress(),
				UpdateReason:       types.ObserverUpdateReason(100),
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
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
