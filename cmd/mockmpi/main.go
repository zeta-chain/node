package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	BSC_MPI = "0xCC3e1C9460B7803d4d79F32342b2b27543362536"
	ETH_MPI = "0xDe8802902Ff3136bdACe5FFC9a2423B1d37F6833"
	ABI_MPI = `[{"inputs":[{"internalType":"address","name":"zetaAddress","type":"address"},{"internalType":"address","name":"oracleAddress","type":"address"},{"internalType":"address","name":"_TSSAddress","type":"address"},{"internalType":"address","name":"_TSSAddressUpdater","type":"address"},{"internalType":"address","name":"_OracleUpdater","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"string","name":"destChainID","type":"string"},{"indexed":false,"internalType":"address","name":"destContract","type":"address"},{"indexed":false,"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"gasLimit","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"},{"indexed":false,"internalType":"bytes32","name":"messageID","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"utxoHash","type":"bytes32"}],"name":"ZetaMessageReceiveEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"string","name":"destChainID","type":"string"},{"indexed":false,"internalType":"string","name":"destContract","type":"string"},{"indexed":false,"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"gasLimit","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"},{"indexed":false,"internalType":"bytes32","name":"messageID","type":"bytes32"}],"name":"ZetaMessageSendEvent","type":"event"},{"inputs":[],"name":"OracleUpdater","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TSSAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TSSAddressUpdater","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"newOracleAddres","type":"address"}],"name":"changeOracle","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getLockedAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceTSSAddressUpdater","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint8","name":"newFlexibility","type":"uint8"}],"name":"updateSupplyOracleFlexibility","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_address","type":"address"}],"name":"updateTSSAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"string","name":"destChainID","type":"string"},{"internalType":"address","name":"destContract","type":"address"},{"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"internalType":"uint256","name":"gasLimit","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"},{"internalType":"bytes32","name":"messageID","type":"bytes32"},{"internalType":"bytes32","name":"sendHash","type":"bytes32"}],"name":"zetaMessageReceive","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"string","name":"destChainID","type":"string"},{"internalType":"string","name":"destContract","type":"string"},{"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"internalType":"uint256","name":"gasLimit","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"},{"internalType":"bytes32","name":"messageID","type":"bytes32"}],"name":"zetaMessageSend","outputs":[],"stateMutability":"nonpayable","type":"function"}]
`
)

