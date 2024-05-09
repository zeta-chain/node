package e2etests

import (
	"fmt"

	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// TestEtherDeposit tests deposit of ethers
func TestEtherDeposit(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestEtherDeposit requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestEtherDeposit.")
	}

	hash := r.DepositEtherWithAmount(false, amount) // in wei
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
}

// TestEtherDepositAndCall tests deposit of ethers calling a example contract
func TestEtherDepositAndCall(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestEtherDepositAndCall requires exactly one argument for the amount.")
	}

	value, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestEtherDepositAndCall.")
	}

	r.Logger.Info("Deploying example contract")
	exampleAddr, _, exampleContract, err := testcontract.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("Example contract deployed")

	// preparing tx
	evmClient := r.EVMClient
	gasLimit := uint64(23000)
	gasPrice, err := evmClient.SuggestGasPrice(r.Ctx)
	if err != nil {
		panic(err)
	}
	nonce, err := evmClient.PendingNonceAt(r.Ctx, r.DeployerAddress)
	if err != nil {
		panic(err)
	}

	data := append(exampleAddr.Bytes(), []byte("hello sailors")...)
	tx := ethtypes.NewTransaction(nonce, r.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := evmClient.NetworkID(r.Ctx)
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(r.DeployerPrivateKey)
	if err != nil {
		panic(err)
	}
	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("Sending a cross-chain call to example contract")
	err = evmClient.SendTransaction(r.Ctx, signedTx)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("tx failed")
	}
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf("expected cctx status to be mined; got %s", cctx.CctxStatus.Status))
	}

	// Checking example contract has been called, bar value should be set to amount
	bar, err := exampleContract.Bar(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if bar.Cmp(value) != 0 {
		panic(fmt.Sprintf("cross-chain call failed bar value %s should be equal to amount %s", bar.String(), value.String()))
	}
	r.Logger.Info("Cross-chain call succeeded")

	r.Logger.Info("Deploying reverter contract")
	reverterAddr, _, _, err := testcontract.DeployReverter(r.ZEVMAuth, r.ZEVMClient)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("Example reverter deployed")

	// preparing tx for reverter
	gasPrice, err = evmClient.SuggestGasPrice(r.Ctx)
	if err != nil {
		panic(err)
	}
	nonce, err = evmClient.PendingNonceAt(r.Ctx, r.DeployerAddress)
	if err != nil {
		panic(err)
	}

	data = append(reverterAddr.Bytes(), []byte("hello sailors")...)
	tx = ethtypes.NewTransaction(nonce, r.TSSAddress, value, gasLimit, gasPrice, data)
	signedTx, err = ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("Sending a cross-chain call to reverter contract")
	err = evmClient.SendTransaction(r.Ctx, signedTx)
	if err != nil {
		panic(err)
	}

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("tx failed")
	}

	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
		panic(fmt.Sprintf("expected cctx status to be reverted; got %s", cctx.CctxStatus.Status))
	}
	r.Logger.Info("Cross-chain call to reverter reverted")

	// check the status message contains revert error hash in case of revert
	// 0xbfb4ebcf is the hash of "Foo()"
	if !strings.Contains(cctx.CctxStatus.StatusMessage, "0xbfb4ebcf") {
		panic(fmt.Sprintf("expected cctx status message to contain revert reason; got %s", cctx.CctxStatus.StatusMessage))
	}
}

func TestDepositAndCallRefund(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestDepositAndCallRefund requires exactly one argument for the amount.")
	}

	value, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestDepositAndCallRefund.")
	}

	evmClient := r.EVMClient

	nonce, err := evmClient.PendingNonceAt(r.Ctx, r.DeployerAddress)
	if err != nil {
		panic(err)
	}

	gasLimit := uint64(23000) // in units
	gasPrice, err := evmClient.SuggestGasPrice(r.Ctx)
	if err != nil {
		panic(err)
	}

	data := append(r.BTCZRC20Addr.Bytes(), []byte("hello sailors")...) // this data
	tx := ethtypes.NewTransaction(nonce, r.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := evmClient.NetworkID(r.Ctx)
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(r.DeployerPrivateKey)
	if err != nil {
		panic(err)
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		panic(err)
	}
	err = evmClient.SendTransaction(r.Ctx, signedTx)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("EVM tx sent: %s; to %s, nonce %d", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("EVM tx receipt: %d", receipt.Status)
	r.Logger.Info("  tx hash: %s", receipt.TxHash.String())
	r.Logger.Info("  to: %s", signedTx.To().String())
	r.Logger.Info("  value: %d", signedTx.Value())
	r.Logger.Info("  block num: %d", receipt.BlockNumber)

	func() {
		cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.Info("cctx status message: %s", cctx.CctxStatus.StatusMessage)
		revertTxHash := cctx.GetCurrentOutboundParam().Hash
		r.Logger.Info("EVM revert tx receipt: status %d", receipt.Status)

		tx, _, err := r.EVMClient.TransactionByHash(r.Ctx, ethcommon.HexToHash(revertTxHash))
		if err != nil {
			panic(err)
		}
		receipt, err := r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(revertTxHash))
		if err != nil {
			panic(err)
		}

		printTxInfo := func() {
			// debug info when test fails
			r.Logger.Info("  tx: %+v", tx)
			r.Logger.Info("  receipt: %+v", receipt)
			r.Logger.Info("cctx http://localhost:1317/zeta-chain/crosschain/cctx/%s", cctx.Index)
		}

		if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
			printTxInfo()
			panic(fmt.Sprintf("expected cctx status to be PendingRevert; got %s", cctx.CctxStatus.Status))
		}

		if receipt.Status == 0 {
			printTxInfo()
			panic("expected the revert tx receipt to have status 1; got 0")
		}

		if *tx.To() != r.DeployerAddress {
			printTxInfo()
			panic(fmt.Sprintf("expected tx to %s; got %s", r.DeployerAddress.Hex(), tx.To().Hex()))
		}

		// the received value must be lower than the original value because of the paid fees for the revert tx
		// we check that the value is still greater than 0
		if tx.Value().Cmp(value) != -1 || tx.Value().Cmp(big.NewInt(0)) != 1 {
			printTxInfo()
			panic(fmt.Sprintf("expected tx value %s; should be non-null and lower than %s", tx.Value().String(), value.String()))
		}

		r.Logger.Info("REVERT tx receipt: %d", receipt.Status)
		r.Logger.Info("  tx hash: %s", receipt.TxHash.String())
		r.Logger.Info("  to: %s", tx.To().String())
		r.Logger.Info("  value: %s", tx.Value().String())
		r.Logger.Info("  block num: %d", receipt.BlockNumber)
	}()
}
