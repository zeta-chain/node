package types_test

import (
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgVoteOnObservedOutboundTx_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgVoteOnObservedOutboundTx
		err  error
	}{
		{
			name: "valid message",
			msg: types.MsgVoteOnObservedOutboundTx{
				Creator:                        sample.AccAddress(),
				CctxHash:                       sample.String(),
				ObservedOutTxHash:              sample.String(),
				ObservedOutTxBlockHeight:       42,
				ObservedOutTxGasUsed:           42,
				ObservedOutTxEffectiveGasPrice: math.NewInt(42),
				ZetaMinted:                     math.NewUint(42),
				Status:                         common.ReceiveStatus_Created,
				OutTxChain:                     42,
				OutTxTssNonce:                  42,
				CoinType:                       common.CoinType_Zeta,
			},
		},
		{
			name: "effective gas price can be nil",
			msg: types.MsgVoteOnObservedOutboundTx{
				Creator:                  sample.AccAddress(),
				CctxHash:                 sample.String(),
				ObservedOutTxHash:        sample.String(),
				ObservedOutTxBlockHeight: 42,
				ObservedOutTxGasUsed:     42,
				ZetaMinted:               math.NewUint(42),
				Status:                   common.ReceiveStatus_Created,
				OutTxChain:               42,
				OutTxTssNonce:            42,
				CoinType:                 common.CoinType_Zeta,
			},
		},
		{
			name: "invalid address",
			msg: types.MsgVoteOnObservedOutboundTx{
				Creator:                        "invalid_address",
				CctxHash:                       sample.String(),
				ObservedOutTxHash:              sample.String(),
				ObservedOutTxBlockHeight:       42,
				ObservedOutTxGasUsed:           42,
				ObservedOutTxEffectiveGasPrice: math.NewInt(42),
				ZetaMinted:                     math.NewUint(42),
				Status:                         common.ReceiveStatus_Created,
				OutTxChain:                     42,
				OutTxTssNonce:                  42,
				CoinType:                       common.CoinType_Zeta,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid chain ID",
			msg: types.MsgVoteOnObservedOutboundTx{
				Creator:                        sample.AccAddress(),
				CctxHash:                       sample.String(),
				ObservedOutTxHash:              sample.String(),
				ObservedOutTxBlockHeight:       42,
				ObservedOutTxGasUsed:           42,
				ObservedOutTxEffectiveGasPrice: math.NewInt(42),
				ZetaMinted:                     math.NewUint(42),
				Status:                         common.ReceiveStatus_Created,
				OutTxChain:                     -1,
				OutTxTssNonce:                  42,
				CoinType:                       common.CoinType_Zeta,
			},
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
		ZetaMinted:                     math.NewUint(42),
		Status:                         common.ReceiveStatus_Created,
		OutTxChain:                     42,
		OutTxTssNonce:                  42,
		CoinType:                       common.CoinType_Zeta,
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
	msg2.Status = common.ReceiveStatus_Failed
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

	// zeta minted used
	msg2 = msg
	msg2.ZetaMinted = math.NewUint(43)
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
	msg2.CoinType = common.CoinType_ERC20
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "coin type should change hash")
}
