package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgCreateTSSVoter_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgCreateTSSVoter
		err  error
	}{
		{
			name: "valid message",
			msg:  types.NewMsgCreateTSSVoter(sample.AccAddress(), "pubkey", 1, common.ReceiveStatus_Created),
		},
		{
			name: "invalid creator address",
			msg:  types.NewMsgCreateTSSVoter("invalid", "pubkey", 1, common.ReceiveStatus_Created),
			err:  sdkerrors.ErrInvalidAddress,
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

func TestMsgCreateTSSVoter_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgCreateTSSVoter
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgCreateTSSVoter(signer, "pubkey", 1, common.ReceiveStatus_Created),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgCreateTSSVoter("invalid", "pubkey", 1, common.ReceiveStatus_Created),
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

func TestMsgCreateTSSVoter_Type(t *testing.T) {
	msg := types.NewMsgCreateTSSVoter(sample.AccAddress(), "pubkey", 1, common.ReceiveStatus_Created)
	assert.Equal(t, types.TypeMsgCreateTSSVoter, msg.Type())
}

func TestMsgCreateTSSVoter_Route(t *testing.T) {
	msg := types.NewMsgCreateTSSVoter(sample.AccAddress(), "pubkey", 1, common.ReceiveStatus_Created)
	assert.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgCreateTSSVoter_GetSignBytes(t *testing.T) {
	msg := types.NewMsgCreateTSSVoter(sample.AccAddress(), "pubkey", 1, common.ReceiveStatus_Created)
	assert.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}

func TestMsgCreateTSSVoter_Digest(t *testing.T) {
	msg := types.NewMsgCreateTSSVoter(sample.AccAddress(), "pubkey", 1, common.ReceiveStatus_Created)
	assert.Equal(t, "1-tss-keygen", msg.Digest())
}
