package types_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

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
			true,
			types.InboundStatus_SUCCESS,
			types.ConfirmationMode_SAFE,
		)
		require.EqualValues(t, types.NewEmptyRevertOptions(), msg.RevertOptions)
	})

	t.Run("can set revert options", func(t *testing.T) {
		revertAddress := sample.EthAddress()
		abortAddress := sample.EthAddress()
		revertMessage := sample.Bytes()

		msg := types.NewMsgVoteInbound(
			sample.AccAddress(),
			sample.AccAddress(),
			31,
			sample.String(),
			sample.String(),
			31,
			math.NewUint(31),
			sample.String(),
			sample.String(),
			31,
			31,
			coin.CoinType_Gas,
			sample.String(),
			31,
			types.ProtocolContractVersion_V2,
			true,
			types.InboundStatus_SUCCESS,
			types.ConfirmationMode_SAFE,
			types.WithRevertOptions(types.RevertOptions{
				RevertAddress:  revertAddress.Hex(),
				CallOnRevert:   true,
				AbortAddress:   abortAddress.Hex(),
				RevertMessage:  revertMessage,
				RevertGasLimit: math.NewUint(21000),
			}),
		)
		require.EqualValues(t, types.RevertOptions{
			RevertAddress:  revertAddress.Hex(),
			CallOnRevert:   true,
			AbortAddress:   abortAddress.Hex(),
			RevertMessage:  revertMessage,
			RevertGasLimit: math.NewUint(21000),
		}, msg.RevertOptions)
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
			true,
			types.InboundStatus_SUCCESS,
			types.ConfirmationMode_SAFE,
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
			true,
			types.InboundStatus_SUCCESS,
			types.ConfirmationMode_SAFE,
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
			true,
			types.InboundStatus_SUCCESS,
			types.ConfirmationMode_SAFE,
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
			true,
			types.InboundStatus_SUCCESS,
			types.ConfirmationMode_SAFE,
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

	t.Run("can set is cross chain call options", func(t *testing.T) {
		// false by default
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
			true,
			types.InboundStatus_SUCCESS,
			types.ConfirmationMode_SAFE,
		)
		require.False(t, msg.IsCrossChainCall)

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
			true,
			types.InboundStatus_SUCCESS,
			types.ConfirmationMode_SAFE,
			types.WithCrossChainCall(true),
		)
		require.True(t, msg.IsCrossChainCall)

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
			true,
			types.InboundStatus_SUCCESS,
			types.ConfirmationMode_SAFE,
			types.WithCrossChainCall(false),
		)
		require.False(t, msg.IsCrossChainCall)
	})

	t.Run("can set inbound status and confirmation mode", func(t *testing.T) {
		expectedInboundStatus := types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE
		expectedConfirmationMode := types.ConfirmationMode_FAST

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
			true,
			expectedInboundStatus,
			expectedConfirmationMode,
		)
		require.Equal(t, expectedInboundStatus, msg.Status)
		require.Equal(t, expectedConfirmationMode, msg.ConfirmationMode)
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
				true,
				types.InboundStatus_SUCCESS,
				types.ConfirmationMode_SAFE,
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
				true,
				types.InboundStatus_SUCCESS,
				types.ConfirmationMode_SAFE,
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
				true,
				types.InboundStatus_SUCCESS,
				types.ConfirmationMode_SAFE,
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
				true,
				types.InboundStatus_SUCCESS,
				types.ConfirmationMode_SAFE,
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
				true,
				types.InboundStatus_SUCCESS,
				types.ConfirmationMode_SAFE,
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

	var (
		creator     = sample.AccAddress()
		sender      = sample.AccAddress()
		txOrigin    = sample.String()
		receiver    = sample.String()
		message     = sample.String()
		inboundHash = sample.String()
		asset       = sample.String()
	)

	// getMsg creates a constant message object
	getMsg := func() types.MsgVoteInbound {
		return types.MsgVoteInbound{
			Creator:            creator,
			Sender:             sender,
			SenderChainId:      42,
			TxOrigin:           txOrigin,
			Receiver:           receiver,
			ReceiverChain:      42,
			Amount:             math.NewUint(42),
			Message:            message,
			InboundHash:        inboundHash,
			InboundBlockHeight: 42,
			CallOptions: &types.CallOptions{
				GasLimit: 42,
			},
			CoinType:                coin.CoinType_Zeta,
			Asset:                   asset,
			EventIndex:              42,
			ProtocolContractVersion: types.ProtocolContractVersion_V1,
			Status:                  types.InboundStatus_SUCCESS,
			ConfirmationMode:        types.ConfirmationMode_SAFE,
		}
	}

	// get original digest
	msg := getMsg()
	hash := msg.Digest()
	require.NotEmpty(t, hash, "hash should not be empty")

	// creator not used
	msg = getMsg()
	msg.Creator = sample.AccAddress()
	hash2 := msg.Digest()
	require.Equal(t, hash, hash2, "creator should not change hash")

	// in block height not used
	msg = getMsg()
	msg.InboundBlockHeight = 43
	hash2 = msg.Digest()
	require.Equal(t, hash, hash2, "in block height should not change hash")

	// sender used
	msg = getMsg()
	msg.Sender = sample.AccAddress()
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "sender should change hash")

	// sender chain ID used
	msg = getMsg()
	msg.SenderChainId = 43
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "sender chain ID should change hash")

	// tx origin used
	msg = getMsg()
	msg.TxOrigin = sample.StringRandom(r, 32)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "tx origin should change hash")

	// receiver used
	msg = getMsg()
	msg.Receiver = sample.StringRandom(r, 32)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "receiver should change hash")

	// receiver chain ID used
	msg = getMsg()
	msg.ReceiverChain = 43
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "receiver chain ID should change hash")

	// amount used
	msg = getMsg()
	msg.Amount = math.NewUint(43)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "amount should change hash")

	// message used
	msg = getMsg()
	msg.Message = sample.StringRandom(r, 32)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "message should change hash")

	// in tx hash used
	msg = getMsg()
	msg.InboundHash = sample.StringRandom(r, 32)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "in tx hash should change hash")

	// gas limit used
	msg = getMsg()
	msg.CallOptions.GasLimit = 43
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "gas limit should change hash")

	// coin type used
	msg = getMsg()
	msg.CoinType = coin.CoinType_ERC20
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "coin type should change hash")

	// asset used
	msg = getMsg()
	msg.Asset = sample.StringRandom(r, 32)
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "asset should change hash")

	// event index used
	msg = getMsg()
	msg.EventIndex = 43
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "event index should change hash")

	// protocol contract version used
	msg = getMsg()
	msg.ProtocolContractVersion = types.ProtocolContractVersion_V2
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "protocol contract version should change hash")

	// inbound status used
	msg = getMsg()
	msg.Status = types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE
	hash2 = msg.Digest()
	require.NotEqual(t, hash, hash2, "inbound status should change hash")

	// confirmation mode not used
	msg = getMsg()
	msg.ConfirmationMode = types.ConfirmationMode_FAST
	hash2 = msg.Digest()
	require.Equal(t, hash, hash2, "confirmation mode should not change hash")
}

func TestMsgVoteInbound_EligibleForFastConfirmation(t *testing.T) {
	tests := []struct {
		name     string
		msg      types.MsgVoteInbound
		eligible bool
	}{
		{
			name: "eligible for fast confirmation",
			msg: func() types.MsgVoteInbound {
				msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)
				msg.ProtocolContractVersion = types.ProtocolContractVersion_V2
				return msg
			}(),
			eligible: true,
		},
		{
			name:     "not eligible for non-fungible coin type",
			msg:      sample.InboundVote(coin.CoinType_NoAssetCall, 1, 7000),
			eligible: false,
		},
		{
			name: "not eligible for protocol contract version V1",
			msg: func() types.MsgVoteInbound {
				msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)
				msg.ProtocolContractVersion = types.ProtocolContractVersion_V1
				return msg
			}(),
			eligible: false,
		},
		{
			name: "not eligible for unknown protocol contract version",
			msg: func() types.MsgVoteInbound {
				msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)
				msg.ProtocolContractVersion = types.ProtocolContractVersion(999)
				return msg
			}(),
			eligible: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eligible := tt.msg.EligibleForFastConfirmation()
			require.Equal(t, tt.eligible, eligible)
		})
	}
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
