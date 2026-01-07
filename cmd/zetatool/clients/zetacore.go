package clients

import (
	"context"
	"fmt"

	zetacorerpc "github.com/zeta-chain/node/pkg/rpc"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// ZetacoreClientAdapter wraps pkg/rpc.Clients to implement ZetacoreReader interface
type ZetacoreClientAdapter struct {
	clients zetacorerpc.Clients
}

// NewZetacoreClientAdapter creates a new ZetacoreClientAdapter from an RPC URL
func NewZetacoreClientAdapter(rpcURL string) (*ZetacoreClientAdapter, error) {
	clients, err := zetacorerpc.NewCometBFTClients(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create zetacore client: %w", err)
	}
	return &ZetacoreClientAdapter{clients: clients}, nil
}

// GetCctxByHash retrieves a CCTX by its hash/index
func (z *ZetacoreClientAdapter) GetCctxByHash(ctx context.Context, hash string) (*crosschaintypes.CrossChainTx, error) {
	return z.clients.GetCctxByHash(ctx, hash)
}

// InboundHashToCctxData retrieves CCTX data from an inbound hash
func (z *ZetacoreClientAdapter) InboundHashToCctxData(ctx context.Context, hash string) (*crosschaintypes.QueryInboundHashToCctxDataResponse, error) {
	req := &crosschaintypes.QueryInboundHashToCctxDataRequest{InboundHash: hash}
	return z.clients.Crosschain.InboundHashToCctxData(ctx, req)
}

// GetOutboundTracker retrieves an outbound tracker by chain ID and nonce
func (z *ZetacoreClientAdapter) GetOutboundTracker(ctx context.Context, chainID int64, nonce uint64) (*crosschaintypes.OutboundTracker, error) {
	return z.clients.GetOutboundTracker(ctx, chainID, nonce)
}

// GetChainParamsForChainID retrieves chain params for a given chain ID
func (z *ZetacoreClientAdapter) GetChainParamsForChainID(ctx context.Context, chainID int64) (*observertypes.ChainParams, error) {
	return z.clients.GetChainParamsForChainID(ctx, chainID)
}

// GetTssAddress retrieves TSS addresses (both EVM and BTC)
func (z *ZetacoreClientAdapter) GetTssAddress(ctx context.Context, btcChainID int64) (*observertypes.QueryGetTssAddressResponse, error) {
	req := &observertypes.QueryGetTssAddressRequest{BitcoinChainId: btcChainID}
	return z.clients.Observer.GetTssAddress(ctx, req)
}

// GetEVMTSSAddress retrieves the current EVM TSS address
func (z *ZetacoreClientAdapter) GetEVMTSSAddress(ctx context.Context) (string, error) {
	return z.clients.GetEVMTSSAddress(ctx)
}

// GetBTCTSSAddress retrieves the current BTC TSS address for a chain
func (z *ZetacoreClientAdapter) GetBTCTSSAddress(ctx context.Context, chainID int64) (string, error) {
	return z.clients.GetBTCTSSAddress(ctx, chainID)
}

// GetBallotByID retrieves a ballot by its identifier
func (z *ZetacoreClientAdapter) GetBallotByID(ctx context.Context, id string) (*observertypes.QueryBallotByIdentifierResponse, error) {
	return z.clients.GetBallotByID(ctx, id)
}
