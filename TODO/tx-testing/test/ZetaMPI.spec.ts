import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { ethers } from "hardhat";
import { expect } from "chai";
import {
  deployZetaEth,
  deployZetaMpiBase,
  deployZetaMpiEth,
  deployZetaMpiNonEth,
  deployZetaNonEth,
  deployZetaReceiverMock,
} from "../lib/contracts.helpers";
import {
  ZetaMPIEth,
  ZetaMPINonEth,
  ZetaMPIBase,
  ZetaEth,
  ZetaNonEth,
} from "../typechain";
import { ZetaReceiverMock } from "../typechain/ZetaReceiverMock";

describe("ZetaMPI tests", () => {
  let zetaTokenEthContract: ZetaEth;
  let zetaTokenNonEthContract: ZetaNonEth;
  let zetaMpiBaseContract: ZetaMPIBase;
  let zetaMpiEthContract: ZetaMPIEth;
  let zetaReceiverMockContract: ZetaReceiverMock;
  let zetaMpiNonEthContract: ZetaMPINonEth;
  let tssUpdater: SignerWithAddress;
  let tssSigner: SignerWithAddress;
  let randomSigner: SignerWithAddress;

  const tssUpdaterApproveMpiEth = async () => {
    await (
      await zetaTokenEthContract.approve(zetaMpiEthContract.address, 100_000)
    ).wait();
  };

  const tssUpdaterApproveMpiNonEth = async () => {
    await (
      await zetaTokenNonEthContract.approve(
        zetaMpiNonEthContract.address,
        100_000
      )
    ).wait();
  };

  const transfer100kZetaEth = async (transferTo: string) => {
    await (await zetaTokenEthContract.transfer(transferTo, 100_000)).wait();
  };

  const transfer100kZetaNonEth = async (transferTo: string) => {
    await (await zetaTokenNonEthContract.transfer(transferTo, 100_000)).wait();
  };

  beforeEach(async () => {
    const accounts = await ethers.getSigners();
    [tssUpdater, tssSigner, randomSigner] = accounts;

    zetaTokenEthContract = await deployZetaEth({
      args: [100_000],
    });

    zetaTokenNonEthContract = await deployZetaNonEth({
      args: [100_000, tssSigner.address, tssUpdater.address],
    });

    zetaReceiverMockContract = await deployZetaReceiverMock();
    zetaMpiBaseContract = await deployZetaMpiBase({
      args: [
        zetaTokenEthContract.address,
        tssSigner.address,
        tssUpdater.address,
      ],
    });
    zetaMpiEthContract = await deployZetaMpiEth({
      args: [
        zetaTokenEthContract.address,
        tssSigner.address,
        tssUpdater.address,
      ],
    });
    zetaMpiNonEthContract = await deployZetaMpiNonEth({
      args: [
        zetaTokenNonEthContract.address,
        tssSigner.address,
        tssUpdater.address,
      ],
    });

    await zetaTokenNonEthContract.updateTSSAndMPIAddresses(
      tssSigner.address,
      zetaMpiNonEthContract.address
    );
  });

  describe("ZetaMPI.base", () => {
    describe("updateTssAddress", () => {
      it("Should revert if the caller is not the TSS updater", async () => {
        await expect(
          zetaMpiBaseContract
            .connect(randomSigner)
            .updateTssAddress(randomSigner.address)
        ).to.revertedWith("ZetaMPI: only TSS updater can call this function");
      });

      it("Should revert if the new TSS address is invalid", async () => {
        await expect(
          zetaMpiBaseContract.updateTssAddress(
            "0x0000000000000000000000000000000000000000"
          )
        ).to.revertedWith("ZetaMPI: invalid tssAddress");
      });

      it("Should change the TSS address if called by TSS updater", async () => {
        await (
          await zetaMpiBaseContract.updateTssAddress(randomSigner.address)
        ).wait();

        const address = await zetaMpiBaseContract.tssAddress();

        expect(address).to.equal(randomSigner.address);
      });
    });

    describe("pause, unpause", () => {
      it("Should revert if not called by the TSS updater", async () => {
        await expect(
          zetaMpiBaseContract.connect(randomSigner).pause()
        ).to.revertedWith("ZetaMPI: only TSS updater can call this function");

        await expect(
          zetaMpiBaseContract.connect(randomSigner).unpause()
        ).to.revertedWith("ZetaMPI: only TSS updater can call this function");
      });

      it("Should pause if called by the TSS updater", async () => {
        await (await zetaMpiBaseContract.pause()).wait();
        const paused1 = await zetaMpiBaseContract.paused();
        expect(paused1).to.equal(true);

        await (await zetaMpiBaseContract.unpause()).wait();
        const paused2 = await zetaMpiBaseContract.paused();
        expect(paused2).to.equal(false);
      });
    });
  });

  describe("ZetaMPI.eth", () => {
    describe("send", () => {
      it("Should revert if the contract is paused", async () => {
        await (await zetaMpiEthContract.pause()).wait();
        const paused1 = await zetaMpiEthContract.paused();
        expect(paused1).to.equal(true);

        await expect(
          zetaMpiEthContract.send({
            destinationAddress: randomSigner.address,
            destinationChainId: 1,
            gasLimit: 2500000,
            message: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            zetaAmount: 1000,
            zetaParams: new ethers.utils.AbiCoder().encode(
              ["string"],
              ["hello"]
            ),
          })
        ).to.revertedWith("Pausable: paused");
      });

      it("Should revert if the sender has no enough zeta", async () => {
        await (
          await zetaTokenEthContract
            .connect(randomSigner)
            .approve(zetaMpiEthContract.address, 100_000)
        ).wait();

        await expect(
          zetaMpiEthContract.connect(randomSigner).send({
            destinationAddress: randomSigner.address,
            destinationChainId: 1,
            gasLimit: 2500000,
            message: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            zetaAmount: 1000,
            zetaParams: new ethers.utils.AbiCoder().encode(
              ["string"],
              ["hello"]
            ),
          })
        ).to.revertedWith("ERC20: transfer amount exceeds balance");
      });

      it("Should revert if the sender didn't allow ZetaMPI to spend Zeta token", async () => {
        await expect(
          zetaMpiEthContract.send({
            destinationAddress: randomSigner.address,
            destinationChainId: 1,
            gasLimit: 2500000,
            message: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            zetaAmount: 1000,
            zetaParams: new ethers.utils.AbiCoder().encode(
              ["string"],
              ["hello"]
            ),
          })
        ).to.revertedWith("ERC20: transfer amount exceeds allowance");
      });

      it("Should transfer Zeta token from the sender account to the MPI contract", async () => {
        const initialBalanceDeployer = await zetaTokenEthContract.balanceOf(
          tssUpdater.address
        );
        const initialBalanceMpi = await zetaTokenEthContract.balanceOf(
          zetaMpiEthContract.address
        );

        expect(initialBalanceDeployer.toString()).to.equal(
          "100000000000000000000000"
        );
        expect(initialBalanceMpi.toString()).to.equal("0");

        await tssUpdaterApproveMpiEth();

        await (
          await zetaMpiEthContract.send({
            destinationAddress: randomSigner.address,
            destinationChainId: 1,
            gasLimit: 2500000,
            message: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            zetaAmount: 1000,
            zetaParams: new ethers.utils.AbiCoder().encode(
              ["string"],
              ["hello"]
            ),
          })
        ).wait();

        const finalBalanceDeployer = await zetaTokenEthContract.balanceOf(
          tssUpdater.address
        );
        const finalBalanceMpi = await zetaTokenEthContract.balanceOf(
          zetaMpiEthContract.address
        );

        expect(finalBalanceDeployer.toString()).to.equal(
          "99999999999999999999000"
        );
        expect(finalBalanceMpi.toString()).to.equal("1000");
      });

      it("Should emit `ZetaSent` on success", async () => {
        const zetaSentFilter = zetaMpiEthContract.filters.ZetaSent();
        const e1 = await zetaMpiEthContract.queryFilter(zetaSentFilter);
        expect(e1.length).to.equal(0);

        await zetaMpiEthContract.send({
          destinationAddress: randomSigner.address,
          destinationChainId: 1,
          gasLimit: 2500000,
          message: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
          zetaAmount: 0,
          zetaParams: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
        });

        const e2 = await zetaMpiEthContract.queryFilter(zetaSentFilter);
        expect(e2.length).to.equal(1);
      });
    });

    describe("onReceive", () => {
      it("Should revert if the contract is paused", async () => {
        await (await zetaMpiEthContract.pause()).wait();
        const paused1 = await zetaMpiEthContract.paused();
        expect(paused1).to.equal(true);

        await expect(
          zetaMpiEthContract.onReceive(
            tssUpdater.address,
            1,
            randomSigner.address,
            1000,
            new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            "0x0000000000000000000000000000000000000000000000000000000000000000"
          )
        ).to.revertedWith("Pausable: paused");
      });

      it("Should revert if not called by TSS address", async () => {
        await expect(
          zetaMpiEthContract.onReceive(
            tssUpdater.address,
            1,
            randomSigner.address,
            1000,
            new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            "0x0000000000000000000000000000000000000000000000000000000000000000"
          )
        ).to.revertedWith("ZetaMPI: only TSS address can call this function");
      });

      it("Should revert if Zeta transfer fails", async () => {
        await expect(
          zetaMpiEthContract
            .connect(tssSigner)
            .onReceive(
              randomSigner.address,
              1,
              randomSigner.address,
              1000,
              new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
              "0x0000000000000000000000000000000000000000000000000000000000000000"
            )
        ).to.revertedWith("ERC20: transfer amount exceeds balance");
      });

      it("Should transfer to the receiver address", async () => {
        await transfer100kZetaEth(zetaMpiEthContract.address);

        const initialBalanceMpi = await zetaTokenEthContract.balanceOf(
          zetaMpiEthContract.address
        );
        const initialBalanceReceiver = await zetaTokenEthContract.balanceOf(
          zetaReceiverMockContract.address
        );
        expect(initialBalanceMpi.toString()).to.equal("100000");
        expect(initialBalanceReceiver.toString()).to.equal("0");

        await (
          await zetaMpiEthContract
            .connect(tssSigner)
            .onReceive(
              randomSigner.address,
              1,
              zetaReceiverMockContract.address,
              1000,
              new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
              "0x0000000000000000000000000000000000000000000000000000000000000000"
            )
        ).wait();

        const finalBalanceMpi = await zetaTokenEthContract.balanceOf(
          zetaMpiEthContract.address
        );
        const finalBalanceReceiver = await zetaTokenEthContract.balanceOf(
          zetaReceiverMockContract.address
        );

        expect(finalBalanceMpi.toString()).to.equal("99000");
        expect(finalBalanceReceiver.toString()).to.equal("1000");
      });

      it("Should emit `ZetaReceived` on success", async () => {
        await transfer100kZetaEth(zetaMpiEthContract.address);

        const zetaReceivedFilter = zetaMpiEthContract.filters.ZetaReceived();
        const e1 = await zetaMpiEthContract.queryFilter(zetaReceivedFilter);
        expect(e1.length).to.equal(0);

        await (
          await zetaMpiEthContract
            .connect(tssSigner)
            .onReceive(
              randomSigner.address,
              1,
              zetaReceiverMockContract.address,
              1000,
              new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
              "0x0000000000000000000000000000000000000000000000000000000000000000"
            )
        ).wait();

        const e2 = await zetaMpiEthContract.queryFilter(zetaReceivedFilter);
        expect(e2.length).to.equal(1);
      });
    });

    describe("onRevert", () => {
      it("Should revert if the contract is paused", async () => {
        await (await zetaMpiEthContract.pause()).wait();
        const paused1 = await zetaMpiEthContract.paused();
        expect(paused1).to.equal(true);

        await expect(
          zetaMpiEthContract.onRevert(
            randomSigner.address,
            1,
            randomSigner.address,
            2,
            1000,
            new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            "0x0000000000000000000000000000000000000000000000000000000000000000"
          )
        ).to.revertedWith("Pausable: paused");
      });

      it("Should revert if not called by TSS address", async () => {
        await expect(
          zetaMpiEthContract.onRevert(
            randomSigner.address,
            1,
            tssUpdater.address,
            1,
            1000,
            new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            "0x0000000000000000000000000000000000000000000000000000000000000000"
          )
        ).to.revertedWith("ZetaMPI: only TSS address can call this function");
      });

      it("Should transfer to the origin address", async () => {
        await transfer100kZetaEth(zetaMpiEthContract.address);

        const initialBalanceMpi = await zetaTokenEthContract.balanceOf(
          zetaMpiEthContract.address
        );
        const initialBalanceSender = await zetaTokenEthContract.balanceOf(
          zetaReceiverMockContract.address
        );
        expect(initialBalanceMpi.toString()).to.equal("100000");
        expect(initialBalanceSender.toString()).to.equal("0");

        await (
          await zetaMpiEthContract
            .connect(tssSigner)
            .onRevert(
              zetaReceiverMockContract.address,
              1,
              randomSigner.address,
              1,
              1000,
              new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
              "0x0000000000000000000000000000000000000000000000000000000000000000"
            )
        ).wait();

        const finalBalanceMpi = await zetaTokenEthContract.balanceOf(
          zetaMpiEthContract.address
        );
        const finalBalanceSender = await zetaTokenEthContract.balanceOf(
          zetaReceiverMockContract.address
        );

        expect(finalBalanceMpi.toString()).to.equal("99000");
        expect(finalBalanceSender.toString()).to.equal("1000");
      });

      it("Should emit `ZetaReverted` on success", async () => {
        await transfer100kZetaEth(zetaMpiEthContract.address);

        const zetaRevertedFilter = zetaMpiEthContract.filters.ZetaReverted();
        const e1 = await zetaMpiEthContract.queryFilter(zetaRevertedFilter);
        expect(e1.length).to.equal(0);

        await (
          await zetaMpiEthContract
            .connect(tssSigner)
            .onRevert(
              zetaReceiverMockContract.address,
              1,
              randomSigner.address,
              1,
              1000,
              new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
              "0x0000000000000000000000000000000000000000000000000000000000000000"
            )
        ).wait();

        const e2 = await zetaMpiEthContract.queryFilter(zetaRevertedFilter);
        expect(e2.length).to.equal(1);
      });
    });
  });

  describe("ZetaMPI.non-eth", () => {
    describe("send", () => {
      it("Should revert if the contract is paused", async () => {
        await (await zetaMpiNonEthContract.pause()).wait();
        const paused1 = await zetaMpiNonEthContract.paused();
        expect(paused1).to.equal(true);

        await expect(
          zetaMpiNonEthContract.send({
            destinationAddress: randomSigner.address,
            destinationChainId: 1,
            gasLimit: 2500000,
            message: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            zetaAmount: 1000,
            zetaParams: new ethers.utils.AbiCoder().encode(
              ["string"],
              ["hello"]
            ),
          })
        ).to.revertedWith("Pausable: paused");
      });

      it("Should revert if the sender has no enough zeta", async () => {
        await expect(
          zetaMpiNonEthContract.connect(randomSigner).send({
            destinationAddress: randomSigner.address,
            destinationChainId: 1,
            gasLimit: 2500000,
            message: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            zetaAmount: 1000,
            zetaParams: new ethers.utils.AbiCoder().encode(
              ["string"],
              ["hello"]
            ),
          })
        ).to.revertedWith("ERC20: burn amount exceeds allowance");
      });

      it("Should revert if the sender didn't allow ZetaMPI to spend Zeta token", async () => {
        await expect(
          zetaMpiNonEthContract.send({
            destinationAddress: randomSigner.address,
            destinationChainId: 1,
            gasLimit: 2500000,
            message: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            zetaAmount: 1000,
            zetaParams: new ethers.utils.AbiCoder().encode(
              ["string"],
              ["hello"]
            ),
          })
        ).to.revertedWith("ERC20: burn amount exceeds allowance");
      });

      it("Should burn Zeta token from the sender account", async () => {
        const initialBalanceDeployer = await zetaTokenNonEthContract.balanceOf(
          tssUpdater.address
        );
        expect(initialBalanceDeployer.toString()).to.equal(
          "100000000000000000000000"
        );

        await tssUpdaterApproveMpiNonEth();

        await (
          await zetaMpiNonEthContract.send({
            destinationAddress: randomSigner.address,
            destinationChainId: 1,
            gasLimit: 2500000,
            message: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            zetaAmount: 1000,
            zetaParams: new ethers.utils.AbiCoder().encode(
              ["string"],
              ["hello"]
            ),
          })
        ).wait();

        const finalBalanceDeployer = await zetaTokenNonEthContract.balanceOf(
          tssUpdater.address
        );
        expect(finalBalanceDeployer.toString()).to.equal(
          "99999999999999999999000"
        );
      });

      it("Should emit `ZetaSent` on success", async () => {
        const zetaSentFilter = zetaMpiNonEthContract.filters.ZetaSent();
        const e1 = await zetaMpiNonEthContract.queryFilter(zetaSentFilter);
        expect(e1.length).to.equal(0);

        await zetaMpiNonEthContract.send({
          destinationAddress: randomSigner.address,
          destinationChainId: 1,
          gasLimit: 2500000,
          message: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
          zetaAmount: 0,
          zetaParams: new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
        });

        const e2 = await zetaMpiNonEthContract.queryFilter(zetaSentFilter);
        expect(e2.length).to.equal(1);
      });
    });

    describe("onReceive", () => {
      it("Should revert if the contract is paused", async () => {
        await (await zetaMpiNonEthContract.pause()).wait();
        const paused1 = await zetaMpiNonEthContract.paused();
        expect(paused1).to.equal(true);

        await expect(
          zetaMpiNonEthContract.onReceive(
            tssUpdater.address,
            1,
            randomSigner.address,
            1000,
            new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            "0x0000000000000000000000000000000000000000000000000000000000000000"
          )
        ).to.revertedWith("Pausable: paused");
      });

      it("Should revert if not called by TSS address", async () => {
        await expect(
          zetaMpiNonEthContract.onReceive(
            tssUpdater.address,
            1,
            randomSigner.address,
            1000,
            new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            "0x0000000000000000000000000000000000000000000000000000000000000000"
          )
        ).to.revertedWith("ZetaMPI: only TSS address can call this function");
      });

      it("Should revert if mint fails", async () => {
        await zetaTokenNonEthContract.updateTSSAndMPIAddresses(
          tssSigner.address,
          randomSigner.address
        );

        await expect(
          zetaMpiNonEthContract
            .connect(tssSigner)
            .onReceive(
              randomSigner.address,
              1,
              zetaReceiverMockContract.address,
              1000,
              new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
              "0x0000000000000000000000000000000000000000000000000000000000000000"
            )
        ).to.revertedWith("ZetaNonEth: only TSSAddress or MPIAddress can mint");
      });

      it("Should mint on the receiver address", async () => {
        const initialBalanceReceiver = await zetaTokenNonEthContract.balanceOf(
          zetaReceiverMockContract.address
        );
        expect(initialBalanceReceiver.toString()).to.equal("0");

        await (
          await zetaMpiNonEthContract
            .connect(tssSigner)
            .onReceive(
              randomSigner.address,
              1,
              zetaReceiverMockContract.address,
              1000,
              new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
              "0x0000000000000000000000000000000000000000000000000000000000000000"
            )
        ).wait();

        const finalBalanceReceiver = await zetaTokenNonEthContract.balanceOf(
          zetaReceiverMockContract.address
        );

        expect(finalBalanceReceiver.toString()).to.equal("1000");
      });

      it("Should emit `ZetaReceived` on success", async () => {
        await transfer100kZetaNonEth(zetaMpiNonEthContract.address);

        const zetaReceivedFilter = zetaMpiNonEthContract.filters.ZetaReceived();
        const e1 = await zetaMpiNonEthContract.queryFilter(zetaReceivedFilter);
        expect(e1.length).to.equal(0);

        await (
          await zetaMpiNonEthContract
            .connect(tssSigner)
            .onReceive(
              randomSigner.address,
              1,
              zetaReceiverMockContract.address,
              1000,
              new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
              "0x0000000000000000000000000000000000000000000000000000000000000000"
            )
        ).wait();

        const e2 = await zetaMpiNonEthContract.queryFilter(zetaReceivedFilter);
        expect(e2.length).to.equal(1);
      });
    });

    describe("onRevert", () => {
      it("Should revert if the contract is paused", async () => {
        await (await zetaMpiNonEthContract.pause()).wait();
        const paused1 = await zetaMpiNonEthContract.paused();
        expect(paused1).to.equal(true);

        await expect(
          zetaMpiNonEthContract.onRevert(
            randomSigner.address,
            1,
            randomSigner.address,
            2,
            1000,
            new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            "0x0000000000000000000000000000000000000000000000000000000000000000"
          )
        ).to.revertedWith("Pausable: paused");
      });

      it("Should revert if not called by TSS address", async () => {
        await expect(
          zetaMpiNonEthContract.onRevert(
            randomSigner.address,
            1,
            tssUpdater.address,
            1,
            1000,
            new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
            "0x0000000000000000000000000000000000000000000000000000000000000000"
          )
        ).to.revertedWith("ZetaMPI: only TSS address can call this function");
      });

      it("Should mint on the origin address", async () => {
        const initialBalanceSender = await zetaTokenNonEthContract.balanceOf(
          zetaReceiverMockContract.address
        );
        expect(initialBalanceSender.toString()).to.equal("0");

        await (
          await zetaMpiNonEthContract
            .connect(tssSigner)
            .onRevert(
              zetaReceiverMockContract.address,
              1,
              randomSigner.address,
              1,
              1000,
              new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
              "0x0000000000000000000000000000000000000000000000000000000000000000"
            )
        ).wait();

        const finalBalanceSender = await zetaTokenNonEthContract.balanceOf(
          zetaReceiverMockContract.address
        );
        expect(finalBalanceSender.toString()).to.equal("1000");
      });

      it("Should emit `ZetaReverted` on success", async () => {
        await transfer100kZetaNonEth(zetaMpiNonEthContract.address);

        const zetaRevertedFilter = zetaMpiNonEthContract.filters.ZetaReverted();
        const e1 = await zetaMpiNonEthContract.queryFilter(zetaRevertedFilter);
        expect(e1.length).to.equal(0);

        await (
          await zetaMpiNonEthContract
            .connect(tssSigner)
            .onRevert(
              zetaReceiverMockContract.address,
              1,
              randomSigner.address,
              1,
              1000,
              new ethers.utils.AbiCoder().encode(["string"], ["hello"]),
              "0x0000000000000000000000000000000000000000000000000000000000000000"
            )
        ).wait();

        const e2 = await zetaMpiNonEthContract.queryFilter(zetaRevertedFilter);
        expect(e2.length).to.equal(1);
      });
    });
  });
});
