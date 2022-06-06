import { assert, expect } from "chai";
import { EVMChain, ZetaChain, wait } from "lib/Blockchain";
import { describe } from "mocha";
import {
  testAllowance,
  testZetaTokenTransfer,
  testSendMessage,
  checkZetaTxStatus,
} from "../lib/TestFunctions";

// The following two lines disable debug and info logging - comment them out to enable
console.debug = function () {}; // Disables Debug Level Logging
// console.info = function () {}; // Disables Info Level Logging


const MPIContractABINonEth = [
  {
    inputs: [
      { internalType: "address", name: "_zetaTokenAddress", type: "address" },
      { internalType: "address", name: "_tssAddress", type: "address" },
      { internalType: "address", name: "_tssAddressUpdater", type: "address" },
    ],
    stateMutability: "nonpayable",
    type: "constructor",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "address",
        name: "account",
        type: "address",
      },
    ],
    name: "Paused",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "address",
        name: "account",
        type: "address",
      },
    ],
    name: "Unpaused",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "bytes",
        name: "originSenderAddress",
        type: "bytes",
      },
      {
        indexed: true,
        internalType: "uint256",
        name: "originChainId",
        type: "uint256",
      },
      {
        indexed: true,
        internalType: "address",
        name: "destinationAddress",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "zetaAmount",
        type: "uint256",
      },
      { indexed: false, internalType: "bytes", name: "message", type: "bytes" },
      {
        indexed: true,
        internalType: "bytes32",
        name: "internalSendHash",
        type: "bytes32",
      },
    ],
    name: "ZetaReceived",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "address",
        name: "originSenderAddress",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "originChainId",
        type: "uint256",
      },
      {
        indexed: true,
        internalType: "uint256",
        name: "destinationChainId",
        type: "uint256",
      },
      {
        indexed: true,
        internalType: "bytes",
        name: "destinationAddress",
        type: "bytes",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "zetaAmount",
        type: "uint256",
      },
      { indexed: false, internalType: "bytes", name: "message", type: "bytes" },
      {
        indexed: true,
        internalType: "bytes32",
        name: "internalSendHash",
        type: "bytes32",
      },
    ],
    name: "ZetaReverted",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "originSenderAddress",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "destinationChainId",
        type: "uint256",
      },
      {
        indexed: false,
        internalType: "bytes",
        name: "destinationAddress",
        type: "bytes",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "zetaAmount",
        type: "uint256",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "gasLimit",
        type: "uint256",
      },
      { indexed: false, internalType: "bytes", name: "message", type: "bytes" },
      {
        indexed: false,
        internalType: "bytes",
        name: "zetaParams",
        type: "bytes",
      },
    ],
    name: "ZetaSent",
    type: "event",
  },
  {
    inputs: [],
    name: "getLockedAmount",
    outputs: [{ internalType: "uint256", name: "", type: "uint256" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      { internalType: "bytes", name: "originSenderAddress", type: "bytes" },
      { internalType: "uint256", name: "originChainId", type: "uint256" },
      { internalType: "address", name: "destinationAddress", type: "address" },
      { internalType: "uint256", name: "zetaAmount", type: "uint256" },
      { internalType: "bytes", name: "message", type: "bytes" },
      { internalType: "bytes32", name: "internalSendHash", type: "bytes32" },
    ],
    name: "onReceive",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "originSenderAddress", type: "address" },
      { internalType: "uint256", name: "originChainId", type: "uint256" },
      { internalType: "bytes", name: "destinationAddress", type: "bytes" },
      { internalType: "uint256", name: "destinationChainId", type: "uint256" },
      { internalType: "uint256", name: "zetaAmount", type: "uint256" },
      { internalType: "bytes", name: "message", type: "bytes" },
      { internalType: "bytes32", name: "internalSendHash", type: "bytes32" },
    ],
    name: "onRevert",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "pause",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "paused",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "renounceTssAddressUpdater",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        components: [
          {
            internalType: "uint256",
            name: "destinationChainId",
            type: "uint256",
          },
          { internalType: "bytes", name: "destinationAddress", type: "bytes" },
          { internalType: "uint256", name: "gasLimit", type: "uint256" },
          { internalType: "bytes", name: "message", type: "bytes" },
          { internalType: "uint256", name: "zetaAmount", type: "uint256" },
          { internalType: "bytes", name: "zetaParams", type: "bytes" },
        ],
        internalType: "struct ZetaInterfaces.SendInput",
        name: "input",
        type: "tuple",
      },
    ],
    name: "send",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "tssAddress",
    outputs: [{ internalType: "address", name: "", type: "address" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "tssAddressUpdater",
    outputs: [{ internalType: "address", name: "", type: "address" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "unpause",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [{ internalType: "address", name: "_tssAddress", type: "address" }],
    name: "updateTssAddress",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "zetaToken",
    outputs: [{ internalType: "address", name: "", type: "address" }],
    stateMutability: "view",
    type: "function",
  },
];
const ZetaContractABINonEth = [
  {
    inputs: [
      { internalType: "uint256", name: "initialSupply", type: "uint256" },
      { internalType: "string", name: "name", type: "string" },
      { internalType: "string", name: "symbol", type: "string" },
      { internalType: "address", name: "_TSSAddress", type: "address" },
      { internalType: "address", name: "_TSSAddressUpdater", type: "address" },
    ],
    stateMutability: "nonpayable",
    type: "constructor",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "owner",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "spender",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "value",
        type: "uint256",
      },
    ],
    name: "Approval",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "burnee",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
    ],
    name: "MBurnt",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "mintee",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
      {
        indexed: true,
        internalType: "bytes32",
        name: "sendHash",
        type: "bytes32",
      },
    ],
    name: "MMinted",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      { indexed: true, internalType: "address", name: "from", type: "address" },
      { indexed: true, internalType: "address", name: "to", type: "address" },
      {
        indexed: false,
        internalType: "uint256",
        name: "value",
        type: "uint256",
      },
    ],
    name: "Transfer",
    type: "event",
  },
  {
    inputs: [],
    name: "MPIAddress",
    outputs: [{ internalType: "address", name: "", type: "address" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "TSSAddress",
    outputs: [{ internalType: "address", name: "", type: "address" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "TSSAddressUpdater",
    outputs: [{ internalType: "address", name: "", type: "address" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "owner", type: "address" },
      { internalType: "address", name: "spender", type: "address" },
    ],
    name: "allowance",
    outputs: [{ internalType: "uint256", name: "", type: "uint256" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "spender", type: "address" },
      { internalType: "uint256", name: "amount", type: "uint256" },
    ],
    name: "approve",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [{ internalType: "address", name: "account", type: "address" }],
    name: "balanceOf",
    outputs: [{ internalType: "uint256", name: "", type: "uint256" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [{ internalType: "uint256", name: "amount", type: "uint256" }],
    name: "burn",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "account", type: "address" },
      { internalType: "uint256", name: "amount", type: "uint256" },
    ],
    name: "burnFrom",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "decimals",
    outputs: [{ internalType: "uint8", name: "", type: "uint8" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "spender", type: "address" },
      { internalType: "uint256", name: "subtractedValue", type: "uint256" },
    ],
    name: "decreaseAllowance",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "spender", type: "address" },
      { internalType: "uint256", name: "addedValue", type: "uint256" },
    ],
    name: "increaseAllowance",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "mintee", type: "address" },
      { internalType: "uint256", name: "value", type: "uint256" },
      { internalType: "bytes32", name: "sendHash", type: "bytes32" },
    ],
    name: "mint",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "name",
    outputs: [{ internalType: "string", name: "", type: "string" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [{ internalType: "address", name: "", type: "address" }],
    name: "nonces",
    outputs: [{ internalType: "uint256", name: "", type: "uint256" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "renounceTSSAddressUpdater",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "symbol",
    outputs: [{ internalType: "string", name: "", type: "string" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "totalSupply",
    outputs: [{ internalType: "uint256", name: "", type: "uint256" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "recipient", type: "address" },
      { internalType: "uint256", name: "amount", type: "uint256" },
    ],
    name: "transfer",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "sender", type: "address" },
      { internalType: "address", name: "recipient", type: "address" },
      { internalType: "uint256", name: "amount", type: "uint256" },
    ],
    name: "transferFrom",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "_tss", type: "address" },
      { internalType: "address", name: "_mpi", type: "address" },
    ],
    name: "updateTSSAndMPIAddresses",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
];
const MPIContractABIEth = [
  {
    inputs: [
      { internalType: "address", name: "_zetaTokenAddress", type: "address" },
      { internalType: "address", name: "_tssAddress", type: "address" },
      { internalType: "address", name: "_tssAddressUpdater", type: "address" },
    ],
    stateMutability: "nonpayable",
    type: "constructor",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "address",
        name: "account",
        type: "address",
      },
    ],
    name: "Paused",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "address",
        name: "account",
        type: "address",
      },
    ],
    name: "Unpaused",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "bytes",
        name: "originSenderAddress",
        type: "bytes",
      },
      {
        indexed: true,
        internalType: "uint256",
        name: "originChainId",
        type: "uint256",
      },
      {
        indexed: true,
        internalType: "address",
        name: "destinationAddress",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "zetaAmount",
        type: "uint256",
      },
      { indexed: false, internalType: "bytes", name: "message", type: "bytes" },
      {
        indexed: true,
        internalType: "bytes32",
        name: "internalSendHash",
        type: "bytes32",
      },
    ],
    name: "ZetaReceived",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "address",
        name: "originSenderAddress",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "originChainId",
        type: "uint256",
      },
      {
        indexed: true,
        internalType: "uint256",
        name: "destinationChainId",
        type: "uint256",
      },
      {
        indexed: true,
        internalType: "bytes",
        name: "destinationAddress",
        type: "bytes",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "zetaAmount",
        type: "uint256",
      },
      { indexed: false, internalType: "bytes", name: "message", type: "bytes" },
      {
        indexed: true,
        internalType: "bytes32",
        name: "internalSendHash",
        type: "bytes32",
      },
    ],
    name: "ZetaReverted",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "originSenderAddress",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "destinationChainId",
        type: "uint256",
      },
      {
        indexed: false,
        internalType: "bytes",
        name: "destinationAddress",
        type: "bytes",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "zetaAmount",
        type: "uint256",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "gasLimit",
        type: "uint256",
      },
      { indexed: false, internalType: "bytes", name: "message", type: "bytes" },
      {
        indexed: false,
        internalType: "bytes",
        name: "zetaParams",
        type: "bytes",
      },
    ],
    name: "ZetaSent",
    type: "event",
  },
  {
    inputs: [],
    name: "getLockedAmount",
    outputs: [{ internalType: "uint256", name: "", type: "uint256" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      { internalType: "bytes", name: "originSenderAddress", type: "bytes" },
      { internalType: "uint256", name: "originChainId", type: "uint256" },
      { internalType: "address", name: "destinationAddress", type: "address" },
      { internalType: "uint256", name: "zetaAmount", type: "uint256" },
      { internalType: "bytes", name: "message", type: "bytes" },
      { internalType: "bytes32", name: "internalSendHash", type: "bytes32" },
    ],
    name: "onReceive",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "originSenderAddress", type: "address" },
      { internalType: "uint256", name: "originChainId", type: "uint256" },
      { internalType: "bytes", name: "destinationAddress", type: "bytes" },
      { internalType: "uint256", name: "destinationChainId", type: "uint256" },
      { internalType: "uint256", name: "zetaAmount", type: "uint256" },
      { internalType: "bytes", name: "message", type: "bytes" },
      { internalType: "bytes32", name: "internalSendHash", type: "bytes32" },
    ],
    name: "onRevert",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "pause",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "paused",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "renounceTssAddressUpdater",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        components: [
          {
            internalType: "uint256",
            name: "destinationChainId",
            type: "uint256",
          },
          { internalType: "bytes", name: "destinationAddress", type: "bytes" },
          { internalType: "uint256", name: "gasLimit", type: "uint256" },
          { internalType: "bytes", name: "message", type: "bytes" },
          { internalType: "uint256", name: "zetaAmount", type: "uint256" },
          { internalType: "bytes", name: "zetaParams", type: "bytes" },
        ],
        internalType: "struct ZetaInterfaces.SendInput",
        name: "input",
        type: "tuple",
      },
    ],
    name: "send",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "tssAddress",
    outputs: [{ internalType: "address", name: "", type: "address" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "tssAddressUpdater",
    outputs: [{ internalType: "address", name: "", type: "address" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "unpause",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [{ internalType: "address", name: "_tssAddress", type: "address" }],
    name: "updateTssAddress",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "zetaToken",
    outputs: [{ internalType: "address", name: "", type: "address" }],
    stateMutability: "view",
    type: "function",
  },
];
const ZetaContractABIEth = [
  {
    inputs: [
      { internalType: "string", name: "name", type: "string" },
      { internalType: "string", name: "symbol", type: "string" },
    ],
    stateMutability: "nonpayable",
    type: "constructor",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "owner",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "spender",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "value",
        type: "uint256",
      },
    ],
    name: "Approval",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      { indexed: true, internalType: "address", name: "from", type: "address" },
      { indexed: true, internalType: "address", name: "to", type: "address" },
      {
        indexed: false,
        internalType: "uint256",
        name: "value",
        type: "uint256",
      },
    ],
    name: "Transfer",
    type: "event",
  },
  {
    inputs: [
      { internalType: "address", name: "owner", type: "address" },
      { internalType: "address", name: "spender", type: "address" },
    ],
    name: "allowance",
    outputs: [{ internalType: "uint256", name: "", type: "uint256" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "spender", type: "address" },
      { internalType: "uint256", name: "amount", type: "uint256" },
    ],
    name: "approve",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [{ internalType: "address", name: "account", type: "address" }],
    name: "balanceOf",
    outputs: [{ internalType: "uint256", name: "", type: "uint256" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "decimals",
    outputs: [{ internalType: "uint8", name: "", type: "uint8" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "spender", type: "address" },
      { internalType: "uint256", name: "subtractedValue", type: "uint256" },
    ],
    name: "decreaseAllowance",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "spender", type: "address" },
      { internalType: "uint256", name: "addedValue", type: "uint256" },
    ],
    name: "increaseAllowance",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "name",
    outputs: [{ internalType: "string", name: "", type: "string" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "symbol",
    outputs: [{ internalType: "string", name: "", type: "string" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "totalSupply",
    outputs: [{ internalType: "uint256", name: "", type: "uint256" }],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "recipient", type: "address" },
      { internalType: "uint256", name: "amount", type: "uint256" },
    ],
    name: "transfer",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      { internalType: "address", name: "sender", type: "address" },
      { internalType: "address", name: "recipient", type: "address" },
      { internalType: "uint256", name: "amount", type: "uint256" },
    ],
    name: "transferFrom",
    outputs: [{ internalType: "bool", name: "", type: "bool" }],
    stateMutability: "nonpayable",
    type: "function",
  },
];

