package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestMsgUpdatePolicies_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgUpdatePolicies
		err  error
	}{
		{
			name: "valid message",
			msg:  types.NewMsgUpdatePolicies(sample.AccAddress(), sample.Policies()),
		},
		{
			name: "invalid creator address",
			msg:  types.NewMsgUpdatePolicies("invalid", sample.Policies()),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid policies",
			msg: types.NewMsgUpdatePolicies(sample.AccAddress(), types.Policies{
				Items: []*types.Policy{
					{
						Address:    "invalid",
						PolicyType: types.PolicyType_groupEmergency,
					},
				},
			}),
			err: sdkerrors.ErrInvalidRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
