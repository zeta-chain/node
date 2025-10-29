package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime/pprof"

	ethmetricsexp "github.com/ethereum/go-ethereum/metrics/exp"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	abciserver "github.com/cometbft/cometbft/abci/server"
	tcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	"github.com/cometbft/cometbft/rpc/client/local"
	cmttypes "github.com/cometbft/cometbft/types"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/evm/indexer"
	evmmempool "github.com/cosmos/evm/mempool"
	evmmetrics "github.com/cosmos/evm/metrics"
	srvflags "github.com/cosmos/evm/server/flags"
	servertypes "github.com/cosmos/evm/server/types"
	ethdebug "github.com/zeta-chain/node/rpc/namespaces/ethereum/debug"
	cosmosevmserverconfig "github.com/zeta-chain/node/server/config"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	pruningtypes "cosmossdk.io/store/pruning/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servergrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	servercmtlog "github.com/cosmos/cosmos-sdk/server/log"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdkmempool "github.com/cosmos/cosmos-sdk/types/mempool"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

// DBOpener is a function to open `application.db`, potentially with customized options.
type DBOpener func(opts types.AppOptions, rootDir string, backend dbm.BackendType) (dbm.DB, error)

type Application interface {
	types.Application
	AppWithPendingTxStream
	GetMempool() sdkmempool.ExtMempool
	SetClientCtx(clientCtx client.Context)
}

// AppCreator is a function that allows us to lazily initialize an application implementing with AppWithPendingTxStream.
type AppCreator func(log.Logger, dbm.DB, io.Writer, types.AppOptions) Application

// StartOptions defines options that can be customized in `StartCmd`
type StartOptions struct {
	AppCreator      types.AppCreator
	DefaultNodeHome string
	DBOpener        DBOpener
}

// NewDefaultStartOptions use the default db opener provided in tm-db.
func NewDefaultStartOptions(appCreator AppCreator, defaultNodeHome string) StartOptions {
	return StartOptions{
		AppCreator: func(l log.Logger, d dbm.DB, w io.Writer, ao types.AppOptions) types.Application {
			return appCreator(l, d, w, ao)
		},
		DefaultNodeHome: defaultNodeHome,
		DBOpener:        cosmosevmserverconfig.OpenDB,
	}
}

