package types

func (m Chain) IsEqual(chain Chain) bool {
	if m.ChainName == chain.ChainName && m.ChainId == chain.ChainId {
		return true
	}
	return false
}
