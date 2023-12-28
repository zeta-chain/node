package smoketests

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/bitcoin"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient"
)

var blockHeaderBTCTimeout = 5 * time.Minute

func TestBTCMerkelProof(sm *runner.SmokeTestRunner) {
	// mine blocks
	stop := sm.MineBlocks()

	utxos, err := sm.BtcRPCClient.ListUnspent()
	if err != nil {
		panic(err)
	}
	spendableAmount := 0.0
	spendableUTXOs := 0
	for _, utxo := range utxos {
		if utxo.Spendable {
			spendableAmount += utxo.Amount
			spendableUTXOs++
		}
	}

	if spendableAmount < 1.1 {
		panic(fmt.Errorf("not enough spendable BTC to run the test; have %f", spendableAmount))
	}
	if spendableUTXOs < 2 {
		panic(fmt.Errorf("not enough spendable BTC UTXOs to run the test; have %d", spendableUTXOs))
	}

	sm.Logger.Info("ListUnspent:")
	sm.Logger.Info("  spendableAmount: %f", spendableAmount)
	sm.Logger.Info("  spendableUTXOs: %d", spendableUTXOs)
	sm.Logger.Info("Now sending two txs to TSS address...")

	// send two transactions to the TSS address
	txHash, err := sm.SendToTSSFromDeployerToDeposit(
		sm.BTCTSSAddress,
		1.1+zetaclient.BtcDepositorFeeMin,
		utxos[:2],
		sm.BtcRPCClient,
		sm.BTCDeployerAddress,
	)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("BTC tx sent: %s", txHash.String())

	// check that the tx is in the block header
	proveBTCTransaction(sm, txHash)

	// stop mining
	stop <- struct{}{}
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
		if time.Since(startTime) > blockHeaderBTCTimeout {
			panic("timed out waiting for block header to show up in observer")
		}

		_, err := sm.ObserverClient.GetBlockHeaderByHash(sm.Ctx, &observertypes.QueryGetBlockHeaderByHashRequest{
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
	res, err := sm.ObserverClient.Prove(sm.Ctx, &observertypes.QueryProveRequest{
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
