package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

type Payload struct {
	sender       []byte
	srcChainID   *big.Int
	destChainID  *big.Int
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
//   uint256 destChainID,
//   bytes  destContract,
//   uint zetaAmount,
//   uint gasLimit,
//   bytes message,
//   bytes zetaParams);
func (cl *ChainETHish) recievePayload(topics []ethcommon.Hash, data []byte) (Payload, error) {
	//("ZetaSent(address,uint16,bytes,uint256,uint256,bytes,bytes)")
	vals, err := cl.mpiAbi.Unpack("ZetaSent", data)
	if err != nil {
		return Payload{}, fmt.Errorf("unpack error %s", err.Error())
	}

	sender := topics[1]
	destChainID := vals[0].(*big.Int)
	destContract := vals[1].([]byte)
	zetaAmount := vals[2].(*big.Int)
	gasLimit := vals[3].(*big.Int)
	message := vals[4].([]byte)
	zetaParams := vals[5].([]byte)

	log.Debug().Msgf("sender %s", sender)
	log.Debug().Msgf("destChainID %d", destChainID)
	log.Debug().Msgf("destContract %s, len %d", hex.EncodeToString(destContract), len(destContract))
	log.Debug().Msgf("zetaAmount %d", zetaAmount)
	log.Debug().Msgf("gasLimit %s", gasLimit)
	log.Debug().Msgf("message %s", message)
	log.Debug().Msgf("zetaParams %s", zetaParams)

	return Payload{
		sender:       sender.Bytes(),
		srcChainID:   cl.chainID,
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
	sendHash, err := hex.DecodeString(MagicHash[2:])
	if err != nil {
		log.Error().Err(err).Msg("sendTransaction: DecodeString err")
	}
	var sendHash32 [32]byte
	copy(sendHash32[:], sendHash[:32])
	data, err := cl.mpiAbi.Pack(
		"onReceive",
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
		log.Err(err).Msg("sendTransaction() Chain ID error; reverting...")
		cl.revertTransaction(payload)
		return
	}

	nonce, err := other.client.PendingNonceAt(context.Background(), cl.tss.EVMAddress())
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
	otherMpi := ethcommon.HexToAddress(other.MpiContract)
	tx := ethtypes.NewTransaction(nonce, otherMpi, big.NewInt(0), GasLimit, gasPrice, data)
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

	// tracking outbound tx:
	go func() {
		log.Info().Msgf("[%s] tracking outbound tx %s", other.name, signedTX.Hash().Hex())
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			receipt, err := other.client.TransactionReceipt(context.TODO(), signedTX.Hash())
			if err != nil {
				log.Debug().Err(err).Msgf("receipt non-existent: chain %s tx %s", other.name, signedTX.Hash())
				continue
			}
			if receipt.Status == 1 { // Successful tx
				log.Info().Msgf("tx %s succeed!", signedTX.Hash())
				return
			}
			log.Info().Msgf("tx %s reverted! initiating revert on origin chain...", signedTX.Hash())
			cl.revertTransaction(payload)
			return
		}
	}()
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
func (cl *ChainETHish) revertTransaction(payload Payload) {
	sendHash, err := hex.DecodeString(MagicHash[2:])
	if err != nil {
		log.Error().Err(err).Msg("revertTransaction: DecodeString err")
	}
	var sendHash32 [32]byte
	copy(sendHash32[:], sendHash[:32])
	data, err := cl.mpiAbi.Pack(
		"onRevert",
		ethcommon.BytesToAddress(payload.sender),
		payload.srcChainID,
		payload.destContract,
		payload.destChainID,
		payload.zetaAmount,
		payload.message,
		sendHash32)
	if err != nil {
		log.Err(err).Msg("revertTransaction() ABI Pack() error")
		return
	}

	nonce, err := cl.client.PendingNonceAt(context.Background(), cl.tss.EVMAddress())
	if err != nil {
		log.Err(err).Msg("revertTransaction() PendingNonceAt error")
		return
	}

	gasPrice, err := cl.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Err(err).Msg("revertTransaction() SuggestGasPrice error")
		return
	}
	GasLimit := payload.gasLimit.Uint64()

	ethSigner := ethtypes.LatestSignerForChainID(cl.id)
	otherMpi := ethcommon.HexToAddress(cl.MpiContract)
	tx := ethtypes.NewTransaction(nonce, otherMpi, big.NewInt(0), GasLimit, gasPrice, data)
	hashBytes := ethSigner.Hash(tx).Bytes()
	sig, err := cl.tss.Sign(hashBytes)
	if err != nil {
		log.Err(err).Msg("revertTransaction() tss.Sign error")
		return
	}

	signedTX, err := tx.WithSignature(ethSigner, sig[:])
	if err != nil {
		log.Err(err).Msg("revertTransaction() tx.WithSignature error")
		return
	}

	err = cl.client.SendTransaction(cl.context, signedTX)
	if err != nil {
		log.Err(err).Msg("revertTransaction() error")
		return
	}

	log.Info().Str("hash", signedTX.Hash().Hex()).Msg("bcast tx done!")

	// tracking outbound tx:
	go func() {
		log.Info().Msgf("[%s] tracking outbound tx %s", cl.name, signedTX.Hash().Hex())

		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			receipt, err := cl.client.TransactionReceipt(context.TODO(), signedTX.Hash())
			if err != nil {
				log.Debug().Err(err).Msgf("revert receipt non-existent: chain %s tx %s", cl.name, signedTX.Hash())
				continue
			}
			if receipt.Status == 1 { // Successful tx
				log.Info().Msgf("onRevert tx %s succeed!", signedTX.Hash())
			} else { // revert
				log.Info().Msgf("onRevert tx %s reverted!", signedTX.Hash())
			}
			return
		}
	}()
}
