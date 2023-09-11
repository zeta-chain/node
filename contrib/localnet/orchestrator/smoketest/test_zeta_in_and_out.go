//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/zeta-chain/zetacore/common"

	ethcommon "github.com/ethereum/go-ethereum/common"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/connectorzevm.sol"
	wzeta "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/wzeta.sol"
)

func (sm *SmokeTest) TestSendZetaIn() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	// ==================== Sending ZETA to ZetaChain ===================
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(100)) // 100 Zeta
	LoudPrintf("Step 3: Sending ZETA to ZetaChain\n")
	tx, err := sm.ZetaEth.Approve(sm.goerliAuth, sm.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Approve tx hash: %s\n", tx.Hash().Hex())
	receipt := MustWaitForTxReceipt(sm.goerliClient, tx)
	fmt.Printf("Approve tx receipt: status %d\n", receipt.Status)
	tx, err = sm.ConnectorEth.Send(sm.goerliAuth, zetaconnectoreth.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(101), // in dev mode, 101 is the  zEVM ChainID
		DestinationAddress:  DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Send tx hash: %s\n", tx.Hash().Hex())
	receipt = MustWaitForTxReceipt(sm.goerliClient, tx)
	fmt.Printf("Send tx receipt: status %d\n", receipt.Status)
	fmt.Printf("  Logs:\n")
	for _, log := range receipt.Logs {
		sentLog, err := sm.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			fmt.Printf("    Dest Addr: %s\n", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			fmt.Printf("    Dest Chain: %d\n", sentLog.DestinationChainId)
			fmt.Printf("    Dest Gas: %d\n", sentLog.DestinationGasLimit)
			fmt.Printf("    Zeta Value: %d\n", sentLog.ZetaValueAndGas)
			fmt.Printf("    Block Num: %d\n", log.BlockNumber)
		}
	}

	sm.wg.Add(1)
	go func() {
		bn, _ := sm.zevmClient.BlockNumber(context.Background())
		initialBal, _ := sm.zevmClient.BalanceAt(context.Background(), DeployerAddress, big.NewInt(int64(bn)))
		fmt.Printf("Zeta block %d, Initial Deployer Zeta balance: %d\n", bn, initialBal)

		defer sm.wg.Done()
		for {
			time.Sleep(5 * time.Second)
			bn, _ = sm.zevmClient.BlockNumber(context.Background())
			bal, _ := sm.zevmClient.BalanceAt(context.Background(), DeployerAddress, big.NewInt(int64(bn)))
			fmt.Printf("Zeta block %d, Deployer Zeta balance: %d\n", bn, bal)

			diff := big.NewInt(0)
			diff.Sub(bal, initialBal)

			if diff.Cmp(amount) == 0 {
				fmt.Printf("Expected zeta balance; success!\n")
				break
			}
		}
	}()
	sm.wg.Wait()
}

func (sm *SmokeTest) TestSendZetaOut() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	zevmClient := sm.zevmClient
	cctxClient := sm.cctxClient

	LoudPrintf("Step 4: Sending ZETA from ZEVM to Ethereum\n")
	ConnectorZEVMAddr := ethcommon.HexToAddress("0x239e96c8f17C85c30100AC26F635Ea15f23E9c67")
	ConnectorZEVM, err := connectorzevm.NewZetaConnectorZEVM(ConnectorZEVMAddr, zevmClient)
	if err != nil {
		panic(err)
	}

	wzetaAddr := ethcommon.HexToAddress("0x5F0b1a82749cb4E2278EC87F8BF6B618dC71a8bf")
	wzeta, err := wzeta.NewWETH9(wzetaAddr, zevmClient)
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

	zauth := sm.zevmAuth
	zauth.Value = amount
	tx, err := wzeta.Deposit(zauth)
	if err != nil {
		panic(err)
	}
	zauth.Value = BigZero

	fmt.Printf("Deposit tx hash: %s\n", tx.Hash().Hex())
	receipt := MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("Deposit tx receipt: status %d\n", receipt.Status)

	tx, err = wzeta.Approve(zauth, ConnectorZEVMAddr, amount)
	if err != nil {
		panic(err)
	}
	fmt.Printf("wzeta.approve tx hash: %s\n", tx.Hash().Hex())
	receipt = MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("approve tx receipt: status %d\n", receipt.Status)
	tx, err = ConnectorZEVM.Send(zauth, connectorzevm.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(1337),
		DestinationAddress:  DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("send tx hash: %s\n", tx.Hash().Hex())
	receipt = MustWaitForTxReceipt(zevmClient, tx)
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

	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		cctx := WaitCctxMinedByInTxHash(tx.Hash().Hex(), cctxClient)
		receipt, err := sm.goerliClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
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
	sm.wg.Wait()
}

func (sm *SmokeTest) TestSendZetaOutBTCRevert() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	zevmClient := sm.zevmClient

	LoudPrintf("Step 5: Should revert when sending ZETA from ZEVM to Bitcoin\n")
	ConnectorZEVMAddr := ethcommon.HexToAddress("0x239e96c8f17C85c30100AC26F635Ea15f23E9c67")
	ConnectorZEVM, err := connectorzevm.NewZetaConnectorZEVM(ConnectorZEVMAddr, zevmClient)
	if err != nil {
		panic(err)
	}

	wzetaAddr := ethcommon.HexToAddress("0x5F0b1a82749cb4E2278EC87F8BF6B618dC71a8bf")
	wzeta, err := wzeta.NewWETH9(wzetaAddr, zevmClient)
	if err != nil {
		panic(err)
	}
	zchainid, err := zevmClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("zevm chainid: %d\n", zchainid)

	zauth := sm.zevmAuth
	zauth.Value = big.NewInt(1e18)
	tx, err := wzeta.Deposit(zauth)
	if err != nil {
		panic(err)
	}
	zauth.Value = BigZero

	fmt.Printf("Deposit tx hash: %s\n", tx.Hash().Hex())
	receipt := MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("Deposit tx receipt: status %d\n", receipt.Status)

	tx, err = wzeta.Approve(zauth, ConnectorZEVMAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	fmt.Printf("wzeta.approve tx hash: %s\n", tx.Hash().Hex())
	receipt = MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("approve tx receipt: status %d\n", receipt.Status)
	tx, err = ConnectorZEVM.Send(zauth, connectorzevm.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(common.BtcRegtestChain().ChainId),
		DestinationAddress:  DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     big.NewInt(1e17),
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("send tx hash: %s\n", tx.Hash().Hex())
	receipt = MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("send tx receipt: status %d\n", receipt.Status)
	if receipt.Status != 0 {
		panic("Was able to send ZETA to BTC")
	}
}
