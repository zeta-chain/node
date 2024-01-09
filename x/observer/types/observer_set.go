package types

func (m *ObserverSet) Len() int {
	return len(m.ObserverList)
}

func (m *ObserverSet) LenUint() uint64 {
	return uint64(len(m.ObserverList))
}
