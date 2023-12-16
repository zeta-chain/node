package smoketests

import (
	"context"
	"fmt"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/connectorzevm.sol"
	wzeta "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/wzeta.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

func TestSendZetaOut(sm *runner.SmokeTestRunner) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	zevmClient := sm.ZevmClient
	cctxClient := sm.CctxClient

	utils.LoudPrintf("Step 4: Sending ZETA from ZEVM to Ethereum\n")
	ConnectorZEVMAddr := ethcommon.HexToAddress("0x239e96c8f17C85c30100AC26F635Ea15f23E9c67")
	ConnectorZEVM, err := connectorzevm.NewZetaConnectorZEVM(ConnectorZEVMAddr, zevmClient)
	if err != nil {
		panic(err)
	}

	wzetaAddr := ethcommon.HexToAddress("0x5F0b1a82749cb4E2278EC87F8BF6B618dC71a8bf")
	wZeta, err := wzeta.NewWETH9(wzetaAddr, zevmClient)
	if err != nil {
		panic(err)
	}
	zchainid, err := zevmClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("zevm chainid: %d\n", zchainid)

	// 10 Zeta
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(10))

	zauth := sm.ZevmAuth
	zauth.Value = amount
	tx, err := wZeta.Deposit(zauth)
	if err != nil {
		panic(err)
	}
	zauth.Value = big.NewInt(0)

	fmt.Printf("Deposit tx hash: %s\n", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("Deposit tx receipt: status %d\n", receipt.Status)

	tx, err = wZeta.Approve(zauth, ConnectorZEVMAddr, amount)
	if err != nil {
		panic(err)
	}
	fmt.Printf("wzeta.approve tx hash: %s\n", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("approve tx receipt: status %d\n", receipt.Status)
	tx, err = ConnectorZEVM.Send(zauth, connectorzevm.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(1337),
		DestinationAddress:  sm.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("send tx hash: %s\n", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("send tx receipt: status %d\n", receipt.Status)
	fmt.Printf("  Logs:\n")
	for _, log := range receipt.Logs {
		sentLog, err := ConnectorZEVM.ParseZetaSent(*log)
		if err == nil {
			fmt.Printf("    Dest Addr: %s\n", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			fmt.Printf("    Dest Chain: %d\n", sentLog.DestinationChainId)
			fmt.Printf("    Dest Gas: %d\n", sentLog.DestinationGasLimit)
			fmt.Printf("    Zeta Value: %d\n", sentLog.ZetaValueAndGas)
		}
	}
	fmt.Printf("waiting for cctx status to change to final...\n")

	sm.WG.Add(1)
	go func() {
		defer sm.WG.Done()
		cctx := utils.WaitCctxMinedByInTxHash(tx.Hash().Hex(), cctxClient)
		receipt, err := sm.GoerliClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
		if err != nil {
			panic(err)
		}
		for _, log := range receipt.Logs {
			event, err := sm.ConnectorEth.ParseZetaReceived(*log)
			if err == nil {
				fmt.Printf("    Dest Addr: %s\n", event.DestinationAddress.Hex())
				fmt.Printf("    sender addr: %x\n", event.ZetaTxSenderAddress)
				fmt.Printf("    Zeta Value: %d\n", event.ZetaValue)
				if event.ZetaValue.Cmp(amount) != -1 {
					panic("wrong zeta value, gas should be paid in the amount")
				}
			}
		}
	}()
	sm.WG.Wait()
}

func TestSendZetaOutBTCRevert(sm *runner.SmokeTestRunner) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	zevmClient := sm.ZevmClient

	utils.LoudPrintf("Step 5: Should revert when sending ZETA from ZEVM to Bitcoin\n")
	ConnectorZEVMAddr := ethcommon.HexToAddress("0x239e96c8f17C85c30100AC26F635Ea15f23E9c67")
	ConnectorZEVM, err := connectorzevm.NewZetaConnectorZEVM(ConnectorZEVMAddr, zevmClient)
	if err != nil {
		panic(err)
	}

	wzetaAddr := ethcommon.HexToAddress("0x5F0b1a82749cb4E2278EC87F8BF6B618dC71a8bf")
	wZeta, err := wzeta.NewWETH9(wzetaAddr, zevmClient)
	if err != nil {
		panic(err)
	}
	zchainid, err := zevmClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("zevm chainid: %d\n", zchainid)

	zauth := sm.ZevmAuth
	zauth.Value = big.NewInt(1e18)
	tx, err := wZeta.Deposit(zauth)
	if err != nil {
		panic(err)
	}
	zauth.Value = big.NewInt(0)

	fmt.Printf("Deposit tx hash: %s\n", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("Deposit tx receipt: status %d\n", receipt.Status)

	tx, err = wZeta.Approve(zauth, ConnectorZEVMAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	fmt.Printf("wzeta.approve tx hash: %s\n", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("approve tx receipt: status %d\n", receipt.Status)
	tx, err = ConnectorZEVM.Send(zauth, connectorzevm.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(common.BtcRegtestChain().ChainId),
		DestinationAddress:  sm.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     big.NewInt(1e17),
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("send tx hash: %s\n", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("send tx receipt: status %d\n", receipt.Status)
	if receipt.Status != 0 {
		panic("Was able to send ZETA to BTC")
	}
}