// StartCmd runs the service passed in, either stand-alone or in-process with
// CometBFT.
func StartCmd(opts StartOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		Long: `Run the full node application with CometBFT in or out of process. By
default, the application will run with CometBFT in process.

Pruning options can be provided via the '--pruning' flag or alternatively with '--pruning-keep-recent',
'pruning-keep-every', and 'pruning-interval' together.

For '--pruning' the options are as follows:

default: the last 100 states are kept in addition to every 500th state; pruning at 10 block intervals
nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
everything: all saved states will be deleted, storing only the current state; pruning at 10 block intervals
custom: allow pruning options to be manually specified through 'pruning-keep-recent', 'pruning-keep-every', and 'pruning-interval'

Node halting configurations exist in the form of two flags: '--halt-height' and '--halt-time'. During
the ABCI Commit phase, the node will check if the current block height is greater than or equal to
the halt-height or if the current block time is greater than or equal to the halt-time. If so, the
node will attempt to gracefully shutdown and the block will not be committed. In addition, the node
will not be able to commit subsequent blocks.

For profiling and benchmarking purposes, CPU profiling can be enabled via the '--cpu-profile' flag
which accepts a path for the resulting pprof file.
`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)

			// Bind flags to the Context's Viper so the app construction can set
			// options accordingly.
			err := serverCtx.Viper.BindPFlags(cmd.Flags())
			if err != nil {
				return err
			}

			_, err = server.GetPruningOptionsFromFlags(serverCtx.Viper)
			return err
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			withbft, _ := cmd.Flags().GetBool(srvflags.WithCometBFT)
			if !withbft {
				serverCtx.Logger.Info("starting ABCI without CometBFT")
				return wrapCPUProfile(serverCtx, func() error {
					return startStandAlone(serverCtx, clientCtx, opts)
				})
			}

			serverCtx.Logger.Info("Unlocking keyring")

			// fire unlock precess for keyring
			keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
			if keyringBackend == keyring.BackendFile {
				_, err = clientCtx.Keyring.List()
				if err != nil {
					return err
				}
			}

			serverCtx.Logger.Info("starting ABCI with CometBFT")

			// amino is needed here for backwards compatibility of REST routes
			err = wrapCPUProfile(serverCtx, func() error {
				return startInProcess(serverCtx, clientCtx, opts)
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, opts.DefaultNodeHome, "The application home directory")
	cmd.Flags().Bool(srvflags.WithCometBFT, true, "Run abci app embedded in-process with CometBFT")
	cmd.Flags().String(srvflags.Address, "tcp://0.0.0.0:26658", "Listen address")
	cmd.Flags().String(srvflags.Transport, "socket", "Transport protocol: socket, grpc")
	cmd.Flags().String(srvflags.TraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().String(server.FlagMinGasPrices, "", "Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum (e.g. 20000000000aatom)") //nolint:lll
	cmd.Flags().IntSlice(server.FlagUnsafeSkipUpgrades, []int{}, "Skip a set of upgrade heights to continue the old binary")
	cmd.Flags().Uint64(server.FlagHaltHeight, 0, "Block height at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Uint64(server.FlagHaltTime, 0, "Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Bool(server.FlagInterBlockCache, true, "Enable inter-block caching")
	cmd.Flags().String(srvflags.CPUProfile, "", "Enable CPU profiling and write to the provided file")
	cmd.Flags().Bool(server.FlagTrace, false, "Provide full stack traces for errors in ABCI Log")
	cmd.Flags().String(server.FlagPruning, pruningtypes.PruningOptionDefault, "Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().Uint64(server.FlagPruningKeepRecent, 0, "Number of recent heights to keep on disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(server.FlagPruningInterval, 0, "Height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom')") //nolint:lll
	cmd.Flags().Uint(server.FlagInvCheckPeriod, 0, "Assert registered invariants every N blocks")
	cmd.Flags().Uint64(server.FlagMinRetainBlocks, 0, "Minimum block height offset during ABCI commit to prune CometBFT blocks")
	cmd.Flags().String(srvflags.AppDBBackend, "", "The type of database for application and snapshots databases")
	cmd.Flags().Int32(server.FlagMempoolMaxTxs, 0, "The maximum number of transactions in the mempool")

	cmd.Flags().Bool(srvflags.GRPCOnly, false, "Start the node in gRPC query only mode without CometBFT process")
	cmd.Flags().Bool(srvflags.GRPCEnable, cosmosevmserverconfig.DefaultGRPCEnable, "Define if the gRPC server should be enabled")
	cmd.Flags().String(srvflags.GRPCAddress, serverconfig.DefaultGRPCAddress, "the gRPC server address to listen on")
	cmd.Flags().Bool(srvflags.GRPCWebEnable, cosmosevmserverconfig.DefaultGRPCWebEnable, "Define if the gRPC-Web server should be enabled. (Note: gRPC must also be enabled.)")
	cmd.Flags().String(srvflags.GRPCWebAddress, cosmosevmserverconfig.DefaultGRPCAddress, "The gRPC-Web server address to listen on")

	cmd.Flags().Bool(srvflags.RPCEnable, cosmosevmserverconfig.DefaultAPIEnable, "Defines if Cosmos-sdk REST server should be enabled")
	cmd.Flags().Bool(srvflags.EnabledUnsafeCors, false, "Defines if CORS should be enabled (unsafe - use it at your own risk)")

	cmd.Flags().Bool(srvflags.JSONRPCEnable, cosmosevmserverconfig.DefaultJSONRPCEnable, "Define if the JSON-RPC server should be enabled")
	cmd.Flags().StringSlice(srvflags.JSONRPCAPI, cosmosevmserverconfig.GetDefaultAPINamespaces(), "Defines a list of JSON-RPC namespaces that should be enabled")
	cmd.Flags().String(srvflags.JSONRPCAddress, cosmosevmserverconfig.DefaultJSONRPCAddress, "the JSON-RPC server address to listen on")
	cmd.Flags().String(srvflags.JSONWsAddress, cosmosevmserverconfig.DefaultJSONRPCWsAddress, "the JSON-RPC WS server address to listen on")
	cmd.Flags().StringSlice(srvflags.JSONRPCWSOrigins, cosmosevmserverconfig.GetDefaultWSOrigins(), "Defines a list of WebSocket origins that should be allowed to connect")
	cmd.Flags().Uint64(srvflags.JSONRPCGasCap, cosmosevmserverconfig.DefaultGasCap, "Sets a cap on gas that can be used in eth_call/estimateGas unit is aatom (0=infinite)")                         //nolint:lll
	cmd.Flags().Bool(srvflags.JSONRPCAllowInsecureUnlock, cosmosevmserverconfig.DefaultJSONRPCAllowInsecureUnlock, "Allow insecure account unlocking when account-related RPCs are exposed by http") //nolint:lll
	cmd.Flags().Float64(srvflags.JSONRPCTxFeeCap, cosmosevmserverconfig.DefaultTxFeeCap, "Sets a cap on transaction fee that can be sent via the RPC APIs (1 = default 1 evmos)")                    //nolint:lll
	cmd.Flags().Int32(srvflags.JSONRPCFilterCap, cosmosevmserverconfig.DefaultFilterCap, "Sets the global cap for total number of filters that can be created")
	cmd.Flags().Duration(srvflags.JSONRPCEVMTimeout, cosmosevmserverconfig.DefaultEVMTimeout, "Sets a timeout used for eth_call (0=infinite)")
	cmd.Flags().Duration(srvflags.JSONRPCHTTPTimeout, cosmosevmserverconfig.DefaultHTTPTimeout, "Sets a read/write timeout for json-rpc http server (0=infinite)")
	cmd.Flags().Duration(srvflags.JSONRPCHTTPIdleTimeout, cosmosevmserverconfig.DefaultHTTPIdleTimeout, "Sets a idle timeout for json-rpc http server (0=infinite)")
	cmd.Flags().Bool(srvflags.JSONRPCAllowUnprotectedTxs, cosmosevmserverconfig.DefaultAllowUnprotectedTxs, "Allow for unprotected (non EIP155 signed) transactions to be submitted via the node's RPC when the global parameter is disabled") //nolint:lll
	cmd.Flags().Int(srvflags.JSONRPCBatchRequestLimit, cosmosevmserverconfig.DefaultBatchRequestLimit, "Maximum number of requests in a batch")
	cmd.Flags().Int(srvflags.JSONRPCBatchResponseMaxSize, cosmosevmserverconfig.DefaultBatchResponseMaxSize, "Maximum size of server response")
	cmd.Flags().Int32(srvflags.JSONRPCLogsCap, cosmosevmserverconfig.DefaultLogsCap, "Sets the max number of results can be returned from single `eth_getLogs` query")
	cmd.Flags().Int32(srvflags.JSONRPCBlockRangeCap, cosmosevmserverconfig.DefaultBlockRangeCap, "Sets the max block range allowed for `eth_getLogs` query")
	cmd.Flags().Int(srvflags.JSONRPCMaxOpenConnections, cosmosevmserverconfig.DefaultMaxOpenConnections, "Sets the maximum number of simultaneous connections for the server listener") //nolint:lll
	cmd.Flags().Bool(srvflags.JSONRPCEnableIndexer, false, "Enable the custom tx indexer for json-rpc")
	cmd.Flags().Bool(srvflags.JSONRPCEnableMetrics, false, "Define if EVM rpc metrics server should be enabled")
	cmd.Flags().Bool(srvflags.JSONRPCEnableProfiling, false, "Enables the profiling in the debug namespace")

	cmd.Flags().String(srvflags.EVMTracer, cosmosevmserverconfig.DefaultEVMTracer, "the EVM tracer type to collect execution traces from the EVM transaction execution (json|struct|access_list|markdown)") //nolint:lll
	cmd.Flags().Uint64(srvflags.EVMMaxTxGasWanted, cosmosevmserverconfig.DefaultMaxTxGasWanted, "the gas wanted for each eth tx returned in ante handler in check tx mode")                                 //nolint:lll
	cmd.Flags().Bool(srvflags.EVMEnablePreimageRecording, cosmosevmserverconfig.DefaultEnablePreimageRecording, "Enables tracking of SHA3 preimages in the EVM (not implemented yet)")                      //nolint:lll
	cmd.Flags().Uint64(srvflags.EVMChainID, cosmosevmserverconfig.DefaultEVMChainID, "the EIP-155 compatible replay protection chain ID")
	cmd.Flags().Uint64(srvflags.EVMMinTip, cosmosevmserverconfig.DefaultEVMMinTip, "the minimum priority fee for the mempool")
	cmd.Flags().String(srvflags.EvmGethMetricsAddress, cosmosevmserverconfig.DefaultGethMetricsAddress, "the address to bind the geth metrics server to")

	cmd.Flags().Uint64(srvflags.EVMMempoolPriceLimit, cosmosevmserverconfig.DefaultMempoolConfig().PriceLimit, "the minimum gas price to enforce for acceptance into the pool (in wei)")
	cmd.Flags().Uint64(srvflags.EVMMempoolPriceBump, cosmosevmserverconfig.DefaultMempoolConfig().PriceBump, "the minimum price bump percentage to replace an already existing transaction (nonce)")
	cmd.Flags().Uint64(srvflags.EVMMempoolAccountSlots, cosmosevmserverconfig.DefaultMempoolConfig().AccountSlots, "the number of executable transaction slots guaranteed per account")
	cmd.Flags().Uint64(srvflags.EVMMempoolGlobalSlots, cosmosevmserverconfig.DefaultMempoolConfig().GlobalSlots, "the maximum number of executable transaction slots for all accounts")
	cmd.Flags().Uint64(srvflags.EVMMempoolAccountQueue, cosmosevmserverconfig.DefaultMempoolConfig().AccountQueue, "the maximum number of non-executable transaction slots permitted per account")
	cmd.Flags().Uint64(srvflags.EVMMempoolGlobalQueue, cosmosevmserverconfig.DefaultMempoolConfig().GlobalQueue, "the maximum number of non-executable transaction slots for all accounts")
	cmd.Flags().Duration(srvflags.EVMMempoolLifetime, cosmosevmserverconfig.DefaultMempoolConfig().Lifetime, "the maximum amount of time non-executable transaction are queued")

	cmd.Flags().String(srvflags.TLSCertPath, "", "the cert.pem file path for the server TLS configuration")
	cmd.Flags().String(srvflags.TLSKeyPath, "", "the key.pem file path for the server TLS configuration")

	cmd.Flags().Uint64(server.FlagStateSyncSnapshotInterval, 0, "State sync snapshot interval")
	cmd.Flags().Uint32(server.FlagStateSyncSnapshotKeepRecent, 2, "State sync snapshot to keep")

	// add support for all CometBFT-specific command line options
	tcmd.AddNodeFlags(cmd)
	return cmd
}

// startStandAlone starts an ABCI server in stand-alone mode.
// Parameters:
// - svrCtx: The context object that holds server configurations, logger, and other stateful information.
// - opts: Options for starting the server, including functions for creating the application and opening the database.
func startStandAlone(svrCtx *server.Context, clientCtx client.Context, opts StartOptions) error {
	addr := svrCtx.Viper.GetString(srvflags.Address)
	transport := svrCtx.Viper.GetString(srvflags.Transport)
	home := svrCtx.Viper.GetString(flags.FlagHome)

	db, err := opts.DBOpener(svrCtx.Viper, home, server.GetAppDBBackend(svrCtx.Viper))
	if err != nil {
		return err
	}

	var app types.Application
	defer func() {
		if app == nil {
			if err := db.Close(); err != nil {
				svrCtx.Logger.Error("error closing db", "error", err.Error())
			}
		}
	}()

	traceWriterFile := svrCtx.Viper.GetString(srvflags.TraceStore)
	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		return err
	}

	app = opts.AppCreator(svrCtx.Logger, db, traceWriter, svrCtx.Viper)
	defer func() {
		if err := app.Close(); err != nil {
			svrCtx.Logger.Error("close application failed", "error", err.Error())
		}
	}()
	evmApp, ok := app.(Application)
	if !ok {
		svrCtx.Logger.Error("failed to get server config", "error", err.Error())
	}
	evmApp.SetClientCtx(clientCtx)

	config, err := cosmosevmserverconfig.GetConfig(svrCtx.Viper)
	if err != nil {
		svrCtx.Logger.Error("failed to get server config", "error", err.Error())
		return err
	}

	if err := config.ValidateBasic(); err != nil {
		svrCtx.Logger.Error("invalid server config", "error", err.Error())
		return err
	}

	_, err = startTelemetry(config)
	if err != nil {
		return err
	}

	cmtApp := server.NewCometABCIWrapper(app)
	svr, err := abciserver.NewServer(addr, transport, cmtApp)
	if err != nil {
		return fmt.Errorf("error creating listener: %v", err)
	}

	svr.SetLogger(servercmtlog.CometLoggerWrapper{Logger: svrCtx.Logger.With("server", "abci")})
	g, ctx := getCtx(svrCtx, false)

	g.Go(func() error {
		if err := svr.Start(); err != nil {
			svrCtx.Logger.Error("failed to start out-of-process ABCI server", "err", err)
			return err
		}

		// Wait for the calling process to be canceled or close the provided context,
		// so we can gracefully stop the ABCI server.
		<-ctx.Done()
		svrCtx.Logger.Info("stopping the ABCI server...")
		return svr.Stop()
	})

	return g.Wait()
}

// legacyAminoCdc is used for the legacy REST API
func startInProcess(svrCtx *server.Context, clientCtx client.Context, opts StartOptions) (err error) {
	cfg := svrCtx.Config
	home := cfg.RootDir
	logger := svrCtx.Logger
	g, ctx := getCtx(svrCtx, true)

	if cpuProfile := svrCtx.Viper.GetString(srvflags.CPUProfile); cpuProfile != "" {
		fp, err := ethdebug.ExpandHome(cpuProfile)
		if err != nil {
			svrCtx.Logger.Debug("failed to get filepath for the CPU profile file", "error", err.Error())
			return err
		}

		f, err := os.Create(fp)
		if err != nil {
			return err
		}

		svrCtx.Logger.Info("starting CPU profiler", "profile", cpuProfile)
		if err := pprof.StartCPUProfile(f); err != nil {
			return err
		}

		defer func() {
			svrCtx.Logger.Info("stopping CPU profiler", "profile", cpuProfile)
			pprof.StopCPUProfile()
			if err := f.Close(); err != nil {
				logger.Error("failed to close CPU profiler file", "error", err.Error())
			}
		}()
	}

	db, err := opts.DBOpener(svrCtx.Viper, home, server.GetAppDBBackend(svrCtx.Viper))
	if err != nil {
		logger.Error("failed to open DB", "error", err.Error())
		return err
	}

	var app types.Application
	defer func() {
		if app == nil {
			if err := db.Close(); err != nil {
				svrCtx.Logger.Error("error closing db", "error", err.Error())
			}
		}
	}()

	traceWriterFile := svrCtx.Viper.GetString(srvflags.TraceStore)
	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		logger.Error("failed to open trace writer", "error", err.Error())
		return err
	}

	config, err := cosmosevmserverconfig.GetConfig(svrCtx.Viper)
	if err != nil {
		logger.Error("failed to get server config", "error", err.Error())
		return err
	}

	if err := config.ValidateBasic(); err != nil {
		logger.Error("invalid server config", "error", err.Error())
		return err
	}

	app = opts.AppCreator(svrCtx.Logger, db, traceWriter, svrCtx.Viper)
	defer func() {
		if err := app.Close(); err != nil {
			logger.Error("close application failed", "error", err.Error())
		}
	}()
	evmApp, ok := app.(Application)
	if !ok {
		svrCtx.Logger.Error("failed to get server config", "error", err.Error())
	}
	evmApp.SetClientCtx(clientCtx)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		logger.Error("failed load or gen node key", "error", err.Error())
		return err
	}

	genDocProvider := GenDocProvider(cfg)

	var (
		bftNode  *node.Node
		gRPCOnly = svrCtx.Viper.GetBool(srvflags.GRPCOnly)
	)

	if gRPCOnly {
		logger.Info("starting node in query only mode; CometBFT is disabled")
		config.GRPC.Enable = true
		config.JSONRPC.EnableIndexer = false
	} else {
		logger.Info("starting node with ABCI CometBFT in-process")

		cmtApp := server.NewCometABCIWrapper(app)
		bftNode, err = node.NewNode(
			cfg,
			pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
			nodeKey,
			proxy.NewLocalClientCreator(cmtApp),
			genDocProvider,
			cmtcfg.DefaultDBProvider,
			node.DefaultMetricsProvider(cfg.Instrumentation),
			servercmtlog.CometLoggerWrapper{Logger: svrCtx.Logger.With("server", "node")},
		)
		if err != nil {
			logger.Error("failed init node", "error", err.Error())
			return err
		}

		if err := bftNode.Start(); err != nil {
			logger.Error("failed start CometBFT server", "error", err.Error())
			return err
		}

		if m, ok := evmApp.GetMempool().(*evmmempool.ExperimentalEVMMempool); ok && m != nil {
			m.SetEventBus(bftNode.EventBus())
		}
		defer func() {
			if bftNode.IsRunning() {
				_ = bftNode.Stop()
			}
		}()
	}

	// Add the tx service to the gRPC router. We only need to register this
	// service if API or gRPC or JSONRPC is enabled, and avoid doing so in the general
	// case, because it spawns a new local CometBFT RPC client.
	if (config.API.Enable || config.GRPC.Enable || config.JSONRPC.Enable || config.JSONRPC.EnableIndexer) && bftNode != nil {
		clientCtx = clientCtx.WithClient(local.New(bftNode))

		app.RegisterTxService(clientCtx)
		app.RegisterTendermintService(clientCtx)
		app.RegisterNodeService(clientCtx, config.Config)
	}

	metrics, err := startTelemetry(config)
	if err != nil {
		return err
	}

	// Enable metrics if JSONRPC is enabled and --metrics is passed
	// Flag not added in config to avoid user enabling in config without passing in CLI
	if config.JSONRPC.Enable && svrCtx.Viper.GetBool(srvflags.JSONRPCEnableMetrics) {
		ethmetricsexp.Setup(config.JSONRPC.MetricsAddress)
	}

	var idxer servertypes.EVMTxIndexer
	if config.JSONRPC.EnableIndexer {
		idxDB, err := OpenIndexerDB(home, server.GetAppDBBackend(svrCtx.Viper))
		if err != nil {
			logger.Error("failed to open evm indexer DB", "error", err.Error())
			return err
		}

		idxLogger := svrCtx.Logger.With("indexer", "evm")
		idxer = indexer.NewKVIndexer(idxDB, idxLogger, clientCtx)
		indexerService := NewEVMIndexerService(idxer, clientCtx.Client.(rpcclient.Client))
		indexerService.SetLogger(servercmtlog.CometLoggerWrapper{Logger: idxLogger})

		g.Go(func() error {
			errCh := make(chan error, 1)
			go func() {
				if err := indexerService.Start(); err != nil {
					errCh <- err
				}
			}()

			select {
			case <-ctx.Done():
				logger.Info("stopping evm indexer service due to context cancellation")
				if err := indexerService.Stop(); err != nil {
					logger.Error("failed to stop evm indexer service", "error", err.Error())
				}
				return ctx.Err()
			case err := <-errCh:
				if err != nil {
					logger.Error("evm indexer service failed", "error", err.Error())
				}
				return err
			}
		})
	}

	if config.API.Enable || config.JSONRPC.Enable {
		genDoc, err := genDocProvider()
		if err != nil {
			return err
		}

		clientCtx = clientCtx.
			WithHomeDir(home).
			WithChainID(genDoc.ChainID)
	}

	grpcSrv, clientCtx, err := startGrpcServer(ctx, svrCtx, clientCtx, g, config.GRPC, app)
	if err != nil {
		return err
	}
	if grpcSrv != nil {
		defer grpcSrv.GracefulStop()
	}

	startAPIServer(ctx, svrCtx, clientCtx, g, config.Config, app, grpcSrv, metrics, config.EVM.GethMetricsAddress)

	if config.JSONRPC.Enable {
		txApp, ok := app.(AppWithPendingTxStream)
		if !ok {
			return fmt.Errorf("json-rpc server requires AppWithPendingTxStream")
		}
		_, err = StartJSONRPC(ctx, svrCtx, clientCtx, g, &config, idxer, txApp, evmApp.GetMempool().(*evmmempool.ExperimentalEVMMempool))
		if err != nil {
			return err
		}
	}

	// At this point it is safe to block the process if we're in query only mode as
	// we do not need to start Rosetta or handle any CometBFT related processes.
	if gRPCOnly {
		// wait for signal capture and gracefully return
		// we are guaranteed to be waiting for the "ListenForQuitSignals" goroutine.
		return g.Wait()
	}

	// wait for signal capture and gracefully return
	// we are guaranteed to be waiting for the "ListenForQuitSignals" goroutine.
	return g.Wait()
}

