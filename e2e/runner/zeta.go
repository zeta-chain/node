package runner

import (
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zetaconnectorzevm.sol"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// WaitForTxReceiptOnZEVM waits for a tx receipt on ZEVM
func (runner *E2ERunner) WaitForTxReceiptOnZEVM(tx *ethtypes.Transaction) {
	defer func() {
		runner.Unlock()
	}()
	runner.Lock()

	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.ZEVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
}

// WaitForMinedCCTX waits for a cctx to be mined from a tx
func (runner *E2ERunner) WaitForMinedCCTX(txHash ethcommon.Hash) {
	defer func() {
		runner.Unlock()
	}()
	runner.Lock()

	cctx := utils.WaitCctxMinedByInboundHash(runner.Ctx, txHash.Hex(), runner.CctxClient, runner.Logger, runner.CctxTimeout)
	if cctx.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf("expected cctx status to be mined; got %s, message: %s",
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage),
		)
	}
}

// WaitForMinedCCTXFromIndex waits for a cctx to be mined from its index
func (runner *E2ERunner) WaitForMinedCCTXFromIndex(index string) {
	defer func() {
		runner.Unlock()
	}()
	runner.Lock()

	cctx := utils.WaitCCTXMinedByIndex(runner.Ctx, index, runner.CctxClient, runner.Logger, runner.CctxTimeout)
	if cctx.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf("expected cctx status to be mined; got %s, message: %s",
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage),
		)
	}
}

// SendZetaOnEvm sends ZETA to an address on EVM
// this allows the ZETA contract deployer to funds other accounts on EVM
func (runner *E2ERunner) SendZetaOnEvm(address ethcommon.Address, zetaAmount int64) *ethtypes.Transaction {
	// the deployer might be sending ZETA in different goroutines
	defer func() {
		runner.Unlock()
	}()
	runner.Lock()

	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(zetaAmount))
	tx, err := runner.ZetaEth.Transfer(runner.EVMAuth, address, amount)
	if err != nil {
		panic(err)
	}
	return tx
}

// DepositZeta deposits ZETA on ZetaChain from the ZETA smart contract on EVM
func (runner *E2ERunner) DepositZeta() ethcommon.Hash {
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(100)) // 100 Zeta

	return runner.DepositZetaWithAmount(runner.DeployerAddress, amount)
}

