package e2etests

import (
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

// TestContextUpgrade tests sending ETH on ZetaChain and check context data
func TestContextUpgrade(sm *runner.E2ERunner) {
	value := big.NewInt(1000000000000000) // in wei (1 eth)
	data := make([]byte, 0, 32)
	data = append(data, sm.ContextAppAddr.Bytes()...)
	data = append(data, []byte("filler")...) // just to make sure that this is a contract call;

	signedTx, err := sm.SendEther(sm.TSSAddress, value, data)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("GOERLI tx sent: %s; to %s, nonce %d", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, signedTx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
	sm.Logger.Info("GOERLI tx receipt: %d", receipt.Status)
	sm.Logger.Info("  tx hash: %s", receipt.TxHash.String())
	sm.Logger.Info("  to: %s", signedTx.To().String())
	sm.Logger.Info("  value: %d", signedTx.Value())
	sm.Logger.Info("  block num: %d", receipt.BlockNumber)
	sm.Logger.Info("  data: %x", signedTx.Data())

	found := false
	for i := 0; i < 10; i++ {
		eventIter, err := sm.ContextApp.FilterContextData(&bind.FilterOpts{
			Start: 0,
			End:   nil,
		})
		if err != nil {
			sm.Logger.Info("filter error: %s", err.Error())
			continue
		}
		for eventIter.Next() {
			sm.Logger.Info("event: ContextData")
			sm.Logger.Info("  origin: %x", eventIter.Event.Origin)
			sm.Logger.Info("  sender: %s", eventIter.Event.Sender.Hex())
			sm.Logger.Info("  chainid: %d", eventIter.Event.ChainID)
			sm.Logger.Info("  msgsender: %s", eventIter.Event.MsgSender.Hex())
			found = true
			if bytes.Compare(eventIter.Event.Origin, sm.DeployerAddress.Bytes()) != 0 {
				panic("origin mismatch")
			}
			chainID, err := sm.GoerliClient.ChainID(sm.Ctx)
			if err != nil {
				panic(err)
			}
			if eventIter.Event.ChainID.Cmp(chainID) != 0 {
				panic("chainID mismatch")
			}

		}
		if found {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if !found {
		panic("event not found")
	}

}
