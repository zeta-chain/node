const { expect } = require("chai");
const { ethers } = require("hardhat");
const { ecsign } = require("ethereumjs-utils");

const keccak256 = ethers.utils.keccak256;
const toUtf8Bytes = ethers.utils.toUtf8Bytes;
const solidityPack = ethers.utils.solidityPack;
const MaxUint256 = ethers.constants.MaxUint256;
const defaultAbiCoder = ethers.utils.defaultAbiCoder;
const hexlify = ethers.utils.hexlify;
const BigNumber = ethers.BigNumber;
const parseUnits = ethers.utils.parseUnits;

describe("Router Contract", function () {
  it("deposit/withdraw router", async function () {
    const signers = await ethers.getSigners();
    deployer = signers[0].address;

    const TestOracle = await ethers.getContractFactory("TestSupplyOracle");
    const oracle = await TestOracle.deploy();
    await oracle.deployed();
    console.log("test oracle deployed at", oracle.address);

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
    expect(await zeta.totalSupply()).to.equal(ethers.utils.parseUnits(supply));
    oracle.setSupply(ethers.utils.parseUnits(supply));
    oracle.setLocked(ethers.utils.parseUnits(supply));

    const DRouter = await ethers.getContractFactory("DepositRouter");
    const drouter = await DRouter.deploy(zeta.address);
    await drouter.deployed();
    zeta.approve(drouter.address, ethers.utils.parseUnits("1337"));
    const tx1 = await drouter.send_m(
      123,
      signers[1].address,
      ethers.utils.parseUnits("1337")
    );
    await tx1.wait();
    expect(await zeta.totalSupply()).to.equal(ethers.utils.parseUnits("8663"));
    oracle.setSupply(ethers.utils.parseUnits("8663"));

    // this one should fail because allowance is used up.
    try {
      const tx2 = await drouter.send_m(
        223,
        signers[1].address,
        ethers.utils.parseUnits("1")
      );
      await tx2.wait();
      expect(1).to.equal(2);
    } catch (error) {
      // console.log(error);
    }

    // test the send_m_by_signature
    const nonce = await zeta.nonces(signers[0].address);
    const deadline = MaxUint256;
    const mAmount = parseUnits("1337");

    const DOMAIN_SEPARATOR = await zeta.DOMAIN_SEPARATOR();
    const sender = signers[0].address;
    const digest = keccak256(
      solidityPack(
        ["bytes2", "bytes32", "bytes32"],
        [
          "0x1901",
          DOMAIN_SEPARATOR,
          keccak256(
            defaultAbiCoder.encode(
              ["bytes32", "address", "uint256", "uint256", "uint256"],
              [await zeta.BURN_TYPEHASH(), sender, mAmount, nonce, deadline]
            )
          ),
        ]
      )
    );

    // this is the prikey of address: 0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266
    const privateKey =
      "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80";
    var { v, r, s } = ecsign(
      Buffer.from(digest.slice(2), "hex"), // msgHash
      Buffer.from(privateKey.slice(2), "hex")
    );
    const tx3 = await drouter.send_m_by_signature(
      123,
      signers[1].address,
      signers[0].address,
      mAmount,
      deadline,
      v,
      hexlify(r),
      hexlify(s)
    );
    const receipt = await tx3.wait();
    expect(receipt.events[2].event).to.equal("SendM");
    expect(await zeta.totalSupply()).to.equal(parseUnits("7326")); // 10000 - 1337 - 1337 = 7326

    // test mint

    const WRouter = await ethers.getContractFactory("ZetaExecutor");
    const wrouter = await WRouter.deploy(zeta.address);
    await drouter.deployed();
    const Receiver = await ethers.getContractFactory("TestReceiver");
    const receiver = await Receiver.deploy();
    await receiver.deployed();
    const wdigest = keccak256(
      solidityPack(
        ["bytes2", "bytes32", "bytes32"],
        [
          "0x1901",
          DOMAIN_SEPARATOR,
          keccak256(
            defaultAbiCoder.encode(
              ["bytes32", "address", "uint256", "uint256", "uint256"],
              [
                await zeta.MINT_TYPEHASH(),
                receiver.address,
                mAmount,
                await zeta.nonces(sender),
                deadline,
              ]
            )
          ),
        ]
      )
    );
    var { v, r, s } = ecsign(
      Buffer.from(wdigest.slice(2), "hex"),
      Buffer.from(privateKey.slice(2), "hex")
    );
    const tx4 = await wrouter.mint_m_by_signature(
      receiver.address,
      mAmount,
      deadline,
      v,
      hexlify(r),
      hexlify(s),
      Buffer.from("hi"),
      ethers.utils.keccak256(Buffer.from("sendhash"))
    );
    const receipt4 = await tx4.wait();
    expect(await zeta.balanceOf(receiver.address)).to.equal(parseUnits("1337"));
    expect(await zeta.totalSupply()).to.equal(ethers.utils.parseUnits("8663"));
  });
});
