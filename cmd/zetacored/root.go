package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"

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
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	cosmosevmcmd "github.com/cosmos/evm/client"
	"github.com/cosmos/evm/crypto/hd"
	cosmosevmserverconfig "github.com/cosmos/evm/server/config"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/app"
	zetacoredconfig "github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/pkg/chains"
	zevmserver "github.com/zeta-chain/node/server"
	zetaserverconfig "github.com/zeta-chain/node/server/config"
)

const EnvPrefix = "zetacore"

// NewRootCmd creates a new root command for zetacored. It is called once in the
// main function.

type emptyAppOptions struct{}

func (ao emptyAppOptions) Get(_ string) interface{} { return nil }

func NewRootCmd() *cobra.Command {
	// need to create this app instance to get autocliopts
	tempApp := app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		"",
		0,
		zetaserverconfig.DefaultEVMChainID, // should be ok to use default just for temp app
		emptyAppOptions{},
	)

	// should be ok to use default just for temp app
	encodingConfig := app.MakeEncodingConfig(zetaserverconfig.DefaultEVMChainID)

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

			// TODO https://github.com/zeta-chain/node/issues/4078
			// need to check about evm chain id, getting it like this is generally fine, but some commands
			// like docs are halting because of it
			if initClientCtx.ChainID == "" {
				return nil
			}

			zetachain, err := chains.ZetaChainFromCosmosChainID(initClientCtx.ChainID)
			if err != nil {
				return err
			}

			//#nosec G115 chain id won't exceed uint64
			customAppTemplate, customAppConfig := InitAppConfig(zetacoredconfig.BaseDenom, uint64(zetachain.ChainId))

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, initTmConfig())
		},
	}

	initRootCmd(rootCmd, encodingConfig)
	rootCmd.AddCommand(
		confixcmd.ConfigCommand(),
	)

	autoCliOpts := tempApp.AutoCliOpts()
	autoCliOpts.ClientCtx = initClientCtx

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd
}

// InitAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func InitAppConfig(denom string, evmChainID uint64) (string, interface{}) {
	ethCfg := evmtypes.DefaultChainConfig(evmChainID)

	configurator := evmtypes.NewEVMConfigurator()
	err := configurator.
		WithExtendedEips(zetacoredconfig.CosmosEVMActivators).
		WithChainConfig(ethCfg).
		WithEVMCoinInfo(evmtypes.EvmCoinInfo{
			Denom:         denom,
			ExtendedDenom: denom,
			DisplayDenom:  denom,
			Decimals:      18,
		}).
		Configure()
	if err != nil {
		panic(err)
	}

	type CustomAppConfig struct {
		serverconfig.Config

		EVM     cosmosevmserverconfig.EVMConfig
		JSONRPC cosmosevmserverconfig.JSONRPCConfig
		TLS     cosmosevmserverconfig.TLSConfig
	}

	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()
	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In this application, we set the min gas prices to 1.
	srvCfg.MinGasPrices = "1" + denom

	evmCfg := cosmosevmserverconfig.DefaultEVMConfig()
	evmCfg.EVMChainID = evmChainID

	customAppConfig := CustomAppConfig{
		Config:  *srvCfg,
		EVM:     *evmCfg,
		JSONRPC: *cosmosevmserverconfig.DefaultJSONRPCConfig(),
		TLS:     *cosmosevmserverconfig.DefaultTLSConfig(),
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate +
		cosmosevmserverconfig.DefaultEVMConfigTemplate

	return customAppTemplate, customAppConfig
}

// initTmConfig overrides the default Tendermint config
func initTmConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	cfg.DBBackend = "pebbledb"

	return cfg
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig testutil.TestEncodingConfig) {
	ac := appCreator{
		encCfg: encodingConfig,
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
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
		cosmosevmcmd.KeyCommands(app.DefaultNodeHome, true),
	)

	// replace the default hd-path for the key add command with Ethereum HD Path
	if err := SetEthereumHDPath(rootCmd); err != nil {
		fmt.Printf("warning: unable to set default HD path: %v\n", err)
	}
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
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
	encCfg testutil.TestEncodingConfig
}

const DefaultMaxTxs = 3000

type PrivValidatorState struct {
	Height string `json:"height"`
	Round  int    `json:"round"`
	Step   int    `json:"step"`
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
		maxTxs = DefaultMaxTxs
	}
	signerExtractor := app.NewEthSignerExtractionAdapter(mempool.NewDefaultSignerExtractionAdapter())
	mpool := mempool.NewPriorityMempool(mempool.PriorityNonceMempoolConfig[int64]{
		TxPriority:      mempool.NewDefaultTxPriority(),
		SignerExtractor: signerExtractor,
		MaxTx:           maxTxs,
	})
	baseappOptions = append(baseappOptions, func(app *baseapp.BaseApp) {
		app.SetMempool(mpool)
		handler := baseapp.NewDefaultProposalHandler(mpool, app)
		handler.SetSignerExtractionAdapter(signerExtractor)
		app.SetPrepareProposal(handler.PrepareProposalHandler())
		app.SetProcessProposal(handler.ProcessProposalHandler())
	})

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	chainID, err := getChainIDFromOpts(appOpts)
	if err != nil {
		panic(err)
	}

	zetachain, err := chains.ZetaChainFromCosmosChainID(chainID)
	if err != nil {
		panic(err)
	}

	// this version is not allowed for validators
	privValStatePath := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "priv_validator_state.json")
	// #nosec G304 -- this is file present on every node
	if data, err := os.ReadFile(privValStatePath); err == nil {
		var state PrivValidatorState
		if err := json.Unmarshal(data, &state); err == nil {
			h, err := strconv.ParseInt(state.Height, 10, 64)
			if err != nil {
				panic(err)
			}
			if h > 0 {
				panic("version not allowed for validators")
			}
		}
	}

	return app.New(logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		uint64(zetachain.ChainId), //#nosec G115 chain id won't exceed uint64
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

	loadLatest := height == -1

	chainID, err := getChainIDFromOpts(appOpts)
	if err != nil {
		panic(err)
	}

	zetachain, err := chains.ZetaChainFromCosmosChainID(chainID)
	if err != nil {
		panic(err)
	}

	zetaApp = app.New(
		logger,
		db,
		traceStore,
		loadLatest,
		map[int64]bool{},
		homePath,
		uint(1),
		uint64(zetachain.ChainId), //#nosec G115 chain id won't exceed uint64
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

// getChainIDFromOpts returns the chain Id from app Opts
// It first tries to get from the chainId flag, if not available
// it will load from home
func getChainIDFromOpts(appOpts servertypes.AppOptions) (chainID string, err error) {
	// Get the chain Id from appOpts
	chainID = cast.ToString(appOpts.Get(flags.FlagChainID))
	if chainID == "" {
		// If not available load from home
		homeDir := cast.ToString(appOpts.Get(flags.FlagHome))
		chainID, err = zetacoredconfig.GetChainIDFromHome(homeDir)
		if err != nil {
			return "", err
		}
	}

	return
}
