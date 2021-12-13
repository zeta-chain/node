package config

// ClientConfiguration
type ClientConfiguration struct {
	ChainHost       string `json:"chain_host" mapstructure:"chain_host"`
	ChainRPC        string `json:"chain_rpc" mapstructure:"chain_rpc"`
	ChainHomeFolder string `json:"chain_home_folder" mapstructure:"chain_home_folder"`
	SignerName      string `json:"signer_name" mapstructure:"signer_name"`
	SignerPasswd    string
}

const (
	ETH_CONFIRMATION_COUNT     = 5
	BSC_CONFIRMATION_COUNT     = 5
	POLYGON_CONFIRMATION_COUNT = 5
)

const (
	TSS_TEST_PRIVKEY = "2082bc9775d6ee5a05ef221a9d1c00b3cc3ecb274a4317acc0a182bc1e05d1bb"
	TSS_TEST_ADDRESS = "0xE80B6467863EbF8865092544f441da8fD3cF6074"
	TEST_RECEIVER    = "0x566bF3b1993FFd4BA134c107A63bb2aebAcCdbA0"
)

// Constants
// #nosec G101
const (
	// Routers
	ETH_ZETALOCK_ADDRESS  = "0x73Bd3bc3f9e54fe0Ca6cfD837e41D20a6e549cF3"
	ETH_TOKEN_ADDRESS     = "0x9BBff31cF0E0B7689470Fb2834295217813eCBfC"
	POLYGON_TOKEN_ADDRESS = "0xF146f7E7b9C404874b7B0f84E62886b721101E49"
	BSC_TOKEN_ADDRESS     = "0xC7e08ceb6bbaCaC1F4Fe1dA7EdBAEA949391a365"

	// API Endpoints
	ETH_ENDPOINT  = "https://eth-goerli.alchemyapi.io/v2/J-W7M8JtqtQI3ckka76fz9kxX-Sa_CSK"
	POLY_ENDPOINT = "https://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/polygon/mumbai/archive"
	BSC_ENDPOINT  = "https://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/bsc/testnet/archive"

	// Ticker timers
	ETH_BLOCK_TIME  = 12
	POLY_BLOCK_TIME = 5
	BSC_BLOCK_TIME  = 5

	// to catch up:
	MAX_BLOCKS_PER_PERIOD       = 64
	TIMEOUT_THRESHOLD_FOR_RETRY = 120 // 120 blocks for Zetachain; roughly 600s or 10min
)

