//go:build PRIVNET
// +build PRIVNET

package update

import (
	"github.com/zeta-chain/zetacore/common"
)

func UpdateConfig() {

	updateEndpoint(common.GoerliLocalNetChain(), "GOERLILOCALNET_ENDPOINT")

	updateMPIAddress(common.GoerliLocalNetChain(), "GOERLILOCALNET_MPI_ENDPOINT")

	updateTokenAddress(common.GoerliLocalNetChain(), "GOERLILOCALNET_ZETA_ENDPOINT")

}
