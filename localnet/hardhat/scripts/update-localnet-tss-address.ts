import { assert, expect } from "chai";
import { EVMChain, ZetaChain, wait } from "../lib/Blockchain";

// The following two lines disable debug and info logging - comment them out to enable
console.debug = function () {}; // Disables Debug Level Logging
console.info = function () {}; // Disables Info Level Logging

// Setup each chain
const eth = new EVMChain("ethLocalNet", "http://localhost:8100", 5);
const bsc = new EVMChain("bscLocalNet", "http://localhost:8120", 97);
const polygon = new EVMChain("polygonLocalNet", "http://localhost:8140", 80001);
// const zeta = new ZetaChain("zetaLocalNet", "http://localhost:1317", 1317);

const tssAddress = process.argv[2];

async function updateTSSAddress(network, tssAddress) {
  await network.initStatus;
  const tx = await network.MPIContract.updateTssAddress(tssAddress);
  const tx2 = await network.ZetaContract.updateTSSAndMPIAddresses(
    tssAddress,
    network.MPIContractAddress
  );

  console.debug(await tx.wait());
  console.debug(await tx2.wait());

  const tx3 = await network.MPIContract.tssAddress();
  assert.equal(tx3.toString(), tssAddress);
  console.log(
    `TSS Address for ${network.name} has been updated to: ${tx3.toString()}`
  );
}

async function updateEthTSSAddress(network, tssAddress) {
  await network.initStatus;
  const tx = await network.MPIContract.updateTssAddress(tssAddress);

  console.debug(await tx.wait());

  const tx2 = await network.MPIContract.tssAddress();
  assert.equal(tx2.toString(), tssAddress);
  console.log(
    `TSS Address for ${network.name} has been updated to: ${tx2.toString()}`
  );
}

updateTSSAddress(bsc, tssAddress);
updateTSSAddress(polygon, tssAddress);
updateEthTSSAddress(eth, tssAddress);
