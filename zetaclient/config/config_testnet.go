//go:build TESTNET
// +build TESTNET

package config

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/zeta-chain/zetacore/common"
)

const (
	TssTestPrivkey = "2082bc9775d6ee5a05ef221a9d1c00b3cc3ecb274a4317acc0a182bc1e05d1bb"
	TssTestAddress = "0xE80B6467863EbF8865092544f441da8fD3cF6074"
	//TestReceiver  = "0x566bF3b1993FFd4BA134c107A63bb2aebAcCdbA0"
)

// Constants
// #nosec G101
const (

	// Ticker timers
	EthBlockTime     = 12
	PolygonBlockTime = 2
	BscBlockTime     = 5

	// to catch up:
	MaxBlocksPerPeriod = 100
)

const (
	ConnectorAbiString = `
[{"inputs":[{"internalType":"address","name":"_zetaTokenAddress","type":"address"},{"internalType":"address","name":"_tssAddress","type":"address"},{"internalType":"address","name":"_tssAddressUpdater","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Paused","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Unpaused","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"originSenderAddress","type":"bytes"},{"indexed":true,"internalType":"uint256","name":"originChainId","type":"uint256"},{"indexed":true,"internalType":"address","name":"destinationAddress","type":"address"},{"indexed":false,"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"},{"indexed":true,"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"ZetaReceived","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"originSenderAddress","type":"address"},{"indexed":false,"internalType":"uint256","name":"originChainId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"indexed":true,"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"},{"indexed":true,"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"ZetaReverted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"originSenderAddress","type":"address"},{"indexed":false,"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"gasLimit","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"zetaParams","type":"bytes"}],"name":"ZetaSent","type":"event"},{"inputs":[],"name":"getLockedAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"originSenderAddress","type":"bytes"},{"internalType":"uint256","name":"originChainId","type":"uint256"},{"internalType":"address","name":"destinationAddress","type":"address"},{"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"},{"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"onReceive","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"originSenderAddress","type":"address"},{"internalType":"uint256","name":"originChainId","type":"uint256"},{"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"},{"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"onRevert","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"pause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"paused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceTssAddressUpdater","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"components":[{"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"internalType":"uint256","name":"gasLimit","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"},{"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"internalType":"bytes","name":"zetaParams","type":"bytes"}],"internalType":"struct ZetaInterfaces.SendInput","name":"input","type":"tuple"}],"name":"send","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"tssAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"tssAddressUpdater","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"unpause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_tssAddress","type":"address"}],"name":"updateTssAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"zetaToken","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`
	ERC20CustodyAbiString = `
[{"inputs":[{"internalType":"address","name":"_TSSAddress","type":"address"},{"internalType":"address","name":"_TSSAddressUpdater","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"InvalidSender","type":"error"},{"inputs":[],"name":"InvalidTSSUpdater","type":"error"},{"inputs":[],"name":"IsPaused","type":"error"},{"inputs":[],"name":"NotPaused","type":"error"},{"inputs":[],"name":"NotWhitelisted","type":"error"},{"inputs":[],"name":"ZeroAddress","type":"error"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"recipient","type":"bytes"},{"indexed":false,"internalType":"address","name":"asset","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"}],"name":"Deposited","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"sender","type":"address"}],"name":"Paused","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"sender","type":"address"}],"name":"Unpaused","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"asset","type":"address"}],"name":"Unwhitelisted","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"asset","type":"address"}],"name":"Whitelisted","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"recipient","type":"address"},{"indexed":false,"internalType":"address","name":"asset","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"Withdrawn","type":"event"},{"inputs":[],"name":"TSSAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TSSAddressUpdater","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"recipient","type":"bytes"},{"internalType":"address","name":"asset","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"}],"name":"deposit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"pause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"paused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceTSSAddressUpdater","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"unpause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"asset","type":"address"}],"name":"unwhitelist","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_address","type":"address"}],"name":"updateTSSAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"asset","type":"address"}],"name":"whitelist","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"whitelisted","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"address","name":"asset","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

var (
	BitconNetParams = &chaincfg.Testnet3Params
)

var ChainsEnabled = []common.Chain{}

var ChainConfigs = map[string]*ChainETHish{
	common.GoerliChain().ChainName.String(): {
		Chain:                    common.GoerliChain(),
		ConnectorContractAddress: "0x851b2446f225266C4EC3cd665f6801D624626c4D",
		ZETATokenContractAddress: "0xfF8dee1305D6200791e26606a0b04e12C5292aD8",
		BlockTime:                EthBlockTime,
		Endpoint:                 "https://eth-goerli-sh285ns91n5975.athens.zetachain.com",
		ConfCount:                15,
	},
	common.BscTestnetChain().ChainName.String(): {
		Chain:                       common.BscTestnetChain(),
		ConnectorContractAddress:    "0xcF1B4B432CA02D6418a818044d38b18CDd3682E9",
		ZETATokenContractAddress:    "0x33580e10212342d0aA66C9de3F6F6a4AfefA144C",
		ERC20CustodyContractAddress: "0x0e141A7e7C0A7E15E7d22713Fc0a6187515Fa9BF",
		BlockTime:                   BscBlockTime,
		Endpoint:                    "https://bsc-sh285ns91n5975.athens.zetachain.com",
		ConfCount:                   15,
	},
	common.MumbaiChain().ChainName.String(): {
		Chain:                       common.MumbaiChain(),
		ConnectorContractAddress:    "0xED4d7f8cA6252Ccf85A1eFB5444d7dB794ddD328",
		ZETATokenContractAddress:    "0xBaEF590c5Aef9881b0a5C86e18D35432218C64D5",
		ERC20CustodyContractAddress: "0x0e141A7e7C0A7E15E7d22713Fc0a6187515Fa9BF",
		BlockTime:                   PolygonBlockTime,
		Endpoint:                    "https://mumbai-sh285ns91n5975.athens.zetachain.com",
		ConfCount:                   128,
	},
	common.BaobabChain().ChainName.String(): {
		Chain:                       common.BaobabChain(),
		ConnectorContractAddress:    "0x000054d3A0Bc83Ec7808F52fCdC28A96c89F6C5c",
		ZETATokenContractAddress:    "0x000080383847bD75F91c168269Aa74004877592f",
		ERC20CustodyContractAddress: "0x0e141A7e7C0A7E15E7d22713Fc0a6187515Fa9BF",
		BlockTime:                   EthBlockTime,
		Endpoint:                    "https://baobab-sh285ns91n5975.athens.zetachain.com",
		ConfCount:                   15,
	},

	common.ZetaChain().ChainName.String(): {
		Chain:                    common.ZetaChain(),
		BlockTime:                6,
		ZETATokenContractAddress: "0x2DD9830f8Ac0E421aFF9B7c8f7E9DF6F65DBF6Ea",
		ConfCount:                3,
	},
}

var BitcoinConfig = &ChainBitcoinish{
	RPCUsername: "smoketest",
	RPCPassword: "123",
	RPCEndpoint: "bitcoin:18443",
	RPCParams:   "regtest",

	WatchInTxPeriod:     5,
	WatchGasPricePeriod: 5,
	WatchUTXOSPeriod:    5,
}

func New() Config {
	return Config{
		ChainsEnabled: ChainsEnabled,
		ChainConfigs:  ChainConfigs,
		BitcoinConfig: BitcoinConfig,
	}
}
