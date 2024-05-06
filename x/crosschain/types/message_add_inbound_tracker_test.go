package types_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgAddInboundTracker_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgAddInboundTracker
		err  error
	}{
		{
			name: "invalid address",
			msg: types.NewMsgAddInboundTracker(
				"invalid_address",
				chains.GoerliChain.ChainId,
				coin.CoinType_Gas,
				"hash",
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid chain id",
			msg: types.NewMsgAddInboundTracker(
				sample.AccAddress(),
				42,
				coin.CoinType_Gas,
				"hash",
			),
			err: errorsmod.Wrapf(types.ErrInvalidChainID, "chain id (%d)", 42),
		},
		{
			name: "invalid proof",
			msg: &types.MsgAddInboundTracker{
				Creator:  sample.AccAddress(),
				ChainId:  chains.ZetaTestnetChain.ChainId,
				CoinType: coin.CoinType_Gas,
				Proof:    &proofs.Proof{},
			},
			err: errorsmod.Wrapf(types.ErrProofVerificationFail, "chain id %d does not support proof-based trackers", chains.ZetaTestnetChain.ChainId),
		},
		{
			name: "invalid coin type",
			msg: &types.MsgAddInboundTracker{
				Creator:  sample.AccAddress(),
				ChainId:  chains.ZetaTestnetChain.ChainId,
				CoinType: 5,
			},
			err: errorsmod.Wrapf(types.ErrProofVerificationFail, "coin-type not supported"),
		},
		{
			name: "valid",
			msg: types.NewMsgAddInboundTracker(
				sample.AccAddress(),
				chains.GoerliChain.ChainId,
				coin.CoinType_Gas,
				"hash",
			),
			err: nil,
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

func TestMsgAddInboundTracker_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgAddInboundTracker
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.NewMsgAddInboundTracker(
				signer,
				chains.GoerliChain.ChainId,
				coin.CoinType_Gas,
				"hash",
			),
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.NewMsgAddInboundTracker(
				"invalid_address",
				chains.GoerliChain.ChainId,
				coin.CoinType_Gas,
				"hash",
			),
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

func TestMsgAddInboundTracker_Type(t *testing.T) {
	msg := types.NewMsgAddInboundTracker(
		sample.AccAddress(),
		chains.GoerliChain.ChainId,
		coin.CoinType_Gas,
		"hash",
	)
	require.Equal(t, types.TypeMsgAddInboundTracker, msg.Type())
}

func TestMsgAddInboundTracker_Route(t *testing.T) {
	msg := types.NewMsgAddInboundTracker(
		sample.AccAddress(),
		chains.GoerliChain.ChainId,
		coin.CoinType_Gas,
		"hash",
	)
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgAddInboundTracker_GetSignBytes(t *testing.T) {
	msg := types.NewMsgAddInboundTracker(
		sample.AccAddress(),
		chains.GoerliChain.ChainId,
		coin.CoinType_Gas,
		"hash",
	)
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
