package smoketests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/connectorzevm.sol"
	wzeta "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/wzeta.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

func TestSendZetaOut(sm *runner.SmokeTestRunner) {
	zevmClient := sm.ZevmClient
	cctxClient := sm.CctxClient

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
	zchainid, err := zevmClient.ChainID(sm.Ctx)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("zevm chainid: %d", zchainid)

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
	sm.Logger.Info("Deposit tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(sm.Ctx, zevmClient, tx, sm.Logger)
	sm.Logger.Info("Deposit tx receipt: status %d", receipt.Status)

	tx, err = wZeta.Approve(zauth, ConnectorZEVMAddr, amount)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("wzeta.approve tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, zevmClient, tx, sm.Logger)
	sm.Logger.Info("approve tx receipt: status %d", receipt.Status)
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
	sm.Logger.Info("send tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, zevmClient, tx, sm.Logger)
	sm.Logger.Info("send tx receipt: status %d", receipt.Status)
	sm.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := ConnectorZEVM.ParseZetaSent(*log)
		if err == nil {
			sm.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			sm.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			sm.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			sm.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
		}
	}
	sm.Logger.Info("waiting for cctx status to change to final...")

	sm.WG.Add(1)
	go func() {
		defer sm.WG.Done()
		cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, tx.Hash().Hex(), cctxClient, sm.Logger)
		receipt, err := sm.GoerliClient.TransactionReceipt(sm.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
		if err != nil {
			panic(err)
		}
		for _, log := range receipt.Logs {
			event, err := sm.ConnectorEth.ParseZetaReceived(*log)
			if err == nil {
				sm.Logger.Info("    Dest Addr: %s", event.DestinationAddress.Hex())
				sm.Logger.Info("    sender addr: %x", event.ZetaTxSenderAddress)
				sm.Logger.Info("    Zeta Value: %d", event.ZetaValue)
				if event.ZetaValue.Cmp(amount) != -1 {
					panic("wrong zeta value, gas should be paid in the amount")
				}
			}
		}
	}()
	sm.WG.Wait()
}

func TestSendZetaOutBTCRevert(sm *runner.SmokeTestRunner) {
	zevmClient := sm.ZevmClient

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
	zchainid, err := zevmClient.ChainID(sm.Ctx)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("zevm chainid: %d", zchainid)

	zauth := sm.ZevmAuth
	zauth.Value = big.NewInt(1e18)
	tx, err := wZeta.Deposit(zauth)
	if err != nil {
		panic(err)
	}
	zauth.Value = big.NewInt(0)

	sm.Logger.Info("Deposit tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, zevmClient, tx, sm.Logger)
	sm.Logger.Info("Deposit tx receipt: status %d", receipt.Status)

	tx, err = wZeta.Approve(zauth, ConnectorZEVMAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("wzeta.approve tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, zevmClient, tx, sm.Logger)
	sm.Logger.Info("approve tx receipt: status %d", receipt.Status)
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
	sm.Logger.Info("send tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, zevmClient, tx, sm.Logger)
	sm.Logger.Info("send tx receipt: status %d", receipt.Status)
	if receipt.Status != 0 {
		panic("Was able to send ZETA to BTC")
	}
}
