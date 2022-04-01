package eth

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/zeta-chain/zetacore/cmd/mockmpi/common"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

type Payload struct {
	Sender       []byte
	SrcChainID   uint16
	DestChainID  uint16
	DestContract []byte
	ZetaAmount   *big.Int
	GasLimit     *big.Int
	Message      []byte
	ZetaParams   []byte
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
						cl.SendTransaction(payload)
					}
				}
			}
		}
	}()
}

// Contract signature:
//
// event ZetaMessageSendEvent(
//   uint16 DestChainID,
//   bytes  DestContract,
//   uint ZetaAmount,
//   uint GasLimit,
//   bytes Message,
//   bytes ZetaParams);
func (cl *ChainETHish) recievePayload(topics []ethcommon.Hash, data []byte) (Payload, error) {
	var log_message string

	vals, err := cl.mpi_abi.Unpack("ZetaMessageSendEvent", data)
	if err != nil {
		return Payload{}, fmt.Errorf("Unpack error %s\n", err)
	}

	sender := topics[1]
	log_message = fmt.Sprintf("Sender %x", sender)
	log.Debug().Msg(log_message)

	destChainID := vals[0].(uint16)
	log_message = fmt.Sprintf("DestChainID %d", destChainID)
	log.Debug().Msg(log_message)

	destContract := vals[1].([]byte)
	log_message = fmt.Sprintf("DestContract %x", destContract)
	log.Debug().Msg(log_message)

	zetaAmount := vals[2].(*big.Int)
	log_message = fmt.Sprintf("ZetaAmount %d", zetaAmount)
	log.Debug().Msg(log_message)

	gasLimit := vals[3].(*big.Int)
	log_message = fmt.Sprintf("GasLimit %d", gasLimit)
	log.Debug().Msg(log_message)

	message := vals[4].([]byte)
	log_message = fmt.Sprintf("Message %s", hex.EncodeToString(message))
	log.Debug().Msg(log_message)

	zetaParams := vals[5].([]byte)
	log_message = fmt.Sprintf("ZetaParams %s", hex.EncodeToString(zetaParams[:]))
	log.Debug().Msg(log_message)

	return Payload{
		Sender:       sender.Bytes(),
		SrcChainID:   cl.Chain_id,
		DestChainID:  destChainID,
		DestContract: destContract,
		ZetaAmount:   zetaAmount,
		GasLimit:     gasLimit,
		Message:      message,
		ZetaParams:   zetaParams,
	}, nil
}

// Contract signature:
//
// function zetaMessageReceive(
//	 bytes Sender,
//	 uint16  DestChainID,
//	 address DestContract,
//	 uint ZetaAmount,
//	 bytes calldata Message,
//	 bytes32 sendHash) external {
func (cl *ChainETHish) SendTransaction(payload Payload) {
	log.Info().Msgf("Sending transaction: %v", payload)
	sendHash, err := hex.DecodeString(MAGIC_HASH[2:])
	if err != nil {
		log.Error().Err(err).Msg("SendTransaction: DecodeString err")
	}
	var sendHash32 [32]byte
	copy(sendHash32[:], sendHash[:32])
	data, err := cl.mpi_abi.Pack(
		"zetaMessageReceive",
		payload.Sender,
		payload.SrcChainID,
		ethcommon.BytesToAddress(payload.DestContract),
		payload.ZetaAmount,
		payload.Message,
		sendHash32)
	if err != nil {
		log.Err(err).Msg("SendTransaction() ABI Pack() error")
		return
	}

	other, err := common.FindChainByID(payload.DestChainID)
	if err != nil {
		log.Err(err).Msg("SendTransaction() Chain ID error")
		return
	}

	nonce, err := other.(*ChainETHish).client.PendingNonceAt(context.Background(), cl.tss.Address())
	if err != nil {
		log.Err(err).Msg("SendTransaction() PendingNonceAt error")
		return
	}

	gasPrice, err := other.(*ChainETHish).client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Err(err).Msg("SendTransaction() SuggestGasPrice error")
		return
	}
	GasLimit := payload.GasLimit.Uint64()

	ethSigner := ethtypes.LatestSignerForChainID(other.(*ChainETHish).id)
	other_mpi := ethcommon.HexToAddress(other.(*ChainETHish).MPI_CONTRACT)
	tx := ethtypes.NewTransaction(nonce, other_mpi, big.NewInt(0), GasLimit, gasPrice, data)
	hashBytes := ethSigner.Hash(tx).Bytes()
	sig, err := cl.tss.Sign(hashBytes)
	if err != nil {
		log.Err(err).Msg("SendTransaction() tss.Sign error")
		return
	}

	signedTX, err := tx.WithSignature(ethSigner, sig[:])
	if err != nil {
		log.Err(err).Msg("SendTransaction() tx.WithSignature error")
		return
	}

	err = other.(*ChainETHish).client.SendTransaction(cl.context, signedTX)
	if err != nil {
		log.Err(err).Msg("SendTransaction() error")
		return
	}

	log.Info().Str("hash", signedTX.Hash().Hex()).Msg("bcast tx done!")
}
