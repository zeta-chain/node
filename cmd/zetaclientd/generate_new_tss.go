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
	height := cfg.KeygenBlock
	pubKeys := cfg.KeyGenPubKeys
	log.Info().Msgf("Keygen at blocknum %d", height)
	bn, err := bridge.GetZetaBlockHeight()
	if err != nil {
		log.Error().Err(err).Msg("GetZetaBlockHeight error")
		return err
	}
	if bn+3 > height {
		return fmt.Errorf(fmt.Sprintf("Keygen at Blocknum %d, but current blocknum %d , Too late to take part in this keygen. Try again at a later block", height, bn))
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
		if currentBlock == height {
			log.Debug().Msgf("Trying to keygen at Block %d", currentBlock)
			break
		}
		if currentBlock > lastBlock {
			lastBlock = currentBlock
			log.Debug().Msgf("Waiting for KeygenBlock %d, Current blocknum %d", height, currentBlock)
		}
	}
	log.Info().Msgf("Keygen with %d TSS signers", len(pubKeys))
	log.Info().Msgf("%s", pubKeys)
	var req keygen.Request
	req = keygen.NewRequest(pubKeys, height, "0.14.0")
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

	if err != nil {
		log.Error().Msgf("SetPubKey fail")
		return err
	}
	tss.CurrentPubkey = res.PubKey // this is only needed for version 0.13.0 leaderless keysign
	tss.Signers = pubKeys
	log.Info().Msgf("TSS address in hex: %s", tss.EVMAddress().Hex())

	return nil
}
