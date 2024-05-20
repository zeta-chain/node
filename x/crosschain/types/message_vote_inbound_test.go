package types_test

import (
	"testing"

	"math/rand"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/authz"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgVoteInbound_ValidateBasic(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	tests := []struct {
		name string
		msg  *types.MsgVoteInbound
		err  error
	}{
		{
			name: "valid message",
			msg: types.NewMsgVoteInbound(
				sample.AccAddress(),
				sample.AccAddress(),
				42,
				sample.String(),
				sample.String(),
				42,
				math.NewUint(42),
				sample.String(),
				sample.String(),
				42,
				42,
				coin.CoinType_Zeta,
				sample.String(),
				42,
			),
		},
		{
			name: "invalid address",
			msg: types.NewMsgVoteInbound(
				"invalid_address",
				sample.AccAddress(),
				42,
				sample.String(),
				sample.String(),
				42,
				math.NewUint(42),
				sample.String(),
				sample.String(),
				42,
				42,
				coin.CoinType_Zeta,
				sample.String(),
				42,
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid sender chain ID",
			msg: types.NewMsgVoteInbound(
				sample.AccAddress(),
				sample.AccAddress(),
				-1,
				sample.String(),
				sample.String(),
				42,
				math.NewUint(42),
				sample.String(),
				sample.String(),
				42,
				42,
				coin.CoinType_Zeta,
				sample.String(),
				42,
			),
			err: types.ErrInvalidChainID,
		},
		{
			name: "invalid receiver chain ID",
			msg: types.NewMsgVoteInbound(
				sample.AccAddress(),
				sample.AccAddress(),
				42,
				sample.String(),
				sample.String(),
				-1,
				math.NewUint(42),
				sample.String(),
				sample.String(),
				42,
				42,
				coin.CoinType_Zeta,
				sample.String(),
				42,
			),
			err: types.ErrInvalidChainID,
		},
		{
			name: "invalid message length",
			msg: types.NewMsgVoteInbound(
				sample.AccAddress(),
				sample.AccAddress(),
				42,
				sample.String(),
				sample.String(),
				42,
				math.NewUint(42),
				sample.StringRandom(r, types.MaxMessageLength+1),
				sample.String(),
				42,
				42,
				coin.CoinType_Zeta,
				sample.String(),
				42,
			),
			err: sdkerrors.ErrInvalidRequest,
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

func TestMsgVoteInbound_Digest(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	msg := types.MsgVoteInbound{
		Creator:            sample.AccAddress(),
		Sender:             sample.AccAddress(),
		SenderChainId:      42,
		TxOrigin:           sample.String(),
		Receiver:           sample.String(),
		ReceiverChain:      42,
		Amount:             math.NewUint(42),
		Message:            sample.String(),
		InboundHash:        sample.String(),
		InboundBlockHeight: 42,
		GasLimit:           42,
		CoinType:           coin.CoinType_Zeta,
		Asset:              sample.String(),
		EventIndex:         42,
	}
	hash := msg.Digest()
	require.NotEmpty(t, hash, "hash should not be empty")

	// creator not used
	msg2 := msg
	msg2.Creator = sample.AccAddress()
	hash2 := msg2.Digest()
	require.Equal(t, hash, hash2, "creator should not change hash")

	// in block height not used
	msg2 = msg
	msg2.InboundBlockHeight = 43
	hash2 = msg2.Digest()
	require.Equal(t, hash, hash2, "in block height should not change hash")

	// sender used
	msg2 = msg
	msg2.Sender = sample.AccAddress()
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "sender should change hash")

	// sender chain ID used
	msg2 = msg
	msg2.SenderChainId = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "sender chain ID should change hash")

	// tx origin used
	msg2 = msg
	msg2.TxOrigin = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "tx origin should change hash")

	// receiver used
	msg2 = msg
	msg2.Receiver = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "receiver should change hash")

	// receiver chain ID used
	msg2 = msg
	msg2.ReceiverChain = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "receiver chain ID should change hash")

	// amount used
	msg2 = msg
	msg2.Amount = math.NewUint(43)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "amount should change hash")

	// message used
	msg2 = msg
	msg2.Message = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "message should change hash")

	// in tx hash used
	msg2 = msg
	msg2.InboundHash = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "in tx hash should change hash")

	// gas limit used
	msg2 = msg
	msg2.GasLimit = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "gas limit should change hash")

	// coin type used
	msg2 = msg
	msg2.CoinType = coin.CoinType_ERC20
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "coin type should change hash")

	// asset used
	msg2 = msg
	msg2.Asset = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "asset should change hash")

	// event index used
	msg2 = msg
	msg2.EventIndex = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "event index should change hash")
}

func TestMsgVoteInbound_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgVoteInbound
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgVoteInbound{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgVoteInbound{
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

func TestMsgVoteInbound_Type(t *testing.T) {
	msg := types.MsgVoteInbound{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, authz.InboundVoter.String(), msg.Type())
}

func TestMsgVoteInbound_Route(t *testing.T) {
	msg := types.MsgVoteInbound{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgVoteInbound_GetSignBytes(t *testing.T) {
	msg := types.MsgVoteInbound{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
