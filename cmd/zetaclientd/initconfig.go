package main

import (
	"cosmossdk.io/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/config"
	zetatss "github.com/zeta-chain/node/zetaclient/tss"
)

// initializeConfigOptions is a set of CLI options for `init` command.
type initializeConfigOptions struct {
	peer               string
	publicIP           string
	logFormat          string
	logSampler         bool
	preParamsPath      string
	chainID            string
	zetacoreURL        string
	authzGranter       string
	authzHotkey        string
	level              int8
	configUpdateTicker uint64

	p2pDiagnostic              bool
	p2pDiagnosticTicker        uint64
	TSSPath                    string
	TestTSSKeySign             bool
	KeyringBackend             string
	RelayerKeyPath             string
	MaxBaseFee                 int64
	MempoolCongestionThreshold int64
}

var initializeConfigOpts initializeConfigOptions

func setupInitializeConfigOptions() {
	f, cfg := InitializeConfigCmd.Flags(), &initializeConfigOpts

	const (
		usagePeer             = "peer address e.g. /dns/tss1/tcp/6668/ipfs/16Uiu2HAmACG5DtqmQsH..."
		usageHotKey           = "hotkey for zetaclient this key is used for TSS and ZetaClient operations"
		usageLogLevel         = "log level (0:debug, 1:info, 2:warn, 3:error, 4:fatal, 5:panic)"
		usageP2PDiag          = "p2p diagnostic ticker (default: 0 means no ticker)"
		usageTicker           = "config update ticker (default: 0 means no ticker)"
		usageKeyring          = "keyring backend to use (test, file)"
		usageMaxBaseFee       = "the maximum base fee allowed to send ZetaChain transactions (0 means no limit)"
		usageMempoolThreshold = "the threshold number of unconfirmed txs in the zetacore mempool to consider it congested (0 means no threshold)"
	)

	f.StringVar(&cfg.peer, "peer", "", usagePeer)
	f.StringVar(&cfg.publicIP, "public-ip", "", "public ip address")
	f.StringVar(&cfg.preParamsPath, "pre-params", "~/preParams.json", "pre-params file path")
	f.StringVar(&cfg.chainID, "chain-id", "athens_7001-1", "chain id")
	f.StringVar(&cfg.zetacoreURL, "zetacore-url", "127.0.0.1", "zetacore node URL")
	f.StringVar(&cfg.authzGranter, "operator", "", "granter for the authorization , this should be operator address")
	f.StringVar(&cfg.authzHotkey, "hotkey", "hotkey", usageHotKey)
	f.Int8Var(&cfg.level, "log-level", int8(zerolog.InfoLevel), usageLogLevel)
	f.StringVar(&cfg.logFormat, "log-format", "json", "log format (json, test)")
	f.BoolVar(&cfg.logSampler, "log-sampler", false, "set to to true to turn on log sampling")
	f.BoolVar(&cfg.p2pDiagnostic, "p2p-diagnostic", false, "enable p2p diagnostic")
	f.Uint64Var(&cfg.p2pDiagnosticTicker, "p2p-diagnostic-ticker", 30, usageP2PDiag)
	f.Uint64Var(&cfg.configUpdateTicker, "config-update-ticker", 5, usageTicker)
	f.StringVar(&cfg.TSSPath, "tss-path", "~/.tss", "path to tss location")
	f.BoolVar(&cfg.TestTSSKeySign, "test-tss", false, "set to to true to run a check for TSS keysign on startup")
	f.StringVar(&cfg.KeyringBackend, "keyring-backend", string(config.KeyringBackendTest), usageKeyring)
	f.StringVar(&cfg.RelayerKeyPath, "relayer-key-path", "~/.zetacored/relayer-keys", "path to relayer keys")
	f.Int64Var(&cfg.MaxBaseFee, "max-base-fee", 0, usageMaxBaseFee)
	f.Int64Var(
		&cfg.MempoolCongestionThreshold,
		"mempool-threshold",
		config.DefaultMempoolCongestionThreshold,
		usageMempoolThreshold,
	)
}

// InitializeConfig creates new config for zetaclientd and saves it to the config file.
func InitializeConfig(_ *cobra.Command, _ []string) error {
	// Create new config struct
	configData := config.New(true)
	opts := &initializeConfigOpts

	// Validate Peer
	// e.g. /ip4/172.0.2.1/tcp/6668/p2p/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp
	if opts.peer != "" {
		if _, err := zetatss.MultiAddressFromString(opts.peer); err != nil {
			return errors.Wrap(err, "invalid peer address")
		}
	}

	// Populate new struct with cli arguments
	configData.Peer = initializeConfigOpts.peer
	configData.PublicIP = opts.publicIP
	configData.PreParamsPath = opts.preParamsPath
	configData.ChainID = opts.chainID
	configData.ZetaCoreURL = opts.zetacoreURL
	configData.AuthzHotkey = opts.authzHotkey
	configData.AuthzGranter = opts.authzGranter
	configData.LogLevel = opts.level
	configData.LogFormat = opts.logFormat
	configData.LogSampler = opts.logSampler
	configData.P2PDiagnostic = opts.p2pDiagnostic
	configData.TssPath = opts.TSSPath
	configData.TestTssKeysign = opts.TestTSSKeySign
	configData.P2PDiagnosticTicker = opts.p2pDiagnosticTicker
	configData.ConfigUpdateTicker = opts.configUpdateTicker
	configData.KeyringBackend = config.KeyringBackend(initializeConfigOpts.KeyringBackend)
	configData.RelayerKeyPath = opts.RelayerKeyPath
	configData.MaxBaseFee = opts.MaxBaseFee
	configData.MempoolCongestionThreshold = opts.MempoolCongestionThreshold
	configData.ComplianceConfig = sample.ComplianceConfig()
	configData.FeatureFlags = sample.FeatureFlags()

	// Save config file
	return config.Save(&configData, globalOpts.ZetacoreHome)
}
