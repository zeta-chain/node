package main

import (
	"encoding/json"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

// solanaTestKey is a local test private key for Solana
// TODO: use separate keys for each zetaclient in Solana E2E tests
// https://github.com/zeta-chain/node/issues/2614
var solanaTestKey = []uint8{
	199, 16, 63, 28, 125, 103, 131, 13, 6, 94, 68, 109, 13, 68, 132, 17,
	71, 33, 216, 51, 49, 103, 146, 241, 245, 162, 90, 228, 71, 177, 32, 199,
	31, 128, 124, 2, 23, 207, 48, 93, 141, 113, 91, 29, 196, 95, 24, 137,
	170, 194, 90, 4, 124, 113, 12, 222, 166, 209, 119, 19, 78, 20, 99, 5,
}

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
	SolanaKey           string
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
	InitCmd.Flags().StringVar(&initArgs.SolanaKey, "solana-key", "solana-key.json", "solana key file name")
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
	configData.SolanaKeyFile = initArgs.SolanaKey
	configData.ComplianceConfig = testutils.ComplianceConfigTest()

	// Save solana fee payer key
	keyFile := path.Join(rootArgs.zetaCoreHome, initArgs.SolanaKey)
	err = createSolanaKeyFile(keyFile)
	if err != nil {
		return err
	}

	// Save config file
	return config.Save(&configData, rootArgs.zetaCoreHome)
}

// createSolanaKeyFile creates a solana key json file
func createSolanaKeyFile(keyFile string) error {
	// marshal the byte array to JSON
	keyBytes, err := json.Marshal(solanaTestKey)
	if err != nil {
		return err
	}

	// create file (or overwrite if it already exists)
	file, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the key bytes to the file
	_, err = file.Write(keyBytes)
	return err
}
