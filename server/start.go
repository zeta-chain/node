package server

import (
	"context"
	"fmt"
	"time"

	pruningtypes "cosmossdk.io/store/pruning/types"
	tcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/server/types"
	cosmosevmserverconfig "github.com/cosmos/evm/server/config"
	srvflags "github.com/cosmos/evm/server/flags"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const (
	KeyIsTestnet                = "is-testnet"
	KeyNewChainID               = "new-chain-ID"
	KeyValidatorConsensusAddr   = "validator-consensus-address"
	KeyValidatorConsensusPubkey = "validator-consensus-pubkey"
	KeyOperatorAddress          = "operator-address"
	FlagSkipConfirmation        = "skip-confirmation"
)

type StartCmdOptions struct {
	DBOpener            func(rootDir string, backendType dbm.BackendType) (dbm.DB, error)
	PostSetup           func(svrCtx *server.Context, clientCtx client.Context, ctx context.Context, g *errgroup.Group) error
	PostSetupStandalone func(svrCtx *server.Context, clientCtx client.Context, ctx context.Context, g *errgroup.Group) error
	AddFlags            func(cmd *cobra.Command)
	StartCommandHandler func(svrCtx *server.Context, clientCtx client.Context, appCreator types.AppCreator, withCmt bool, opts StartCmdOptions) error
}

// StartCmd runs the service passed in, either stand-alone or in-process with
// CometBFT.
func StartCmd(appCreator types.AppCreator, defaultNodeHome string) *cobra.Command {
	return StartCmdWithOptions(appCreator, defaultNodeHome, StartCmdOptions{
		DBOpener:            openDB,
		StartCommandHandler: start,
	})
}

func StartCmdWithOptions(appCreator types.AppCreator, defaultNodeHome string, opts StartCmdOptions) *cobra.Command {
	if opts.DBOpener == nil || opts.StartCommandHandler == nil {
		panic("DBOpener and StartCommandHandler must be provided")
	}
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

			serverCtx.Logger.Info("Unlocking keyring")

			// fire unlock precess for keyring
			keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
			if keyringBackend == keyring.BackendFile {
				_, err = clientCtx.Keyring.List()
				if err != nil {
					return err
				}
			}
			// Default value is false we always use default value
			skipOverwrite := false
			if val, err := cmd.Flags().GetBool(FlagSkipConfigOverwrite); err == nil {
				skipOverwrite = val
			}

			if !skipOverwrite {
				err := overWriteConfig(cmd)
				if err != nil {
					return fmt.Errorf("failed to overwrite config: %w", err)
				}
			}

			withCmt, _ := cmd.Flags().GetBool(srvflags.WithCometBFT)

			err = opts.StartCommandHandler(serverCtx, clientCtx, appCreator, withCmt, opts)

			serverCtx.Logger.Debug("received quit signal")
			graceDuration, _ := cmd.Flags().GetDuration(server.FlagShutdownGrace)
			if graceDuration > 0 {
				serverCtx.Logger.Info("graceful shutdown start", server.FlagShutdownGrace, graceDuration)
				<-time.After(graceDuration)
				serverCtx.Logger.Info("graceful shutdown complete")
			}

			return err
		},
	}

	cmd.Flags().
		Bool(FlagSkipConfigOverwrite, false, "Skip running the config configuration overwrite handler.This is used for testing purposes only and skips using the default timeouts hardcoded and uses the config file instead")
	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().Bool(srvflags.WithCometBFT, true, "Run abci app embedded in-process with CometBFT")
	cmd.Flags().String(srvflags.Address, "tcp://0.0.0.0:26658", "Listen address")
	cmd.Flags().String(srvflags.Transport, "socket", "Transport protocol: socket, grpc")
	cmd.Flags().String(srvflags.TraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().
		String(server.FlagMinGasPrices, "", "Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum (e.g. 20000000000azeta)")
	cmd.Flags().Duration(server.FlagShutdownGrace, 3*time.Second, "On Shutdown, duration to wait for resource clean up")
	cmd.Flags().
		IntSlice(server.FlagUnsafeSkipUpgrades, []int{}, "Skip a set of upgrade heights to continue the old binary")
	cmd.Flags().
		Uint64(server.FlagHaltHeight, 0, "Block height at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().
		Uint64(server.FlagHaltTime, 0, "Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Bool(server.FlagInterBlockCache, true, "Enable inter-block caching")
	cmd.Flags().String(srvflags.CPUProfile, "", "Enable CPU profiling and write to the provided file")
	cmd.Flags().Bool(server.FlagTrace, false, "Provide full stack traces for errors in ABCI Log")
	cmd.Flags().
		String(server.FlagPruning, pruningtypes.PruningOptionDefault, "Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().
		Uint64(server.FlagPruningKeepRecent, 0, "Number of recent heights to keep on disk (ignored if pruning is not 'custom')")
	cmd.Flags().
		Uint64(server.FlagPruningInterval, 0, "Height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom')")

	cmd.Flags().Uint(server.FlagInvCheckPeriod, 0, "Assert registered invariants every N blocks")
	cmd.Flags().
		Uint64(server.FlagMinRetainBlocks, 0, "Minimum block height offset during ABCI commit to prune CometBFT blocks")
	cmd.Flags().String(srvflags.AppDBBackend, "", "The type of database for application and snapshots databases")

	cmd.Flags().Bool(srvflags.GRPCOnly, false, "Start the node in gRPC query only mode without CometBFT process")
	cmd.Flags().
		Bool(srvflags.GRPCEnable, cosmosevmserverconfig.DefaultGRPCEnable, "Define if the gRPC server should be enabled")
	cmd.Flags().String(srvflags.GRPCAddress, serverconfig.DefaultGRPCAddress, "the gRPC server address to listen on")
	cmd.Flags().
		Bool(srvflags.GRPCWebEnable, cosmosevmserverconfig.DefaultGRPCWebEnable, "Define if the gRPC-Web server should be enabled. (Note: gRPC must also be enabled.)")
	cmd.Flags().
		String(srvflags.GRPCWebAddress, cosmosevmserverconfig.DefaultGRPCAddress, "The gRPC-Web server address to listen on")

	cmd.Flags().
		Bool(srvflags.RPCEnable, cosmosevmserverconfig.DefaultAPIEnable, "Defines if Cosmos-sdk REST server should be enabled")
	cmd.Flags().
		Bool(srvflags.EnabledUnsafeCors, false, "Defines if CORS should be enabled (unsafe - use it at your own risk)")

	cmd.Flags().
		Bool(srvflags.JSONRPCEnable, cosmosevmserverconfig.DefaultJSONRPCEnable, "Define if the JSON-RPC server should be enabled")
	cmd.Flags().
		StringSlice(srvflags.JSONRPCAPI, cosmosevmserverconfig.GetDefaultAPINamespaces(), "Defines a list of JSON-RPC namespaces that should be enabled")
	cmd.Flags().
		String(srvflags.JSONRPCAddress, cosmosevmserverconfig.DefaultJSONRPCAddress, "the JSON-RPC server address to listen on")
	cmd.Flags().
		String(srvflags.JSONWsAddress, cosmosevmserverconfig.DefaultJSONRPCWsAddress, "the JSON-RPC WS server address to listen on")
	cmd.Flags().
		Uint64(srvflags.JSONRPCGasCap, cosmosevmserverconfig.DefaultGasCap, "Sets a cap on gas that can be used in eth_call/estimateGas unit is aatom (0=infinite)")

	cmd.Flags().
		Bool(srvflags.JSONRPCAllowInsecureUnlock, cosmosevmserverconfig.DefaultJSONRPCAllowInsecureUnlock, "Allow insecure account unlocking when account-related RPCs are exposed by http")
	cmd.Flags().
		Float64(srvflags.JSONRPCTxFeeCap, cosmosevmserverconfig.DefaultTxFeeCap, "Sets a cap on transaction fee that can be sent via the RPC APIs (1 = default 1 evmos)")
	cmd.Flags().
		Int32(srvflags.JSONRPCFilterCap, cosmosevmserverconfig.DefaultFilterCap, "Sets the global cap for total number of filters that can be created")
	cmd.Flags().
		Duration(srvflags.JSONRPCEVMTimeout, cosmosevmserverconfig.DefaultEVMTimeout, "Sets a timeout used for eth_call (0=infinite)")
	cmd.Flags().
		Duration(srvflags.JSONRPCHTTPTimeout, cosmosevmserverconfig.DefaultHTTPTimeout, "Sets a read/write timeout for json-rpc http server (0=infinite)")
	cmd.Flags().
		Duration(srvflags.JSONRPCHTTPIdleTimeout, cosmosevmserverconfig.DefaultHTTPIdleTimeout, "Sets a idle timeout for json-rpc http server (0=infinite)")
	cmd.Flags().
		Bool(srvflags.JSONRPCAllowUnprotectedTxs, cosmosevmserverconfig.DefaultAllowUnprotectedTxs, "Allow for unprotected (non EIP155 signed) transactions to be submitted via the node's RPC when the global parameter is disabled")

	cmd.Flags().
		Int32(srvflags.JSONRPCLogsCap, cosmosevmserverconfig.DefaultLogsCap, "Sets the max number of results can be returned from single `eth_getLogs` query")
	cmd.Flags().
		Int32(srvflags.JSONRPCBlockRangeCap, cosmosevmserverconfig.DefaultBlockRangeCap, "Sets the max block range allowed for `eth_getLogs` query")
	cmd.Flags().
		Int(srvflags.JSONRPCMaxOpenConnections, cosmosevmserverconfig.DefaultMaxOpenConnections, "Sets the maximum number of simultaneous connections for the server listener")

	cmd.Flags().Bool(srvflags.JSONRPCEnableIndexer, false, "Enable the custom tx indexer for json-rpc")
	cmd.Flags().Bool(srvflags.JSONRPCEnableMetrics, false, "Define if EVM rpc metrics server should be enabled")

	cmd.Flags().
		String(srvflags.EVMTracer, cosmosevmserverconfig.DefaultEVMTracer, "the EVM tracer type to collect execution traces from the EVM transaction execution (json|struct|access_list|markdown)")
	cmd.Flags().
		Uint64(srvflags.EVMMaxTxGasWanted, cosmosevmserverconfig.DefaultMaxTxGasWanted, "the gas wanted for each eth tx returned in ante handler in check tx mode")
	cmd.Flags().
		Bool(srvflags.EVMEnablePreimageRecording, cosmosevmserverconfig.DefaultEnablePreimageRecording, "Enables tracking of SHA3 preimages in the EVM (not implemented yet)")
	cmd.Flags().
		Uint64(srvflags.EVMChainID, cosmosevmserverconfig.DefaultEVMChainID, "the EIP-155 compatible replay protection chain ID")

	cmd.Flags().String(srvflags.TLSCertPath, "", "the cert.pem file path for the server TLS configuration")
	cmd.Flags().String(srvflags.TLSKeyPath, "", "the key.pem file path for the server TLS configuration")

	cmd.Flags().Uint64(server.FlagStateSyncSnapshotInterval, 0, "State sync snapshot interval")
	cmd.Flags().Uint32(server.FlagStateSyncSnapshotKeepRecent, 2, "State sync snapshot to keep")

	cmd.Flags().Bool(KeyIsTestnet, false, "Enable testnet mode to fork from existing state")
	cmd.Flags().String(KeyNewChainID, "", "New chain ID to use when running in testnet mode")

	if opts.AddFlags != nil {
		opts.AddFlags(cmd)
	}

	// add support for all CometBFT-specific command line options
	tcmd.AddNodeFlags(cmd)
	return cmd
}
