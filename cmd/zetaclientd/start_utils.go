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
	logger.Debug().Msg("Waiting for ZetaCore to open 9090 port...")
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
	//Perform deep copy to avoid modifying original config
	maskedCfg := config.Config{
		Peer:                cfg.Peer,
		PublicIP:            cfg.PublicIP,
		LogFormat:           cfg.LogFormat,
		LogLevel:            cfg.LogLevel,
		LogSampler:          cfg.LogSampler,
		PreParamsPath:       cfg.PreParamsPath,
		ZetaCoreHome:        cfg.ZetaCoreHome,
		ChainID:             cfg.ChainID,
		ZetaCoreURL:         cfg.ZetaCoreURL,
		AuthzGranter:        cfg.AuthzGranter,
		AuthzHotkey:         cfg.AuthzHotkey,
		P2PDiagnostic:       cfg.P2PDiagnostic,
		ConfigUpdateTicker:  cfg.ConfigUpdateTicker,
		P2PDiagnosticTicker: cfg.P2PDiagnosticTicker,
		TssPath:             cfg.TssPath,
		TestTssKeysign:      cfg.TestTssKeysign,
		KeyringBackend:      cfg.KeyringBackend,
		HsmMode:             cfg.HsmMode,
		HsmHotKey:           cfg.HsmHotKey,
	}

	maskedCfg.BitcoinConfig = config.BTCConfig{
		RPCUsername: cfg.BitcoinConfig.RPCUsername,
		RPCPassword: cfg.BitcoinConfig.RPCPassword,
		RPCHost:     cfg.BitcoinConfig.RPCHost,
		RPCParams:   cfg.BitcoinConfig.RPCParams,
	}

	restrictedAddresses := make([]string, len(cfg.ComplianceConfig.RestrictedAddresses))
	copy(restrictedAddresses, cfg.ComplianceConfig.RestrictedAddresses)
	maskedCfg.ComplianceConfig = config.ComplianceConfig{
		LogPath:             cfg.ComplianceConfig.LogPath,
		RestrictedAddresses: restrictedAddresses,
	}

	maskedCfg.EVMChainConfigs = map[int64]config.EVMConfig{}
	for key, val := range cfg.EVMChainConfigs {
		endpoints := make([]string, len(val.Endpoints))
		copy(endpoints, val.Endpoints)

		maskedCfg.EVMChainConfigs[key] = config.EVMConfig{
			Chain:     val.Chain,
			Endpoints: endpoints,
		}
	}

	// Mask Sensitive data
	for _, chain := range maskedCfg.EVMChainConfigs {
		if len(chain.Endpoints) == 0 {
			continue
		}
		for i, endpoint := range chain.Endpoints {
			endpointURL, err := url.Parse(endpoint)
			if err != nil {
				continue
			}
			chain.Endpoints[i] = endpointURL.Hostname()
		}
	}
	maskedCfg.BitcoinConfig.RPCUsername = ""
	maskedCfg.BitcoinConfig.RPCPassword = ""

	return maskedCfg.String()
}
