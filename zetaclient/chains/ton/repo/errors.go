package repo

import "errors"

var (
	ErrFetchGasPrice      = errors.New("failed to fetch gas price")
	ErrParseGasPrice      = errors.New("failed to parse gas price")
	ErrGetMasterchainInfo = errors.New("failed to get masterchain info")

	ErrPostVoteGasPrice = errors.New("failed to post vote for gas price")

	ErrGetTransactions = errors.New("failed to get transactions")
	ErrNoTransactions  = errors.New("found no transactions")

	ErrTransactionEncoding  = errors.New("invalid encoding for transaction")
	ErrGetTransactionsSince = errors.New("failed to get transactions")

	ErrGetInboundTrackers = errors.New("failed to get inbound trackers")
)
