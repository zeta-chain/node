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
	OutboundTxSuccessful = "crosschain/OutboundTxSuccessful"
	OutboundTxFailed     = "crosschain/OutboundTxFailed"
	CctxCreated          = "crosschain/CctxCreated"
	InboundFinalized     = "crosschain/InboundFinalized"
	StatusChanged        = "crosschain/StatusChanged"
	CctxFinalized        = "crosschain/CctxFinalized"
	CctxScrubbed         = "crosschain/CCTXScrubbed"
)
