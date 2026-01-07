package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
	btcclient "github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	zetaclientConfig "github.com/zeta-chain/node/zetaclient/config"
)

const (
	mempoolAddressAPIMainnet  = "https://mempool.space/api/address/%s"
	mempoolAddressAPITestnet3 = "https://mempool.space/testnet/api/address/%s"
	mempoolAddressAPISignet   = "https://mempool.space/signet/api/address/%s"
	mempoolAddressAPITestnet4 = "https://mempool.space/testnet4/api/address/%s"
	satoshisPerBitcoin        = 100_000_000
	httpClientTimeout         = 30 * time.Second
)

// BTCAddressStats represents the response from mempool.space address API
type BTCAddressStats struct {
	Address    string `json:"address"`
	ChainStats struct {
		FundedTxoCount int   `json:"funded_txo_count"`
		FundedTxoSum   int64 `json:"funded_txo_sum"`
		SpentTxoCount  int   `json:"spent_txo_count"`
		SpentTxoSum    int64 `json:"spent_txo_sum"`
		TxCount        int   `json:"tx_count"`
	} `json:"chain_stats"`
	MempoolStats struct {
		FundedTxoCount int   `json:"funded_txo_count"`
		FundedTxoSum   int64 `json:"funded_txo_sum"`
		SpentTxoCount  int   `json:"spent_txo_count"`
		SpentTxoSum    int64 `json:"spent_txo_sum"`
		TxCount        int   `json:"tx_count"`
	} `json:"mempool_stats"`
}

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

// GetBTCBalance fetches the BTC balance for a given address using mempool.space API
// Returns the balance in BTC (not satoshis)
func GetBTCBalance(ctx context.Context, address string, chainID int64) (float64, error) {
	apiURL := getMempoolAddressAPIURL(chainID, address)
	if apiURL == "" {
		return 0, fmt.Errorf("unsupported Bitcoin chain ID: %d", chainID)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	btcClient := &http.Client{Timeout: httpClientTimeout}
	resp, err := btcClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch address stats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("mempool.space API returned status %d", resp.StatusCode)
	}

	var stats BTCAddressStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return 0, fmt.Errorf("failed to decode address stats: %w", err)
	}

	balanceSatoshis := stats.ChainStats.FundedTxoSum - stats.ChainStats.SpentTxoSum

	return float64(balanceSatoshis) / satoshisPerBitcoin, nil
}

// getMempoolAddressAPIURL returns the mempool.space address API URL for the given chain ID
func getMempoolAddressAPIURL(chainID int64, address string) string {
	switch chainID {
	case 8332: // Bitcoin mainnet
		return fmt.Sprintf(mempoolAddressAPIMainnet, address)
	case 18332: // Bitcoin testnet3
		return fmt.Sprintf(mempoolAddressAPITestnet3, address)
	case 18333: // Bitcoin signet
		return fmt.Sprintf(mempoolAddressAPISignet, address)
	case 18334: // Bitcoin testnet4
		return fmt.Sprintf(mempoolAddressAPITestnet4, address)
	default:
		return ""
	}
}

// GetBTCChainID returns the Bitcoin chain ID for the given network
func GetBTCChainID(network string) int64 {
	switch network {
	case config.NetworkMainnet:
		return 8332
	case config.NetworkTestnet:
		return 18332
	case config.NetworkLocalnet:
		return 18444
	default:
		panic("invalid network")
	}
}
