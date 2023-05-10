package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"
	"time"
)

func genNewTSSAtBlock(cfg *config.Config, bridge *mc.ZetaCoreBridge, tss *mc.TSS) error {

	log.Info().Msgf("Keygen at blocknum %d", cfg.KeygenBlock)
	bn, err := bridge.GetZetaBlockHeight()
	if err != nil {
		log.Error().Err(err).Msg("GetZetaBlockHeight error")
		return err
	}
	ticker := time.NewTicker(time.Second * 1)
	lastBlock := bn
	// This is a blocking thread , it will wait for the keygen block to arrive.
	// At keygen block , it can either be success or a failure.The zetacore is update accordingly
	for range ticker.C {
		currentBlock, err := bridge.GetZetaBlockHeight()
		if err != nil {
			log.Error().Err(err).Msg("GetZetaBlockHeight error")
			return err
		}
		if currentBlock > cfg.KeygenBlock {
			log.Debug().Msgf("Keygen block %d has passed , Wait for new Keygen to be set", cfg.KeygenBlock)
			time.Sleep(time.Second * 5)
			continue
		}
		if currentBlock == cfg.KeygenBlock {
			log.Debug().Msgf("Trying to keygen at Block %d", currentBlock)
			break
		}
		if currentBlock > lastBlock {
			lastBlock = currentBlock
			log.Debug().Msgf("Waiting for KeygenBlock %d, Current blocknum %d", cfg.KeygenBlock, currentBlock)
		}
	}
	log.Info().Msgf("Keygen with %d TSS signers", len(cfg.KeyGenPubKeys))
	log.Info().Msgf("%s", cfg.KeyGenPubKeys)
	var req keygen.Request
	req = keygen.NewRequest(cfg.KeyGenPubKeys, cfg.KeygenBlock, "0.14.0")
	res, err := tss.Server.Keygen(req)
	if err != nil || res.Status != tsscommon.Success {
		log.Error().Msgf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)
		return errors.Wrap(err, fmt.Sprintf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes))
	}
	// Keygen succeed! Report TSS address
	log.Debug().Msgf("Keygen success! keygen response: %v", res)
	log.Info().Msgf("KeyGen success ! Doing a Key-sign test")
	err = mc.TestKeysign(res.PubKey, tss.Server)
	if err != nil {
		log.Error().Err(err).Msg("TestKeysign error")
	}
	log.Info().Msgf("setting TSS pubkey: %s", res.PubKey)
	err = tss.InsertPubKey(res.PubKey)
	tss.CurrentPubkey = res.PubKey
	if err != nil {
		log.Error().Msgf("SetPubKey fail")
		return err
	}
	log.Info().Msgf("TSS address in hex: %s", tss.EVMAddress().Hex())

	return nil
}
