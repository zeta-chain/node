package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"
)

func keygenTss(cfg *config.Config, tss *mc.TSS, logger zerolog.Logger) (error, bool) {
	keygenLogger := logger.With().Str("module", "keygen").Logger()
	keygenLogger.Info().Msgf("Keygen at blocknum %d , TSS signers %s ", cfg.KeygenBlock, cfg.KeyGenPubKeys)
	reportKeyGenFail := true
	var req keygen.Request
	req = keygen.NewRequest(cfg.KeyGenPubKeys, cfg.KeygenBlock, "0.14.0")
	res, err := tss.Server.Keygen(req)
	if err != nil || res.Status != tsscommon.Success {
		keygenLogger.Error().Msgf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)
		return errors.Wrap(err, fmt.Sprintf("Keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)), true
	}
	// Keygen succeed! Report TSS address
	keygenLogger.Debug().Msgf("Keygen success! keygen response: %v", res)

	// Do not report Failed Keygen to ZetaCore after this point , TSS test may fail even if keygen is successful
	reportKeyGenFail = false
	keygenLogger.Info().Msgf("KeyGen success ! Doing a Key-sign test")
	// KeySign can fail even if TSS keygen is successful , just logging the error here to break out of outer loop and report TSS
	err = mc.TestKeysign(res.PubKey, tss.Server)
	if err != nil {
		keygenLogger.Error().Err(err).Msg("TestKeysign error")
	}
	keygenLogger.Info().Msgf("setting TSS pubkey: %s", res.PubKey)
	err = tss.InsertPubKey(res.PubKey)
	tss.CurrentPubkey = res.PubKey
	if err != nil {
		keygenLogger.Error().Msgf("SetPubKey fail")
		return err, reportKeyGenFail
	}
	keygenLogger.Info().Msgf("TSS address in hex: %s", tss.EVMAddress().Hex())
	return nil, reportKeyGenFail
}
