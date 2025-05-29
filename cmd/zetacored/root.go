package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"

	"cosmossdk.io/log"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	tmcfg "github.com/cometbft/cometbft/config"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ethermintclient "github.com/zeta-chain/ethermint/client"
	"github.com/zeta-chain/ethermint/crypto/hd"
	evmenc "github.com/zeta-chain/ethermint/encoding"
	"github.com/zeta-chain/ethermint/types"

	"github.com/zeta-chain/node/app"
	zetacoredconfig "github.com/zeta-chain/node/cmd/zetacored/config"
	zetamempool "github.com/zeta-chain/node/pkg/mempool"
	zevmserver "github.com/zeta-chain/node/server"
	servercfg "github.com/zeta-chain/node/server/config"
)

const EnvPrefix = "zetacore"

const DefaultTimeoutCommit = 5 * time.Second

// NewRootCmd creates a new root command for wasmd. It is called once in the
// main function.

type emptyAppOptions struct{}

func (ao emptyAppOptions) Get(_ string) interface{} { return nil }

func NewRootCmd() (*cobra.Command, types.EncodingConfig) {
	encodingConfig := app.MakeEncodingConfig()

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(app.DefaultNodeHome).
		WithKeyringOptions(hd.EthSecp256k1Option()).
		WithViper(EnvPrefix)

	rootCmd := &cobra.Command{
		Use:   zetacoredconfig.AppName,
		Short: "Zetacore Daemon (server)",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			// This needs to go after ReadFromClientConfig, as that function
			// sets the RPC client needed for SIGN_MODE_TEXTUAL. This sign mode
			// is only available if the client is online.
			if !initClientCtx.Offline {
				enabledSignModes := slices.Clone(tx.DefaultSignModes)
				enabledSignModes = append(enabledSignModes, signing.SignMode_SIGN_MODE_TEXTUAL)
				txConfigOpts := tx.ConfigOptions{
					EnabledSignModes:           enabledSignModes,
					TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(initClientCtx),
				}
				txConfig, err := tx.NewTxConfigWithOptions(
					initClientCtx.Codec,
					txConfigOpts,
				)
				if err != nil {
					return err
				}

				initClientCtx = initClientCtx.WithTxConfig(txConfig)
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, initTmConfig())
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	// need to create this app instance to get autocliopts
	tempApp := app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		"",
		0,
		evmenc.MakeConfig(),
		emptyAppOptions{},
	)
	autoCliOpts := tempApp.AutoCliOpts()
	autoCliOpts.ClientCtx = initClientCtx

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd, encodingConfig
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	return servercfg.AppConfig(zetacoredconfig.BaseDenom)
}

// initTmConfig overrides the default Tendermint config
func initTmConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	cfg.DBBackend = "pebbledb"

	return cfg
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig types.EncodingConfig) {
	ac := appCreator{
		encCfg: encodingConfig,
	}

	rootCmd.AddCommand(
		ethermintclient.ValidateChainID(
			genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
		),
		genutilcli.CollectGenTxsCmd(
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
			genutiltypes.DefaultMessageValidator,
			encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec(),
		),
		genutilcli.GenTxCmd(
			app.ModuleBasics,
			encodingConfig.TxConfig,
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
			encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec(),
		),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		AddObserverListCmd(),
		CmdParseGenesisFile(),
		GetPubKeyCmd(),
		CollectObserverInfoCmd(),
		AddrConversionCmd(),
		UpgradeHandlerVersionCmd(),
		tmcli.NewCompletionCmd(rootCmd, true),
		ethermintclient.NewTestnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),

		debug.Cmd(),
		snapshot.Cmd(ac.newApp),
	)

	zevmserver.AddCommands(
		rootCmd,
		zevmserver.NewDefaultStartOptions(ac.newApp, app.DefaultNodeHome),
		ac.appExport,
		addModuleInitFlags,
	)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		queryCommand(),
		txCommand(),
		docsCommand(),
		ethermintclient.KeyCommands(app.DefaultNodeHome),
	)

	rootCmd.AddCommand(
		confixcmd.ConfigCommand(),
	)

	// override config.toml default values, can be skipped with a flag (eg: in localnet, e2e tests etc.)
	for i, cmd := range rootCmd.Commands() {
		if cmd.Name() == "start" {
			startRunE := cmd.RunE

			cmd.RunE = func(cmd *cobra.Command, args []string) error {
				serverCtx := server.GetServerContextFromCmd(cmd)
				skipConfigDefaults := serverCtx.Viper.GetBool(FlagSkipConfigOverride)

				if !skipConfigDefaults {
					if err := overrideConfigFile(serverCtx); err != nil {
						return fmt.Errorf("failed to override config file: %w", err)
					}
				}

				return startRunE(cmd, args)
			}

			rootCmd.Commands()[i] = cmd
			break
		}
	}

	// replace the default hd-path for the key add command with Ethereum HD Path
	if err := SetEthereumHDPath(rootCmd); err != nil {
		fmt.Printf("warning: unable to set default HD path: %v\n", err)
	}
}

