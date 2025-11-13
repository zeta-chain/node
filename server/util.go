package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	abciserver "github.com/cometbft/cometbft/abci/server"
	tmcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	cmtcfg "github.com/cometbft/cometbft/config"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	cmtstate "github.com/cometbft/cometbft/proto/tendermint/state"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/proxy"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	"github.com/cometbft/cometbft/rpc/client/local"
	jsonrpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	sm "github.com/cometbft/cometbft/state"
	"github.com/cometbft/cometbft/store"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servergrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	servercmtlog "github.com/cosmos/cosmos-sdk/server/log"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/version"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/evm/indexer"
	"github.com/cosmos/evm/server/config"
	srvflags "github.com/cosmos/evm/server/flags"
	cosmosevmtypes "github.com/cosmos/evm/types"
	ethmetricsexp "github.com/ethereum/go-ethereum/metrics/exp"
	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/spf13/cobra"
	"golang.org/x/net/netutil"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ethdebug "github.com/zeta-chain/node/rpc/namespaces/ethereum/debug"
)

// AddCommands adds server commands
func AddCommands(
	rootCmd *cobra.Command,
	appCreator types.AppCreator,
	defaultNodeHome string,
	appExport types.AppExporter,
	addStartFlags types.ModuleInitFlags,
) {
	cometbftCmd := &cobra.Command{
		Use:     "comet",
		Aliases: []string{"cometbft", "tendermint"},
		Short:   "CometBFT subcommands",
	}

	cometbftCmd.AddCommand(
		server.ShowNodeIDCmd(),
		server.ShowValidatorCmd(),
		server.ShowAddressCmd(),
		server.VersionCmd(),
		tmcmd.ResetAllCmd,
		tmcmd.ResetStateCmd,
		server.BootstrapStateCmd(appCreator),
	)

	startCmd := StartCmd(appCreator, defaultNodeHome)
	addStartFlags(startCmd)

	rootCmd.AddCommand(
		startCmd,
		cometbftCmd,
		server.ExportCmd(appExport, defaultNodeHome),
		version.NewVersionCommand(),
		server.NewRollbackCmd(appCreator, defaultNodeHome),

		// custom tx indexer command
		NewIndexTxCmd(),
	)
}

// ConnectTmWS connects to a Tendermint WebSocket (WS) server.
// Parameters:
// - tmRPCAddr: The RPC address of the Tendermint server.
// - tmEndpoint: The WebSocket endpoint on the Tendermint server.
// - logger: A logger instance used to log debug and error messages.
func ConnectTmWS(tmRPCAddr, tmEndpoint string, logger log.Logger) *jsonrpcclient.WSClient {
	tmWsClient, err := jsonrpcclient.NewWS(tmRPCAddr, tmEndpoint,
		jsonrpcclient.MaxReconnectAttempts(256),
		jsonrpcclient.ReadWait(120*time.Second),
		jsonrpcclient.WriteWait(120*time.Second),
		jsonrpcclient.PingPeriod(50*time.Second),
		jsonrpcclient.OnReconnect(func() {
			logger.Debug("EVM RPC reconnects to Tendermint WS", "address", tmRPCAddr+tmEndpoint)
		}),
	)

	if err != nil {
		logger.Error(
			"Tendermint WS client could not be created",
			"address", tmRPCAddr+tmEndpoint,
			"error", err,
		)
	} else if err := tmWsClient.OnStart(); err != nil {
		logger.Error(
			"Tendermint WS client could not start",
			"address", tmRPCAddr+tmEndpoint,
			"error", err,
		)
	}

	return tmWsClient
}

// MountGRPCWebServices mounts gRPC-Web services on specific HTTP POST routes.
// Parameters:
// - router: The HTTP router instance to mount the routes on (using mux.Router).
// - grpcWeb: The wrapped gRPC-Web server that will handle incoming gRPC-Web and WebSocket requests.
// - grpcResources: A list of resource endpoints (URLs) that should be mounted for gRPC-Web POST requests.
// - logger: A logger instance used to log information about the mounted resources.
func MountGRPCWebServices(
	router *mux.Router,
	grpcWeb *grpcweb.WrappedGrpcServer,
	grpcResources []string,
	logger log.Logger,
) {
	for _, res := range grpcResources {
		logger.Info("[GRPC Web] HTTP POST mounted", "resource", res)

		s := router.Methods("POST").Subrouter()
		s.HandleFunc(res, func(resp http.ResponseWriter, req *http.Request) {
			if grpcWeb.IsGrpcWebSocketRequest(req) {
				grpcWeb.HandleGrpcWebsocketRequest(resp, req)
				return
			}

			if grpcWeb.IsGrpcWebRequest(req) {
				grpcWeb.HandleGrpcWebRequest(resp, req)
				return
			}
		})
	}
}

