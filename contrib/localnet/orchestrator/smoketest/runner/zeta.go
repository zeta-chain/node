package runner

import (
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// WaitForMinedCCTX waits for a cctx to be mined from a tx
func (sm *SmokeTestRunner) WaitForMinedCCTX(txHash ethcommon.Hash) {
	defer func() {
		sm.Unlock()
	}()
	sm.Lock()

	cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, txHash.Hex(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	if cctx.CctxStatus.Status != types.CctxStatus_OutboundMined {
		panic(fmt.Sprintf("expected cctx status to be mined; got %s, message: %s",
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage),
		)
	}
}

// SendZetaOnEvm sends ZETA to an address on EVM
// this allows the ZETA contract deployer to funds other accounts on EVM
func (sm *SmokeTestRunner) SendZetaOnEvm(address ethcommon.Address, zetaAmount int64) *ethtypes.Transaction {
	// the deployer might be sending ZETA in different goroutines
	defer func() {
		sm.Unlock()
	}()
	sm.Lock()

	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(zetaAmount))
	tx, err := sm.ZetaEth.Transfer(sm.GoerliAuth, address, amount)
	if err != nil {
		panic(err)
	}
	return tx
}

// DepositZeta deposits ZETA on ZetaChain from the ZETA smart contract on EVM
func (sm *SmokeTestRunner) DepositZeta() ethcommon.Hash {
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(100)) // 100 Zeta

	return sm.DepositZetaWithAmount(amount)
}

// DepositZetaWithAmount deposits ZETA on ZetaChain from the ZETA smart contract on EVM with the specified amount
func (sm *SmokeTestRunner) DepositZetaWithAmount(amount *big.Int) ethcommon.Hash {
	tx, err := sm.ZetaEth.Approve(sm.GoerliAuth, sm.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	sm.Logger.EVMReceipt(*receipt, "approve")
	if receipt.Status != 1 {
		panic("approve tx failed")
	}

	tx, err = sm.ConnectorEth.Send(sm.GoerliAuth, zetaconnectoreth.ZetaInterfacesSendInput{
		// TODO: allow user to specify destination chain id
		// https://github.com/zeta-chain/node-private/issues/41
		DestinationChainId:  big.NewInt(common.ZetaChainMainnet().ChainId),
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

	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	sm.Logger.EVMReceipt(*receipt, "send")
	if receipt.Status != 1 {
		panic(fmt.Sprintf("expected tx receipt status to be 1; got %d", receipt.Status))
	}

	sm.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := sm.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			sm.Logger.Info("    Connector: %s", sm.ConnectorEthAddr.String())
			sm.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			sm.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			sm.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			sm.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
			sm.Logger.Info("    Block Num: %d", log.BlockNumber)
		}
	}

	return tx.Hash()
}
