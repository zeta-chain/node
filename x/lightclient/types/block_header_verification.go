package types

func (b *BlockHeaderVerification) EnableChain(chainID int64) {
	found := false
	for i, enabledChain := range b.EnabledChains {
		if enabledChain.ChainId == chainID {
			b.EnabledChains[i].Enabled = true
			found = true
		}
	}
	if !found {
		b.EnabledChains = append(b.EnabledChains, EnabledChain{
			ChainId: chainID,
			Enabled: true,
		})
	}
}

func (b *BlockHeaderVerification) DisableChain(chainID int64) {
	found := false
	for i, v := range b.EnabledChains {
		if v.ChainId == chainID {
			b.EnabledChains[i].Enabled = false
			found = true
		}
	}
	if !found {
		b.EnabledChains = append(b.EnabledChains, EnabledChain{
			ChainId: chainID,
			Enabled: false,
		})
	}
}

func (b *BlockHeaderVerification) IsChainEnabled(chainID int64) bool {
	for _, v := range b.EnabledChains {
		if v.ChainId == chainID {
			return v.Enabled
		}
	}
	return false
}

func (b *BlockHeaderVerification) GetEnabledChainIDList() []int64 {
	var enabledChains []int64
	for _, v := range b.EnabledChains {
		if v.Enabled {
			enabledChains = append(enabledChains, v.ChainId)
		}
	}
	return enabledChains
}

func (b *BlockHeaderVerification) GetVerificationFlags() []EnabledChain {
	if b == nil || b.EnabledChains == nil {
		return []EnabledChain{}
	}
	return b.EnabledChains
}
