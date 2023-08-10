//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/x/crosschain/types"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// wait until cctx is mined; returns the cctxIndex
func WaitCctxMinedByInTxHash(inTxHash string, cctxClient types.QueryClient) *types.CrossChainTx {
	var cctxIndex string
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf("Waiting for cctx to be mined by inTxHash: %s\n", inTxHash)
		res, err := cctxClient.InTxHashToCctx(context.Background(), &types.QueryGetInTxHashToCctxRequest{InTxHash: inTxHash})
		if err != nil {
			continue
		}
		cctxIndex = res.InTxHashToCctx.CctxIndex
		fmt.Printf("Deposit receipt cctx index: %s\n", cctxIndex)
		break
	}
	for {
		time.Sleep(5 * time.Second)
		{
			res, err := cctxClient.OutTxTrackerAll(context.Background(), &types.QueryAllOutTxTrackerRequest{})
			if err != nil {
				fmt.Printf("OutTxTrackerAll err: %s\n", err.Error())
				continue
			}
			for _, tracker := range res.OutTxTracker {
				for _, hash := range tracker.HashList {
					fmt.Printf("OutTxTracker: %s proved %v\n", hash.TxHash, hash.Proved)
				}
			}
		}
		res, err := cctxClient.Cctx(context.Background(), &types.QueryGetCctxRequest{Index: cctxIndex})
		if err == nil && IsTerminalStatus(res.CrossChainTx.CctxStatus.Status) {
			fmt.Printf("Deposit receipt cctx status: %+v; The cctx is processed\n", res.CrossChainTx.CctxStatus.Status.String())
			return res.CrossChainTx
		}

	}
}

func IsTerminalStatus(status types.CctxStatus) bool {
	return status == types.CctxStatus_OutboundMined || status == types.CctxStatus_Aborted || status == types.CctxStatus_Reverted
}

func LoudPrintf(format string, a ...any) {
	fmt.Println("=======================================")
	fmt.Printf(format, a...)
	fmt.Println("=======================================")
}

func CheckNonce(client *ethclient.Client, addr ethcommon.Address, expectedNonce uint64) error {
	nonce, err := client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return err
	}
	if nonce != expectedNonce {
		return fmt.Errorf("want nonce %d; got %d", expectedNonce, nonce)
	}
	return nil
}

// wait until a broadcasted tx to be mined and return its receipt
// timeout and panic after 30s.
func MustWaitForTxReceipt(client *ethclient.Client, tx *ethtypes.Transaction) *ethtypes.Receipt {
	start := time.Now()
	for {
		if time.Since(start) > 30*time.Second {
			panic("waiting tx receipt timeout")
		}
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			continue
		}
		if receipt != nil {
			return receipt
		}
		time.Sleep(1 * time.Second)
	}
}

// scriptPK is a hex string for P2WPKH script
func ScriptPKToAddress(scriptPKHex string) string {
	pkh, err := hex.DecodeString(scriptPKHex[4:])
	if err == nil {
		addr, err := btcutil.NewAddressWitnessPubKeyHash(pkh, &chaincfg.RegressionNetParams)
		if err == nil {
			return addr.EncodeAddress()
		}
	}
	return ""
}