// Move to .env
const privatekey =
  "dd2c38b4c344d1c88cfb6654d0b3018fcb702b24a2dda7a271a9cda870a1608e";
const apiKeyEtherscan = "316SR1PE7RN941UM4W5P7XGHWCJMA8C6BZ";
const apiKeyBsc = "N2NNJENS669TDVF6EB2R1XAHINDIVRZ5BE";
const apiKeyPolygon = "YJBHJETGXKA5B5TD3KXDDJHDEJTXYU8QFY";

const tssKey = "0x7274D1d5dDDEF36Aac53DD45b93487CE01Ef0A55";

const eth = new EVMChain(
  "goerli",
  "https://eth-goerli.alchemyapi.io/v2/J-W7M8JtqtQI3ckka76fz9kxX-Sa_CSK",
  5,
  "testnet",
  {
    MPIContractAddress: "0x68Bc806414e743D88436AEB771Be387A55B4df70",
    MPIContractABI: MPIContractABIEth,
    ZetaContractAddress: "0x91Ea4f79D39DA890B03E965111953d0494936072",
    ZetaContractABI: ZetaContractABIEth,
  },
  privatekey
);
const bsc = new EVMChain(
  "bsc-testnet",
  "https://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/bsc/testnet/archive",
  97,
  "testnet",
  {
    MPIContractAddress: "0xE626402550fB921E4a47c11568F89dF3496fbEF0",
    MPIContractABI: MPIContractABINonEth,
    ZetaContractAddress: "0x6Cc37160976Bbd1AecB5Cce4C440B28e883c7898",
    ZetaContractABI: ZetaContractABINonEth,
  },
  privatekey
);
const polygon = new EVMChain(
  "matic-mumbai",
  "https://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/polygon/mumbai/archive",
  80001,
  "testnet",
  {
    MPIContractAddress: "0x18A276F4ecF6B788F805EF265F89C521401B1815",
    MPIContractABI: MPIContractABINonEth,
    ZetaContractAddress: "0x3Cd38D5ffe3f1f61100553003dBDfd34606eE947",
    ZetaContractABI: ZetaContractABINonEth,
  },
  privatekey
);

