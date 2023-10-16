//go:build !PRIVNET && !TESTNET && !MOCK_MAINNET
// +build !PRIVNET,!TESTNET,!MOCK_MAINNET

package types

import (
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/zeta-chain/zetacore/common"
)

func GetCoreParams() common.CoreParamsList {
	params := common.CoreParamsList{
		CoreParams: []*common.CoreParams{
			{
				ChainId:                     common.EthChain().ChainId,
				ConfirmationCount:           14,
				ZetaTokenContractAddress:    "",
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
				ChainId:                     common.BscMainnetChain().ChainId,
				ConfirmationCount:           14,
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
				ChainId:                     common.BtcMainnetChain().ChainId,
				ConfirmationCount:           2,
				ZetaTokenContractAddress:    "",
				ConnectorContractAddress:    "",
				Erc20CustodyContractAddress: "",
				WatchUtxoTicker:             30,
				InTxTicker:                  120,
				OutTxTicker:                 60,
				GasPriceTicker:              30,
				OutboundTxScheduleInterval:  30,
				OutboundTxScheduleLookahead: 60,
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
