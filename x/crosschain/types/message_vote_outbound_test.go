package types_test

import (
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/authz"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
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
				types.ConfirmationMode_SAFE,
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
				types.ConfirmationMode_SAFE,
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
				types.ConfirmationMode_SAFE,
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
		ConfirmationMode:                  types.ConfirmationMode_SAFE,
	}
	hash := msg.Digest()
	require.NotEmpty(t, hash, "hash should not be empty")

	// creator not used
	msgNew := msg
	msgNew.Creator = sample.AccAddress()
	hash2 := msgNew.Digest()
	require.Equal(t, hash, hash2, "creator should not change hash")

	// status not used
	msgNew = msg
	msgNew.Status = chains.ReceiveStatus_failed
	hash2 = msgNew.Digest()
	require.Equal(t, hash, hash2, "status should not change hash")

	// cctx hash used
	msgNew = msg
	msgNew.CctxHash = sample.StringRandom(r, 32)
	hash2 = msgNew.Digest()
	require.NotEqual(t, hash, hash2, "cctx hash should change hash")

	// observed outbound tx hash used
	msgNew = msg
	msgNew.ObservedOutboundHash = sample.StringRandom(r, 32)
	hash2 = msgNew.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx hash should change hash")

	// observed outbound tx block height used
	msgNew = msg
	msgNew.ObservedOutboundBlockHeight = 43
	hash2 = msgNew.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx block height should change hash")

	// observed outbound tx gas used used
	msgNew = msg
	msgNew.ObservedOutboundGasUsed = 43
	hash2 = msgNew.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx gas used should change hash")

	// observed outbound tx effective gas price used
	msgNew = msg
	msgNew.ObservedOutboundEffectiveGasPrice = math.NewInt(43)
	hash2 = msgNew.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx effective gas price should change hash")

	// observed outbound tx effective gas limit used
	msgNew = msg
	msgNew.ObservedOutboundEffectiveGasLimit = 43
	hash2 = msgNew.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx effective gas limit should change hash")

	// zeta minted used
	msgNew = msg
	msgNew.ValueReceived = math.NewUint(43)
	hash2 = msgNew.Digest()
	require.NotEqual(t, hash, hash2, "zeta minted should change hash")

	// out tx chain used
	msgNew = msg
	msgNew.OutboundChain = 43
	hash2 = msgNew.Digest()
	require.NotEqual(t, hash, hash2, "out tx chain should change hash")

	// out tx tss nonce used
	msgNew = msg
	msgNew.OutboundTssNonce = 43
	hash2 = msgNew.Digest()
	require.NotEqual(t, hash, hash2, "out tx tss nonce should change hash")

	// coin type used
	msgNew = msg
	msgNew.CoinType = coin.CoinType_ERC20
	hash2 = msgNew.Digest()
	require.NotEqual(t, hash, hash2, "coin type should change hash")

	// confirmation mode not used
	msgNew = msg
	msgNew.ConfirmationMode = types.ConfirmationMode_FAST
	hash2 = msgNew.Digest()
	require.Equal(t, hash, hash2, "confirmation mode should not change hash")
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