// Listen starts a net.Listener on the tcp network on the given address.
// If there is a specified MaxOpenConnections in the config, it will also set the limitListener.
func Listen(addr string, config *config.Config) (net.Listener, error) {
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	if config.JSONRPC.MaxOpenConnections > 0 {
		ln = netutil.LimitListener(ln, config.JSONRPC.MaxOpenConnections)
	}
	return ln, err
}

func start(
	svrCtx *server.Context,
	clientCtx client.Context,
	appCreator types.AppCreator,
	withCmt bool,
	opts StartCmdOptions,
) error {
	serverConfig, err := config.GetConfig(svrCtx.Viper)
	if err != nil {
		return fmt.Errorf("failed to get server config: %w", err)
	}

	if err := serverConfig.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid server config: %w", err)
	}

	app, appCleanupFn, err := setupApp(svrCtx, appCreator, opts)
	if err != nil {
		return err
	}
	defer appCleanupFn()

	metrics, err := startTelemetry(serverConfig)
	if err != nil {
		return err
	}

	if !withCmt {
		svrCtx.Logger.Info("starting ABCI without CometBFT")
		return startStandAlone(svrCtx, clientCtx, app, opts)
	}

	svrCtx.Logger.Info("starting ABCI with CometBFT")
	return startInProcess(svrCtx, &serverConfig, clientCtx, app, metrics, opts)
}

func setupApp(
	svrCtx *server.Context,
	appCreator types.AppCreator,
	opts StartCmdOptions,
) (types.Application, func(), error) {
	traceWriter, traceCleanupFn, err := setupTraceWriter(svrCtx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup trace writer: %w", err)
	}

	home := svrCtx.Config.RootDir
	db, err := opts.DBOpener(home, server.GetAppDBBackend(svrCtx.Viper))
	if err != nil {
		traceCleanupFn()
		return nil, nil, fmt.Errorf("failed to open DB: %w", err)
	}

	var app types.Application
	if isDevnet, ok := svrCtx.Viper.Get(KeyIsDevnet).(bool); ok && isDevnet {
		svrCtx.Logger.Info("starting in devnet mode, applying devnetify")
		app, err = devnetify(svrCtx, appCreator, db, traceWriter)
		if err != nil {
			traceCleanupFn()
			if closeErr := db.Close(); closeErr != nil {
				svrCtx.Logger.Error("error closing db after devnetify failure", "error", closeErr)
			}
			return nil, nil, fmt.Errorf("failed to devnetify: %w", err)
		}
		err := initAppForDevnet(svrCtx, app)
		if err != nil {
			traceCleanupFn()
			if closeErr := db.Close(); closeErr != nil {
				svrCtx.Logger.Error("error closing db after devnetify failure", "error", closeErr)
			}
			return nil, nil, fmt.Errorf("failed to init app for devnet: %w", err)
		}
	} else {
		app = appCreator(svrCtx.Logger, db, traceWriter, svrCtx.Viper)
	}

	cleanupFn := func() {
		traceCleanupFn()
		if localErr := app.Close(); localErr != nil {
			svrCtx.Logger.Error("error closing app", "error", localErr)
		}
		// Note: db.Close() is not called here because app.Close() already closes the database.
		// Calling db.Close() again would cause "pebble: closed" panic.
	}

	return app, cleanupFn, nil
}

func setupTraceWriter(svrCtx *server.Context) (traceWriter io.WriteCloser, cleanup func(), err error) {
	// clean up the traceWriter when the server is shutting down
	cleanup = func() {}

	traceWriterFile := svrCtx.Viper.GetString(srvflags.TraceStore)
	traceWriter, err = openTraceWriter(traceWriterFile)
	if err != nil {
		return traceWriter, cleanup, err
	}

	// if flagTraceStore is not used then traceWriter is nil
	if traceWriter != nil {
		cleanup = func() {
			if err = traceWriter.Close(); err != nil {
				svrCtx.Logger.Error("failed to close trace writer", "err", err)
			}
		}
	}

	return traceWriter, cleanup, nil
}

