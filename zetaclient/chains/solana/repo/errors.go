package repo

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrFoundNoSignatures     = errors.New("found no signatures for the gateway")
	ErrUnsupportedTxVersion  = errors.New("unsupported transaction version")
	ErrAddressLookupTableRPC = errors.New("solana AddressLookupTable requires *rpc.Client")
)

var (
	ErrClient = errors.New("error calling a solana client function")

	ErrClientGetAccountInfo              = errors.New("failed to get account info")
	ErrClientGetBalance                  = errors.New("failed to get balance")
	ErrClientGetBlockTime                = errors.New("failed to get block time")
	ErrClientGetLatestBlockhash          = errors.New("failed to get latest blockhash")
	ErrClientGetRecentPrioritizationFees = errors.New("failed to get recent prioritization fees")
	ErrClientGetSignaturesForAddress     = errors.New("failed to get signatures for address")
	ErrClientGetSlot                     = errors.New("failed to get slot")
	ErrClientGetTransaction              = errors.New("failed to get transaction")
	ErrClientSendTransaction             = errors.New("failed to send transaction")

	ErrGetAddressLookupTableState = errors.New("failed to get address-lookup-table state")
)

// newClientError joins the ErrClient error with an outer error (ErrClient...) and the actual
// inner client error.
func newClientError(outer, inner error) error {
	return fmt.Errorf("%w (%w): %w", ErrClient, outer, inner)
}
