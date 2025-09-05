package rpc

import (
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	rpcclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	"google.golang.org/grpc"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
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
	// Upgrade is a cosmossdk.io/x/upgrade/types QueryClient
	Upgrade upgradetypes.QueryClient
	// Distribution is a "github.com/cosmos/cosmos-sdk/x/distribution/types" QueryClient
	Distribution distributiontypes.QueryClient

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
	// EmissionsClient is a github.com/zeta-chain/zetacore/x/emissions/types QueryClient
	Emissions emissionstypes.QueryClient

	// Evm specific clients

	// Evm is a github.com/zeta-chain/zetacore/rpc/types QueryClient
	// EvmFeeMarket is a github.com/zeta-chain/evm/x/feemarket/types QueryClient
	EvmFeeMarket feemarkettypes.QueryClient

	// Tendermint specific clients

	// Tendermint is a github.com/cosmos/cosmos-sdk/client/grpc/cmtservice QueryClient
	Tendermint cmtservice.ServiceClient
}

func newClients(ctx client.Context) (Clients, error) {
	return Clients{
		// Cosmos SDK clients
		Auth:         authtypes.NewQueryClient(ctx),
		Bank:         banktypes.NewQueryClient(ctx),
		Staking:      stakingtypes.NewQueryClient(ctx),
		Upgrade:      upgradetypes.NewQueryClient(ctx),
		Authority:    authoritytypes.NewQueryClient(ctx),
		Distribution: distributiontypes.NewQueryClient(ctx),
		// ZetaCore specific clients
		Crosschain:  crosschaintypes.NewQueryClient(ctx),
		Fungible:    fungibletypes.NewQueryClient(ctx),
		Observer:    observertypes.NewQueryClient(ctx),
		Lightclient: lightclienttypes.NewQueryClient(ctx),
		Emissions:   emissionstypes.NewQueryClient(ctx),
		// Evm specific clients
		EvmFeeMarket: feemarkettypes.NewQueryClient(ctx),
		// Tendermint specific clients
		Tendermint: cmtservice.NewServiceClient(ctx),
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

// NewGRPCClients creates a Clients which uses gRPC as the transport
func NewGRPCClients(url string, opts ...grpc.DialOption) (Clients, error) {
	grpcConn, err := grpc.Dial(url, opts...)
	if err != nil {
		return Clients{}, err
	}
	clientCtx := client.Context{}.WithGRPCClient(grpcConn)
	return newClients(clientCtx)
}
