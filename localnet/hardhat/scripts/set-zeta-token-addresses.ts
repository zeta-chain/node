import { ethers, network } from "hardhat";
import { getAddress, isNetworkName } from "../lib/networks";
import { ZetaNonEth__factory as ZetaNonEthFactory } from "../typechain";

async function main() {
  if (!isNetworkName(network.name)) {
    throw new Error(`network.name: ${network.name} isn't supported.`);
  }

  const factory = (await ethers.getContractFactory(
    "ZetaNonEth"
  )) as ZetaNonEthFactory;

  const contract = factory.attach(getAddress("zetaToken"));

  console.log("Updating");
  await (
    await contract.updateTSSAndMPIAddresses(
      getAddress("tss"),
      getAddress("mpi")
    )
  ).wait();
  console.log("Updated");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
