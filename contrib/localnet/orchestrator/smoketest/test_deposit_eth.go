//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	zrc20 "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/ethereum"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient"
)

// this tests sending ZETA out of ZetaChain to Ethereum
func (sm *SmokeTest) TestDepositEtherIntoZRC20() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	goerliClient := sm.goerliClient
	LoudPrintf("Deposit Ether into ZEVM\n")
	bn, err := goerliClient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI block number: %d\n", bn)
	bal, err := goerliClient.BalanceAt(context.Background(), DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI deployer balance: %s\n", bal.String())

	systemContract := sm.SystemContract
	if err != nil {
		panic(err)
	}
	ethZRC20Addr, err := systemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.GoerliChain().ChainId))
	if err != nil {
		panic(err)
	}
	sm.ETHZRC20Addr = ethZRC20Addr
	fmt.Printf("eth zrc20 address: %s\n", ethZRC20Addr.String())
	ethZRC20, err := zrc20.NewZRC20(ethZRC20Addr, sm.zevmClient)
	if err != nil {
		panic(err)
	}
	sm.ETHZRC20 = ethZRC20
	initialBalance, err := ethZRC20.BalanceOf(nil, DeployerAddress)
	if err != nil {
		panic(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	signedTx, err := sm.SendEther(TSSAddress, value, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("GOERLI tx sent: %s; to %s, nonce %d\n", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := MustWaitForTxReceipt(sm.goerliClient, signedTx)
	fmt.Printf("GOERLI tx receipt: %d\n", receipt.Status)
	fmt.Printf("  tx hash: %s\n", receipt.TxHash.String())
	fmt.Printf("  to: %s\n", signedTx.To().String())
	fmt.Printf("  value: %d\n", signedTx.Value())
	fmt.Printf("  block num: %d\n", receipt.BlockNumber)

	{
		LoudPrintf("Merkle Proof\n")
		txHash := receipt.TxHash
		blockHash := receipt.BlockHash
		txIndex := int(receipt.TransactionIndex)

		block, err := sm.goerliClient.BlockByHash(context.Background(), blockHash)
		if err != nil {
			panic(err)
		}
		i := 0
		for {
			if i > 20 {
				panic("block header not found")
			}
			_, err := sm.observerClient.GetBlockHeaderByHash(context.Background(), &observertypes.QueryGetBlockHeaderByHashRequest{
				BlockHash: blockHash.Bytes(),
			})
			if err != nil {
				fmt.Printf("WARN: block header not found; retrying...\n")
				time.Sleep(2 * time.Second)
			} else {
				fmt.Printf("OK: block header found\n")
				break
			}
			i++
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
		res, err := sm.observerClient.Prove(context.Background(), &observertypes.QueryProveRequest{
			BlockHash: blockHash.Hex(),
			TxIndex:   int64(txIndex),
			TxHash:    txHash.Hex(),
			Proof:     common.NewEthereumProof(*txProof),
			ChainId:   0,
		})
		if err != nil {
			panic(err)
		}
		if !res.Valid {
			panic("txProof invalid") // FIXME: don't do this in production
		}
		fmt.Printf("OK: txProof verified\n")
	}

	{
		tx, err := sm.SendEther(TSSAddress, big.NewInt(101000000000000000), []byte(zetaclient.DonationMessage))
		if err != nil {
			panic(err)
		}
		receipt := MustWaitForTxReceipt(sm.goerliClient, tx)
		fmt.Printf("GOERLI donation tx receipt: %d\n", receipt.Status)
	}

	c := make(chan any)
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		WaitCctxMinedByInTxHash(signedTx.Hash().Hex(), sm.cctxClient)
		c <- 0
	}()
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		<-c

		currentBalance, err := ethZRC20.BalanceOf(nil, DeployerAddress)
		if err != nil {
			panic(err)
		}
		diff := big.NewInt(0)
		diff.Sub(currentBalance, initialBalance)
		fmt.Printf("eth zrc20 balance: %s\n", currentBalance.String())
		if diff.Cmp(value) != 0 {
			fmt.Printf("eth zrc20 bal wanted %d, got %d\n", value, diff)
			panic("bal mismatch")
		}

	}()
	sm.wg.Wait()
}

func (sm *SmokeTest) TestDepositAndCallRefund() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	LoudPrintf("Deposit ZRC20 into ZEVM and call a contract that reverts; should refund\n")

	goerliClient := sm.goerliClient
	bn, err := goerliClient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI block number: %d\n", bn)
	bal, err := goerliClient.BalanceAt(context.Background(), DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI deployer balance: %s\n", bal.String())
	nonce, err := goerliClient.PendingNonceAt(context.Background(), DeployerAddress)
	if err != nil {
		panic(err)
	}
	value := big.NewInt(100000000000000000) // in wei (1 eth)
	gasLimit := uint64(23000)               // in units
	gasPrice, err := goerliClient.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	data := append(sm.BTCZRC20Addr.Bytes(), []byte("hello sailors")...) // this data
	tx := ethtypes.NewTransaction(nonce, TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := goerliClient.NetworkID(context.Background())
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		panic(err)
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		panic(err)
	}
	err = goerliClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI tx sent: %s; to %s, nonce %d\n", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := MustWaitForTxReceipt(sm.goerliClient, signedTx)
	fmt.Printf("GOERLI tx receipt: %d\n", receipt.Status)
	fmt.Printf("  tx hash: %s\n", receipt.TxHash.String())
	fmt.Printf("  to: %s\n", signedTx.To().String())
	fmt.Printf("  value: %d\n", signedTx.Value())
	fmt.Printf("  block num: %d\n", receipt.BlockNumber)

	func() {
		cctx := WaitCctxMinedByInTxHash(signedTx.Hash().Hex(), sm.cctxClient)
		fmt.Printf("cctx status message: %s", cctx.CctxStatus.StatusMessage)
		revertTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
		fmt.Printf("GOERLI revert tx receipt: status %d\n", receipt.Status)
		tx, _, err := sm.goerliClient.TransactionByHash(context.Background(), ethcommon.HexToHash(revertTxHash))
		if err != nil {
			panic(err)
		}
		receipt, err := sm.goerliClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(revertTxHash))
		if err != nil {
			panic(err)
		}
		if cctx.CctxStatus.Status != types.CctxStatus_Reverted || receipt.Status == 0 || *tx.To() != DeployerAddress || tx.Value().Cmp(value) != 0 {
			// debug info when test fails
			fmt.Printf("  tx: %+v\n", tx)
			fmt.Printf("  receipt: %+v\n", receipt)
			fmt.Printf("cctx http://localhost:1317/zeta-chain/crosschain/cctx/%s\n", cctx.Index)
			panic(fmt.Sprintf("expected cctx status PendingRevert; got %s", cctx.CctxStatus.Status))
		}
	}()
}

func (sm *SmokeTest) SendEther(to ethcommon.Address, value *big.Int, data []byte) (*ethtypes.Transaction, error) {
	goerliClient := sm.goerliClient

	nonce, err := goerliClient.PendingNonceAt(context.Background(), DeployerAddress)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(30000) // in units
	gasPrice, err := goerliClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	tx := ethtypes.NewTransaction(nonce, TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := goerliClient.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	deployerPrivkey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		return nil, err
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		return nil, err
	}
	err = goerliClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
