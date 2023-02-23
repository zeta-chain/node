//go:build TESTNET
// +build TESTNET

package main

import "github.com/zeta-chain/zetacore/common"

func updateConfig() {

	updateEndpoint(common.GoerliChain(), "GOERLI_ENDPOINT")
	updateEndpoint(common.BscTestnetChain(), "BSCTESTNET_ENDPOINT")
	updateEndpoint(common.MumbaiChain(), "MUMBAI_ENDPOINT")
	updateEndpoint(common.BaobabChain(), "BAOBAB_ENDPOINT")

	updateMPIAddress(common.GoerliChain(), "GOERLI_MPI_ADDRESS")
	updateMPIAddress(common.BscTestnetChain(), "BSCTESTNET_MPI_ADDRESS")
	updateMPIAddress(common.MumbaiChain(), "MUMBAI_MPI_ADDRESS")
	updateMPIAddress(common.BaobabChain(), "BAOBAB_MPI_ADDRESS")

	updateTokenAddress(common.GoerliChain(), "GOERLI_ZETA_ADDRESS")
	updateTokenAddress(common.BscTestnetChain(), "BSCTESTNET_ZETA_ADDRESS")
	updateTokenAddress(common.MumbaiChain(), "MUMBAI_ZETA_ADDRESS")
	updateTokenAddress(common.BaobabChain(), "BAOBAB_ZETA_ADDRESS")
}
