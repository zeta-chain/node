const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Zeta Token Contract", function () {
  it("supply, burn check", async function () {
    const TestOracle = await ethers.getContractFactory("TestSupplyOracle");
    const oracle = await TestOracle.deploy();
    await oracle.deployed();
    console.log("test oracle deployed at", oracle.address);

    const signers = await ethers.getSigners();
    const deployer = signers[0].address;

    const Zeta = await ethers.getContractFactory("Zeta");
    const supply = "10000";
    const zeta = await Zeta.deploy(
      supply,
      "Zeta",
      "ZETA",
      oracle.address,
      deployer,
      deployer,
      deployer
    );
    await zeta.deployed();
    console.log("contract deployed at address", zeta.address);
    console.log("deployer of contract", zeta.deployTransaction.from);

    expect(signers[0].address).to.equal(zeta.deployTransaction.from);
    expect(await zeta.totalSupply()).to.equal(ethers.utils.parseUnits(supply));

    oracle.setSupply(ethers.utils.parseUnits(supply));
    oracle.setLocked(ethers.utils.parseUnits(supply));

    // const tx = await zeta.burnFrom(signers[0].address, "1234");
    const tx = await zeta.burn(ethers.utils.parseEther("1234"));
    await tx.wait();
    expect(await zeta.totalSupply()).to.equal(ethers.utils.parseUnits("8766"));
    oracle.setSupply(ethers.utils.parseUnits("8766"));

    const MINT_TYPEHASH = await zeta.MINT_TYPEHASH();
    const BURN_TYPEHASH = await zeta.BURN_TYPEHASH();
    const mt = ethers.utils.keccak256(
      ethers.utils.toUtf8Bytes(
        "Mint(address mintee,uint256 value,uint256 nonce,uint256 deadline)"
      )
    );
    const bt = ethers.utils.keccak256(
      ethers.utils.toUtf8Bytes(
        "Burn(address burnee,uint256 value,uint256 nonce,uint256 deadline)"
      )
    );
    expect(mt).to.equal(MINT_TYPEHASH);
    expect(bt).to.equal(BURN_TYPEHASH);

    const message = `
Message Calls
Contracts can call other contracts or send Ether to non-contract accounts by the means of message calls. Message calls are similar to transactions, in that they have a source, a target, data payload, Ether, gas and return data. In fact, every transaction consists of a top-level message call which in turn can create further message calls.
A contract can decide how much of its remaining gas should be sent with the inner message call and how much it wants to retain. If an out-of-gas exception happens in the inner call (or any other exception), this will be signaled by an error value put onto the stack. In this case, only the gas sent together with the call is used up. In Solidity, the calling contract causes a manual exception by default in such situations, so that exceptions “bubble up” the call stack.
As already said, the called contract (which can be the same as the caller) will receive a freshly cleared instance of memory and has access to the call payload - which will be provided in a separate area called the calldata. After it has finished execution, it can return data which will be stored at a location in the caller’s memory preallocated by the caller. All such calls are fully synchronous.
Calls are limited to a depth of 1024, which means that for more complex operations, loops should be preferred over recursive calls. Furthermore, only 63/64th of the gas can be forwarded in a message call, which causes a depth limit of a little less than 1000 in practice.";
`;
    const amt = ethers.utils.parseUnits("1337");
    const tx2 = await zeta.burnSend(
      signers[1].address,
      amt,
      amt,
      "BSC",
      Buffer.from(message, "utf-8")
    );
    const receipt2 = await tx2.wait();
    args = receipt2.events[1].args;
    console.log(args);
    expect(receipt2.events[1].event).to.equal("BurnSend");
    expect(args[0]).to.equal(signers[0].address);
    expect(args[1]).to.equal(signers[1].address);
    expect(args[2]).to.equal(amt);
    expect(args[4]).to.equal("BSC");

    // console.log("message", args.message);
    expect(await zeta.balanceOf(signers[0].address)).to.equal(
      ethers.utils.parseUnits("7429")
    );

    const tx3 = await zeta.mint(
      signers[1].address,
      amt,
      ethers.utils.keccak256(Buffer.from("hello"))
    );
    const receipt3 = await tx3.wait();
    // console.log(receipt3);
    args = receipt3.events[1].args;
    // console.log(args)
    expect(
      args.sendHash ==
        ethers.utils.keccak256(Buffer.from("hello")).toString("hex")
    );
    expect(await zeta.balanceOf(signers[1].address)).to.equal(amt);
  });
});
