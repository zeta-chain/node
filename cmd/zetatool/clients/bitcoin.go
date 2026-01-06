package clients

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
	btcclient "github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	zetaclientConfig "github.com/zeta-chain/node/zetaclient/config"
)

// BitcoinClientAdapter wraps btcclient.Client to implement BitcoinClient interface
type BitcoinClientAdapter struct {
	client *btcclient.Client
}

// NewBitcoinClientAdapter creates a new BitcoinClientAdapter
func NewBitcoinClientAdapter(cfg *config.Config, chain chains.Chain, logger zerolog.Logger) (*BitcoinClientAdapter, error) {
	params, err := chains.BitcoinNetParamsFromChainID(chain.ChainId)
	if err != nil {
		return nil, fmt.Errorf("unable to get bitcoin net params: %w", err)
	}

	connCfg := zetaclientConfig.BTCConfig{
		RPCUsername: cfg.BtcUser,
		RPCPassword: cfg.BtcPassword,
		RPCHost:     cfg.BtcHost,
		RPCParams:   params.Name,
	}

	client, err := btcclient.New(connCfg, chain.ChainId, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create bitcoin client: %w", err)
	}

	return &BitcoinClientAdapter{
		client: client,
	}, nil
}

// Ping tests the connection to the Bitcoin server
func (b *BitcoinClientAdapter) Ping(ctx context.Context) error {
	return b.client.Ping(ctx)
}

// GetRawTransactionVerbose returns detailed information about a transaction
func (b *BitcoinClientAdapter) GetRawTransactionVerbose(ctx context.Context, txHash *chainhash.Hash) (*btcjson.TxRawResult, error) {
	return b.client.GetRawTransactionVerbose(ctx, txHash)
}

// GetBlockVerbose returns detailed information about a block
func (b *BitcoinClientAdapter) GetBlockVerbose(ctx context.Context, blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error) {
	return b.client.GetBlockVerbose(ctx, blockHash)
}

// GetRawClient returns the underlying btcclient.Client for advanced operations
func (b *BitcoinClientAdapter) GetRawClient() *btcclient.Client {
	return b.client
}
