package runner

import (
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
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

	cctx := utils.WaitCctxMinedByInTxHash(runner.Ctx, txHash.Hex(), runner.CctxClient, runner.Logger, runner.CctxTimeout)
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
