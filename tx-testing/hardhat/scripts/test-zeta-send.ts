import { AbiCoder } from "ethers/lib/utils";
import { ethers } from "hardhat";
import { getAddress, getChainId } from "../lib/networks";
import { ZetaMPIEth__factory as ZetaMPIEthFactory } from "../typechain";

const encoder = new AbiCoder();

async function main() {
  const factory = (await ethers.getContractFactory(
    "ZetaMPIEth"
  )) as ZetaMPIEthFactory;

  const contract = factory.attach(getAddress("mpi"));

  console.log("Sending");
  await (
    await contract.send({
      destinationChainId: getChainId("bsc-testnet"),
      destinationAddress: encoder.encode(
        ["address"],
        ["0x09b80BEcBe709Dd354b1363727514309d1Ac3C7b"]
      ),
      gasLimit: 1_000_000,
      message: encoder.encode(
        ["address"],
        ["0x09b80BEcBe709Dd354b1363727514309d1Ac3C7b"]
      ),
      zetaAmount: 0,
      zetaParams: [],
    })
  ).wait();
  console.log("Sent");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
