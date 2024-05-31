package types

import "fmt"

func (b *BlockHeaderVerification) Validate() error {
	detectDuplicates := make(map[int64]bool)
	for _, chain := range b.HeaderSupportedChains {
		if _, ok := detectDuplicates[chain.ChainId]; ok {
			return fmt.Errorf("duplicated chain id for block header verification")
		}
		detectDuplicates[chain.ChainId] = true
	}
	return nil
}

// EnableChain enables block header verification for a specific chain
func (b *BlockHeaderVerification) EnableChain(chainID int64) {
	found := false
	for i, enabledChain := range b.HeaderSupportedChains {
		if enabledChain.ChainId == chainID {
			b.HeaderSupportedChains[i].Enabled = true
			found = true
		}
	}
	if !found {
		b.HeaderSupportedChains = append(b.HeaderSupportedChains, HeaderSupportedChain{
			ChainId: chainID,
			Enabled: true,
		})
	}
}

// DisableChain disables block header verification for a specific chain
// This function does not remove the chain from the list of enabled chains, it just disables it
// This keeps track of the chains tha support block header verification and also the ones that currently disabled or enabled
func (b *BlockHeaderVerification) DisableChain(chainID int64) {
	found := false
	for i, v := range b.HeaderSupportedChains {
		if v.ChainId == chainID {
			b.HeaderSupportedChains[i].Enabled = false
			found = true
		}
	}
	if !found {
		b.HeaderSupportedChains = append(b.HeaderSupportedChains, HeaderSupportedChain{
			ChainId: chainID,
			Enabled: false,
		})
	}
}

// IsChainEnabled checks if block header verification is enabled for a specific chain
// It returns true if the chain is enabled, false otherwise
// If the chain is not found in the list of chains, it returns false
func (b *BlockHeaderVerification) IsChainEnabled(chainID int64) bool {
	for _, v := range b.HeaderSupportedChains {
		if v.ChainId == chainID {
			return v.Enabled
		}
	}
	return false
}

// GetHeaderEnabledChainIDs returns a list of chain IDs that have block header verification enabled
func (b *BlockHeaderVerification) GetHeaderEnabledChainIDs() []int64 {
	var enabledChains []int64
	for _, v := range b.HeaderSupportedChains {
		if v.Enabled {
			enabledChains = append(enabledChains, v.ChainId)
		}
	}
	return enabledChains
}

// GetHeaderSupportedChainsList returns a list of chains that support block header verification
func (b *BlockHeaderVerification) GetHeaderSupportedChainsList() []HeaderSupportedChain {
	if b == nil || b.HeaderSupportedChains == nil {
		return []HeaderSupportedChain{}
	}
	return b.HeaderSupportedChains
}

// GetHeaderEnabledChains returns a list of chains that have block header verification enabled
func (b *BlockHeaderVerification) GetHeaderEnabledChains() []HeaderSupportedChain {
	var chains []HeaderSupportedChain
	if b == nil || b.HeaderSupportedChains == nil {
		return []HeaderSupportedChain{}
	}
	for _, chain := range b.HeaderSupportedChains {
		if chain.Enabled {
			chains = append(chains, chain)
		}
	}
	return chains
}
