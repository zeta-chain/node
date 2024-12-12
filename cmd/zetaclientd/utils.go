package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	authz2 "github.com/zeta-chain/node/pkg/authz"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/authz"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

func createZetacoreClient(cfg config.Config, hotkeyPassword string, logger zerolog.Logger) (*zetacore.Client, error) {
	hotKey := cfg.AuthzHotkey

	chainIP := cfg.ZetaCoreURL

	kb, _, err := keys.GetKeyringKeybase(cfg, hotkeyPassword)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get keyring base")
	}

	granterAddress, err := sdk.AccAddressFromBech32(cfg.AuthzGranter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get granter address")
	}

	k := keys.NewKeysWithKeybase(kb, granterAddress, cfg.AuthzHotkey, hotkeyPassword)

	// All votes broadcasts to zetacore are wrapped in authz.
	// This is to ensure that the user does not need to keep their operator key online, and can use a cold key to sign votes
	signerAddress, err := k.GetAddress()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signer address")
	}

	authz.SetupAuthZSignerList(k.GetOperatorAddress().String(), signerAddress)

	client, err := zetacore.NewClient(k, chainIP, hotKey, cfg.ChainID, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create zetacore client")
	}

	return client, nil
}

// waitForBlocks waits for zetacore to be ready (i.e. producing blocks)
func waitForBlocks(ctx context.Context, zc *zetacore.Client, logger zerolog.Logger) error {
	const (
		interval = 5 * time.Second
		attempts = 15
	)

	var (
		retryCount = 0
		start      = time.Now()
	)

	task := func(ctx context.Context, t *ticker.Ticker) error {
		blockHeight, err := zc.GetBlockHeight(ctx)

		if err == nil && blockHeight > 1 {
			logger.Info().Msgf("Zetacore is ready, block height: %d", blockHeight)
			t.Stop()
			return nil
		}

		retryCount++
		if retryCount > attempts {
			return fmt.Errorf("zetacore is not ready, timeout %s", time.Since(start).String())
		}

		logger.Info().Msgf("Failed to get block number, retry: %d/%d", retryCount, attempts)
		return nil
	}

	return ticker.Run(ctx, interval, task)
}

// prepareZetacoreClient prepares the zetacore client for use.
// EXITS THE PROGRAM IF THIS NODE IS NOT AN OBSERVER.
func prepareZetacoreClient(ctx context.Context, zc *zetacore.Client, cfg *config.Config, logger zerolog.Logger) error {
	// Set grantee account number and sequence number
	if err := zc.SetAccountNumber(authz2.ZetaClientGranteeKey); err != nil {
		return errors.Wrap(err, "failed to set account number")
	}

	res, err := zc.GetNodeInfo(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get node info")
	}

	network := res.GetDefaultNodeInfo().Network
	if network != cfg.ChainID {
		logger.Warn().
			Str("got", cfg.ChainID).
			Str("want", network).
			Msg("Zetacore chain id config mismatch. Forcing update from the network")

		cfg.ChainID = network
		if err = zc.UpdateChainID(cfg.ChainID); err != nil {
			return errors.Wrap(err, "failed to update chain id")
		}
	}

	isObserver, err := isObserverNode(ctx, zc)
	switch {
	case err != nil:
		return errors.Wrap(err, "failed to check if this node is an observer")
	case !isObserver:
		addr := zc.GetKeys().GetOperatorAddress().String()
		logger.Info().Str("operator_address", addr).Msg("This node is not an observer. Exit 0")
		os.Exit(0)
	}

	return nil
}

// isObserverNode checks whether THIS node is an observer node.
func isObserverNode(ctx context.Context, zc *zetacore.Client) (bool, error) {
	observers, err := zc.GetObserverList(ctx)
	if err != nil {
		return false, errors.Wrap(err, "unable to get observers list")
	}

	operatorAddress := zc.GetKeys().GetOperatorAddress().String()

	for _, observer := range observers {
		if observer == operatorAddress {
			return true, nil
		}
	}

	return false, nil
}

func isEnvFlagEnabled(flag string) bool {
	v, _ := strconv.ParseBool(os.Getenv(flag))
	return v
}

func btcChainIDsFromContext(app *zctx.AppContext) []int64 {
	var (
		btcChains   = app.FilterChains(zctx.Chain.IsBitcoin)
		btcChainIDs = make([]int64, len(btcChains))
	)

	for i, chain := range btcChains {
		btcChainIDs[i] = chain.ID()
	}

	return btcChainIDs
}
