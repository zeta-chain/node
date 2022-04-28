import { ethers } from "hardhat";
import {
  ZetaEth,
  ZetaEth__factory as ZetaEthFactory,
  ZetaMPIBase,
  ZetaMPIBase__factory as ZetaMPIBaseFactory,
  ZetaMPIEth,
  ZetaMPIEth__factory as ZetaMPIEthFactory,
  ZetaMPINonEth,
  ZetaMPINonEth__factory as ZetaMPINonEthFactory,
  ZetaNonEth,
  ZetaNonEth__factory as ZetaNonEthFactory,
  ZetaReceiverMock,
  ZetaReceiverMock__factory as ZetaReceiverMockFactory,
} from "../typechain";

export const deployZetaMpiBase = async ({
  args,
}: {
  args: Parameters<ZetaMPIBaseFactory["deploy"]>;
}) => {
  const Factory = (await ethers.getContractFactory(
    "ZetaMPIBase"
  )) as ZetaMPIBaseFactory;

  const zetaMpiContract = (await Factory.deploy(...args)) as ZetaMPIBase;

  await zetaMpiContract.deployed();

  return zetaMpiContract;
};

export const deployZetaMpiEth = async ({
  args,
}: {
  args: Parameters<ZetaMPIEthFactory["deploy"]>;
}) => {
  const Factory = (await ethers.getContractFactory(
    "ZetaMPIEth"
  )) as ZetaMPIEthFactory;

  const zetaMpiContract = (await Factory.deploy(...args)) as ZetaMPIEth;

  await zetaMpiContract.deployed();

  return zetaMpiContract;
};

export const deployZetaMpiNonEth = async ({
  args,
}: {
  args: Parameters<ZetaMPINonEthFactory["deploy"]>;
}) => {
  const Factory = (await ethers.getContractFactory(
    "ZetaMPINonEth"
  )) as ZetaMPINonEthFactory;

  const zetaMpiContract = (await Factory.deploy(...args)) as ZetaMPINonEth;

  await zetaMpiContract.deployed();

  return zetaMpiContract;
};

export const deployZetaReceiverMock = async () => {
  const Factory = (await ethers.getContractFactory(
    "ZetaReceiverMock"
  )) as ZetaReceiverMockFactory;

  const zetaReceiverMock = (await Factory.deploy()) as ZetaReceiverMock;

  await zetaReceiverMock.deployed();

  return zetaReceiverMock;
};

export const deployZetaEth = async ({
  args,
}: {
  args: Parameters<ZetaEthFactory["deploy"]>;
}) => {
  const Factory = (await ethers.getContractFactory(
    "ZetaEth"
  )) as ZetaEthFactory;

  const zetaEthContract = (await Factory.deploy(...args)) as ZetaEth;

  await zetaEthContract.deployed();

  return zetaEthContract;
};

export const deployZetaNonEth = async ({
  args,
}: {
  args: Parameters<ZetaNonEthFactory["deploy"]>;
}) => {
  const Factory = (await ethers.getContractFactory(
    "ZetaNonEth"
  )) as ZetaNonEthFactory;

  const zetaEthContract = (await Factory.deploy(...args)) as ZetaNonEth;

  await zetaEthContract.deployed();

  return zetaEthContract;
};
