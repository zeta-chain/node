package types

import (
	"fmt"

	"github.com/zeta-chain/node/pkg/chains"
)

// DefaultChainInfo returns the structure with an empty list of chains
func DefaultChainInfo() ChainInfo {
	return ChainInfo{
		Chains: []chains.Chain{},
	}
}

// Validate performs basic validation of chain info
// It checks all chains are valid and they're all of external type
// The structure is used to store external chain information
func (ci ChainInfo) Validate() error {
	for _, chain := range ci.Chains {
		if err := chain.Validate(); err != nil {
			return err
		}
		if !chain.IsExternal {
			return fmt.Errorf("chain %d is not external", chain.ChainId)
		}
	}

	return nil
}