// OpenIndexerDB opens the custom eth indexer db, using the same db backend as the main app
func OpenIndexerDB(rootDir string, backendType dbm.BackendType) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("evmindexer", backendType, dataDir)
}

// openTraceWriter opens a trace writer if a trace store file is specified.
// Parameters:
// - traceWriterFile: The path to the trace store file. If this is an empty string, no file will be opened.
func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile == "" {
		return
	}

	filePath := filepath.Clean(traceWriterFile)
	return os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0o600,
	)
}

func startTelemetry(cfg cosmosevmserverconfig.Config) (*telemetry.Metrics, error) {
	if !cfg.Telemetry.Enabled {
		return nil, nil
	}
	return telemetry.New(cfg.Telemetry)
}

// wrapCPUProfile runs callback in a goroutine, then wait for quit signals.
func wrapCPUProfile(ctx *server.Context, callback func() error) error {
	if cpuProfile := ctx.Viper.GetString(srvflags.CPUProfile); cpuProfile != "" {
		fp, err := ethdebug.ExpandHome(cpuProfile)
		if err != nil {
			ctx.Logger.Debug("failed to get filepath for the CPU profile file", "error", err.Error())
			return err
		}
		f, err := os.Create(fp)
		if err != nil {
			return err
		}

		ctx.Logger.Info("starting CPU profiler", "profile", cpuProfile)
		if err := pprof.StartCPUProfile(f); err != nil {
			return err
		}

		defer func() {
			ctx.Logger.Info("stopping CPU profiler", "profile", cpuProfile)
			pprof.StopCPUProfile()
			if err := f.Close(); err != nil {
				ctx.Logger.Info("failed to close cpu-profile file", "profile", cpuProfile, "err", err.Error())
			}
		}()
	}

	return callback()
}

