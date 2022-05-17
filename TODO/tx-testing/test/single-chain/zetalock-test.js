const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Zeta Locker Contract", function () {
  it("lock, unlock", async function () {
    const Zeta = await ethers.getContractFactory("ZetaEth");
    const supply = "2100000000";
    const zeta = await Zeta.deploy(supply, "Zeta", "ZETA");
    await zeta.deployed();
    console.log("contract deployed at address", zeta.address);
    console.log("deployer of contract", zeta.deployTransaction.from);
    const signers = await ethers.getSigners();
    expect(signers[0].address).to.equal(zeta.deployTransaction.from);
    expect(await zeta.totalSupply()).to.equal(ethers.utils.parseUnits(supply));

    const TestOracle = await ethers.getContractFactory("TestSupplyOracle");
    const oracle = await TestOracle.deploy();
    await oracle.deployed();
    console.log("test oracle deployed at", oracle.address);

    const Lock = await ethers.getContractFactory("ZetaLocker");
    const deployer = signers[0].address;
    const locker = await Lock.deploy(
      zeta.address,
      oracle.address,
      deployer,
      deployer,
      deployer
    );
    await locker.deployed();
    console.log("ZetaLocker contract deployed at address", locker.address);
    expect(await locker.getLockedAmount()).to.equal(0);

    let amt = ethers.utils.parseUnits("1337");
    const msg = Buffer.from("please accept my gift", "utf-8");
    await zeta.approve(locker.address, amt);
    const tx2 = await locker.lockSend(
      signers[1].address,
      amt,
      amt,
      "Ethereum#137",
      msg
    );
    const receipt2 = await tx2.wait();
    // console.log(receipt2.events[2].args);
    expect(receipt2.events[2].args[5]).to.equal("0x" + msg.toString("hex"));
    expect(await locker.getLockedAmount()).to.equal(amt);
    console.log("lockSend gas cost:", receipt2.gasUsed.toString());

    amt = ethers.utils.parseUnits("1204");
    await oracle.setSupply(amt);

    const amt2 = ethers.utils.parseUnits("133");
    const tx3 = await locker.unlock(
      signers[2].address,
      amt2,
      ethers.utils.keccak256(Buffer.from("hello"))
    );
    const receipt3 = await tx3.wait();
    console.log("unlock cost gas: ", receipt3.gasUsed.toString());
    expect(await zeta.balanceOf(signers[2].address)).to.equal(amt2);

    const amt3 = ethers.utils.parseUnits("28"); // should fail to unlock
    try {
      await locker.unlock(
        signers[2].address,
        amt3,
        ethers.utils.keccak256(Buffer.from("hello"))
      );
    } catch (error) {
      console.log(error);
    }

    await locker.updateTSSAddress(signers[1].address);
    expect(await locker.TSSAddress()).to.equal(signers[1].address);

    try {
      await locker.unlock(
        signers[1].address,
        amt,
        ethers.utils.keccak256(Buffer.from("hello"))
      );
      expect(1).to.equal(2);
    } catch (error) {
      console.log(error);
    }
  });
});
