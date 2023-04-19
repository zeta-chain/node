//go:build TESTNET
// +build TESTNET

package types

import "github.com/zeta-chain/zetacore/common"

func GetClientParams() map[int64]ClientParams {
	var ChainConfigs = map[int64]ClientParams{
		common.GoerliChain().ChainId: {
			ConfirmationCount:           14,
			GasPriceTicker:              5,
			ZetaTokenContractAddress:    "0xfF8dee1305D6200791e26606a0b04e12C5292aD8",
			ConnectorContractAddress:    "0x851b2446f225266C4EC3cd665f6801D624626c4D",
			Erc20CustodyContractAddress: "",
		},
		common.BscTestnetChain().ChainId: {
			ConfirmationCount:           14,
			GasPriceTicker:              5,
			ZetaTokenContractAddress:    "0x33580e10212342d0aA66C9de3F6F6a4AfefA144C",
			ConnectorContractAddress:    "0xcF1B4B432CA02D6418a818044d38b18CDd3682E9",
			Erc20CustodyContractAddress: "0x0e141A7e7C0A7E15E7d22713Fc0a6187515Fa9BF",
		},
		common.MumbaiChain().ChainId: {
			ConfirmationCount:           14,
			GasPriceTicker:              5,
			ZetaTokenContractAddress:    "0xfF8dee1305D6200791e26606a0b04e12C5292aD8",
			ConnectorContractAddress:    "0x851b2446f225266C4EC3cd665f6801D624626c4D",
			Erc20CustodyContractAddress: "",
		},
		common.BaobabChain().ChainId: {
			ConfirmationCount:           14,
			GasPriceTicker:              5,
			ZetaTokenContractAddress:    "0xfF8dee1305D6200791e26606a0b04e12C5292aD8",
			ConnectorContractAddress:    "0x851b2446f225266C4EC3cd665f6801D624626c4D",
			Erc20CustodyContractAddress: "",
		},

		common.ZetaChain().ChainId: {
			ConfirmationCount:           14,
			GasPriceTicker:              5,
			ZetaTokenContractAddress:    "0xfF8dee1305D6200791e26606a0b04e12C5292aD8",
			ConnectorContractAddress:    "0x851b2446f225266C4EC3cd665f6801D624626c4D",
			Erc20CustodyContractAddress: "",
		},
	}

}
