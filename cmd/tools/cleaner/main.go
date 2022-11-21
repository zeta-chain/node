package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/cmd"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
	"github.com/zeta-chain/zetacore/contracts/evm"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"math/big"
	"os"
	"os/signal"
	"path"
	"sort"
	"strings"
	"syscall"
)

var (
	NewIndex bool
)

func main() {
	nodeIP := flag.String("node-ip", "", "IP address of the zetacore node")
	signerName := flag.String("signer-name", "admin", "name of the signer")
	enabledChains := flag.String("enable-chains", "GOERLI,BSCTESTNET,MUMBAI", "enable chains, comma separated list")
	index := flag.String("index", "", "index [TSS address]: collect all txs originating from TSS address and put them into sqlite3 db [chain].sqlite3")
	newIndex := flag.Bool("new-index", false, "first index")
	flag.Parse()
	chains := strings.Split(*enabledChains, ",")
	NewIndex = *newIndex

	indexerMap := make(map[string]*Indexer)
	if len(*index) > 0 {
		log.Info().Msg("Index mode...")
		for _, chain := range chains {
			endpoint := os.Getenv(fmt.Sprintf("%s_ENDPOINT", chain))
			if endpoint == "" {
				log.Fatal().Msgf("envvar %s_ENDPOINT is not set", chain)
				return
			}
			indexer, err := NewIndexer(chain, endpoint, *index)
			if err != nil {
				log.Fatal().Err(err).Msgf("NewIndexer error")
				return
			}
			indexerMap[chain] = indexer
			indexer.Start()
		}

		// wait....
		log.Info().Msgf("awaiting the os.Interrupt, syscall.SIGTERM signals...")
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		sig := <-ch
		log.Info().Msgf("stop signal received: %s", sig)

		for _, chain := range chains {
			indexerMap[chain].Stop()
		}
		return
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("fail to get user home dir")
	}
	coreDir := path.Join(userHome, ".zetacored")
	signerPass := "password"
	kb, _, err := mc.GetKeyringKeybase(coreDir, *signerName, signerPass)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to get keyring keybase")
		return
	}

	k := mc.NewKeysWithKeybase(kb, *signerName, signerPass)
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount(cmd.Bech32PrefixAccAddr, cmd.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(cmd.Bech32PrefixValAddr, cmd.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(cmd.Bech32PrefixConsAddr, cmd.Bech32PrefixConsPub)
	cmd.CHAINID = "athens_7001-1"
	bridge, err := mc.NewZetaCoreBridge(k, *nodeIP, *signerName)
	if err != nil {
		log.Fatal().Err(err).Msg("NewZetaCoreBridge")
		return
	}

	pendingCctx, err := bridge.GetAllPendingCctx()
	if err != nil {
		log.Error().Err(err).Msg("GetAllPendingCctx")
		return
	}

	nonceToCctx := make(map[string]*types.CrossChainTx)

	log.Info().Msgf("pending cctx: %d", len(pendingCctx))
	sendMap := splitAndSortSendListByChain(pendingCctx)
	for _, c := range chains {
		list := sendMap[c]
		var nonces []uint64
		for _, cctx := range list {
			nonce := cctx.OutBoundTxParams.OutBoundTxTSSNonce
			nonces = append(nonces, nonce)
			outTxID := fmt.Sprintf("%s-%d", c, nonce)
			nonceToCctx[outTxID] = cctx
		}
		chain, err := common.ParseChain(c)
		if err != nil {
			panic(err)
		}
		outTxTracker, err := bridge.GetAllOutTxTrackerByChain(chain)
		if err != nil {
			panic(err)
		}
		log.Info().Msgf("outTxTracker: %d", len(outTxTracker))

		log.Info().Msgf("chain %s has %d pending cctx, divided into %d intervals", c, len(list), len(BreakSortedSequenceIntoIntervals(nonces)))
		intervals := BreakSortedSequenceIntoIntervals(nonces)
		for idx, interval := range intervals {
			log.Info().Msgf("  interval[%d]: %d - %d", idx, interval[0], interval[len(interval)-1])
		}

		ethClient, err := ethclient.Dial(os.Getenv(fmt.Sprintf("%s_ENDPOINT", c)))
		if err != nil {
			log.Error().Err(err).Msgf("fail to dial %s", c)
			continue
		}
		connectorAddr := ethcommon.HexToAddress(os.Getenv(fmt.Sprintf("%s_CONNECTOR", c)))
		if connectorAddr == (ethcommon.Address{}) {
			log.Error().Msgf("envvar %s_CONNECTOR is not set", c)
			continue
		}
		conn, err := evm.NewConnector(connectorAddr, ethClient)
		if err != nil {
			log.Error().Err(err).Msgf("fail to create connector %s", c)
			continue
		}
		for idx, interval := range intervals {
			if idx == 0 {
				//if idx < len(intervals)-1 {
				for _, nonce := range interval {
					outTxID := fmt.Sprintf("%s-%d", c, nonce)
					log.Info().Msgf("  fixing %s", outTxID)
					cctx, found := nonceToCctx[outTxID]
					if !found {
						log.Error().Msgf("  cctx not found for %s", outTxID)
						continue
					}
					sendHash, err := hex.DecodeString(cctx.Index[2:])
					if err != nil {
						log.Error().Err(err).Msgf("  fail to decode sendHash %s", cctx.Index)
						continue
					}
					var sendHashB32 [32]byte
					copy(sendHashB32[:32], sendHash[:32])
					bn, err := ethClient.BlockNumber(context.Background())
					if err != nil {
						log.Error().Err(err).Msgf("  fail to get block number")
						continue
					}
					logs, err := conn.FilterZetaReceived(&bind.FilterOpts{
						Start:   0,
						End:     &bn,
						Context: context.TODO(),
					}, []*big.Int{}, []ethcommon.Address{}, [][32]byte{sendHashB32})
					if err != nil {
						log.Error().Err(err).Msg("  fail to filter ZetaReceived")
						continue
					}
					for logs.Next() {
						log.Info().Msgf("  found zeta received: %s", logs.Event.Raw.TxHash.Hex())
						//txhash := logs.Event.Raw.TxHash.Hex()
						//
						//zTxHash, err := bridge.AddTxHashToOutTxTracker(chain.String(), nonce, txhash)
						//if err != nil {
						//	log.Error().Err(err).Msgf("  fail to add txhash to outTxTracker")
						//	continue
						//}
						//log.Info().Msgf("  outTxTracker tx: %s", zTxHash)
					}
				}
			}
		}
	}

}

func splitAndSortSendListByChain(sendList []*types.CrossChainTx) map[string][]*types.CrossChainTx {
	sendMap := make(map[string][]*types.CrossChainTx)
	for _, send := range sendList {
		targetChain := mc.GetTargetChain(send)
		if targetChain == "" {
			continue
		}
		if _, found := sendMap[targetChain]; !found {
			sendMap[targetChain] = make([]*types.CrossChainTx, 0)
		}
		sendMap[targetChain] = append(sendMap[targetChain], send)
	}
	for _, sends := range sendMap {
		sort.Slice(sends, func(i, j int) bool {
			return sends[i].OutBoundTxParams.OutBoundTxTSSNonce < sends[j].OutBoundTxParams.OutBoundTxTSSNonce
		})
	}
	return sendMap
}

func BreakSortedSequenceIntoIntervals(seq []uint64) [][]uint64 {
	var intervals [][]uint64
	for i := 0; i < len(seq); i++ {
		if i == 0 {
			intervals = append(intervals, []uint64{seq[i]})
		} else {
			if seq[i] == seq[i-1]+1 {
				intervals[len(intervals)-1] = append(intervals[len(intervals)-1], seq[i])
			} else {
				intervals = append(intervals, []uint64{seq[i]})
			}
		}
	}

	return intervals
}
