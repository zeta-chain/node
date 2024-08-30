package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestNewMsgVoteBlameMsg_ValidateBasic(t *testing.T) {
	keeper.SetConfig(false)
	tests := []struct {
		name  string
		msg   *types.MsgVoteBlame
		error bool
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgVoteBlameMsg(
				"invalid_address",
				1,
				sample.BlameRecordsList(t, 1)[0],
			),
			error: true,
		},
		{
			name: "valid",
			msg: types.NewMsgVoteBlameMsg(
				sample.AccAddress(),
				5,
				sample.BlameRecordsList(t, 1)[0],
			),
			error: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.error {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgVoteBlameMsg_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgVoteBlame
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgVoteBlame{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgVoteBlame{
				Creator: "invalid",
			},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				require.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				require.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestNewMsgVoteBlameMsg_Type(t *testing.T) {
	msg := types.MsgVoteBlame{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgVoteBlame, msg.Type())
}

func TestNewMsgVoteBlameMsg_Route(t *testing.T) {
	msg := types.MsgVoteBlame{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestNewMsgVoteBlameMsg_GetSignBytes(t *testing.T) {
	msg := types.MsgVoteBlame{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}

func TestNewMsgVoteBlameMsg_Digest(t *testing.T) {
	msg := types.MsgVoteBlame{
		Creator: sample.AccAddress(),
	}

	digest := msg.Digest()
	msg.Creator = ""
	expectedDigest := crypto.Keccak256Hash([]byte(msg.String()))
	require.Equal(t, expectedDigest.Hex(), digest)
}