func getCtx(svrCtx *server.Context, block bool) (*errgroup.Group, context.Context) {
	ctx, cancelFn := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	// listen for quit signals so the calling parent process can gracefully exit
	server.ListenForQuitSignals(g, block, cancelFn, svrCtx.Logger)
	return g, ctx
}

// startGrpcServer starts a gRPC server based on the provided configuration.
func startGrpcServer(
	ctx context.Context,
	svrCtx *server.Context,
	clientCtx client.Context,
	g *errgroup.Group,
	config serverconfig.GRPCConfig,
	app types.Application,
) (*grpc.Server, client.Context, error) {
	if !config.Enable {
		// return grpcServer as nil if gRPC is disabled
		return nil, clientCtx, nil
	}
	_, _, err := net.SplitHostPort(config.Address)
	if err != nil {
		return nil, clientCtx, errorsmod.Wrapf(err, "invalid grpc address %s", config.Address)
	}

	maxSendMsgSize := config.MaxSendMsgSize
	if maxSendMsgSize == 0 {
		maxSendMsgSize = serverconfig.DefaultGRPCMaxSendMsgSize
	}

	maxRecvMsgSize := config.MaxRecvMsgSize
	if maxRecvMsgSize == 0 {
		maxRecvMsgSize = serverconfig.DefaultGRPCMaxRecvMsgSize
	}

	// if gRPC is enabled, configure gRPC client for gRPC gateway and json-rpc
	grpcClient, err := grpc.NewClient(
		config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.ForceCodec(codec.NewProtoCodec(clientCtx.InterfaceRegistry).GRPCCodec()),
			grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
			grpc.MaxCallSendMsgSize(maxSendMsgSize),
		),
	)
	if err != nil {
		return nil, clientCtx, err
	}
	// Set `GRPCClient` to `clientCtx` to enjoy concurrent grpc query.
	// only use it if gRPC server is enabled.
	clientCtx = clientCtx.WithGRPCClient(grpcClient)
	svrCtx.Logger.Debug("gRPC client assigned to client context", "address", config.Address)

	grpcSrv, err := servergrpc.NewGRPCServer(clientCtx, app, config)
	if err != nil {
		return nil, clientCtx, err
	}

	// Start the gRPC server in a goroutine. Note, the provided ctx will ensure
	// that the server is gracefully shut down.
	g.Go(func() error {
		return servergrpc.StartGRPCServer(ctx, svrCtx.Logger.With("module", "grpc-server"), config, grpcSrv)
	})
	return grpcSrv, clientCtx, nil
}

