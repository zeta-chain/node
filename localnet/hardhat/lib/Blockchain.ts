import * as fs from "fs";
import { Contract, ethers, Signer } from "ethers";
import { PrivateKey } from "web3/eth/accounts";
import { getAddress, isNetworkName, saveAddress } from "./networks";
import { Axios, AxiosInstance } from "axios";
const axios = require("axios").default;

// console.debug = function () {}; // Disables Debug Level Logging
// console.info = function () {}; // Disables Info Level Logging
export class Blockchain {
  name: string;
  chainId: number;
  rpcEndpoint: string;
  initStatus: Promise<void>;
  p: any; // TODO - ethers.provider
  accounts: any;
  api: AxiosInstance;
  type: "local" | "testnet" | "mainnet" = "local";
  wallet: ethers.Wallet;
  privateKey: string;

  constructor(
    name: string,
    rpcEndpoint: string,
    chainId: number,
    type: "local" | "testnet" | "mainnet" = "local",
    privateKey: string = null
  ) {
    this.name = name;
    this.chainId = chainId;
    this.rpcEndpoint = rpcEndpoint;
    this.type = type;
    this.privateKey = privateKey;
  }

  async promiseTest(ms: number) {
    const promise = new Promise((resolve) => setTimeout(resolve, ms));
    await promise;
  }
}

export class EVMChain extends Blockchain {
  MPIContractAddress: string;
  MPIContractABI: string[];
  MPIContract: Contract; // How to avoid errors on all sub functions and values

  ZetaContractAddress: string;
  ZetaContractABI: string[];
  ZetaContract: Contract;

  signer: Signer;
  constructor(
    name: string,
    rpcEndpoint: string,
    chainId: number,
    type: "local" | "testnet" | "mainnet" = "local",
    contractArgs: {
      MPIContractAddress: string;
      MPIContractABI: any[];
      ZetaContractAddress: string;
      ZetaContractABI: any[];
    } = null,
    privateKey: string = null
  ) {
    super(name, rpcEndpoint, chainId, type, privateKey);

    if (contractArgs) {
      this.MPIContractAddress = contractArgs.MPIContractAddress;
      this.MPIContractABI = contractArgs.MPIContractABI;
      this.ZetaContractAddress = contractArgs.ZetaContractAddress;
      this.ZetaContractABI = contractArgs.ZetaContractABI;
    }
    this.initStatus = this.init(privateKey);
  }

  async init(walletPrivateKey = null) {
    this.p = new ethers.providers.JsonRpcProvider(this.rpcEndpoint);

    if (this.type === "local") {
      this.accounts = await this.p.listAccounts();
      this.signer = this.p.getSigner(this.accounts[0]);
      this.MPIContractAddress = await this.getLocalNetContractAddress(
        "zetaMPI"
      );
      this.MPIContractABI = await this.loadMPIContractABI();
      this.ZetaContractAddress = await this.getLocalNetContractAddress("zeta");
      this.ZetaContractABI = await this.loadTokenContractABI();
    } else {
      this.wallet = new ethers.Wallet(walletPrivateKey, this.p);
      this.accounts = [this.wallet.address];
      this.signer = this.wallet;
    }

    this.ZetaContract = new ethers.Contract(
      this.ZetaContractAddress,
      this.ZetaContractABI,
      this.signer
    );
    this.MPIContract = new ethers.Contract(
      this.MPIContractAddress,
      this.MPIContractABI,
      this.signer
    );

    this.api = await axios.create({
      baseURL: `${this.rpcEndpoint}/`,
      timeout: 15000,
      jsonrpc: "2.0",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    });
  }

  async getAccountBalance(account = this.accounts) {
    await this.initStatus;
    const balance = await this.p.getBalance(account);
    console.info(`balance of ${account}: ${balance}`);
    return balance;
  }

  async getLocalNetContractAddress(contractName: string) {
    try {
      const data = fs.readFileSync(
        `localnet-addresses/${this.name}-${contractName}-address`,
        "utf8"
      );
      return data;
    } catch (err) {
      console.error(err);
    }
  }

  async loadMPIContractABI() {
    let contractName;
    if (this.name === "ethLocalNet") {
      contractName = "evm/ZetaMPI.eth.sol/ZetaMPIEth.json";
    } else {
      contractName = "evm/ZetaMPI.non-eth.sol/ZetaMPINonEth.json";
    }
    return require(`../artifacts/contracts/${contractName}`).abi;
  }

  async loadTokenContractABI() {
    let contractName: string;
    if (this.name === "ethLocalNet") {
      contractName = "evm/ZetaEth.sol/ZetaEth.json";
    } else {
      contractName = "evm/ZetaNonEth.sol/ZetaNonEth.json";
    }
    return require(`../artifacts/contracts/${contractName}`).abi;
  }

  async getMPIEvents(
    eventName = "ZetaMessageReceiveEvent",
    filterKeyPairs = null
  ) {
    await this.initStatus;

    const sentFilter = await this.MPIContract.filters[eventName]();
    const sent = await this.MPIContract.queryFilter(sentFilter);
    console.debug(
      `Recent '${eventName}' Events On ${this.name} Network: ${sent}`
    );
  }
}

export async function wait(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

export class ZetaChain extends Blockchain {
  api: AxiosInstance;

  constructor(name: string, rpcEndpoint: string, chainId: number) {
    super(name, rpcEndpoint, chainId);
    this.initStatus = this.init();
  }

  async init() {
    this.api = await axios.create({
      baseURL: `${this.rpcEndpoint}/Meta-Protocol/zetacore/zetacore/`,
      timeout: 10000,
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    });
  }

  //   Working but timeout function could be improved
  async getTxWithHash(txHash: string, timeout: number = 18) {
    // TODO - Better Timeout handling
    let i = 0;
    let response;
    console.debug(
      "Checking Zetachain For Transaction With Source Hash: " + txHash
    );
    do {
      try {
        if (i >= timeout) {
          throw new Error(
            "TX not found received within " + i * 10 + " seconds"
          );
        }
        response = await this.api.get(`inTxRich/${txHash}`);
        // console.debug(response);

        if (response.data.tx != null) {
          console.debug(
            `Found Transaction ${txHash} from ${response.data.tx.senderChain}`
          );
          return response.data.tx;
        }
        i++;
        await wait(10000); // Wait 10 seconds between checks
      } catch (err) {
        console.error(err);
        throw err;
      }
    } while (response.data.tx == null);
  }

  async getEvent(sendHash: string) {
    const result = await this.api.get(`send/${sendHash}`);
    return result;
  }
}