// DepositZetaWithAmount deposits ZETA on ZetaChain from the ZETA smart contract on EVM with the specified amount
func (runner *E2ERunner) DepositZetaWithAmount(to ethcommon.Address, amount *big.Int) ethcommon.Hash {
	tx, err := runner.ZetaEth.Approve(runner.EVMAuth, runner.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	runner.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.EVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	runner.Logger.EVMReceipt(*receipt, "approve")
	if receipt.Status != 1 {
		panic("approve tx failed")
	}

	// query the chain ID using zevm client
	zetaChainID, err := runner.ZEVMClient.ChainID(runner.Ctx)
	if err != nil {
		panic(err)
	}

	tx, err = runner.ConnectorEth.Send(runner.EVMAuth, zetaconnectoreth.ZetaInterfacesSendInput{
		// TODO: allow user to specify destination chain id
		// https://github.com/zeta-chain/node-private/issues/41
		DestinationChainId:  zetaChainID,
		DestinationAddress:  to.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	runner.Logger.Info("Send tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(runner.Ctx, runner.EVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	runner.Logger.EVMReceipt(*receipt, "send")
	if receipt.Status != 1 {
		panic(fmt.Sprintf("expected tx receipt status to be 1; got %d", receipt.Status))
	}

	runner.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := runner.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			runner.Logger.Info("    Connector: %s", runner.ConnectorEthAddr.String())
			runner.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			runner.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			runner.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			runner.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
			runner.Logger.Info("    Block Num: %d", log.BlockNumber)
		}
	}

	return tx.Hash()
}

// DepositAndApproveWZeta deposits and approves WZETA on ZetaChain from the ZETA smart contract on ZEVM
func (runner *E2ERunner) DepositAndApproveWZeta(amount *big.Int) {
	runner.ZEVMAuth.Value = amount
	tx, err := runner.WZeta.Deposit(runner.ZEVMAuth)
	if err != nil {
		panic(err)
	}
	runner.ZEVMAuth.Value = big.NewInt(0)
	runner.Logger.Info("wzeta deposit tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.ZEVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	runner.Logger.EVMReceipt(*receipt, "wzeta deposit")
	if receipt.Status == 0 {
		panic("deposit failed")
	}

	tx, err = runner.WZeta.Approve(runner.ZEVMAuth, runner.ConnectorZEVMAddr, amount)
	if err != nil {
		panic(err)
	}
	runner.Logger.Info("wzeta approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(runner.Ctx, runner.ZEVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	runner.Logger.EVMReceipt(*receipt, "wzeta approve")
	if receipt.Status == 0 {
		panic(fmt.Sprintf("approve failed, logs: %+v", receipt.Logs))
	}
}

// WithdrawZeta withdraws ZETA from ZetaChain to the ZETA smart contract on EVM
// waitReceipt specifies whether to wait for the tx receipt and check if the tx was successful
func (runner *E2ERunner) WithdrawZeta(amount *big.Int, waitReceipt bool) *ethtypes.Transaction {
	chainID, err := runner.EVMClient.ChainID(runner.Ctx)
	if err != nil {
		panic(err)
	}

	tx, err := runner.ConnectorZEVM.Send(runner.ZEVMAuth, connectorzevm.ZetaInterfacesSendInput{
		DestinationChainId:  chainID,
		DestinationAddress:  runner.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	runner.Logger.Info("send tx hash: %s", tx.Hash().Hex())

	if waitReceipt {
		receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.ZEVMClient, tx, runner.Logger, runner.ReceiptTimeout)
		runner.Logger.EVMReceipt(*receipt, "send")
		if receipt.Status == 0 {
			panic(fmt.Sprintf("send failed, logs: %+v", receipt.Logs))

		}

		runner.Logger.Info("  Logs:")
		for _, log := range receipt.Logs {
			sentLog, err := runner.ConnectorZEVM.ParseZetaSent(*log)
			if err == nil {
				runner.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
				runner.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
				runner.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
				runner.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
			}
		}
	}

	return tx
}

// WithdrawEther withdraws Ether from ZetaChain to the ZETA smart contract on EVM
func (runner *E2ERunner) WithdrawEther(amount *big.Int) *ethtypes.Transaction {
	// withdraw
	tx, err := runner.ETHZRC20.Withdraw(runner.ZEVMAuth, runner.DeployerAddress.Bytes(), amount)
	if err != nil {
		panic(err)
	}
	runner.Logger.EVMTransaction(*tx, "withdraw")

	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.ZEVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("withdraw failed")
	}
	runner.Logger.EVMReceipt(*receipt, "withdraw")
	runner.Logger.ZRC20Withdrawal(runner.ETHZRC20, *receipt, "withdraw")

	return tx
}

// WithdrawERC20 withdraws an ERC20 token from ZetaChain to the ZETA smart contract on EVM
func (runner *E2ERunner) WithdrawERC20(amount *big.Int) *ethtypes.Transaction {
	tx, err := runner.ERC20ZRC20.Withdraw(runner.ZEVMAuth, runner.DeployerAddress.Bytes(), amount)
	if err != nil {
		panic(err)
	}
	runner.Logger.EVMTransaction(*tx, "withdraw")

	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.ZEVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	runner.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := runner.ERC20ZRC20.ParseWithdrawal(*log)
		if err != nil {
			continue
		}
		runner.Logger.Info(
			"  logs: from %s, to %x, value %d, gasfee %d",
			event.From.Hex(),
			event.To,
			event.Value,
			event.Gasfee,
		)
	}

	return tx
}
