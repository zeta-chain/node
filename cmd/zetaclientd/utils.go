package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zeta-chain/node/zetaclient/authz"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

func createAuthzSigner(granter string, grantee sdk.AccAddress) {
	authz.SetupAuthZSignerList(granter, grantee)
}

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

	client, err := zetacore.NewClient(k, chainIP, hotKey, cfg.ChainID, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create zetacore client")
	}

	return client, nil
}

func waitForZetaCore(config config.Config, logger zerolog.Logger) {
	const (
		port  = 9090
		retry = 5 * time.Second
	)

	var (
		url = fmt.Sprintf("%s:%d", config.ZetaCoreURL, port)
		opt = grpc.WithTransportCredentials(insecure.NewCredentials())
	)

	// wait until zetacore is up
	logger.Debug().Msgf("Waiting for zetacore to open %d port...", port)

	for {
		if _, err := grpc.Dial(url, opt); err != nil {
			logger.Warn().Err(err).Msg("grpc dial fail")
			time.Sleep(retry)
		} else {
			break
		}
	}
}

func waitForZetacoreToCreateBlocks(ctx context.Context, zc interfaces.ZetacoreClient, logger zerolog.Logger) error {
	const (
		interval = 5 * time.Second
		attempts = 15
	)

	var (
		retryCount = 0
		start      = time.Now()
	)

	for {
		blockHeight, err := zc.GetBlockHeight(ctx)
		if err == nil && blockHeight > 1 {
			logger.Info().Msgf("Zeta block height: %d", blockHeight)
			return nil
		}

		retryCount++
		if retryCount > attempts {
			return fmt.Errorf("zetacore is not ready, timeout %s", time.Since(start).String())
		}

		logger.Info().Msgf("Failed to get block number, retry : %d/%d", retryCount, attempts)
		time.Sleep(interval)
	}
}

func validatePeer(seedPeer string) error {
	parsedPeer := strings.Split(seedPeer, "/")

	if len(parsedPeer) < 7 {
		return errors.New("seed peer missing IP or ID or both, seed: " + seedPeer)
	}

	seedIP := parsedPeer[2]
	seedID := parsedPeer[6]

	if net.ParseIP(seedIP) == nil {
		return errors.New("invalid seed IP address format, seed: " + seedPeer)
	}

	if len(seedID) == 0 {
		return errors.New("seed id is empty, seed: " + seedPeer)
	}

	return nil
}
