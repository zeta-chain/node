package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

type Payload struct {
	sender       ethcommon.Address
	destChainID  string
	destContract string
	zetaAmount   *big.Int
	gasLimit     *big.Int
	message      []byte
	msgid        [32]byte
}

func (cl *ChainETHish) Listen() {
	log.Info().Msg(fmt.Sprintf("begining listening to %s log...", cl.chain))

	go func() {
		for {
			select {
			case err := <-cl.subscription.Err():
				log.Fatal().Err(err).Msg("sub error")
			case log := <-cl.channel:
				fmt.Printf("txhash %s\n", log.TxHash)
				payload, err := cl.recievePayload(log.Data)
				if err == nil {
					cl.sendTransaction(payload)
				}
			}
		}
	}()
}

func (cl *ChainETHish) recievePayload(data []byte) (Payload, error) {
	vals, err := cl.mpi_abi.Unpack("ZetaMessageSendEvent", data)
	if err != nil {
		return Payload{}, fmt.Errorf("Unpack error %s\n", err)
	}
	//    event ZetaMessageSendEvent(address sender, string destChainID, string  destContract, uint zetaAmount, uint gasLimit, bytes message, bytes32 messageID);

	sender := vals[0].(ethcommon.Address)
	fmt.Printf("sender %s\n", sender)

	destChainID := vals[1].(string)
	fmt.Printf("destChainID %s\n", destChainID)

	destContract := vals[2].(string)
	if destContract == "" {
		destContract = cl.DEFAULT_DESTINATION_CONTRACT
	}
	fmt.Printf("destContract %s\n", destContract)

	zetaAmount := vals[3].(*big.Int)
	fmt.Printf("zetaAmount %d\n", zetaAmount)

	gasLimit := vals[4].(*big.Int)
	fmt.Printf("gasLimit %d\n", gasLimit)

	message := vals[5].([]byte)
	fmt.Printf("message %s\n", hex.EncodeToString(message))

	msgid := vals[6].([32]byte)
	fmt.Printf("messageID %s\n", hex.EncodeToString(msgid[:]))

	return Payload{
		sender:       sender,
		destChainID:  destChainID,
		destContract: destContract,
		zetaAmount:   zetaAmount,
		gasLimit:     gasLimit,
		message:      message,
		msgid:        msgid,
	}, nil
}

// Contract signature:
//
// function zetaMessageReceive(
//	 address sender,
//	 string calldata destChainID,
//	 address destContract,
//	 uint zetaAmount,
//	 uint gasLimit,
//	 bytes calldata message,
//	 bytes32 messageID,
//	 bytes32 sendHash) external {
func (cl *ChainETHish) sendTransaction(payload Payload) {
	data, err := cl.mpi_abi.Pack(
		"zetaMessageReceive",
		payload.sender,
		payload.destChainID,
		ethcommon.HexToAddress(payload.destContract),
		payload.zetaAmount,
		payload.gasLimit,
		payload.message,
		payload.msgid,
		payload.msgid)
	if err != nil {
		fmt.Printf("Pack error %s\n", err)
		return
	}

	pair, err := FindChainByID(payload.destChainID)
	if err != nil {
		fmt.Printf("Chain ID error %s\n", err)
		return
	}

	nonce, err := pair.client.PendingNonceAt(context.Background(), cl.tss.Address())
	if err != nil {
		fmt.Printf("PendingNonceAt error %s\n", err)
		return
	}

	gasPrice, err := pair.client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Printf("SuggestGasPrice error %s\n", err)
		return
	}
	GasLimit := payload.gasLimit.Uint64()
	ethSigner := ethtypes.LatestSignerForChainID(pair.id)
	pair_mpi := ethcommon.HexToAddress(pair.contract)
	tx := ethtypes.NewTransaction(nonce, pair_mpi, big.NewInt(0), GasLimit, gasPrice, data)
	hashBytes := ethSigner.Hash(tx).Bytes()
	sig, err := cl.tss.Sign(hashBytes)
	if err != nil {
		fmt.Printf("tss.Sign error %s\n", err)
		return
	}

	signedTX, err := tx.WithSignature(ethSigner, sig[:])
	if err != nil {
		fmt.Printf("tx.WithSignature error %s\n", err)
		return
	}

	err = pair.client.SendTransaction(cl.context, signedTX)
	if err != nil {
		fmt.Printf("SendTransaction error %s\n", err)
		return
	}
	fmt.Printf("bcast tx %s done!\n", signedTX.Hash().Hex())
}
