package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient"
)

const (
	BSC_MPI_CONTRACT = "0xCC3e1C9460B7803d4d79F32342b2b27543362536"
	ETH_MPI_CONTRACT = "0xDe8802902Ff3136bdACe5FFC9a2423B1d37F6833"

	// What is this?
	MAGIC_HASH = "0x38f8fa9ce079e7e087c700936fd84330f80123e22a6aea6e125b55e95dcd585a"
)

type ChainETHish struct {
	tss          zetaclient.TSSSigner
	mpi_abi      abi.ABI
	contract     string
	context      context.Context
	client       *ethclient.Client
	chain        common.Chain
	id           *big.Int
	topics       [][]ethcommon.Hash
	channel      chan types.Log
	subscription ethereum.Subscription
}

func (cl *ChainETHish) Setup() {
	cl.tss = GetZetaTestSignature()

	_abi, err := abi.JSON(strings.NewReader(ABI_MPI))
	if err != nil {
		log.Err(err).Msg("abi.JSON")
		os.Exit(1)
	}
	cl.mpi_abi = _abi

	cl.context = context.TODO()

	chain, err := zetaclient.NewChainObserver(cl.chain, nil, cl.tss, "")
	cl.client = chain.Client

	_id, _ := cl.client.ChainID(cl.context)
	fmt.Printf("BSC chain id %d\n", _id)
	cl.id, err = cl.client.ChainID(context.TODO())
	if err != nil {
		fmt.Printf("Chain.id error %s\n", err)
		os.Exit(1)
	}

	cl.topics = make([][]ethcommon.Hash, 1)
	cl.topics[0] = []ethcommon.Hash{ethcommon.HexToHash(MAGIC_HASH)}
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(cl.contract)},
		Topics:    cl.topics,
	}

	cl.channel = make(chan types.Log)

	_subscription, err := cl.client.SubscribeFilterLogs(cl.context, query, cl.channel)
	if err != nil {
		log.Printf("SubscribeFilterLogs error %s\n", err)
		os.Exit(1)
	}
	cl.subscription = _subscription
}

func (cl *ChainETHish) Listen(pair *ChainETHish) {
	log.Info().Msg(fmt.Sprintf("begining listening to %s log...", cl.chain))

	go func() {
		for {
			select {
			case err := <-cl.subscription.Err():
				log.Fatal().Err(err).Msg("sub error")
			case log := <-cl.channel:
				fmt.Printf("txhash %s\n", log.TxHash)
				vals, err := cl.mpi_abi.Unpack("ZetaMessageSendEvent", log.Data)
				if err != nil {
					fmt.Printf("Unpack error %s\n", err)
					continue
				}
				//    event ZetaMessageSendEvent(address sender, string destChainID, string  destContract, uint zetaAmount, uint gasLimit, bytes message, bytes32 messageID);

				sender := vals[0].(ethcommon.Address)
				fmt.Printf("sender %s\n", sender)

				destChainID := vals[1].(string)
				fmt.Printf("destChainID %s\n", destChainID)

				destContract := vals[2].(string)
				if destContract == "" {
					destContract = "0xF47bd84B86d1667e7621c38c72C6905Ca8710b0d"
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
				data, err := cl.mpi_abi.Pack(
					"zetaMessageReceive",
					sender,
					destChainID,
					ethcommon.HexToAddress(destContract),
					zetaAmount,
					gasLimit,
					message,
					msgid,
					msgid)
				if err != nil {
					fmt.Printf("Pack error %s\n", err)
					continue
				}

				nonce, err := pair.client.PendingNonceAt(context.Background(), cl.tss.Address())
				if err != nil {
					fmt.Printf("PendingNonceAt error %s\n", err)
					continue
				}

				gasPrice, err := pair.client.SuggestGasPrice(context.Background())
				if err != nil {
					fmt.Printf("SuggestGasPrice error %s\n", err)
					continue
				}
				GasLimit := gasLimit.Uint64()
				ethSigner := ethtypes.LatestSignerForChainID(pair.id)
				pair_mpi := ethcommon.HexToAddress(pair.contract)
				tx := ethtypes.NewTransaction(nonce, pair_mpi, big.NewInt(0), GasLimit, gasPrice, data)
				hashBytes := ethSigner.Hash(tx).Bytes()
				sig, err := cl.tss.Sign(hashBytes)
				if err != nil {
					fmt.Printf("tss.Sign error %s\n", err)
					continue
				}

				signedTX, err := tx.WithSignature(ethSigner, sig[:])
				if err != nil {
					fmt.Printf("tx.WithSignature error %s\n", err)
					continue
				}

				err = pair.client.SendTransaction(cl.context, signedTX)
				if err != nil {
					fmt.Printf("SendTransaction error %s\n", err)
					continue
				}
				fmt.Printf("bcast tx %s done!\n", signedTX.Hash().Hex())
			}
		}
	}()
}

func main() {
	bsc := &ChainETHish{chain: common.Chain("BSC"), contract: BSC_MPI_CONTRACT}
	bsc.Setup()

	eth := &ChainETHish{chain: common.Chain("ETH"), contract: ETH_MPI_CONTRACT}
	eth.Setup()

	bsc.Listen(eth)
	eth.Listen(bsc)

	ch3 := make(chan os.Signal, 1)
	signal.Notify(ch3, syscall.SIGINT, syscall.SIGTERM)
	<-ch3
	log.Info().Msg("stop signal received")
}
