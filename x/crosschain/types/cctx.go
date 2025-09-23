package types

import (
	"encoding/hex"
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// GetConnectedChainID returns the connected chain ID for the CCTX.
// If the CCTX is outgoing, this is the receiver chain ID.
// If the CCTX is incoming, this is the sender chain ID.
// Second argument is boolean, true if the CCTX is outgoing, false if incoming.
func (m CrossChainTx) GetConnectedChainID() (int64, bool, error) {
	if m.InboundParams == nil {
		return 0, false, fmt.Errorf("inbound params cannot be nil")
	}

	// If the sender chain ID is ZetaChain, this is an outgoing CCTX.
	// Note: additional chains argument is empty, all ZetaChain IDs are hardcoded in the codebase.
	if chains.IsZetaChain(m.InboundParams.SenderChainId, []chains.Chain{}) {
		if len(m.OutboundParams) < 1 || m.OutboundParams[0] == nil {
			return 0, false, fmt.Errorf("outbound params cannot be nil")
		}
		return m.OutboundParams[0].ReceiverChainId, true, nil
	}
	return m.InboundParams.SenderChainId, false, nil
}

// GetEVMRevertAddress returns the EVM revert address
// If a revert address is specified in the revert options, it returns the address
// Otherwise returns sender address
func (m CrossChainTx) GetEVMRevertAddress() ethcommon.Address {
	addr, valid := m.RevertOptions.GetEVMRevertAddress()
	if valid {
		return addr
	}
	return ethcommon.HexToAddress(m.InboundParams.Sender)
}

// GetEVMAbortAddress returns the EVM abort address
// If an abort address is specified in the revert options, it returns the address
// Otherwise returns sender address
func (m CrossChainTx) GetEVMAbortAddress() ethcommon.Address {
	addr, valid := m.RevertOptions.GetEVMAbortAddress()
	if valid {
		return addr
	}
	return ethcommon.HexToAddress(m.InboundParams.Sender)
}

// GetCurrentOutboundParam returns the current outbound params.
// There can only be one active outbound.
// OutboundParams[0] is the original outbound, if it reverts, then
// OutboundParams[1] is the new outbound.
func (m CrossChainTx) GetCurrentOutboundParam() *OutboundParams {
	// TODO: Deprecated (V21) gasLimit should be removed and CallOptions should be mandatory
	// this should never happen, but since it is optional, adding it just in case
	if len(m.OutboundParams) == 0 {
		return &OutboundParams{CallOptions: &CallOptions{}}
	}

	outboundParams := m.OutboundParams[len(m.OutboundParams)-1]
	if outboundParams.CallOptions == nil {
		outboundParams.CallOptions = &CallOptions{
			GasLimit: outboundParams.GasLimit,
		}
	}
	return outboundParams
}

// IsCurrentOutboundRevert returns true if the current outbound is the revert tx.
func (m CrossChainTx) IsCurrentOutboundRevert() bool {
	return len(m.OutboundParams) >= 2
}

// OriginalDestinationChainID returns the original destination of the outbound tx, reverted or not
// If there is no outbound tx, return -1
func (m CrossChainTx) OriginalDestinationChainID() int64 {
	if len(m.OutboundParams) == 0 {
		return -1
	}
	return m.OutboundParams[0].ReceiverChainId
}

// Validate checks if the CCTX is valid.
func (m CrossChainTx) Validate() error {
	if m.InboundParams == nil {
		return fmt.Errorf("inbound tx params cannot be nil")
	}
	if m.OutboundParams == nil {
		return fmt.Errorf("outbound tx params cannot be nil")
	}
	if m.CctxStatus == nil {
		return fmt.Errorf("cctx status cannot be nil")
	}
	if len(m.OutboundParams) > 2 {
		return fmt.Errorf("outbound tx params cannot be more than 2")
	}
	if m.Index != "" {
		err := ValidateCCTXIndex(m.Index)
		if err != nil {
			return err
		}
	}
	err := m.InboundParams.Validate()
	if err != nil {
		return err
	}
	for _, outboundParam := range m.OutboundParams {
		err = outboundParam.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// AddRevertOutbound does the following things in one function:
//  1. create a new OutboundTxParams for the revert
//  2. append the new OutboundTxParams to the current OutboundTxParams
//  3. update the TxFinalizationStatus of the current OutboundTxParams to Executed.
func (m *CrossChainTx) AddRevertOutbound(gasLimit uint64) error {
	if m.IsCurrentOutboundRevert() {
		return fmt.Errorf("cannot revert a revert tx")
	}
	if len(m.OutboundParams) == 0 {
		return fmt.Errorf("cannot revert before trying to process an outbound tx")
	}

	// Use revert address from RevertOptions if available, otherwise use the sender address
	revertReceiver := m.InboundParams.Sender
	if m.ProtocolContractVersion == ProtocolContractVersion_V2 {
		switch {
		case chains.IsBitcoinChain(m.InboundParams.SenderChainId, []chains.Chain{}):
			if m.RevertOptions.RevertAddress != "" {
				revertReceiver = m.RevertOptions.RevertAddress
			}
		case chains.IsSolanaChain(m.InboundParams.SenderChainId, []chains.Chain{}):
			revertAddress, valid := m.RevertOptions.GetSOLRevertAddress()
			if valid {
				revertReceiver = revertAddress.String()
			}
		default:
			revertAddress, valid := m.RevertOptions.GetEVMRevertAddress()
			if valid {
				revertReceiver = revertAddress.Hex()
			}
		}
	}

	revertTxParams := &OutboundParams{
		Receiver:        revertReceiver,
		ReceiverChainId: m.InboundParams.SenderChainId,
		Amount:          m.GetCurrentOutboundParam().Amount,
		CallOptions: &CallOptions{
			GasLimit: gasLimit,
		},
		TssPubkey: m.GetCurrentOutboundParam().TssPubkey,
		// Inherit same confirmation mode from original outbound as placeholder.
		// It will be overwritten by actual confirmation mode in the outbound vote message
		ConfirmationMode: m.GetCurrentOutboundParam().ConfirmationMode,
	}

	// TODO : Refactor to move FungibleTokenCoinType field to the CCTX object directly : https://github.com/zeta-chain/node/issues/1943
	if m.InboundParams != nil {
		revertTxParams.CoinType = m.InboundParams.CoinType
	}
	// The original outbound has been finalized, the new outbound is pending
	m.GetCurrentOutboundParam().TxFinalizationStatus = TxFinalizationStatus_Executed
	m.OutboundParams = append(m.OutboundParams, revertTxParams)
	return nil
}

// AddOutbound adds a new outbound tx to the CCTX.
func (m *CrossChainTx) AddOutbound(
	ctx sdk.Context,
	msg MsgVoteOutbound,
	ballotStatus observertypes.BallotStatus,
) error {
	if ballotStatus != observertypes.BallotStatus_BallotFinalized_FailureObservation {
		if !msg.ValueReceived.Equal(m.GetCurrentOutboundParam().Amount) {
			ctx.Logger().Error(fmt.Sprintf("VoteOutbound: Mint mismatch: %s value received vs %s cctx amount",
				msg.ValueReceived,
				m.GetCurrentOutboundParam().Amount))
			return cosmoserrors.Wrap(
				sdkerrors.ErrInvalidRequest,
				fmt.Sprintf(
					"ValueReceived %s does not match sent value %s",
					msg.ValueReceived,
					m.GetCurrentOutboundParam().Amount,
				),
			)
		}
	}
	// Update CCTX values
	m.GetCurrentOutboundParam().Hash = msg.ObservedOutboundHash
	m.GetCurrentOutboundParam().GasUsed = msg.ObservedOutboundGasUsed
	m.GetCurrentOutboundParam().EffectiveGasPrice = msg.ObservedOutboundEffectiveGasPrice
	m.GetCurrentOutboundParam().EffectiveGasLimit = msg.ObservedOutboundEffectiveGasLimit
	m.GetCurrentOutboundParam().ObservedExternalHeight = msg.ObservedOutboundBlockHeight
	m.GetCurrentOutboundParam().ConfirmationMode = msg.ConfirmationMode
	return nil
}

// SetAbort sets the CCTX status to Aborted with the given error message.
func (m CrossChainTx) SetAbort(messages StatusMessages) {
	m.CctxStatus.UpdateStatusAndErrorMessages(CctxStatus_Aborted, messages)
}

// SetPendingRevert sets the CCTX status to PendingRevert with the given error message.
func (m CrossChainTx) SetPendingRevert(messages StatusMessages) {
	m.CctxStatus.UpdateStatusAndErrorMessages(CctxStatus_PendingRevert, messages)
}

// SetPendingOutbound sets the CCTX status to PendingOutbound with the given error message.
func (m CrossChainTx) SetPendingOutbound(messages StatusMessages) {
	m.CctxStatus.UpdateStatusAndErrorMessages(CctxStatus_PendingOutbound, messages)
}

// SetOutboundMined sets the CCTX status to OutboundMined with the given error message.
func (m CrossChainTx) SetOutboundMined() {
	m.CctxStatus.UpdateStatusAndErrorMessages(CctxStatus_OutboundMined, StatusMessages{})
}

// SetReverted sets the CCTX status to Reverted with the given error message.
func (m CrossChainTx) SetReverted() {
	m.CctxStatus.UpdateStatusAndErrorMessages(CctxStatus_Reverted, StatusMessages{})
}

func (m CrossChainTx) GetCCTXIndexBytes() ([32]byte, error) {
	sendHash := [32]byte{}
	if len(m.Index) < 2 {
		return [32]byte{}, fmt.Errorf("decode CCTX %s index too short", m.Index)
	}
	decodedIndex, err := hex.DecodeString(m.Index[2:]) // remove the leading 0x
	if err != nil || len(decodedIndex) != 32 {
		return [32]byte{}, fmt.Errorf("decode CCTX %s error", m.Index)
	}
	copy(sendHash[:32], decodedIndex[:32])
	return sendHash, nil
}

// SetOutboundBallotIndex sets the outbound ballot index for the most recent outbound.
func (m CrossChainTx) SetOutboundBallotIndex(index string) {
	m.GetCurrentOutboundParam().BallotIndex = index
}

// GetCctxIndexFromBytes returns the CCTX index from a byte array.
func GetCctxIndexFromBytes(sendHash [32]byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(sendHash[:]))
}

// GetCctxIndexFromArbitraryBytes converts an arbitrary byte slice to a CCTX index string.
// Returns an error if the input slice is less than 32 bytes.
func GetCctxIndexFromArbitraryBytes(sendHash []byte) (string, error) {
	if len(sendHash) < 32 {
		return "", fmt.Errorf("input byte slice length %d is less than required 32 bytes", len(sendHash))
	}

	var indexBytes [32]byte
	copy(indexBytes[:], sendHash[:32])
	return GetCctxIndexFromBytes(indexBytes), nil
}

// IsWithdrawAndCall returns true if the CCTX is performing a withdraw and call operation.
func (m CrossChainTx) IsWithdrawAndCall() bool {
	if m.InboundParams == nil || m.CctxStatus == nil {
		return false
	}
	return m.InboundParams.IsCrossChainCall && m.CctxStatus.Status == CctxStatus_PendingOutbound
}

// NewCCTX creates a new CCTX from a MsgVoteInbound message and a TSS pubkey.
// It also validates the created cctx
func NewCCTX(ctx sdk.Context, msg MsgVoteInbound, tssPubkey string) (CrossChainTx, error) {
	index := msg.Digest()

	if msg.TxOrigin == "" {
		msg.TxOrigin = msg.Sender
	}
	inboundParams := &InboundParams{
		Sender:                 msg.Sender,
		SenderChainId:          msg.SenderChainId,
		TxOrigin:               msg.TxOrigin,
		Asset:                  msg.Asset,
		Amount:                 msg.Amount,
		ObservedHash:           msg.InboundHash,
		ObservedExternalHeight: msg.InboundBlockHeight,
		FinalizedZetaHeight:    0,
		BallotIndex:            index,
		CoinType:               msg.CoinType,
		IsCrossChainCall:       msg.IsCrossChainCall,
		Status:                 msg.Status,
		ConfirmationMode:       msg.ConfirmationMode,
		ErrorMessage:           msg.ErrorMessage,
	}

	outboundParams := &OutboundParams{
		Receiver:        msg.Receiver,
		ReceiverChainId: msg.ReceiverChain,
		Hash:            "",
		TssNonce:        0,
		CallOptions: &CallOptions{
			IsArbitraryCall: msg.CallOptions.IsArbitraryCall,
			GasLimit:        msg.CallOptions.GasLimit,
		},
		GasPrice:               "",
		GasPriorityFee:         "",
		BallotIndex:            "",
		ObservedExternalHeight: 0,
		Amount:                 sdkmath.ZeroUint(),
		TssPubkey:              tssPubkey,
		CoinType:               msg.CoinType,
		// use SAFE confirmation mode as default value.
		// it will be overwritten by actual confirmation mode in the outbound vote message.
		ConfirmationMode: ConfirmationMode_SAFE,
	}

	status := &Status{
		Status:              CctxStatus_PendingInbound,
		StatusMessage:       "",
		CreatedTimestamp:    ctx.BlockHeader().Time.Unix(),
		LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
		IsAbortRefunded:     false,
	}
	cctx := CrossChainTx{
		Creator:                 msg.Creator,
		Index:                   index,
		ZetaFees:                sdkmath.ZeroUint(),
		RelayedMessage:          msg.Message,
		CctxStatus:              status,
		InboundParams:           inboundParams,
		OutboundParams:          []*OutboundParams{outboundParams},
		ProtocolContractVersion: msg.ProtocolContractVersion,
		RevertOptions:           msg.RevertOptions,
	}

	// TODO: remove this validate call
	// https://github.com/zeta-chain/node/issues/2236
	err := cctx.Validate()
	if err != nil {
		return CrossChainTx{}, err
	}

	return cctx, nil
}