// startStandAlone starts an ABCI server in stand-alone mode
func startStandAlone(
	svrCtx *server.Context,
	clientCtx client.Context,
	app types.Application,
	opts StartCmdOptions,
) error {
	addr := svrCtx.Viper.GetString(srvflags.Address)
	transport := svrCtx.Viper.GetString(srvflags.Transport)

	cmtApp := server.NewCometABCIWrapper(app)
	svr, err := abciserver.NewServer(addr, transport, cmtApp)
	if err != nil {
		return fmt.Errorf("error creating listener: %v", err)
	}

	svr.SetLogger(servercmtlog.CometLoggerWrapper{Logger: svrCtx.Logger.With("server", "abci")})
	g, ctx := getCtx(svrCtx, false)

	if opts.PostSetupStandalone != nil {
		if err := opts.PostSetupStandalone(svrCtx, clientCtx, ctx, g); err != nil {
			return err
		}
	}

	g.Go(func() error {
		if err := svr.Start(); err != nil {
			svrCtx.Logger.Error("failed to start out-of-process ABCI server", "err", err)
			return err
		}

		<-ctx.Done()
		svrCtx.Logger.Info("stopping the ABCI server...")
		return svr.Stop()
	})

	return g.Wait()
}

func startInProcess(
	svrCtx *server.Context,
	config *config.Config,
	clientCtx client.Context,
	app types.Application,
	metrics *telemetry.Metrics,
	opts StartCmdOptions,
) error {
	cfg := svrCtx.Config
	logger := svrCtx.Logger
	g, ctx := getCtx(svrCtx, true)

	if err := setupCPUProfiling(svrCtx); err != nil {
		return err
	}
	defer stopCPUProfiling(svrCtx)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return fmt.Errorf("failed to load or gen node key: %w", err)
	}

	genDocProvider := GenDocProvider(cfg)

	var cmtNode *node.Node
	gRPCOnly := svrCtx.Viper.GetBool(srvflags.GRPCOnly)

	if gRPCOnly {
		logger.Info("starting node in query only mode; CometBFT is disabled")
		config.GRPC.Enable = true
		config.JSONRPC.EnableIndexer = false
	} else {
		logger.Info("starting node with ABCI CometBFT in-process")

		cmtApp := server.NewCometABCIWrapper(app)
		cmtNode, err = node.NewNode(
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
			return fmt.Errorf("failed init node: %w", err)
		}

		if err := cmtNode.Start(); err != nil {
			return fmt.Errorf("failed start CometBFT server: %w", err)
		}

		defer func() {
			if cmtNode.IsRunning() {
				_ = cmtNode.Stop()
			}
		}()
	}

	// Add the tx service to the gRPC router. We only need to register this
	// service if API or gRPC or JSONRPC is enabled, and avoid doing so in the general
	// case, because it spawns a new local CometBFT RPC client.
	if (config.API.Enable || config.GRPC.Enable || config.JSONRPC.Enable || config.JSONRPC.EnableIndexer) &&
		cmtNode != nil {
		clientCtx = clientCtx.WithClient(local.New(cmtNode))
		app.RegisterTxService(clientCtx)
		app.RegisterTendermintService(clientCtx)
		app.RegisterNodeService(clientCtx, config.Config)
	}

	// Enable metrics if JSONRPC is enabled and --metrics is passed
	// Flag not added in config to avoid user enabling in config without passing in CLI
	if config.JSONRPC.Enable && svrCtx.Viper.GetBool(srvflags.JSONRPCEnableMetrics) {
		ethmetricsexp.Setup(config.JSONRPC.MetricsAddress)
	}

	var evmIndexer cosmosevmtypes.EVMTxIndexer
	if config.JSONRPC.EnableIndexer {
		indexDB, err := OpenIndexerDB(cfg.RootDir, server.GetAppDBBackend(svrCtx.Viper))
		if err != nil {
			return fmt.Errorf("failed to open evm indexer DB: %w", err)
		}

		idxLogger := svrCtx.Logger.With("indexer", "evm")
		evmIndexer = indexer.NewKVIndexer(indexDB, idxLogger, clientCtx)
		indexerService := NewEVMIndexerService(evmIndexer, clientCtx.Client.(rpcclient.Client))
		indexerService.SetLogger(servercmtlog.CometLoggerWrapper{Logger: idxLogger})

		g.Go(func() error {
			return indexerService.Start()
		})
	}

	if config.API.Enable || config.JSONRPC.Enable {
		genDoc, err := genDocProvider()
		if err != nil {
			return err
		}
		clientCtx = clientCtx.
			WithHomeDir(cfg.RootDir).
			WithChainID(genDoc.ChainID)
	}

	grpcSrv, clientCtx, err := startGrpcServer(ctx, svrCtx, clientCtx, g, config.GRPC, app)
	if err != nil {
		return err
	}
	if grpcSrv != nil {
		defer grpcSrv.GracefulStop()
	}

	apiSrv := startAPIServer(ctx, svrCtx, clientCtx, g, config.Config, app, grpcSrv, metrics)
	if apiSrv != nil {
		defer apiSrv.Close()
	}

	clientCtx, httpSrv, httpSrvDone, err := startJSONRPCServer(
		svrCtx,
		clientCtx,
		config,
		genDocProvider,
		cfg.RPC.ListenAddress,
		evmIndexer,
	)
	if err != nil {
		return err
	}
	if httpSrv != nil {
		defer func() {
			shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFn()
			if err := httpSrv.Shutdown(shutdownCtx); err != nil {
				logger.Error("HTTP server shutdown produced a warning", "error", err.Error())
			} else {
				logger.Info("HTTP server shut down, waiting 5 sec")
				select {
				case <-time.Tick(5 * time.Second):
				case <-httpSrvDone:
				}
			}
		}()
	}

	if opts.PostSetup != nil {
		if err := opts.PostSetup(svrCtx, clientCtx, ctx, g); err != nil {
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

func setupCPUProfiling(svrCtx *server.Context) error {
	cpuProfile := svrCtx.Viper.GetString(srvflags.CPUProfile)
	if cpuProfile == "" {
		return nil
	}

	fp, err := ethdebug.ExpandHome(cpuProfile)
	if err != nil {
		return fmt.Errorf("failed to expand cpu profile path: %w", err)
	}

	f, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("failed to create cpu profile file: %w", err)
	}

	svrCtx.Logger.Info("starting CPU profiler", "profile", cpuProfile)
	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return fmt.Errorf("failed to start cpu profile: %w", err)
	}

	svrCtx.Viper.Set("_cpu_profile_file", f)
	return nil
}

