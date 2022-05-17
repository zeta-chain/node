// The network name is passed in as an argument
// import { getAddress, isNetworkName, saveAddress } from "./lib/networks";

// const deploymentNetwork = process.argv[2];

// let _TSSAddress = process.argv[3]

let _TSSAddress = "0xC897825288aB5EC3CE86437042004dCB0a190962"


// process.env.HARDHAT_NETWORK = deploymentNetwork;
const hre = require("hardhat");
const networkName = hre.network.name;

const defaultOverrideOptions = {
  gasLimit: 2500000,
  gasPrice: 20000000000,
};

async function sendGas() {
  // If this script is run directly using `node` you may want to call compile to make sure everything is compiled
  // await hre.run('compile'); //

  const accounts = await hre.ethers.getSigners();

  console.log(accounts[0].address);

  // Send gas to TSS Address
  const sendGasTx = {
    from: accounts[0].address,
    to: _TSSAddress,
    value: hre.ethers.utils.parseEther("1.0"),
  };
  await accounts[0].sendTransaction(sendGasTx);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
sendGas()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
