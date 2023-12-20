package runner

import (
	"context"
	"fmt"
	"github.com/zeta-chain/zetacore/common"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// DepositZeta deposits ZETA on ZetaChain from the ZETA smart contract on EVM
func (sm *SmokeTestRunner) DepositZeta() {
	sm.Logger.Print("⏳ depositing ZETA into ZEVM")
	startTime := time.Now()
	defer func() {
		sm.Logger.Print("✅ ZETA deposited in %s", time.Since(startTime))
	}()

	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(100)) // 100 Zeta
	sm.Logger.InfoLoud("Sending ZETA to ZetaChain\n")
	tx, err := sm.ZetaEth.Approve(sm.GoerliAuth, sm.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("Approve tx receipt: status %d", receipt.Status)

	tx, err = sm.ConnectorEth.Send(sm.GoerliAuth, zetaconnectoreth.ZetaInterfacesSendInput{
		// TODO: allow user to specify destination chain id
		// https://github.com/zeta-chain/node-private/issues/41
		DestinationChainId:  big.NewInt(common.ZetaPrivnetChain().ChainId),
		DestinationAddress:  sm.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("Send tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	if receipt.Status != 1 {
		panic(fmt.Sprintf("expected tx receipt status to be 1; got %d", receipt.Status))
	}
	sm.Logger.Info("Send tx receipt: status %d", receipt.Status)
	sm.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := sm.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			sm.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			sm.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			sm.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			sm.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
			sm.Logger.Info("    Block Num: %d", log.BlockNumber)
		}
	}

	// wait for cctx to be mined and check balance
	cctx := utils.WaitCctxMinedByInTxHash(tx.Hash().Hex(), sm.CctxClient, sm.Logger)
	if cctx.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf("expected cctx status to be mined; got %s, message: %s",
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage),
		)
	}

	bn, err := sm.ZevmClient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	// #nosec G701 smoketest - always in range
	bal, err := sm.ZevmClient.BalanceAt(context.Background(), sm.DeployerAddress, big.NewInt(int64(bn)))
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Zeta block %d, Deployer %s Zeta balance: %d", bn, sm.DeployerAddress.Hex(), bal)

	if bal.Int64() == 0 {
		panic("Deployer ZETA balance is 0 after deposit")
	}
}
