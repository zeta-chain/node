package e2etests

import (
	"math/big"

	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zetaconnectorzevm.sol"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

func TestZetaWithdrawBTCRevert(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestZetaWithdrawBTCRevert requires exactly one argument for the withdrawal.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestZetaWithdrawBTCRevert.")
	}

	r.ZEVMAuth.Value = amount
	tx, err := r.WZeta.Deposit(r.ZEVMAuth)
	if err != nil {
		panic(err)
	}
	r.ZEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("Deposit tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "Deposit")
	if receipt.Status != 1 {
		panic("Deposit failed")
	}

	tx, err = r.WZeta.Approve(r.ZEVMAuth, r.ConnectorZEVMAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	r.Logger.Info("wzeta.approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "Approve")
	if receipt.Status != 1 {
		panic("Approve failed")
	}

	lessThanAmount := amount.Div(amount, big.NewInt(10)) // 1/10 of amount
	tx, err = r.ConnectorZEVM.Send(r.ZEVMAuth, connectorzevm.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(chains.BtcRegtestChain.ChainId),
		DestinationAddress:  r.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		ZetaValueAndGas:     lessThanAmount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	r.Logger.Info("send tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "send")
	if receipt.Status != 0 {
		panic("Was able to send ZETA to BTC")
	}
}
