package cctx

type Status int

const (
	Unknown                     Status = iota
	PendingInboundConfirmation  Status = 1
	PendingInboundVoting        Status = 2
	PendingOutbound             Status = 3
	OutboundMined               Status = 4
	PendingRevert               Status = 5
	Reverted                    Status = 6
	PendingOutboundConfirmation Status = 7
	PendingRevertConfirmation   Status = 8
	PendingRevertVoting         Status = 10
	Aborted                     Status = 11
	PendingOutboundSigning      Status = 12
	PendingRevertSigning        Status = 13
	PendingOutboundVoting       Status = 14
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
