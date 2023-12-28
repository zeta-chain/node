package smoketests

import (
	"context"
	"math/big"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/ethereum"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

var blockHeaderETHTimeout = 5 * time.Minute

func TestEthereumMerkelProof(sm *runner.SmokeTestRunner) {
	// send eth to TSS address
	value := big.NewInt(100000000000000000) // in wei (0.1 eth)
	signedTx, err := sm.SendEther(sm.TSSAddress, value, nil)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("GOERLI tx sent: %s; to %s, nonce %d", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, signedTx, sm.Logger)
	if receipt.Status == 0 {
		panic("deposit failed")
	}

	// check that the tx is in the block header
	proveEthTransaction(sm, receipt)
}

func proveEthTransaction(sm *runner.SmokeTestRunner, receipt *ethtypes.Receipt) {
	startTime := time.Now()

	txHash := receipt.TxHash
	blockHash := receipt.BlockHash

	// #nosec G701 smoketest - always in range
	txIndex := int(receipt.TransactionIndex)

	block, err := sm.GoerliClient.BlockByHash(context.Background(), blockHash)
	if err != nil {
		panic(err)
	}
	for {
		// check timeout
		if time.Since(startTime) > blockHeaderETHTimeout {
			panic("timeout waiting for block header")
		}

		_, err := sm.ObserverClient.GetBlockHeaderByHash(context.Background(), &observertypes.QueryGetBlockHeaderByHashRequest{
			BlockHash: blockHash.Bytes(),
		})
		if err != nil {
			sm.Logger.Info("WARN: block header not found; retrying... error: %s", err.Error())
		} else {
			sm.Logger.Info("OK: block header found")
			break
		}

		time.Sleep(2 * time.Second)
	}

	trie := ethereum.NewTrie(block.Transactions())
	if trie.Hash() != block.Header().TxHash {
		panic("tx root hash & block tx root mismatch")
	}
	txProof, err := trie.GenerateProof(txIndex)
	if err != nil {
		panic("error generating txProof")
	}
	val, err := txProof.Verify(block.TxHash(), txIndex)
	if err != nil {
		panic("error verifying txProof")
	}
	var txx ethtypes.Transaction
	err = txx.UnmarshalBinary(val)
	if err != nil {
		panic("error unmarshalling txProof'd tx")
	}
	res, err := sm.ObserverClient.Prove(context.Background(), &observertypes.QueryProveRequest{
		BlockHash: blockHash.Hex(),
		TxIndex:   int64(txIndex),
		TxHash:    txHash.Hex(),
		Proof:     common.NewEthereumProof(txProof),
		ChainId:   common.GoerliLocalnetChain().ChainId,
	})
	if err != nil {
		panic(err)
	}
	if !res.Valid {
		panic("txProof invalid") // FIXME: don't do this in production
	}
	sm.Logger.Info("OK: txProof verified")
}
