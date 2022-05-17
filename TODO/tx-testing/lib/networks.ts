import troyAddresses from "../addresses/addresses.troy.json";
import spartaAddresses from "../addresses/addresses.sparta.json";
import athensAddresses from "../addresses/addresses.athens.json";
import mainnetAddresses from "../addresses/addresses.mainnet.json";
import { network } from "hardhat";
import { readFileSync, writeFileSync } from "fs";
import { join } from "path";
import { execSync } from "child_process";
import dotenv from "dotenv";

type ZetaAddress = "mpi" | "tss" | "tssUpdater" | "zetaToken";
type NetworkAddresses = Record<ZetaAddress, string>;
const zetaAddresses = new Set<ZetaAddress>([
  "mpi",
  "tss",
  "tssUpdater",
  "zetaToken",
]);

export const isZetaAddress = (a: string | undefined): a is ZetaAddress =>
  zetaAddresses.has(a as ZetaAddress);

/**
 * @description Testnet
 */
type TestnetNetworkName = "goerli" | "bsc-testnet" | "matic-mumbai";
type ZetaTestnetNetworkName = "troy" | "sparta" | "athens";
type TestnetListItem = Record<TestnetNetworkName, NetworkAddresses>;
const isTestnetNetworkName = (
  networkName: string
): networkName is TestnetNetworkName =>
  networkName === "goerli" ||
  networkName === "bsc-testnet" ||
  networkName === "matic-mumbai";
const isZetaTestnet = (
  networkName: string | undefined
): networkName is ZetaTestnetNetworkName =>
  networkName === "troy" ||
  networkName === "sparta" ||
  networkName === "athens";

const testnetList: Record<ZetaTestnetNetworkName, TestnetListItem> = {
  athens: athensAddresses,
  sparta: spartaAddresses,
  troy: troyAddresses,
};

/**
 * @description Mainnet
 */
type MainnetNetworkName = "eth-mainnet";
type ZetaMainnetNetworkName = "mainnet";
type MainnetListItem = Record<MainnetNetworkName, NetworkAddresses>;
const isMainnetNetworkName = (
  networkName: string
): networkName is MainnetNetworkName => networkName === "eth-mainnet";
const isZetaMainnet = (
  networkName: string | undefined
): networkName is ZetaMainnetNetworkName => networkName === "mainnet";

const mainnetList: Record<ZetaMainnetNetworkName, MainnetListItem> = {
  mainnet: mainnetAddresses,
};

/**
 * @description Shared
 */

type NetworkName = TestnetNetworkName | MainnetNetworkName;

export const getChainId = (networkName: NetworkName) => {
  const chainIds: Record<NetworkName, number> = {
    goerli: 5,
    "bsc-testnet": 97,
    "matic-mumbai": 80001,
    "eth-mainnet": 1,
  };

  return chainIds[networkName];
};

export const isNetworkName = (str: string): str is NetworkName =>
  isTestnetNetworkName(str) || isMainnetNetworkName(str);

export const getScanVariable = (): string => {
  if (!isNetworkName(network.name)) throw new Error();
  dotenv.config();

  const v = {
    "bsc-testnet": process.env.BSCSCAN_API_KEY || "",
    "eth-mainnet": process.env.ETHERSCAN_API_KEY || "",
    goerli: process.env.ETHERSCAN_API_KEY || "",
    "matic-mumbai": process.env.POLYGONSCAN_API_KEY || "",
  };

  return v[network.name];
};

export const getAddress = (address: ZetaAddress): string => {
  const { ZETA_NETWORK } = process.env;
  const { name: networkName } = network;

  console.log(
    `Getting ${address} address from ${ZETA_NETWORK}: ${networkName}.`
  );

  if (isZetaTestnet(ZETA_NETWORK) && isTestnetNetworkName(networkName)) {
    return testnetList[ZETA_NETWORK][networkName][address];
  }

  if (isZetaMainnet(ZETA_NETWORK) && isMainnetNetworkName(networkName)) {
    return mainnetList[ZETA_NETWORK][networkName][address];
  }

  throw new Error(
    `Invalid ZETA_NETWORK + network combination ${ZETA_NETWORK} ${networkName}.`
  );
};

export const saveAddress = (addressName: ZetaAddress, newAddress: string) => {
  const { ZETA_NETWORK } = process.env;
  const { name: networkName } = network;

  console.log(
    `Updating ${addressName} address on ${ZETA_NETWORK}: ${networkName}.`
  );

  const dirname = join(
    __dirname,
    `../addresses/addresses.${ZETA_NETWORK}.json`
  );

  if (isZetaTestnet(ZETA_NETWORK) && isTestnetNetworkName(networkName)) {
    const originalAddresses: TestnetListItem = JSON.parse(
      readFileSync(dirname, "utf8")
    );

    const newAddresses = {
      ...originalAddresses,
    };
    newAddresses[networkName][addressName] = newAddress;

    writeFileSync(dirname, JSON.stringify(newAddresses, null, 2));

    console.log(`Updated, new address: ${newAddress}.`);

    return;
  }

  if (isZetaMainnet(ZETA_NETWORK) && isMainnetNetworkName(networkName)) {
    const originalAddresses: MainnetListItem = JSON.parse(
      readFileSync(dirname, "utf8")
    );

    const newAddresses = {
      ...originalAddresses,
    };
    newAddresses[networkName][addressName] = newAddress;

    writeFileSync(dirname, JSON.stringify(newAddresses, null, 2));

    console.log(`Updated, new address: ${newAddress}.`);

    return;
  }

  throw new Error(
    `Invalid ZETA_NETWORK + network combination ${ZETA_NETWORK} ${networkName}.`
  );
};

export const verifyContract = (addressName: ZetaAddress) => {
  const { ZETA_NETWORK } = process.env;
  const { name: networkName } = network;

  console.log(
    `Verifying ${addressName} address on ${ZETA_NETWORK}: ${networkName}.`
  );

  const command = `ZETA_NETWORK=${ZETA_NETWORK} SCAN_API_KEY=${getScanVariable()} npx hardhat verify --network ${networkName} --constructor-args lib/args/${addressName}.js ${getAddress(
    addressName
  )}`;

  execSync(command);
};