func stopCPUProfiling(svrCtx *server.Context) {
	if f := svrCtx.Viper.Get("_cpu_profile_file"); f != nil {
		svrCtx.Logger.Info("stopping CPU profiler")
		pprof.StopCPUProfile()
		if file, ok := f.(*os.File); ok {
			if err := file.Close(); err != nil {
				svrCtx.Logger.Error("failed to close CPU profiler file", "error", err.Error())
			}
		}
	}
}

func OpenIndexerDB(rootDir string, backendType dbm.BackendType) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("evmindexer", backendType, dataDir)
}

func openTraceWriter(traceWriterFile string) (w io.WriteCloser, err error) {
	if traceWriterFile == "" {
		return
	}
	return os.OpenFile(
		traceWriterFile,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0o666,
	)
}

func startTelemetry(cfg config.Config) (*telemetry.Metrics, error) {
	if !cfg.Telemetry.Enabled {
		return nil, nil
	}
	return telemetry.New(cfg.Telemetry)
}

func getCtx(svrCtx *server.Context, block bool) (*errgroup.Group, context.Context) {
	ctx, cancelFn := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	// Enable graceful shutdown if grace duration is set
	graceDuration := svrCtx.Viper.GetDuration(server.FlagShutdownGrace)
	wrappedCancelFn := cancelFn
	if graceDuration > 0 {
		wrappedCancelFn = func() {
			cancelFn()
			svrCtx.Logger.Info(
				"graceful shutdown start, waiting for services to stop",
				server.FlagShutdownGrace,
				graceDuration,
			)
			time.Sleep(graceDuration)
			svrCtx.Logger.Info("graceful shutdown complete")
		}
	}

	server.ListenForQuitSignals(g, block, wrappedCancelFn, svrCtx.Logger)
	return g, ctx
}

