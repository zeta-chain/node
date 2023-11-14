package types

import (
	"fmt"

	"github.com/zeta-chain/zetacore/common"
)

func GetCoreParams() (CoreParamsList, error) {
	params := CoreParamsList{
		CoreParams: []*CoreParams{
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

	// check all core params correspond to a chain
	chainMap := make(map[int64]struct{})
	chainList := common.ExternalChainList()
	for _, chain := range chainList {
		chainMap[chain.ChainId] = struct{}{}
	}
	for _, param := range params.CoreParams {
		if _, ok := chainMap[param.ChainId]; !ok {
			return params, fmt.Errorf("chain id %d not found in chain list", param.ChainId)
		}
	}

	return params, nil
}
