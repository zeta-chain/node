package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"google.golang.org/grpc"
)

func waitForZetaCore(config config.Config, logger zerolog.Logger) {
	// wait until zetacore is up
	logger.Debug().Msg("Waiting for zetacore to open 9090 port...")
	for {
		_, err := grpc.Dial(
			fmt.Sprintf("%s:9090", config.ZetaCoreURL),
			grpc.WithInsecure(),
		)
		if err != nil {
			logger.Warn().Err(err).Msg("grpc dial fail")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
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

// maskCfg sensitive fields are masked, currently only the EVM endpoints and bitcoin credentials,
//
//	other fields can be added.
func maskCfg(cfg config.Config) string {
	maskedCfg := cfg

	maskedCfg.BitcoinConfig = config.BTCConfig{
		RPCUsername: cfg.BitcoinConfig.RPCUsername,
		RPCPassword: cfg.BitcoinConfig.RPCPassword,
		RPCHost:     cfg.BitcoinConfig.RPCHost,
		RPCParams:   cfg.BitcoinConfig.RPCParams,
	}
	maskedCfg.EVMChainConfigs = map[int64]config.EVMConfig{}
	for key, val := range cfg.EVMChainConfigs {
		maskedCfg.EVMChainConfigs[key] = config.EVMConfig{
			Chain:    val.Chain,
			Endpoint: val.Endpoint,
		}
	}

	// Mask Sensitive data
	for _, chain := range maskedCfg.EVMChainConfigs {
		if chain.Endpoint == "" {
			continue
		}
		endpointURL, err := url.Parse(chain.Endpoint)
		if err != nil {
			continue
		}
		chain.Endpoint = endpointURL.Hostname()
	}

	maskedCfg.BitcoinConfig.RPCUsername = ""
	maskedCfg.BitcoinConfig.RPCPassword = ""

	return maskedCfg.String()
}
