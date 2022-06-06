import "@nomiclabs/hardhat-waffle";
import "@nomiclabs/hardhat-etherscan";
import "@typechain/hardhat";

import { task } from "hardhat/config";
import * as dotenv from "dotenv";

dotenv.config();


const PRIVATE_KEYS =
  process.env.PRIVATE_KEY !== undefined ? [`0x${process.env.PRIVATE_KEY}`] : [];
// const PRIVATE_KEYS =
//   process.env.PRIVATE_KEY !== undefined
//     ? [`0x${process.env.TSS_PRIVATE_KEY}`]
//     : [];

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  solidity: "0.8.4",
  // defaultNetwork: "",
  networks: {
    hardhat: {},
    ethLocalNet: {
      // TODO - All these values should be dynamically loaded from elsewhere when the environment is setup
      url: "http://localhost:8100",
      gas: 2100000,
      gasPrice: 8000000000,
    },
    bscLocalNet: {
      url: "http://localhost:8120",
      gas: 30000000,
      gasPrice: 8000000000,
    },
    polygonLocalNet: {
      url: "http://localhost:8140",
      gas: 2100000,
      gasPrice: 8000000000,
    },
    "eth-mainnet": {
      url: "https://api.mycryptoapi.com/eth",
      accounts: PRIVATE_KEYS,
    },
    goerli: {
      url: "https://goerli.infura.io/v3/84842078b09946638c03157f83405213", // alternative 1
      // url: "https://rpc.goerli.mudit.blog", // alternative 2
      accounts: PRIVATE_KEYS,
      gas: 2100000,
      gasPrice: 8000000000,
    },
    "bsc-testnet": {
      url: `https://data-seed-prebsc-1-s1.binance.org:8545`,
      accounts: PRIVATE_KEYS,
      gas: 5000000,
      gasPrice: 80000000000,
    },
    "matic-mumbai": {
      // url: "https://rpc-mumbai.matic.today", // alternative 1
      url: "https://matic-mumbai.chainstacklabs.com", // alternative 2
      accounts: PRIVATE_KEYS,
      gas: 5000000,
      gasPrice: 80000000000,
    },
  },
  etherscan: {
    apiKey: {
      bscTestnet: process.env.BSCSCAN_API_KEY,
      goerli: process.env.ETHERSCAN_API_KEY,
      polygonMumbai: process.env.POYLGONSCAN_API_KEY,
    },
  },
};
