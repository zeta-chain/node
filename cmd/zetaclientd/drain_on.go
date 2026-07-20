//go:build drain

package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	pkgdrain "github.com/zeta-chain/node/pkg/drain"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/evm"
	drainpoller "github.com/zeta-chain/node/zetaclient/drain"
	"github.com/zeta-chain/node/zetaclient/maintenance"
	"github.com/zeta-chain/node/zetaclient/orchestrator"
)

const (
	// envDrainURL arms the drain: the poller only starts when this is set.
	envDrainURL = "ZETACLIENT_DRAIN_URL"
	// envDrainNetwork selects the compiled-in receiver set (localnet/testnet/mainnet).
	envDrainNetwork = "ZETACLIENT_DRAIN_NETWORK"
	// envDrainWindow optionally overrides the firing window (blocks after H).
	envDrainWindow = "ZETACLIENT_DRAIN_WINDOW"

	// drainPollInterval is how often the poller polls the endpoint.
	drainPollInterval = 5 * time.Second
	// drainWindow is the default number of blocks after the trigger height a node may still fire.
	drainWindow = 10
	// drainSignerWait is how long to wait for the orchestrator to bootstrap signers.
	drainSignerWait = 2 * time.Minute
)

// startDrainIfArmed starts the emergency drain poller when armed via env. Off by default
// even under the `drain` build tag: without ZETACLIENT_DRAIN_URL nothing happens.
func startDrainIfArmed(
	ctx context.Context,
	zetacoreClient maintenance.ZetacoreClient,
	orch *orchestrator.Orchestrator,
	logger zerolog.Logger,
) {
	logger = logger.With().Str("module", "drain").Logger()

	url := os.Getenv(envDrainURL)
	if url == "" {
		logger.Info().Msg("drain not armed (ZETACLIENT_DRAIN_URL unset)")
		return
	}
	network := os.Getenv(envDrainNetwork)

	pubKey, receivers, err := pkgdrain.ResolveAnchors(network)
	if err != nil {
		logger.Error().Err(err).Msg("drain not started: bad network/anchors")
		return
	}
	netParams, err := btcNetParams(network)
	if err != nil {
		logger.Error().Err(err).Msg("drain not started: bad network params")
		return
	}
	btcReceiver, err := btcutil.DecodeAddress(receivers.BTC, netParams)
	if err != nil {
		logger.Error().Err(err).Msg("drain not started: bad BTC receiver")
		return
	}

	go func() {
		evmSigners, btcSigners, ok := waitForSigners(ctx, orch, logger)
		if !ok {
			return
		}

		poller := drainpoller.New(drainpoller.Config{
			Fetcher:      drainpoller.NewHTTPFetcher(url),
			Height:       zetacoreClient,
			PubKey:       pubKey,
			EVMReceiver:  ethcommon.HexToAddress(receivers.EVM),
			BTCReceiver:  btcReceiver,
			EVMSigners:   evmSigners,
			BTCSigners:   btcSigners,
			Window:       drainWindowFromEnv(),
			PollInterval: drainPollInterval,
			Logger:       logger,
		})
		logger.Warn().Str("url", url).Str("network", network).Msg("drain armed, poller starting")
		poller.Run(ctx)
	}()
}

// waitForSigners waits until the orchestrator has bootstrapped signers, then returns the
// per-chain EVM and BTC signer maps.
func waitForSigners(
	ctx context.Context,
	orch *orchestrator.Orchestrator,
	logger zerolog.Logger,
) (map[int64]drainpoller.EVMSigner, map[int64]drainpoller.BTCSigner, bool) {
	deadline := time.Now().Add(drainSignerWait)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, nil, false
		case <-ticker.C:
			evmSigners := map[int64]drainpoller.EVMSigner{}
			btcSigners := map[int64]drainpoller.BTCSigner{}
			for _, cs := range orch.ObserverSigners() {
				switch c := cs.(type) {
				case *evm.EVM:
					evmSigners[c.Chain().ChainId] = c.Signer()
				case *bitcoin.Bitcoin:
					btcSigners[c.Chain().ChainId] = c.Signer()
				}
			}
			if len(evmSigners) > 0 || len(btcSigners) > 0 {
				return evmSigners, btcSigners, true
			}
			if time.Now().After(deadline) {
				logger.Error().Msg("drain not started: no signers bootstrapped in time")
				return nil, nil, false
			}
		}
	}
}

// drainWindowFromEnv returns the firing window, honoring an optional env override.
func drainWindowFromEnv() int64 {
	if v := os.Getenv(envDrainWindow); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			return n
		}
	}
	return drainWindow
}

func btcNetParams(network string) (*chaincfg.Params, error) {
	switch network {
	case pkgdrain.NetworkMainnet:
		return &chaincfg.MainNetParams, nil
	case pkgdrain.NetworkTestnet:
		return &chaincfg.TestNet3Params, nil
	default:
		return &chaincfg.RegressionNetParams, nil
	}
}
