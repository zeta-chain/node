package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"
	"time"
)

func keygenTss(cfg *config.Config, bridge *mc.ZetaCoreBridge, tss *mc.TSS, logger zerolog.Logger) error {
	keygenLogger := logger.With().Str("module", "keygen").Logger()
	keygenLogger.Info().Msgf("Keygen at blocknum %d", cfg.KeygenBlock)
	bn, err := bridge.GetZetaBlockHeight()
	if err != nil {
		keygenLogger.Error().Err(err).Msg("GetZetaBlockHeight RPC error")
		return err
	}
	ticker := time.NewTicker(time.Second * 1)
	lastBlock := bn
	// This is a blocking thread , it will wait for the keygen block to arrive.
	// At keygen block , it can either be success or a failure.The zetacore is update accordingly
	// This ticker waits for the keygen block to arrive
	for range ticker.C {
		currentBlock, err := bridge.GetZetaBlockHeight()
		if err != nil {
			keygenLogger.Error().Err(err).Msg("GetZetaBlockHeight RPC  error")
			return err
		}
		if currentBlock == cfg.KeygenBlock {
			log.Debug().Msgf("Trying to keygen at Block %d", currentBlock)
			break
		}
		if currentBlock > cfg.KeygenBlock {
			return errors.New("Keygen block has passed , Wait for new Keygen to be set")
		}
		// This is the only condition which triggers the debug message and causes this thread to wait
		if currentBlock > lastBlock {
			lastBlock = currentBlock
			log.Debug().Msgf("Waiting for KeygenBlock %d, Current blocknum %d", cfg.KeygenBlock, currentBlock)
		}
	}
	keygenLogger.Info().Msgf("Keygen with TSS signers %s ", cfg.KeyGenPubKeys)
	var req keygen.Request
	req = keygen.NewRequest(cfg.KeyGenPubKeys, cfg.KeygenBlock, "0.14.0")
	res, err := tss.Server.Keygen(req)
	if err != nil || res.Status != tsscommon.Success {
		keygenLogger.Error().Msgf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)
		_, err = bridge.SetTSS("", cfg.KeygenBlock, common.ReceiveStatus_Failed)
		if err != nil {
			keygenLogger.Error().Err(err).Msg("Failed to broadcast Failed TSS Vote to zetacore")
		}
		return errors.Wrap(err, fmt.Sprintf("Keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes))
	}
	// Keygen succeed! Report TSS address
	keygenLogger.Debug().Msgf("Keygen success! keygen response: %v", res)
	keygenLogger.Info().Msgf("KeyGen success ! Doing a Key-sign test")
	// KeySign can fail even if TSS keygen is successful , just loggin the error here to break out of outer loop and report TSS
	err = mc.TestKeysign(res.PubKey, tss.Server)
	if err != nil {
		keygenLogger.Error().Err(err).Msg("TestKeysign error")
	}
	keygenLogger.Info().Msgf("setting TSS pubkey: %s", res.PubKey)
	err = tss.InsertPubKey(res.PubKey)
	tss.CurrentPubkey = res.PubKey
	if err != nil {
		keygenLogger.Error().Msgf("SetPubKey fail")
		return err
	}
	keygenLogger.Info().Msgf("TSS address in hex: %s", tss.EVMAddress().Hex())
	return nil
}
