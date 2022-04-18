const hardhat = require("hardhat");

const { ZETA_NETWORK } = process.env;

const addresses = require(`../../addresses/addresses.${ZETA_NETWORK}.json`);

module.exports = [
  addresses[hardhat.network.name].zetaToken,
  addresses[hardhat.network.name].tss,
  addresses[hardhat.network.name].tssUpdater,
];
