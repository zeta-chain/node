package clients

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
)

// Factory creates RPC clients for different chains
type Factory interface {
	NewEVMClient(chain chains.Chain) (EVMClient, error)
	NewBitcoinClient(chain chains.Chain) (BitcoinClient, error)
	NewSolanaClient() (SolanaClient, error)
}

// DefaultFactory implements Factory using real RPC connections
type DefaultFactory struct {
	cfg    *config.Config
	logger zerolog.Logger
}

// NewDefaultFactory creates a new DefaultFactory
func NewDefaultFactory(cfg *config.Config, logger zerolog.Logger) *DefaultFactory {
	return &DefaultFactory{
		cfg:    cfg,
		logger: logger,
	}
}

// NewEVMClient creates a new EVM client for the given chain
func (f *DefaultFactory) NewEVMClient(chain chains.Chain) (EVMClient, error) {
	rpcURL := f.getEVMRPCURL(chain)
	if rpcURL == "" {
		return nil, fmt.Errorf("no RPC URL configured for chain %d", chain.ChainId)
	}
	return NewEVMClientAdapter(rpcURL)
}

// NewBitcoinClient creates a new Bitcoin client for the given chain
func (f *DefaultFactory) NewBitcoinClient(chain chains.Chain) (BitcoinClient, error) {
	return NewBitcoinClientAdapter(f.cfg, chain, f.logger)
}

// NewSolanaClient creates a new Solana client
func (f *DefaultFactory) NewSolanaClient() (SolanaClient, error) {
	if f.cfg.SolanaRPC == "" {
		return nil, fmt.Errorf("solana RPC URL not configured")
	}
	return NewSolanaClientAdapter(f.cfg.SolanaRPC)
}

// getEVMRPCURL returns the RPC URL for the given EVM chain
func (f *DefaultFactory) getEVMRPCURL(chain chains.Chain) string {
	switch chain.Network {
	case chains.Network_eth:
		return f.cfg.EthereumRPC
	case chains.Network_bsc:
		return f.cfg.BscRPC
	case chains.Network_polygon:
		return f.cfg.PolygonRPC
	case chains.Network_base:
		return f.cfg.BaseRPC
	case chains.Network_arbitrum:
		return f.cfg.ArbitrumRPC
	case chains.Network_optimism:
		return f.cfg.OptimismRPC
	case chains.Network_avalanche:
		return f.cfg.AvalancheRPC
	case chains.Network_worldchain:
		return f.cfg.WorldRPC
	default:
		return ""
	}
}
