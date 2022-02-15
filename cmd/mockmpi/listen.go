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
	log.Info().Msg(fmt.Sprintf("begining listening to %s log...", cl.name))

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

// Contract signature:
//
// event ZetaMessageSendEvent(
//   address sender,
//   string destChainID,
//   string  destContract,
//   uint zetaAmount,
//   uint gasLimit,
//   bytes message,
//   bytes32 messageID);
func (cl *ChainETHish) recievePayload(data []byte) (Payload, error) {
	var log_message string

	vals, err := cl.mpi_abi.Unpack("ZetaMessageSendEvent", data)
	if err != nil {
		return Payload{}, fmt.Errorf("Unpack error %s\n", err)
	}

	sender := vals[0].(ethcommon.Address)
	log_message = fmt.Sprintf("sender %s", sender)

	destChainID := vals[1].(string)
	log_message = fmt.Sprintf("destChainID %s", destChainID)
	log.Debug().Msg(log_message)

	destContract := vals[2].(string)
	if destContract == "" {
		destContract = cl.DEFAULT_DESTINATION_CONTRACT
	}
	log_message = fmt.Sprintf("destContract %s", destContract)
	log.Debug().Msg(log_message)

	zetaAmount := vals[3].(*big.Int)
	log_message = fmt.Sprintf("zetaAmount %d", zetaAmount)
	log.Debug().Msg(log_message)

	gasLimit := vals[4].(*big.Int)
	log_message = fmt.Sprintf("gasLimit %d", gasLimit)
	log.Debug().Msg(log_message)

	message := vals[5].([]byte)
	log_message = fmt.Sprintf("message %s", hex.EncodeToString(message))
	log.Debug().Msg(log_message)

	msgid := vals[6].([32]byte)
	log_message = fmt.Sprintf("messageID %s", hex.EncodeToString(msgid[:]))
	log.Debug().Msg(log_message)

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
		log.Err(err).Msg("sendTransaction() ABI Pack() error")
		return
	}

	other, err := FindChainByID(payload.destChainID)
	if err != nil {
		log.Err(err).Msg("sendTransaction() Chain ID error")
		return
	}

	nonce, err := other.client.PendingNonceAt(context.Background(), cl.tss.Address())
	if err != nil {
		log.Err(err).Msg("sendTransaction() PendingNonceAt error")
		return
	}

	gasPrice, err := other.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Err(err).Msg("sendTransaction() SuggestGasPrice error")
		return
	}
	GasLimit := payload.gasLimit.Uint64()

	ethSigner := ethtypes.LatestSignerForChainID(other.id)
	other_mpi := ethcommon.HexToAddress(other.MPI_CONTRACT)
	tx := ethtypes.NewTransaction(nonce, other_mpi, big.NewInt(0), GasLimit, gasPrice, data)
	hashBytes := ethSigner.Hash(tx).Bytes()
	sig, err := cl.tss.Sign(hashBytes)
	if err != nil {
		log.Err(err).Msg("sendTransaction() tss.Sign error")
		return
	}

	signedTX, err := tx.WithSignature(ethSigner, sig[:])
	if err != nil {
		log.Err(err).Msg("sendTransaction() tx.WithSignature error")
		return
	}

	err = other.client.SendTransaction(cl.context, signedTX)
	if err != nil {
		log.Err(err).Msg("sendTransaction() error")
		return
	}

	log.Info().Str("hash", signedTX.Hash().Hex()).Msg("bcast tx done!")
}
