package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"os/signal"
	"path"
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

	log.Info().Msgf("pending cctx: %d", len(pendingCctx))
	sendMap := mc.SplitAndSortSendListByChain(pendingCctx)
	for _, c := range chains {
		list := sendMap[c]
		var nonces []uint64
		for _, cctx := range list {
			nonces = append(nonces, cctx.OutBoundTxParams.OutBoundTxTSSNonce)
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
		for _, interval := range BreakSortedSequenceIntoIntervals(nonces) {
			log.Info().Msgf("  interval: %d - %d", interval[0], interval[len(interval)-1])
			outTxTracker, err := bridge.GetOutTxTracker(chain, interval[0])
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					log.Warn().Err(err).Msg("unknown gRPC error code")
				} else {
					if st.Code() == codes.NotFound {
						log.Info().Msgf("  outTxTracker not found for nonce %d", interval[0])
					}
				}
			} else { // found
				log.Info().Msgf("  outTxTracker found for nonce %d: %v", interval[0], outTxTracker.HashList)
			}

		}

	}

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
