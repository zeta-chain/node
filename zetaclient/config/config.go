package config

import (
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/types"
	"math/big"
)

// ClientConfiguration
type ClientConfiguration struct {
	ChainHost       string `json:"chain_host" mapstructure:"chain_host"`
	ChainRPC        string `json:"chain_rpc" mapstructure:"chain_rpc"`
	ChainHomeFolder string `json:"chain_home_folder" mapstructure:"chain_home_folder"`
	SignerName      string `json:"signer_name" mapstructure:"signer_name"`
	SignerPasswd    string
}

const (
	ETH_CONFIRMATION_COUNT     = 3
	BSC_CONFIRMATION_COUNT     = 5
	POLYGON_CONFIRMATION_COUNT = 5
)

var (
	// API Endpoints
	//ETH_ENDPOINT  = "https://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/eth/goerli/archive"
	ETH_ENDPOINT  = "https://eth-goerli.alchemyapi.io/v2/J-W7M8JtqtQI3ckka76fz9kxX-Sa_CSK"
	POLY_ENDPOINT = "https://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/polygon/mumbai/archive"
	BSC_ENDPOINT  = "https://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/bsc/testnet/archive"
)

const (
	TSS_TEST_PRIVKEY = "2082bc9775d6ee5a05ef221a9d1c00b3cc3ecb274a4317acc0a182bc1e05d1bb"
	TSS_TEST_ADDRESS = "0xE80B6467863EbF8865092544f441da8fD3cF6074"
	TEST_RECEIVER    = "0x566bF3b1993FFd4BA134c107A63bb2aebAcCdbA0"
)

// Constants
// #nosec G101
const (

	// Ticker timers
	ETH_BLOCK_TIME  = 12
	POLY_BLOCK_TIME = 10
	BSC_BLOCK_TIME  = 10

	// to catch up:
	MAX_BLOCKS_PER_PERIOD = 2000
)

const (
	MPI_ABI_STRING = `
[{"inputs":[{"internalType":"address","name":"_zetaTokenAddress","type":"address"},{"internalType":"address","name":"_tssAddress","type":"address"},{"internalType":"address","name":"_tssAddressUpdater","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Paused","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Unpaused","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"originSenderAddress","type":"bytes"},{"indexed":true,"internalType":"uint256","name":"originChainId","type":"uint256"},{"indexed":true,"internalType":"address","name":"destinationAddress","type":"address"},{"indexed":false,"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"},{"indexed":true,"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"ZetaReceived","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"originSenderAddress","type":"address"},{"indexed":false,"internalType":"uint256","name":"originChainId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"indexed":true,"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"},{"indexed":true,"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"ZetaReverted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"originSenderAddress","type":"address"},{"indexed":false,"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"gasLimit","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"zetaParams","type":"bytes"}],"name":"ZetaSent","type":"event"},{"inputs":[],"name":"getLockedAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"originSenderAddress","type":"bytes"},{"internalType":"uint256","name":"originChainId","type":"uint256"},{"internalType":"address","name":"destinationAddress","type":"address"},{"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"},{"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"onReceive","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"originSenderAddress","type":"address"},{"internalType":"uint256","name":"originChainId","type":"uint256"},{"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"},{"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"onRevert","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"pause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"paused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceTssAddressUpdater","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"components":[{"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"internalType":"uint256","name":"gasLimit","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"},{"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"internalType":"bytes","name":"zetaParams","type":"bytes"}],"internalType":"struct ZetaInterfaces.SendInput","name":"input","type":"tuple"}],"name":"send","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"tssAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"tssAddressUpdater","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"unpause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_tssAddress","type":"address"}],"name":"updateTssAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"zetaToken","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]
`
)

var Chains = map[string]*types.ChainETHish{
	"ETH": {
		Name:               common.ETHChain,
		MPIContractAddress: "0x132b042bD5198a48E4D273f46b979E5f13Bd9239",
		ChainID:            big.NewInt(5),
	},
	"BSC": {
		Name:               common.BSCChain,
		MPIContractAddress: "0x96cE47e42A73649CFe33d93D93ACFbEc6FD5ee14",
		ChainID:            big.NewInt(97),
	},
	"POLYGON": {
		Name:               common.POLYGONChain,
		MPIContractAddress: "0x692E8A48634B530b4BFF1e621FC18C82F471892c",
		ChainID:            big.NewInt(80001),
	},
}

func FindChainByID(id *big.Int) string {
	for _, v := range Chains {
		if v.ChainID.Cmp(id) == 0 {
			return v.Name.String()
		}
	}
	return ""
}
