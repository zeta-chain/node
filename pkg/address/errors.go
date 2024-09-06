package address

import (
	"fmt"
)

// InvalidAddressError represents an error for an invalid address.
type InvalidAddressError struct {
	Address string
	Msg     string
}

func (e *InvalidAddressError) Error() string {
	return fmt.Sprintf("Invalid address: %s, %s", e.Address, e.Msg)
}

// InvalidChainError represents an error for an invalid chain ID.
type InvalidChainError struct {
	ChainID int64
	Msg     string
}

func (e *InvalidChainError) Error() string {
	return fmt.Sprintf("Invalid chain ID: %d, %s", e.ChainID, e.Msg)
}

// InvalidNetworkError represents an error for an invalid network.
type InvalidNetworkError struct {
	Network string
	Msg     string
}

func (e *InvalidNetworkError) Error() string {
	return fmt.Sprintf("Invalid network: %s, %s", e.Network, e.Msg)
}
