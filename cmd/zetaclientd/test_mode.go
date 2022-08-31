package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p-peerstore/addr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	metrics2 "github.com/zeta-chain/zetacore/zetaclient/metrics"
	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"
)

func startTestMode(validatorName string, peers addr.AddrList, zetacoreHome string) {
	SetupConfigForTest() // setup meta-prefix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	chainIP := os.Getenv("CHAIN_IP")
	if chainIP == "" {
		chainIP = "127.0.0.1"
	}

	// wait until zetacore is up
	log.Info().Msg("Waiting for ZetaCore to open 9090 port...")
	for {
		_, err := grpc.Dial(
			fmt.Sprintf("%s:9090", chainIP),
			grpc.WithInsecure(),
		)
		if err != nil {
			log.Warn().Err(err).Msg("grpc dial fail")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	log.Info().Msgf("ZetaCore to open 9090 port...")

	// setup 2 metabridges
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Err(err).Msg("UserHomeDir error")
		return
	}
	chainHomeFoler := filepath.Join(homeDir, zetacoreHome)

	// first signer & bridge
	signerName := validatorName
	signerPass := "password"
	bridge1, done := CreateZetaBridge(chainHomeFoler, signerName, signerPass)
	if done {
		return
	}

	bridgePk, err := bridge1.GetKeys().GetPrivateKey()
	if err != nil {
		log.Error().Err(err).Msg("GetKeys GetPrivateKey error:")
	}
	if len(bridgePk.Bytes()) != 32 {
		log.Error().Msgf("key bytes len %d != 32", len(bridgePk.Bytes()))
		return
	}
	var priKey secp256k1.PrivKey
	priKey = bridgePk.Bytes()[:32]

	log.Info().Msgf("NewTSS: with peer pubkey %s", bridgePk.PubKey())
	tss, err := mc.NewTSS(peers, priKey, preParams)
	if err != nil {
		log.Error().Err(err).Msg("NewTSS error")
		return
	}

	consKey := ""
	pubkeySet, err := bridge1.GetKeys().GetPubKeySet()
	if err != nil {
		log.Error().Err(err).Msgf("Get Pubkey Set Error")
	}
	ztx, err := bridge1.SetNodeKey(pubkeySet, consKey)
	log.Info().Msgf("SetNodeKey: %s by node %s zeta tx %s", pubkeySet.Secp256k1.String(), consKey, ztx)
	if err != nil {
		log.Error().Err(err).Msgf("SetNodeKey error")
	}

	log.Info().Msg("wait for 20s for all node to SetNodeKey")
	time.Sleep(12 * time.Second)

	if keygenBlock > 0 {
		log.Info().Msgf("Keygen at blocknum %d", keygenBlock)
		bn, err := bridge1.GetZetaBlockHeight()
		if err != nil {
			log.Error().Err(err).Msg("GetZetaBlockHeight error")
			return
		}
		if int64(bn)+3 > keygenBlock {
			log.Warn().Msgf("Keygen at blocknum %d, but current blocknum %d", keygenBlock, bn)
			return
		}
		nodeAccounts, err := bridge1.GetAllNodeAccounts()
		if err != nil {
			log.Error().Err(err).Msg("GetAllNodeAccounts error")
			return
		}
		pubkeys := make([]string, 0)
		for _, na := range nodeAccounts {
			pubkeys = append(pubkeys, na.PubkeySet.Secp256k1.String())
		}
		ticker := time.NewTicker(time.Second * 2)
		for range ticker.C {
			bn, err := bridge1.GetZetaBlockHeight()
			if err != nil {
				log.Error().Err(err).Msg("GetZetaBlockHeight error")
				return
			}
			if int64(bn) == keygenBlock {
				break
			}
		}
		log.Info().Msgf("Keygen with %d TSS signers", len(nodeAccounts))
		log.Info().Msgf("%s", pubkeys)
		var req keygen.Request
		req = keygen.NewRequest(pubkeys, keygenBlock, "0.14.0")
		res, err := tss.Server.Keygen(req)
		if err != nil || res.Status != tsscommon.Success {
			log.Error().Msgf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)
			return
		}
		// Keygen succeed! Report TSS address
		log.Info().Msgf("Keygen success! keygen response: %v...", res)

		log.Info().Msgf("doing a keysign test...")
		err = mc.TestKeysign(res.PubKey, tss.Server)
		if err != nil {
			log.Error().Err(err).Msg("TestKeysign error")
			return
		}

		log.Info().Msgf("setting TSS pubkey: %s", res.PubKey)
		err = tss.InsertPubKey(res.PubKey)
		tss.CurrentPubkey = res.PubKey
		if err != nil {
			log.Error().Msgf("SetPubKey fail")
			return
		}
		log.Info().Msgf("TSS address in hex: %s", tss.Address().Hex())
		return
	}

	metrics, err := metrics2.NewMetrics()
	if err != nil {
		log.Error().Err(err).Msg("NewMetrics")
		return
	}
	metrics.Start()

	log.Info().Msg("starting keysign tests...")
	go startKeysignTest(bridge1, tss)

	// wait....
	log.Info().Msgf("awaiting the os.Interrupt, syscall.SIGTERM signals...")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	log.Info().Msgf("stop signal received: %s", sig)

}

func startKeysignTest(bridge *mc.ZetaCoreBridge, tss *mc.TSS) {
	ticker := time.NewTicker(2 * time.Second)
	var lastZetaBlock uint64 = 0
	var numConcurrentKeysign int64 = 0
	var totalKeysign int64 = 0
	var successfulKeysign int64 = 0
	startBlock, err := bridge.GetZetaBlockHeight()
	if err != nil {
		log.Error().Err(err).Msg("GetZetaBlockHeight error")
		return
	}
	startBlock = (startBlock + 5) / 5 * 5
	for range ticker.C {
		bn, err := bridge.GetZetaBlockHeight()
		if err != nil {
			log.Error().Err(err).Msg("GetZetaBlockHeight error")
			continue
		}
		if bn > lastZetaBlock {
			if bn > 0 {
				for idx := 0; idx < 10; idx++ {
					go func(idx int) {
						atomic.AddInt64(&numConcurrentKeysign, 1)
						log.Info().Msgf("doing a keysign test at block %d... numConcurrentKeysign %d, idx %d", bn, numConcurrentKeysign, idx)
						testMsg := fmt.Sprintf("test message at block %d num %d", bn, idx)
						msgHash := crypto.Keccak256Hash([]byte(testMsg))
						_, err := tss.Sign(msgHash.Bytes())
						atomic.AddInt64(&totalKeysign, 1)
						if err != nil {
							log.Error().Err(err).Msg("Sign error")
						} else {
							log.Info().Msgf("sign success")
							atomic.AddInt64(&successfulKeysign, 1)
						}
						atomic.AddInt64(&numConcurrentKeysign, -1)
						log.Info().Msgf("done a keysign test at block %d numConcurrentKeysign %d, idx %d", bn, numConcurrentKeysign, idx)
					}(idx)
				}
			}
			lastZetaBlock = bn
			log.Info().Msgf("current block %d, success/total keysign: %d/%d", bn, successfulKeysign, totalKeysign)
		}
	}
}
