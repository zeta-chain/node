import { assert, expect } from "chai";
import { EVMChain, ZetaChain, wait } from "./Blockchain";

const defaultOverrideOptions = {
  gasLimit: 2500000,
  gasPrice: 20000000000,
};

export async function testAllowance(
  network: EVMChain,
  allowanceAmount = "2000000000000000000" // 2 ZETA
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
  transferAmount = "10000000000" // 0.01
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
  zeta: ZetaChain,
  type: "local" | "testnet" | "mockmpi" = "local",
  zetaAmount = "500000000000000000", // .5
  gasLimit = "300000",
  message = [],
  zetaParams = []
) {
  await sourceNetwork.initStatus;
  await destinationNetwork.initStatus;

  const input = {
    destinationChainId: destinationNetwork.chainId,
    destinationAddress: destinationNetwork.MPIContractAddress,
    gasLimit: gasLimit,
    message: message,
    zetaAmount: zetaAmount,
    zetaParams: zetaParams,
  };

  const tx = await sourceNetwork.MPIContract.send(
    input,
    defaultOverrideOptions
  );

  tx.wait();
  console.info(
    `Outbound Hash for MPI Message From ${sourceNetwork.name} to ${destinationNetwork.name} ${tx.hash}`
  );

  console.debug(tx);
  return tx;
}

export async function getTxIndexFromSourceTxHash(network, hash) {
  const tx = await network.getTxWithHash(hash);
  console.info(
    `Found Zeta TX Index Hash From ${tx.data.Send.senderChain} to ${tx.data.Send.receiverChain}: ${tx.index}`
  );
  return tx.index;
}

export async function checkZetaTxStatus(network, zetaIndexHash: string) {
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
