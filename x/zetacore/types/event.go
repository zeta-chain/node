package types

const (
	// event value
	OutboundTxSuccessful string = "OutboundTxSuccessful"
	OutboundTxFailed            = "OutboundTxFailed"
	InboundCreated              = "InboundCreated"
	InboundFinalized            = "InboundFinalized"

	// event key
	SubType       = "SubType"
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
