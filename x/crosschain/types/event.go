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
	OutboundTxSuccessful = "crosschain/OutboundTxSuccessful"
	OutboundTxFailed     = "crosschain/OutboundTxFailed"
	ZrcWithdrawCreated   = "crosschain/ZrcWithdrawCreated"
	BallotCreated        = "crosschain/BallotCreated"
	ZetaWithdrawCreated  = "crosschain/ZetaWithdrawCreated"
	InboundFinalized     = "crosschain/InboundFinalized"
	StatusChanged        = "crosschain/StatusChanged"
	CctxScrubbed         = "crosschain/CCTXScrubbed"
	CctxNewKeygenBlock   = "crosschain/CCTXNewKeygenBlock"
)
