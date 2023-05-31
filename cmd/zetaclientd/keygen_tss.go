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

func keygenTss(cfg *config.Config, tss *mc.TSS, logger zerolog.Logger) error {
	keygenLogger := logger.With().Str("module", "keygen").Logger()
	keygenLogger.Info().Msgf("Keygen at blocknum %d , TSS signers %s ", cfg.KeygenBlock, cfg.KeyGenPubKeys)
	var req keygen.Request
	req = keygen.NewRequest(cfg.KeyGenPubKeys, cfg.KeygenBlock, "0.14.0")
	res, err := tss.Server.Keygen(req)
	if res.Status != tsscommon.Success || res.PubKey == "" {
		keygenLogger.Error().Msgf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)
		return errors.New(fmt.Sprintf("Keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes))
	}
	if err != nil {
		keygenLogger.Error().Msgf("keygen fail: reason %s ", err.Error())
		return err
	}
	tss.CurrentPubkey = res.PubKey

	// Keygen succeed! Report TSS address
	keygenLogger.Debug().Msgf("Keygen success! keygen response: %v", res)
	return nil
}

func SetTSSPubKey(tss *mc.TSS, logger zerolog.Logger) error {
	keygenLogger := logger.With().Str("module", "set-keygen").Logger()
	keygenLogger.Info().Msgf("setting TSS pubkey: %s", tss.CurrentPubkey)
	err := tss.InsertPubKey(tss.CurrentPubkey)
	if err != nil {
		keygenLogger.Error().Msgf("SetPubKey fail")
		return err
	}
	keygenLogger.Info().Msgf("TSS address in hex: %s", tss.EVMAddress().Hex())
	return nil

}
func TestTSS(tss *mc.TSS, logger zerolog.Logger) error {
	keygenLogger := logger.With().Str("module", "test-keygen").Logger()
	keygenLogger.Info().Msgf("KeyGen success ! Doing a Key-sign test")
	// KeySign can fail even if TSS keygen is successful , just logging the error here to break out of outer loop and report TSS
	err := mc.TestKeysign(tss.CurrentPubkey, tss.Server)
	if err != nil {
		return err
	}
	return nil
}