// overrideConfigFile handles the configuration of config.toml file,
// it either creates a new config with optimized defaults or updates an existing one.
func overrideConfigFile(serverCtx *server.Context) error {
	rootDir := serverCtx.Viper.GetString(tmcli.HomeFlag)
	configDir := filepath.Join(rootDir, "config")
	configPath := filepath.Join(configDir, "config.toml")

	// check if config file exists
	_, err := os.Stat(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to check config file %s: %w", configPath, err)
		}
		// config doesn't exist - will be created with defaults by server.InterceptConfigsPreRunHandler
		return nil
	}

	fileInfo, err := os.Stat(configPath)
	if err != nil {
		return fmt.Errorf("failed to check config file permissions: %w", err)
	}

	// if config.toml is not writable skip next steps
	if fileInfo.Mode()&os.FileMode(0200) == 0 {
		fmt.Println("warning: config.toml is not writable")
		return nil
	}

	// config exists - read and update it
	viperConfig := viper.New()
	viperConfig.SetConfigType("toml")
	viperConfig.SetConfigName("config")
	viperConfig.AddConfigPath(configDir)

	if err := viperConfig.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// update consensus timeout
	serverCtx.Config.Consensus.TimeoutCommit = DefaultTimeoutCommit

	// writeConfigFile panics on error, so we need to recover
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("warning: failed to write config file: %v\n", r)
			}
		}()
		tmcfg.WriteConfigFile(configPath, serverCtx.Config)
	}()
	return nil
}

// FlagSkipConfigOverride defines a flag to skip automatic config.toml override.
// This should be used for localnet, e2e tests etc.
const FlagSkipConfigOverride = "skip-config-override"

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
	startCmd.Flags().Bool(FlagSkipConfigOverride, false, "Skip automatic config.toml override")
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.ValidatorCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockCmd(),
		server.QueryBlocksCmd(),
		server.QueryBlockResultsCmd(),
	)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg types.EncodingConfig
}

func (ac appCreator) newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)
	maxTxs := cast.ToInt(appOpts.Get(server.FlagMempoolMaxTxs))
	if maxTxs <= 0 {
		maxTxs = zetamempool.DefaultMaxTxs
	}
	baseappOptions = append(baseappOptions, func(app *baseapp.BaseApp) {
		app.SetMempool(zetamempool.NewPriorityMempool(zetamempool.PriorityNonceWithMaxTx(maxTxs)))
	})
	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	return app.New(logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		//cosmoscmd.EncodingConfig(ac.encCfg),
		ac.encCfg,
		appOpts,
		baseappOptions...,
	)
}

// appExport is used to export the state of the application for a genesis file.
func (ac appCreator) appExport(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	var zetaApp *app.App

	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	loadLatest := false
	if height == -1 {
		loadLatest = true
	}

	zetaApp = app.New(
		logger,
		db,
		traceStore,
		loadLatest,
		map[int64]bool{},
		homePath,
		uint(1),
		ac.encCfg,
		appOpts,
	)

	// If height is -1, it means we are using the latest height.
	// For all other cases, we load the specified height from the Store
	if !loadLatest {
		if err := zetaApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return zetaApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
