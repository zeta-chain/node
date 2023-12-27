package smoketests

import (
	"context"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/bitcoin"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"time"
)

var blockHeaderTimeout = 30 * time.Second

func TestBTCMerkelProof(sm *runner.SmokeTestRunner) {

}

func proveBTCTransaction(sm *runner.SmokeTestRunner, txHash *chainhash.Hash) {
	// get tx result
	btc := sm.BtcRPCClient
	txResult, err := btc.GetTransaction(txHash)
	if err != nil {
		panic("should get outTx result")
	}
	if txResult.Confirmations <= 0 {
		panic("outTx should have already confirmed")
	}
	txBytes, err := hex.DecodeString(txResult.Hex)
	if err != nil {
		panic(err)
	}

	// get the block with verbose transactions
	blockHash, err := chainhash.NewHashFromStr(txResult.BlockHash)
	if err != nil {
		panic(err)
	}
	blockVerbose, err := btc.GetBlockVerboseTx(blockHash)
	if err != nil {
		panic("should get block verbose tx")
	}

	// get the block header
	header, err := btc.GetBlockHeader(blockHash)
	if err != nil {
		panic("should get block header")
	}

	// collect all the txs in the block
	txns := []*btcutil.Tx{}
	for _, res := range blockVerbose.Tx {
		txBytes, err := hex.DecodeString(res.Hex)
		if err != nil {
			panic(err)
		}
		tx, err := btcutil.NewTxFromBytes(txBytes)
		if err != nil {
			panic(err)
		}
		txns = append(txns, tx)
	}

	// build merkle proof
	mk := bitcoin.NewMerkle(txns)
	path, index, err := mk.BuildMerkleProof(int(txResult.BlockIndex))
	if err != nil {
		panic("should build merkle proof")
	}

	// verify merkle proof statically
	pass := bitcoin.Prove(*txHash, header.MerkleRoot, path, index)
	if !pass {
		panic("should verify merkle proof")
	}

	// wait for block header to show up in ZetaChain
	startTime := time.Now()
	hash := header.BlockHash()
	for {
		// timeout
		if time.Since(startTime) > blockHeaderTimeout {
			panic("timed out waiting for block header to show up in observer")
		}

		_, err := sm.ObserverClient.GetBlockHeaderByHash(context.Background(), &observertypes.QueryGetBlockHeaderByHashRequest{
			BlockHash: hash.CloneBytes(),
		})
		if err != nil {
			sm.Logger.Info("waiting for block header to show up in observer... current hash %s; err %s", hash.String(), err.Error())
		}
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// verify merkle proof through RPC
	res, err := sm.ObserverClient.Prove(context.Background(), &observertypes.QueryProveRequest{
		ChainId:   common.BtcRegtestChain().ChainId,
		TxHash:    txHash.String(),
		BlockHash: blockHash.String(),
		Proof:     common.NewBitcoinProof(txBytes, path, index),
		TxIndex:   0, // bitcoin doesn't use txIndex
	})
	if err != nil {
		panic(err)
	}
	if !res.Valid {
		panic("txProof should be valid")
	}
	sm.Logger.Info("OK: txProof verified for inTx: %s", txHash.String())
}
