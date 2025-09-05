package types

import (
	"encoding/hex"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
)

// ParseGatewayEvent parses the event from the gateway contract
func ParseGatewayEvent(
	log ethtypes.Log,
	gatewayAddr ethcommon.Address,
) (*gatewayzevm.GatewayZEVMWithdrawn, *gatewayzevm.GatewayZEVMCalled, *gatewayzevm.GatewayZEVMWithdrawnAndCalled, error) {
	if len(log.Topics) == 0 {
		return nil, nil, nil, errors.New("ParseGatewayCallEvent: invalid log - no topics")
	}
	filterer, err := gatewayzevm.NewGatewayZEVMFilterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, nil, nil, err
	}
	withdrawalEvent, err := ParseGatewayWithdrawalEvent(log, gatewayAddr, filterer)
	if err == nil {
		return withdrawalEvent, nil, nil, nil
	}
	callEvent, err := ParseGatewayCallEvent(log, gatewayAddr, filterer)
	if err == nil {
		return nil, callEvent, nil, nil
	}
	withdrawAndCallEvent, err := ParseGatewayWithdrawAndCallEvent(log, gatewayAddr, filterer)
	if err == nil {
		return nil, nil, withdrawAndCallEvent, nil
	}
	return nil, nil, nil, errors.New("ParseGatewayEvent: invalid log - no event found")
}

// ParseGatewayWithdrawalEvent parses the GatewayZEVMWithdrawal event from the log
func ParseGatewayWithdrawalEvent(
	log ethtypes.Log,
	gatewayAddr ethcommon.Address,
	filterer *gatewayzevm.GatewayZEVMFilterer,
) (*gatewayzevm.GatewayZEVMWithdrawn, error) {
	event, err := filterer.ParseWithdrawn(log)
	if err != nil {
		return nil, err
	}
	if event.Raw.Address != gatewayAddr {
		return nil, errors.New("ParseGatewayWithdrawalEvent: invalid log - wrong contract address")
	}

	return event, nil
}

// ParseGatewayCallEvent parses the GatewayZEVMCall event from the log
func ParseGatewayCallEvent(
	log ethtypes.Log,
	gatewayAddr ethcommon.Address,
	filterer *gatewayzevm.GatewayZEVMFilterer,
) (*gatewayzevm.GatewayZEVMCalled, error) {
	event, err := filterer.ParseCalled(log)
	if err != nil {
		return nil, err
	}
	if event.Raw.Address != gatewayAddr {
		return nil, errors.New("ParseGatewayCallEvent: invalid log - wrong contract address")
	}
	return event, nil
}

// ParseGatewayWithdrawAndCallEvent parses the GatewayZEVMWithdrawAndCall event from the log
func ParseGatewayWithdrawAndCallEvent(
	log ethtypes.Log,
	gatewayAddr ethcommon.Address,
	filterer *gatewayzevm.GatewayZEVMFilterer,
) (*gatewayzevm.GatewayZEVMWithdrawnAndCalled, error) {
	event, err := filterer.ParseWithdrawnAndCalled(log)
	if err != nil {
		return nil, err
	}
	if event.Raw.Address != gatewayAddr {
		return nil, errors.New("ParseGatewayWithdrawAndCallEvent: invalid log - wrong contract address")
	}
	return event, nil
}

// NewWithdrawalInbound creates a new inbound object for a withdrawal
// currently inbound data is represented with a MsgVoteInbound message
// TODO: replace with a more appropriate object
// https://github.com/zeta-chain/node/issues/2658
func NewWithdrawalInbound(
	ctx sdk.Context,
	txOrigin string,
	coinType coin.CoinType,
	asset string,
	event *gatewayzevm.GatewayZEVMWithdrawn,
	receiverChain chains.Chain,
	gasLimitQueried *big.Int,
) (*MsgVoteInbound, error) {
	senderChain, err := chains.ZetaChainFromCosmosChainID(ctx.ChainID())
	if err != nil {
		return nil, errors.Wrapf(err, "ProcessZEVMInboundV2: failed to convert chainID %s", ctx.ChainID())
	}

	toAddr, err := receiverChain.EncodeAddress(event.Receiver)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot encode address %v", event.Receiver)
	}

	gasLimit := event.CallOptions.GasLimit.Uint64()
	if gasLimit == 0 {
		gasLimit = gasLimitQueried.Uint64()
	}

	// if the message is not empty, specify cross-chain call for backward compatibility with the Withdraw event
	isCrossChainCall := len(event.Message) > 0

	return NewMsgVoteInbound(
		"",
		event.Sender.Hex(),
		senderChain.ChainId,
		txOrigin,
		toAddr,
		receiverChain.ChainId,
		math.NewUintFromBigInt(event.Value),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.String(),
		event.Raw.BlockNumber,
		gasLimit,
		coinType,
		asset,
		uint64(event.Raw.Index),
		ProtocolContractVersion_V2,
		event.CallOptions.IsArbitraryCall,
		InboundStatus_SUCCESS,
		ConfirmationMode_SAFE,
		WithZEVMRevertOptions(event.RevertOptions),
		WithCrossChainCall(isCrossChainCall),
	), nil
}

