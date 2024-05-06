package types_test

import (
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/authz"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgVoteOutbound_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgVoteOutbound
		err  error
	}{
		{
			name: "valid message",
			msg: types.NewMsgVoteOutbound(
				sample.AccAddress(),
				sample.String(),
				sample.String(),
				42,
				42,
				math.NewInt(42),
				42,
				math.NewUint(42),
				chains.ReceiveStatus_created,
				42,
				42,
				coin.CoinType_Zeta,
			),
		},
		{
			name: "invalid address",
			msg: types.NewMsgVoteOutbound(
				"invalid_address",
				sample.String(),
				sample.String(),
				42,
				42,
				math.NewInt(42),
				42,
				math.NewUint(42),
				chains.ReceiveStatus_created,
				42,
				42,
				coin.CoinType_Zeta,
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid chain ID",
			msg: types.NewMsgVoteOutbound(
				sample.AccAddress(),
				sample.String(),
				sample.String(),
				42,
				42,
				math.NewInt(42),
				42,
				math.NewUint(42),
				chains.ReceiveStatus_created,
				-1,
				42,
				coin.CoinType_Zeta,
			),
			err: types.ErrInvalidChainID,
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

func TestMsgVoteOutbound_Digest(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	msg := types.MsgVoteOutbound{
		Creator:                           sample.AccAddress(),
		CctxHash:                          sample.String(),
		ObservedOutboundHash:              sample.String(),
		ObservedOutboundBlockHeight:       42,
		ObservedOutboundGasUsed:           42,
		ObservedOutboundEffectiveGasPrice: math.NewInt(42),
		ObservedOutboundEffectiveGasLimit: 42,
		ValueReceived:                     math.NewUint(42),
		Status:                            chains.ReceiveStatus_created,
		OutboundChain:                     42,
		OutboundTssNonce:                  42,
		CoinType:                          coin.CoinType_Zeta,
	}
	hash := msg.Digest()
	require.NotEmpty(t, hash, "hash should not be empty")

	// creator not used
	msg2 := msg
	msg2.Creator = sample.AccAddress()
	hash2 := msg2.Digest()
	require.Equal(t, hash, hash2, "creator should not change hash")

	// status not used
	msg2 = msg
	msg2.Status = chains.ReceiveStatus_failed
	hash2 = msg2.Digest()
	require.Equal(t, hash, hash2, "status should not change hash")

	// cctx hash used
	msg2 = msg
	msg2.CctxHash = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "cctx hash should change hash")

	// observed outbound tx hash used
	msg2 = msg
	msg2.ObservedOutboundHash = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx hash should change hash")

	// observed outbound tx block height used
	msg2 = msg
	msg2.ObservedOutboundBlockHeight = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx block height should change hash")

	// observed outbound tx gas used used
	msg2 = msg
	msg2.ObservedOutboundGasUsed = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx gas used should change hash")

	// observed outbound tx effective gas price used
	msg2 = msg
	msg2.ObservedOutboundEffectiveGasPrice = math.NewInt(43)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx effective gas price should change hash")

	// observed outbound tx effective gas limit used
	msg2 = msg
	msg2.ObservedOutboundEffectiveGasLimit = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx effective gas limit should change hash")

	// zeta minted used
	msg2 = msg
	msg2.ValueReceived = math.NewUint(43)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "zeta minted should change hash")

	// out tx chain used
	msg2 = msg
	msg2.OutboundChain = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "out tx chain should change hash")

	// out tx tss nonce used
	msg2 = msg
	msg2.OutboundTssNonce = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "out tx tss nonce should change hash")

	// coin type used
	msg2 = msg
	msg2.CoinType = coin.CoinType_ERC20
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "coin type should change hash")
}

func TestMsgVoteOutbound_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgVoteOutbound
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgVoteOutbound{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgVoteOutbound{
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

func TestMsgVoteOutbound_Type(t *testing.T) {
	msg := types.MsgVoteOutbound{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, authz.OutboundVoter.String(), msg.Type())
}

func TestMsgVoteOutbound_Route(t *testing.T) {
	msg := types.MsgVoteOutbound{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgVoteOutbound_GetSignBytes(t *testing.T) {
	msg := types.MsgVoteOutbound{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
