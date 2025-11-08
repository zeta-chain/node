package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/spf13/cobra"
	"golang.org/x/net/netutil"

	tmcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"

	"github.com/zeta-chain/node/server/config"

	"cosmossdk.io/log"

	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/version"
)

// AddCommands adds server commands
func AddCommands(
	rootCmd *cobra.Command,
	opts StartOptions,
	appExport types.AppExporter,
	addStartFlags types.ModuleInitFlags,
) {
	cometbftCmd := &cobra.Command{
		Use:     "comet",
		Aliases: []string{"cometbft"},
		Short:   "CometBFT subcommands",
	}

	cometbftCmd.AddCommand(
		sdkserver.ShowNodeIDCmd(),
		sdkserver.ShowValidatorCmd(),
		sdkserver.ShowAddressCmd(),
		sdkserver.VersionCmd(),
		tmcmd.ResetAllCmd,
		tmcmd.ResetStateCmd,
		sdkserver.BootstrapStateCmd(opts.AppCreator),
	)

	startCmd := StartCmd(opts)
	addStartFlags(startCmd)

	rootCmd.AddCommand(
		startCmd,
		cometbftCmd,
		sdkserver.ExportCmd(appExport, opts.DefaultNodeHome),
		version.NewVersionCommand(),
		sdkserver.NewRollbackCmd(opts.AppCreator, opts.DefaultNodeHome),

		// custom tx indexer command
		NewIndexTxCmd(),
	)
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
