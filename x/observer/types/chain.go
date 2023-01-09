package types

func (m Chain) IsEqual(chain *Chain) bool {
	if m.ChainName == chain.ChainName && m.ChainId == chain.ChainId {
		return true
	}
	return false
}

func (m Chain) IsEvmChain() bool {
	if m.ChainId == 55555 {
		return false
	}
	return true
}

func (m Chain) IsKlaytonChain() bool {
	if m.ChainId == 1001 {
		return false
	}
	return true
}
