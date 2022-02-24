package main

const ABI_MPI = `
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
