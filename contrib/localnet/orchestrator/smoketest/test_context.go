package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"time"
)

func (sm *SmokeTest) TestSystemContractUpdateAddress() {

}

// this tests sending ZETA out of ZetaChain to Ethereum
func (sm *SmokeTest) TestContextUpgrade() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	goerliClient := sm.goerliClient
	LoudPrintf("Test ContextApp\n")
	bn, err := goerliClient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI block number: %d\n", bn)

	value := big.NewInt(1000000000000000) // in wei (1 eth)
	data := make([]byte, 0, 32)
	data = append(data, sm.ContextAppAddr.Bytes()...)
	data = append(data, []byte("filler")...) // just to make sure that this is a contract call;

	signedTx, err := sm.SendEther(TSSAddress, value, data)
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
	fmt.Printf("  data: %x\n", signedTx.Data())

	found := false
	for i := 0; i < 10; i++ {
		eventIter, err := sm.ContextApp.FilterContextData(&bind.FilterOpts{
			Start: 0,
			End:   nil,
		})
		if err != nil {
			fmt.Printf("filter error: %s\n", err.Error())
			continue
		}
		for eventIter.Next() {
			fmt.Printf("event: ContextData\n")
			fmt.Printf("  origin: %x\n", eventIter.Event.Origin)
			fmt.Printf("  sender: %s\n", eventIter.Event.Sender.Hex())
			fmt.Printf("  chainid: %d\n", eventIter.Event.ChainID)
			fmt.Printf("  msgsender: %s\n", eventIter.Event.MsgSender.Hex())
			found = true
			if bytes.Compare(eventIter.Event.Origin, DeployerAddress.Bytes()) != 0 {
				panic("origin mismatch")
			}
			chainID, err := sm.goerliClient.ChainID(context.Background())
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
