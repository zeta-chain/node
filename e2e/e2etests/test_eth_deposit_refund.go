package e2etests

import (
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestEtherDepositAndCallRefund(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestEtherDepositAndCallRefund requires exactly one argument for the amount.")
	}

	value, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestEtherDepositAndCallRefund.")
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
		cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.Info("cctx status message: %s", cctx.CctxStatus.StatusMessage)
		revertTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
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
