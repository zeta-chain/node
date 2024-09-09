package types_test

import (
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"
	"math/big"
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/authz"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestNewMsgVoteInbound(t *testing.T) {
	t.Run("empty revert options by default", func(t *testing.T) {
		msg := types.NewMsgVoteInbound(
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
			types.ProtocolContractVersion_V1,
		)
		require.EqualValues(t, types.NewEmptyRevertOptions(), msg.RevertOptions)
	})

	t.Run("can set ZEVM revert options", func(t *testing.T) {
		revertAddress := sample.EthAddress()
		abortAddress := sample.EthAddress()
		revertMessage := sample.Bytes()

		msg := types.NewMsgVoteInbound(
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
			types.ProtocolContractVersion_V1,
			types.WithZEVMRevertOptions(gatewayzevm.RevertOptions{
				RevertAddress:    revertAddress,
				CallOnRevert:     true,
				AbortAddress:     abortAddress,
				RevertMessage:    revertMessage,
				OnRevertGasLimit: big.NewInt(1000),
			}),
		)
		require.EqualValues(t, types.RevertOptions{
			RevertAddress:  revertAddress.Hex(),
			CallOnRevert:   true,
			AbortAddress:   abortAddress.Hex(),
			RevertMessage:  revertMessage,
			RevertGasLimit: math.NewUint(1000),
		}, msg.RevertOptions)

		// if revertGasLimit not specified, it should be zero
		msg = types.NewMsgVoteInbound(
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
			types.ProtocolContractVersion_V1,
			types.WithZEVMRevertOptions(gatewayzevm.RevertOptions{
				RevertAddress: revertAddress,
				CallOnRevert:  true,
				AbortAddress:  abortAddress,
				RevertMessage: revertMessage,
			}),
		)
		require.EqualValues(t, types.RevertOptions{
			RevertAddress:  revertAddress.Hex(),
			CallOnRevert:   true,
			AbortAddress:   abortAddress.Hex(),
			RevertMessage:  revertMessage,
			RevertGasLimit: math.ZeroUint(),
		}, msg.RevertOptions)
	})

	t.Run("can set EVM revert options", func(t *testing.T) {
		revertAddress := sample.EthAddress()
		abortAddress := sample.EthAddress()
		revertMessage := sample.Bytes()

		msg := types.NewMsgVoteInbound(
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
			types.ProtocolContractVersion_V1,
			types.WithEVMRevertOptions(gatewayevm.RevertOptions{
				RevertAddress:    revertAddress,
				CallOnRevert:     true,
				AbortAddress:     abortAddress,
				RevertMessage:    revertMessage,
				OnRevertGasLimit: big.NewInt(1000),
			}),
		)
		require.EqualValues(t, types.RevertOptions{
			RevertAddress:  revertAddress.Hex(),
			CallOnRevert:   true,
			AbortAddress:   abortAddress.Hex(),
			RevertMessage:  revertMessage,
			RevertGasLimit: math.NewUint(1000),
		}, msg.RevertOptions)

		msg = types.NewMsgVoteInbound(
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
			types.ProtocolContractVersion_V1,
			types.WithEVMRevertOptions(gatewayevm.RevertOptions{
				RevertAddress: revertAddress,
				CallOnRevert:  true,
				AbortAddress:  abortAddress,
				RevertMessage: revertMessage,
			}),
		)
		require.EqualValues(t, types.RevertOptions{
			RevertAddress:  revertAddress.Hex(),
			CallOnRevert:   true,
			AbortAddress:   abortAddress.Hex(),
			RevertMessage:  revertMessage,
			RevertGasLimit: math.ZeroUint(),
		}, msg.RevertOptions)
	})
}

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
				types.ProtocolContractVersion_V1,
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
				types.ProtocolContractVersion_V1,
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
				types.ProtocolContractVersion_V1,
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
				types.ProtocolContractVersion_V1,
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
				types.ProtocolContractVersion_V1,
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
		Creator:                 sample.AccAddress(),
		Sender:                  sample.AccAddress(),
		SenderChainId:           42,
		TxOrigin:                sample.String(),
		Receiver:                sample.String(),
		ReceiverChain:           42,
		Amount:                  math.NewUint(42),
		Message:                 sample.String(),
		InboundHash:             sample.String(),
		InboundBlockHeight:      42,
		GasLimit:                42,
		CoinType:                coin.CoinType_Zeta,
		Asset:                   sample.String(),
		EventIndex:              42,
		ProtocolContractVersion: types.ProtocolContractVersion_V1,
	}
	hash := msg.Digest()
	require.NotEmpty(t, hash, "hash should not be empty")

	// creator not used
	msg = msg
	msg.Creator = sample.AccAddress()
	hash2 := msg.Digest()
	require.Equal(t, hash, hash2, "creator should not change hash")

	// in block height not used
	msg = msg
	msg.InboundBlockHeight = 43
	hash2 = msg.Digest()
	require.Equal(t, hash, hash2, "in block height should not change hash")

	// sender used
	msg = msg
	msg.Sender = sample.AccAddress()
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "sender should change hash")

	// sender chain ID used
	msg = msg
	msg.SenderChainId = 43
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "sender chain ID should change hash")

	// tx origin used
	msg = msg
	msg.TxOrigin = sample.StringRandom(r, 32)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "tx origin should change hash")

	// receiver used
	msg = msg
	msg.Receiver = sample.StringRandom(r, 32)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "receiver should change hash")

	// receiver chain ID used
	msg = msg
	msg.ReceiverChain = 43
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "receiver chain ID should change hash")

	// amount used
	msg = msg
	msg.Amount = math.NewUint(43)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "amount should change hash")

	// message used
	msg = msg
	msg.Message = sample.StringRandom(r, 32)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "message should change hash")

	// in tx hash used
	msg = msg
	msg.InboundHash = sample.StringRandom(r, 32)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "in tx hash should change hash")

	// gas limit used
	msg = msg
	msg.GasLimit = 43
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "gas limit should change hash")

	// coin type used
	msg = msg
	msg.CoinType = coin.CoinType_ERC20
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "coin type should change hash")

	// asset used
	msg = msg
	msg.Asset = sample.StringRandom(r, 32)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "asset should change hash")

	// event index used
	msg = msg
	msg.EventIndex = 43
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "event index should change hash")

	// protocol contract version used
	msg = msg
	msg.ProtocolContractVersion = types.ProtocolContractVersion_V2
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "protocol contract version should change hash")
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
