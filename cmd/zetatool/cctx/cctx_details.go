package cctx

import crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"

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
	Aborted                     Status = 9
)

type CCTXDetails struct {
	CCCTXIdentifier         string   `json:"cctx_identifier"`
	Status                  Status   `json:"status"`
	OutboundChainID         int64    `json:"outbound_chain_id"`
	OutboundTrackerHashList []string `json:"outbound_tracker_hash_list"`
}

func NewCCTXDetails() CCTXDetails {
	return CCTXDetails{
		CCCTXIdentifier: "",
		Status:          Unknown,
	}
}

func (c *CCTXDetails) UpdateStatusFromZetacoreCCTX(status crosschaintypes.CctxStatus) {
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

func (c *CCTXDetails) IsPendingConfirmation() bool {
	return c.Status == PendingOutbound || c.Status == PendingRevert
}
