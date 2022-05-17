// The network name is passed in as an argument
// import { getAddress, isNetworkName, saveAddress } from "./lib/networks";

const deploymentNetwork = process.argv[2];

let mainnet;
if (process.argv[3] === "true") {
  mainnet = true;
} else {
  mainnet = false;
}

process.env.HARDHAT_NETWORK = deploymentNetwork;
const hre = require("hardhat");
const assert = require("assert");
const networkName = hre.network.name;

const defaultOverrideOptions = {
  gasLimit: 2500000,
  gasPrice: 20000000000,
};

async function deployZeta() {
  // If this script is run directly using `node` you may want to call compile to make sure everything is compiled
  // await hre.run('compile'); //

  const accounts = await hre.ethers.getSigners();
  let Zeta;
  let ZetaMPI;
  console.log(accounts[0].address);
  console.log("Deploying Zeta Contract to the " + networkName + " network");

  const _TSSAddress = "0x80AFF18eE2862980C859aF47Bc263A338c3d8399";
  const _TSSAddressUpdater = accounts[0].address;
  const zetaTotalSupply = "1000000000000000000000000000"; // 10^27

  if (networkName === "eth") {
    Zeta = await hre.ethers.getContractFactory("ZetaEth");
    ZetaMPI = await hre.ethers.getContractFactory("ZetaMPIEth");
  } else {
    Zeta = await hre.ethers.getContractFactory("ZetaNonEth");
    ZetaMPI = await hre.ethers.getContractFactory("ZetaMPINonEth");
  }

  const zeta = await Zeta.deploy(
    zetaTotalSupply,
    _TSSAddress,
    _TSSAddressUpdater
  );

  const zetaMPI = await ZetaMPI.deploy(
    zeta.address,
    _TSSAddress,
    _TSSAddressUpdater
  );

  await zeta.deployTransaction.wait();
  await zetaMPI.deployTransaction.wait();

  // approve MPI Contract to Spend Zeta Tokens
  zeta.approve(zetaMPI.address, "1000000000000000000000000"); // 10^24
  // zeta.transfer(zetaMPI.address, "1000000000000000000000000"); // 10^24

  // whitelist MPI contract in the token contract:
  console.log("Whitelisting MPI contract in ZETA contract...");
  const wl_tx = await zeta.updateTSSAndMPIAddresses(
    _TSSAddress,
    zetaMPI.address,
    defaultOverrideOptions
  );
  await wl_tx.wait();

  // Send gas to TSS Address
  const sendGasTx = {
    from: accounts[0].address,
    to: _TSSAddress,
    value: hre.ethers.utils.parseEther("1.0"),
  };
  await accounts[0].sendTransaction(sendGasTx);

  console.log("Zeta Contract deployed to:", zeta.address);
  console.log("Zeta MPI Contract deployed to:", zetaMPI.address);

  const totalSupply = await zeta.totalSupply();

  assert.equal(
    hre.ethers.utils.formatEther(totalSupply),
    Number(zetaTotalSupply)
  );
  console.log("Zeta total supply", hre.ethers.utils.formatEther(totalSupply));

  // Write values for the Zeta Contracts
  // This should be doable within Node

  const fs = require("fs");

  const zetaFilename = networkName + "-zeta-address";
  try {
    fs.writeFileSync("./localnet-addresses/" + zetaFilename, zeta.address);
    // file written successfully
  } catch (err) {
    console.error(err);
  }
  const zetaMPIFilename = networkName + "-zetaMPI-address";
  try {
    fs.writeFileSync(
      "./localnet-addresses/" + zetaMPIFilename,
      zetaMPI.address
    );
    // file written successfully
  } catch (err) {
    console.error(err);
  }

  // saveAddress("zetaToken", zeta.address);
  // saveAddress("mpi", zetaMPI.address);

  // Mint Tokens to Accounts[0] if MainNet = False
  if (!mainnet) {
    console.log("Minting Zeta Tokens to account[0]");
    zeta.mint(
      accounts[0].address,
      "1000000",
      "0x7588a7ce78e90f8c192b04f526819880aa84dcd09a93240d300fe865f3c237f3",
      { from: accounts[0].address }
    );

    console.log("Approving Zeta MPI Contract to account[0]");
    zeta.approve(zetaMPI.address, "1000000", { from: accounts[0].address });
  }
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
deployZeta()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
