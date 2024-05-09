package zetabridge

const (
	// DefaultGasLimit is the default gas limit used for broadcasting txs
	DefaultGasLimit = 200_000

	// PostGasPriceGasLimit is the gas limit for voting new gas price
	PostGasPriceGasLimit = 1_500_000

	// AddTxHashToOutboundTrackerGasLimit is the gas limit for adding tx hash to out tx tracker
	AddTxHashToOutboundTrackerGasLimit = 200_000

	// PostBlameDataGasLimit is the gas limit for voting on blames
	PostBlameDataGasLimit = 200_000

	// DefaultRetryCount is the number of retries for broadcasting a tx
	DefaultRetryCount = 5

	// ExtendedRetryCount is an extended number of retries for broadcasting a tx, used in keygen operations
	ExtendedRetryCount = 15

	// DefaultRetryInterval is the interval between retries in seconds
	DefaultRetryInterval = 5

	// MonitorVoteInboundTxResultInterval is the interval between retries for monitoring tx result in seconds
	MonitorVoteInboundTxResultInterval = 5

	// MonitorVoteInboundTxResultRetryCount is the number of retries to fetch monitoring tx result
	MonitorVoteInboundTxResultRetryCount = 20

	// PostVoteOutboundGasLimit is the gas limit for voting on observed outbound tx
	PostVoteOutboundGasLimit = 400_000

	// PostVoteOutboundRevertGasLimit is the gas limit for voting on observed outbound tx for revert (when outbound fails)
	// The value needs to be higher because reverting implies interacting with the EVM to perform swaps for the gas token
	PostVoteOutboundRevertGasLimit = 1_500_000

	// MonitorVoteOutboundTxResultInterval is the interval between retries for monitoring tx result in seconds
	MonitorVoteOutboundTxResultInterval = 5

	// MonitorVoteOutboundTxResultRetryCount is the number of retries to fetch monitoring tx result
	MonitorVoteOutboundTxResultRetryCount = 20
)
