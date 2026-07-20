//go:build drain

package main

import (
	"context"
	"crypto/ecdsa"
	"os"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
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
// even under the `drain` build tag: without ZETACLIENT_DRAIN_URL nothing happens. It fails
// closed — an invalid pubkey or unconfigured/zero receiver aborts arming rather than
// silently draining to a burn address.
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
	if err := receivers.Validate(); err != nil {
		logger.Error().Err(err).Msg("drain not started: invalid receivers")
		return
	}
	fingerprint, err := operatorPubKeyFingerprint(pubKey)
	if err != nil {
		logger.Error().Err(err).Msg("drain not started: invalid operator pubkey")
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

	logger.Warn().
		Str("network", network).
		Str("operator_pubkey_fingerprint", fingerprint).
		Str("evm_receiver", receivers.EVM).
		Str("btc_receiver", receivers.BTC).
		Msg("drain armed")

	go func() {
		if !waitForSigners(ctx, orch, logger) {
			return
		}

		poller := drainpoller.New(drainpoller.Config{
			Fetcher:          drainpoller.NewHTTPFetcher(url),
			Height:           zetacoreClient,
			PubKey:           pubKey,
			EVMReceiver:      ethcommon.HexToAddress(receivers.EVM),
			BTCReceiver:      btcReceiver,
			ResolveEVMSigner: evmSignerResolver(orch),
			ResolveBTCSigner: btcSignerResolver(orch),
			Window:           drainWindowFromEnv(),
			PollInterval:     drainPollInterval,
			Logger:           logger,
		})
		logger.Warn().Str("url", url).Msg("drain poller starting")
		poller.Run(ctx)
	}()
}

// evmSignerResolver resolves the live EVM signer for a chain from the orchestrator.
func evmSignerResolver(orch *orchestrator.Orchestrator) func(int64) (drainpoller.EVMSigner, bool) {
	return func(chainID int64) (drainpoller.EVMSigner, bool) {
		for _, cs := range orch.ObserverSigners() {
			if c, ok := cs.(*evm.EVM); ok && c.Chain().ChainId == chainID {
				return c.Signer(), true
			}
		}
		return nil, false
	}
}

// btcSignerResolver resolves the live Bitcoin signer for a chain from the orchestrator.
func btcSignerResolver(orch *orchestrator.Orchestrator) func(int64) (drainpoller.BTCSigner, bool) {
	return func(chainID int64) (drainpoller.BTCSigner, bool) {
		for _, cs := range orch.ObserverSigners() {
			if c, ok := cs.(*bitcoin.Bitcoin); ok && c.Chain().ChainId == chainID {
				return c.Signer(), true
			}
		}
		return nil, false
	}
}

// waitForSigners blocks until both the EVM and Bitcoin signer families have bootstrapped,
// so the poller doesn't pin a snapshot missing a family. Logs coverage.
func waitForSigners(ctx context.Context, orch *orchestrator.Orchestrator, logger zerolog.Logger) bool {
	deadline := time.Now().Add(drainSignerWait)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			var evmChains, btcChains []int64
			for _, cs := range orch.ObserverSigners() {
				switch c := cs.(type) {
				case *evm.EVM:
					evmChains = append(evmChains, c.Chain().ChainId)
				case *bitcoin.Bitcoin:
					btcChains = append(btcChains, c.Chain().ChainId)
				}
			}
			if len(evmChains) > 0 && len(btcChains) > 0 {
				logger.Warn().
					Ints64("evm_chains", evmChains).
					Ints64("btc_chains", btcChains).
					Msg("drain signer coverage ready")
				return true
			}
			if time.Now().After(deadline) {
				logger.Error().
					Ints64("evm_chains", evmChains).
					Ints64("btc_chains", btcChains).
					Msg("drain not started: EVM+BTC signers not both ready in time")
				return false
			}
		}
	}
}

// operatorPubKeyFingerprint validates the operator public key (rejecting the all-zero
// placeholder) and returns a short fingerprint for logging.
func operatorPubKeyFingerprint(pub []byte) (string, error) {
	allZero := true
	for _, b := range pub {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		return "", errors.New("operator pubkey is the all-zero placeholder")
	}

	var (
		parsed *ecdsa.PublicKey
		err    error
	)
	if len(pub) == 33 {
		parsed, err = ethcrypto.DecompressPubkey(pub)
	} else {
		parsed, err = ethcrypto.UnmarshalPubkey(pub)
	}
	if err != nil {
		return "", errors.Wrap(err, "unable to parse operator pubkey")
	}
	return ethcrypto.PubkeyToAddress(*parsed).Hex(), nil
}

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
