package smoketests

import (
	"fmt"
	"math/big"
	"strings"

	"cosmossdk.io/math"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

// TestEtherDeposit tests deposit of ethers
func TestEtherDeposit(sm *runner.SmokeTestRunner) {
	sm.DepositEther(false)
}

// TestEtherDepositAndCall tests deposit of ethers calling a example contract
func TestEtherDepositAndCall(sm *runner.SmokeTestRunner) {
	sm.Logger.Info("Deploying example contract")
	exampleAddr, _, exampleContract, err := testcontract.DeployExample(sm.ZevmAuth, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Example contract deployed")

	// preparing tx
	goerliClient := sm.GoerliClient
	value := big.NewInt(1e18)
	gasLimit := uint64(23000)
	gasPrice, err := goerliClient.SuggestGasPrice(sm.Ctx)
	if err != nil {
		panic(err)
	}
	nonce, err := goerliClient.PendingNonceAt(sm.Ctx, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}

	data := append(exampleAddr.Bytes(), []byte("hello sailors")...)
	tx := ethtypes.NewTransaction(nonce, sm.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := goerliClient.NetworkID(sm.Ctx)
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(sm.DeployerPrivateKey)
	if err != nil {
		panic(err)
	}
	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("Sending a cross-chain call to example contract")
	err = goerliClient.SendTransaction(sm.Ctx, signedTx)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, signedTx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("tx failed")
	}
	cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, signedTx.Hash().Hex(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
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
	sm.Logger.Info("Cross-chain call succeeded")

	sm.Logger.Info("Deploying reverter contract")
	reverterAddr, _, _, err := testcontract.DeployReverter(sm.ZevmAuth, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Example reverter deployed")

	// preparing tx for reverter
	gasPrice, err = goerliClient.SuggestGasPrice(sm.Ctx)
	if err != nil {
		panic(err)
	}
	nonce, err = goerliClient.PendingNonceAt(sm.Ctx, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}

	data = append(reverterAddr.Bytes(), []byte("hello sailors")...)
	tx = ethtypes.NewTransaction(nonce, sm.TSSAddress, value, gasLimit, gasPrice, data)
	signedTx, err = ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("Sending a cross-chain call to reverter contract")
	err = goerliClient.SendTransaction(sm.Ctx, signedTx)
	if err != nil {
		panic(err)
	}

	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, signedTx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("tx failed")
	}
	cctx = utils.WaitCctxMinedByInTxHash(sm.Ctx, signedTx.Hash().Hex(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
		panic(fmt.Sprintf("expected cctx status to be reverted; got %s", cctx.CctxStatus.Status))
	}
	sm.Logger.Info("Cross-chain call to reverter reverted")

	// check the status message contains revert error hash in case of revert
	// 0xbfb4ebcf is the hash of "Foo()"
	if !strings.Contains(cctx.CctxStatus.StatusMessage, "0xbfb4ebcf") {
		panic(fmt.Sprintf("expected cctx status message to contain revert reason; got %s", cctx.CctxStatus.StatusMessage))
	}
}

func TestDepositAndCallRefund(sm *runner.SmokeTestRunner) {
	goerliClient := sm.GoerliClient

	// in wei (10 eth)
	value := big.NewInt(1e18)
	value = value.Mul(value, big.NewInt(10))

	nonce, err := goerliClient.PendingNonceAt(sm.Ctx, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}

	gasLimit := uint64(23000) // in units
	gasPrice, err := goerliClient.SuggestGasPrice(sm.Ctx)
	if err != nil {
		panic(err)
	}

	data := append(sm.BTCZRC20Addr.Bytes(), []byte("hello sailors")...) // this data
	tx := ethtypes.NewTransaction(nonce, sm.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := goerliClient.NetworkID(sm.Ctx)
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(sm.DeployerPrivateKey)
	if err != nil {
		panic(err)
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		panic(err)
	}
	err = goerliClient.SendTransaction(sm.Ctx, signedTx)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("GOERLI tx sent: %s; to %s, nonce %d", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, signedTx, sm.Logger, sm.ReceiptTimeout)
	sm.Logger.Info("GOERLI tx receipt: %d", receipt.Status)
	sm.Logger.Info("  tx hash: %s", receipt.TxHash.String())
	sm.Logger.Info("  to: %s", signedTx.To().String())
	sm.Logger.Info("  value: %d", signedTx.Value())
	sm.Logger.Info("  block num: %d", receipt.BlockNumber)

	func() {
		cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, signedTx.Hash().Hex(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
		sm.Logger.Info("cctx status message: %s", cctx.CctxStatus.StatusMessage)
		revertTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
		sm.Logger.Info("GOERLI revert tx receipt: status %d", receipt.Status)

		tx, _, err := sm.GoerliClient.TransactionByHash(sm.Ctx, ethcommon.HexToHash(revertTxHash))
		if err != nil {
			panic(err)
		}
		receipt, err := sm.GoerliClient.TransactionReceipt(sm.Ctx, ethcommon.HexToHash(revertTxHash))
		if err != nil {
			panic(err)
		}

		printTxInfo := func() {
			// debug info when test fails
			sm.Logger.Info("  tx: %+v", tx)
			sm.Logger.Info("  receipt: %+v", receipt)
			sm.Logger.Info("cctx http://localhost:1317/zeta-chain/crosschain/cctx/%s", cctx.Index)
		}

		if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
			printTxInfo()
			panic(fmt.Sprintf("expected cctx status to be PendingRevert; got %s", cctx.CctxStatus.Status))
		}

		if receipt.Status == 0 {
			printTxInfo()
			panic("expected the revert tx receipt to have status 1; got 0")
		}

		if *tx.To() != sm.DeployerAddress {
			printTxInfo()
			panic(fmt.Sprintf("expected tx to %s; got %s", sm.DeployerAddress.Hex(), tx.To().Hex()))
		}

		// the received value must be lower than the original value because of the paid fees for the revert tx
		// we check that the value is still greater than 0
		if tx.Value().Cmp(value) != -1 || tx.Value().Cmp(big.NewInt(0)) != 1 {
			printTxInfo()
			panic(fmt.Sprintf("expected tx value %s; should be non-null and lower than %s", tx.Value().String(), value.String()))
		}

		sm.Logger.Info("REVERT tx receipt: %d", receipt.Status)
		sm.Logger.Info("  tx hash: %s", receipt.TxHash.String())
		sm.Logger.Info("  to: %s", tx.To().String())
		sm.Logger.Info("  value: %s", tx.Value().String())
		sm.Logger.Info("  block num: %d", receipt.BlockNumber)
	}()
}

// TestDepositEtherLiquidityCap tests depositing Ethers in a context where a liquidity cap is set
func TestDepositEtherLiquidityCap(sm *runner.SmokeTestRunner) {
	supply, err := sm.ETHZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	// Set a liquidity cap slightly above the current supply
	sm.Logger.Info("Setting a liquidity cap")
	liquidityCap := math.NewUintFromBigInt(supply).Add(math.NewUint(1e16))
	msg := fungibletypes.NewMsgUpdateZRC20LiquidityCap(
		sm.ZetaTxServer.GetAccountAddress(0),
		sm.ETHZRC20Addr.Hex(),
		liquidityCap,
	)
	res, err := sm.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("set liquidity cap tx hash: %s", res.TxHash)

	sm.Logger.Info("Depositing more than liquidity cap should make cctx reverted")
	signedTx, err := sm.SendEther(sm.TSSAddress, big.NewInt(1e17), nil)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, signedTx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposit eth tx failed")
	}
	cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, signedTx.Hash().Hex(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
		panic(fmt.Sprintf("expected cctx status to be Reverted; got %s", cctx.CctxStatus.Status))
	}
	sm.Logger.Info("CCTX has been reverted")

	sm.Logger.Info("Depositing less than liquidity cap should still succeed")
	initialBal, err := sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	signedTx, err = sm.SendEther(sm.TSSAddress, big.NewInt(1e15), nil)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, signedTx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposit eth tx failed")
	}
	utils.WaitCctxMinedByInTxHash(sm.Ctx, signedTx.Hash().Hex(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	expectedBalance := big.NewInt(0).Add(initialBal, big.NewInt(1e15))

	bal, err := sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	if bal.Cmp(expectedBalance) != 0 {
		panic(fmt.Sprintf("expected balance to be %s; got %s", expectedBalance.String(), bal.String()))
	}
	sm.Logger.Info("Deposit succeeded")

	sm.Logger.Info("Removing the liquidity cap")
	msg = fungibletypes.NewMsgUpdateZRC20LiquidityCap(
		sm.ZetaTxServer.GetAccountAddress(0),
		sm.ETHZRC20Addr.Hex(),
		math.ZeroUint(),
	)
	res, err = sm.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("remove liquidity cap tx hash: %s", res.TxHash)
	initialBal, err = sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	signedTx, err = sm.SendEther(sm.TSSAddress, big.NewInt(1e17), nil)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, signedTx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposit eth tx failed")
	}
	utils.WaitCctxMinedByInTxHash(sm.Ctx, signedTx.Hash().Hex(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	expectedBalance = big.NewInt(0).Add(initialBal, big.NewInt(1e17))

	bal, err = sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	if bal.Cmp(expectedBalance) != 0 {
		panic(fmt.Sprintf("expected balance to be %s; got %s", expectedBalance.String(), bal.String()))
	}
	sm.Logger.Info("New deposit succeeded")
}
