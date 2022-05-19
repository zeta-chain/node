package types

const (
	// event key
	SubTypeKey    = "SubTypeKey"
	SendHash      = "SendHash"
	OutTxHash     = "OutTxHash"
	ZetaMint      = "ZetaMint"
	ZetaBurnt     = "ZetaBurnt"
	Chain         = "Chain"
	OldStatus     = "OldStatus"
	NewStatus     = "NewStatus"
	Sender        = "Sender"
	SenderChain   = "SenderChain"
	Receiver      = "Receiver"
	ReceiverChain = "ReceiverChain"
	Message       = "Message"
	InTxHash      = "InTxHash"
	InBlockHeight = "InBlockHeight"
)

type SubType string

const (
	OutboundTxSuccessful SubType = "OutboundTxSuccessful"
	OutboundTxFailed             = "OutboundTxFailed"
	InboundCreated               = "InboundCreated"
	InboundFinalized             = "InboundFinalized"
)
