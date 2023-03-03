package zetaclient

import (
	"github.com/zeta-chain/zetacore/common"
)

// Modify to update this from the core later
func GetSupportedChains() []*common.Chain {
	return common.DefaultChainsList()
}
