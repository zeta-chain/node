package types

const (
	// event key
	SubTypeKey    = "SubTypeKey"
	CctxIndex     = "CctxIndex"
	KeyGenBlock   = "KeyGenBlock"
	KeyGenPubKeys = "KeyGenPubKeys"

	Sender        = "Sender"
	SenderChain   = "SenderChain"
	TxOrigin      = "TxOrigin"
	InTxHash      = "InTxObservedHash"
	InBlockHeight = "InTxObservedBlockHeight"

	Receiver      = "Receiver"
	ReceiverChain = "ReceiverChain"
	OutTxHash     = "OutTxObservedHash"

	ZetaMint         = "ZetaMint"
	Amount           = "Amount"
	Asset            = "Asset"
	OutTXVotingChain = "OutTxVotingChain"
	OutBoundChain    = "OutBoundChain"
	OldStatus        = "OldStatus"
	NewStatus        = "NewStatus"
	StatusMessage    = "StatusMessage"
	RelayedMessage   = "RelayedMessage"
	Identifiers      = "LogIdentifiers"

	BallotIdentifier       = "BallotIdentifier"
	CCTXIndex              = "CCTXIndex"
	BallotObservationHash  = "BallotObservationHash"
	BallotObservationChain = "BallotObservationChain"
	BallotType             = "BallotType"
)

const (
	OutboundTxSuccessful = "OutboundTxSuccessful"
	OutboundTxFailed     = "OutboundTxFailed"
	ZrcWithdrawCreated   = "ZrcWithdrawCreated"
	BallotCreated        = "BallotCreated"
	ZetaWithdrawCreated  = "ZetaWithdrawCreated"
	InboundFinalized     = "InboundFinalized"
	StatusChanged        = "StatusChanged"
	CctxScrubbed         = "CCTXScrubbed"
	CctxNewKeygenBlock   = "CCTXNewKeygenBlock"
)