func startGrpcServer(
	ctx context.Context,
	svrCtx *server.Context,
	clientCtx client.Context,
	g *errgroup.Group,
	config serverconfig.GRPCConfig,
	app types.Application,
) (*grpc.Server, client.Context, error) {
	if !config.Enable {
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

	clientCtx = clientCtx.WithGRPCClient(grpcClient)
	svrCtx.Logger.Debug("gRPC client assigned to client context", "address", config.Address)

	grpcSrv, err := servergrpc.NewGRPCServer(clientCtx, app, config)
	if err != nil {
		return nil, clientCtx, err
	}

	g.Go(func() error {
		return servergrpc.StartGRPCServer(ctx, svrCtx.Logger.With("module", "grpc-server"), config, grpcSrv)
	})
	return grpcSrv, clientCtx, nil
}

func startAPIServer(
	ctx context.Context,
	svrCtx *server.Context,
	clientCtx client.Context,
	g *errgroup.Group,
	svrCfg serverconfig.Config,
	app types.Application,
	grpcSrv *grpc.Server,
	metrics *telemetry.Metrics,
) *api.Server {
	if !svrCfg.API.Enable {
		return nil
	}

	apiSrv := api.New(clientCtx, svrCtx.Logger.With("server", "api"), grpcSrv)
	app.RegisterAPIRoutes(apiSrv, svrCfg.API)

	if svrCfg.Telemetry.Enabled {
		apiSrv.SetTelemetry(metrics)
	}

	g.Go(func() error {
		return apiSrv.Start(ctx, svrCfg)
	})
	return apiSrv
}

func startJSONRPCServer(
	svrCtx *server.Context,
	clientCtx client.Context,
	config *config.Config,
	genDocProvider node.GenesisDocProvider,
	cmtRPCAddr string,
	idxer cosmosevmtypes.EVMTxIndexer,
) (ctx client.Context, httpSrv *http.Server, httpSrvDone chan struct{}, err error) {
	ctx = clientCtx
	if !config.JSONRPC.Enable {
		return
	}

	genDoc, err := genDocProvider()
	if err != nil {
		return ctx, httpSrv, httpSrvDone, err
	}

	ctx = clientCtx.WithChainID(genDoc.ChainID)
	cmtEndpoint := "/websocket"

	httpSrv, httpSrvDone, err = StartJSONRPC(svrCtx, ctx, cmtRPCAddr, cmtEndpoint, config, idxer)
	if err != nil {
		return ctx, nil, nil, err
	}

	return ctx, httpSrv, httpSrvDone, nil
}

func openDB(rootDir string, backendType dbm.BackendType) (dbm.DB, error) {
	return config.OpenDB(server.NewDefaultContext().Viper, rootDir, backendType)
}

// devnetify modifies both state and blockStore, allowing the provided operator address and local validator key to control the network
// that the state in the data folder represents. The chainID of the local genesis file is modified to match the provided chainID.
// this function focuses on modifying the CometBFT state to match the new chainID and validator info.
func devnetify(
	ctx *server.Context,
	devnetAppCreator types.AppCreator,
	db dbm.DB,
	traceWriter io.WriteCloser,
) (types.Application, error) {
	nodeConfig := ctx.Config

	// Modify P2P config to prevent connections to other peers.
	nodeConfig.P2P.Seeds = ""
	nodeConfig.P2P.PersistentPeers = ""
	nodeConfig.P2P.MaxNumInboundPeers = 0
	nodeConfig.P2P.MaxNumOutboundPeers = 0
	cmtcfg.WriteConfigFile(filepath.Join(nodeConfig.RootDir, "config", "config.toml"), nodeConfig)

	newChainID, ok := ctx.Viper.Get(KeyNewChainID).(string)
	if !ok {
		return nil, fmt.Errorf("expected string for key %s", KeyNewChainID)
	}

	// Modify app genesis chain ID and save to genesis file.
	genFilePath := nodeConfig.GenesisFile()
	appGen, err := genutiltypes.AppGenesisFromFile(genFilePath)
	if err != nil {
		return nil, err
	}
	appGen.ChainID = newChainID
	if err := appGen.ValidateAndComplete(); err != nil {
		return nil, err
	}
	if err := appGen.SaveAs(genFilePath); err != nil {
		return nil, err
	}

	// Regenerate addrbook.json to prevent peers on old network from causing error logs.
	addrBookPath := filepath.Join(nodeConfig.RootDir, "config", "addrbook.json")
	if err := os.Remove(addrBookPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to remove existing addrbook.json: %w", err)
	}

	emptyAddrBook := []byte("{}")
	if err := os.WriteFile(addrBookPath, emptyAddrBook, 0o600); err != nil {
		return nil, fmt.Errorf("failed to create empty addrbook.json: %w", err)
	}

	// Load the comet genesis doc provider.
	genDocProvider := node.DefaultGenesisDocProviderFunc(nodeConfig)

	// Initialize blockStore and stateDB.
	blockStoreDB, err := cmtcfg.DefaultDBProvider(&cmtcfg.DBContext{ID: "blockstore", Config: nodeConfig})
	if err != nil {
		return nil, err
	}
	blockStore := store.NewBlockStore(blockStoreDB)

	stateDB, err := cmtcfg.DefaultDBProvider(&cmtcfg.DBContext{ID: "state", Config: nodeConfig})
	if err != nil {
		return nil, err
	}

	defer blockStore.Close()
	defer stateDB.Close()

	privValidator := pvm.LoadOrGenFilePV(nodeConfig.PrivValidatorKeyFile(), nodeConfig.PrivValidatorStateFile())
	userPubKey, err := privValidator.GetPubKey()
	if err != nil {
		return nil, err
	}
	validatorAddress := userPubKey.Address()

	stateStore := sm.NewStore(stateDB, sm.StoreOptions{
		DiscardABCIResponses: nodeConfig.Storage.DiscardABCIResponses,
	})

	cmtState, genDoc, err := node.LoadStateFromDBOrGenesisDocProvider(stateDB, genDocProvider)
	if err != nil {
		return nil, err
	}

	// This is used later when modifying the application state.
	ctx.Viper.Set(KeyValidatorConsensusAddr, validatorAddress.Bytes())
	ctx.Viper.Set(KeyValidatorConsensusPubkey, userPubKey.Bytes())
	devnetApp := devnetAppCreator(ctx.Logger, db, traceWriter, ctx.Viper)

	// We need to create a temporary proxyApp to get the initial state of the application.
	// Depending on how the node was stopped, the application height can differ from the blockStore height.
	// This height difference changes how we go about modifying the state.
	cmtApp := server.NewCometABCIWrapper(devnetApp)
	_, proxyAppContext := getCtx(ctx, true)

	clientCreator := proxy.NewLocalClientCreator(cmtApp)
	metrics := node.DefaultMetricsProvider(cmtcfg.DefaultConfig().Instrumentation)
	//nolint:dogsled // external function, we only need proxyMetrics
	_, _, _, _, proxyMetrics, _, _ := metrics(genDoc.ChainID)
	proxyApp := proxy.NewAppConns(clientCreator, proxyMetrics)
	if err := proxyApp.Start(); err != nil {
		return nil, fmt.Errorf("error starting proxy app connections: %w", err)
	}
	res, err := proxyApp.Query().Info(proxyAppContext, proxy.RequestInfo)
	if err != nil {
		return nil, fmt.Errorf("error calling Info: %w", err)
	}
	err = proxyApp.Stop()
	if err != nil {
		return nil, fmt.Errorf("failed to stop proxy app: %w", err)
	}
	appHash := res.LastBlockAppHash
	appHeight := res.LastBlockHeight

	var block *cmttypes.Block
	switch {
	case appHeight == blockStore.Height():
		block = blockStore.LoadBlock(blockStore.Height())
		if cmtState.LastBlockHeight != appHeight {
			cmtState.LastBlockHeight = appHeight
			block.AppHash = appHash
			cmtState.AppHash = appHash
		} else {
			// Node was likely stopped via SIGTERM, delete the next block's seen commit
			err := blockStoreDB.Delete(fmt.Appendf(nil, "SC:%v", blockStore.Height()+1))
			if err != nil {
				return nil, err
			}
		}
	case blockStore.Height() > cmtState.LastBlockHeight:
		// This state usually occurs when we gracefully stop the node.
		err = blockStore.DeleteLatestBlock()
		if err != nil {
			return nil, err
		}
		block = blockStore.LoadBlock(blockStore.Height())
	default:
		// If there is any other state, we just load the block
		block = blockStore.LoadBlock(blockStore.Height())
	}

	block.ChainID = newChainID
	cmtState.ChainID = newChainID

	block.LastBlockID = cmtState.LastBlockID
	block.LastCommit.BlockID = cmtState.LastBlockID

	// Create a vote from our validator
	vote := cmttypes.Vote{
		Type:             cmtproto.PrecommitType,
		Height:           cmtState.LastBlockHeight,
		Round:            0,
		BlockID:          cmtState.LastBlockID,
		Timestamp:        time.Now(),
		ValidatorAddress: validatorAddress,
		ValidatorIndex:   0,
		Signature:        []byte{},
	}

	ctx.Viper.Set(KeyAppBlockHeight, cmtState.LastBlockHeight)

	// Sign the vote, and copy the proto changes from the act of signing to the vote itself
	voteProto := vote.ToProto()
	privValidator.Reset()

	err = privValidator.SignVote(newChainID, voteProto)
	if err != nil {
		return nil, err
	}

	vote.Signature = voteProto.Signature
	vote.Timestamp = voteProto.Timestamp

	// Modify the block's lastCommit to be signed only by our validator
	block.LastCommit.Signatures[0].ValidatorAddress = validatorAddress
	block.LastCommit.Signatures[0].Signature = vote.Signature
	block.LastCommit.Signatures = []cmttypes.CommitSig{block.LastCommit.Signatures[0]}

	// Load the seenCommit of the lastBlockHeight and modify it to be signed from our validator
	seenCommit := blockStore.LoadSeenCommit(cmtState.LastBlockHeight)
	seenCommit.BlockID = cmtState.LastBlockID
	seenCommit.Round = vote.Round
	seenCommit.Signatures[0].Signature = vote.Signature
	seenCommit.Signatures[0].ValidatorAddress = validatorAddress
	seenCommit.Signatures[0].Timestamp = vote.Timestamp
	seenCommit.Signatures = []cmttypes.CommitSig{seenCommit.Signatures[0]}
	err = blockStore.SaveSeenCommit(cmtState.LastBlockHeight, seenCommit)
	if err != nil {
		return nil, err
	}

	// Create ValidatorSet struct containing just our validator. This is updated into the comet bft state.
	// The voting power is set to a high value to avoid any potential proposer selection issues.
	newVal := &cmttypes.Validator{
		Address:     validatorAddress,
		PubKey:      userPubKey,
		VotingPower: 300000000000000,
	}
	newValSet := &cmttypes.ValidatorSet{
		Validators: []*cmttypes.Validator{newVal},
		Proposer:   newVal,
	}

	// Replace all valSets in state to be the valSet with just one validator.
	cmtState.Validators = newValSet
	cmtState.LastValidators = newValSet
	cmtState.NextValidators = newValSet
	cmtState.LastHeightValidatorsChanged = blockStore.Height()

	err = stateStore.Save(cmtState)
	if err != nil {
		return nil, err
	}

	valSet, err := cmtState.Validators.ToProto()
	if err != nil {
		return nil, err
	}
	valInfo := &cmtstate.ValidatorsInfo{
		ValidatorSet:      valSet,
		LastHeightChanged: cmtState.LastBlockHeight,
	}
	buf, err := valInfo.Marshal()
	if err != nil {
		return nil, err
	}

	err = stateDB.Set(fmt.Appendf(nil, "validatorsKey:%v", blockStore.Height()), buf)
	if err != nil {
		return nil, err
	}

	err = stateDB.Set(fmt.Appendf(nil, "validatorsKey:%v", blockStore.Height()-1), buf)
	if err != nil {
		return nil, err
	}

	err = stateDB.Set(fmt.Appendf(nil, "validatorsKey:%v", blockStore.Height()+1), buf)
	if err != nil {
		return nil, err
	}

	// Since we modified the chainID, we set the new genesisDoc in the stateDB.
	b, err := cmtjson.Marshal(genDoc)
	if err != nil {
		return nil, err
	}
	if err := stateDB.SetSync([]byte("genesisDoc"), b); err != nil {
		return nil, err
	}

	return devnetApp, err
}

func GenDocProvider(cfg *cmtcfg.Config) func() (*cmttypes.GenesisDoc, error) {
	return func() (*cmttypes.GenesisDoc, error) {
		appGenesis, err := genutiltypes.AppGenesisFromFile(cfg.GenesisFile())
		if err != nil {
			return nil, err
		}

		return appGenesis.ToGenesisDoc()
	}
}