// startAPIServer starts an API server based on the provided configuration and application context.
// Parameters:
// - ctx: The context used for managing the server's lifecycle, allowing for graceful shutdown.
// - svrCtx: The server context containing configuration, logger, and other stateful components.
// - clientCtx: The client context, which provides necessary information for API operations.
// - g: An errgroup.Group for managing goroutines and handling errors concurrently.
// - svrCfg: The server configuration that specifies whether the API server is enabled and other settings.
// - app: The application instance that registers API routes.
// - grpcSrv: A pointer to the gRPC server, which may be used by the API server.
// - metrics: A telemetry metrics instance for monitoring API server performance.
func startAPIServer(
	ctx context.Context,
	svrCtx *server.Context,
	clientCtx client.Context,
	g *errgroup.Group,
	svrCfg serverconfig.Config,
	app types.Application,
	grpcSrv *grpc.Server,
	metrics *telemetry.Metrics,
	gethMetricsAddress string,
) {
	if !svrCfg.API.Enable {
		return
	}

	apiSrv := api.New(clientCtx, svrCtx.Logger.With("server", "api"), grpcSrv)
	app.RegisterAPIRoutes(apiSrv, svrCfg.API)

	if svrCfg.Telemetry.Enabled {
		apiSrv.SetTelemetry(metrics)
		g.Go(func() error {
			return evmmetrics.StartGethMetricServer(ctx, svrCtx.Logger.With("server", "geth_metrics"), gethMetricsAddress)
		})
	}

	g.Go(func() error {
		return apiSrv.Start(ctx, svrCfg)
	})
}

// GenDocProvider returns a function which returns the genesis doc from the genesis file.
func GenDocProvider(cfg *cmtcfg.Config) func() (*cmttypes.GenesisDoc, error) {
	return func() (*cmttypes.GenesisDoc, error) {
		appGenesis, err := genutiltypes.AppGenesisFromFile(cfg.GenesisFile())
		if err != nil {
			return nil, err
		}

		return appGenesis.ToGenesisDoc()
	}
}
