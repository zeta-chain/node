package main

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Configuration",
	RunE:  Initialize,
}

var initArgs = initArguments{}

type initArguments struct {
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

	p2pDiagnostic       bool
	p2pDiagnosticTicker uint64
	TssPath             string
	TestTssKeysign      bool
	KeyringBackend      string
	HsmMode             bool
	HsmHotKey           string
}

func init() {
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(VersionCmd)

	InitCmd.Flags().
		StringVar(&initArgs.peer, "peer", "", "peer address, e.g. /dns/tss1/tcp/6668/ipfs/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp")
	InitCmd.Flags().StringVar(&initArgs.publicIP, "public-ip", "", "public ip address")
	InitCmd.Flags().StringVar(&initArgs.preParamsPath, "pre-params", "~/preParams.json", "pre-params file path")
	InitCmd.Flags().StringVar(&initArgs.chainID, "chain-id", "athens_7001-1", "chain id")
	InitCmd.Flags().StringVar(&initArgs.zetacoreURL, "zetacore-url", "127.0.0.1", "zetacore node URL")
	InitCmd.Flags().
		StringVar(&initArgs.authzGranter, "operator", "", "granter for the authorization , this should be operator address")
	InitCmd.Flags().
		StringVar(&initArgs.authzHotkey, "hotkey", "hotkey", "hotkey for zetaclient this key is used for TSS and ZetaClient operations")
	InitCmd.Flags().
		Int8Var(&initArgs.level, "log-level", int8(zerolog.InfoLevel), "log level (0:debug, 1:info, 2:warn, 3:error, 4:fatal, 5:panic , 6: NoLevel , 7: Disable)")
	InitCmd.Flags().StringVar(&initArgs.logFormat, "log-format", "json", "log format (json, test)")
	InitCmd.Flags().BoolVar(&initArgs.logSampler, "log-sampler", false, "set to to true to turn on log sampling")
	InitCmd.Flags().BoolVar(&initArgs.p2pDiagnostic, "p2p-diagnostic", false, "enable p2p diagnostic")
	InitCmd.Flags().
		Uint64Var(&initArgs.p2pDiagnosticTicker, "p2p-diagnostic-ticker", 30, "p2p diagnostic ticker (default: 0 means no ticker)")
	InitCmd.Flags().
		Uint64Var(&initArgs.configUpdateTicker, "config-update-ticker", 5, "config update ticker (default: 0 means no ticker)")
	InitCmd.Flags().StringVar(&initArgs.TssPath, "tss-path", "~/.tss", "path to tss location")
	InitCmd.Flags().
		BoolVar(&initArgs.TestTssKeysign, "test-tss", false, "set to to true to run a check for TSS keysign on startup")
	InitCmd.Flags().
		StringVar(&initArgs.KeyringBackend, "keyring-backend", string(config.KeyringBackendTest), "keyring backend to use (test, file)")
	InitCmd.Flags().BoolVar(&initArgs.HsmMode, "hsm-mode", false, "enable hsm signer, default disabled")
	InitCmd.Flags().
		StringVar(&initArgs.HsmHotKey, "hsm-hotkey", "hsm-hotkey", "name of hotkey associated with hardware security module")
}

func Initialize(_ *cobra.Command, _ []string) error {
	err := setHomeDir()
	if err != nil {
		return err
	}

	//Create new config struct
	configData := config.New(true)

	//Validate Peer eg. /ip4/172.0.2.1/tcp/6668/p2p/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp
	if len(initArgs.peer) != 0 {
		err := validatePeer(initArgs.peer)
		if err != nil {
			return err
		}
	}

	//Populate new struct with cli arguments
	configData.Peer = initArgs.peer
	configData.PublicIP = initArgs.publicIP
	configData.PreParamsPath = initArgs.preParamsPath
	configData.ChainID = initArgs.chainID
	configData.ZetaCoreURL = initArgs.zetacoreURL
	configData.AuthzHotkey = initArgs.authzHotkey
	configData.AuthzGranter = initArgs.authzGranter
	configData.LogLevel = initArgs.level
	configData.LogFormat = initArgs.logFormat
	configData.LogSampler = initArgs.logSampler
	configData.P2PDiagnostic = initArgs.p2pDiagnostic
	configData.TssPath = initArgs.TssPath
	configData.P2PDiagnosticTicker = initArgs.p2pDiagnosticTicker
	configData.ConfigUpdateTicker = initArgs.configUpdateTicker
	configData.KeyringBackend = config.KeyringBackend(initArgs.KeyringBackend)
	configData.HsmMode = initArgs.HsmMode
	configData.HsmHotKey = initArgs.HsmHotKey
	configData.ComplianceConfig = testutils.ComplianceConfigTest()

	//Save config file
	return config.Save(&configData, rootArgs.zetaCoreHome)
}
