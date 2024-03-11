package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
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

func TestMsgUpdatePolicies_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgUpdatePolicies
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgUpdatePolicies(signer, sample.Policies()),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgUpdatePolicies("invalid", sample.Policies()),
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				assert.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				assert.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestMsgUpdatePolicies_Type(t *testing.T) {
	msg := types.NewMsgUpdatePolicies(sample.AccAddress(), sample.Policies())
	assert.Equal(t, types.TypeMsgUpdatePolicies, msg.Type())
}

func TestMsgUpdatePolicies_Route(t *testing.T) {
	msg := types.NewMsgUpdatePolicies(sample.AccAddress(), sample.Policies())
	assert.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdatePolicies_GetSignBytes(t *testing.T) {
	msg := types.NewMsgUpdatePolicies(sample.AccAddress(), sample.Policies())
	assert.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
