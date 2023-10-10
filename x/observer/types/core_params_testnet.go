//go:build TESTNET
// +build TESTNET

package types

import (
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/zeta-chain/zetacore/common"
)

func GetCoreParams() CoreParamsList {
	params := CoreParamsList{
		CoreParams: []*CoreParams{
			{
				ChainId:           common.GoerliChain().ChainId,
				ConfirmationCount: 6,
				// This is the actual Zeta token Goerli testnet, we need to specify this address for the integration tests to pass
				ZetaTokenContractAddress:    "0x0000c304d2934c00db1d51995b9f6996affd17c0",
				ConnectorContractAddress:    "",
				Erc20CustodyContractAddress: "",
				InTxTicker:                  12,
				OutTxTicker:                 15,
				WatchUtxoTicker:             0,
				GasPriceTicker:              30,
				OutboundTxScheduleInterval:  30,
				OutboundTxScheduleLookahead: 60,
			},
			{
				ChainId:                     common.BscTestnetChain().ChainId,
				ConfirmationCount:           6,
				ZetaTokenContractAddress:    "",
				ConnectorContractAddress:    "",
				Erc20CustodyContractAddress: "",
				InTxTicker:                  5,
				OutTxTicker:                 15,
				WatchUtxoTicker:             0,
				GasPriceTicker:              30,
				OutboundTxScheduleInterval:  30,
				OutboundTxScheduleLookahead: 60,
			},
			{
				ChainId:                     common.MumbaiChain().ChainId,
				ConfirmationCount:           12,
				ZetaTokenContractAddress:    "",
				ConnectorContractAddress:    "",
				Erc20CustodyContractAddress: "",
				InTxTicker:                  2,
				OutTxTicker:                 15,
				WatchUtxoTicker:             0,
				GasPriceTicker:              30,
				OutboundTxScheduleInterval:  30,
				OutboundTxScheduleLookahead: 60,
			},
			{
				ChainId:                     common.BtcTestNetChain().ChainId,
				ConfirmationCount:           2,
				ZetaTokenContractAddress:    "",
				ConnectorContractAddress:    "",
				Erc20CustodyContractAddress: "",
				WatchUtxoTicker:             30,
				InTxTicker:                  120,
				OutTxTicker:                 12,
				GasPriceTicker:              30,
				OutboundTxScheduleInterval:  30,
				OutboundTxScheduleLookahead: 100,
			},
		},
	}
	chainList := common.ExternalChainList()
	requiredParams := len(chainList)
	availableParams := 0
	for _, chain := range chainList {
		for _, param := range params.CoreParams {
			if chain.ChainId == param.ChainId {
				availableParams++
			}
		}
	}
	if availableParams != requiredParams {
		panic(fmt.Sprintf("Core params are not available for all chains , DefaultChains : %s , CoreParams : %s",
			types.PrettyPrintStruct(chainList), params.String()))
	}
	return params
}
