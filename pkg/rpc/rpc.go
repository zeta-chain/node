package rpc

import (
	"fmt"

	rpcclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"

	etherminttypes "github.com/zeta-chain/zetacore/rpc/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// Clients contains RPC client interfaces to interact with zetacored
type Clients struct {
	// Cosmos SDK clients

	// Auth is a github.com/cosmos/cosmos-sdk/x/auth/types QueryClient
	Auth authtypes.QueryClient
	// Bank is a github.com/cosmos/cosmos-sdk/x/bank/types QueryClient
	Bank banktypes.QueryClient

	// ZetaCore specific clients

	// Authority is a github.com/zeta-chain/zetacore/x/authority/types QueryClient
	Authority authoritytypes.QueryClient
	// Crosschain is a github.com/zeta-chain/zetacore/x/crosschain/types QueryClient
	Crosschain crosschaintypes.QueryClient
	// Fungible is a github.com/zeta-chain/zetacore/x/fungible/types QueryClient
	Fungible fungibletypes.QueryClient
	// Observer is a github.com/zeta-chain/zetacore/x/observer/types QueryClient
	Observer observertypes.QueryClient
	// Light is a github.com/zeta-chain/zetacore/x/lightclient/types QueryClient
	Light lightclienttypes.QueryClient

	// Ethermint specific clients

	// Ethermint is a github.com/zeta-chain/zetacore/rpc/types QueryClient
	Ethermint *etherminttypes.QueryClient
}

func newClients(ctx client.Context) (Clients, error) {
	return Clients{
		Authority:  authoritytypes.NewQueryClient(ctx),
		Crosschain: crosschaintypes.NewQueryClient(ctx),
		Fungible:   fungibletypes.NewQueryClient(ctx),
		Auth:       authtypes.NewQueryClient(ctx),
		Bank:       banktypes.NewQueryClient(ctx),
		Observer:   observertypes.NewQueryClient(ctx),
		Light:      lightclienttypes.NewQueryClient(ctx),
		Ethermint:  etherminttypes.NewQueryClient(ctx),
	}, nil
}

// NewCometBFTClients creates a Clients which uses cometbft abci_query as the transport
func NewCometBFTClients(url string) (Clients, error) {
	cometRPCClient, err := rpcclient.New(url, "/websocket")
	if err != nil {
		return Clients{}, fmt.Errorf("create cometbft rpc client: %w", err)
	}
	clientCtx := client.Context{}.WithClient(cometRPCClient)

	return newClients(clientCtx)
}

// NewGRPCClient creates a Clients which uses gRPC as the transport
func NewGRPCClients(url string, opts ...grpc.DialOption) (Clients, error) {
	grpcConn, err := grpc.Dial(url, opts...)
	if err != nil {
		return Clients{}, err
	}
	clientCtx := client.Context{}.WithGRPCClient(grpcConn)
	return newClients(clientCtx)
}
