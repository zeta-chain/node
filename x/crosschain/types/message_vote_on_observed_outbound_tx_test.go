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

func TestMsgVoteOnObservedOutboundTx_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgVoteOnObservedOutboundTx
		err  error
	}{
		{
			name: "valid message",
			msg: types.NewMsgVoteOnObservedOutboundTx(
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
			msg: types.NewMsgVoteOnObservedOutboundTx(
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
			msg: types.NewMsgVoteOnObservedOutboundTx(
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

func TestMsgVoteOnObservedOutboundTx_Digest(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	msg := types.MsgVoteOnObservedOutboundTx{
		Creator:                        sample.AccAddress(),
		CctxHash:                       sample.String(),
		ObservedOutTxHash:              sample.String(),
		ObservedOutTxBlockHeight:       42,
		ObservedOutTxGasUsed:           42,
		ObservedOutTxEffectiveGasPrice: math.NewInt(42),
		ObservedOutTxEffectiveGasLimit: 42,
		ValueReceived:                  math.NewUint(42),
		Status:                         chains.ReceiveStatus_created,
		OutTxChain:                     42,
		OutTxTssNonce:                  42,
		CoinType:                       coin.CoinType_Zeta,
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
	msg2.ObservedOutTxHash = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx hash should change hash")

	// observed outbound tx block height used
	msg2 = msg
	msg2.ObservedOutTxBlockHeight = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx block height should change hash")

	// observed outbound tx gas used used
	msg2 = msg
	msg2.ObservedOutTxGasUsed = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx gas used should change hash")

	// observed outbound tx effective gas price used
	msg2 = msg
	msg2.ObservedOutTxEffectiveGasPrice = math.NewInt(43)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx effective gas price should change hash")

	// observed outbound tx effective gas limit used
	msg2 = msg
	msg2.ObservedOutTxEffectiveGasLimit = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "observed outbound tx effective gas limit should change hash")

	// zeta minted used
	msg2 = msg
	msg2.ValueReceived = math.NewUint(43)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "zeta minted should change hash")

	// out tx chain used
	msg2 = msg
	msg2.OutTxChain = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "out tx chain should change hash")

	// out tx tss nonce used
	msg2 = msg
	msg2.OutTxTssNonce = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "out tx tss nonce should change hash")

	// coin type used
	msg2 = msg
	msg2.CoinType = coin.CoinType_ERC20
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "coin type should change hash")
}

func TestMsgVoteOnObservedOutboundTx_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgVoteOnObservedOutboundTx
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgVoteOnObservedOutboundTx{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgVoteOnObservedOutboundTx{
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

func TestMsgVoteOnObservedOutboundTx_Type(t *testing.T) {
	msg := types.MsgVoteOnObservedOutboundTx{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, authz.OutboundVoter.String(), msg.Type())
}

func TestMsgVoteOnObservedOutboundTx_Route(t *testing.T) {
	msg := types.MsgVoteOnObservedOutboundTx{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgVoteOnObservedOutboundTx_GetSignBytes(t *testing.T) {
	msg := types.MsgVoteOnObservedOutboundTx{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
