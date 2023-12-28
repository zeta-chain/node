package types

import (
	"fmt"

	"github.com/zeta-chain/zetacore/common"
)

// GetCoreParams returns a list of default core params
func GetCoreParams() CoreParamsList {
	return CoreParamsList{
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
				ConfirmationCount:           1,
				ZetaTokenContractAddress:    "",
				ConnectorContractAddress:    "",
				Erc20CustodyContractAddress: "",
				GasPriceTicker:              5,
				WatchUtxoTicker:             1,
				InTxTicker:                  1,
				OutTxTicker:                 1,
				OutboundTxScheduleInterval:  1,
				OutboundTxScheduleLookahead: 5,
			},
			{
				ChainId:                     common.GoerliLocalnetChain().ChainId,
				ConfirmationCount:           1,
				ZetaTokenContractAddress:    "0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9",
				ConnectorContractAddress:    "0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca",
				Erc20CustodyContractAddress: "0xff3135df4F2775f4091b81f4c7B6359CfA07862a",
				InTxTicker:                  1,
				OutTxTicker:                 1,
				WatchUtxoTicker:             0,
				GasPriceTicker:              5,
				OutboundTxScheduleInterval:  1,
				OutboundTxScheduleLookahead: 5,
			},
		},
	}
}

// Validate checks all core params correspond to a chain and there is no duplicate chain id
func (cpl CoreParamsList) Validate() error {
	// check all core params correspond to a chain
	externalChainMap := make(map[int64]struct{})
	existingChainMap := make(map[int64]struct{})

	externalChainList := common.ExternalChainList()
	for _, chain := range externalChainList {
		externalChainMap[chain.ChainId] = struct{}{}
	}

	for _, param := range cpl.CoreParams {
		if _, ok := externalChainMap[param.ChainId]; !ok {
			return fmt.Errorf("chain id %d not found in chain list", param.ChainId)
		}
		if _, ok := existingChainMap[param.ChainId]; ok {
			return fmt.Errorf("duplicated chain id %d found", param.ChainId)
		}
		existingChainMap[param.ChainId] = struct{}{}
	}
	return nil
}
