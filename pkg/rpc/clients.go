package rpc

import (
	"fmt"

	rpcclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	feemarkettypes "github.com/zeta-chain/ethermint/x/feemarket/types"
	"google.golang.org/grpc"

	etherminttypes "github.com/zeta-chain/node/rpc/types"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// Clients contains RPC client interfaces to interact with ZetaCore
//
// Clients also has some high level wrappers for the clients
type Clients struct {
	// Cosmos SDK clients

	// Auth is a github.com/cosmos/cosmos-sdk/x/auth/types QueryClient
	Auth authtypes.QueryClient
	// Bank is a github.com/cosmos/cosmos-sdk/x/bank/types QueryClient
	Bank banktypes.QueryClient
	// Bank is a github.com/cosmos/cosmos-sdk/x/staking/types QueryClient
	Staking stakingtypes.QueryClient
	// Upgrade is a github.com/cosmos/cosmos-sdk/x/upgrade/types QueryClient
	Upgrade upgradetypes.QueryClient

	// ZetaCore specific clients

	// Authority is a github.com/zeta-chain/zetacore/x/authority/types QueryClient
	Authority authoritytypes.QueryClient
	// Crosschain is a github.com/zeta-chain/zetacore/x/crosschain/types QueryClient
	Crosschain crosschaintypes.QueryClient
	// Fungible is a github.com/zeta-chain/zetacore/x/fungible/types QueryClient
	Fungible fungibletypes.QueryClient
	// Observer is a github.com/zeta-chain/zetacore/x/observer/types QueryClient
	Observer observertypes.QueryClient
	// Lightclient is a github.com/zeta-chain/zetacore/x/lightclient/types QueryClient
	Lightclient lightclienttypes.QueryClient

	// Ethermint specific clients

	// Ethermint is a github.com/zeta-chain/zetacore/rpc/types QueryClient
	Ethermint *etherminttypes.QueryClient
	// EthermintFeeMarket is a github.com/zeta-chain/ethermint/x/feemarket/types QueryClient
	EthermintFeeMarket feemarkettypes.QueryClient

	// Tendermint specific clients

	// Tendermint is a github.com/cosmos/cosmos-sdk/client/grpc/tmservice QueryClient
	Tendermint tmservice.ServiceClient
}

func newClients(ctx client.Context) (Clients, error) {
	return Clients{
		// Cosmos SDK clients
		Auth:      authtypes.NewQueryClient(ctx),
		Bank:      banktypes.NewQueryClient(ctx),
		Staking:   stakingtypes.NewQueryClient(ctx),
		Upgrade:   upgradetypes.NewQueryClient(ctx),
		Authority: authoritytypes.NewQueryClient(ctx),
		// ZetaCore specific clients
		Crosschain:  crosschaintypes.NewQueryClient(ctx),
		Fungible:    fungibletypes.NewQueryClient(ctx),
		Observer:    observertypes.NewQueryClient(ctx),
		Lightclient: lightclienttypes.NewQueryClient(ctx),
		// Ethermint specific clients
		Ethermint:          etherminttypes.NewQueryClient(ctx),
		EthermintFeeMarket: feemarkettypes.NewQueryClient(ctx),
		// Tendermint specific clients
		Tendermint: tmservice.NewServiceClient(ctx),
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
