package config

import (
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/types"
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

	// MPI contract addresses
	ETH_MPI_ADDRESS     = "0xAA8e5c4b142d5c8ab75e7710E579C8e76cc15b85"
	POLYGON_MPI_ADDRESS = "0x6b8344F178CaaAc2068967c7ee6a067c3F9dC9AC"
	BSC_MPI_ADDRESS     = "0xdb9828f2bB2f9ab7C876768618768cC51fE10BAc"
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
	MAX_BLOCKS_PER_PERIOD       = 2000
	TIMEOUT_THRESHOLD_FOR_RETRY = 12 // 120 blocks for Zetachain; roughly 600s or 10min
)

const (
	ETH_ZETA_ABI = `[
			{
				"inputs": [
					{
						"internalType": "uint256",
						"name": "initialSupply",
						"type": "uint256"
					},
					{
						"internalType": "string",
						"name": "name",
						"type": "string"
					},
					{
						"internalType": "string",
						"name": "symbol",
						"type": "string"
					}
				],
				"stateMutability": "nonpayable",
				"type": "constructor"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "owner",
						"type": "address"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "spender",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "value",
						"type": "uint256"
					}
				],
				"name": "Approval",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "from",
						"type": "address"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "to",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "value",
						"type": "uint256"
					}
				],
				"name": "Transfer",
				"type": "event"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "owner",
						"type": "address"
					},
					{
						"internalType": "address",
						"name": "spender",
						"type": "address"
					}
				],
				"name": "allowance",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "spender",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "approve",
				"outputs": [
					{
						"internalType": "bool",
						"name": "",
						"type": "bool"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "account",
						"type": "address"
					}
				],
				"name": "balanceOf",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "decimals",
				"outputs": [
					{
						"internalType": "uint8",
						"name": "",
						"type": "uint8"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "name",
				"outputs": [
					{
						"internalType": "string",
						"name": "",
						"type": "string"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "symbol",
				"outputs": [
					{
						"internalType": "string",
						"name": "",
						"type": "string"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "totalSupply",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "recipient",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "transfer",
				"outputs": [
					{
						"internalType": "bool",
						"name": "",
						"type": "bool"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"internalType": "address",
						"name": "recipient",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "transferFrom",
				"outputs": [
					{
						"internalType": "bool",
						"name": "",
						"type": "bool"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			}
		]`
	MPI_ABI_STRING = `
[
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "zetaAddress",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "_TSSAddress",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "_TSSAddressUpdater",
				"type": "address"
			}
		],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			}
		],
		"name": "Paused",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			}
		],
		"name": "Unpaused",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "bytes",
				"name": "sender",
				"type": "bytes"
			},
			{
				"indexed": true,
				"internalType": "uint16",
				"name": "srcChainID",
				"type": "uint16"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "destContract",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "zetaAmount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "bytes",
				"name": "message",
				"type": "bytes"
			},
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "sendHash",
				"type": "bytes32"
			}
		],
		"name": "ZetaMessageReceiveEvent",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint16",
				"name": "destChainID",
				"type": "uint16"
			},
			{
				"indexed": false,
				"internalType": "bytes",
				"name": "destContract",
				"type": "bytes"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "zetaAmount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "gasLimit",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "bytes",
				"name": "message",
				"type": "bytes"
			},
			{
				"indexed": false,
				"internalType": "bytes",
				"name": "zetaParams",
				"type": "bytes"
			}
		],
		"name": "ZetaMessageSendEvent",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "TSSAddress",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "TSSAddressUpdater",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "ZETA_TOKEN",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "pause",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "paused",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "renounceTSSAddressUpdater",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "unpause",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_address",
				"type": "address"
			}
		],
		"name": "updateTSSAddress",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes",
				"name": "srcContract",
				"type": "bytes"
			},
			{
				"internalType": "uint16",
				"name": "srcChainID",
				"type": "uint16"
			},
			{
				"internalType": "address",
				"name": "destContract",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "zetaAmount",
				"type": "uint256"
			},
			{
				"internalType": "bytes",
				"name": "message",
				"type": "bytes"
			},
			{
				"internalType": "bytes32",
				"name": "sendHash",
				"type": "bytes32"
			}
		],
		"name": "zetaMessageReceive",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint16",
				"name": "destChainID",
				"type": "uint16"
			},
			{
				"internalType": "bytes",
				"name": "destContract",
				"type": "bytes"
			},
			{
				"internalType": "uint256",
				"name": "zetaAmount",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "gasLimit",
				"type": "uint256"
			},
			{
				"internalType": "bytes",
				"name": "message",
				"type": "bytes"
			},
			{
				"internalType": "bytes",
				"name": "zetaParams",
				"type": "bytes"
			}
		],
		"name": "zetaMessageSend",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`
)

var Chains = map[string]*types.ChainETHish{
	"ETH": {
		Name:               common.ETHChain,
		MPIContractAddress: "0x132b042bD5198a48E4D273f46b979E5f13Bd9239",
		ChainID:            5,
	},
	"BSC": {
		Name:               common.BSCChain,
		MPIContractAddress: "0x96cE47e42A73649CFe33d93D93ACFbEc6FD5ee14",
		ChainID:            97,
	},
	"POLYGON": {
		Name:               common.POLYGONChain,
		MPIContractAddress: "0x692E8A48634B530b4BFF1e621FC18C82F471892c",
		ChainID:            8001, // Should be 80001, but the chainid is uint16 (should be uint32 instead??)
	},
}

func FindChainByID(chainID uint16) string {
	for _, v := range Chains {
		if v.ChainID == chainID {
			return v.Name.String()
		}
	}
	return ""
}
