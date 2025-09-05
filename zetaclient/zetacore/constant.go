package zetacore

const (
	// DefaultBaseGasPrice is the default base gas price
	DefaultBaseGasPrice = 1_000_000

	// DefaultGasLimit is the default gas limit used for broadcasting txs
	DefaultGasLimit = 200_000

	// PostGasPriceGasLimit is the gas limit for voting new gas price
	PostGasPriceGasLimit = 1_500_000

	// PostVoteInboundGasLimit is the gas limit for voting on observed inbound tx (for zetachain itself)
	PostVoteInboundGasLimit = 500_000

	// PostTSSGasLimit is the gas limit for voting on TSS keygen
	PostTSSGasLimit = 500_000

	// PostVoteInboundExecutionGasLimit is the gas limit for voting on observed inbound tx and executing it
	PostVoteInboundExecutionGasLimit = 7_000_000

	// PostVoteInboundMessagePassingExecutionGasLimit is the gas limit for voting on, and executing ,observed inbound tx related to message passing (coin_type == zeta)
	PostVoteInboundMessagePassingExecutionGasLimit = 4_000_000

	// PostVoteInboundCallOptionsGasLimit is the gas limit for inbound call options
	PostVoteInboundCallOptionsGasLimit uint64 = 1_500_000

	// AddOutboundTrackerGasLimit is the gas limit for adding tx hash to out tx tracker
	AddOutboundTrackerGasLimit = 400_000

	// PostBlameDataGasLimit is the gas limit for voting on blames
	PostBlameDataGasLimit = 200_000

	// PostVoteOutboundGasLimit is the gas limit for voting on observed outbound tx (for zetachain itself)
	PostVoteOutboundGasLimit = 500_000

	// PostVoteOutboundRevertGasLimit is the gas limit for voting on observed outbound tx for revert (when outbound fails)
	// The value is set to 7M because in case of onRevert call, it might consume lot of gas
	PostVoteOutboundRevertGasLimit = 7_000_000

	// PostVoteOutboundRevertGasLimit is the retry gas limit for voting on observed outbound tx for success outbound
	PostVoteOutboundRetryGasLimit uint64 = 1_000_000
)
