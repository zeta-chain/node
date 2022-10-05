package zetaobserver

import (
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
)

// TODO : modify this list with actual data before MainNet
func GetGenesisObserverList() []types.ObserverMapper {
	return []types.ObserverMapper{
		{
			ObserverChain:   types.ObserverChain_EthChainObserver,
			ObservationType: types.ObservationType_InBoundTx,
			ObserverList:    []string{"zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk", "zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50"},
		},
	}
}
