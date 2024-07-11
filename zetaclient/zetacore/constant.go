package zetacore

import "time"

const (
	// DefaultBaseGasPrice is the default base gas price
	DefaultBaseGasPrice = 1_000_000

	// DefaultGasLimit is the default gas limit used for broadcasting txs
	DefaultGasLimit = 200_000

	// PostGasPriceGasLimit is the gas limit for voting new gas price
	PostGasPriceGasLimit = 1_500_000

	// PostVoteInboundGasLimit is the gas limit for voting on observed inbound tx (for zetachain itself)
	PostVoteInboundGasLimit = 400_000

	// PostVoteInboundExecutionGasLimit is the gas limit for voting on observed inbound tx and executing it
	PostVoteInboundExecutionGasLimit = 4_000_000

	// PostVoteInboundMessagePassingExecutionGasLimit is the gas limit for voting on, and executing ,observed inbound tx related to message passing (coin_type == zeta)
	PostVoteInboundMessagePassingExecutionGasLimit = 4_000_000

	// AddOutboundTrackerGasLimit is the gas limit for adding tx hash to out tx tracker
	AddOutboundTrackerGasLimit = 200_000

	// PostBlameDataGasLimit is the gas limit for voting on blames
	PostBlameDataGasLimit = 200_000

	// DefaultRetryCount is the number of retries for broadcasting a tx
	DefaultRetryCount = 5

	// ExtendedRetryCount is an extended number of retries for broadcasting a tx, used in keygen operations
	ExtendedRetryCount = 15

	// DefaultRetryInterval is the interval between retries in seconds
	DefaultRetryInterval = 5

	// PostVoteOutboundGasLimit is the gas limit for voting on observed outbound tx (for zetachain itself)
	PostVoteOutboundGasLimit = 400_000

	// PostVoteOutboundRevertGasLimit is the gas limit for voting on observed outbound tx for revert (when outbound fails)
	// The value needs to be higher because reverting implies interacting with the EVM to perform swaps for the gas token
	PostVoteOutboundRevertGasLimit = 1_500_000
)

// constants for monitoring tx results
const (
	monitorInterval   = 5 * time.Second
	monitorRetryCount = 20
)
