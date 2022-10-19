package types

const (
	// event key
	SubTypeKey = "SubTypeKey"
	CctxIndex  = "CctxIndex"

	Sender        = "Sender"
	SenderChain   = "SenderChain"
	InTxHash      = "InTxObservedHash"
	InBlockHeight = "InTxObservedBlockHeight"

	Receiver      = "Receiver"
	ReceiverChain = "ReceiverChain"
	OutTxHash     = "OutTxObservedHash"

	ZetaMint         = "ZetaMint"
	ZetaBurnt        = "ZetaBurnt"
	OutTXVotingChain = "OutTxVotingChain"
	OutBoundChain    = "OutBoundChain"
	OldStatus        = "OldStatus"
	NewStatus        = "NewStatus"
	StatusMessage    = "StatusMessage"
	RelayedMessage   = "RelayedMessage"
	Identifiers      = "LogIdentifiers"
)

const (
	OutboundTxSuccessful = "zetacore/OutboundTxSuccessful"
	OutboundTxFailed     = "zetacore/OutboundTxFailed"
	CctxCreated          = "zetacore/CctxCreated"
	InboundFinalized     = "zetacore/InboundFinalized"
	StatusChanged        = "zetacore/StatusChanged"
	CctxFinalized        = "zetacore/CctxFinalized"
	SendScrubbed         = "SendScrubbed"
)
