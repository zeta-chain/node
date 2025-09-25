package repo

import "errors"

// TODO: organize
var (
	ErrHealthCheck = errors.New("unable to check TON client health")

	ErrFetchGasPrice      = errors.New("unable to fetch gas price")
	ErrParseGasPrice      = errors.New("unable to parse gas price")
	ErrGetMasterchainInfo = errors.New("unable to get masterchain info")
	ErrPostVoteGasPrice   = errors.New("unable to post vote for gas price")

	ErrGetTransaction       = errors.New("unable to get transaction")
	ErrGetTransactions      = errors.New("unable to get transactions (by index)")
	ErrGetTransactionsSince = errors.New("unable to get transactions (by last transaction)")
	ErrGetInboundTrackers   = errors.New("unable to get inbound trackers")

	ErrNoTransactions = errors.New("found no transactions")
	ErrEncoding       = errors.New("invalid transaction hash encoding")
)
