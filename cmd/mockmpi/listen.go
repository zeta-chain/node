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
	sender       []byte
	srcChainID   uint16
	destChainID  uint16
	destContract []byte
	zetaAmount   *big.Int
	gasLimit     *big.Int
	message      []byte
	zetaParams   []byte
}

func (cl *ChainETHish) Listen() {
	log.Info().Msg(fmt.Sprintf("beginning listening to %s log...", cl.name))
	dedup := make(map[ethcommon.Hash]bool)
	go func() {
		for {
			select {
			case err := <-cl.subscription.Err():
				log.Fatal().Err(err).Msg("sub error")
			case log := <-cl.channel:
				fmt.Printf("txhash %s\n", log.TxHash)
				if _, ok := dedup[log.TxHash]; ok {
					fmt.Printf("txhash %s already processed!\n", log.TxHash)
				} else {
					dedup[log.TxHash] = true
					payload, err := cl.recievePayload(log.Topics, log.Data)
					if err == nil {
						cl.sendTransaction(payload)
					}
				}
			}
		}
	}()
}

// Contract signature:
//
// event ZetaMessageSendEvent(
//   uint16 destChainID,
//   bytes  destContract,
//   uint zetaAmount,
//   uint gasLimit,
//   bytes message,
//   bytes zetaParams);
func (cl *ChainETHish) recievePayload(topics []ethcommon.Hash, data []byte) (Payload, error) {
	var log_message string

	vals, err := cl.mpi_abi.Unpack("ZetaMessageSendEvent", data)
	if err != nil {
		return Payload{}, fmt.Errorf("Unpack error %s\n", err)
	}

	sender := topics[1]
	log_message = fmt.Sprintf("sender %x", sender)
	log.Debug().Msg(log_message)

	destChainID := vals[0].(uint16)
	log_message = fmt.Sprintf("destChainID %d", destChainID)
	log.Debug().Msg(log_message)

	destContract := vals[1].([]byte)
	log_message = fmt.Sprintf("destContract %x", destContract)
	log.Debug().Msg(log_message)

	zetaAmount := vals[2].(*big.Int)
	log_message = fmt.Sprintf("zetaAmount %d", zetaAmount)
	log.Debug().Msg(log_message)

	gasLimit := vals[3].(*big.Int)
	log_message = fmt.Sprintf("gasLimit %d", gasLimit)
	log.Debug().Msg(log_message)

	message := vals[4].([]byte)
	log_message = fmt.Sprintf("message %s", hex.EncodeToString(message))
	log.Debug().Msg(log_message)

	zetaParams := vals[5].([]byte)
	log_message = fmt.Sprintf("zetaParams %s", hex.EncodeToString(zetaParams[:]))
	log.Debug().Msg(log_message)

	return Payload{
		sender:       sender.Bytes(),
		srcChainID:   cl.chain_id,
		destChainID:  destChainID,
		destContract: destContract,
		zetaAmount:   zetaAmount,
		gasLimit:     gasLimit,
		message:      message,
		zetaParams:   zetaParams,
	}, nil
}

// Contract signature:
//
// function zetaMessageReceive(
//	 bytes sender,
//	 uint16  destChainID,
//	 address destContract,
//	 uint zetaAmount,
//	 bytes calldata message,
//	 bytes32 sendHash) external {
func (cl *ChainETHish) sendTransaction(payload Payload) {
	sendHash, err := hex.DecodeString(MAGIC_HASH[2:])
	if err != nil {
		log.Error().Err(err).Msg("sendTransaction: DecodeString err")
	}
	var sendHash32 [32]byte
	copy(sendHash32[:], sendHash[:32])
	data, err := cl.mpi_abi.Pack(
		"zetaMessageReceive",
		payload.sender,
		payload.srcChainID,
		ethcommon.BytesToAddress(payload.destContract),
		payload.zetaAmount,
		payload.message,
		sendHash32)
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
