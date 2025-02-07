package cctx

// Status represents the status of a CCTX transaction, it is more granular than the status present on zetacore
type Status int

const (
	Unknown Status = iota
	// Zetacore statuses
	PendingOutbound Status = 1
	OutboundMined   Status = 2
	PendingRevert   Status = 3
	Reverted        Status = 4
	Aborted         Status = 5
	// Zetatool only statuses
	// PendingInboundConfirmation the inbound transaction is pending confirmation on the inbound chain
	PendingInboundConfirmation Status = 6
	// PendingInboundVoting the inbound transaction is confirmed on the inbound chain, and we are waiting for observers to vote
	PendingInboundVoting Status = 7
	// PendingOutboundSigning the outbound transaction is pending signing by the tss
	PendingOutboundSigning Status = 8
	// PendingRevertSigning the revert transaction is pending signing by the tss
	PendingRevertSigning Status = 9
	// PendingOutboundConfirmation the outbound transaction
	// broadcast by the tss is pending confirmation on the outbound chain
	PendingOutboundConfirmation Status = 10
	// PendingRevertConfirmation the revert transaction broadcast by the tss is pending confirmation on the outbound chain
	PendingRevertConfirmation Status = 11
	// PendingOutboundVoting the outbound transaction is confirmed on the outbound chain,
	//and we are waiting for observers to vote
	PendingOutboundVoting Status = 12
	// PendingRevertVoting the revert transaction is confirmed on the outbound chain,
	//and we are waiting for observers to vote
	PendingRevertVoting Status = 13
)

func (s Status) String() string {
	switch s {
	case PendingInboundConfirmation:
		return "PendingInboundConfirmation"
	case PendingInboundVoting:
		return "PendingInboundVoting"
	case PendingOutbound:
		return "PendingOutbound"
	case OutboundMined:
		return "OutboundMined"
	case PendingRevert:
		return "PendingRevert"
	case Reverted:
		return "Reverted"
	case PendingOutboundConfirmation:
		return "PendingOutboundConfirmation"
	case PendingRevertConfirmation:
		return "PendingRevertConfirmation"
	case PendingRevertVoting:
		return "PendingRevertVoting"
	case Aborted:
		return "Aborted"
	case PendingOutboundSigning:
		return "PendingOutboundSigning"
	case PendingRevertSigning:
		return "PendingRevertSigning"
	case PendingOutboundVoting:
		return "PendingOutboundVoting"
	default:
		return "Unknown"
	}
}
