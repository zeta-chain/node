//go:build !PRIVNET && !TESTNET
// +build !PRIVNET,!TESTNET

package update

import (
	"github.com/zeta-chain/zetacore/common"
)

func UpdateConfig() {

	updateEndpoint(common.EthChain(), "ETH_ENDPOINT")
	updateEndpoint(common.BscMainnetChain(), "BSC")
	updateEndpoint(common.PolygonChain(), "POLYGON_ENDPOINT")

	updateMPIAddress(common.EthChain(), "ETH_MPI_ADDRESS")
	updateMPIAddress(common.BscMainnetChain(), "BSC_MPI_ADDRESS")
	updateMPIAddress(common.PolygonChain(), "POLYGON_MPI_ADDRESS")

	updateTokenAddress(common.EthChain(), "ETH_ZETA_ADDRESS")
	updateTokenAddress(common.BscMainnetChain(), "BSC_ZETA_ADDRESS")
	updateTokenAddress(common.PolygonChain(), "POLYGON_ZETA_ADDRESS")
}
