package e2etests

import (
	"fmt"

	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/connectorzevm.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestZetaWithdraw(r *runner.E2ERunner) {
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(10)) // 10 Zeta

	r.ZevmAuth.Value = amount
	tx, err := r.WZeta.Deposit(r.ZevmAuth)
	if err != nil {
		panic(err)
	}
	r.ZevmAuth.Value = big.NewInt(0)
	r.Logger.Info("wzeta deposit tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta deposit")
	if receipt.Status == 0 {
		panic("deposit failed")
	}

	chainID, err := r.GoerliClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}

	tx, err = r.WZeta.Approve(r.ZevmAuth, r.ConnectorZEVMAddr, amount)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("wzeta approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta approve")
	if receipt.Status == 0 {
		panic(fmt.Sprintf("approve failed, logs: %+v", receipt.Logs))
	}

	tx, err = r.ConnectorZEVM.Send(r.ZevmAuth, connectorzevm.ZetaInterfacesSendInput{
		DestinationChainId:  chainID,
		DestinationAddress:  r.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	r.Logger.Info("send tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "send")
	if receipt.Status == 0 {
		panic(fmt.Sprintf("send failed, logs: %+v", receipt.Logs))

	}

	r.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := r.ConnectorZEVM.ParseZetaSent(*log)
		if err == nil {
			r.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			r.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			r.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			r.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
		}
	}
	r.Logger.Info("waiting for cctx status to change to final...")

	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta withdraw")
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_OutboundMined {
		panic(fmt.Errorf(
			"expected cctx status to be %s; got %s, message %s",
			cctxtypes.CctxStatus_OutboundMined,
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage,
		))
	}
}

func TestZetaWithdrawBTCRevert(r *runner.E2ERunner) {
	r.ZevmAuth.Value = big.NewInt(1e18) // 1 Zeta
	tx, err := r.WZeta.Deposit(r.ZevmAuth)
	if err != nil {
		panic(err)
	}
	r.ZevmAuth.Value = big.NewInt(0)
	r.Logger.Info("Deposit tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "Deposit")
	if receipt.Status != 1 {
		panic("Deposit failed")
	}

	tx, err = r.WZeta.Approve(r.ZevmAuth, r.ConnectorZEVMAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	r.Logger.Info("wzeta.approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "Approve")
	if receipt.Status != 1 {
		panic("Approve failed")
	}

	tx, err = r.ConnectorZEVM.Send(r.ZevmAuth, connectorzevm.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(common.BtcRegtestChain().ChainId),
		DestinationAddress:  r.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		ZetaValueAndGas:     big.NewInt(1e17),
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	r.Logger.Info("send tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "send")
	if receipt.Status != 0 {
		panic("Was able to send ZETA to BTC")
	}
}
