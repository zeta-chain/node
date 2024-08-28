package rpc

import (
	"fmt"

	rpcclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"

	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// Clients contains RPC client interfaces to interact with zetacored
type Clients struct {
	AuthorityClient authoritytypes.QueryClient
	CctxClient      crosschaintypes.QueryClient
	FungibleClient  fungibletypes.QueryClient
	AuthClient      authtypes.QueryClient
	BankClient      banktypes.QueryClient
	ObserverClient  observertypes.QueryClient
	LightClient     lightclienttypes.QueryClient
}

func newClients(ctx client.Context) (Clients, error) {
	return Clients{
		AuthorityClient: authoritytypes.NewQueryClient(ctx),
		CctxClient:      crosschaintypes.NewQueryClient(ctx),
		FungibleClient:  fungibletypes.NewQueryClient(ctx),
		AuthClient:      authtypes.NewQueryClient(ctx),
		BankClient:      banktypes.NewQueryClient(ctx),
		ObserverClient:  observertypes.NewQueryClient(ctx),
		LightClient:     lightclienttypes.NewQueryClient(ctx),
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