// NewCallInbound creates a new inbound object for a call
// currently inbound data is represented with a MsgVoteInbound message
// TODO: replace with a more appropriate object
// https://github.com/zeta-chain/node/issues/2658
func NewCallInbound(
	ctx sdk.Context,
	txOrigin string,
	event *gatewayzevm.GatewayZEVMCalled,
	receiverChain chains.Chain,
	gasLimitQueried *big.Int,
) (*MsgVoteInbound, error) {
	senderChain, err := chains.ZetaChainFromCosmosChainID(ctx.ChainID())
	if err != nil {
		return nil, errors.Wrapf(err, "ProcessZEVMInboundV2: failed to convert chainID %s", ctx.ChainID())
	}

	toAddr, err := receiverChain.EncodeAddress(event.Receiver)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot encode address %v", event.Receiver)
	}

	gasLimit := event.CallOptions.GasLimit.Uint64()
	if gasLimit == 0 {
		gasLimit = gasLimitQueried.Uint64()
	}

	return NewMsgVoteInbound(
		"",
		event.Sender.Hex(),
		senderChain.ChainId,
		txOrigin,
		toAddr,
		receiverChain.ChainId,
		math.ZeroUint(),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.String(),
		event.Raw.BlockNumber,
		gasLimit,
		coin.CoinType_NoAssetCall,
		"",
		uint64(event.Raw.Index),
		ProtocolContractVersion_V2,
		event.CallOptions.IsArbitraryCall,
		InboundStatus_SUCCESS,
		ConfirmationMode_SAFE,
		WithZEVMRevertOptions(event.RevertOptions),
	), nil
}

// NewWithdrawAndCallInbound creates a new inbound object for a withdraw and call
// currently inbound data is represented with a MsgVoteInbound message
// TODO: replace with a more appropriate object
// https://github.com/zeta-chain/node/issues/2658
func NewWithdrawAndCallInbound(
	ctx sdk.Context,
	txOrigin string,
	coinType coin.CoinType,
	asset string,
	event *gatewayzevm.GatewayZEVMWithdrawnAndCalled,
	receiverChain chains.Chain,
	gasLimitQueried *big.Int,
) (*MsgVoteInbound, error) {
	senderChain, err := chains.ZetaChainFromCosmosChainID(ctx.ChainID())
	if err != nil {
		return nil, errors.Wrapf(err, "ProcessZEVMInboundV2: failed to convert chainID %s", ctx.ChainID())
	}

	toAddr, err := receiverChain.EncodeAddress(event.Receiver)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot encode address %v", event.Receiver)
	}

	// Temporary NOTE for PR review: Don't think this check is necessary as CallOptions.GasLimit is always set
	// If its zero we should probably return as the event is directly from the smart contract
	gasLimit := event.CallOptions.GasLimit.Uint64()
	if gasLimit == 0 {
		gasLimit = gasLimitQueried.Uint64()
	}

	return NewMsgVoteInbound(
		"",
		event.Sender.Hex(),
		senderChain.ChainId,
		txOrigin,
		toAddr,
		receiverChain.ChainId,
		math.NewUintFromBigInt(event.Value),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.String(),
		event.Raw.BlockNumber,
		gasLimit,
		coinType,
		asset,
		uint64(event.Raw.Index),
		ProtocolContractVersion_V2,
		event.CallOptions.IsArbitraryCall,
		InboundStatus_SUCCESS,
		ConfirmationMode_SAFE,
		WithZEVMRevertOptions(event.RevertOptions),
		WithCrossChainCall(true),
	), nil
}