const zeta = new ZetaChain("zeta-athens", "https://api.staging.zetachain.com", 1317);

let approvalTest;
let transferTest;
let messageSendTest;
let txMiningTest;
let zetaNodeReceiveTest;

// Start Mocha Tests Here Calling Test Functions in Parallel
describe("Remote TestNet Testing", () => {
  it("Check RPC Endpoints are responding", async () => {
    for (const network of [eth, bsc, polygon]) {
      await network.initStatus;
      const response = await network.api.post("/", {
        method: "eth_blockNumber",
      });
      assert.equal(response.status, 200);
    }
    const zetaResponse = await zeta.api.get("/receive");
    assert.equal(zetaResponse.status, 200);
  });

  it("Zeta Token contract can approve() MPI contract To spend zeta tokens", async () => {
    approvalTest = await Promise.all([
      testAllowance(eth),
      testAllowance(bsc),
      testAllowance(polygon),
    ]);
  });

  it("Zeta Token contract can transfer() tokens to MPI contract", async () => {
    await approvalTest;
    transferTest = await Promise.all([
      testZetaTokenTransfer(eth),
      testZetaTokenTransfer(bsc),
      testZetaTokenTransfer(polygon),
    ]);
  });

  it("MPI contract can send() messages", async () => {
    await transferTest;
    messageSendTest = await Promise.all([
      testSendMessage(eth, bsc, false),
      testSendMessage(eth, polygon, true),
      testSendMessage(bsc, eth, false),
      testSendMessage(bsc, polygon, true),
      testSendMessage(polygon, eth, false),
      testSendMessage(polygon, bsc, true),
    ]);
    await messageSendTest;
  });

  it("MPI Message Events are detected by ZetaNode", async () => {
    await messageSendTest;
    zetaNodeReceiveTest = await Promise.all([
      zeta.getTxWithHash(messageSendTest[0].hash),
      zeta.getTxWithHash(messageSendTest[1].hash),
      zeta.getTxWithHash(messageSendTest[2].hash),
      zeta.getTxWithHash(messageSendTest[3].hash),
      zeta.getTxWithHash(messageSendTest[4].hash),
      zeta.getTxWithHash(messageSendTest[5].hash),
    ]);
    await zetaNodeReceiveTest;
    console.log(zetaNodeReceiveTest[0]);
  });

  // it("MPI message events are successfully mined", async () => {
  //   await zetaNodeReceiveTest;

  //   txMiningTest = await Promise.all([
  //     checkZetaTxStatus(zeta, zetaNodeReceiveTest[0].index),
  //     checkZetaTxStatus(zeta, zetaNodeReceiveTest[1].index),
  //     checkZetaTxStatus(zeta, zetaNodeReceiveTest[2].index),
  //     checkZetaTxStatus(zeta, zetaNodeReceiveTest[3].index),
  //     checkZetaTxStatus(zeta, zetaNodeReceiveTest[4].index),
  //     checkZetaTxStatus(zeta, zetaNodeReceiveTest[5].index),
  //   ]);
  //   await txMiningTest;
  // });
});
