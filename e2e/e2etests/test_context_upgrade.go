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
func TestContextUpgrade(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestContextUpgrade requires exactly one argument for the value.")
	}

	// parse the value from the provided arguments
	value, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid value specified for TestContextUpgrade.")
	}

	data := make([]byte, 0, 32)
	data = append(data, r.ContextAppAddr.Bytes()...)
	data = append(data, []byte("filler")...) // just to make sure that this is a contract call;

	signedTx, err := r.SendEther(r.TSSAddress, value, data)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("EVM tx sent: %s; to %s, nonce %d", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
	r.Logger.Info("EVM tx receipt: %d", receipt.Status)
	r.Logger.Info("  tx hash: %s", receipt.TxHash.String())
	r.Logger.Info("  to: %s", signedTx.To().String())
	r.Logger.Info("  value: %d", signedTx.Value())
	r.Logger.Info("  block num: %d", receipt.BlockNumber)
	r.Logger.Info("  data: %x", signedTx.Data())

	found := false
	for i := 0; i < 10; i++ {
		eventIter, err := r.ContextApp.FilterContextData(&bind.FilterOpts{
			Start: 0,
			End:   nil,
		})
		if err != nil {
			r.Logger.Info("filter error: %s", err.Error())
			continue
		}
		for eventIter.Next() {
			r.Logger.Info("event: ContextData")
			r.Logger.Info("  origin: %x", eventIter.Event.Origin)
			r.Logger.Info("  sender: %s", eventIter.Event.Sender.Hex())
			r.Logger.Info("  chainid: %d", eventIter.Event.ChainID)
			r.Logger.Info("  msgsender: %s", eventIter.Event.MsgSender.Hex())
			found = true
			if bytes.Compare(eventIter.Event.Origin, r.DeployerAddress.Bytes()) != 0 {
				panic("origin mismatch")
			}
			chainID, err := r.EVMClient.ChainID(r.Ctx)
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
