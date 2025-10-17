package zrepo

import (
	"errors"
	"fmt"
)

var (
	ErrNotRPCError = errors.New("not a RPC error")

	ErrClientGetCCTXByHash = errors.New("failed to get CCTX by hash")
	ErrClientGetBallotByID = errors.New("failed to get ballot by ID")
)

var (
	ErrClient = errors.New("error calling a zetacore client function")

	ErrClientGetBTCTSSAddress        = errors.New("failed to get BTC TSS address")
	ErrClientGetCCTX                 = errors.New("failed to get CCTX for nonce")
	ErrClientGetInboundTrackers      = errors.New("failed to get inbound trackers")
	ErrClientGetOutboundTrackers     = errors.New("failed to get outbound trackers")
	ErrClientGetPendingCCTXs         = errors.New("failed to get pending CCTXs")
	ErrClientGetPendingNonces        = errors.New("failed to get pending nonces")
	ErrClientGetForeignCoinsForAsset = errors.New("failed to get foreign coins for asset")

	ErrClientPostOutboundTracker = errors.New("failed to post outbound tracker")
	ErrClientVoteGasPrice        = errors.New("failed to post gas price vote")
	ErrClientVoteInbound         = errors.New("failed to post inbound vote")
	ErrClientVoteOutbound        = errors.New("failed to post outbound vote")

	ErrClientNewBlockSubscriber = errors.New("failed to create new block subscriber")
)

var (
	ErrGetKeysAddress = errors.New("failed to get address for the observer's keys")
)

// newClientError joins the ErrClient error with an outer error (ErrClient...) and the actual
// inner client error.
func newClientError(outer, inner error) error {
	return fmt.Errorf("%w (%w): %w", ErrClient, outer, inner)
}