func main() {
	pkstring := os.Getenv("PRIVKEY")
	if pkstring == "" {
		log.Fatal().Msg("missing env variable PRIVKEY")
		return
	}
	privateKey, err := crypto.HexToECDSA(pkstring)
	if err != nil {
		log.Err(err).Msg("TEST private key error")
		return
	}
	tss := mc.TestSigner{
		PrivKey: privateKey,
	}
	eth_chain, err := mc.NewChainObserver(common.ETHChain, nil, tss, "")
	bsc_chain, err := mc.NewChainObserver(common.BSCChain, nil, tss, "")

	_ = eth_chain
	fmt.Printf("tss key address: %s\n", tss.Address())

	abi_mpi, err := abi.JSON(strings.NewReader(ABI_MPI))
	if err != nil {
		log.Err(err).Msg("abi.JSON")
		return
	}

	ctxt := context.TODO()

	bsc_client := bsc_chain.Client
	eth_client := eth_chain.Client
	cid, _ := bsc_client.ChainID(ctxt)
	fmt.Printf("chain id %d\n", cid)
	eth_chainID, err := eth_client.ChainID(context.TODO())
	if err != nil {
		fmt.Printf("eth_client.ChainID error %s\n", err)
		return
	}

	bsc_topics := 	make([][]ethcommon.Hash,1)
	bsc_topics[0] = []ethcommon.Hash{ethcommon.HexToHash("0x38f8fa9ce079e7e087c700936fd84330f80123e22a6aea6e125b55e95dcd585a")}
	bsc_query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(BSC_MPI)},
		Topics: bsc_topics,
	}

	ch := make(chan types.Log)
	sub, err := bsc_client.SubscribeFilterLogs(ctxt, bsc_query, ch)
	if err != nil {
		fmt.Printf("SubscribeFilterLogs error %s\n", err)
		return
	}

	fmt.Println("begining listening to BSC log...")
	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatal().Err(err).Msg("sub error")
			case log := <-ch:
				//fmt.Printf("%#v\n", log)
				fmt.Printf("txhash %s\n", log.TxHash)
				vals, err := abi_mpi.Unpack("ZetaMessageSendEvent", log.Data)
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
				//     function zetaMessageReceive(address sender, string calldata destChainID, address destContract, uint zetaAmount, uint gasLimit, bytes calldata message, bytes32 messageID, bytes32 sendHash) external {
				data, err := abi_mpi.Pack("zetaMessageReceive", sender, destChainID, ethcommon.HexToAddress(destContract), zetaAmount, gasLimit, message, msgid, msgid)
				if err != nil {
					fmt.Printf("Pack error %s\n", err)
					continue
				}
				nonce, err := eth_client.PendingNonceAt(context.Background(), tss.Address())
				if err != nil {
					fmt.Printf("PendingNonceAt error %s\n", err)
					continue
				}
				gasPrice, err := eth_client.SuggestGasPrice(context.Background())
				if err != nil {
					fmt.Printf("SuggestGasPrice error %s\n", err)
					continue
				}
				//GasLimit := uint64(300000) // in units
				GasLimit := gasLimit.Uint64()
				ethSigner := ethtypes.LatestSignerForChainID(eth_chainID)
				eth_mpi := ethcommon.HexToAddress(ETH_MPI)
				tx := ethtypes.NewTransaction(nonce, eth_mpi, big.NewInt(0), GasLimit, gasPrice, data)
				hashBytes := ethSigner.Hash(tx).Bytes()
				sig, err := tss.Sign(hashBytes)
				if err != nil {
					fmt.Printf("tss.Sign error %s\n", err)
					continue
				}
				signedTX, err := tx.WithSignature(ethSigner, sig[:])
				if err != nil {
					fmt.Printf("tx.WithSignature error %s\n", err)
					continue
				}
				err = eth_client.SendTransaction(ctxt, signedTX)
				if err != nil {
					fmt.Printf("SendTransaction error %s\n", err)
					continue
				}
				fmt.Printf("bcast tx %s done!\n", signedTX.Hash().Hex())
			}
		}
	}()


	eth_cid, _ := eth_client.ChainID(ctxt)
	fmt.Printf("eth chain id %d\n", eth_cid)
	bsc_chainID, err := bsc_client.ChainID(context.TODO())
	if err != nil {
		fmt.Printf("eth_client.ChainID error %s\n", err)
		return
	}

	eth_topics := 	make([][]ethcommon.Hash,1)
	eth_topics[0] = []ethcommon.Hash{ethcommon.HexToHash("0x38f8fa9ce079e7e087c700936fd84330f80123e22a6aea6e125b55e95dcd585a")}
	eth_filter := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(ETH_MPI)},
		Topics: bsc_topics,
	}

	ch2 := make(chan types.Log)
	sub2, err := eth_client.SubscribeFilterLogs(ctxt, eth_filter, ch2)
	if err != nil {
		fmt.Printf("SubscribeFilterLogs error %s\n", err)
		return
	}
	fmt.Println("begining listening to Goerli log...")
	go func() {
		for {
			select {
			case err := <-sub2.Err():
				log.Fatal().Err(err).Msg("sub error")
			case log := <-ch2:
				//fmt.Printf("%#v\n", log)
				fmt.Printf("txhash %s\n", log.TxHash)
				vals, err := abi_mpi.Unpack("ZetaMessageSendEvent", log.Data)
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
					destContract = "0xFf6B270ac3790589A1Fe90d0303e9D4d9A54FD1A"
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
				//     function zetaMessageReceive(address sender, string calldata destChainID, address destContract, uint zetaAmount, uint gasLimit, bytes calldata message, bytes32 messageID, bytes32 sendHash) external {
				data, err := abi_mpi.Pack("zetaMessageReceive", sender, destChainID, ethcommon.HexToAddress(destContract), zetaAmount, gasLimit, message, msgid, msgid)
				if err != nil {
					fmt.Printf("Pack error %s\n", err)
					continue
				}
				nonce, err := bsc_client.PendingNonceAt(context.Background(), tss.Address())
				if err != nil {
					fmt.Printf("PendingNonceAt error %s\n", err)
					continue
				}
				gasPrice, err := bsc_client.SuggestGasPrice(context.Background())
				if err != nil {
					fmt.Printf("SuggestGasPrice error %s\n", err)
					continue
				}
				//GasLimit := uint64(300000) // in units
				GasLimit := gasLimit.Uint64()
				ethSigner := ethtypes.LatestSignerForChainID(bsc_chainID)
				bsc_mpi := ethcommon.HexToAddress(BSC_MPI)
				tx := ethtypes.NewTransaction(nonce, bsc_mpi, big.NewInt(0), GasLimit, gasPrice, data)
				hashBytes := ethSigner.Hash(tx).Bytes()
				sig, err := tss.Sign(hashBytes)
				if err != nil {
					fmt.Printf("tss.Sign error %s\n", err)
					continue
				}
				signedTX, err := tx.WithSignature(ethSigner, sig[:])
				if err != nil {
					fmt.Printf("tx.WithSignature error %s\n", err)
					continue
				}
				err = bsc_client.SendTransaction(ctxt, signedTX)
				if err != nil {
					fmt.Printf("SendTransaction error %s\n", err)
					continue
				}
				fmt.Printf("bcast tx %s done!\n", signedTX.Hash().Hex())
			}
		}
	}()

	ch3 := make(chan os.Signal, 1)
	signal.Notify(ch3, syscall.SIGINT, syscall.SIGTERM)
	<-ch3
	log.Info().Msg("stop signal received")
}
