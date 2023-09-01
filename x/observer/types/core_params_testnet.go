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
				ChainId:                     common.GoerliChain().ChainId,
				ConfirmationCount:           6,
				ZetaTokenContractAddress:    "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				ConnectorContractAddress:    "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				Erc20CustodyContractAddress: "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
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
				ZetaTokenContractAddress:    "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				ConnectorContractAddress:    "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				Erc20CustodyContractAddress: "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
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
				ZetaTokenContractAddress:    "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				ConnectorContractAddress:    "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				Erc20CustodyContractAddress: "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				InTxTicker:                  2,
				OutTxTicker:                 15,
				WatchUtxoTicker:             0,
				GasPriceTicker:              30,
				OutboundTxScheduleInterval:  30,
				OutboundTxScheduleLookahead: 60,
			},
			{
				ChainId:                     common.ZetaChain().ChainId,
				ConfirmationCount:           3,
				ZetaTokenContractAddress:    "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				ConnectorContractAddress:    "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				Erc20CustodyContractAddress: "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				InTxTicker:                  2,
				OutTxTicker:                 3,
				WatchUtxoTicker:             0,
				GasPriceTicker:              5,
			},
			{
				ChainId:                     common.BtcTestNetChain().ChainId,
				ConfirmationCount:           2,
				ZetaTokenContractAddress:    "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				ConnectorContractAddress:    "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				Erc20CustodyContractAddress: "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea",
				WatchUtxoTicker:             30,
				InTxTicker:                  120,
				OutTxTicker:                 60,
				GasPriceTicker:              30,
				OutboundTxScheduleInterval:  30,
				OutboundTxScheduleLookahead: 60,
			},
		},
	}
	chainList := common.DefaultChainsList()
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
