package repo

import "errors"

// TODO: organize
var (
	ErrHealthCheck = errors.New("failed to check TON client health")

	ErrFetchGasPrice      = errors.New("failed to fetch gas price")
	ErrParseGasPrice      = errors.New("failed to parse gas price")
	ErrGetMasterchainInfo = errors.New("failed to get masterchain info")
	ErrPostVoteGasPrice   = errors.New("failed to post vote for gas price")

	ErrGetTransaction       = errors.New("failed to get transaction")
	ErrGetTransactions      = errors.New("failed to get transactions (by index)")
	ErrGetTransactionsSince = errors.New("failed to get transactions (by last transaction)")
	ErrGetInboundTrackers   = errors.New("failed to get inbound trackers")

	ErrNoTransactions = errors.New("found no transactions")
	ErrEncoding       = errors.New("invalid transaction hash encoding")
)
