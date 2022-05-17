import { Contract } from "ethers";
import { network } from "hardhat";
import {
  deployZetaMpiEth,
  deployZetaMpiNonEth,
} from "../lib/contracts.helpers";
import { getAddress, isNetworkName, saveAddress } from "../lib/networks";

async function main() {
  if (!isNetworkName(network.name)) {
    throw new Error(`network.name: ${network.name} isn't supported.`);
  }

  let contract: Contract;
  console.log(`Deploying ZetaMPI to ${network.name}`);

  if (network.name === "goerli" || network.name === "eth-mainnet") {
    console.log([
      getAddress("zetaToken"),
      getAddress("tss"),
      getAddress("tssUpdater"),
    ]);
    contract = await deployZetaMpiEth({
      args: [
        getAddress("zetaToken"),
        getAddress("tss"),
        getAddress("tssUpdater"),
      ],
    });
  } else {
    contract = await deployZetaMpiNonEth({
      args: [
        getAddress("zetaToken"),
        getAddress("tss"),
        getAddress("tssUpdater"),
      ],
    });
  }

  saveAddress("mpi", contract.address);
  console.log("Deployed ZetaMPI. Address:", contract.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
