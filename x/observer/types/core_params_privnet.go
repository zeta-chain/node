//go:build PRIVNET
// +build PRIVNET

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
				ConfirmationCount:           2,
				ZetaTokenContractAddress:    "0xA8D5060feb6B456e886F023709A2795373691E63",
				ConnectorContractAddress:    "0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9",
				Erc20CustodyContractAddress: "0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca",
				InTxTicker:                  2,
				OutTxTicker:                 2,
				WatchUtxoTicker:             0,
				GasPriceTicker:              5,
				OutboundTxScheduleInterval:  2,
				OutboundTxScheduleLookahead: 5,
			},
			{
				ChainId:                     common.BtcRegtestChain().ChainId,
				ConfirmationCount:           2,
				ZetaTokenContractAddress:    "",
				ConnectorContractAddress:    "",
				Erc20CustodyContractAddress: "",
				GasPriceTicker:              5,
				WatchUtxoTicker:             1,
				InTxTicker:                  1,
				OutTxTicker:                 2,
				OutboundTxScheduleInterval:  2,
				OutboundTxScheduleLookahead: 5,
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
