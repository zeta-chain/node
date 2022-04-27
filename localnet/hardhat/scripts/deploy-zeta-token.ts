import { Contract } from "ethers";
import { ethers, network } from "hardhat";
import { deployZetaEth, deployZetaNonEth } from "../lib/contracts.helpers";
import { getAddress, isNetworkName, saveAddress } from "../lib/networks";

async function main() {
  if (!isNetworkName(network.name)) {
    throw new Error(`network.name: ${network.name} isn't supported.`);
  }

  let contract: Contract;
  console.log(`Deploying Zeta Token to ${network.name}`);

  if (network.name === "goerli" || network.name === "eth-mainnet") {
    contract = await deployZetaEth({
      args: [ethers.utils.parseEther("10000000")],
    });
  } else {
    contract = await deployZetaNonEth({
      args: [0, getAddress("tss"), getAddress("tssUpdater")],
    });
  }

  saveAddress("zetaToken", contract.address);
  console.log("Deployed Zeta to:", contract.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
