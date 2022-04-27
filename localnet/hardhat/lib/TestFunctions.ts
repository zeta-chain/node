import { assert, expect } from "chai";
import { EVMChain, ZetaChain, wait } from "./Blockchain";
import { Axios, AxiosInstance } from "axios";
const axios = require("axios").default;

const defaultDestinationAddress = "0x04dA1034E7d84c004092671bBcEb6B1c8DCda7AE";
const defaultOverrideOptions = {
  gasLimit: 2500000,
  gasPrice: 20000000000,
};

export async function testAllowance(
  network: EVMChain,
  allowanceAmount = "1000000000000000" // 0.001 ZETA
) {
  await network.initStatus;

  const tx = await network.ZetaContract.approve(
    network.MPIContractAddress,
    allowanceAmount,
    defaultOverrideOptions
  );
  await tx.wait(1);

  // Check that the allowance was set correctly
  const allowance = await network.ZetaContract.allowance(
    network.accounts[0],
    network.MPIContractAddress,
    defaultOverrideOptions
  );

  if (allowance.toString() !== allowanceAmount) {
    console.error("Allowance Amounts do not match");
    console.error("Allowance Amount: " + allowance.toString());
    console.error(`${network.name} allowance() Tx Hash: ${tx.hash}`);
  }

  assert.equal(
    allowance.toString() === allowanceAmount,
    true,
    network.name + "Allowance Amounts do not match"
  );
}

export async function testZetaTokenTransfer(
  network: EVMChain,
  transferAmount = "10000000" // 0.00000000001
) {
  await network.initStatus;

  const preTxBalance = await network.ZetaContract.balanceOf(
    network.MPIContractAddress,
    defaultOverrideOptions
  );

  const tx = await network.ZetaContract.transfer(
    network.MPIContractAddress,
    transferAmount,
    defaultOverrideOptions
  );

  await tx.wait(2);
  const locked = await network.ZetaContract.balanceOf(
    network.MPIContractAddress,
    defaultOverrideOptions
  );
  console.info(
    `Amount of ZETA locked in ${
      network.name
    } MPI contract: ${locked.toString()}`
  );

  const correctTotal = preTxBalance.add(transferAmount);

  console.debug("PreTxBalance: " + preTxBalance.toString());
  console.debug("TransferAmount: " + transferAmount);
  console.debug("Locked: " + locked.toString());
  console.debug("Correct Total: " + correctTotal.toString());

  assert.equal(
    locked.toString() === correctTotal.toString(),
    true,
    network.name + " Locked Amount Value Not As Expected"
  );
}

export async function testSendMessage(
  sourceNetwork: EVMChain,
  destinationNetwork: EVMChain,
  incrementNone: boolean = false,
  zetaAmount = "10000000", // 0.00000000001
  gasLimit = "300000",
  message: any = [],
  zetaParams:any = []
) {
  await sourceNetwork.initStatus;
  await destinationNetwork.initStatus;

  const input = {
    destinationChainId: destinationNetwork.chainId,
    // destinationAddress: destinationNetwork.MPIContractAddress,
    destinationAddress: defaultDestinationAddress,
    gasLimit: gasLimit,
    message: message,
    zetaAmount: zetaAmount,
    zetaParams: zetaParams,
  };

  const overrideOptions = defaultOverrideOptions;

  if (incrementNone) {
    overrideOptions.nonce = (
      (await sourceNetwork.wallet.getTransactionCount()) + 1
    ).toString();
  }
  // console.log(overrideOptions);
  const tx = await sourceNetwork.MPIContract.send(
    input,
    defaultOverrideOptions
  );

  await tx.wait();
  console.info(
    `Outbound Hash for MPI Message From ${sourceNetwork.name} to ${destinationNetwork.name} ${tx.hash}`
  );

  console.debug(tx);
  return await tx;
}

export async function getTxIndexFromSourceTxHash(network, hash) {
  const tx = await network.getTxWithHash(hash);
  console.info(
    `Found Zeta TX Index Hash From ${tx.data.Send.senderChain} to ${tx.data.Send.receiverChain}: ${tx.index}`
  );
  return tx.index;
}

export async function checkZetaTxStatus(network: ZetaChain, zetaIndexHash: string) {
  console.info(`Checking TX Status For Zeta Index Hash: ${zetaIndexHash}`);

  let tx = await network.getEvent(zetaIndexHash);

  if (tx.data.Send.status !== "Finalized") {
    await wait(15000); // If the TX status isn't finalized wait 15 seconds and try again
    tx = await network.getEvent(zetaIndexHash);
  }
  console.debug(
    `Zeta TX Status From ${tx.data.Send.senderChain} to ${tx.data.Send.receiverChain}: ${tx.data.Send.status}`
  );
  assert.equal(
    tx.data.Send.status,
    "Finalized",
    `Zeta transaction status not mined for index hash ${zetaIndexHash} from ${tx.data.Send.senderChain} to ${tx.data.Send.receiverChain}`
  );
}


// // TODO - Store Explorer Endpoint, API Key, etc as part of the network
// export async function checkForMpiTransactions(
//   network: EVMChain,
//   tssAddress: string,
//   apiEndpoint: string,
//   apiKey: string,
//   blocksToCheck: number = 85
// ) {
//   let success = false;
//   const endpoint = apiEndpoint;
//   const api = await axios.create({
//     baseURL: `${endpoint}/`,
//     timeout: 15000,
//     jsonrpc: "2.0",
//     headers: {
//       Accept: "application/json",
//       "Content-Type": "application/json",
//     },
//   });

//   const response = await network.api.post("/", {
//     method: "eth_blockNumber",
//   });

//   const startBlockNumber = parseInt(response.data.result, 16) - blocksToCheck;

//   const r = await api.get(
//     `/api?module=account&action=txlist&address=${network.MPIContractAddress}&sort=desc&apikey=${apiKey}&startblock=${startBlockNumber}`
//   );

//   console.debug(r.data.result[0]);

//   for (const tx in r.data.result) {
//     if (
//       (r.data.result[tx].from = tssAddress) &&
//       (r.data.result[tx].isError = "0")
//     ) {
//       console.info(
//         `Successful Message Received on ${network.name} -- TSS Address: ${tssAddress} -- Tx Hash ${r.data.result[tx].hash}`
//       );
//       success = true;
//       break;
//     }
//   }
//   assert.equal(
//     success,
//     true,
//     `No Successful Messages From TSS Address Detected To The MPI Contract Have Been Detected In The Last ${blocksToCheck} Blocks`
//   );

//   return success;
// }
