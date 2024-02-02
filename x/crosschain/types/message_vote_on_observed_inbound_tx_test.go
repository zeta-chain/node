package types_test

import (
	"testing"

	"math/rand"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgVoteOnObservedInboundTx_ValidateBasic(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	tests := []struct {
		name string
		msg  types.MsgVoteOnObservedInboundTx
		err  error
	}{
		{
			name: "valid message",
			msg: types.MsgVoteOnObservedInboundTx{
				Creator:       sample.AccAddress(),
				Sender:        sample.AccAddress(),
				SenderChainId: 42,
				TxOrigin:      sample.String(),
				Receiver:      sample.String(),
				ReceiverChain: 42,
				Amount:        math.NewUint(42),
				Message:       sample.String(),
				InTxHash:      sample.String(),
				InBlockHeight: 42,
				GasLimit:      42,
				CoinType:      common.CoinType_Zeta,
				Asset:         sample.String(),
				EventIndex:    42,
			},
		},
		{
			name: "invalid address",
			msg: types.MsgVoteOnObservedInboundTx{
				Creator:       "invalid_address",
				Sender:        sample.AccAddress(),
				SenderChainId: 42,
				TxOrigin:      sample.String(),
				Receiver:      sample.String(),
				ReceiverChain: 42,
				Amount:        math.NewUint(42),
				Message:       sample.String(),
				InTxHash:      sample.String(),
				InBlockHeight: 42,
				GasLimit:      42,
				CoinType:      common.CoinType_Zeta,
				Asset:         sample.String(),
				EventIndex:    42,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid sender chain ID",
			msg: types.MsgVoteOnObservedInboundTx{
				Creator:       sample.AccAddress(),
				Sender:        sample.AccAddress(),
				SenderChainId: -1,
				TxOrigin:      sample.String(),
				Receiver:      sample.String(),
				ReceiverChain: 42,
				Amount:        math.NewUint(42),
				Message:       sample.String(),
				InTxHash:      sample.String(),
				InBlockHeight: 42,
				GasLimit:      42,
				CoinType:      common.CoinType_Zeta,
				Asset:         sample.String(),
				EventIndex:    42,
			},
			err: types.ErrInvalidChainID,
		},
		{
			name: "invalid receiver chain ID",
			msg: types.MsgVoteOnObservedInboundTx{
				Creator:       sample.AccAddress(),
				Sender:        sample.AccAddress(),
				SenderChainId: 42,
				TxOrigin:      sample.String(),
				Receiver:      sample.String(),
				ReceiverChain: -1,
				Amount:        math.NewUint(42),
				Message:       sample.String(),
				InTxHash:      sample.String(),
				InBlockHeight: 42,
				GasLimit:      42,
				CoinType:      common.CoinType_Zeta,
				Asset:         sample.String(),
				EventIndex:    42,
			},
			err: types.ErrInvalidChainID,
		},
		{
			name: "invalid message length",
			msg: types.MsgVoteOnObservedInboundTx{
				Creator:       sample.AccAddress(),
				Sender:        sample.AccAddress(),
				SenderChainId: 42,
				TxOrigin:      sample.String(),
				Receiver:      sample.String(),
				ReceiverChain: 42,
				Amount:        math.NewUint(42),
				Message:       sample.StringRandom(r, types.MaxMessageLength+1),
				InTxHash:      sample.String(),
				InBlockHeight: 42,
				GasLimit:      42,
				CoinType:      common.CoinType_Zeta,
				Asset:         sample.String(),
				EventIndex:    42,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestMsgVoteOnObservedInboundTx_Digest(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	msg := types.MsgVoteOnObservedInboundTx{
		Creator:       sample.AccAddress(),
		Sender:        sample.AccAddress(),
		SenderChainId: 42,
		TxOrigin:      sample.String(),
		Receiver:      sample.String(),
		ReceiverChain: 42,
		Amount:        math.NewUint(42),
		Message:       sample.String(),
		InTxHash:      sample.String(),
		InBlockHeight: 42,
		GasLimit:      42,
		CoinType:      common.CoinType_Zeta,
		Asset:         sample.String(),
		EventIndex:    42,
	}
	hash := msg.Digest()
	assert.NotEmpty(t, hash, "hash should not be empty")

	// creator not used
	msg2 := msg
	msg2.Creator = sample.AccAddress()
	hash2 := msg2.Digest()
	assert.Equal(t, hash, hash2, "creator should not change hash")

	// in block height not used
	msg2 = msg
	msg2.InBlockHeight = 43
	hash2 = msg2.Digest()
	assert.Equal(t, hash, hash2, "in block height should not change hash")

	// sender used
	msg2 = msg
	msg2.Sender = sample.AccAddress()
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "sender should change hash")

	// sender chain ID used
	msg2 = msg
	msg2.SenderChainId = 43
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "sender chain ID should change hash")

	// tx origin used
	msg2 = msg
	msg2.TxOrigin = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "tx origin should change hash")

	// receiver used
	msg2 = msg
	msg2.Receiver = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "receiver should change hash")

	// receiver chain ID used
	msg2 = msg
	msg2.ReceiverChain = 43
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "receiver chain ID should change hash")

	// amount used
	msg2 = msg
	msg2.Amount = math.NewUint(43)
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "amount should change hash")

	// message used
	msg2 = msg
	msg2.Message = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "message should change hash")

	// in tx hash used
	msg2 = msg
	msg2.InTxHash = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "in tx hash should change hash")

	// gas limit used
	msg2 = msg
	msg2.GasLimit = 43
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "gas limit should change hash")

	// coin type used
	msg2 = msg
	msg2.CoinType = common.CoinType_ERC20
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "coin type should change hash")

	// asset used
	msg2 = msg
	msg2.Asset = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "asset should change hash")

	// event index used
	msg2 = msg
	msg2.EventIndex = 43
	hash2 = msg2.Digest()
	assert.NotEqual(t, hash, hash2, "event index should change hash")
}
