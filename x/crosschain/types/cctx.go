package types

import (
	"encoding/hex"
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// GetCurrentOutTxParam returns the current outbound tx params.
// There can only be one active outtx.
// OutboundTxParams[0] is the original outtx, if it reverts, then
// OutboundTxParams[1] is the new outtx.
func (m CrossChainTx) GetCurrentOutTxParam() *OutboundTxParams {
	if len(m.OutboundTxParams) == 0 {
		return &OutboundTxParams{}
	}
	return m.OutboundTxParams[len(m.OutboundTxParams)-1]
}

// IsCurrentOutTxRevert returns true if the current outbound tx is the revert tx.
func (m CrossChainTx) IsCurrentOutTxRevert() bool {
	return len(m.OutboundTxParams) >= 2
}

// OriginalDestinationChainID returns the original destination of the outbound tx, reverted or not
// If there is no outbound tx, return -1
func (m CrossChainTx) OriginalDestinationChainID() int64 {
	if len(m.OutboundTxParams) == 0 {
		return -1
	}
	return m.OutboundTxParams[0].ReceiverChainId
}

// Validate checks if the CCTX is valid.
func (m CrossChainTx) Validate() error {
	if m.InboundTxParams == nil {
		return fmt.Errorf("inbound tx params cannot be nil")
	}
	if m.OutboundTxParams == nil {
		return fmt.Errorf("outbound tx params cannot be nil")
	}
	if m.CctxStatus == nil {
		return fmt.Errorf("cctx status cannot be nil")
	}
	if len(m.OutboundTxParams) > 2 {
		return fmt.Errorf("outbound tx params cannot be more than 2")
	}
	if m.Index != "" {
		err := ValidateZetaIndex(m.Index)
		if err != nil {
			return err
		}
	}
	err := m.InboundTxParams.Validate()
	if err != nil {
		return err
	}
	for _, outboundTxParam := range m.OutboundTxParams {
		err = outboundTxParam.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

/*
AddRevertOutbound does the following things in one function:

	1. create a new OutboundTxParams for the revert

	2. append the new OutboundTxParams to the current OutboundTxParams

	3. update the TxFinalizationStatus of the current OutboundTxParams to Executed.
*/

func (m *CrossChainTx) AddRevertOutbound(gasLimit uint64) error {
	if m.IsCurrentOutTxRevert() {
		return fmt.Errorf("cannot revert a revert tx")
	}
	revertTxParams := &OutboundTxParams{
		Receiver:           m.InboundTxParams.Sender,
		ReceiverChainId:    m.InboundTxParams.SenderChainId,
		Amount:             m.InboundTxParams.Amount,
		OutboundTxGasLimit: gasLimit,
		TssPubkey:          m.GetCurrentOutTxParam().TssPubkey,
	}
	// The original outbound has been finalized, the new outbound is pending
	m.GetCurrentOutTxParam().TxFinalizationStatus = TxFinalizationStatus_Executed
	m.OutboundTxParams = append(m.OutboundTxParams, revertTxParams)
	return nil
}

// AddOutbound adds a new outbound tx to the CCTX.
func (m *CrossChainTx) AddOutbound(ctx sdk.Context, msg MsgVoteOnObservedOutboundTx, ballotStatus observertypes.BallotStatus) error {
	if ballotStatus != observertypes.BallotStatus_BallotFinalized_FailureObservation {
		if !msg.ValueReceived.Equal(m.GetCurrentOutTxParam().Amount) {
			ctx.Logger().Error(fmt.Sprintf("VoteOnObservedOutboundTx: Mint mismatch: %s value received vs %s cctx amount",
				msg.ValueReceived,
				m.GetCurrentOutTxParam().Amount))
			return cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ValueReceived %s does not match sent value %s", msg.ValueReceived, m.GetCurrentOutTxParam().Amount))
		}
	}
	// Update CCTX values
	m.GetCurrentOutTxParam().OutboundTxHash = msg.ObservedOutTxHash
	m.GetCurrentOutTxParam().OutboundTxGasUsed = msg.ObservedOutTxGasUsed
	m.GetCurrentOutTxParam().OutboundTxEffectiveGasPrice = msg.ObservedOutTxEffectiveGasPrice
	m.GetCurrentOutTxParam().OutboundTxEffectiveGasLimit = msg.ObservedOutTxEffectiveGasLimit
	m.GetCurrentOutTxParam().OutboundTxObservedExternalHeight = msg.ObservedOutTxBlockHeight
	m.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()
	return nil
}

// SetAbort sets the CCTX status to Aborted with the given error message.
func (m CrossChainTx) SetAbort(message string) {
	m.CctxStatus.ChangeStatus(CctxStatus_Aborted, message)
}

// SetPendingRevert sets the CCTX status to PendingRevert with the given error message.
func (m CrossChainTx) SetPendingRevert(message string) {
	m.CctxStatus.ChangeStatus(CctxStatus_PendingRevert, message)
}

// SetPendingOutbound sets the CCTX status to PendingOutbound with the given error message.
func (m CrossChainTx) SetPendingOutbound(message string) {
	m.CctxStatus.ChangeStatus(CctxStatus_PendingOutbound, message)
}

// SetOutBoundMined sets the CCTX status to OutboundMined with the given error message.
func (m CrossChainTx) SetOutBoundMined(message string) {
	m.CctxStatus.ChangeStatus(CctxStatus_OutboundMined, message)
}

// SetReverted sets the CCTX status to Reverted with the given error message.
func (m CrossChainTx) SetReverted(message string) {
	m.CctxStatus.ChangeStatus(CctxStatus_Reverted, message)
}

func (m CrossChainTx) GetCCTXIndexBytes() ([32]byte, error) {
	sendHash := [32]byte{}
	decodedIndex, err := hex.DecodeString(m.Index[2:]) // remove the leading 0x
	if err != nil || len(decodedIndex) != 32 {
		return [32]byte{}, fmt.Errorf("decode CCTX %s error", m.Index)
	}
	copy(sendHash[:32], decodedIndex[:32])
	return sendHash, nil
}

func GetCctxIndexFromBytes(sendHash [32]byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(sendHash[:]))
}

// NewCCTX creates a new CCTX.From a MsgVoteOnObservedInboundTx message and a TSS pubkey.
// It also validates the created cctx
func NewCCTX(ctx sdk.Context, msg MsgVoteOnObservedInboundTx, tssPubkey string) (CrossChainTx, error) {
	index := msg.Digest()

	if msg.TxOrigin == "" {
		msg.TxOrigin = msg.Sender
	}
	inboundParams := &InboundTxParams{
		Sender:                          msg.Sender,
		SenderChainId:                   msg.SenderChainId,
		TxOrigin:                        msg.TxOrigin,
		Asset:                           msg.Asset,
		Amount:                          msg.Amount,
		InboundTxObservedHash:           msg.InTxHash,
		InboundTxObservedExternalHeight: msg.InBlockHeight,
		InboundTxFinalizedZetaHeight:    0,
		InboundTxBallotIndex:            index,
		CoinType:                        msg.CoinType,
	}

	outBoundParams := &OutboundTxParams{
		Receiver:                         msg.Receiver,
		ReceiverChainId:                  msg.ReceiverChain,
		OutboundTxHash:                   "",
		OutboundTxTssNonce:               0,
		OutboundTxGasLimit:               msg.GasLimit,
		OutboundTxGasPrice:               "",
		OutboundTxBallotIndex:            "",
		OutboundTxObservedExternalHeight: 0,
		Amount:                           sdkmath.ZeroUint(),
		TssPubkey:                        tssPubkey,
		CoinType:                         msg.CoinType,
	}
	status := &Status{
		Status:              CctxStatus_PendingInbound,
		StatusMessage:       "",
		LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
		IsAbortRefunded:     false,
	}
	cctx := CrossChainTx{
		Creator:          msg.Creator,
		Index:            index,
		ZetaFees:         sdkmath.ZeroUint(),
		RelayedMessage:   msg.Message,
		CctxStatus:       status,
		InboundTxParams:  inboundParams,
		OutboundTxParams: []*OutboundTxParams{outBoundParams},
	}
	err := cctx.Validate()
	if err != nil {
		return CrossChainTx{}, err
	}
	return cctx, nil
}