const (
	NONETH_ZETA_ABI = `[
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
					},
					{
						"internalType": "address",
						"name": "oracleAddress",
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
					},
					{
						"internalType": "address",
						"name": "_OracleUpdater",
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
						"indexed": false,
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "string",
						"name": "receiver",
						"type": "string"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "wanted",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "string",
						"name": "chainid",
						"type": "string"
					},
					{
						"indexed": false,
						"internalType": "bytes",
						"name": "message",
						"type": "bytes"
					}
				],
				"name": "BurnSend",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": false,
						"internalType": "address",
						"name": "burnee",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "MBurnt",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": false,
						"internalType": "address",
						"name": "mintee",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "bytes32",
						"name": "sendHash",
						"type": "bytes32"
					}
				],
				"name": "MMinted",
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
				"inputs": [],
				"name": "BURN_TYPEHASH",
				"outputs": [
					{
						"internalType": "bytes32",
						"name": "",
						"type": "bytes32"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "DOMAIN_SEPARATOR",
				"outputs": [
					{
						"internalType": "bytes32",
						"name": "",
						"type": "bytes32"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "MAX_TOTAL_SUPPLY",
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
				"name": "MINT_TYPEHASH",
				"outputs": [
					{
						"internalType": "bytes32",
						"name": "",
						"type": "bytes32"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "OracleUpdater",
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
				"inputs": [
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "burn",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "burnFrom",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "string",
						"name": "receiver",
						"type": "string"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "wanted",
						"type": "uint256"
					},
					{
						"internalType": "string",
						"name": "chainid",
						"type": "string"
					},
					{
						"internalType": "bytes",
						"name": "message",
						"type": "bytes"
					}
				],
				"name": "burnSend",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "newOracleAddres",
						"type": "address"
					}
				],
				"name": "changeOracle",
				"outputs": [],
				"stateMutability": "nonpayable",
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
				"inputs": [
					{
						"internalType": "address",
						"name": "spender",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "subtractedValue",
						"type": "uint256"
					}
				],
				"name": "decreaseAllowance",
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
						"name": "spender",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "addedValue",
						"type": "uint256"
					}
				],
				"name": "increaseAllowance",
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
						"name": "mintee",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "value",
						"type": "uint256"
					},
					{
						"internalType": "bytes32",
						"name": "sendHash",
						"type": "bytes32"
					}
				],
				"name": "mint",
				"outputs": [],
				"stateMutability": "nonpayable",
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
				"inputs": [
					{
						"internalType": "address",
						"name": "",
						"type": "address"
					}
				],
				"name": "nonces",
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
						"name": "burnee",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "value",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "deadline",
						"type": "uint256"
					},
					{
						"internalType": "uint8",
						"name": "v",
						"type": "uint8"
					},
					{
						"internalType": "bytes32",
						"name": "r",
						"type": "bytes32"
					},
					{
						"internalType": "bytes32",
						"name": "s",
						"type": "bytes32"
					}
				],
				"name": "permitBurn",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "mintee",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "value",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "deadline",
						"type": "uint256"
					},
					{
						"internalType": "uint8",
						"name": "v",
						"type": "uint8"
					},
					{
						"internalType": "bytes32",
						"name": "r",
						"type": "bytes32"
					},
					{
						"internalType": "bytes32",
						"name": "s",
						"type": "bytes32"
					},
					{
						"internalType": "bytes32",
						"name": "sendHash",
						"type": "bytes32"
					}
				],
				"name": "permitMint",
				"outputs": [],
				"stateMutability": "nonpayable",
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
			},
			{
				"inputs": [
					{
						"internalType": "uint8",
						"name": "newFlexibility",
						"type": "uint8"
					}
				],
				"name": "updateSupplyOracleFlexibility",
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
			}
		]`
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
	ETH_ZETALOCK_ABI = `[
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "zetaAddress",
						"type": "address"
					},
					{
						"internalType": "address",
						"name": "oracleAddress",
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
					},
					{
						"internalType": "address",
						"name": "_OracleUpdater",
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
					},
					{
						"indexed": false,
						"internalType": "string",
						"name": "receiver",
						"type": "string"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "wanted",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "string",
						"name": "chainid",
						"type": "string"
					},
					{
						"indexed": false,
						"internalType": "bytes",
						"name": "message",
						"type": "bytes"
					}
				],
				"name": "LockSend",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": false,
						"internalType": "address",
						"name": "receiver",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "bytes32",
						"name": "sendHash",
						"type": "bytes32"
					}
				],
				"name": "Unlock",
				"type": "event"
			},
			{
				"inputs": [],
				"name": "OracleUpdater",
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
				"inputs": [
					{
						"internalType": "address",
						"name": "newOracleAddres",
						"type": "address"
					}
				],
				"name": "changeOracle",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "getLockedAmount",
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
						"internalType": "string",
						"name": "receiver",
						"type": "string"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "wanted",
						"type": "uint256"
					},
					{
						"internalType": "string",
						"name": "chainid",
						"type": "string"
					},
					{
						"internalType": "bytes",
						"name": "message",
						"type": "bytes"
					}
				],
				"name": "lockSend",
				"outputs": [],
				"stateMutability": "nonpayable",
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
				"inputs": [
					{
						"internalType": "address",
						"name": "receiver",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					},
					{
						"internalType": "bytes32",
						"name": "sendHash",
						"type": "bytes32"
					}
				],
				"name": "unlock",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "uint8",
						"name": "newFlexibility",
						"type": "uint8"
					}
				],
				"name": "updateSupplyOracleFlexibility",
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
			}
		]`
)
