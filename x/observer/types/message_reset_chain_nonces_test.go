package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgResetChainNonces_ValidateBasic(t *testing.T) {
	chainList := chains.ExternalChainList([]chains.Chain{})

	tests := []struct {
		name    string
		msg     types.MsgResetChainNonces
		wantErr bool
	}{
		{
			name: "valid message chain nonce high greater than nonce low",
			msg: types.MsgResetChainNonces{
				Creator:        sample.AccAddress(),
				ChainId:        chainList[0].ChainId,
				ChainNonceLow:  1,
				ChainNonceHigh: 5,
			},
			wantErr: false,
		},
		{
			name: "valid message chain nonce high same as nonce low",
			msg: types.MsgResetChainNonces{
				Creator:        sample.AccAddress(),
				ChainId:        chainList[0].ChainId,
				ChainNonceLow:  1,
				ChainNonceHigh: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid address",
			msg: types.MsgResetChainNonces{
				Creator: "invalid_address",
				ChainId: chainList[0].ChainId,
			},
			wantErr: true,
		},
		{
			name: "invalid chain nonce low",
			msg: types.MsgResetChainNonces{
				Creator:       sample.AccAddress(),
				ChainId:       chainList[0].ChainId,
				ChainNonceLow: -1,
			},
			wantErr: true,
		},
		{
			name: "invalid chain nonce high",
			msg: types.MsgResetChainNonces{
				Creator:        sample.AccAddress(),
				ChainId:        chainList[0].ChainId,
				ChainNonceLow:  1,
				ChainNonceHigh: -1,
			},
			wantErr: true,
		},
		{
			name: "invalid chain nonce low greater than chain nonce high",
			msg: types.MsgResetChainNonces{
				Creator:        sample.AccAddress(),
				ChainId:        chainList[0].ChainId,
				ChainNonceLow:  1,
				ChainNonceHigh: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgResetChainNonces_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgResetChainNonces
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgResetChainNonces(signer, 5, 1, 5),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgResetChainNonces("invalid", 5, 1, 5),
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

func TestMsgResetChainNonces_Type(t *testing.T) {
	msg := types.NewMsgResetChainNonces(sample.AccAddress(), 5, 1, 5)
	require.Equal(t, types.TypeMsgResetChainNonces, msg.Type())
}

func TestMsgResetChainNonces_Route(t *testing.T) {
	msg := types.NewMsgResetChainNonces(sample.AccAddress(), 5, 1, 5)
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgResetChainNonces_GetSignBytes(t *testing.T) {
	msg := types.NewMsgResetChainNonces(sample.AccAddress(), 5, 1, 5)
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
