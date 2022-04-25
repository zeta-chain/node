import { assert, expect } from "chai";
import { EVMChain, ZetaChain, wait } from "../lib/Blockchain";
import { describe } from "mocha";
import {
  testAllowance,
  testZetaTokenTransfer,
  testSendMessage,
  checkZetaTxStatus,
} from "../lib/TestFunctions";

// The following two lines disable debug and info logging - comment them out to enable
console.debug = function () {}; // Disables Debug Level Logging
// console.info = function () {}; // Disables Info Level Logging

// Setup each chain
const eth = new EVMChain("ethLocalNet", "http://localhost:8100", 5);
const bsc = new EVMChain("bscLocalNet", "http://localhost:8120", 97);
const polygon = new EVMChain("polygonLocalNet", "http://localhost:8140", 80001);
const zeta = new ZetaChain("zetaLocalNet", "http://localhost:1317", 1317);

let approvalTest;
let transferTest;
let messageSendTest;
let txMiningTest;
let zetaNodeReceiveTest;

// Start Mocha Tests Here Calling Test Functions in Parallel
describe("LocalNet Testing", () => {
  it("Check RPC Endpoints are responding", async () => {
    for (const network of [eth, bsc, polygon]) {
      await network.initStatus;
      const response = await network.api.post("/", {
        method: "eth_blockNumber",
      });
      assert.equal(response.status, 200);
    }
    const zetaResponse = await zeta.api.get("/receive");
    assert.equal(zetaResponse.status, 200);
  });

  it("Zeta Token contract can approve() MPI contract To spend zeta tokens", async () => {
    approvalTest = await Promise.all([
      testAllowance(eth),
      testAllowance(bsc),
      testAllowance(polygon),
    ]);
  });

  it("Zeta Token contract can transfer() tokens to MPI contract", async () => {
    await approvalTest;
    transferTest = await Promise.all([
      testZetaTokenTransfer(eth),
      testZetaTokenTransfer(bsc),
      testZetaTokenTransfer(polygon),
    ]);
  });

  it("MPI contract can send() messages", async () => {
    await transferTest;
    messageSendTest = await Promise.all([
      testSendMessage(eth, bsc, false),
      testSendMessage(eth, polygon, true),
      testSendMessage(bsc, eth, false),
      testSendMessage(bsc, polygon, true),
      testSendMessage(polygon, eth, false),
      testSendMessage(polygon, bsc, true),
    ]);
    await messageSendTest;
  });

  it("MPI Message Events are detected by ZetaNode", async () => {
    await messageSendTest;
    zetaNodeReceiveTest = await Promise.all([
      zeta.getTxWithHash(messageSendTest[0].hash),
      zeta.getTxWithHash(messageSendTest[1].hash),
      zeta.getTxWithHash(messageSendTest[2].hash),
      zeta.getTxWithHash(messageSendTest[3].hash),
      zeta.getTxWithHash(messageSendTest[4].hash),
      zeta.getTxWithHash(messageSendTest[5].hash),
    ]);
    await zetaNodeReceiveTest;
    console.log(zetaNodeReceiveTest[0]);
  });

  it("MPI message events are successfully mined", async () => {
    await zetaNodeReceiveTest;

    txMiningTest = await Promise.all([
      checkZetaTxStatus(zeta, zetaNodeReceiveTest[0].index),
      checkZetaTxStatus(zeta, zetaNodeReceiveTest[1].index),
      checkZetaTxStatus(zeta, zetaNodeReceiveTest[2].index),
      checkZetaTxStatus(zeta, zetaNodeReceiveTest[3].index),
      checkZetaTxStatus(zeta, zetaNodeReceiveTest[4].index),
      checkZetaTxStatus(zeta, zetaNodeReceiveTest[5].index),
    ]);
    await txMiningTest;
  });
});
