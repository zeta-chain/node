package zetacore

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/authz"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// This file contains some high level logic for creating a zetacore client
// when starting zetaclientd in cmd/zetaclientd/start.go

// NewFromConfig creates a new client from the given config.
// It also makes sure that the zetacore (i.e. the node) is ready to be used.
func NewFromConfig(
	ctx context.Context,
	cfg *config.Config,
	hotkeyPassword string,
	logger zerolog.Logger,
) (*Client, error) {
	kb, _, err := keys.GetKeyringKeybase(*cfg, hotkeyPassword)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get keyring base")
	}

	granterAddress, err := sdk.AccAddressFromBech32(cfg.AuthzGranter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get granter address")
	}

	k := keys.NewKeysWithKeybase(kb, granterAddress, cfg.AuthzHotkey, hotkeyPassword)

	// All votes broadcasts to zetacore are wrapped in authz.
	// This is to ensure that the user does not need to keep their operator key online,
	// and can use a cold key to sign votes
	signerAddress, err := k.GetAddress()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signer address")
	}

	authz.SetupAuthZSignerList(k.GetOperatorAddress().String(), signerAddress)

	// Create client
	client, err := NewClient(k, *cfg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the client")
	}

	// Make sure that the node produces blocks
	if err = ensureBlocksProduction(ctx, client); err != nil {
		return nil, errors.Wrap(err, "zetacore unavailable")
	}

	// Prepare the client
	if err = prepareZetacoreClient(ctx, client, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to prepare the client")
	}

	return client, nil
}

// ensureBlocksProduction waits for zetacore to be ready (i.e. producing blocks)
func ensureBlocksProduction(ctx context.Context, zc *Client) error {
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
			zc.logger.Info().Int64(logs.FieldBlock, blockHeight).Msg("zetacore is ready")
			t.Stop()
			return nil
		}

		retryCount++
		if retryCount > attempts {
			return fmt.Errorf("zetacore is not ready, timeout %s", time.Since(start).String())
		}

		zc.logger.Info().
			Int("current_attempt", retryCount).
			Uint("max_attempts", attempts).
			Msg("failed to get block number")
		return nil
	}

	return ticker.Run(ctx, interval, task)
}

// prepareZetacoreClient prepares the zetacore client for use.
func prepareZetacoreClient(ctx context.Context, zc *Client, cfg *config.Config) error {
	res, err := zc.GetNodeInfo(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get node info")
	}

	network := res.GetDefaultNodeInfo().Network
	if network != cfg.ChainID {
		zc.logger.Warn().
			Str("got", cfg.ChainID).
			Str("want", network).
			Msg("zetacore chain id config mismatch; forcing update from the network")

		cfg.ChainID = network
		if err = zc.UpdateChainID(cfg.ChainID); err != nil {
			return errors.Wrap(err, "failed to update chain id")
		}
	}

	return nil
}
