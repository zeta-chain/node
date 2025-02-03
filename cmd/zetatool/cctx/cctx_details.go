package cctx

import (
	"fmt"

	"github.com/zeta-chain/node/cmd/zetatool/context"
	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TrackingDetails represents the status of a CCTX transaction
type TrackingDetails struct {
	CCTXIdentifier          string       `json:"cctx_identifier"`
	Status                  Status       `json:"status"`
	OutboundChain           chains.Chain `json:"outbound_chain_id"`
	OutboundTssNonce        uint64       `json:"outbound_tss_nonce"`
	OutboundTrackerHashList []string     `json:"outbound_tracker_hash_list"`
	Message                 string       `json:"message"`
}

func NewCCTXDetails() *TrackingDetails {
	return &TrackingDetails{
		CCTXIdentifier: "",
		Status:         Unknown,
	}
}

func (c *TrackingDetails) UpdateStatusFromZetacoreCCTX(status crosschaintypes.CctxStatus) {
	switch status {
	case crosschaintypes.CctxStatus_PendingOutbound:
		c.Status = PendingOutbound
	case crosschaintypes.CctxStatus_OutboundMined:
		c.Status = OutboundMined
	case crosschaintypes.CctxStatus_Reverted:
		c.Status = Reverted
	case crosschaintypes.CctxStatus_PendingRevert:
		c.Status = PendingRevert
	case crosschaintypes.CctxStatus_Aborted:
		c.Status = Aborted
	default:
		c.Status = Unknown
	}
}

func (c *TrackingDetails) IsPendingOutbound() bool {
	return c.Status == PendingOutbound || c.Status == PendingRevert
}

func (c *TrackingDetails) IsPendingConfirmation() bool {
	return c.Status == PendingOutboundConfirmation || c.Status == PendingRevertConfirmation
}

func (c *TrackingDetails) Print() string {
	return fmt.Sprintf("CCTX: %s Status: %s", c.CCTXIdentifier, c.Status.String())
}

func (c *TrackingDetails) DebugPrint() string {
	return fmt.Sprintf("CCTX: %s Status: %s Message: %s", c.CCTXIdentifier, c.Status.String(), c.Message)
}

// UpdateCCTXStatus updates the TrackingDetails with status from zetacore
func (c *TrackingDetails) UpdateCCTXStatus(ctx *context.Context) {
	var (
		zetacoreClient = ctx.GetZetaCoreClient()
		goCtx          = ctx.GetContext()
	)

	CCTX, err := zetacoreClient.GetCctxByHash(goCtx, c.CCTXIdentifier)
	if err != nil {
		c.Message = fmt.Sprintf("failed to get cctx: %v", err)
		return
	}

	c.UpdateStatusFromZetacoreCCTX(CCTX.CctxStatus.Status)

	return
}

func (c *TrackingDetails) UpdateCCTXOutboundDetails(ctx *context.Context) {
	var (
		zetacoreClient = ctx.GetZetaCoreClient()
		goCtx          = ctx.GetContext()
	)
	CCTX, err := zetacoreClient.GetCctxByHash(goCtx, c.CCTXIdentifier)
	if err != nil {
		c.Message = fmt.Sprintf("failed to get cctx: %v", err)
	}
	chainID := CCTX.GetCurrentOutboundParam().ReceiverChainId

	// This is almost impossible to happen as the cctx would not have been created if the chain was not supported
	chain, found := chains.GetChainFromChainID(chainID, []chains.Chain{})
	if !found {
		c.Message = fmt.Sprintf("receiver chain not supported,chain id: %d", chainID)
	}
	c.OutboundChain = chain
	c.OutboundTssNonce = CCTX.GetCurrentOutboundParam().TssNonce
	return
}

func (c *TrackingDetails) UpdateHashListAndPendingStatus(ctx *context.Context) {
	var (
		zetacoreClient = ctx.GetZetaCoreClient()
		goCtx          = ctx.GetContext()
		outboundChain  = c.OutboundChain
		outboundNonce  = c.OutboundTssNonce
	)

	if !c.IsPendingOutbound() {
		return
	}

	tracker, err := zetacoreClient.GetOutboundTracker(goCtx, outboundChain, outboundNonce)
	// tracker is found that means the outbound has been broadcast, but we are waiting for confirmations
	if err == nil && tracker != nil {
		c.updateOutboundConfirmation()
		var hashList []string
		for _, hash := range tracker.HashList {
			hashList = append(hashList, hash.TxHash)
		}
		c.OutboundTrackerHashList = hashList
		return
	}
	// the cctx is in pending state by the outbound signing has not been done
	c.updateOutboundSigning()
	return
}

// 0 - Inbound Confirmation
func (c *TrackingDetails) updateInboundConfirmation(isConfirmed bool) {
	c.Status = PendingInboundConfirmation
	if isConfirmed {
		c.Status = PendingInboundVoting
	}
}

// 1 - Outbound Signing
func (c *TrackingDetails) updateOutboundSigning() {
	switch {
	case c.Status == PendingOutbound:
		c.Status = PendingOutboundSigning
	case c.Status == PendingRevert:
		c.Status = PendingRevertSigning
	}
}

// 2 - Outbound Confirmation
func (c *TrackingDetails) updateOutboundConfirmation() {
	switch {
	case c.Status == PendingOutbound:
		c.Status = PendingOutboundConfirmation
	case c.Status == PendingRevert:
		c.Status = PendingRevertConfirmation
	}
}

// 3 - Outbound Voting
func (c *TrackingDetails) updateOutboundVoting() {
	switch {
	case c.Status == PendingOutboundConfirmation:
		c.Status = PendingOutboundVoting
	case c.Status == PendingRevertConfirmation:
		c.Status = PendingRevertVoting
	}
}
