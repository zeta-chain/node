package zrepo

import "errors"

var (
	ErrClient = errors.New("zetacore client error")

	ErrClientGetCCTXByNonce      = errors.New("failed to get CCTX by nonce from zetacore")
	ErrClientGetOutboundTrackers = errors.New("failed to get outbound trackers from zetacore")
	ErrClientPostVoteOutbound    = errors.New("failed to post outbound vote with zetacore")
)

// newClientError joins the ErrClient error with a description error (ErrClient...) and an inner
// client error.
func newClientError(err, inner error) error {
	return errors.Join(ErrClient, err, inner)
}
